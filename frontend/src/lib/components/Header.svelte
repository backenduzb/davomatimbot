<script lang="ts">
    import { user } from "../../stores/auth";
    import { logout } from "$lib/auth";
    import { goto } from "$app/navigation";
    import { t } from "$lib/i18n";
    import { toggleMenu } from "../../stores/menu";
    import Language from "$lib/components/Language.svelte";
    import Theme from "$lib/components/Theme.svelte";

    let userDropdownOpen = false;
    let langDropdownOpen = false;

    function toggleUserDropdown() {
        userDropdownOpen = !userDropdownOpen;
        langDropdownOpen = false;
    }

    function handleClickOutside(node: HTMLElement) {
        const onClick = (e: MouseEvent) => {
            if (!node.contains(e.target as Node)) {
                userDropdownOpen = false;
            }
        };
        document.addEventListener("click", onClick, true);
        return {
            destroy() {
                document.removeEventListener("click", onClick, true);
            },
        };
    }

    $: username = $user?.username ?? "";
    $: initials = (username[0] ?? "U").toUpperCase();
</script>

<header
    class="sticky top-0 z-30 h-16 w-full
           bg-white dark:bg-slate-800
           border-b border-slate-200 dark:border-slate-700
           shadow-sm flex items-center justify-between px-4 gap-3"
>
    <button
        onclick={toggleMenu}
        aria-label={$t("common.toggle_sidebar")}
        class="flex items-center justify-center w-9 h-9 rounded-lg
               text-slate-500 dark:text-slate-400
               hover:bg-slate-100 dark:hover:bg-slate-700
               transition cursor-pointer shrink-0"
    >
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" class="size-5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5" />
        </svg>
    </button>

    <div class="flex items-center gap-1 sm:gap-2 ml-auto">
        <Language bind:open={langDropdownOpen} />
        <Theme />

        <div class="relative" use:handleClickOutside>
            <button
                onclick={toggleUserDropdown}
                class="flex items-center gap-2 pl-2 pr-1.5 py-1.5 rounded-xl
                       hover:bg-slate-100 dark:hover:bg-slate-700
                       transition cursor-pointer"
            >
                <span class="text-sm font-medium text-slate-700 dark:text-slate-300 hidden sm:block max-w-64 truncate">
                    @{username}
                </span>

                <div class="w-8 h-8 rounded-full bg-blue-600 dark:bg-blue-500 flex items-center justify-center text-white text-sm font-semibold select-none shrink-0">
                    {initials}
                </div>
            </button>

            {#if userDropdownOpen}
                <div class="absolute right-0 top-full mt-2 w-64 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl shadow-lg overflow-hidden z-50">
                    <div class="px-4 py-3 border-b border-slate-100 dark:border-slate-700 bg-slate-50 dark:bg-slate-700/50">
                        <p class="text-xs text-slate-500 dark:text-slate-400 font-medium uppercase tracking-wide">
                            {$t("common.account")}
                        </p>
                        <p class="text-sm font-semibold text-slate-800 dark:text-slate-200 mt-0.5 truncate">
                            @{username}
                        </p>
                    </div>

                    <ul class="py-1">
                        <li>
                            <button
                                onclick={() => {
                                    userDropdownOpen = false;
                                    goto("/profile");
                                }}
                                class="w-full flex items-center gap-3 px-4 py-2.5 text-sm text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700 transition cursor-pointer"
                            >
                                {$t("header.account_settings")}
                            </button>
                        </li>
                        <li class="border-t border-slate-100 dark:border-slate-700 mt-1 pt-1">
                            <button
                                onclick={() => {
                                    userDropdownOpen = false;
                                    logout();
                                }}
                                class="w-full flex items-center gap-3 px-4 py-2.5 text-sm text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 transition cursor-pointer"
                            >
                                {$t("common.logout")}
                            </button>
                        </li>
                    </ul>
                </div>
            {/if}
        </div>
    </div>
</header>
