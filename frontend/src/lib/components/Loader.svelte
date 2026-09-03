<script lang="ts">
    import Logo from "$lib/assets/logo.svg";
    import { t } from "$lib/i18n";
</script>

<div class="w-full min-h-screen flex items-center justify-center bg-blue-950">
    <div class="flex flex-col items-center">
        <img
            class="w-40 h-40 rounded-full"
            src={Logo}
            alt={$t("app.logo_alt")}
        />
        <div class="w-30 h-20 flex items-center justify-around">
            <span
                class="rounded-full w-4 h-4 bg-blue-600 animate-scale animate-scale-1"
            ></span>
            <span
                class="rounded-full w-4 h-4 bg-blue-600 animate-scale animate-scale-2"
            ></span>
            <span
                class="rounded-full w-4 h-4 bg-blue-600 animate-scale animate-scale-3"
            ></span>
        </div>
    </div>
</div>
