package student

type Student struct {
	ID uint `json:"id"`
	FullName string `json:"full_name"`
	ClassID uint `json:"class_id"`
}