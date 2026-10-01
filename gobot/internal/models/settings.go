package models

// BotSetting — botning kalit-qiymat ko'rinishidagi sozlamalari
// (masalan, mehmonlar uchun tabrik funksiyasini yoqish/o'chirish).
// Jadvalni faqat gobot yuritadi.
type BotSetting struct {
	Key   string `gorm:"primaryKey;size:64"`
	Value string `gorm:"size:16;not null"`
}
