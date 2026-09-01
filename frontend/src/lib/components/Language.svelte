<script lang="ts">
    import { locale, t, type Locale } from "$lib/i18n";

    export let open = false;
    export let wrapperClass = "relative";
    export let buttonClass =
        "flex items-center gap-1.5 h-9 px-2 sm:px-3 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-700 transition cursor-pointer";
    export let labelClass =
        "text-sm font-medium text-slate-700 dark:text-slate-300 hidden sm:block";
    export let menuClass =
        "absolute right-0 top-full mt-2 w-44 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl shadow-lg overflow-hidden z-50";

    const languages: { code: Locale; label: string }[] = [
        { code: "uz", label: "O'zbek" },
        { code: "ru", label: "Русский" },
        { code: "en", label: "English" },
        { code: "kaa", label: "Qaraqalpaqsha" },
    ];

    function toggle() {
        open = !open;
    }

    function changeLocale(newLocale: Locale) {
        locale.set(newLocale);
        open = false;
    }

    function handleClickOutside(node: HTMLElement) {
        const onClick = (e: MouseEvent) => {
            if (!node.contains(e.target as Node)) {
                open = false;
            }
        };
        document.addEventListener("click", onClick, true);
        return {
            destroy() {
                document.removeEventListener("click", onClick, true);
            },
        };
    }
</script>

<div class={wrapperClass} use:handleClickOutside>
    <button onclick={toggle} class={buttonClass} title={$t("header.language")}>
        <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.8"
            class="size-5 text-slate-600 dark:text-slate-400 shrink-0"
        >
            <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="m10.5 21 5.25-11.25L21 21m-9-3h7.5M3 5.621a48.474 48.474 0 0 1 6-.371m0 0c1.12 0 2.233.038 3.334.114M9 5.25V3m3.334 2.364C11.176 10.658 7.69 15.08 3 17.502m9.334-12.138c.896.061 1.785.147 2.666.257m-4.589 8.495a18.023 18.023 0 0 1-3.827-5.802"
            />
        </svg>
        <span class={labelClass}>
            {$t(`languages.${$locale}`)}
        </span>
    </button>

    {#if open}
        <div class={menuClass}>
            <div
                class="px-3 py-2 border-b border-slate-100 dark:border-slate-700 bg-slate-50 dark:bg-slate-700/50"
            >
                <p
                    class="text-xs text-slate-500 dark:text-slate-400 font-medium uppercase tracking-wide"
                >
                    {$t("header.language")}
                </p>
            </div>
            <ul class="py-1">
                {#each languages as lang}
                    <li>
                        <button
                            onclick={() => changeLocale(lang.code)}
                            class={`w-full text-left px-4 py-2 text-sm hover:bg-slate-100 dark:hover:bg-slate-700 transition cursor-pointer ${
                                $locale === lang.code
                                    ? "bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400 font-medium"
                                    : "text-slate-700 dark:text-slate-300"
                            }`}
                        >
                            {lang.label}
                        </button>
                    </li>
                {/each}
            </ul>
        </div>
    {/if}
</div>
