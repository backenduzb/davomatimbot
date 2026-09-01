<script lang="ts">
    import { t } from "$lib/i18n";
    import { user } from "../../stores/auth";
    import { importXlsx } from "$lib/api/importer";
    import type { ImportResult } from "$lib/types";

    let fileInput!: HTMLInputElement;
    let file: File | null = null;
    let fileName = "";
    let uploading = false;
    let error = "";
    let message = "";
    let result: ImportResult | null = null;

    function pickFile() {
        fileInput?.click();
    }

    function onFileSelected(e: Event) {
        const input = e.currentTarget as HTMLInputElement;
        const selected = input.files?.[0] ?? null;
        input.value = "";
        setFile(selected);
    }

    function setFile(selected: File | null) {
        error = "";
        message = "";
        result = null;
        if (!selected) {
            file = null;
            fileName = "";
            return;
        }
        if (!/\.xlsx$/i.test(selected.name)) {
            error = $t("importer.only_xlsx");
            return;
        }
        file = selected;
        fileName = selected.name;
    }

    async function uploadFile() {
        if (!file) {
            error = $t("importer.file_required");
            return;
        }
        uploading = true;
        error = "";
        message = "";
        result = null;
        try {
            const res = await importXlsx(file);
            message = res.message ?? $t("importer.success");
            result = res.result;
            file = null;
            fileName = "";
        } catch (e: unknown) {
            error = e instanceof Error ? e.message : $t("common.save_failed");
        } finally {
            uploading = false;
        }
    }
</script>

<div class="min-h-full bg-slate-50 dark:bg-slate-900 p-4 md:p-6 space-y-6">
    <div>
        <h1 class="text-2xl font-bold text-slate-800 dark:text-slate-100">
            {$t("dashboard.title")}
        </h1>
        <p class="text-sm text-slate-500 dark:text-slate-400 mt-1">
            {$t("dashboard.subtitle")}
        </p>
    </div>

    {#if $user?.is_admin}
        <div
            class="max-w-xl bg-white dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 p-5"
        >
            <h2
                class="text-base font-semibold text-slate-800 dark:text-slate-100"
            >
                {$t("importer.upload_title")}
            </h2>
            <p class="text-sm text-slate-500 dark:text-slate-400 mt-1">
                {$t("importer.upload_subtitle")}
            </p>

            <input
                bind:this={fileInput}
                type="file"
                accept=".xlsx,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
                class="hidden"
                onchange={onFileSelected}
            />

            <div class="mt-4 flex flex-wrap items-center gap-3">
                <button
                    type="button"
                    onclick={pickFile}
                    class="h-10 px-4 rounded-lg border border-slate-300 dark:border-slate-600 text-slate-700 dark:text-slate-300 text-sm font-medium hover:bg-slate-50 dark:hover:bg-slate-700/40 transition cursor-pointer"
                >
                    {fileName || $t("importer.choose_file")}
                </button>
                <button
                    type="button"
                    onclick={uploadFile}
                    disabled={!file || uploading}
                    class="h-10 px-4 rounded-lg bg-blue-600 text-white text-sm font-medium hover:bg-blue-700 disabled:opacity-60 transition cursor-pointer"
                >
                    {uploading
                        ? $t("importer.uploading")
                        : $t("importer.upload_button")}
                </button>
            </div>

            {#if error}
                <div
                    class="mt-4 rounded-xl border border-red-200 dark:border-red-900 bg-red-50 dark:bg-red-900/20 px-4 py-3 text-red-700 dark:text-red-400 text-sm"
                >
                    {error}
                </div>
            {/if}

            {#if message}
                <div
                    class="mt-4 rounded-xl border border-green-200 dark:border-green-900 bg-green-50 dark:bg-green-900/20 px-4 py-3 text-green-700 dark:text-green-400 text-sm"
                >
                    {message}
                </div>
            {/if}

            {#if result}
                <div class="mt-4 flex flex-wrap gap-x-6 gap-y-2 text-sm text-slate-600 dark:text-slate-300">
                    <span>{$t("importer.students_created")}: <strong>{result.students_created}</strong></span>
                    <span>{$t("importer.students_linked")}: <strong>{result.students_linked}</strong></span>
                    <span>{$t("importer.classes_created")}: <strong>{result.classes_created}</strong></span>
                    <span>{$t("importer.class_names_created")}: <strong>{result.class_names_created}</strong></span>
                </div>
            {/if}
        </div>
    {/if}
</div>
