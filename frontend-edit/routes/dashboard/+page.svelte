<script>
  // ============================================================================
  //  DASHBOARD — .xlsx fayl import qilish bo'limi
  //  Manzil: frontend/src/routes/dashboard/+page.svelte
  //
  //  Backend endpoint:  POST /api/import/xlsx   (multipart/form-data, field: "file")
  //  Auth:              Authorization: Bearer <token>
  //  Javob:             { message: string, result: { ... } }
  //
  //  Eslatma: `token` kalitini o'z loyihangizdagi localStorage kalitiga moslang.
  //  Agar Svelte 4 ishlatayotgan bo'lsangiz, runes o'rniga klassik
  //  `let` reaktivligiga o'tkazish kerak (pastga qarang).
  // ============================================================================
  import { onMount } from 'svelte';

  let file = $state(null);
  let fileName = $state('');
  let isDragging = $state(false);
  let isLoading = $state(false);
  let message = $state('');
  let error = $state('');
  let result = $state(null);
  let token = $state('');

  // Token saqlanishi mumkin bo'lgan kalitlar (loyihangizdagisiga moslang)
  const TOKEN_KEYS = ['token', 'auth_token', 'jwt', 'access_token'];

  onMount(() => {
    for (const key of TOKEN_KEYS) {
      const value = localStorage.getItem(key);
      if (value) {
        token = value;
        break;
      }
    }
  });

  function pick(files) {
    error = '';
    message = '';
    result = null;

    if (!files || files.length === 0) return;
    const selected = files[0];

    if (!/\.xlsx$/i.test(selected.name)) {
      error = "Iltimos, faqat .xlsx fayl tanlang.";
      return;
    }

    file = selected;
    fileName = selected.name;
  }

  async function upload() {
    if (!file) {
      error = "Avval .xlsx faylni tanlang.";
      return;
    }
    if (!token) {
      error = "Avtorizatsiya tokeni topilmadi — qayta tizimga kiring.";
      return;
    }

    isLoading = true;
    error = '';
    message = '';

    try {
      const formData = new FormData();
      formData.append('file', file);

      const res = await fetch('/api/import/xlsx', {
        method: 'POST',
        headers: { Authorization: `Bearer ${token}` },
        body: formData,
      });

      const data = await res.json().catch(() => ({}));

      if (!res.ok) {
        error = data.error || `Xatolik yuz berdi (${res.status}).`;
        return;
      }

      message = data.message || "Import muvaffaqiyatli yakunlandi ✅";
      result = data.result || null;
      file = null;
      fileName = '';
    } catch (e) {
      error = "Server bilan bog'lanib bo'lmadi. Internet yoki backendni tekshiring.";
    } finally {
      isLoading = false;
    }
  }

  function onDragOver(e) {
    e.preventDefault();
    isDragging = true;
  }
  function onDragLeave() {
    isDragging = false;
  }
  function onDrop(e) {
    e.preventDefault();
    isDragging = false;
    pick(e.dataTransfer.files);
  }
</script>

<svelte:head>
  <title>Boshqaruv paneli</title>
</svelte:head>

