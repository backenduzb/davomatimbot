package importer

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"

	"admin/internal/database"
	"admin/internal/models"

	"github.com/360EntSecGroup-Skylar/excelize"
	"gorm.io/gorm"
)

// Result — import jarayonining natijalari.
type Result struct {
	RowsProcessed     int   `json:"rows_processed"`
	StudentsCreated   int64 `json:"students_created"`
	StudentsLinked    int64 `json:"students_linked"`
	ClassesCreated    int64 `json:"classes_created"`
	ClassNamesCreated int64 `json:"class_names_created"`
}

// rowData — xlsx faylning bitta satridan olingan ma'lumot.
type rowData struct {
	FullName  string
	ClassName string
}

// Import xlsx faylni o'qiydi, sinf nomlarini normal ko'rinishga keltiradi va
// o'quvchilarni, sinflarni hamda sinf nomlarini bazaga yozadi (atomik tranzaksiyada).
func Import(r io.Reader) (*Result, error) {
	data, err := parseRows(r)
	if err != nil {
		return nil, err
	}

	result := &Result{RowsProcessed: len(data)}

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// Sinf nomlari: yagona nomlarni yig'ib, bir so'rovda tekshiramiz.
		classNames := make([]string, 0, len(data))
		classSet := make(map[string]bool, len(data))
		for _, d := range data {
			if !classSet[d.ClassName] {
				classSet[d.ClassName] = true
				classNames = append(classNames, d.ClassName)
			}
		}

		cnIDs, createdCN, err := getOrCreateClassNames(tx, classNames)
		if err != nil {
			return err
		}
		result.ClassNamesCreated = createdCN

		classIDByName, createdClass, err := getOrCreateClasses(tx, cnIDs)
		if err != nil {
			return err
		}
		result.ClassesCreated = createdClass

		createdSt, linkedSt, err := upsertStudents(tx, data, classIDByName)
		if err != nil {
			return err
		}
		result.StudentsCreated = createdSt
		result.StudentsLinked = linkedSt

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// parseRows xlsx faylni ochib, "I.F.Sh" (yoki "F.I.SH") va "Sinf" ustunlarini
// aniqlab, normallashtirilgan qatorlar ro'yxatini qaytaradi.
func parseRows(r io.Reader) ([]rowData, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, fmt.Errorf("faylni o'qishda xato: %w", err)
	}

	sheet := f.GetSheetName(1)
	if sheet == "" {
		return nil, errors.New("faylda hech qanday varaq topilmadi")
	}

	rows := f.GetRows(sheet)
	if len(rows) < 2 {
		return nil, errors.New("fayl bo'sh yoki ma'lumotlar yetarli emas")
	}

	fishIdx, sinfIdx := findColumns(rows[0])
	if fishIdx < 0 || sinfIdx < 0 {
		return nil, errors.New("Excel faylda 'I.F.Sh' yoki 'Sinf' ustuni topilmadi")
	}

	data := make([]rowData, 0, len(rows)-1)
	seen := make(map[string]bool, len(rows)-1)

	for _, row := range rows[1:] {
		fullName := cellValue(row, fishIdx)
		className := cellValue(row, sinfIdx)
		if fullName == "" || className == "" {
			continue
		}

		fullName = titleWords(toLatin(fullName))
		className = NormalizeClassName(className)
		if fullName == "" || className == "" {
			continue
		}

		key := fullName + "\x00" + className
		if seen[key] {
			continue
		}
		seen[key] = true
		data = append(data, rowData{FullName: fullName, ClassName: className})
	}

	if len(data) == 0 {
		return nil, errors.New("mos ma'lumot topilmadi (I.F.Sh va Sinf ustunlarini tekshiring)")
	}

	return data, nil
}

// findColumns sarlavha qatoridan F.I.Sh va Sinf ustunlarining indeksini topadi.
// Ustun nomlari harf-registr, nuqta va bo'shliqlardan qat'i nazar aniqlanadi.
func findColumns(header []string) (fishIdx, sinfIdx int) {
	fishIdx, sinfIdx = -1, -1
	for i, cell := range header {
		switch normalizeHeader(cell) {
		case "ifsh", "fish", "fio", "ismfamilya":
			if fishIdx < 0 {
				fishIdx = i
			}
		case "sinf", "sinfi", "class":
			if sinfIdx < 0 {
				sinfIdx = i
			}
		}
	}
	return fishIdx, sinfIdx
}

