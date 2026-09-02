package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"scheduler/internal/repository/attendance"
	"scheduler/internal/services/report"
)

// Kunlik reja (barcha vaqtlar TIMEZONE sohasida):
//
//	00:00 — yangi kun boshlanadi (kunlik davomat sikli)
//	16:00 — kunlik davomat hisoboti (davomat.xlsx) Telegram'ga yuboriladi
const (
	DayStartHour = 0
	ReportHour   = 16
)

// ReportState oxirgi yuborilgan hisobot sanasini saqlaydi o'qiydi.
// (Takroriy yuborishdan himoya uchun — qarang RunReportJob.)
type ReportState interface {
	GetLastReportDate() time.Time
	SetLastReportDate(date time.Time) error
}

// DocumentSender xotiradagi fayl ma'lumotini hujjat sifatida chatga
// yuboradi.
type DocumentSender interface {
	SendDocument(chatID int64, fileName string, data []byte, caption string) error
}

// ListAttendanceFunc berilgan sanadagi davomat qatorlarini qaytaradi.
type ListAttendanceFunc func(date time.Time) ([]attendance.ReportRow, error)

// Scheduler — 16:00 (hisobot) va 00:00 (kun o'tishi) voqealarini soat
// sohasi asosida rejalashtirib turadigan fon ishchi jarayoni.
type Scheduler struct {
	loc      *time.Location
	chatID   int64
	sender   DocumentSender
	listRows ListAttendanceFunc
	state    ReportState
	nowFn    func() time.Time
}

// New scheduler yaratadi. nowFn default holda time.Now — testlarda
// voqeani to'g'ridan-to'g'ri RunReportJob/RunDayTransition orqali chaqirish
// orqali vaqt kutmasdan sinash mumkin.
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

// NextEventTime `now` dan keyingi (yoki aynan shu zahotdagi) rejalashtirilgan
// voqeani qaytaradi: (vaqt, hisobotmi). Qaytariladigan vaqtlar har doim
// soat sohasidagi devor soati (wall-clock) asosida hisoblanadi — DST
// sohalarida ham "soat 16:00" ma'nosini saqlaydi.
func NextEventTime(now time.Time, loc *time.Location) (time.Time, bool) {
	now = now.In(loc)
	base := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	for d := 0; d < 3; d++ {
		day := base.AddDate(0, 0, d)
		midnight := time.Date(day.Year(), day.Month(), day.Day(), DayStartHour, 0, 0, 0, loc)
		report := time.Date(day.Year(), day.Month(), day.Day(), ReportHour, 0, 0, 0, loc)
		if !midnight.Before(now) {
			return midnight, false
		}
		if !report.Before(now) {
			return report, true
		}
	}
	// Har kunda 2 ta voqea bo'lgani uchun 3 kun ichida albatta topiladi.
	panic("NextEventTime: voqea topilmadi")
}

// Run scheduler asosiy sikli: keyingi rejalashgan vaqtni hisoblaydi, shu
// vaqtni kutadi va voqeani bajaradi. Xatolar jarayonni to'xtatmaydi —
// yozib, keyingi voqea kutiladi. ctx bekor qilinganda jarayon to'xtaydi.
func (s *Scheduler) Run(ctx context.Context) error {
	for {
		now := s.nowFn()
		next, isReport := NextEventTime(now, s.loc)

		what := "day transition"
		if isReport {
			what = "report"
		}
		log.Printf("next %s: %s", what, next.Format("2006-01-02 15:04:05 -07 (MST)"))

		timer := time.NewTimer(next.Sub(now))
		select {
		case <-ctx.Done():
			timer.Stop()
			log.Println("scheduler stopped")
			return nil
		case <-timer.C:
		}

		// Voqea bajarildi — keyingi siklda xatolar keyingi vaqtga
		// o'tkaziladi (jarayon to'xtamaydi).
		if isReport {
			if err := s.RunReportJob(s.nowFn()); err != nil {
				log.Printf("error: kunlik hisobot ishi bajarilmadi: %v", err)
			}
		} else {
			s.RunDayTransition(s.nowFn())
		}
	}
}

// RunReportJob bitta kunlik hisobot ishini bajaradi:
//
//  1. bugungi (soat sohasidagi) sanani aniqlaydi;
//  2. restart himoyasi — shu sana uchun hisobot allaqachon yuborilgan
//     bo'lsa, takrorlamaydi;
//  3. bazadan bugungi davomat qatorlarini o'qiydi;
//  4. davomat.xlsx generatsiya qiladi;
//  5. faylni Telegram guruhiga yuboradi;
//  6. muvaffaqiyatli yuborilgan sanani bazaga saqlaydi.
//
// Hisobotdan keyin davomat ma'lumotlari o'chirilmaydi — tarixiy yozuvlar
// saqlanib qoladi.
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
		// Hisobot yuborilgan — takrorlash xavfli, shuning uchun faqat
		// ogohlantiramiz.
		log.Printf("warning: oxirgi hisobot sanasi bazaga saqlanamadi: %v", err)
	}
	return nil
}

// RunDayTransition 00:00 voqeasi. Mevcud bazada davomat yozuvlari
// (o'quvchi, sana) unikal bo'lgani va yozuvlar o'qituvchi/bot tomonidan
// zarurat bo'lganda (kiritilganda) yaratilgani uchun, yangi kunga o'tishda
// hech qanday o'chirish yoki yaratish talab etilmaydi — yangi kunning
// davomati o'sha kunning sanasi bilan yangi yozuvlar sifatida keladi.
// Shu voqea shunchaki siklni belgilaydi va log'ga yozadi.
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
