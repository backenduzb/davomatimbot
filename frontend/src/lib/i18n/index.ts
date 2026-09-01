import { writable, derived } from "svelte/store";
import { browser } from "$app/environment";
import uz from "./locales/uz";
import en from "./locales/en";
import ru from "./locales/ru";
import kaa from "./locales/kaa";

export type Locale = "uz" | "en" | "ru" | "kaa";

const translations = { uz, en, ru, kaa };

const stored = browser ? (localStorage.getItem("locale") as Locale) : null;
export const locale = writable<Locale>(stored ?? "uz");

locale.subscribe((val) => {
  if (!browser) return;
  localStorage.setItem("locale", val);
});

export const t = derived(locale, ($locale) => {
  return (key: string): string => {
    const keys = key.split(".");
    let result: any = translations[$locale];
    for (const k of keys) {
      result = result?.[k];
    }
    return result ?? key;
  };
});
