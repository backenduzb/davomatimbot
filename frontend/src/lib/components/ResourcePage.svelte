<script lang="ts">
    import { onDestroy, onMount } from "svelte";
    import { apiFetch, fetchAllResults, API_BASE_URL } from "$lib/api";
    import { t } from "$lib/i18n";
    import type { Column, Field, Filter } from "$lib/types/resource";

    type PaginationItem = number | "ellipsis";

    export let title = "";
    export let titleKey = "";
    export let endpoint = "";
    export let columns: Column[] = [];
    export let formFields: Field[] = [];
    export let filters: Filter[] = [];
    export let initialFilterValues: Record<string, string> = {};

    let items: any[] = [];
    let filteredItems: any[] = [];
    let loading = true;
    let saving = false;
    let errorKey = "";
    let editingId: number | null = null;
    let search = "";
    let searchDebounce: ReturnType<typeof setTimeout> | null = null;
    let filterValues: Record<string, string> = {};
    let pageNum = 1;
    let pageSize = 30;
    let totalPages = 0;
    let totalCount = 0;
    let paginationPages: PaginationItem[] = [];
    let deleteConfirmOpen = false;
    let itemToDelete: any = null;
    let selectOptions: Record<
        string,
        { label?: string; labelKey?: string; value: string | number }[]
    > = {};

    const labelFor = (label?: string, labelKey?: string) => {
        if (labelKey) return $t(labelKey);
        return label ?? "";
    };

    onMount(async () => {
        filterValues = { ...initialFilterValues };
        await loadOptions();
        await loadItems();
    });

    onDestroy(() => {
        if (searchDebounce) {
            clearTimeout(searchDebounce);
        }
    });

    async function loadItems() {
        loading = true;
        errorKey = "";
        try {
            const params = new URLSearchParams();
            params.set("page", String(pageNum));
            params.set("page_size", String(pageSize));
            const normalizedSearch = search.trim();
            if (normalizedSearch) {
                params.set("search", normalizedSearch);
            }
            Object.entries(filterValues).forEach(([key, value]) => {
                if (value) {
                    params.set(key, value);
                }
            });
            const query = params.toString();
            const res = await apiFetch(
                query ? `${endpoint}?${query}` : endpoint,
            );
            const data = await res.json();
            if (Array.isArray(data)) {
                items = data;
                filteredItems = data;
                totalCount = data.length;
                totalPages = data.length ? 1 : 0;
            } else {
                items = data?.results ?? [];
                filteredItems = items;
                totalCount =
                    data?.count ??
                    data?.total ??
                    data?.total_count ??
                    items.length;
                totalPages =
                    data?.total_pages ??
                    (totalCount ? Math.ceil(totalCount / pageSize) : 0);
            }
            if (totalPages > 0 && pageNum > totalPages) {
                pageNum = totalPages;
                await loadItems();
                return;
            }
        } catch (e) {
            errorKey = "common.load_failed";
        } finally {
            loading = false;
        }
    }

    async function loadOptions() {
        const selectFields = formFields.filter((f) => f.type === "select");
        for (const field of selectFields) {
            if (field.options) {
                selectOptions[field.key] = field.options;
                continue;
            }
            if (!field.optionsEndpoint) continue;
            try {
                const list = await fetchAllResults<any>(field.optionsEndpoint);
                
                selectOptions[field.key] = list.map((item: any) => {
                    let label = "";
                    if (field.optionLabel && item[field.optionLabel] !== undefined) {
                        label = item[field.optionLabel];
                    } else if (item.start_date && item.end_date) {
                        label = `${item.start_date} - ${item.end_date}`;
                        if (item.course && item.course.name) label += ` (${item.course.name})`;
                        else if (item.course_name) label += ` (${item.course_name})`;
                    } else if (item.date && item.topic) {
                        label = `${item.date} - ${item.topic}`;
                    } else if (item.first_name || item.last_name) {
                        label = `${item.first_name ?? ""} ${item.last_name ?? ""}`.trim();
                    } else {
                        label = item.name ?? item.topic ?? item.username ?? item.date ?? item.id;
                    }

                    return {
                        label,
                        value: field.optionValue ? item[field.optionValue] : item.id,
                    };
                });
            } catch {
                selectOptions[field.key] = [];
            }
        }
    }

    // initialFormValues boshlang'ich forma qiymatlarini qaytaradi.
    // "create" rejimida filter orqali tanlangan qiymat (masalan, sinf sahifasidan
    // otilgan class_id) mos select maydoniga avtomatik qo'yiladi.
    function initialFormValues(): Record<string, any> {
        const initial: Record<string, any> = {};
        formFields.forEach((f) => {
            if (f.type === "file") {
                initial[f.key] = null;
                return;
            }
            if (f.type === "select") {
                const filtered = initialFilterValues[f.key];
                if (filtered !== undefined && filtered !== "") {
                    const opt = (f.options ?? []).find(
                        (o) => String(o.value) === String(filtered),
                    );
                    initial[f.key] = opt ? opt.value : filtered;
                }
            }
        });
        return initial;
    }

    // toPayloadValue maydon qiymatini backend sxemasiga mos turlarga keltiradi.
    // DOM <select> qiymat string qaytaradi, lekin backend (masalan, class_id)
    // son kutadi — optionning asl tipi (number) saqlanib qolishi kerak.
    function toPayloadValue(field: Field, value: any): any {
        if (field.type === "select") {
            const opt = (selectOptions[field.key] ?? []).find(
                (o) => String(o.value) === String(value),
            );
            if (opt && typeof opt.value === "number") return opt.value;
        }
        return value;
    }

    function resetForm() {
        form = initialFormValues();
        editingId = null;
    }

    let form: Record<string, any> = initialFormValues();

    function startEdit(item: any) {
        form = { ...item };
        formFields.forEach((f) => {
            if (f.type === "file") form[f.key] = null;
            
            if (f.type === "select" && form[f.key] && typeof form[f.key] === "object") {
                form[f.key] = form[f.key].id;
            }
        });
        editingId = item.id;
    }

    async function submitForm(e?: Event) {
        e?.preventDefault();
        saving = true;
        errorKey = "";
        try {
            const method = editingId ? "PATCH" : "POST";
            const url = editingId ? `${endpoint}/${editingId}` : endpoint;
            const hasFile = formFields.some((f) => f.type === "file");

            if (hasFile) {
                const token = localStorage.getItem("token");
                const fd = new FormData();
                formFields.forEach((f) => {
                    const value = form[f.key];
                    if (value === undefined || value === null || value === "")
                        return;
                    if (f.type === "file") {
                        fd.append(f.key, value);
                    } else {
                        fd.append(f.key, String(value));
                    }
                });
                const res = await fetch(`${API_BASE_URL}/api/${url}`, {
                    method,
                    headers: {
                        Authorization: token ? `Bearer ${token}` : "",
                    },
                    body: fd,
                });
                if (!res.ok) throw new Error("save_failed");
            } else {
                const payload: Record<string, any> = {};
                formFields.forEach((f) => {
                    if (form[f.key] === undefined || form[f.key] === "") return;
                    payload[f.key] = toPayloadValue(f, form[f.key]);
                });
                const res = await apiFetch(url, {
                    method,
                    body: JSON.stringify(payload),
                });
                if (!res.ok) throw new Error("save_failed");
            }

            await loadItems();
            resetForm();
        } catch (e) {
            errorKey = "common.save_failed";
        } finally {
            saving = false;
        }
    }

    function removeItem(item: any) {
        itemToDelete = item;
        deleteConfirmOpen = true;
    }

    async function handleConfirmDelete() {
        if (!itemToDelete) return;
        deleteConfirmOpen = false;
        try {
            const res = await apiFetch(`${endpoint}/${itemToDelete.id}`, {
                method: "DELETE",
            });
            if (!res.ok) throw new Error("delete_failed");
            await loadItems();
        } catch {
            errorKey = "common.delete_failed";
        } finally {
            itemToDelete = null;
        }
    }

    function renderValue(value: any, key?: string) {
        if (value === null || value === undefined || value === "") return "-";
        if (typeof value === "boolean")
            return value ? $t("common.yes") : $t("common.no");
        if (Array.isArray(value)) return value.length ? value.length : "-";
        if (key && selectOptions[key]) {
            const hit = selectOptions[key].find(
                (opt) => String(opt.value) === String(value),
            );
            if (hit) return labelFor(hit.label, hit.labelKey);
        }
        if (typeof value === "object") {
            if (value.name) return value.name;
            if (value.username) return value.username;
            if (value.label) return value.label;
            if (value.id) return value.id;
            return JSON.stringify(value);
        }
        return String(value);
    }

    function setPage(nextPage: number) {
        if (nextPage < 1 || (totalPages && nextPage > totalPages)) return;
        if (nextPage === pageNum) return;
        pageNum = nextPage;
        loadItems();
    }

    function reloadFromFirstPage() {
        pageNum = 1;
        loadItems();
    }

    function handleSearchInput() {
        if (searchDebounce) {
            clearTimeout(searchDebounce);
        }
        searchDebounce = setTimeout(() => {
            reloadFromFirstPage();
        }, 300);
    }

    function handleFilterChange() {
        reloadFromFirstPage();
    }

    function buildPagination(current: number, total: number): PaginationItem[] {
        if (total <= 1) {
            return [];
        }

        const range = (start: number, end: number) => {
            const normalizedStart = Math.max(1, Math.min(start, total));
            const normalizedEnd = Math.max(
                normalizedStart,
                Math.min(end, total),
            );
            const list: number[] = [];
            for (let i = normalizedStart; i <= normalizedEnd; i += 1) {
                list.push(i);
            }
            return list;
        };

        if (total <= 7) {
            return range(1, total);
        }

        if (current <= 4) {
            return [...range(1, 5), "ellipsis", total];
        }

        if (current >= total - 3) {
            return [1, "ellipsis", ...range(total - 4, total)];
        }

        return [
            1,
            "ellipsis",
            ...range(current - 1, current + 1),
            "ellipsis",
            total,
        ];
    }

    $: paginationPages = buildPagination(pageNum, totalPages);

    $: filteredItems = items;
