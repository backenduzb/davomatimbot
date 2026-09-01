# Frontend dashboard — .xlsx import tugmasi

Bu papkadagi fayl sizning **alohida frontend repoyingizga** (bu repoda `frontend/`
papkasi bo'sh git-submodule bo'lgani uchun uni shu yerdan o'zgartirib bo'lmadi)
tayyor drop-in kod sifatida tayyorlandi.

## O'rnatish

1. Faylni quyidagi manzilga nusxalang:

   ```
   frontend/src/routes/dashboard/+page.svelte
   ```

   Agar dashboard sahifangizda boshqa kontent bo'lsa, faqat
   `<section class="upload-card">…</section>` blokini va `<script>` qismidagi
   upload logic'ni o'z sahifangizga qo'shib oling.

2. Token kalitini moslang. `upload()` funksiyasi token'ni `localStorage`dan
   quyidagi kalitlar bo'yicha qidiradi (birinchisi topilsa ishlatiladi):

   ```
   token | auth_token | jwt | access_token
   ```

   Loyihangizda token boshqa kalitda (yoki Svelte store'da) saqlansa,
   `TOKEN_KEYS` ro'yxatini yoki `onMount` ichidagi o'qishni o'zgartiring.

3. API manzili. Kod `fetch('/api/import/xlsx', …)` — nisbiy URL ishlatadi.
   Agar frontend va backend turli origin'da bo'lsa, `'/api/import/xlsx'` o'rniga
   to'liq manzil yozing, masalan `'https://api.sizning-sayt.uz/api/import/xlsx'`.

## Backend shartnomasi (admin)

- Endpoint: `POST /api/import/xlsx`
- Format: `multipart/form-data`, maydon nomi **`file`**
- Auth: `Authorization: Bearer <token>` (faqat admin)
- Muvaffaqiyat:

  ```json
  {
    "message": "12 ta o'quvchi muvaffaqiyatli import qilindi ✅",
    "result": {
      "rows_processed": 12,
      "students_created": 10,
      "students_linked": 2,
      "classes_created": 1,
      "class_names_created": 1
    }
  }
  ```

- Xato: `{ "error": "..." }` (400/401/403/409)

## Svelte 4 eslatma

Kod Svelte 5 runes (`$state`, `on:...` o'rniga `onclick/ondrop`) bilan yozilgan.
Svelte 4 ishlatsangiz, `$state(…)` larni oddiy `let` ga almashtirib,
`onclick={upload}` ni `on:click={upload}`, `onchange={…}` ni `on:change={…}`,
`ondrop/ondragover/ondragleave` ni `on:drop/on:dragover/on:dragleave` qilib
o'zgartiring.
