package attendance

type Record struct {
	ID         uint   `json:"id"`
	StudentID  uint   `json:"student_id"`
	StudentName string `json:"student_name"`
	ClassID    uint   `json:"class_id"`
	Date       string `json:"date"`
	Status     string `json:"status"`
	Reason     string `json:"reason,omitempty"`
}

type TodayResponse struct {
	Date    string   `json:"date"`
	ClassID uint     `json:"class_id,omitempty"`
	Records []Record `json:"records"`
}

type ClassStats struct {
	ClassID           uint    `json:"class_id"`
	ClassName         string  `json:"class_name"`
	TotalStudents     int64   `json:"total_students"`
	Present           int64   `json:"present"`
	Absent            int64   `json:"absent"`
	Excused           int64   `json:"excused"`
	Late              int64   `json:"late"`
	NotMarked         int64   `json:"not_marked"`
	AttendancePercent float64 `json:"attendance_percent"`
}

type TodayStatistics struct {
	Date               string       `json:"date"`
	TotalClasses       int64        `json:"total_classes"`
	TotalStudents      int64        `json:"total_students"`
	Present            int64        `json:"present"`
	Absent             int64        `json:"absent"`
	Excused            int64        `json:"excused"`
	Late               int64        `json:"late"`
	NotMarked          int64        `json:"not_marked"`
	AttendancePercent  float64      `json:"attendance_percent"`
	Classes            []ClassStats `json:"classes"`
}