</script>

<section class="p-6 space-y-6">
    <div class="flex items-center justify-between gap-4">
        <div>
            <h1
                class="text-2xl font-semibold text-slate-900 dark:text-slate-100"
            >
                {titleKey ? $t(titleKey) : title}
            </h1>
        </div>
        <div class="flex items-center gap-2 flex-wrap">
            {#each filters as f}
                {#if f.type === "select"}
                    <select
                        class="h-10 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-800 dark:text-slate-100 text-sm"
                        bind:value={filterValues[f.key]}
                        on:change={handleFilterChange}
                    >
                        <option value="">
                            {labelFor(f.label, f.labelKey)}
                        </option>
                        {#each f.options ?? [] as opt}
                            <option value={opt.value}>
                                {labelFor(opt.label, opt.labelKey)}
                            </option>
                        {/each}
                    </select>
                {:else if f.type === "date"}
                    <input
                        type="date"
                        class="h-10 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-800 dark:text-slate-100 text-sm"
                        bind:value={filterValues[f.key]}
                        on:change={handleFilterChange}
                    />
                {/if}
            {/each}
            <input
                type="text"
                placeholder={$t("common.search_placeholder")}
                class="h-10 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-800 dark:text-slate-100 text-sm"
                bind:value={search}
                on:input={handleSearchInput}
            />
            <button
                class="px-3 py-2 rounded-lg text-sm font-medium border border-slate-200 dark:border-slate-700 text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800"
                on:click={reloadFromFirstPage}
            >
                {$t("common.refresh")}
            </button>
        </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div
            class="lg:col-span-2 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-2xl overflow-hidden"
        >
            <div
                class="px-5 py-4 border-b border-slate-200 dark:border-slate-700 text-sm font-medium text-slate-700 dark:text-slate-300"
            >
                {$t("common.list")}
            </div>
            <div class="overflow-x-auto">
                {#if loading}
                    <div class="p-6 text-sm text-slate-500">
                        {$t("common.loading")}
                    </div>
                {:else if errorKey}
                    <div class="p-6 text-sm text-red-500">
                        {$t(errorKey)}
                    </div>
                {:else if filteredItems.length === 0}
                    <div class="p-6 text-sm text-slate-500">
                        {$t("common.no_data")}
                    </div>
                {:else}
                    <table class="min-w-full text-sm">
                        <thead
                            class="bg-slate-50 dark:bg-slate-900/40 text-slate-500 dark:text-slate-400"
                        >
                            <tr>
                                {#each columns as col}
                                    <th class="px-4 py-3 text-left font-medium"
                                        >{labelFor(col.label, col.labelKey)}</th
                                    >
                                {/each}
                                <th class="px-4 py-3 text-right font-medium"
                                    >{$t("common.actions")}</th
                                >
                            </tr>
                        </thead>
                        <tbody>
                            {#each filteredItems as item (item.id)}
                                <tr
                                    class="border-t border-slate-200 dark:border-slate-700"
                                >
                                    {#each columns as col}
                                        <td
                                            class="px-4 py-3 text-slate-700 dark:text-slate-200"
                                        >
                                            {#if col.link}
                                                <a
                                                    href={col.link(item)}
                                                    class="text-blue-600 dark:text-blue-400 hover:underline"
                                                >
                                                    {col.formatter
                                                        ? col.formatter(
                                                              item[col.key],
                                                              item,
                                                          )
                                                        : renderValue(
                                                              item[col.key],
                                                              col.key,
                                                          )}
                                                </a>
                                            {:else if col.formatter}
                                                {col.formatter(
                                                    item[col.key],
                                                    item,
                                                )}
                                            {:else}
                                                {renderValue(
                                                    item[col.key],
                                                    col.key,
                                                )}
                                            {/if}
                                        </td>
                                    {/each}
                                    <td
                                        class="px-4 py-3 w-50 flex items-center justify-baseline text-right space-x-2"
                                    >
                                        <button
                                            class="px-4 py-2 text-xs rounded-md border hover:cursor-pointer border-slate-200 dark:border-slate-700 text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700"
                                            on:click={() => startEdit(item)}
                                        >
                                            {$t("common.edit")}
                                        </button>
                                        <button
                                            class="px-3 py-2 text-xs rounded-md border hover:cursor-pointer border-red-200 text-red-600 hover:bg-red-50 dark:border-red-700 dark:text-red-400 dark:hover:bg-red-900/30"
                                            on:click={() => removeItem(item)}
                                        >
                                            {$t("common.delete")}
                                        </button>
                                    </td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                {/if}
            </div>
            {#if totalPages > 1}
                <div
                    class="flex items-center justify-between px-5 py-4 border-t border-slate-200 dark:border-slate-700"
                >
                    <button
                        class="h-9 px-4 rounded-lg border border-slate-200 dark:border-slate-700 text-sm text-slate-600 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-700/40 disabled:opacity-60"
                        on:click={() => setPage(pageNum - 1)}
                        disabled={pageNum <= 1}
                    >
                        {$t("common.prev")}
                    </button>
                    <div
                        class="flex flex-1 items-center justify-center gap-2 px-2"
                    >
                        {#each paginationPages as pageEntry}
                            {#if pageEntry === "ellipsis"}
                                <span
                                    class="h-9 px-3 rounded-lg text-sm text-slate-400 dark:text-slate-500"
                                    aria-hidden="true"
                                >
                                    …
                                </span>
                            {:else}
                                <button
                                    type="button"
                                    class={`h-9 px-3 rounded-lg text-sm font-semibold transition ${
                                        pageEntry === pageNum
                                            ? "bg-blue-600 text-white border-blue-600 shadow-sm"
                                            : "border border-slate-200 dark:border-slate-700 text-slate-600 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-700/40"
                                    }`}
                                    on:click={() => setPage(pageEntry)}
                                    aria-current={pageEntry === pageNum
                                        ? "page"
                                        : undefined}
                                >
                                    {pageEntry}
                                </button>
                            {/if}
                        {/each}
                    </div>
                    <button
                        class="h-9 px-4 rounded-lg border border-slate-200 dark:border-slate-700 text-sm text-slate-600 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-700/40 disabled:opacity-60"
                        on:click={() => setPage(pageNum + 1)}
                        disabled={pageNum >= totalPages}
                    >
                        {$t("common.next")}
                    </button>
                </div>
            {/if}
        </div>

        <div
            class="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-2xl p-5"
        >
            <div
                class="text-sm font-medium text-slate-700 dark:text-slate-300 mb-4"
            >
                {editingId ? $t("common.edit") : $t("common.create")}
            </div>
            <form class="space-y-4" on:submit={submitForm}>
                {#each formFields as field}
                    <div class="space-y-2">
                        <label
                            class="text-xs font-medium text-slate-500 dark:text-slate-400"
                        >
                            {labelFor(field.label, field.labelKey)}
                        </label>

                        {#if field.type === "select"}
                            <select
                                class="w-full h-10 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-800 dark:text-slate-100"
                                bind:value={form[field.key]}
                                required={field.required}
                            >
                                <option value="" disabled selected
                                    >{$t("common.select_placeholder")}</option
                                >
                                {#each selectOptions[field.key] ?? [] as opt}
                                    <option value={opt.value}
                                        >{labelFor(
                                            opt.label,
                                            opt.labelKey,
                                        )}</option
                                    >
                                {/each}
                            </select>
                        {:else if field.type === "checkbox"}
                            <label
                                class="inline-flex items-center gap-2 text-sm text-slate-700 dark:text-slate-200"
                            >
                                <input
                                    type="checkbox"
                                    class="accent-blue-600"
                                    bind:checked={form[field.key]}
                                />
                                <span
                                    >{labelFor(
                                        field.label,
                                        field.labelKey,
                                    )}</span
                                >
                            </label>
                        {:else if field.type === "file"}
                            <input
                                type="file"
                                class="w-full text-sm text-slate-700 dark:text-slate-200"
                                on:change={(e) => {
                                    const target =
                                        e.currentTarget as HTMLInputElement;
                                    form[field.key] = target.files?.[0] ?? null;
                                }}
                            />
                        {:else}
                            <input
                                type={field.type}
                                class="w-full h-10 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-800 dark:text-slate-100"
                                bind:value={form[field.key]}
                                required={field.required}
                            />
                        {/if}
                    </div>
                {/each}

                <div class="flex items-center gap-2 pt-2">
                    <button
                        type="submit"
                        class="px-4 py-2 rounded-lg bg-blue-600 text-white text-sm font-medium hover:bg-blue-700 disabled:opacity-60"
                        disabled={saving}
                    >
                        {saving ? $t("common.saving") : $t("common.save")}
                    </button>
                    {#if editingId}
                        <button
                            type="button"
                            class="px-3 py-2 rounded-lg text-sm border border-slate-200 dark:border-slate-700 text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700"
                            on:click={resetForm}
                        >
                            {$t("common.cancel")}
                        </button>
                    {/if}
                </div>
            </form>
        </div>
    </div>
</section>

{#if deleteConfirmOpen}
    <div class="fixed inset-0 z-50 flex items-center justify-center">
        <div
            class="absolute inset-0 bg-black/50"
            on:click={() => (deleteConfirmOpen = false)}
        ></div>
        <div
            class="relative w-full max-w-md mx-4 rounded-2xl bg-white dark:bg-slate-800
                   border border-slate-200 dark:border-slate-700 shadow-xl p-6"
        >
            <h3
                class="text-lg font-semibold text-slate-800 dark:text-slate-100"
            >
                {$t("common.delete_confirm")}
            </h3>
            <p class="text-sm text-slate-500 dark:text-slate-400 mt-2">
                {$t("common.delete_warning") || "Ushbu ma'lumotni o'chirishni tasdiqlaysizmi?"}
            </p>
            <div class="mt-6 flex items-center gap-2 justify-end">
                <button
                    class="shrink-0 h-10 px-5 rounded-lg dark:text-slate-300 border dark:border-slate-700 border-slate-300 text-slate-600 dark:text-slate-400 text-sm active:scale-[0.98] transition cursor-pointer hover:bg-slate-50 dark:hover:bg-slate-700"
                    on:click={() => (deleteConfirmOpen = false)}
                >
                    {$t("common.cancel")}
                </button>
                <button
                    class="shrink-0 h-10 px-5 rounded-lg bg-red-600 hover:bg-red-700 text-white text-sm active:scale-[0.98] transition cursor-pointer"
                    on:click={handleConfirmDelete}
                >
                    {$t("common.confirm") || "Tasdiqlash"}
                </button>
            </div>
        </div>
    </div>
{/if}
