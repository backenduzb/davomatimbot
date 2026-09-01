<script lang="ts">
    import { onMount } from "svelte";
    import { loading, isAuthenticated } from "../stores/auth";
    import { checkAuth } from "$lib/auth";
    import { goto } from "$app/navigation";
    import { page } from "$app/stores";
    import Loader from "$lib/components/Loader.svelte";
    import Header from "$lib/components/Header.svelte";
    import Sidebar from "$lib/components/Menu.svelte";
    import "../app.css";

    let authChecked = false;

    onMount(async () => {
        await checkAuth();
        authChecked = true;
    });

    $: if (authChecked) {
        const path = $page.url.pathname;

        if (path !== "/authenticate" && !$isAuthenticated) {
            const currentPath = path + $page.url.search;
            goto(`/authenticate?redirectTo=${encodeURIComponent(currentPath)}`);
        }

        if ($isAuthenticated && path === "/") {
            goto("/dashboard");
        }
    }
</script>

{#if $loading || !authChecked}
    <Loader />
{:else if $isAuthenticated}
    <div class="flex min-h-screen bg-slate-50 dark:bg-slate-900">
        <Sidebar />
        <div class="flex flex-col flex-1 min-w-0">
            <Header />
            <main class="flex-1 overflow-x-hidden">
                <slot />
            </main>
        </div>
    </div>
{:else}
    <slot />
{/if}
