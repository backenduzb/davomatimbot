import { apiRequest, fetchAllResults } from "./client";
import type { PaginatedResponse, Student } from "$lib/types";

export const studentsApi = {
  getAll(classId?: number) {
    const query = classId ? `?class_id=${classId}` : "";
    return fetchAllResults<Student>(`students${query}`);
  },

  getPage(page = 1, pageSize = 30) {
    return apiRequest<PaginatedResponse<Student>>(
      `students?page=${page}&page_size=${pageSize}`,
    );
  },

  create(data: Omit<Student, "id">) {
    return apiRequest<Student>("students", {
      method: "POST",
      body: JSON.stringify(data),
    });
  },

  update(id: number, data: Partial<Student>) {
    return apiRequest<Student>(`students/${id}`, {
      method: "PATCH",
      body: JSON.stringify(data),
    });
  },

  delete(id: number) {
    return apiRequest<void>(`students/${id}`, { method: "DELETE" });
  },
};
