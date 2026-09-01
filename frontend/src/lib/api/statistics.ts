import { apiRequest } from "./client";
import type { TodayStatistics } from "$lib/types";

export const statisticsApi = {
  getToday(date?: string) {
    const query = date ? `?date=${date}` : "";
    return apiRequest<TodayStatistics>(`statistics/today${query}`);
  },
};
