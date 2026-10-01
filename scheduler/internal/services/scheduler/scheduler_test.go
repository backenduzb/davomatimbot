package scheduler

import (
	"testing"
	"time"

	"scheduler/internal/repository/attendance"
	"scheduler/internal/repository/classes"
)

var tashkent, _ = time.LoadLocation("Asia/Tashkent")

func TestNextEventTime(t *testing.T) {
	tests := []struct {
		name      string
		now       time.Time
		want      time.Time
		wantKind  EventKind
	}{
		{
			name:     "ertalab 09:45 dan oldin — bugungi eslatma",
			now:      time.Date(2026, 9, 2, 9, 0, 0, 0, tashkent),
			want:     time.Date(2026, 9, 2, 9, 45, 0, 0, tashkent),
			wantKind: EventReminder,
		},
		{
			name:     "aynan 09:45 — eslatma o'z zahotida ishga tushadi",
			now:      time.Date(2026, 9, 2, 9, 45, 0, 0, tashkent),
			want:     time.Date(2026, 9, 2, 9, 45, 0, 0, tashkent),
			wantKind: EventReminder,
		},
		{
			name:     "09:45 o'tib ketgan — bugungi hisobot",
			now:      time.Date(2026, 9, 2, 9, 46, 0, 0, tashkent),
			want:     time.Date(2026, 9, 2, 16, 0, 0, 0, tashkent),
			wantKind: EventReport,
		},
		{
			name:     "16:00 da bir soniya oldin",
			now:      time.Date(2026, 9, 2, 15, 59, 0, 0, tashkent),
			want:     time.Date(2026, 9, 2, 16, 0, 0, 0, tashkent),
			wantKind: EventReport,
		},
		{
			name:     "aynan 16:00 — hisobot o'z zahotida ishga tushadi",
			now:      time.Date(2026, 9, 2, 16, 0, 0, 0, tashkent),
			want:     time.Date(2026, 9, 2, 16, 0, 0, 0, tashkent),
			wantKind: EventReport,
		},
		{
			name:     "16:00 o'tib ketgan — kechasi 00:00",
			now:      time.Date(2026, 9, 2, 16, 1, 0, 0, tashkent),
			want:     time.Date(2026, 9, 3, 0, 0, 0, 0, tashkent),
			wantKind: EventDayTransition,
		},
		{
			name:     "tungi soatlar — ertangi 00:00",
			now:      time.Date(2026, 9, 2, 23, 59, 0, 0, tashkent),
			want:     time.Date(2026, 9, 3, 0, 0, 0, 0, tashkent),
			wantKind: EventDayTransition,
		},
		{
			name:     "aynan 00:00 — kun o'tishi o'z zahotida ishga tushadi",
			now:      time.Date(2026, 9, 3, 0, 0, 0, 0, tashkent),
			want:     time.Date(2026, 9, 3, 0, 0, 0, 0, tashkent),
			wantKind: EventDayTransition,
		},
		{
			name:     "kun boshlanib ketgan — bugungi eslatma",
			now:      time.Date(2026, 9, 3, 0, 0, 1, 0, tashkent),
			want:     time.Date(2026, 9, 3, 9, 45, 0, 0, tashkent),
			wantKind: EventReminder,
		},
		{
			name:     "juma 16:01 — shanba eslatmasiz o'tib dushanba emas, shanba 00:00",
			now:      time.Date(2026, 9, 4, 16, 1, 0, 0, tashkent),
			want:     time.Date(2026, 9, 5, 0, 0, 0, 0, tashkent), 
			wantKind: EventDayTransition,
		},
		{
			name:     "shanba kuni 09:00 — eslatma yo'q, bugungi hisobot",
			now:      time.Date(2026, 9, 5, 9, 0, 0, 0, tashkent), 
			want:     time.Date(2026, 9, 5, 16, 0, 0, 0, tashkent),
			wantKind: EventReport,
		},
		{
			name:     "yakshanba 16:01 — dushanba 00:00",
			now:      time.Date(2026, 9, 6, 16, 1, 0, 0, tashkent),
			want:     time.Date(2026, 9, 7, 0, 0, 0, 0, tashkent),
			wantKind: EventDayTransition,
		},
		{
			name:     "dushanba ertalab — eslatma",
			now:      time.Date(2026, 9, 7, 8, 0, 0, 0, tashkent),
			want:     time.Date(2026, 9, 7, 9, 45, 0, 0, tashkent),
			wantKind: EventReminder,
		},
		{
			name:     "UTC vaqt bilan chaqirilsa ham Tashkent devor soati bilan hisoblanadi",
			now:      time.Date(2026, 9, 2, 11, 30, 0, 0, time.UTC).In(tashkent), // = 16:30 +05
			want:     time.Date(2026, 9, 3, 0, 0, 0, 0, tashkent),
			wantKind: EventDayTransition,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotKind := NextEventTime(tt.now, tashkent)
			if !got.Equal(tt.want) {
				t.Errorf("NextEventTime() = %v, want %v", got, tt.want)
			}
			if gotKind != tt.wantKind {
				t.Errorf("NextEventTime() kind = %v, want %v", gotKind, tt.wantKind)
			}
		})
	}
}

