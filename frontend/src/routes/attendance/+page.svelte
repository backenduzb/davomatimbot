<script lang="ts">
    import { onMount } from "svelte";
    import { attendanceApi } from "$lib/api/attendance";
    import { classesApi, classNamesApi } from "$lib/api/classes";
    import type { AttendanceRecord, AttendanceStatus, Class, ClassName } from "$lib/types";
    import { t } from "$lib/i18n";

    let date = new Date().toISOString().slice(0, 10);
    let classes: Class[] = [];
    let classNames: ClassName[] = [];
    let selectedClassId: number | null = null;
    let records: AttendanceRecord[] = [];
    let loading = false;
    let saving = false;
    let error = "";
    let success = "";

    const statusOptions: AttendanceStatus[] = [
        "present",
        "absent",
        "excused",
        "late",
        "not_marked",
    ];

    function classLabel(classId: number) {
        const cls = classes.find((c) => c.id === classId);
        if (!cls) return `#${classId}`;
        const name = classNames.find((n) => n.id === cls.class_name_id);
        return name?.name ?? `#${classId}`;
    }

    function statusLabel(status: AttendanceStatus) {
        return $t(`attendance.status_${status}`);
    }

    async function loadClasses() {
        [classes, classNames] = await Promise.all([
            classesApi.getAll(),
            classNamesApi.getAll(),
        ]);
        if (!selectedClassId && classes.length > 0) {
            selectedClassId = classes[0].id;
        }
    }

    async function loadAttendance() {
        if (!selectedClassId) {
            records = [];
            return;
        }
        loading = true;
        error = "";
        try {
            const data = await attendanceApi.getToday(date, selectedClassId);
            records = data.records;
        } catch {
            error = $t("common.load_failed");
            records = [];
        } finally {
            loading = false;
        }
    }

    async function saveAttendance() {
        if (!selectedClassId) return;
        saving = true;
        error = "";
        success = "";
        try {
            const toSave = records
                .filter((r) => r.status !== "not_marked")
                .map((r) => ({
                    student_id: r.student_id,
                    status: r.status,
                    reason: r.reason,
                }));

            await attendanceApi.saveBatch({
                date,
                class_id: selectedClassId,
                records: toSave,
            });
            success = $t("attendance.save_success");
            await loadAttendance();
        } catch {
            error = $t("common.save_failed");
        } finally {
            saving = false;
        }
    }

    function updateStatus(index: number, status: AttendanceStatus) {
        records[index].status = status;
        if (status !== "excused") {
            records[index].reason = "";
        }
    }

    onMount(async () => {
        try {
            await loadClasses();
            await loadAttendance();
        } catch {
            error = $t("common.load_failed");
        }
    });
</script>

<div class="min-h-full bg-slate-50 dark:bg-slate-900 p-4 md:p-6 space-y-6">
    <div class="flex flex-col lg:flex-row lg:items-end lg:justify-between gap-4">
        <div>
            <h1 class="text-2xl font-bold text-slate-800 dark:text-slate-100">
                {$t("dashboard.daily_attendance")}
            </h1>
            <p class="text-sm text-slate-500 dark:text-slate-400 mt-1">
                {$t("attendance.subtitle")}
            </p>
        </div>
        <div class="flex flex-col sm:flex-row gap-3">
            <input
                type="date"
                bind:value={date}
                onchange={loadAttendance}
                class="h-10 px-3 rounded-lg border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-800 text-slate-800 dark:text-slate-100"
            />
            <select
                bind:value={selectedClassId}
                onchange={loadAttendance}
                class="h-10 px-3 rounded-lg border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-800 text-slate-800 dark:text-slate-100 min-w-40"
            >
                {#each classes as cls}
                    <option value={cls.id}>{classLabel(cls.id)}</option>
                {/each}
            </select>
            <button
                onclick={saveAttendance}
                disabled={saving || !selectedClassId}
                class="h-10 px-4 rounded-lg bg-blue-600 text-white text-sm font-medium hover:bg-blue-700 disabled:opacity-60"
            >
                {saving ? $t("common.saving") : $t("common.save")}
            </button>
        </div>
    </div>

    {#if error}
        <div class="rounded-xl border border-red-200 dark:border-red-900 bg-red-50 dark:bg-red-900/20 px-4 py-3 text-red-700 dark:text-red-400 text-sm">
            {error}
        </div>
    {/if}

    {#if success}
        <div class="rounded-xl border border-green-200 dark:border-green-900 bg-green-50 dark:bg-green-900/20 px-4 py-3 text-green-700 dark:text-green-400 text-sm">
            {success}
        </div>
    {/if}

    <div class="bg-white dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 overflow-hidden">
        {#if loading}
            <div class="p-8 text-center text-slate-500 dark:text-slate-400">
                {$t("attendance.loading")}
            </div>
        {:else if records.length === 0}
            <div class="p-8 text-center text-slate-500 dark:text-slate-400">
                {$t("dashboard.no_data")}
            </div>
        {:else}
            <div class="overflow-x-auto">
                <table class="w-full text-sm">
                    <thead class="bg-slate-50 dark:bg-slate-900/50 text-slate-500 dark:text-slate-400">
                        <tr>
                            <th class="text-left px-5 py-3 font-medium">{$t("dashboard.student")}</th>
                            <th class="text-left px-5 py-3 font-medium">{$t("dashboard.class")}</th>
                            <th class="text-left px-5 py-3 font-medium">{$t("common.status")}</th>
                            <th class="text-left px-5 py-3 font-medium">{$t("attendance.reason")}</th>
                        </tr>
                    </thead>
                    <tbody>
                        {#each records as record, index}
                            <tr class="border-t border-slate-100 dark:border-slate-700">
                                <td class="px-5 py-3 text-slate-800 dark:text-slate-100">{record.student_name}</td>
                                <td class="px-5 py-3 text-slate-600 dark:text-slate-400">{classLabel(record.class_id)}</td>
                                <td class="px-5 py-3">
                                    <select
                                        value={record.status}
                                        onchange={(e) => updateStatus(index, (e.currentTarget as HTMLSelectElement).value as AttendanceStatus)}
                                        class="h-9 px-2 rounded-lg border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-700 text-slate-800 dark:text-slate-100"
                                    >
                                        {#each statusOptions as status}
                                            <option value={status}>{statusLabel(status)}</option>
                                        {/each}
                                    </select>
                                </td>
                                <td class="px-5 py-3">
                                    {#if record.status === "excused"}
                                        <input
                                            type="text"
                                            bind:value={record.reason}
                                            placeholder={$t("attendance.reason_placeholder")}
                                            class="h-9 w-full min-w-48 px-2 rounded-lg border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-700 text-slate-800 dark:text-slate-100"
                                        />
                                    {:else}
                                        <span class="text-slate-400">—</span>
                                    {/if}
                                </td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
        {/if}
    </div>
</div>
