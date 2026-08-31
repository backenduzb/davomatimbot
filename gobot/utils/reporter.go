package utils

import (
	"fmt"
	"strings"
	"time"

	"bot/internal/repository/sessions"
)

func GenerateFullReport(telegramID uint, savedToDB bool) string {
	session := sessions.GetSession(telegramID)
	if session == nil {
		return "⚠️ Sessiya topilmadi."
	}

	var teacherName string = "Noma'lum"
	var className string = "Noma'lum"

	if detail, exists := session.ClassInfo[telegramID]; exists {
		teacherName = detail.TeacherFullName
		className = detail.ClassName
	}

	now := time.Now()
	dateStr := now.Format("02.01.2006")
	timeStr := now.Format("15:04")

	hashtagClass := strings.ReplaceAll(className, " ", "_")

	var b strings.Builder

	b.WriteString("📊 <b>DAVOMAT YAKUNLANDI</b>\n\n")
	fmt.Fprintf(&b, "👤 <b>O'qituvchi:</b> %s\n", teacherName)
	fmt.Fprintf(&b, "🏫 <b>Sinf:</b> %s\n\n", className)

	b.WriteString("🚫 <b>Sababsiz kelmaganlar:</b>\n")
	if len(session.UnexcusedStudents) == 0 {
		b.WriteString(" — Yo'q\n")
	} else {
		for _, name := range session.UnexcusedStudents {
			fmt.Fprintf(&b, " • %s\n", name)
		}
	}

	b.WriteString("\n📝 <b>Sababli kelmaganlar:</b>\n")
	if len(session.ExcusedStudents) == 0 {
		b.WriteString(" — Yo'q\n")
	} else {
		for _, detail := range session.ExcusedStudents {
			fmt.Fprintf(&b, " • %s <i>(Sabab: %s)</i>\n", detail.Name, detail.Reason)
		}
	}

	fmt.Fprintf(&b, "<code>🕒 %s — %s</code>\n", dateStr, timeStr)

	if savedToDB {
		b.WriteString("\n✅ <b>Ma'lumotlar bazasida yangilandi</b>\n")
	}

	fmt.Fprintf(&b, "\n#davomat #sinf_%s\n", hashtagClass)

	return b.String()
}
