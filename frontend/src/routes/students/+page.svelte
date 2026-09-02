<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/stores";
    import ResourcePage from "$lib/components/ResourcePage.svelte";
    import type { Column, Field } from "$lib/types/resource";
    import { classesApi, classNamesApi } from "$lib/api/classes";
    import type { Class, ClassName } from "$lib/types";
    import { t } from "$lib/i18n";

    const columns: Column[] = [
        { key: "id", labelKey: "common.id" },
        { key: "full_name", labelKey: "students.full_name" },
        { key: "class_id", labelKey: "students.class" },
    ];

    let classes: Class[] = [];
    let classNames: ClassName[] = [];
    let ready = false;
    let className = "";

    $: classId = $page.url.searchParams.get("class_id") ?? "";

    // Sinf nomini sinf ID orqali topamiz (Class -> ClassName bog'lanishi).
    function classLabel(cls: Class): string {
        const name = classNames.find((n) => n.id === cls.class_name_id);
        return name?.name ?? `#${cls.id}`;
    }

    $: {
        const cls = classId
            ? classes.find((c) => c.id === Number(classId))
            : undefined;
        className = cls ? classLabel(cls) : classId ? `#${classId}` : "";
    }

    // "Sinf" maydoni: ko'rinadigan nom — sinf nomi, qiymat — backend
    // talab qiladigan class ID.
    $: formFields = [
        { key: "full_name", labelKey: "students.full_name", type: "text", required: true },
        {
            key: "class_id",
            labelKey: "students.class",
            type: "select",
            required: true,
            options: classes.map((c) => ({ label: classLabel(c), value: c.id })),
        },
    ] satisfies Field[];

    onMount(async () => {
        try {
            [classes, classNames] = await Promise.all([
                classesApi.getAll(),
                classNamesApi.getAll(),
            ]);
        } catch {
            // Ro'yxat yuklanmasa ham sahifa ishlaydi; sinf nomlari ID
            // ko'rinishida (#id) ko'rsatiladi.
        } finally {
            ready = true;
        }
    });
</script>

<div class="space-y-4">
    {#if classId}
        <div
            class="flex flex-wrap items-center gap-3 rounded-xl border border-blue-200 dark:border-blue-900 bg-blue-50 dark:bg-blue-900/20 px-4 py-3 text-sm text-blue-700 dark:text-blue-400"
        >
            <span>
                {$t("students.class")}: <strong>{className || classId}</strong>
            </span>
            <a
                href="/students"
                class="text-blue-600 dark:text-blue-400 hover:underline"
            >
                {$t("students.show_all")}
            </a>
        </div>
    {/if}

    {#if ready}
        {#key classId}
            <ResourcePage
                titleKey="nav.students"
                endpoint="students"
                {columns}
                {formFields}
                initialFilterValues={{ class_id: classId }}
            />
        {/key}
    {:else}
        <div class="p-6 text-sm text-slate-500">
            {$t("common.loading")}
        </div>
    {/if}
</div>
