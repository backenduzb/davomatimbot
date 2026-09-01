<script lang="ts">
    import { onMount } from "svelte";
    import { Chart, registerables } from "chart.js";
    import { statisticsApi } from "$lib/api/statistics";
    import { importXlsx } from "$lib/api/importer";
    import { user } from "../../stores/auth";
    import type { TodayStatistics, ImportResult } from "$lib/types";
    import { t } from "$lib/i18n";

    Chart.register(...registerables);

    let data: TodayStatistics | null = null;
    let loading = true;
    let error = "";
    let distributionChart: Chart | null = null;
    let comparisonChart: Chart | null = null;
    let distributionCanvas: HTMLCanvasElement;
    let comparisonCanvas: HTMLCanvasElement;

    let fileInput!: HTMLInputElement;
    let file: File | null = null;
    let fileName = "";
    let uploading = false;
    let importError = "";
    let importMessage = "";
    let importResult: ImportResult | null = null;

    function todayString() {
        return new Date().toISOString().slice(0, 10);
    }

    async function loadData() {
        loading = true;
        error = "";
        try {
            data = await statisticsApi.getToday();
            await renderCharts();
        } catch (e) {
            error = $t("common.load_failed");
            data = null;
        } finally {
            loading = false;
        }
    }

    function destroyCharts() {
        distributionChart?.destroy();
        comparisonChart?.destroy();
        distributionChart = null;
        comparisonChart = null;
    }

    async function renderCharts() {
        if (!data) return;
        destroyCharts();

        const isDark = document.documentElement.classList.contains("dark");
        const textColor = isDark ? "#cbd5e1" : "#475569";
        const gridColor = isDark ? "#334155" : "#e2e8f0";

        distributionChart = new Chart(distributionCanvas, {
            type: "doughnut",
            data: {
                labels: [
                    $t("dashboard.present"),
                    $t("dashboard.absent"),
                    $t("dashboard.excused"),
                    $t("dashboard.late"),
                    $t("dashboard.not_marked"),
                ],
                datasets: [
                    {
                        data: [
                            data.present,
                            data.absent,
                            data.excused,
                            data.late,
                            data.not_marked,
                        ],
                        backgroundColor: [
                            "#22c55e",
                            "#ef4444",
                            "#f59e0b",
                            "#3b82f6",
                            "#94a3b8",
                        ],
                    },
                ],
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        position: "bottom",
                        labels: { color: textColor },
                    },
                },
            },
        });

        comparisonChart = new Chart(comparisonCanvas, {
            type: "bar",
            data: {
                labels: data.classes.map((c) => c.class_name || `#${c.class_id}`),
                datasets: [
                    {
                        label: $t("dashboard.attendance_rate"),
                        data: data.classes.map((c) => c.attendance_percent),
                        backgroundColor: "#3b82f6",
                        borderRadius: 6,
                    },
                ],
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    x: {
                        ticks: { color: textColor },
                        grid: { color: gridColor },
                    },
                    y: {
                        beginAtZero: true,
                        max: 100,
                        ticks: {
                            color: textColor,
                            callback: (v) => `${v}%`,
                        },
                        grid: { color: gridColor },
                    },
                },
                plugins: {
                    legend: { display: false },
                },
            },
        });
    }

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
        importError = "";
        importMessage = "";
        importResult = null;
        if (!selected) {
            file = null;
            fileName = "";
            return;
        }
        if (!/\.xlsx$/i.test(selected.name)) {
            importError = $t("importer.only_xlsx");
            return;
        }
        file = selected;
        fileName = selected.name;
    }

    async function uploadFile() {
        if (!file) {
            importError = $t("importer.file_required");
            return;
        }
        uploading = true;
        importError = "";
        importMessage = "";
        importResult = null;
        try {
            const res = await importXlsx(file);
            importMessage = res.message ?? $t("importer.success");
            importResult = res.result;
            file = null;
            fileName = "";
            await loadData();
        } catch (e: unknown) {
            importError = e instanceof Error ? e.message : $t("common.save_failed");
        } finally {
            uploading = false;
        }
    }

    onMount(() => {
        loadData();
        return () => destroyCharts();
    });
</script>

