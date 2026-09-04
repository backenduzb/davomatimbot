package report

import (
	"strings"
	"unicode"
)

// NormalizeName F.I.SH ni hisobotga yozishdan oldin bir xil ko'rinishga
// keltiradi:
//
//   - boshi/oxiridagi va so'zlar orasidagi ortiqcha bo'shliqlar olib
//     tashlanadi (ko'p bo'shliq bitta bo'shliqqa siqiladi);
//   - har bir so'z alohida Title() qilinadi — birinchi harf katta, qolgani
//     kichik.
//
// Shu tufayli bazadagi turli xil yozuvlar ("ALIYEV ALISHER",
// "aliyev  alisher", "AliYev alisher") hisobotda bir xil
// "Aliyev Alisher" bo'lib chiqadi.
//
// Diqqat: apostrof va chiziqcha so'z ichida qoladi va undan keyingi harf
// KATTA qilinmaydi — o'zbek ismlari to'g'ri yoziladi:
//
//	"O'KTAM"      -> "O'ktam"    ("O'Ktam" emas)
//	"G'ANIYEVA"   -> "G'aniyeva"
//	"ABDULLA-AZIZ"-> "Abdulla-aziz"
//
// Kirill yozuvi ham qo'llab-quvvatlanadi (unicode paketi orqali):
//
//	"КАРИМОВА ДИЛОРОМ" -> "Каримова Дилором"
func NormalizeName(name string) string {
	// strings.Fields boshi/oxiri va orasidagi barcha bo'shliq turlarini
	// (probel, tab, NBSP emas — quyida alohida qaraladi) o'zi tozalaydi.
	fields := strings.FieldsFunc(name, func(r rune) bool {
		// Oddiy bo'shliqlardan tashqari, Excel/Word'dan ko'chirilganda
		// tushib qoladigan uzilmas probel (U+00A0) ham ajratuvchi
		// sifatida qaraladi.
		return unicode.IsSpace(r) || r == '\u00A0'
	})
	if len(fields) == 0 {
		return ""
	}

	for i, word := range fields {
		fields[i] = titleWord(word)
	}

	return strings.Join(fields, " ")
}

// titleWord bitta so'zning birinchi harfini katta, qolganini kichik
// qiladi. Harf bo'lmagan boshlang'ich belgilar (masalan qavs) o'tkazib
// yuboriladi — birinchi HARF katta bo'ladi.
func titleWord(word string) string {
	runes := []rune(word)
	upperDone := false

	for i, r := range runes {
		if !upperDone && unicode.IsLetter(r) {
			runes[i] = unicode.ToUpper(r)
			upperDone = true
			continue
		}
		runes[i] = unicode.ToLower(r)
	}

	return string(runes)
}
