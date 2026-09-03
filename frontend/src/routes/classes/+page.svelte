<script lang="ts">
    import { onMount } from "svelte";
    import ResourcePage from "$lib/components/ResourcePage.svelte";
    import type { Column, Field } from "$lib/types/resource";
    import {
        classesApi,
        classNamesApi,
        classPromotionApi,
        type PromotionPlan,
    } from "$lib/api/classes";
    import { fetchAllResults } from "$lib/api/client";
    import type { Class, ClassName, Student } from "$lib/types";
    import { t } from "$lib/i18n";

    let classNames: ClassName[] = [];
    let studentCounts: Record<number, number> = {};

    $: columns = [
        { key: "id", labelKey: "common.id" },
        {
            key: "class_name_id",
            labelKey: "classes.class_name",
            formatter: (value: unknown) =>
                classNames.find((n) => n.id === value)?.name ?? String(value ?? ""),
            link: (row: Record<string, unknown>) =>
                `/students?class_id=${row.id}`,
        },
        { key: "teacher_full_name", labelKey: "classes.teacher" },
        { key: "teacher_telegram_id", labelKey: "classes.telegram_id" },
        {
            key: "id",
            labelKey: "classes.student_count",
            formatter: (_value: unknown, row: Record<string, unknown>) =>
                String(studentCounts[row.id as number] ?? 0),
        },
        {
            key: "updated",
            labelKey: "classes.updated",
            formatter: (value: unknown) =>
                value ? $t("common.yes") : $t("common.no"),
        },
    ] satisfies Column[];

    const formFields: Field[] = [
        {
            key: "class_name_id",
            labelKey: "classes.class_name",
            type: "select",
            required: true,
            optionsEndpoint: "class/names",
            optionLabel: "name",
            optionValue: "id",
            editableOption: {
                endpoint: "class/names",
                field: "name",
                toggleLabelKey: "classes.edit_class_name",
                inputLabelKey: "classes.class_name",
            },
        },
        { key: "teacher_full_name", labelKey: "classes.teacher", type: "text", required: true },
        { key: "teacher_telegram_id", labelKey: "classes.telegram_id", type: "text" },
        { key: "updated", labelKey: "classes.updated", type: "checkbox" },
    ];

    async function reloadClassNames() {
        classNames = await classNamesApi.getAll();
    }

    // --- Sinflarni keyingi o'quv yiliga oshirish ---

    let promotionOpen = false;
    let promotionLoading = false;
    let promotionRunning = false;
    let promotionError = "";
    let promotionDone = "";
    let promotionPlans: PromotionPlan[] = [];
    let promotionCounts = { promote: 0, graduate: 0, skip: 0 };
    // Tasodifan bosilmasligi uchun foydalanuvchi qo'shimcha belgi qo'yadi.
    let promotionAcknowledged = false;
    let reloadKey = 0;

    $: promotePlans = promotionPlans.filter((p) => p.action === "promote");
    $: keepPlans = promotionPlans.filter((p) => p.action !== "promote");

    // Tugma bosilganda darhol o'zgartirmaydi — avval preview olib,
    // tasdiqlash oynasini ko'rsatadi.
    async function openPromotion() {
        promotionOpen = true;
        promotionLoading = true;
        promotionError = "";
        promotionDone = "";
        promotionAcknowledged = false;
        promotionPlans = [];
        try {
            const preview = await classPromotionApi.preview();
            promotionPlans = preview.plans ?? [];
            promotionCounts = {
                promote: preview.promote ?? 0,
                graduate: preview.graduate ?? 0,
                skip: preview.skip ?? 0,
            };
        } catch (e) {
            promotionError =
                e instanceof Error ? e.message : $t("common.load_failed");
        } finally {
            promotionLoading = false;
        }
    }

    function closePromotion() {
        if (promotionRunning) return;
        promotionOpen = false;
    }

    async function confirmPromotion() {
        if (!promotionAcknowledged || promotionRunning) return;
        promotionRunning = true;
        promotionError = "";
        try {
            const result = await classPromotionApi.promote();
            promotionDone = $t("classes.promote_success")
                .replace("{promoted}", String(result.promoted))
                .replace("{graduated}", String(result.graduated));
            await reloadClassNames();
            // Jadvalni qayta yuklaymiz.
            reloadKey += 1;
        } catch (e) {
            promotionError =
                e instanceof Error ? e.message : $t("common.save_failed");
        } finally {
            promotionRunning = false;
        }
    }

    onMount(async () => {
        const [names, students] = await Promise.all([
            classNamesApi.getAll(),
            fetchAllResults<Student>("students"),
        ]);
        classNames = names;
        studentCounts = students.reduce<Record<number, number>>((acc, student) => {
            acc[student.class_id] = (acc[student.class_id] ?? 0) + 1;
            return acc;
        }, {});
    });
