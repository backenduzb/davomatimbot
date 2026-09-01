import { apiRequest, fetchAllResults } from "./client";
import type { PaginatedResponse, User } from "$lib/types";

export const usersApi = {
  getAll() {
    return fetchAllResults<User>("users");
  },

  getPage(page = 1, pageSize = 30) {
    return apiRequest<PaginatedResponse<User>>(
      `users?page=${page}&page_size=${pageSize}`,
    );
  },

  update(id: number, data: Partial<User>) {
    return apiRequest<User>(`users/${id}`, {
      method: "PATCH",
      body: JSON.stringify(data),
    });
  },

  delete(id: number) {
    return apiRequest<void>(`users/${id}`, { method: "DELETE" });
  },
};
