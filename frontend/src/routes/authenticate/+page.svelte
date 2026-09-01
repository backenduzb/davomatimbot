<script lang="ts">
    import Logo from "$lib/assets/logo.jpg";
    import LoginComponent from "$lib/components/authorization/Login.svelte";
    import { t } from "$lib/i18n";
    import Language from "$lib/components/Language.svelte";
    import Theme from "$lib/components/Theme.svelte";
</script>

<div class="w-full min-h-screen flex bg-white dark:bg-slate-900">
    <div
        class="md:flex md:items-center min-h-screen md:justify-end w-full sm:w-auto md:h-screen md:w-1/2 py-8 px-4 sm:p-12 md:p-16 sm:rounded-2xl md:rounded-none sm:shadow md:shadow-none sm:bg-white dark:sm:bg-slate-800 ng-tns-c535595194-1"
    >
        <div class="h-30 sm:hidden"></div>
        <div class="relativ flex items-center justify-center flex-col gap-8">
            <h1 class="text-4xl font-bold text-slate-800 dark:text-slate-100">
                {$t("auth.login_title")}
            </h1>

            <div
                class="w-full max-w-80 sm:w-80 mx-auto sm:mx-0 flex items-center justify-between"
            >
                <Language
                    wrapperClass="relative"
                    buttonClass="flex items-center gap-2 h-9 px-3 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-700 transition cursor-pointer"
                    menuClass="absolute left-1/4 z-10 top-50 mt-2 w-44 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl shadow-lg overflow-hidden"
                />
                <Theme />
            </div>
            <LoginComponent />
        </div>
    </div>
    <div
        class="hidden relative w-1/2 bg-blue-950 dark:bg-blue-950 sm:flex items-center justify-center"
    >
        <svg
            viewBox="0 0 960 540"
            width="100%"
            height="100%"
            preserveAspectRatio="xMidYMax slice"
            xmlns="http://www.w3.org/2000/svg"
            class="absolute z-0 inset-0 pointer-events-none ng-tns-c535595194-1"
            ><g
                fill="none"
                stroke="currentColor"
                stroke-width="100"
                class="text-gray-700 opacity-25 ng-tns-c535595194-1"
                ><circle r="234" cx="196" cy="23" class="ng-tns-c535595194-1"
                ></circle><circle
                    r="234"
                    cx="790"
                    cy="491"
                    class="ng-tns-c535595194-1"
                ></circle></g
            ></svg
        >
        <div class="w-full h-1/2 p-40">
            <svg
                viewBox="0 0 220 192"
                width="220"
                height="192"
                fill="none"
                class="absolute -top-16 -right-16 text-gray-700 ng-tns-c535595194-1"
                ><defs class="ng-tns-c535595194-1"
                    ><pattern
                        id="837c3e70-6c3a-44e6-8854-cc48c737b659"
                        x="0"
                        y="0"
                        width="20"
                        height="20"
                        patternUnits="userSpaceOnUse"
                        class="ng-tns-c535595194-1"
                        ><rect
                            x="0"
                            y="0"
                            width="4"
                            height="4"
                            fill="currentColor"
                            class="ng-tns-c535595194-1"
                        ></rect></pattern
                    ></defs
                ><rect
                    width="220"
                    height="192"
                    fill="url(#837c3e70-6c3a-44e6-8854-cc48c737b659)"
                    class="ng-tns-c535595194-1"
                ></rect></svg
            >
            <h2 class="text-5xl font-bold text-white z-20">
                {$t("auth.welcome")}
            </h2>
            <p
                class="mt-6 text-lg tracking-tight leading-6 text-gray-400 ng-tns-c535595194-1"
            >
                {$t("auth.welcome_desc")}
            </p>
            <div class="w-full h-10 mt-15 flex items-center gap-5">
                <img
                    class="rounded-full h-10 w-10 z-20"
                    src={Logo}
                    alt={$t("app.logo_alt")}
                />
                <p
                    class="text-lg tracking-tight z-20 leading-6 text-gray-400 ng-tns-c535595194-1"
                >
                    {$t("auth.org_name")}
                </p>
            </div>
        </div>
    </div>
</div>
