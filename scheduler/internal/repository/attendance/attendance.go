package attendance

import (
	"time"

	"scheduler/internal/database"
)

// ReportRow — kunlik hisobotdagi bitta qator: sinf nomi, o'quvchi F.I.Sh,
// sana, davomat holati va sabab.
type ReportRow struct {
	ClassName string
	Student   string
	Date      time.Time
	Status    string
	Reason    string
}

// ListForDate berilgan sanadagi barcha davomat yozuvlarini sinf va o'quvchi
// nomlari bilan birga qaytaradi (sinf nomi bo'yicha, so'ng ism bo'yicha
// tartiblangan).
//
// Ilova modeli soft-delete (gorm.Model) ishlatgani uchun o'chirilgan
// o'quvchi/sinf yozuvlari hisobotga kirib kirmasligi uchun deleted_at
// shartlari qo'shilgan — admin interfeysidagi ro'yxatlar bilan bir xil
// xulosa.
func ListForDate(date time.Time) ([]ReportRow, error) {
	var rows []ReportRow

	// Sana string literal (YYYY-MM-DD) sifatida beriladi: Postgres'da
	// timestamptz -> date qayta o'zgartirish sessiya soat sohasiga
	// bog'liq bo'lib, kunning bir kun siljishiga olib kelishi mumkin.
	err := database.DB.
		Table("attendances a").
		Select("cn.name AS class_name, s.full_name AS student, a.date AS date, a.status AS status, a.reason AS reason").
		Joins("JOIN students s ON s.id = a.student_id").
		Joins("JOIN classes c ON c.id = a.class_id").
		Joins("JOIN class_names cn ON cn.id = c.class_name_id").
		Where("a.date = ?", date.Format("2006-01-02")).
		Where("a.deleted_at IS NULL").
		Where("s.deleted_at IS NULL").
		Where("c.deleted_at IS NULL").
		Where("cn.deleted_at IS NULL").
		Order("cn.name ASC, s.full_name ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}
