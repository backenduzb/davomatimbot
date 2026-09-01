<script lang="ts">
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";
    import { t } from "$lib/i18n";
    import { menuOpen, closeMenu } from "../../stores/menu";
    import { user } from "../../stores/auth";

    const navItems = [
        {
            href: "/dashboard",
            labelKey: "nav.dashboard",
            icon: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" class="size-5"><path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6A2.25 2.25 0 0 1 6 3.75h2.25A2.25 2.25 0 0 1 10.5 6v2.25a2.25 2.25 0 0 1-2.25 2.25H6a2.25 2.25 0 0 1-2.25-2.25V6ZM3.75 15.75A2.25 2.25 0 0 1 6 13.5h2.25a2.25 2.25 0 0 1 2.25 2.25V18a2.25 2.25 0 0 1-2.25 2.25H6A2.25 2.25 0 0 1 3.75 18v-2.25ZM13.5 6a2.25 2.25 0 0 1 2.25-2.25H18A2.25 2.25 0 0 1 20.25 6v2.25A2.25 2.25 0 0 1 18 10.5h-2.25a2.25 2.25 0 0 1-2.25-2.25V6ZM13.5 15.75a2.25 2.25 0 0 1 2.25-2.25H18a2.25 2.25 0 0 1 2.25 2.25V18A2.25 2.25 0 0 1 18 20.25h-2.25A2.25 2.25 0 0 1 13.5 18v-2.25Z" /></svg>`,
        },
        {
            href: "/attendance",
            labelKey: "nav.attendances",
            icon: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" class="size-5"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4M3 7h18M5 7v12a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V7" /></svg>`,
        },
        {
            href: "/students",
            labelKey: "nav.students",
            icon: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" class="size-5"><path stroke-linecap="round" stroke-linejoin="round" d="M15 19.128a9.38 9.38 0 0 0 2.625.372 9.337 9.337 0 0 0 4.121-.952 4.125 4.125 0 0 0-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 0 1 8.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0 1 11.964-3.07M12 6.375a3.375 3.375 0 1 1-6.75 0 3.375 3.375 0 0 1 6.75 0Zm8.25 2.25a2.625 2.625 0 1 1-5.25 0 2.625 2.625 0 0 1 5.25 0Z" /></svg>`,
        },
        {
            href: "/classes",
            labelKey: "nav.classes",
            icon: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" class="size-5"><path stroke-linecap="round" stroke-linejoin="round" d="M3 21h18M6 21V7.5A2.25 2.25 0 0 1 8.25 5.25h7.5A2.25 2.25 0 0 1 18 7.5V21M9 10.5h6M9 13.5h6M9 16.5h6" /></svg>`,
        },
        {
            href: "/profile",
            labelKey: "nav.profile",
            icon: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" class="size-5"><path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0ZM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.933 17.933 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632Z" /></svg>`,
        },
    ];

    const adminItems = [
        {
            href: "/admin/users",
            labelKey: "nav.users",
            icon: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" class="size-5"><path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0ZM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.933 17.933 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632Z" /></svg>`,
        },
        {
            href: "/admin/class-names",
            labelKey: "nav.class_names",
            icon: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" class="size-5"><path stroke-linecap="round" stroke-linejoin="round" d="M4 5h16M4 10h16M4 15h10" /></svg>`,
        },
    ];

    function handleNavClick(href: string) {
        if (window.innerWidth < 768) closeMenu();
        goto(href);
    }

    $: isActive = (href: string) => $page.url.pathname.startsWith(href);
    $: adminActive = adminItems.some((item) => isActive(item.href));
    $: isAdmin = $user?.is_admin ?? false;
    let adminOpen = false;
    $: if (adminActive) adminOpen = true;
</script>

{#if $menuOpen}
    <div
        class="fixed inset-0 bg-black/40 z-30 md:hidden"
        onclick={closeMenu}
        aria-hidden="true"
    ></div>

    <aside
        class={`
        fixed md:static inset-y-0 left-0
        z-40 md:z-auto
        w-64 shrink-0
        flex flex-col
        bg-white dark:bg-slate-800
        border-r border-slate-200 dark:border-slate-700
        min-h-screen
        transform transition-transform duration-300 ease-in-out
        ${$menuOpen ? "translate-x-0" : "-translate-x-full md:translate-x-0"}
    `}
    >
        <div
            class="h-16 flex items-center justify-between px-5 border-b border-slate-200 dark:border-slate-700 shrink-0"
        >
            <span
                class="text-base font-bold text-slate-800 dark:text-slate-100 tracking-tight"
            >
                {$t("app.title")}
            </span>

            <button
                onclick={closeMenu}
                class="md:hidden flex items-center justify-center w-8 h-8 rounded-lg text-slate-500 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-700 transition cursor-pointer"
                aria-label={$t("common.close_menu")}
            >
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="size-5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
                </svg>
            </button>
        </div>
        <nav class="flex-1 p-3 space-y-0.5 overflow-y-auto">
            {#each navItems as item}
                <button
                    onclick={() => handleNavClick(item.href)}
                    class={`
                    w-full flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium
                    transition-all duration-150 cursor-pointer text-left
                    ${
                        isActive(item.href)
                            ? "bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400"
                            : "text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-700 hover:text-slate-800 dark:hover:text-slate-200"
                    }
                `}
                >
                    <span class="shrink-0">{@html item.icon}</span>
                    <span>{$t(item.labelKey)}</span>
                </button>
            {/each}

            {#if isAdmin}
                <div class="pt-1">
                    <button
                        onclick={() => (adminOpen = !adminOpen)}
                        class={`
                        w-full flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium
                        transition-all duration-150 cursor-pointer text-left
                        ${
                            adminActive || adminOpen
                                ? "bg-slate-100 dark:bg-slate-700 text-slate-800 dark:text-slate-200"
                                : "text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-700"
                        }
                    `}
                    >
                        <span class="shrink-0">
                            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" class="size-5">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75 11.25 15 15 9.75m-3-7.036A11.959 11.959 0 0 1 3.598 6 11.99 11.99 0 0 0 3 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285Z" />
                            </svg>
                        </span>
                        <span>{$t("nav.admin")}</span>
                    </button>

                    {#if adminOpen}
                        <div class="mt-1 space-y-0.5 pl-3">
                            {#each adminItems as item}
                                <button
                                    onclick={() => handleNavClick(item.href)}
                                    class={`
                                    w-full flex items-center gap-3 px-3 py-2 rounded-xl text-sm font-medium
                                    transition-all duration-150 cursor-pointer text-left
                                    ${
                                        isActive(item.href)
                                            ? "bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400"
                                            : "text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-700"
                                    }
                                `}
                                >
                                    <span class="shrink-0">{@html item.icon}</span>
                                    <span>{$t(item.labelKey)}</span>
                                </button>
                            {/each}
                        </div>
                    {/if}
                </div>
            {/if}
        </nav>

        <div class="px-5 py-4 border-t border-slate-200 dark:border-slate-700 shrink-0">
            <p class="text-xs text-slate-400 dark:text-slate-500">v1.0.0</p>
        </div>
    </aside>
{/if}
