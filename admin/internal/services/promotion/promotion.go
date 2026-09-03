package promotion

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"admin/internal/models"

	"gorm.io/gorm"
)

// MaxGrade — maktabdagi eng yuqori sinf. Bundan yuqorisi "bitiruvchi".
const MaxGrade = 11

// leadingGrade sinf nomining boshidagi raqamni ajratib oladi.
// Masalan: "8A1" -> 8, "10T" -> 10, "9TSINF" -> 9, "KETGAN" -> 0.
var leadingGrade = regexp.MustCompile(`^\s*(\d{1,2})\s*(.*)$`)

// Action — bitta sinf uchun rejalashtirilgan amal turi.
type Action string

const (
	ActionPromote Action = "promote" // sinf oshiriladi (8A1 -> 9A1)
	// ActionDelete — oshirilganda 11-sinfdan yuqori bo'lib ketadigan sinf.
	// Bunday sinf o'quvchilari va davomat yozuvlari bilan birga o'chiriladi.
	ActionDelete Action = "delete"
	ActionSkip   Action = "skip" // nomida raqam yo'q (masalan "KETGAN")
)

// Plan — bitta sinf bo'yicha o'zgarish rejasi.
type Plan struct {
	ClassID      uint   `json:"class_id"`
	CurrentName  string `json:"current_name"`
	NextName     string `json:"next_name"`
	Action       Action `json:"action"`
	StudentCount int64  `json:"student_count"`
	Reason       string `json:"reason,omitempty"`
}

// Result — oshirish natijasi.
type Result struct {
	Promoted        int    `json:"promoted"`
	Deleted         int    `json:"deleted"`
	StudentsDeleted int64  `json:"students_deleted"`
	Skipped         int    `json:"skipped"`
	ClassNamesMade  int    `json:"class_names_created"`
	Plans           []Plan `json:"plans"`
}

// NextClassName sinf nomidan keyingi o'quv yili nomini hisoblaydi.
// Ikkinchi qiymat — amal turi.
func NextClassName(name string) (string, Action) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return trimmed, ActionSkip
	}

	match := leadingGrade.FindStringSubmatch(trimmed)
	if match == nil {
		// Nomida boshlanuvchi raqam yo'q ("KETGAN", "BITIRUVCHI" va h.k.)
		return trimmed, ActionSkip
	}

	grade, err := strconv.Atoi(match[1])
	if err != nil || grade <= 0 {
		return trimmed, ActionSkip
	}

	suffix := match[2]
	next := grade + 1

	// Oshirilgandan keyin 11-sinfdan yuqori bo'lib ketadigan sinf o'chiriladi.
	if next > MaxGrade {
		return fmt.Sprintf("%d%s", next, suffix), ActionDelete
	}

	return fmt.Sprintf("%d%s", next, suffix), ActionPromote
}

// BuildPlan bazadagi sinflar uchun oshirish rejasini tuzadi (hech narsa
// o'zgartirmaydi). Bu "preview" uchun ishlatiladi — frontend confirm
// oynasida nima o'zgarishini ko'rsatadi.
//
// Optimizatsiya: barcha sinflar va o'quvchilar soni ATIGI 2 ta so'rovda
// olinadi (N+1 yo'q).
func BuildPlan(db *gorm.DB) ([]Plan, error) {
	type classRow struct {
		ClassID   uint
		ClassName string
		Total     int64
	}

	var rows []classRow
	if err := db.Table("classes").
		Select("classes.id as class_id, COALESCE(class_names.name, '') as class_name, COUNT(students.id) as total").
		Joins("LEFT JOIN class_names ON class_names.id = classes.class_name_id").
		Joins("LEFT JOIN students ON students.class_id = classes.id AND students.deleted_at IS NULL").
		Where("classes.deleted_at IS NULL").
		Group("classes.id, class_names.name").
		Order("class_names.name ASC NULLS LAST").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	plans := make([]Plan, 0, len(rows))
	for _, row := range rows {
		nextName, action := NextClassName(row.ClassName)

		plan := Plan{
			ClassID:      row.ClassID,
			CurrentName:  row.ClassName,
			NextName:     nextName,
			Action:       action,
			StudentCount: row.Total,
		}

		switch action {
		case ActionDelete:
			plan.NextName = ""
			plan.Reason = fmt.Sprintf(
				"Bitiruvchi sinf — sinf va %d ta o'quvchi o'chiriladi",
				row.Total,
			)
		case ActionSkip:
			plan.Reason = "Nomida sinf raqami yo'q — o'zgartirilmaydi"
		}

		plans = append(plans, plan)
	}

	return plans, nil
}

