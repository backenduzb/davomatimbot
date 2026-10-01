package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"scheduler/internal/repository/attendance"
	"scheduler/internal/repository/classes"
	"scheduler/internal/services/report"
)

const (
	DayStartHour   = 0
	ReminderHour   = 9
	ReminderMinute = 45
	ReportHour     = 16
)

type EventKind string

const (
	EventDayTransition EventKind = "day_transition"
	EventReminder      EventKind = "reminder"
	EventReport        EventKind = "report"
)

type ReportState interface {
	GetLastReportDate() time.Time
	SetLastReportDate(date time.Time) error
}

type ReminderState interface {
	GetLastReminderDate() time.Time
	SetLastReminderDate(date time.Time) error
}

type DocumentSender interface {
	SendDocument(chatID int64, fileName string, data []byte, caption string) error
}

type MessageSender interface {
	SendMessage(chatID int64, text string) error
}

type ListAttendanceFunc func(date time.Time) ([]attendance.ReportRow, error)

type ListUnsubmittedFunc func(date time.Time) ([]classes.UnsubmittedClass, error)

type Scheduler struct {
	loc      *time.Location
	chatID   int64
	sender   DocumentSender
	listRows ListAttendanceFunc
	state    ReportState
	nowFn    func() time.Time

	reminderSender MessageSender
	listUnsubmitted ListUnsubmittedFunc
	reminderState   ReminderState
}

func New(loc *time.Location, chatID int64, sender DocumentSender, listRows ListAttendanceFunc, state ReportState) *Scheduler {
	return &Scheduler{
		loc:      loc,
		chatID:   chatID,
		sender:   sender,
		listRows: listRows,
		state:    state,
		nowFn:    time.Now,
	}
}

func (s *Scheduler) WithReminder(sender MessageSender, list ListUnsubmittedFunc, state ReminderState) *Scheduler {
	s.reminderSender = sender
	s.listUnsubmitted = list
	s.reminderState = state
	return s
}

func isWeekday(t time.Time) bool {
	wd := t.Weekday()
	return wd >= time.Monday && wd <= time.Friday
}

func NextEventTime(now time.Time, loc *time.Location) (time.Time, EventKind) {
	now = now.In(loc)
	base := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	for d := 0; d < 8; d++ {
		day := base.AddDate(0, 0, d)
		midnight := time.Date(day.Year(), day.Month(), day.Day(), DayStartHour, 0, 0, 0, loc)
		if !midnight.Before(now) {
			return midnight, EventDayTransition
		}
		if isWeekday(day) {
			reminder := time.Date(day.Year(), day.Month(), day.Day(), ReminderHour, ReminderMinute, 0, 0, loc)
			if !reminder.Before(now) {
				return reminder, EventReminder
			}
		}
		report := time.Date(day.Year(), day.Month(), day.Day(), ReportHour, 0, 0, 0, loc)
		if !report.Before(now) {
			return report, EventReport
		}
	}
	panic("NextEventTime: voqea topilmadi")
}

func (s *Scheduler) Run(ctx context.Context) error {
	for {
		now := s.nowFn()
		next, kind := NextEventTime(now, s.loc)

		log.Printf("next %s: %s", string(kind), next.Format("2006-01-02 15:04:05 -07 (MST)"))

		timer := time.NewTimer(next.Sub(now))
		select {
		case <-ctx.Done():
			timer.Stop()
			log.Println("scheduler stopped")
			return nil
		case <-timer.C:
		}

		switch kind {
		case EventReport:
			if err := s.RunReportJob(s.nowFn()); err != nil {
				log.Printf("error: kunlik hisobot ishi bajarilmadi: %v", err)
			}
		case EventReminder:
			if err := s.RunReminderJob(s.nowFn()); err != nil {
				log.Printf("error: eslatma ishi bajarilmadi: %v", err)
			}
		default:
			s.RunDayTransition(s.nowFn())
		}
	}
}

