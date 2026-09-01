import { apiRequest, fetchAllResults } from "./client";
import type { Class, ClassName, PaginatedResponse } from "$lib/types";

export const classesApi = {
  getAll() {
    return fetchAllResults<Class>("classes");
  },

  getPage(page = 1, pageSize = 30) {
    return apiRequest<PaginatedResponse<Class>>(
      `classes?page=${page}&page_size=${pageSize}`,
    );
  },

  create(data: Omit<Class, "id">) {
    return apiRequest<Class>("classes", {
      method: "POST",
      body: JSON.stringify(data),
    });
  },

  update(id: number, data: Partial<Class>) {
    return apiRequest<Class>(`classes/${id}`, {
      method: "PATCH",
      body: JSON.stringify(data),
    });
  },

  delete(id: number) {
    return apiRequest<void>(`classes/${id}`, { method: "DELETE" });
  },
};

export const classNamesApi = {
  getAll() {
    return fetchAllResults<ClassName>("class/names");
  },

  create(data: Omit<ClassName, "id">) {
    return apiRequest<ClassName>("class/names", {
      method: "POST",
      body: JSON.stringify(data),
    });
  },

  update(id: number, data: Partial<ClassName>) {
    return apiRequest<ClassName>(`class/names/${id}`, {
      method: "PATCH",
      body: JSON.stringify(data),
    });
  },

  delete(id: number) {
    return apiRequest<void>(`class/names/${id}`, { method: "DELETE" });
  },
};
