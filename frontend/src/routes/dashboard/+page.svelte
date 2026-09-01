<script lang="ts">
    import { onMount, onDestroy, tick } from "svelte";
    import { statisticsApi } from "$lib/api/statistics";
    import { importXlsx } from "$lib/api/importer";
    import type { TodayStatistics, ImportResult } from "$lib/types";
    import type { Chart } from "chart.js";
    import { t } from "$lib/i18n";
    import { user } from "../../stores/auth";

    let stats: TodayStatistics | null = null;
    let loading = true;
    let error = "";

    let distributionChart: Chart | null = null;
    let comparisonChart: Chart | null = null;
    let distributionCanvas!: HTMLCanvasElement;
    let comparisonCanvas!: HTMLCanvasElement;

    // --- XLSX import holati ---
    let fileInput!: HTMLInputElement;
    let file: File | null = null;
    let fileName = "";
    let isDragging = false;
    let uploading = false;
    let importError = "";
    let importMessage = "";
    let importResult: ImportResult | null = null;

    const statusColor = (status: string) => {
        switch (status) {
            case "present":
                return "#2563eb";
            case "absent":
                return "#dc2626";
            case "excused":
                return "#16a34a";
            case "late":
                return "#f59e0b";
            default:
                return "#94a3b8";
        }
    };

    $: statCards = stats
        ? [
              {
                  labelKey: "dashboard.total_classes",
                  value: String(stats.total_classes),
                  icon: "🏫",
                  accent: "text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-900/20",
              },
              {
                  labelKey: "dashboard.total_students",
                  value: String(stats.total_students),
                  icon: "🎓",
                  accent: "text-indigo-600 dark:text-indigo-400 bg-indigo-50 dark:bg-indigo-900/20",
              },
              {
                  labelKey: "dashboard.present",
                  value: String(stats.present),
                  icon: "✅",
                  accent: "text-emerald-600 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-900/20",
              },
              {
                  labelKey: "dashboard.absent",
                  value: String(stats.absent),
                  icon: "🚫",
                  accent: "text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-900/20",
              },
              {
                  labelKey: "dashboard.excused",
                  value: String(stats.excused),
                  icon: "📝",
                  accent: "text-teal-600 dark:text-teal-400 bg-teal-50 dark:bg-teal-900/20",
              },
              {
                  labelKey: "dashboard.late",
                  value: String(stats.late),
                  icon: "⏰",
                  accent: "text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-900/20",
              },
              {
                  labelKey: "dashboard.not_marked",
                  value: String(stats.not_marked),
                  icon: "❔",
                  accent: "text-slate-600 dark:text-slate-400 bg-slate-100 dark:bg-slate-700/40",
              },
              {
                  labelKey: "dashboard.attendance_rate",
                  value: `${Math.round(stats.attendance_percent * 10) / 10}%`,
                  icon: "📈",
                  accent: "text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-900/20",
              },
          ]
        : [];

    async function loadStats() {
        loading = true;
        error = "";
        try {
            stats = await statisticsApi.getToday();
            await tick();
            await renderCharts();
        } catch {
            error = $t("common.load_failed");
        } finally {
            loading = false;
        }
    }

    async function renderCharts() {
        destroyCharts();
        if (!stats) return;

        const { default: Chart } = await import("chart.js/auto");

        if (distributionCanvas) {
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
                                stats.present,
                                stats.absent,
                                stats.excused,
                                stats.late,
                                stats.not_marked,
                            ],
                            backgroundColor: [
                                "#2563eb",
                                "#dc2626",
                                "#16a34a",
                                "#f59e0b",
                                "#94a3b8",
                            ],
                            borderWidth: 2,
                            borderColor: "transparent",
                        },
                    ],
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    cutout: "62%",
                    plugins: {
                        legend: {
                            position: "bottom",
                            labels: { usePointStyle: true, padding: 14 },
                        },
                    },
                },
            });
        }

        if (comparisonCanvas && stats.classes.length > 0) {
            comparisonChart = new Chart(comparisonCanvas, {
                type: "bar",
                data: {
                    labels: stats.classes.map((c) => c.class_name || `#${c.class_id}`),
                    datasets: [
                        {
                            label: $t("dashboard.attendance_rate"),
                            data: stats.classes.map((c) => c.attendance_percent),
                            backgroundColor: "#2563eb",
                            borderRadius: 6,
                            maxBarThickness: 42,
                        },
                    ],
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    scales: {
                        y: {
                            beginAtZero: true,
                            max: 100,
                            ticks: { callback: (v) => `${v}%` },
                        },
                    },
                    plugins: {
                        legend: { display: false },
                    },
                },
            });
        }
    }

    function destroyCharts() {
        distributionChart?.destroy();
        comparisonChart?.destroy();
        distributionChart = null;
        comparisonChart = null;
    }

    // --- XLSX import ---
    function pickFile() {
        fileInput?.click();
    }

    function onFileSelected(e: Event) {
        const input = e.currentTarget as HTMLInputElement;
        setFile(input.files?.[0] ?? null);
        input.value = "";
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

    function onDragOver(e: DragEvent) {
        e.preventDefault();
        isDragging = true;
    }

    function onDragLeave(e: DragEvent) {
        e.preventDefault();
        isDragging = false;
    }

    function onDrop(e: DragEvent) {
        e.preventDefault();
        isDragging = false;
        setFile(e.dataTransfer?.files?.[0] ?? null);
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
            await loadStats();
        } catch (e: unknown) {
            importError = e instanceof Error ? e.message : $t("common.save_failed");
        } finally {
            uploading = false;
        }
    }

    onMount(loadStats);
    onDestroy(destroyCharts);