<div class="dashboard">
  <header class="page-head">
    <div>
      <h1>Boshqaruv paneli</h1>
      <p>O'quvchilarni Excel (.xlsx) fayl orqali import qiling</p>
    </div>
  </header>

  <!-- Bu kartani mavjud dashboard sahifangiz ichiga ham qo'yishingiz mumkin -->
  <section class="upload-card">
    <div
      class="dropzone"
      class:dragging={isDragging}
      ondrop={onDrop}
      ondragover={onDragOver}
      ondragleave={onDragLeave}
    >
      <div class="dropzone-icon" aria-hidden="true">📄</div>

      <p class="dropzone-title">
        {fileName ? fileName : "Excel faylni shu yerga tashlang"}
      </p>
      <p class="dropzone-hint">
        yoki tugma orqali tanlang — faqat <strong>.xlsx</strong> formati qabul qilinadi
      </p>

      <input
        id="xlsx-input"
        type="file"
        accept=".xlsx,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
        onchange={(e) => pick(e.currentTarget.files)}
      />

      <div class="actions">
        <label for="xlsx-input" class="btn btn-choose">📁 Fayl tanlash</label>

        <button class="btn btn-upload" disabled={!file || isLoading} onclick={upload}>
          {#if isLoading}
            <span class="spinner" aria-hidden="true"></span>
            <span>Yuklanmoqda…</span>
          {:else}
            <span>⬆️ Yuklab import qilish</span>
          {/if}
        </button>
      </div>
    </div>

    {#if error}
      <div class="notice notice-error" role="alert">
        <span class="notice-icon">⚠️</span>
        <span>{error}</span>
      </div>
    {/if}

    {#if message}
      <div class="notice notice-success" role="status">
        <span class="notice-icon">✅</span>
        <span>{message}</span>
      </div>

      {#if result}
        <div class="result-grid">
          <div class="stat">
            <span class="stat-value">{result.students_created ?? 0}</span>
            <span class="stat-label">yangi o'quvchi</span>
          </div>
          <div class="stat">
            <span class="stat-value">{result.students_linked ?? 0}</span>
            <span class="stat-label">sinfga biriktirildi</span>
          </div>
          <div class="stat">
            <span class="stat-value">{result.classes_created ?? 0}</span>
            <span class="stat-label">yangi sinf</span>
          </div>
          <div class="stat">
            <span class="stat-value">{result.class_names_created ?? 0}</span>
            <span class="stat-label">yangi sinf nomi</span>
          </div>
        </div>
      {/if}
    {/if}
  </section>

  <p class="footnote">
    Import faqat <strong>admin</strong> foydalanuvchilar uchun ruxsat etilgan.
    Sinf nomlari avtomatik normallashtiriladi (masalan, <code>10-a2</code> → <code>10A2</code>).
  </p>
</div>

<style>
  .dashboard {
    min-height: 100vh;
    padding: clamp(20px, 4vw, 48px);
    background:
      radial-gradient(1200px 400px at 10% -10%, rgba(99, 102, 241, 0.08), transparent 60%),
      radial-gradient(1000px 400px at 100% 0%, rgba(139, 92, 246, 0.08), transparent 60%),
      #f1f5f9;
    font-family: 'Inter', 'Segoe UI', system-ui, -apple-system, sans-serif;
    color: #0f172a;
    box-sizing: border-box;
  }

  .page-head {
    max-width: 760px;
    margin: 0 auto 24px;
  }

  .page-head h1 {
    margin: 0;
    font-size: clamp(22px, 3vw, 30px);
    font-weight: 700;
    letter-spacing: -0.02em;
  }

  .page-head p {
    margin: 6px 0 0;
    color: #64748b;
    font-size: 15px;
  }

  .upload-card {
    max-width: 760px;
    margin: 0 auto;
    background: #ffffff;
    border: 1px solid #e2e8f0;
    border-radius: 20px;
    box-shadow:
      0 1px 2px rgba(15, 23, 42, 0.04),
      0 12px 32px rgba(15, 23, 42, 0.06);
    padding: 24px;
  }

  /* --- Dropzone --- */
  .dropzone {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    text-align: center;
    padding: 36px 24px;
    border: 2px dashed #cbd5e1;
    border-radius: 14px;
    background: #f8fafc;
    transition: border-color 0.18s ease, background 0.18s ease, transform 0.18s ease;
    cursor: default;
  }

  .dropzone.dragging {
    border-color: #6366f1;
    background: #eef2ff;
    transform: scale(1.005);
  }

  .dropzone-icon {
    font-size: 44px;
    line-height: 1;
    filter: drop-shadow(0 4px 8px rgba(99, 102, 241, 0.18));
  }

  .dropzone-title {
    margin: 4px 0 0;
    font-size: 16px;
    font-weight: 600;
    color: #1e293b;
    max-width: 100%;
    word-break: break-word;
  }

  .dropzone-hint {
    margin: 0 0 12px;
    font-size: 13px;
    color: #94a3b8;
  }

  .dropzone input[type='file'] {
    display: none;
  }

  .actions {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    justify-content: center;
    margin-top: 6px;
  }

  /* --- Tugmalar --- */
  .btn {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 12px 20px;
    border-radius: 12px;
    font-size: 15px;
    font-weight: 600;
    line-height: 1;
    cursor: pointer;
    transition: transform 0.12s ease, box-shadow 0.12s ease, background 0.12s ease, opacity 0.12s ease;
    border: none;
  }

  .btn:active {
    transform: translateY(1px);
  }

  .btn-choose {
    background: #ffffff;
    color: #4f46e5;
    border: 1px solid #c7d2fe;
    box-shadow: 0 1px 2px rgba(15, 23, 42, 0.05);
  }

  .btn-choose:hover {
    background: #eef2ff;
    border-color: #a5b4fc;
  }

  .btn-upload {
    background: linear-gradient(135deg, #6366f1, #8b5cf6);
    color: #ffffff;
    box-shadow: 0 8px 20px rgba(99, 102, 241, 0.32);
  }

  .btn-upload:hover:not(:disabled) {
    transform: translateY(-1px);
    box-shadow: 0 12px 26px rgba(99, 102, 241, 0.4);
  }

  .btn-upload:disabled {
    opacity: 0.55;
    cursor: not-allowed;
    box-shadow: none;
  }

  /* --- Spinner --- */
  .spinner {
    width: 15px;
    height: 15px;
    border: 2px solid rgba(255, 255, 255, 0.4);
    border-top-color: #ffffff;
    border-radius: 50%;
    animation: spin 0.7s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  /* --- Xabar/status --- */
  .notice {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-top: 16px;
    padding: 12px 16px;
    border-radius: 12px;
    font-size: 14px;
    font-weight: 500;
  }

  .notice-error {
    background: #fef2f2;
    color: #b91c1c;
    border: 1px solid #fecaca;
  }

  .notice-success {
    background: #ecfdf5;
    color: #047857;
    border: 1px solid #a7f3d0;
  }

  /* --- Natija statistikasi --- */
  .result-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
    gap: 12px;
    margin-top: 16px;
  }

  .stat {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 4px;
    padding: 16px 12px;
    background: #f8fafc;
    border: 1px solid #e2e8f0;
    border-radius: 12px;
  }

  .stat-value {
    font-size: 24px;
    font-weight: 700;
    color: #4f46e5;
  }

  .stat-label {
    font-size: 12px;
    color: #64748b;
    text-align: center;
  }

  .footnote {
    max-width: 760px;
    margin: 16px auto 0;
    font-size: 12.5px;
    color: #94a3b8;
    text-align: center;
  }

  .footnote code {
    background: #e2e8f0;
    color: #334155;
    padding: 1px 5px;
    border-radius: 5px;
    font-family: 'JetBrains Mono', ui-monospace, monospace;
  }

  @media (max-width: 520px) {
    .upload-card {
      padding: 16px;
    }
    .dropzone {
      padding: 26px 16px;
    }
    .actions {
      flex-direction: column;
      width: 100%;
    }
    .btn {
      justify-content: center;
    }
  }
</style>
