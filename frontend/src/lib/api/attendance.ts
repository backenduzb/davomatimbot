import { apiRequest } from "./client";
import type {
  AttendanceStatus,
  TodayAttendanceResponse,
} from "$lib/types";

export const attendanceApi = {
  getToday(date?: string, classId?: number) {
    const params = new URLSearchParams();
    if (date) params.set("date", date);
    if (classId) params.set("class_id", String(classId));
    const query = params.toString();
    return apiRequest<TodayAttendanceResponse>(
      `attendance/today${query ? `?${query}` : ""}`,
    );
  },

  saveBatch(input: {
    date: string;
    class_id: number;
    records: { student_id: number; status: AttendanceStatus; reason?: string }[];
  }) {
    return apiRequest<{ message: string }>("attendance/batch", {
      method: "POST",
      body: JSON.stringify(input),
    });
  },
};
