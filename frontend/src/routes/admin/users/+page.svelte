<script lang="ts">
    import ResourcePage from "$lib/components/ResourcePage.svelte";
    import type { Column, Field } from "$lib/types/resource";

    const columns: Column[] = [
        { key: "id", labelKey: "common.id" },
        { key: "username", labelKey: "common.username" },
        {
            key: "is_admin",
            labelKey: "users.is_admin",
            formatter: (value) => (value ? "✓" : "—"),
        },
        {
            key: "is_banned",
            labelKey: "users.is_banned",
            formatter: (value) => (value ? "✓" : "—"),
        },
        { key: "telegram_id", labelKey: "classes.telegram_id" },
    ];

    const formFields: Field[] = [
        { key: "username", labelKey: "common.username", type: "text", required: true },
        { key: "telegram_id", labelKey: "classes.telegram_id", type: "text" },
        { key: "is_banned", labelKey: "users.is_banned", type: "checkbox" },
    ];
</script>

<ResourcePage
    titleKey="nav.users"
    endpoint="users"
    columns={columns}
    {formFields}
/>
