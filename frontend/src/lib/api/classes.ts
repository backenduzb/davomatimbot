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

// --- Sinflarni keyingi o'quv yiliga oshirish ---

export type PromotionAction = "promote" | "graduate" | "skip";

export type PromotionPlan = {
  class_id: number;
  current_name: string;
  next_name: string;
  action: PromotionAction;
  student_count: number;
  reason?: string;
};

export type PromotionPreview = {
  plans: PromotionPlan[];
  promote: number;
  graduate: number;
  skip: number;
  total: number;
};

export type PromotionResult = {
  promoted: number;
  graduated: number;
  skipped: number;
  class_names_created: number;
};

export const classPromotionApi = {
  // preview — hech narsani o'zgartirmaydi, faqat rejani qaytaradi.
  preview() {
    return apiRequest<PromotionPreview>("class-promotion/preview");
  },

  // promote — tasdiqlangandan keyin amalni bajaradi.
  promote() {
    return apiRequest<PromotionResult>("class-promotion", {
      method: "POST",
      body: JSON.stringify({ confirm: true }),
    });
  },
};
