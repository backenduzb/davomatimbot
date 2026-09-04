<script lang="ts">
    import { user } from "../../stores/auth";
    import { logout } from "$lib/auth";
    import { t } from "$lib/i18n";
    import { changePassword } from "$lib/api/auth";

    function formatDate(value?: string) {
        if (!value) return "—";
        return new Date(value).toLocaleString();
    }

    const MIN_PASSWORD_LENGTH = 4;

    let currentPassword = $state("");
    let newPassword = $state("");
    let confirmPassword = $state("");
    let showPasswords = $state(false);
    let saving = $state(false);
    let errorKey = $state("");
    let errorText = $state("");
    let successMessage = $state("");

    // Formani yuborish mumkinmi (tugmani o'chirib qo'yish uchun).
    let canSubmit = $derived(
        !saving &&
            currentPassword.length > 0 &&
            newPassword.length >= MIN_PASSWORD_LENGTH &&
            confirmPassword.length > 0,
    );

    function resetFeedback() {
        errorKey = "";
        errorText = "";
        successMessage = "";
    }

    function resetForm() {
        currentPassword = "";
        newPassword = "";
        confirmPassword = "";
        showPasswords = false;
    }

    async function submitPassword(event: SubmitEvent) {
        event.preventDefault();
        if (saving) return;
        resetFeedback();

        if (newPassword.length < MIN_PASSWORD_LENGTH) {
            errorKey = "profile.password_too_short";
            return;
        }
        if (newPassword !== confirmPassword) {
            errorKey = "profile.password_mismatch";
            return;
        }
        if (newPassword === currentPassword) {
            errorKey = "profile.password_same";
            return;
        }

        saving = true;
        try {
            await changePassword({
                current_password: currentPassword,
                new_password: newPassword,
            });
            resetForm();
            successMessage = $t("profile.password_changed");
        } catch (err) {
            // Backend o'zbekcha aniq xabar qaytaradi (masalan "Joriy parol
            // noto'g'ri") — uni to'g'ridan-to'g'ri ko'rsatamiz.
            errorText =
                err instanceof Error && err.message
                    ? err.message
                    : $t("profile.password_error");
        } finally {
            saving = false;
        }
    }
</script>

