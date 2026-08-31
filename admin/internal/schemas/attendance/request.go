package attendance

type BatchRecord struct {
	StudentID uint   `json:"student_id" binding:"required"`
	Status    string `json:"status" binding:"required"`
	Reason    string `json:"reason"`
}

type BatchRequest struct {
	Date    string        `json:"date" binding:"required"`
	ClassID uint          `json:"class_id" binding:"required"`
	Records []BatchRecord `json:"records" binding:"required"`
}
