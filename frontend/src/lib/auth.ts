import { user, isAuthenticated, loading } from "../stores/auth";
import { getCurrentUser } from "$lib/api/auth";
import { goto } from "$app/navigation";
import { browser } from "$app/environment";

export async function checkAuth() {
  loading.set(true);

  if (!browser || !localStorage.getItem("token")) {
    isAuthenticated.set(false);
    user.set(null);
    loading.set(false);
    return;
  }

  try {
    const data = await getCurrentUser();
    user.set(data);
    isAuthenticated.set(true);
  } catch {
    localStorage.removeItem("token");
    user.set(null);
    isAuthenticated.set(false);
  } finally {
    loading.set(false);
  }
}

export function logout() {
  if (browser) {
    localStorage.removeItem("token");
  }
  user.set(null);
  isAuthenticated.set(false);
  goto("/authenticate");
}