<div class="min-h-full bg-slate-50 dark:bg-slate-900">
    <main class="max-w-xl mx-auto px-4 py-8 space-y-5">
        <div class="flex items-center gap-4 bg-white dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 shadow-sm px-6 py-5">
            <div class="w-14 h-14 rounded-full bg-blue-600 dark:bg-blue-500 flex items-center justify-center text-white text-xl font-bold select-none shrink-0">
                {($user?.username?.[0] ?? "U").toUpperCase()}
            </div>
            <div>
                <p class="text-lg font-semibold text-slate-800 dark:text-slate-100">
                    @{$user?.username ?? ""}
                </p>
                <p class="text-sm text-slate-500 dark:text-slate-400">
                    {$user?.is_admin ? $t("users.is_admin") : $t("common.user")}
                </p>
            </div>
        </div>

        <div class="bg-white dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 shadow-sm overflow-hidden">
            <div class="px-6 py-4 border-b border-slate-100 dark:border-slate-700">
                <h2 class="text-base font-semibold text-slate-800 dark:text-slate-100">
                    {$t("profile.account_info")}
                </h2>
            </div>
            <div class="px-6 py-5 space-y-4 text-sm">
                <div class="flex justify-between gap-4">
                    <span class="text-slate-500 dark:text-slate-400">{$t("common.username")}</span>
                    <span class="text-slate-800 dark:text-slate-100">@{$user?.username}</span>
                </div>
                <div class="flex justify-between gap-4">
                    <span class="text-slate-500 dark:text-slate-400">{$t("classes.telegram_id")}</span>
                    <span class="text-slate-800 dark:text-slate-100">{$user?.telegram_id || "—"}</span>
                </div>
                <div class="flex justify-between gap-4">
                    <span class="text-slate-500 dark:text-slate-400">{$t("profile.last_seen")}</span>
                    <span class="text-slate-800 dark:text-slate-100">{formatDate($user?.last_seen)}</span>
                </div>
                <div class="flex justify-between gap-4">
                    <span class="text-slate-500 dark:text-slate-400">{$t("users.is_admin")}</span>
                    <span class="text-slate-800 dark:text-slate-100">{$user?.is_admin ? $t("common.yes") : $t("common.no")}</span>
                </div>
            </div>
        </div>

        <div class="bg-white dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 shadow-sm overflow-hidden">
            <div class="px-6 py-4 border-b border-slate-100 dark:border-slate-700">
                <h2 class="text-base font-semibold text-slate-800 dark:text-slate-100">
                    {$t("profile.change_password")}
                </h2>
                <p class="text-xs text-slate-400 dark:text-slate-500 mt-0.5">
                    {$t("profile.change_password_desc")}
                </p>
            </div>

            <form class="px-6 py-5 space-y-4" onsubmit={submitPassword}>
                <div class="space-y-1.5">
                    <label
                        for="current-password"
                        class="block text-sm text-slate-600 dark:text-slate-300"
                    >
                        {$t("profile.current_password")}
                    </label>
                    <input
                        id="current-password"
                        type={showPasswords ? "text" : "password"}
                        bind:value={currentPassword}
                        autocomplete="current-password"
                        disabled={saving}
                        class="w-full h-10 px-3 rounded-lg border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-900 text-sm text-slate-800 dark:text-slate-100 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-60"
                    />
                </div>

                <div class="space-y-1.5">
                    <label
                        for="new-password"
                        class="block text-sm text-slate-600 dark:text-slate-300"
                    >
                        {$t("profile.new_password")}
                    </label>
                    <input
                        id="new-password"
                        type={showPasswords ? "text" : "password"}
                        bind:value={newPassword}
                        autocomplete="new-password"
                        disabled={saving}
                        class="w-full h-10 px-3 rounded-lg border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-900 text-sm text-slate-800 dark:text-slate-100 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-60"
                    />
                    <p class="text-xs text-slate-400 dark:text-slate-500">
                        {$t("profile.password_hint")}
                    </p>
                </div>

                <div class="space-y-1.5">
                    <label
                        for="confirm-password"
                        class="block text-sm text-slate-600 dark:text-slate-300"
                    >
                        {$t("profile.confirm_password")}
                    </label>
                    <input
                        id="confirm-password"
                        type={showPasswords ? "text" : "password"}
                        bind:value={confirmPassword}
                        autocomplete="new-password"
                        disabled={saving}
                        class="w-full h-10 px-3 rounded-lg border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-900 text-sm text-slate-800 dark:text-slate-100 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-60"
                    />
                </div>

                <label
                    class="flex items-center gap-2 text-sm text-slate-600 dark:text-slate-300 select-none cursor-pointer"
                >
                    <input
                        type="checkbox"
                        bind:checked={showPasswords}
                        class="w-4 h-4 rounded border-slate-300 dark:border-slate-600 cursor-pointer"
                    />
                    {$t("profile.show_password")}
                </label>

                {#if errorKey || errorText}
                    <div
                        class="rounded-lg border border-red-200 dark:border-red-900/50 bg-red-50 dark:bg-red-900/20 px-3 py-2 text-sm text-red-700 dark:text-red-400"
                    >
                        {errorKey ? $t(errorKey) : errorText}
                    </div>
                {/if}

                {#if successMessage}
                    <div
                        class="rounded-lg border border-emerald-200 dark:border-emerald-900/50 bg-emerald-50 dark:bg-emerald-900/20 px-3 py-2 text-sm text-emerald-700 dark:text-emerald-400"
                    >
                        {successMessage}
                    </div>
                {/if}

                <div class="flex justify-end pt-1">
                    <button
                        type="submit"
                        disabled={!canSubmit}
                        class="h-10 px-5 rounded-lg bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium transition disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
                    >
                        {saving
                            ? $t("profile.password_saving")
                            : $t("profile.save_password")}
                    </button>
                </div>
            </form>
        </div>

        <div class="bg-white dark:bg-slate-800 rounded-2xl border border-red-200 dark:border-red-900/50 shadow-sm overflow-hidden">
            <div class="px-6 py-5 flex items-center justify-between gap-4">
                <div>
                    <p class="text-sm font-medium text-slate-700 dark:text-slate-300">
                        {$t("profile.logout_title")}
                    </p>
                    <p class="text-xs text-slate-400 dark:text-slate-500 mt-0.5">
                        {$t("profile.logout_desc")}
                    </p>
                </div>
                <button
                    onclick={logout}
                    class="shrink-0 h-10 px-5 rounded-lg border border-red-300 dark:border-red-800 text-red-600 dark:text-red-400 text-sm font-medium hover:bg-red-50 dark:hover:bg-red-900/20 transition cursor-pointer"
                >
                    {$t("common.logout")}
                </button>
            </div>
        </div>
    </main>
</div>