</script>

<div class="min-h-full bg-slate-50 dark:bg-slate-900 p-4 md:p-6 space-y-6">
    <div class="flex flex-col lg:flex-row lg:items-end lg:justify-between gap-4">
        <div>
            <h1 class="text-2xl font-bold text-slate-800 dark:text-slate-100">
                {$t("dashboard.title")}
            </h1>
            <p class="text-sm text-slate-500 dark:text-slate-400 mt-1">
                {$t("dashboard.subtitle")}
            </p>
        </div>
        <p class="text-sm text-slate-500 dark:text-slate-400">
            {$t("dashboard.date")}: {stats?.date ?? "—"}
        </p>
    </div>

    {#if error}
        <div
            class="rounded-xl border border-red-200 dark:border-red-900 bg-red-50 dark:bg-red-900/20 px-4 py-3 text-red-700 dark:text-red-400 text-sm"
        >
            {error}
        </div>
    {/if}

    {#if loading}
        <div class="p-8 text-center text-slate-500 dark:text-slate-400">
            {$t("dashboard.loading")}
        </div>
    {:else if stats}
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
            {#each statCards as card}
                <div
                    class="bg-white dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 p-4 flex items-center gap-3"
                >
                    <div
                        class={`w-11 h-11 rounded-xl flex items-center justify-center text-xl shrink-0 ${card.accent}`}
                    >
                        {card.icon}
                    </div>
                    <div class="min-w-0">
                        <p class="text-xl font-bold text-slate-800 dark:text-slate-100 leading-none">
                            {card.value}
                        </p>
                        <p class="text-xs text-slate-500 dark:text-slate-400 mt-1 truncate">
                            {$t(card.labelKey)}
                        </p>
                    </div>
                </div>
            {/each}
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div
                class="bg-white dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 p-5"
            >
                <h2
                    class="text-base font-semibold text-slate-800 dark:text-slate-100 mb-4"
                >
                    {$t("dashboard.distribution_chart")}
                </h2>
                <div class="h-72 relative">
                    <canvas bind:this={distributionCanvas}></canvas>
                </div>
            </div>

            <div
                class="bg-white dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 p-5"
            >
                <h2
                    class="text-base font-semibold text-slate-800 dark:text-slate-100 mb-4"
                >
                    {$t("dashboard.comparison_chart")}
                </h2>
                <div class="h-72 relative">
                    <canvas bind:this={comparisonCanvas}></canvas>
                </div>
            </div>
        </div>

        <div
            class="bg-white dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 overflow-hidden"
        >
            <div
                class="px-6 py-4 border-b border-slate-100 dark:border-slate-700"
            >
                <h2
                    class="text-base font-semibold text-slate-800 dark:text-slate-100"
                >
                    {$t("dashboard.class_report")}
                </h2>
            </div>
            {#if stats.classes.length === 0}
                <div class="p-8 text-center text-slate-500 dark:text-slate-400">
                    {$t("dashboard.no_data")}
                </div>
            {:else}
                <div class="overflow-x-auto">
                    <table class="w-full text-sm">
                        <thead
                            class="bg-slate-50 dark:bg-slate-900/50 text-slate-500 dark:text-slate-400"
                        >
                            <tr>
                                <th class="text-left px-5 py-3 font-medium">
                                    {$t("dashboard.class")}
                                </th>
                                <th class="text-left px-5 py-3 font-medium">
                                    {$t("dashboard.total_students")}
                                </th>
                                <th class="text-left px-5 py-3 font-medium">
                                    {$t("dashboard.present")}
                                </th>
                                <th class="text-left px-5 py-3 font-medium">
                                    {$t("dashboard.absent")}
                                </th>
                                <th class="text-left px-5 py-3 font-medium">
                                    {$t("dashboard.excused")}
                                </th>
                                <th class="text-left px-5 py-3 font-medium">
                                    {$t("dashboard.late")}
                                </th>
                                <th class="text-left px-5 py-3 font-medium">
                                    {$t("dashboard.not_marked")}
                                </th>
                                <th class="text-left px-5 py-3 font-medium">
                                    {$t("dashboard.attendance_rate")}
                                </th>
                            </tr>
                        </thead>
                        <tbody>
                            {#each stats.classes as cls}
                                <tr
                                    class="border-t border-slate-100 dark:border-slate-700"
                                >
                                    <td
                                        class="px-5 py-3 font-medium text-slate-800 dark:text-slate-100"
                                    >
                                        {cls.class_name || `#${cls.class_id}`}
                                    </td>
                                    <td class="px-5 py-3 text-slate-600 dark:text-slate-400">
                                        {cls.total_students}
                                    </td>
                                    <td class="px-5 py-3 text-slate-600 dark:text-slate-400">
                                        {cls.present}
                                    </td>
                                    <td class="px-5 py-3 text-slate-600 dark:text-slate-400">
                                        {cls.absent}
                                    </td>
                                    <td class="px-5 py-3 text-slate-600 dark:text-slate-400">
                                        {cls.excused}
                                    </td>
                                    <td class="px-5 py-3 text-slate-600 dark:text-slate-400">
                                        {cls.late}
                                    </td>
                                    <td class="px-5 py-3 text-slate-600 dark:text-slate-400">
                                        {cls.not_marked}
                                    </td>
                                    <td class="px-5 py-3">
                                        <span
                                            class="inline-flex items-center gap-1.5 font-medium"
                                            style:color={statusColor(
                                                cls.attendance_percent >= 75
                                                    ? "present"
                                                    : cls.attendance_percent >= 50
                                                      ? "late"
                                                      : "absent",
                                            )}
                                        >
                                            {cls.attendance_percent}%
                                        </span>
                                    </td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                </div>
            {/if}
        </div>

        {#if $user?.is_admin}
            <div
                class="bg-white dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 p-5"
            >
                <h2
                    class="text-base font-semibold text-slate-800 dark:text-slate-100"
                >
                    {$t("importer.upload_title")}
                </h2>
                <p class="text-sm text-slate-500 dark:text-slate-400 mt-1">
                    {$t("importer.upload_subtitle")}
                </p>

                <div
                    class={`mt-4 rounded-xl border-2 border-dashed p-8 text-center transition cursor-pointer ${
                        isDragging
                            ? "border-blue-500 bg-blue-50 dark:bg-blue-900/20"
                            : "border-slate-300 dark:border-slate-600 bg-slate-50 dark:bg-slate-900/50 hover:border-blue-400"
                    }`}
                    role="button"
                    tabindex="0"
                    onclick={pickFile}
                    onkeydown={(e) => e.key === "Enter" && pickFile()}
                    ondragover={onDragOver}
                    ondragleave={onDragLeave}
                    ondrop={onDrop}
                >
                    <input
                        bind:this={fileInput}
                        type="file"
                        accept=".xlsx,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
                        class="hidden"
                        onchange={onFileSelected}
                    />

                    <div class="flex flex-col items-center gap-2">
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="1.6"
                            class="size-10 text-blue-600 dark:text-blue-400"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5M16.5 12 12 16.5m0 0L7.5 12m4.5 4.5V3"
                            ></path>
                        </svg>

                        <p
                            class="text-sm font-medium text-slate-700 dark:text-slate-300"
                        >
                            {fileName || $t("importer.drop_hint")}
                        </p>
                        <p
                            class="text-xs text-slate-400 dark:text-slate-500"
                        >
                            {$t("importer.file_hint")}
                        </p>
                    </div>
                </div>

                <div class="mt-4 flex flex-wrap items-center gap-3">
                    <button
                        type="button"
                        onclick={pickFile}
                        class="h-10 px-4 rounded-lg border border-slate-300 dark:border-slate-600 text-slate-700 dark:text-slate-300 text-sm font-medium hover:bg-slate-50 dark:hover:bg-slate-700/40 transition cursor-pointer"
                    >
                        {$t("importer.choose_file")}
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

                {#if importError}
                    <div
                        class="mt-4 rounded-xl border border-red-200 dark:border-red-900 bg-red-50 dark:bg-red-900/20 px-4 py-3 text-red-700 dark:text-red-400 text-sm"
                    >
                        {importError}
                    </div>
                {/if}

                {#if importMessage}
                    <div
                        class="mt-4 rounded-xl border border-green-200 dark:border-green-900 bg-green-50 dark:bg-green-900/20 px-4 py-3 text-green-700 dark:text-green-400 text-sm"
                    >
                        {importMessage}
                    </div>
                {/if}

                {#if importResult}
                    <div
                        class="mt-4 grid grid-cols-2 md:grid-cols-4 gap-4"
                    >
                        <div
                            class="rounded-xl bg-slate-50 dark:bg-slate-900/50 border border-slate-200 dark:border-slate-700 p-4 text-center"
                        >
                            <p
                                class="text-xl font-bold text-blue-600 dark:text-blue-400"
                            >
                                {importResult.students_created}
                            </p>
                            <p
                                class="text-xs text-slate-500 dark:text-slate-400 mt-1"
                            >
                                {$t("importer.students_created")}
                            </p>
                        </div>
                        <div
                            class="rounded-xl bg-slate-50 dark:bg-slate-900/50 border border-slate-200 dark:border-slate-700 p-4 text-center"
                        >
                            <p
                                class="text-xl font-bold text-emerald-600 dark:text-emerald-400"
                            >
                                {importResult.students_linked}
                            </p>
                            <p
                                class="text-xs text-slate-500 dark:text-slate-400 mt-1"
                            >
                                {$t("importer.students_linked")}
                            </p>
                        </div>
                        <div
                            class="rounded-xl bg-slate-50 dark:bg-slate-900/50 border border-slate-200 dark:border-slate-700 p-4 text-center"
                        >
                            <p
                                class="text-xl font-bold text-indigo-600 dark:text-indigo-400"
                            >
                                {importResult.classes_created}
                            </p>
                            <p
                                class="text-xs text-slate-500 dark:text-slate-400 mt-1"
                            >
                                {$t("importer.classes_created")}
                            </p>
                        </div>
                        <div
                            class="rounded-xl bg-slate-50 dark:bg-slate-900/50 border border-slate-200 dark:border-slate-700 p-4 text-center"
                        >
                            <p
                                class="text-xl font-bold text-teal-600 dark:text-teal-400"
                            >
                                {importResult.class_names_created}
                            </p>
                            <p
                                class="text-xs text-slate-500 dark:text-slate-400 mt-1"
                            >
                                {$t("importer.class_names_created")}
                            </p>
                        </div>
                    </div>
                {/if}
            </div>
        {/if}
    {/if}
</div>
