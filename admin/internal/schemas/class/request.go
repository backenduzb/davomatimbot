package class

type ClassRequest struct {
	ID uint `json:"id"`
	Updated bool `json:"updated"`
	ClassNameID uint `json:"class_name_id"`
	TeacherFullName string `json:"teacher_full_name"`
	TeacherTelegramId string `json:"teacher_telegram_id"`
}