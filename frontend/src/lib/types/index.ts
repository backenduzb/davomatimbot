export interface PaginatedResponse<T> {
  count: number;
  total_pages: number;
  current_page: number;
  page_size: number;
  results: T[];
}

export interface User {
  id: number;
  username: string;
  created_at: string;
  is_online: boolean;
  is_banned: boolean;
  is_admin: boolean;
  telegram_id: string;
  last_seen: string;
}

export interface ClassName {
  id: number;
  name: string;
}

export interface Class {
  id: number;
  updated: boolean;
  class_name_id: number;
  teacher_full_name: string;
  teacher_telegram_id: string;
}

export interface Student {
  id: number;
  full_name: string;
  class_id: number;
}

export interface AttendanceRecord {
  id: number;
  student_id: number;
  student_name: string;
  class_id: number;
  date: string;
  status: AttendanceStatus;
  reason?: string;
}

export type AttendanceStatus =
  | "present"
  | "absent"
  | "excused"
  | "late"
  | "not_marked";

export interface TodayAttendanceResponse {
  date: string;
  class_id?: number;
  records: AttendanceRecord[];
}

export interface ClassStats {
  class_id: number;
  class_name: string;
  total_students: number;
  present: number;
  absent: number;
  excused: number;
  late: number;
  not_marked: number;
  attendance_percent: number;
}

export interface TodayStatistics {
  date: string;
  total_classes: number;
  total_students: number;
  present: number;
  absent: number;
  excused: number;
  late: number;
  not_marked: number;
  attendance_percent: number;
  classes: ClassStats[];
}

export interface LoginResponse {
  token: string;
}

export interface ApiError {
  error: string;
}

export interface ImportResult {
  rows_processed: number;
  students_created: number;
  students_linked: number;
  classes_created: number;
  class_names_created: number;
}

export interface ImportResponse {
  message: string;
  result: ImportResult;
}
