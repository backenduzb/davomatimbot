package report

import "testing"

func TestNormalizeName(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		// Asosiy holat: hammasi katta harf.
		{"hammasi katta", "ALIYEV ALISHER BAXTIYOROVICH", "Aliyev Alisher Baxtiyorovich"},
		// Hammasi kichik.
		{"hammasi kichik", "aliyev alisher", "Aliyev Alisher"},
		// Aralash — "boshi katta, o'rtasida yana katta" muammosi.
		{"aralash harflar", "AliYev ALisher", "Aliyev Alisher"},
		{"ichida katta", "KaRIMOVA dILOROM", "Karimova Dilorom"},

		// Ortiqcha bo'shliqlar.
		{"ikki bo'shliq", "Aliyev  Alisher", "Aliyev Alisher"},
		{"ko'p bo'shliq", "Aliyev   Alisher    Baxtiyorovich", "Aliyev Alisher Baxtiyorovich"},
		{"boshi va oxiri", "   Aliyev Alisher   ", "Aliyev Alisher"},
		{"tab va newline", "Aliyev\tAlisher\nBaxtiyorovich", "Aliyev Alisher Baxtiyorovich"},
		{"uzilmas probel", "Aliyev\u00A0Alisher", "Aliyev Alisher"},

		// O'zbek apostrofi — keyingi harf KATTA bo'lmasligi kerak.
		{"o'zbek apostrof", "O'KTAM O'G'LI", "O'ktam O'g'li"},
		{"g' harfi", "G'ANIYEVA GULNORA", "G'aniyeva Gulnora"},
		{"burchakli apostrof", "O‘KTAMOV", "O‘ktamov"},

		// Chiziqcha.
		{"chiziqcha", "ABDULLA-AZIZ QODIROV", "Abdulla-aziz Qodirov"},

		// Kirill.
		{"kirill", "КАРИМОВА ДИЛОРОМ", "Каримова Дилором"},
		{"kirill aralash", "рАХИМОВ жАСУР", "Рахимов Жасур"},

		// Chekka holatlar.
		{"bo'sh", "", ""},
		{"faqat bo'shliq", "     ", ""},
		{"bitta so'z", "aliyev", "Aliyev"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := NormalizeName(tc.in); got != tc.want {
				t.Errorf("NormalizeName(%q) = %q, kutilgan %q", tc.in, got, tc.want)
			}
		})
	}
}

// Bir xil ism turli ko'rinishda yozilgan bo'lsa ham natija bir xil
// bo'lishini tekshiradi (foydalanuvchi shikoyat qilgan asosiy holat).
func TestNormalizeNameIdempotentAcrossVariants(t *testing.T) {
	variants := []string{
		"ALIYEV ALISHER",
		"aliyev alisher",
		"Aliyev  Alisher",
		"  AliYeV   aLIShER ",
		"Aliyev Alisher",
	}

	want := "Aliyev Alisher"
	for _, v := range variants {
		if got := NormalizeName(v); got != want {
			t.Errorf("NormalizeName(%q) = %q, kutilgan %q", v, got, want)
		}
	}
}

// Natijani qayta normalize qilish uni o'zgartirmasligi kerak.
func TestNormalizeNameStable(t *testing.T) {
	in := "  O'KTAM   g'aniyev  "
	once := NormalizeName(in)
	twice := NormalizeName(once)
	if once != twice {
		t.Errorf("barqaror emas: %q -> %q", once, twice)
	}
	if once != "O'ktam G'aniyev" {
		t.Errorf("kutilmagan natija: %q", once)
	}
}