// Promote sinflarni keyingi o'quv yiliga oshiradi.
//
// Optimizatsiya: butun amal BITTA tranzaksiyada bajariladi va
// - kerakli sinf nomlari bitta so'rovda o'qiladi,
// - yetishmayotganlari bitta batch INSERT bilan yaratiladi,
// - sinflar yangi nomga bir xil nom bo'yicha guruhlab, har bir nom uchun
//   bitta UPDATE ... WHERE id IN (...) so'rovi bilan yangilanadi.
// Ya'ni sinflar soni qancha bo'lishidan qat'i nazar so'rovlar soni kichik
// bo'lib qoladi — hech qanday lag bo'lmaydi.
func Promote(db *gorm.DB) (*Result, error) {
	plans, err := BuildPlan(db)
	if err != nil {
		return nil, err
	}

	result := &Result{Plans: plans}

	// Oshiriladigan va o'chiriladigan sinflarni ajratamiz.
	targets := make([]Plan, 0, len(plans))
	deleteIDs := make([]uint, 0)
	neededNames := make(map[string]struct{})
	for _, plan := range plans {
		switch plan.Action {
		case ActionPromote:
			targets = append(targets, plan)
			neededNames[plan.NextName] = struct{}{}
		case ActionDelete:
			deleteIDs = append(deleteIDs, plan.ClassID)
		case ActionSkip:
			result.Skipped++
		}
	}

	if len(targets) == 0 && len(deleteIDs) == 0 {
		return result, nil
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// 0) Bitiruvchi sinflarni o'chiramiz: avval davomat yozuvlari, keyin
		//    o'quvchilar, oxirida sinfning o'zi. Har bir bosqich bitta
		//    so'rov — o'quvchilar soni qancha bo'lishidan qat'i nazar tez.
		if len(deleteIDs) > 0 {
			if err := tx.Where("class_id IN ?", deleteIDs).
				Delete(&models.Attendance{}).Error; err != nil {
				return err
			}

			studentRes := tx.Where("class_id IN ?", deleteIDs).
				Delete(&models.Student{})
			if studentRes.Error != nil {
				return studentRes.Error
			}
			result.StudentsDeleted = studentRes.RowsAffected

			classRes := tx.Where("id IN ?", deleteIDs).Delete(&models.Class{})
			if classRes.Error != nil {
				return classRes.Error
			}
			result.Deleted = int(classRes.RowsAffected)
		}

		if len(targets) == 0 {
			return nil
		}

		nameList := make([]string, 0, len(neededNames))
		for name := range neededNames {
			nameList = append(nameList, name)
		}

		// 1) Mavjud sinf nomlarini bitta so'rovda olamiz.
		var existing []models.ClassName
		if err := tx.Where("name IN ?", nameList).Find(&existing).Error; err != nil {
			return err
		}
		nameToID := make(map[string]uint, len(existing))
		for _, item := range existing {
			nameToID[item.Name] = item.ID
		}

		// 2) Yetishmayotgan nomlarni bitta batch INSERT bilan yaratamiz.
		missing := make([]models.ClassName, 0)
		for _, name := range nameList {
			if _, ok := nameToID[name]; !ok {
				missing = append(missing, models.ClassName{Name: name})
			}
		}
		if len(missing) > 0 {
			if err := tx.CreateInBatches(&missing, 100).Error; err != nil {
				return err
			}
			for _, item := range missing {
				nameToID[item.Name] = item.ID
			}
			result.ClassNamesMade = len(missing)
		}

		// 3) Sinflarni yangi nom bo'yicha guruhlaymiz va har bir guruh uchun
		//    bitta UPDATE bajaramiz.
		idsByNewName := make(map[uint][]uint, len(nameList))
		for _, plan := range targets {
			newID, ok := nameToID[plan.NextName]
			if !ok {
				return fmt.Errorf("sinf nomi yaratilmadi: %s", plan.NextName)
			}
			idsByNewName[newID] = append(idsByNewName[newID], plan.ClassID)
		}

		for newNameID, classIDs := range idsByNewName {
			if err := tx.Model(&models.Class{}).
				Where("id IN ?", classIDs).
				Updates(map[string]any{
					"class_name_id": newNameID,
					// Yangi o'quv yili boshlanganda davomat holati tozalanadi.
					"updated": false,
				}).Error; err != nil {
				return err
			}
			result.Promoted += len(classIDs)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