// normalizeHeader sarlavhani faqat harf/raqamlarga keltirib, kichik harfda qaytaradi.
// Masalan: "I.F.Sh" -> "ifsh", "F.I.SH" -> "fish", "Sinf" -> "sinf".
func normalizeHeader(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// NormalizeClassName sinf nomini yagona ko'rinishga keltiradi:
// "10A2", "10a2", "10-A2", "10-a2", "10 A2" -> "10A2".
func NormalizeClassName(s string) string {
	var b strings.Builder
	for _, r := range strings.TrimSpace(s) {
		if unicode.IsSpace(r) || r == '-' || r == '_' || r == '.' || r == '\u00a0' || r == '–' || r == '—' {
			continue
		}
		b.WriteRune(unicode.ToUpper(r))
	}
	return b.String()
}

// toLatin kirillcha (o'zbekcha) matnni lotin alifbosiga o'tkazadi.
func toLatin(text string) string {
	var b strings.Builder
	for _, r := range text {
		if v, ok := latinMap[r]; ok {
			b.WriteString(v)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// titleWords so'zning bosh harfini katta, qolganini kichik qiladi
// (Python'dagi str.title() kabi ishlaydi).
func titleWords(s string) string {
	var b strings.Builder
	prevLetter := false
	for _, r := range s {
		if unicode.IsLetter(r) {
			if prevLetter {
				b.WriteRune(unicode.ToLower(r))
			} else {
				b.WriteRune(unicode.ToUpper(r))
			}
			prevLetter = true
		} else {
			b.WriteRune(r)
			prevLetter = false
		}
	}
	return b.String()
}

// cellValue qatordan indeks bo'yicha qiymat oladi (chegaradan himoyalangan).
func cellValue(row []string, idx int) string {
	if idx < 0 || idx >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[idx])
}

// getOrCreateClassNames yagona sinf nomlarini bitta so'rovda tekshirib,
// mavjud bo'lmaganlarini ommaviy (batch) yaratadi.
func getOrCreateClassNames(tx *gorm.DB, names []string) (map[string]uint, int64, error) {
	idByName := make(map[string]uint, len(names))

	var existing []models.ClassName
	if err := tx.Where("name IN ?", names).Find(&existing).Error; err != nil {
		return nil, 0, err
	}
	for _, cn := range existing {
		idByName[cn.Name] = cn.ID
	}

	var toCreate []models.ClassName
	for _, name := range names {
		if _, ok := idByName[name]; !ok {
			toCreate = append(toCreate, models.ClassName{Name: name})
		}
	}

	var created int64
	if len(toCreate) > 0 {
		if err := tx.CreateInBatches(&toCreate, 500).Error; err != nil {
			return nil, 0, err
		}
		created = int64(len(toCreate))
		for _, cn := range toCreate {
			idByName[cn.Name] = cn.ID
		}
	}

	return idByName, created, nil
}

// getOrCreateClasses har bir sinf nomi uchun bitta Class yozuvini ta'minlaydi.
// O'qituvchi hali biriktirilmagan sinflar uchun standart qiymatlar ishlatiladi.
func getOrCreateClasses(tx *gorm.DB, cnIDs map[string]uint) (map[string]uint, int64, error) {
	nameByID := make(map[uint]string, len(cnIDs))
	ids := make([]uint, 0, len(cnIDs))
	for name, id := range cnIDs {
		nameByID[id] = name
		ids = append(ids, id)
	}

	classIDByName := make(map[string]uint, len(cnIDs))
	var existing []models.Class
	if err := tx.Where("class_name_id IN ?", ids).Find(&existing).Error; err != nil {
		return nil, 0, err
	}
	for _, c := range existing {
		if name, ok := nameByID[c.ClassNameID]; ok {
			classIDByName[name] = c.ID
		}
	}

	var toCreate []models.Class
	for name, id := range cnIDs {
		if _, ok := classIDByName[name]; !ok {
			toCreate = append(toCreate, models.Class{
				ClassNameID:       id,
				TeacherFullName:   "Noma'lum",
				TeacherTelegramId: "",
				Updated:           false,
			})
		}
	}

	var created int64
	if len(toCreate) > 0 {
		if err := tx.CreateInBatches(&toCreate, 500).Error; err != nil {
			return nil, 0, err
		}
		created = int64(len(toCreate))
		for _, c := range toCreate {
			if name, ok := nameByID[c.ClassNameID]; ok {
				classIDByName[name] = c.ID
			}
		}
	}

	return classIDByName, created, nil
}

// upsertStudents o'quvchilarni ommaviy yaratadi; mavjud o'quvchiga sinf
// hali biriktirilmagan bo'lsa (class_id = 0), tanlangan sinfga biriktiradi.
func upsertStudents(tx *gorm.DB, data []rowData, classIDByName map[string]uint) (int64, int64, error) {
	names := make([]string, 0, len(data))
	nameSet := make(map[string]bool, len(data))
	for _, d := range data {
		if !nameSet[d.FullName] {
			nameSet[d.FullName] = true
			names = append(names, d.FullName)
		}
	}

	existingByName := make(map[string]models.Student, len(names))
	var existing []models.Student
	if err := tx.Where("full_name IN ?", names).Find(&existing).Error; err != nil {
		return 0, 0, err
	}
	for _, s := range existing {
		existingByName[s.FullName] = s
	}

	var toCreate []models.Student
	pending := make(map[string]bool, len(names))
	linkedNames := make(map[string]bool, len(names))
	var linked int64

	for _, d := range data {
		classID := classIDByName[d.ClassName]
		if classID == 0 {
			continue
		}

		if s, ok := existingByName[d.FullName]; ok {
			if s.ClassID == 0 && !linkedNames[d.FullName] {
				if err := tx.Model(&models.Student{}).Where("id = ?", s.ID).Update("class_id", classID).Error; err != nil {
					return 0, 0, err
				}
				linkedNames[d.FullName] = true
				linked++
			}
			continue
		}

		if pending[d.FullName] {
			continue
		}
		pending[d.FullName] = true
		toCreate = append(toCreate, models.Student{FullName: d.FullName, ClassID: classID})
	}

	var created int64
	if len(toCreate) > 0 {
		if err := tx.CreateInBatches(&toCreate, 500).Error; err != nil {
			return 0, 0, err
		}
		created = int64(len(toCreate))
	}

	return created, linked, nil
}

// latinMap — kirillcha (o'zbekcha) harflarning lotincha ekvivalentlari.
var latinMap = map[rune]string{
	'А': "A", 'а': "a",
	'Б': "B", 'б': "b",
	'В': "V", 'в': "v",
	'Г': "G", 'г': "g",
	'Д': "D", 'д': "d",
	'Е': "E", 'е': "e",
	'Ё': "Yo", 'ё': "yo",
	'Ж': "J", 'ж': "j",
	'З': "Z", 'з': "z",
	'И': "I", 'и': "i",
	'Й': "Y", 'й': "y",
	'К': "K", 'к': "k",
	'Л': "L", 'л': "l",
	'М': "M", 'м': "m",
	'Н': "N", 'н': "n",
	'О': "O", 'о': "o",
	'П': "P", 'п': "p",
	'Р': "R", 'р': "r",
	'С': "S", 'с': "s",
	'Т': "T", 'т': "t",
	'У': "U", 'у': "u",
	'Ф': "F", 'ф': "f",
	'Х': "X", 'х': "x",
	'Ц': "Ts", 'ц': "ts",
	'Ч': "Ch", 'ч': "ch",
	'Ш': "Sh", 'ш': "sh",
	'Щ': "Sh", 'щ': "sh",
	'Ъ': "", 'ъ': "",
	'Ы': "I", 'ы': "i",
	'Э': "E", 'э': "e",
	'Ю': "Yu", 'ю': "yu",
	'Я': "Ya", 'я': "ya",
	'Ғ': "G‘", 'ғ': "g‘",
	'Ў': "O‘", 'ў': "o‘",
	'Қ': "Q", 'қ': "q",
	'Ҳ': "H", 'ҳ': "h",
	'’': "'", '‘': "'",
}