func TestNextEventTimeDST(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skipf("timezone yo'q: %v", err)
	}

	now := time.Date(2027, 3, 13, 9, 0, 0, 0, loc)
	got, kind := NextEventTime(now, loc)
	want := time.Date(2027, 3, 13, 16, 0, 0, 0, loc)
	if !got.Equal(want) || kind != EventReport {
		t.Errorf("NextEventTime() = %v (kind=%v), want %v (kind=report)", got, kind, want)
	}
}

type fakeState struct {
	last    time.Time
	setErr  error
	setCalls int
	reminderLast time.Time
}

func (f *fakeState) GetLastReportDate() time.Time { return f.last }
func (f *fakeState) SetLastReportDate(date time.Time) error {
	f.setCalls++
	f.last = date
	return f.setErr
}

func (f *fakeState) GetLastReminderDate() time.Time { return f.reminderLast }
func (f *fakeState) SetLastReminderDate(date time.Time) error {
	f.reminderLast = date
	return nil
}

type fakeSender struct {
	calls       int
	lastChat    int64
	lastFile    string
	lastData    []byte
	lastCaption string
	err         error
	msgCalls    int
	msgChats    []int64
	msgTexts    []string
	msgErr      error
}

func (f *fakeSender) SendDocument(chatID int64, fileName string, data []byte, caption string) error {
	f.calls++
	f.lastChat = chatID
	f.lastFile = fileName
	f.lastData = data
	f.lastCaption = caption
	return f.err
}

func (f *fakeSender) SendMessage(chatID int64, text string) error {
	f.msgCalls++
	f.msgChats = append(f.msgChats, chatID)
	f.msgTexts = append(f.msgTexts, text)
	return f.msgErr
}

func sampleRows(date time.Time) []attendance.ReportRow {
	return []attendance.ReportRow{
		{ClassName: "10-A", Student: "Karimova Dilorom", Date: date, Status: "absent", Reason: "Illness"},
		{ClassName: "10-A", Student: "Aliyev Vali", Date: date, Status: "present"},
	}
}

func TestRunReportJobSendsOnceAndSkipsDuplicate(t *testing.T) {
	const chatID = -1003013595617
	date := time.Date(2026, 9, 2, 16, 0, 0, 0, tashkent)

	state := &fakeState{}
	sender := &fakeSender{}
	s := New(tashkent, chatID, sender,
		func(d time.Time) ([]attendance.ReportRow, error) { return sampleRows(d), nil },
		state)

	if err := s.RunReportJob(date); err != nil {
		t.Fatalf("RunReportJob: %v", err)
	}
	if sender.calls != 1 {
		t.Fatalf("yuborishlar soni = %d, want 1", sender.calls)
	}
	if sender.lastChat != chatID {
		t.Errorf("chatID = %d, want %d", sender.lastChat, chatID)
	}
	if sender.lastFile != "davomat.xlsx" {
		t.Errorf("fayl nomi = %q, want davomat.xlsx", sender.lastFile)
	}
	wantCaption := "davomat.xlsx\n\n📅 Sana: 2026/09/02\n#2026_09_02"
	if sender.lastCaption != wantCaption {
		t.Errorf("caption = %q, want %q", sender.lastCaption, wantCaption)
	}
	if state.last.IsZero() || state.last.Format("2006-01-02") != "2026-09-02" {
		t.Errorf("oxirgi sana saqlanmagan: %v", state.last)
	}

	if err := s.RunReportJob(time.Date(2026, 9, 2, 16, 1, 0, 0, tashkent)); err != nil {
		t.Fatalf("RunReportJob (restart): %v", err)
	}
	if sender.calls != 1 {
		t.Fatalf("restart dan keyin yuborishlar soni = %d, want 1 (takrorlanmasligi kerak)", sender.calls)
	}

	if err := s.RunReportJob(time.Date(2026, 9, 3, 16, 0, 0, 0, tashkent)); err != nil {
		t.Fatalf("RunReportJob (ertasi): %v", err)
	}
	if sender.calls != 2 {
		t.Fatalf("ertasi kun yuborishlar soni = %d, want 2", sender.calls)
	}
}

func TestRunReportJobSenderErrorKeepsStateUnset(t *testing.T) {
	date := time.Date(2026, 9, 2, 16, 0, 0, 0, tashkent)
	state := &fakeState{}
	sender := &fakeSender{err: errSendFailed}
	s := New(tashkent, 1, sender,
		func(d time.Time) ([]attendance.ReportRow, error) { return sampleRows(d), nil },
		state)

	if err := s.RunReportJob(date); err == nil {
		t.Fatal("RunReportJob xatoni qaytarmadi")
	}
	if !state.last.IsZero() {
		t.Error("muvaffaqiyatsiz yuborishdan keyin oxirgi sana saqlanmasligi kerak edi")
	}

	sender.err = nil
	if err := s.RunReportJob(time.Date(2026, 9, 2, 16, 5, 0, 0, tashkent)); err != nil {
		t.Fatalf("RunReportJob (qayta urinish): %v", err)
	}
	if sender.calls != 2 || state.last.IsZero() {
		t.Fatalf("qayta urinish ishlamadi: calls=%d last=%v", sender.calls, state.last)
	}
}

