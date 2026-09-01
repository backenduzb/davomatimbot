<script lang="ts">
    import ResourcePage from "$lib/components/ResourcePage.svelte";
    import type { Column, Field } from "$lib/types/resource";

    const columns: Column[] = [
        { key: "id", labelKey: "common.id" },
        { key: "name", labelKey: "common.name" },
    ];

    const formFields: Field[] = [
        { key: "name", labelKey: "common.name", type: "text", required: true },
    ];
</script>

<ResourcePage
    titleKey="nav.class_names"
    endpoint="class/names"
    {columns}
    {formFields}
/>
