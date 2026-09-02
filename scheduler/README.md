# Scheduler — kunlik davomat hisoboti

Sarguzasht (background) Go xizmati. `admin/` (HTTP API) va `gobot/` (Telegram bot) xizmatlaridan mustaqil, lekin ular bilan **bir xil PostgreSQL bazasini** ishlatadi.

## Nima qiladi?

1. **Har kuni 16:00 (`Asia/Tashkent`)** — bugungi davomat yozuvlarini bazadan o'qiydi, `davomat.xlsx` hisobotini generatsiya qiladi va Telegram guruhiga hujjat sifatida yuboradi.
2. **Har kuni 00:00 (`Asia/Tashkent`)** — yangi kunlik davomat siklini belgilaydi (log'ga yozadi).

Xizmatning yagona vazifasi shu — Telegram komandlari, o'quvchi/sinf boshqaruvi yoki boshqa bot mantig'i bunda yo'q (ular `gobot/` da qoladi).

## Kunlik sikl

```text
09:00  O'qituvchilar bugungi davomatni kiritadi (bot/admin orqali)
       ↓
16:00  Scheduler bugungi (o'sha sana) yozuvlarni o'qiydi
       ↓
       davomat.xlsx generatsiya qilinadi va Telegram'ga yuboriladi
       ↓
00:00  Yangi kun boshlanadi
       ↓
       O'qituvchilar yangi kunning davomatini kiritadi
       ↓
16:00  Yangi hisobot ...
```

**Muhim:** hisobot yuborilgach davomat yozuvlari **o'chirilmaydi**. Baza dizaynida davomat `(o'quvchi, sana)` unikal yozuv ko'rinishida saqlanadi (har bir kun o'z sanasi bilan ajralgan), shuning uchun 00:00 da hech qanday o'chirish/yaratish talab etilmaydi — tarixiy ma'lumotlar saqlanib qoladi.

## Restart himoyasi (takroriy hisobotdan saqlanish)

Oxirgi muvaffaqiyatli yuborilgan hisobot sanasi bazadagi `scheduler_states` jadvalida saqlanadi. Masalan, 16:00 da hisobot yuborilib, 16:01 da server restart qilingan bo'lsa, scheduler shu sana uchun hisobot allaqachon yuborilganini ko'rib, takrorlamaydi. Jadval faqat scheduler xizmati tomonidan yuritiladi va ilova jadvallariga ta'sir qilmaydi.

## Konfiguratsiya (environment)

| O'zgaruvchi | Matn | Izoh |
| --- | --- | --- |
| `DATABASE_URL` | **majburiy** | PostgreSQL ulanish DSN (admin/gobot bilan bir xil baza) |
| `TELEGRAM_BOT_TOKEN` | **majburiy** | Bot tokeni (hech qachon kod ichida yozilmaydi) |
| `TELEGRAM_CHAT_ID` | `-1003013595617` | Hisobot yuboriladigan guruh ID |
| `TIMEZONE` | `Asia/Tashkent` | Barcha rejalashtirish shu soat sohasida |

Lokal ishda `.env` faylini `scheduler/.env.example` asosida yaratish mumkin.

## XLSX hisobot

| Class | Student | Date | Attendance | Reason |
| --- | --- | --- | --- | --- |
| 10-A | Aliyev Vali | 2026/05/07 | Present |  |
| 10-A | Karimova Dilorom | 2026/05/07 | Absent | Illness |

- Sarlavha qatori: qalin, fonli, markazda, **muzlatilgan** (frozen)
- Barcha ustunlar bo'yicha **Excel filtri** (Class/Student/Date/Attendance/Reason)
- Ustunlar kengligi, chegara va joylashuv sozlangan; sana `yyyy/mm/dd` formatda
- Fayl xotirada generatsiya qilinadi — konteynerda vaqtincha fayllar to'planmaydi

## Telegram xabari

```text
davomat.xlsx

📅 Sana: 2026/09/02
#2026_09_02
```

Sana har doim joriy sana asosida dinamik quriladi.

## Qurish va ishga tushirish

```bash
# Lokal
go mod download
go run ./cmd/main.go

# Docker
docker build -t scheduler .
docker run --env-file .env scheduler

# Testlar (integratsiya sinovlari uchun DATABASE_URL kerak)
go test ./...
```

## Xizmat strukturasi

```text
scheduler/
├── cmd/main.go                  # kirish nuqtasi: soat, DB, bot, sikl
├── config/settings/base.go      # environment yuklash
├── internal/
│   ├── database/db.go           # PostgreSQL ulanish (admin/gobot naqshi)
│   ├── models/                  # mavjud baza modellari (takrorlanmagan, mos saqlangan)
│   ├── repository/
│   │   ├── attendance/          # kunlik davomat qatorlarini o'qish
│   │   └── schedulerstate/      # oxirgi hisobot sanasi (restart himoyasi)
│   └── services/
│       ├── report/              # davomat.xlsx generatsiya (excelize)
│       ├── scheduler/           # 16:00 / 00:00 rejalashtirish va ishlar
│       └── telegram/            # hujjat yuborish (gotgbot)
├── Dockerfile
├── .env.example
└── go.mod
```

## Qo'shimcha

- **XLSX kutubxonasi:** `github.com/xuri/excelize/v2` — admin xizmatidagi `excelize` (360EntSecGroup, v1.4.1) xonadonining (linage) zamonaviy nashri. v1.4.1 da faqat o'qish uchun kerak bo'lgan API bor edi; yozish uchun (avtofiltr, panes, uslublar) zamonaviy v2 API shart bo'ldi.
- **Telegram kutubxonasi:** `gobot/` da ishlatilayotgan `PaulSonOfLars/gotgbot/v2` ning shu versiyasi — faqat `SendDocument` uchun, polling/handlerlar yo'q.