func TestRunReportJobListError(t *testing.T) {
	date := time.Date(2026, 9, 2, 16, 0, 0, 0, tashkent)
	state := &fakeState{}
	sender := &fakeSender{}
	s := New(tashkent, 1, sender,
		func(d time.Time) ([]attendance.ReportRow, error) { return nil, errDB },
		state)

	if err := s.RunReportJob(date); err == nil {
		t.Fatal("DB xatosida RunReportJob xatoni qaytarmadi")
	}
	if sender.calls != 0 {
		t.Error("DB xatosida hujjat yuborilmagani kerak edi")
	}
}

func TestRunDayTransition(t *testing.T) {
	s := New(tashkent, 1, &fakeSender{},
		func(d time.Time) ([]attendance.ReportRow, error) { return nil, nil },
		&fakeState{})
	s.RunDayTransition(time.Date(2026, 9, 3, 0, 0, 0, 0, tashkent))
}

func TestRunReminderJobSendsToUnsubmitted(t *testing.T) {
	
	now := time.Date(2026, 9, 2, 9, 45, 0, 0, tashkent)
	state := &fakeState{}
	sender := &fakeSender{}
	pending := []classes.UnsubmittedClass{
		{ClassID: 1, ClassName: "5-A", TeacherFullName: "Aliyeva Nodira", TeacherTelegramID: 111},
		{ClassID: 2, ClassName: "5-B", TeacherFullName: "Karimov Botir", TeacherTelegramID: 222},
	}
	s := New(tashkent, 1, sender,
		func(d time.Time) ([]attendance.ReportRow, error) { return nil, nil },
		state,
	).WithReminder(sender,
		func(d time.Time) ([]classes.UnsubmittedClass, error) { return pending, nil },
		state,
	)

	if err := s.RunReminderJob(now); err != nil {
		t.Fatalf("RunReminderJob: %v", err)
	}
	if sender.msgCalls != 2 {
		t.Fatalf("eslatmalar soni = %d, want 2", sender.msgCalls)
	}
	if sender.msgChats[0] != 111 || sender.msgChats[1] != 222 {
		t.Errorf("chatlar = %v, want [111 222]", sender.msgChats)
	}
	if state.reminderLast.Format("2006-01-02") != "2026-09-02" {
		t.Errorf("eslatma sanasi saqlanmagan: %v", state.reminderLast)
	}

	if err := s.RunReminderJob(time.Date(2026, 9, 2, 9, 50, 0, 0, tashkent)); err != nil {
		t.Fatalf("RunReminderJob (restart): %v", err)
	}
	if sender.msgCalls != 2 {
		t.Fatalf("restart dan keyin eslatmalar = %d, want 2", sender.msgCalls)
	}
}

func TestRunReminderJobSkipsWeekend(t *testing.T) {
	now := time.Date(2026, 9, 5, 9, 45, 0, 0, tashkent)
	state := &fakeState{}
	sender := &fakeSender{}
	s := New(tashkent, 1, sender,
		func(d time.Time) ([]attendance.ReportRow, error) { return nil, nil },
		state,
	).WithReminder(sender,
		func(d time.Time) ([]classes.UnsubmittedClass, error) {
			t.Error("dam olish kunida ro'yxat so'ralmasligi kerak edi")
			return nil, nil
		},
		state,
	)
	if err := s.RunReminderJob(now); err != nil {
		t.Fatalf("RunReminderJob (shanba): %v", err)
	}
	if sender.msgCalls != 0 {
		t.Errorf("shanba eslatma yuborilmasligi kerak edi, yuborildi: %d", sender.msgCalls)
	}
}

func TestRunReminderJobListError(t *testing.T) {
	now := time.Date(2026, 9, 2, 9, 45, 0, 0, tashkent)
	state := &fakeState{}
	sender := &fakeSender{}
	s := New(tashkent, 1, sender,
		func(d time.Time) ([]attendance.ReportRow, error) { return nil, nil },
		state,
	).WithReminder(sender,
		func(d time.Time) ([]classes.UnsubmittedClass, error) { return nil, errDB },
		state,
	)
	if err := s.RunReminderJob(now); err == nil {
		t.Fatal("DB xatosida RunReminderJob xatoni qaytarmadi")
	}
	if sender.msgCalls != 0 {
		t.Error("DB xatosida eslatma yuborilmasligi kerak edi")
	}
}

var (
	errSendFailed = &testError{"telegram xato"}
	errDB         = &testError{"db xato"}
)

type testError struct{ msg string }

func (e *testError) Error() string { return e.msg }
