import { goto } from "$app/navigation";
import { browser } from "$app/environment";
import type { PaginatedResponse } from "$lib/types";

export const API_BASE_URL = (
  import.meta.env.PUBLIC_API_URL ?? "http://34.134.26.219:8000"
).replace(/\/$/, "");

const API_URL = `${API_BASE_URL}/api`;

export class ApiClientError extends Error {
  status: number;

  constructor(message: string, status: number) {
    super(message);
    this.status = status;
    this.name = "ApiClientError";
  }
}

function getToken(): string | null {
  if (!browser) return null;
  return localStorage.getItem("token");
}

function clearAuth() {
  if (!browser) return;
  localStorage.removeItem("token");
}

function handleUnauthorized() {
  clearAuth();
  if (!browser) return;
  const path = window.location.pathname + window.location.search;
  if (!path.startsWith("/authenticate")) {
    goto(`/authenticate?redirectTo=${encodeURIComponent(path)}`);
  }
}

export async function apiFetch(
  path: string,
  options: RequestInit = {},
): Promise<Response> {
  const token = getToken();
  const normalizedPath = path.replace(/^\/+/, "");

  return fetch(`${API_URL}/${normalizedPath}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...options.headers,
    },
  });
}

export async function apiRequest<T>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  const res = await apiFetch(path, options);

  if (res.status === 401) {
    handleUnauthorized();
    throw new ApiClientError("Unauthorized", 401);
  }

  if (!res.ok) {
    let message = "Request failed";
    try {
      const data = await res.json();
      message = data?.error ?? message;
    } catch {
      // ignore parse errors
    }
    throw new ApiClientError(message, res.status);
  }

  if (res.status === 204) {
    return undefined as T;
  }

  return (await res.json()) as T;
}

export async function fetchAllResults<T>(
  path: string,
  pageSize = 100,
): Promise<T[]> {
  const items: T[] = [];
  let page = 1;
  let totalPages = 1;

  while (page <= totalPages) {
    const separator = path.includes("?") ? "&" : "?";
    const res = await apiFetch(
      `${path}${separator}page=${page}&page_size=${pageSize}`,
    );

    if (res.status === 401) {
      handleUnauthorized();
      throw new ApiClientError("Unauthorized", 401);
    }

    if (!res.ok) {
      throw new ApiClientError("load_failed", res.status);
    }

    const data = (await res.json()) as PaginatedResponse<T> | T[];
    if (Array.isArray(data)) {
      return data;
    }

    items.push(...(data.results ?? []));
    totalPages = data.total_pages || 1;
    page++;
  }

  return items;
}
