package scheduler

import (
	"testing"
	"time"

	"scheduler/internal/repository/attendance"
)

var tashkent, _ = time.LoadLocation("Asia/Tashkent")

func TestNextEventTime(t *testing.T) {
	tests := []struct {
		name       string
		now        time.Time
		want       time.Time
		wantReport bool
	}{
		{
			name: "ertalabdan 16:00 oldin — bugungi hisobot",
			now:  time.Date(2026, 9, 2, 9, 0, 0, 0, tashkent),
			want: time.Date(2026, 9, 2, 16, 0, 0, 0, tashkent),
		},
		{
			name: "16:00 da bir soniya oldin",
			now:  time.Date(2026, 9, 2, 15, 59, 0, 0, tashkent),
			want: time.Date(2026, 9, 2, 16, 0, 0, 0, tashkent),
		},
		{
			name:       "aynan 16:00 — hisobot o'z zahotida ishga tushadi",
			now:        time.Date(2026, 9, 2, 16, 0, 0, 0, tashkent),
			want:       time.Date(2026, 9, 2, 16, 0, 0, 0, tashkent),
			wantReport: true,
		},
		{
			name: "16:00 o'tib ketgan — kechasi 00:00",
			now:  time.Date(2026, 9, 2, 16, 1, 0, 0, tashkent),
			want: time.Date(2026, 9, 3, 0, 0, 0, 0, tashkent),
		},
		{
			name: "tungi soatlar — ertangi 00:00",
			now:  time.Date(2026, 9, 2, 23, 59, 0, 0, tashkent),
			want: time.Date(2026, 9, 3, 0, 0, 0, 0, tashkent),
		},
		{
			name: "aynan 00:00 — kun o'tishi o'z zahotida ishga tushadi",
			now:  time.Date(2026, 9, 3, 0, 0, 0, 0, tashkent),
			want: time.Date(2026, 9, 3, 0, 0, 0, 0, tashkent),
		},
		{
			name:       "kun boshlanib ketgan — bugungi 16:00",
			now:        time.Date(2026, 9, 3, 0, 0, 1, 0, tashkent),
			want:       time.Date(2026, 9, 3, 16, 0, 0, 0, tashkent),
			wantReport: true,
		},
		{
			name:       "UTC vaqt bilan chaqirilsa ham Tashkent devor soati bilan hisoblanadi",
			now:        time.Date(2026, 9, 2, 11, 30, 0, 0, time.UTC).In(tashkent), // = 16:30 +05
			want:       time.Date(2026, 9, 3, 0, 0, 0, 0, tashkent),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotReport := NextEventTime(tt.now, tashkent)
			if !got.Equal(tt.want) {
				t.Errorf("NextEventTime() = %v, want %v", got, tt.want)
			}
			if gotReport != tt.wantReport {
				t.Errorf("NextEventTime() isReport = %v, want %v", gotReport, tt.wantReport)
			}
		})
	}
}

func TestNextEventTimeDST(t *testing.T) {
	// DST sohasi (America/New_York): 2027-yil 13-mart — soatlar oldinga
	// suriladi. 16:00 voqea devor soatiga rioya qilishi kerak.
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skipf("timezone yo'q: %v", err)
	}

	now := time.Date(2027, 3, 13, 9, 0, 0, 0, loc)
	got, isReport := NextEventTime(now, loc)
	want := time.Date(2027, 3, 13, 16, 0, 0, 0, loc)
	if !got.Equal(want) || !isReport {
		t.Errorf("NextEventTime() = %v (report=%v), want %v (report=true)", got, isReport, want)
	}
}

type fakeState struct {
	last    time.Time
	setErr  error
	setCalls int
}

func (f *fakeState) GetLastReportDate() time.Time { return f.last }
func (f *fakeState) SetLastReportDate(date time.Time) error {
	f.setCalls++
	f.last = date
	return f.setErr
}

type fakeSender struct {
	calls       int
	lastChat    int64
	lastFile    string
	lastData    []byte
	lastCaption string
	err         error
}

func (f *fakeSender) SendDocument(chatID int64, fileName string, data []byte, caption string) error {
	f.calls++
	f.lastChat = chatID
	f.lastFile = fileName
	f.lastData = data
	f.lastCaption = caption
	return f.err
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

	// 16:00 — birinchi hisobot.
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

	// 16:01 da restart — takroriy hisobot bo'lmasligi kerak.
	if err := s.RunReportJob(time.Date(2026, 9, 2, 16, 1, 0, 0, tashkent)); err != nil {
		t.Fatalf("RunReportJob (restart): %v", err)
	}
	if sender.calls != 1 {
		t.Fatalf("restart dan keyin yuborishlar soni = %d, want 1 (takrorlanmasligi kerak)", sender.calls)
	}

	// Keyingi kun — yangi hisobot yuboriladi.
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

	// Xato tuzatildi — keyingi urinishda hisobot yuboriladi va sana saqlanadi.
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
	// Panika bo'lmasligi va xato qaytarmasligi uchun chaqiriladi.
	s.RunDayTransition(time.Date(2026, 9, 3, 0, 0, 0, 0, tashkent))
}

var (
	errSendFailed = &testError{"telegram xato"}
	errDB         = &testError{"db xato"}
)

type testError struct{ msg string }

func (e *testError) Error() string { return e.msg }