<div class="min-h-full bg-slate-50 dark:bg-slate-900 p-4 md:p-6 space-y-6">
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
            <h1 class="text-2xl font-bold text-slate-800 dark:text-slate-100">
                {$t("dashboard.title")}
            </h1>
            <p class="text-sm text-slate-500 dark:text-slate-400 mt-1">
                {data?.date ?? todayString()}
            </p>
        </div>
        <button
            onclick={loadData}
            disabled={loading}
            class="h-10 px-4 rounded-xl bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-sm font-medium text-slate-700 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-700/60 shadow-sm transition disabled:opacity-60 flex items-center justify-center gap-2"
        >
            <svg class="w-4 h-4 {loading ? 'animate-spin' : ''}" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            {$t("common.refresh")}
        </button>
    </div>

    {#if $user?.is_admin}
        <div class="bg-white dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 p-5 shadow-sm transition">
            <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
                <div>
                    <div class="flex items-center gap-2">
                        <span class="p-2 rounded-lg bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400">
                            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                            </svg>
                        </span>
                        <h2 class="text-base font-semibold text-slate-800 dark:text-slate-100">
                            {$t("importer.upload_title")}
                        </h2>
                    </div>
                    <p class="text-sm text-slate-500 dark:text-slate-400 mt-1">
                        {$t("importer.upload_subtitle")}
                    </p>
                </div>

                <input
                    bind:this={fileInput}
                    type="file"
                    accept=".xlsx,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
                    class="hidden"
                    onchange={onFileSelected}
                />

                <div class="flex items-center gap-3">
                    <button
                        type="button"
                        onclick={pickFile}
                        class="h-10 px-4 rounded-xl border border-slate-300 dark:border-slate-600 text-slate-700 dark:text-slate-300 text-sm font-medium hover:bg-slate-50 dark:hover:bg-slate-700/50 transition cursor-pointer flex items-center gap-2"
                    >
                        <svg class="w-4 h-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                        </svg>
                        <span class="max-w-[150px] truncate">{fileName || $t("importer.choose_file")}</span>
                    </button>
                    <button
                        type="button"
                        onclick={uploadFile}
                        disabled={!file || uploading}
                        class="h-10 px-5 rounded-xl bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium shadow-sm shadow-blue-500/20 disabled:opacity-50 transition cursor-pointer flex items-center gap-2"
                    >
                        {#if uploading}
                            <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                            </svg>
                        {/if}
                        {uploading ? $t("importer.uploading") : $t("importer.upload_button")}
                    </button>
                </div>
            </div>

            {#if importError}
                <div class="mt-4 rounded-xl border border-red-200 dark:border-red-900/50 bg-red-50 dark:bg-red-900/20 px-4 py-3 text-red-700 dark:text-red-400 text-sm flex items-center gap-2">
                    <svg class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    {importError}
                </div>
            {/if}

            {#if importMessage}
                <div class="mt-4 rounded-xl border border-green-200 dark:border-green-900/50 bg-green-50 dark:bg-green-900/20 px-4 py-3 text-green-700 dark:text-green-400 text-sm flex items-center gap-2">
                    <svg class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                    </svg>
                    {importMessage}
                </div>
            {/if}

            {#if importResult}
                <div class="mt-4 pt-4 border-t border-slate-100 dark:border-slate-700/50 grid grid-cols-2 sm:grid-cols-4 gap-3 text-xs text-slate-600 dark:text-slate-300">
                    <div class="bg-slate-50 dark:bg-slate-900/50 p-2.5 rounded-xl border border-slate-100 dark:border-slate-800">
                        <span class="block text-slate-400">{$t("importer.students_created")}</span>
                        <strong class="text-sm text-slate-800 dark:text-slate-100 mt-0.5 block">{importResult.students_created}</strong>
                    </div>
                    <div class="bg-slate-50 dark:bg-slate-900/50 p-2.5 rounded-xl border border-slate-100 dark:border-slate-800">
                        <span class="block text-slate-400">{$t("importer.students_linked")}</span>
                        <strong class="text-sm text-slate-800 dark:text-slate-100 mt-0.5 block">{importResult.students_linked}</strong>
                    </div>
                    <div class="bg-slate-50 dark:bg-slate-900/50 p-2.5 rounded-xl border border-slate-100 dark:border-slate-800">
                        <span class="block text-slate-400">{$t("importer.classes_created")}</span>
                        <strong class="text-sm text-slate-800 dark:text-slate-100 mt-0.5 block">{importResult.classes_created}</strong>
                    </div>
                    <div class="bg-slate-50 dark:bg-slate-900/50 p-2.5 rounded-xl border border-slate-100 dark:border-slate-800">
                        <span class="block text-slate-400">{$t("importer.class_names_created")}</span>
                        <strong class="text-sm text-slate-800 dark:text-slate-100 mt-0.5 block">{importResult.class_names_created}</strong>
                    </div>
                </div>
            {/if}
        </div>
    {/if}

    {#if loading}
        <div class="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-6 gap-4">
            {#each Array(6) as _}
                <div class="h-24 rounded-2xl bg-slate-200 dark:bg-slate-800 animate-pulse"></div>
            {/each}
        </div>
        <p class="text-sm text-slate-500 dark:text-slate-400">{$t("dashboard.loading")}</p>
    {:else if error}
        <div class="rounded-2xl border border-red-200 dark:border-red-900 bg-red-50 dark:bg-red-900/20 p-6 text-red-700 dark:text-red-400">
            {error}
        </div>
    {:else if data}
        <div class="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-6 gap-4">
            <div class="bg-white dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 p-4 shadow-sm">
                <p class="text-xs text-slate-500 dark:text-slate-400">{$t("dashboard.total_classes")}</p>
                <p class="text-2xl font-bold text-slate-800 dark:text-slate-100 mt-1">{data.total_classes}</p>
            </div>
            <div class="bg-white dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 p-4 shadow-sm">
                <p class="text-xs text-slate-500 dark:text-slate-400">{$t("dashboard.total_students")}</p>
                <p class="text-2xl font-bold text-slate-800 dark:text-slate-100 mt-1">{data.total_students}</p>
            </div>
            <div class="bg-white dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 p-4 shadow-sm">
                <p class="text-xs text-green-600 dark:text-green-400">{$t("dashboard.present")}</p>
                <p class="text-2xl font-bold text-green-600 dark:text-green-400 mt-1">{data.present}</p>
            </div>
            <div class="bg-white dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 p-4 shadow-sm">
                <p class="text-xs text-red-600 dark:text-red-400">{$t("dashboard.absent")}</p>
                <p class="text-2xl font-bold text-red-600 dark:text-red-400 mt-1">{data.absent}</p>
            </div>
            <div class="bg-white dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 p-4 shadow-sm">
                <p class="text-xs text-blue-600 dark:text-blue-400">{$t("dashboard.late")}</p>
                <p class="text-2xl font-bold text-blue-600 dark:text-blue-400 mt-1">{data.late}</p>
            </div>
            <div class="bg-white dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 p-4 shadow-sm">
                <p class="text-xs text-slate-500 dark:text-slate-400">{$t("dashboard.attendance_rate")}</p>
                <p class="text-2xl font-bold text-blue-600 dark:text-blue-400 mt-1">{data.attendance_percent}%</p>
            </div>
        </div>

        <div class="grid grid-cols-1 xl:grid-cols-2 gap-6">
            <div class="bg-white dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 p-5 shadow-sm">
                <h2 class="text-base font-semibold text-slate-800 dark:text-slate-100 mb-4">
                    {$t("dashboard.distribution_chart")}
                </h2>
                <div class="h-72">
                    <canvas bind:this={distributionCanvas}></canvas>
                </div>
            </div>
            <div class="bg-white dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 p-5 shadow-sm">
                <h2 class="text-base font-semibold text-slate-800 dark:text-slate-100 mb-4">
                    {$t("dashboard.comparison_chart")}
                </h2>
                <div class="h-72">
                    <canvas bind:this={comparisonCanvas}></canvas>
                </div>
            </div>
        </div>

        <!-- Class Report Table -->
        <div class="bg-white dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 overflow-hidden shadow-sm">
            <div class="px-5 py-4 border-b border-slate-200 dark:border-slate-700">
                <h2 class="text-base font-semibold text-slate-800 dark:text-slate-100">
                    {$t("dashboard.class_report")}
                </h2>
            </div>
            <div class="overflow-x-auto">
                <table class="w-full text-sm">
                    <thead class="bg-slate-50 dark:bg-slate-900/50 text-slate-500 dark:text-slate-400">
                        <tr>
                            <th class="text-left px-5 py-3 font-medium">{$t("dashboard.class")}</th>
                            <th class="text-right px-5 py-3 font-medium">{$t("dashboard.students")}</th>
                            <th class="text-right px-5 py-3 font-medium">{$t("dashboard.present")}</th>
                            <th class="text-right px-5 py-3 font-medium">{$t("dashboard.absent")}</th>
                            <th class="text-right px-5 py-3 font-medium">{$t("dashboard.excused")}</th>
                            <th class="text-right px-5 py-3 font-medium">{$t("dashboard.late")}</th>
                            <th class="text-right px-5 py-3 font-medium">{$t("dashboard.not_marked")}</th>
                            <th class="text-right px-5 py-3 font-medium">{$t("dashboard.attendance_rate")}</th>
                        </tr>
                    </thead>
                    <tbody>
                        {#if data.classes.length === 0}
                            <tr>
                                <td colspan="8" class="px-5 py-8 text-center text-slate-500 dark:text-slate-400">
                                    {$t("dashboard.no_data")}
                                </td>
                            </tr>
                        {:else}
                            {#each data.classes as row}
                                <tr class="border-t border-slate-100 dark:border-slate-700/60 hover:bg-slate-50/50 dark:hover:bg-slate-700/30 transition">
                                    <td class="px-5 py-3 font-medium text-slate-800 dark:text-slate-100">
                                        {row.class_name || `#${row.class_id}`}
                                    </td>
                                    <td class="px-5 py-3 text-right">{row.total_students}</td>
                                    <td class="px-5 py-3 text-right text-green-600 dark:text-green-400">{row.present}</td>
                                    <td class="px-5 py-3 text-right text-red-600 dark:text-red-400">{row.absent}</td>
                                    <td class="px-5 py-3 text-right text-amber-600 dark:text-amber-400">{row.excused}</td>
                                    <td class="px-5 py-3 text-right text-blue-600 dark:text-blue-400">{row.late}</td>
                                    <td class="px-5 py-3 text-right text-slate-500">{row.not_marked}</td>
                                    <td class="px-5 py-3 text-right font-semibold">{row.attendance_percent}%</td>
                                </tr>
                            {/each}
                        {/if}
                    </tbody>
                </table>
            </div>
        </div>
    {/if}
</div>