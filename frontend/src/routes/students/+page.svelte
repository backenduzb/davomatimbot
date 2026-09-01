<script lang="ts">
    import ResourcePage from "$lib/components/ResourcePage.svelte";
    import type { Column, Field } from "$lib/types/resource";
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
</script>

<ResourcePage
    titleKey="nav.students"
    endpoint="students"
    {columns}
    {formFields}
/>
