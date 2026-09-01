<script lang="ts">
    import { page } from "$app/stores";
    import ResourcePage from "$lib/components/ResourcePage.svelte";
    import type { Column, Field } from "$lib/types/resource";
    import { classesApi, classNamesApi } from "$lib/api/classes";
    import { t } from "$lib/i18n";

    const columns: Column[] = [
        { key: "id", labelKey: "common.id" },
        { key: "full_name", labelKey: "students.full_name" },
        { key: "class_id", labelKey: "students.class" },
    ];

    const formFields: Field[] = [
        { key: "full_name", labelKey: "students.full_name", type: "text", required: true },
        {
            key: "class_id",
            labelKey: "students.class",
            type: "select",
            required: true,
            optionsEndpoint: "classes",
            optionLabel: "teacher_full_name",
            optionValue: "id",
        },
    ];

    let classId = "";
    let className = "";

    $: classId = $page.url.searchParams.get("class_id") ?? "";

    async function resolveClassName() {
        if (!classId) {
            className = "";
            return;
        }
        try {
            const [classes, names] = await Promise.all([
                classesApi.getAll(),
                classNamesApi.getAll(),
            ]);
            const cls = classes.find((c) => c.id === Number(classId));
            if (cls) {
                const name = names.find((n) => n.id === cls.class_name_id);
                className = name?.name ?? `#${classId}`;
            }
        } catch {
            className = `#${classId}`;
        }
    }

    $: if (classId) resolveClassName();
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

    {#key classId}
        <ResourcePage
            titleKey="nav.students"
            endpoint="students"
            {columns}
            {formFields}
            initialFilterValues={{ class_id: classId }}
        />
    {/key}
</div>
