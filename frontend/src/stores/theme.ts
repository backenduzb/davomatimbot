import { writable } from "svelte/store";
import { browser } from "$app/environment";

const stored = browser ? localStorage.getItem("theme") : null;
const initialValue = stored === "dark";

export const isDark = writable<boolean>(initialValue);

isDark.subscribe((val) => {
  if (!browser) return;

  localStorage.setItem("theme", val ? "dark" : "light");

  if (val) {
    document.documentElement.classList.add("dark");
  } else {
    document.documentElement.classList.remove("dark");
  }
});

export function toggleTheme() {
  isDark.update((val) => !val);
}


