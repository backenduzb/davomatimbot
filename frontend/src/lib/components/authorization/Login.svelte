<script lang="ts">
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";
    import { login } from "$lib/api/auth";
    import { checkAuth } from "$lib/auth";
    import { t } from "$lib/i18n";

    let username = "";
    let password = "";
    let showPassword = false;
    let redirectTo = "/dashboard";
    let validCredentials = true;
    let submitting = false;

    $: {
        const param = $page.url.searchParams.get("redirectTo");
        if (
            param &&
            param.startsWith("/") &&
            !param.includes("undefined") &&
            !param.startsWith("/authenticate")
        ) {
            redirectTo = param;
        }
    }

    async function handleLogin(e: Event) {
        e.preventDefault();
        submitting = true;
        validCredentials = true;

        try {
            const data = await login(username, password);
            localStorage.setItem("token", data.token);
            await checkAuth();
            goto(redirectTo);
        } catch {
            validCredentials = false;
        } finally {
            submitting = false;
        }
    }
</script>

<div class="w-full max-w-80 sm:w-80 mx-auto sm:mx-0">
    <form class="space-y-5" onsubmit={handleLogin}>
        <div class="flex flex-col gap-2">
            <label
                for="login"
                class="text-slate-700 dark:text-slate-300 text-sm font-medium"
            >
                {$t("common.username")}
            </label>
            <input
                id="login"
                type="text"
                bind:value={username}
                class={`w-full h-11 px-3 rounded-lg border transition bg-white dark:bg-slate-700 text-slate-900 dark:text-slate-100
                    focus:outline-none focus:ring-3
                    ${
                        !validCredentials
                            ? "border-red-500 focus:ring-red-500/50 focus:border-red-500 dark:border-red-400"
                            : "border-slate-300 dark:border-slate-600 focus:ring-blue-500/50 focus:border-blue-500"
                    }`}
            />
            {#if !validCredentials}
                <label
                    for="login"
                    class="text-red-500 dark:text-red-400 text-sm font-medium"
                >
                    {$t("auth.wrong_username")}
                </label>
            {/if}
        </div>
        <div class="flex flex-col gap-2">
            <label
                for="password"
                class="text-slate-700 dark:text-slate-300 text-sm font-medium"
            >
                {$t("common.password")}
            </label>
            <div class="relative">
                <input
                    id="password"
                    type={showPassword ? "text" : "password"}
                    bind:value={password}
                    class={`w-full h-11 px-3 rounded-lg border transition bg-white dark:bg-slate-700 text-slate-900 dark:text-slate-100
                        focus:outline-none focus:ring-3
                        ${
                            !validCredentials
                                ? "border-red-500 focus:ring-red-500/50 focus:border-red-500 dark:border-red-400"
                                : "border-slate-300 dark:border-slate-600 focus:ring-blue-500/50 focus:border-blue-500"
                        }`}
                />
                <button
                    type="button"
                    tabindex="-1"
                    onclick={() => (showPassword = !showPassword)}
                    class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 dark:text-slate-500 hover:text-slate-600 dark:hover:text-slate-300 transition cursor-pointer"
                    title={showPassword
                        ? $t("common.hide_password")
                        : $t("common.show_password")}
                >
                    {#if showPassword}
                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" class="size-5">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M3.98 8.223A10.477 10.477 0 0 0 1.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.451 10.451 0 0 1 12 4.5c4.756 0 8.773 3.162 10.065 7.498a10.522 10.522 0 0 1-4.293 5.774M6.228 6.228 3 3m3.228 3.228 3.65 3.65m7.894 7.894L21 21m-3.228-3.228-3.65-3.65m0 0a3 3 0 1 0-4.243-4.243m4.242 4.242L9.88 9.88" />
                        </svg>
                    {:else}
                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" class="size-5">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M2.036 12.322a1.012 1.012 0 0 1 0-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178Z" />
                            <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
                        </svg>
                    {/if}
                </button>
            </div>
            {#if !validCredentials}
                <label
                    for="password"
                    class="text-red-500 dark:text-red-400 text-sm font-medium"
                >
                    {$t("auth.wrong_password")}
                </label>
            {/if}
        </div>
        <button
            type="submit"
            disabled={submitting}
            class="w-full h-11 mt-5 rounded-lg bg-blue-600 dark:bg-blue-500 text-white font-medium
                   hover:bg-blue-700 dark:hover:bg-blue-600 hover:cursor-pointer active:scale-[0.98] transition disabled:opacity-60"
        >
            {submitting ? $t("common.loading") : $t("auth.login_btn")}
        </button>
    </form>
</div>
