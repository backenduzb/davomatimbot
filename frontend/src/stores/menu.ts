import { writable } from "svelte/store";
import { browser } from "$app/environment";

const getInitialState = (): boolean => {
  if (!browser) return false;
  return window.innerWidth >= 768;
};

export const menuOpen = writable<boolean>(getInitialState());

export function toggleMenu() {
  menuOpen.update((v) => !v);
}

export function closeMenu() {
  menuOpen.set(false);
}

export function openMenu() {
  menuOpen.set(true);
}
