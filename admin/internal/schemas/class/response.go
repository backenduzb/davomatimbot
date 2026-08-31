package class

type ClassName struct {
	ID uint `json:"id"`
	Name string `json:"name"`
}

type Class struct {
	ID uint `json:"id"`
	Updated bool `json:"updated"`
	ClassNameID uint `json:"class_name_id"`
	TeacherFullName string `json:"teacher_full_name"`
	TeacherTelegramId string `json:"teacher_telegram_id"`
}