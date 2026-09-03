package sessions

import (
	"fmt"
	"strings"

	"bot/internal/repository/classes"
)

var currentSession = make(map[uint]*AttendanceSession)

type AttendanceSession struct {
	UnexcusedStudents       map[uint]string
	ExcusedStudents         map[uint]ExcusedDetail
	// LateStudents — kech kelgan o'quvchilar (ID -> F.I.Sh).
	LateStudents map[uint]string
	LastSelectedStudentID   uint
	LastSelectedStudentName string
	ClassInfo               map[uint]ClassDetail
	ClassID                 uint
	IsAdmin                 bool
	// PendingClassName — bir xil nomli sinflar orasidan o'qituvchi bo'yicha
	// tanlash bosqichida saqlanadigan sinf nomi.
	PendingClassName string
}

type ClassDetail struct {
	TeacherFullName string
	ClassName       string
}

type ExcusedDetail struct {
	Name   string
	Reason string
}

func NewSession(userID uint) {
	currentSession[userID] = &AttendanceSession{
		UnexcusedStudents: make(map[uint]string),
		ExcusedStudents:   make(map[uint]ExcusedDetail),
		LateStudents:      make(map[uint]string),
		ClassInfo:         make(map[uint]ClassDetail),
	}
}

func GetSession(userID uint) *AttendanceSession {
	return currentSession[userID]
}

func DeleteSession(userID uint) {
	delete(currentSession, userID)
}

func (s *AttendanceSession) ResolveClassID(telegramUserID uint) uint {
	if s != nil && s.ClassID != 0 {
		return s.ClassID
	}
	return classes.GetClassID(telegramUserID)
}

func (s *AttendanceSession) IsAlreadyMarked(studentID uint) bool {
	_, inUnexcused := s.UnexcusedStudents[studentID]
	_, inExcused := s.ExcusedStudents[studentID]
	_, inLate := s.LateStudents[studentID]
	return inUnexcused || inExcused || inLate
}

func (s *AttendanceSession) GenerateInfoText() string {
	var sb strings.Builder
	sb.WriteString("\n 🔴 **Sababsiz kelmaganlar:**\n")
	if len(s.UnexcusedStudents) == 0 {
		sb.WriteString(" - Yo'q 📑\n")
	} else {
		for _, name := range s.UnexcusedStudents {
			sb.WriteString(fmt.Sprintf("• %s\n", name))
		}
	}

	sb.WriteString("\n🟢 **Sababli kelmaganlar:**\n")
	if len(s.ExcusedStudents) == 0 {
		sb.WriteString("— Yo'q\n")
	} else {
		for _, detail := range s.ExcusedStudents {
			sb.WriteString(fmt.Sprintf("• %s — *%s*\n", detail.Name, detail.Reason))
		}
	}

	sb.WriteString("\n⏰ **Kech kelganlar:**\n")
	if len(s.LateStudents) == 0 {
		sb.WriteString("— Yo'q\n")
	} else {
		for _, name := range s.LateStudents {
			sb.WriteString(fmt.Sprintf("• %s\n", name))
		}
	}

	return sb.String()
}

func (s *AttendanceSession) SetClassContext(userID, classID uint) {
	if s == nil || classID == 0 {
		return
	}
	s.ClassID = classID
	teacherName, className := classes.GetClassInfoByID(classID)
	s.ClassInfo[userID] = ClassDetail{
		TeacherFullName: teacherName,
		ClassName:       className,
	}
}
