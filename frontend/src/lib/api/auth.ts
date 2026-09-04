import { API_BASE_URL, apiRequest } from "./client";
import type { LoginResponse, User } from "$lib/types";

export async function login(username: string, password: string) {
  const res = await fetch(`${API_BASE_URL}/api/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });

  if (!res.ok) {
    throw new Error("login_failed");
  }

  return (await res.json()) as LoginResponse;
}

export async function getCurrentUser() {
  const token = localStorage.getItem("token");
  const res = await fetch(`${API_BASE_URL}/api/me`, {
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
  });

  if (!res.ok) {
    throw new Error("unauthorized");
  }

  return (await res.json()) as User;
}

export interface ChangePasswordPayload {
  current_password: string;
  new_password: string;
}

/**
 * Foydalanuvchi o'z parolini almashtiradi.
 * Backend joriy parolni tekshiradi va xato bo'lsa ApiClientError qaytaradi
 * (message — backenddagi o'zbekcha matn).
 */
export async function changePassword(payload: ChangePasswordPayload) {
  return apiRequest<{ message: string }>("me/password", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}