func (s *Scheduler) RunReportJob(now time.Time) error {
	date := midnight(now.In(s.loc))

	if last := s.state.GetLastReportDate(); !last.IsZero() && sameDay(last, date) {
		log.Printf("report for %s was already sent — skipping (restart himoyasi)", date.Format("2006-01-02"))
		return nil
	}

	log.Println("generating attendance report")

	rows, err := s.listRows(date)
	if err != nil {
		return fmt.Errorf("davomat o'qishda xato: %w", err)
	}
	log.Printf("attendance rows loaded: %d", len(rows))

	data, err := report.GenerateXLSX(rows)
	if err != nil {
		return fmt.Errorf("xlsx generatsiya qilinda xato: %w", err)
	}
	log.Println("report generated successfully")

	log.Println("sending report to Telegram")
	if err := s.sender.SendDocument(s.chatID, report.FileName, data, report.CaptionFor(date)); err != nil {
		return err
	}
	log.Println("report sent successfully")

	if err := s.state.SetLastReportDate(date); err != nil {
		log.Printf("warning: oxirgi hisobot sanasi bazaga saqlanamadi: %v", err)
	}
	return nil
}

func ReminderMessageFor(className, teacherName string, date time.Time) string {
	dateStr := date.Format("02.01.2006")
	if teacherName != "" && teacherName != "-" {
		return fmt.Sprintf("⏰ Hurmatli %s!\n\nIltimos, bugungi (<b>%s</b>) <b>%s</b> sinfi davomatini topshiring.\n\nBotda 👇 \"📋 Davomat topshirish\" tugmasini bosing yoki /start yuboring.", teacherName, dateStr, className)
	}
	return fmt.Sprintf("⏰ Iltimos, bugungi (<b>%s</b>) <b>%s</b> sinfi davomatini topshiring.\n\nBotda 👇 \"📋 Davomat topshirish\" tugmasini bosing yoki /start yuboring.", dateStr, className)
}

func (s *Scheduler) RunReminderJob(now time.Time) error {
	today := now.In(s.loc)
	if !isWeekday(today) {
		log.Printf("reminder: %s dam olish kuni — eslatma yuborilmaydi", today.Format("2006-01-02 Mon"))
		return nil
	}
	date := midnight(today)

	if s.reminderSender == nil || s.listUnsubmitted == nil || s.reminderState == nil {
		log.Println("reminder: bog'liqliklar ulanmagan — eslatma o'tkazib yuborildi")
		return nil
	}

	if last := s.reminderState.GetLastReminderDate(); !last.IsZero() && sameDay(last, date) {
		log.Printf("reminder for %s was already sent — skipping (restart himoyasi)", date.Format("2006-01-02"))
		return nil
	}

	pending, err := s.listUnsubmitted(date)
	if err != nil {
		return fmt.Errorf("topshirilmagan sinflarni o'qishda xato: %w", err)
	}
	log.Printf("reminder: topshirilmagan sinflar: %d", len(pending))

	sent := 0
	for _, c := range pending {
		text := ReminderMessageFor(c.ClassName, c.TeacherFullName, date)
		if err := s.reminderSender.SendMessage(c.TeacherTelegramID, text); err != nil {
			log.Printf("reminder: %s (%d) ga yuborilmadi: %v", c.ClassName, c.TeacherTelegramID, err)
			continue
		}
		sent++
	}
	log.Printf("reminder sent: %d/%d", sent, len(pending))

	if err := s.reminderState.SetLastReminderDate(date); err != nil {
		log.Printf("warning: oxirgi eslatma sanasi bazaga saqlanamadi: %v", err)
	}
	return nil
}

func (s *Scheduler) RunDayTransition(now time.Time) {
	today := now.In(s.loc)
	log.Printf("new attendance day started: %s", today.Format("2006-01-02"))
	log.Println("day transition: davomat yozuvlari sana asosida saqlanadi, tarixiy ma'lumotlar o'chirilmaydi")
}

func midnight(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}