</script>

<div class="px-6 pt-6">
    <button
        type="button"
        class="inline-flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium
               bg-amber-500 hover:bg-amber-600 text-white shadow-sm transition
               active:scale-[0.98] cursor-pointer disabled:opacity-60"
        on:click={openPromotion}
    >
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M5 10l7-7m0 0l7 7m-7-7v18"
            />
        </svg>
        {$t("classes.promote")}
    </button>
    <p class="mt-2 text-xs text-slate-500 dark:text-slate-400">
        {$t("classes.promote_hint")}
    </p>
</div>

{#key reloadKey}
    <ResourcePage
        titleKey="nav.classes"
        endpoint="classes"
        {columns}
        {formFields}
        on:optionsUpdated={reloadClassNames}
    />
{/key}

{#if promotionOpen}
    <div class="fixed inset-0 z-50 flex items-center justify-center">
        <div
            class="absolute inset-0 bg-black/50"
            on:click={closePromotion}
        ></div>
        <div
            class="relative w-full max-w-2xl mx-4 rounded-2xl bg-white dark:bg-slate-800
                   border border-slate-200 dark:border-slate-700 shadow-xl p-6
                   max-h-[85vh] overflow-y-auto"
        >
            <h3 class="text-lg font-semibold text-slate-800 dark:text-slate-100">
                {$t("classes.promote_confirm_title")}
            </h3>

            {#if promotionLoading}
                <p class="mt-4 text-sm text-slate-500 dark:text-slate-400">
                    {$t("common.loading")}
                </p>
            {:else if promotionDone}
                <div
                    class="mt-4 rounded-xl border border-green-200 dark:border-green-900/50
                           bg-green-50 dark:bg-green-900/20 px-4 py-3
                           text-green-700 dark:text-green-400 text-sm"
                >
                    {promotionDone}
                </div>
                <div class="mt-6 flex justify-end">
                    <button
                        class="h-10 px-5 rounded-lg bg-blue-600 hover:bg-blue-700 text-white text-sm cursor-pointer"
                        on:click={() => (promotionOpen = false)}
                    >
                        {$t("common.close")}
                    </button>
                </div>
            {:else}
                <p class="mt-2 text-sm text-slate-500 dark:text-slate-400">
                    {$t("classes.promote_warning")}
                </p>

                {#if promotionError}
                    <div
                        class="mt-4 rounded-xl border border-red-200 dark:border-red-900
                               bg-red-50 dark:bg-red-900/20 px-4 py-3 text-red-600 dark:text-red-400 text-sm"
                    >
                        {promotionError}
                    </div>
                {/if}

                <div class="mt-4 grid grid-cols-3 gap-3 text-center text-xs">
                    <div class="rounded-xl bg-slate-50 dark:bg-slate-900/50 p-3">
                        <span class="block text-slate-400">{$t("classes.promote_count")}</span>
                        <strong class="text-lg text-slate-800 dark:text-slate-100"
                            >{promotionCounts.promote}</strong
                        >
                    </div>
                    <div class="rounded-xl bg-slate-50 dark:bg-slate-900/50 p-3">
                        <span class="block text-slate-400">{$t("classes.graduate_count")}</span>
                        <strong class="text-lg text-slate-800 dark:text-slate-100"
                            >{promotionCounts.graduate}</strong
                        >
                    </div>
                    <div class="rounded-xl bg-slate-50 dark:bg-slate-900/50 p-3">
                        <span class="block text-slate-400">{$t("classes.skip_count")}</span>
                        <strong class="text-lg text-slate-800 dark:text-slate-100"
                            >{promotionCounts.skip}</strong
                        >
                    </div>
                </div>

                {#if promotePlans.length}
                    <div class="mt-5">
                        <div class="text-xs font-medium text-slate-500 dark:text-slate-400 mb-2">
                            {$t("classes.promote_changes")}
                        </div>
                        <div
                            class="rounded-xl border border-slate-200 dark:border-slate-700 divide-y
                                   divide-slate-100 dark:divide-slate-700 max-h-60 overflow-y-auto"
                        >
                            {#each promotePlans as plan}
                                <div class="flex items-center justify-between px-4 py-2 text-sm">
                                    <span class="text-slate-700 dark:text-slate-200">
                                        <strong>{plan.current_name}</strong>
                                        <span class="mx-2 text-slate-400">→</span>
                                        <strong class="text-green-600 dark:text-green-400"
                                            >{plan.next_name}</strong
                                        >
                                    </span>
                                    <span class="text-xs text-slate-400">
                                        {plan.student_count} {$t("dashboard.students")}
                                    </span>
                                </div>
                            {/each}
                        </div>
                    </div>
                {/if}

                {#if keepPlans.length}
                    <div class="mt-4">
                        <div class="text-xs font-medium text-slate-500 dark:text-slate-400 mb-2">
                            {$t("classes.promote_unchanged")}
                        </div>
                        <div
                            class="rounded-xl border border-slate-200 dark:border-slate-700 divide-y
                                   divide-slate-100 dark:divide-slate-700 max-h-40 overflow-y-auto"
                        >
                            {#each keepPlans as plan}
                                <div class="flex items-center justify-between px-4 py-2 text-sm">
                                    <span class="text-slate-700 dark:text-slate-200"
                                        >{plan.current_name}</span
                                    >
                                    <span class="text-xs text-slate-400">{plan.reason}</span>
                                </div>
                            {/each}
                        </div>
                    </div>
                {/if}

                <label
                    class="mt-5 flex items-start gap-2 text-sm text-slate-700 dark:text-slate-200 cursor-pointer"
                >
                    <input
                        type="checkbox"
                        class="accent-amber-600 h-4 w-4 mt-0.5"
                        bind:checked={promotionAcknowledged}
                    />
                    <span>{$t("classes.promote_acknowledge")}</span>
                </label>

                <div class="mt-6 flex items-center gap-2 justify-end">
                    <button
                        class="h-10 px-5 rounded-lg border dark:border-slate-700 border-slate-300
                               text-slate-600 dark:text-slate-400 text-sm cursor-pointer
                               hover:bg-slate-50 dark:hover:bg-slate-700 disabled:opacity-60"
                        on:click={closePromotion}
                        disabled={promotionRunning}
                    >
                        {$t("common.cancel")}
                    </button>
                    <button
                        class="h-10 px-5 rounded-lg bg-amber-600 hover:bg-amber-700 text-white text-sm
                               cursor-pointer transition active:scale-[0.98]
                               disabled:opacity-50 disabled:cursor-not-allowed"
                        on:click={confirmPromotion}
                        disabled={!promotionAcknowledged ||
                            promotionRunning ||
                            promotionCounts.promote === 0}
                    >
                        {promotionRunning
                            ? $t("classes.promote_running")
                            : $t("classes.promote_confirm")}
                    </button>
                </div>
            {/if}
        </div>
    </div>
{/if}
