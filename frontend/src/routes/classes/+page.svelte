<script lang="ts">
    import { onMount } from "svelte";
    import ResourcePage from "$lib/components/ResourcePage.svelte";
    import type { Column, Field } from "$lib/types/resource";
    import { classesApi, classNamesApi } from "$lib/api/classes";
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
        },
        { key: "teacher_full_name", labelKey: "classes.teacher", type: "text", required: true },
        { key: "teacher_telegram_id", labelKey: "classes.telegram_id", type: "text" },
        { key: "updated", labelKey: "classes.updated", type: "checkbox" },
    ];

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

<ResourcePage
    titleKey="nav.classes"
    endpoint="classes"
    {columns}
    {formFields}
/>
