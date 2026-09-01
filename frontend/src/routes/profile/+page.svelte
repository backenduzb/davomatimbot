<script lang="ts">
    import { user } from "../../stores/auth";
    import { logout } from "$lib/auth";
    import { t } from "$lib/i18n";

    function formatDate(value?: string) {
        if (!value) return "—";
        return new Date(value).toLocaleString();
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
