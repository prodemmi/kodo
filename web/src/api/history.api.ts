import { History, TrendData, RecentChanges, Compare } from "../types/stat";
import api from "../utils/api";

export const loadHistory = async (): Promise<History> => {
  const response = await api.get<History>(`history/history`);
  return response.data;
};

export const loadTrends = async (): Promise<TrendData> => {
  const response = await api.get<TrendData>(`history/trends`);
  return response.data;
};

export const loadChanges = async (): Promise<RecentChanges> => {
  const response = await api.get<RecentChanges>(`history/changes`);
  return response.data;
};

export const loadComparison = async (): Promise<Compare> => {
  const response = await api.get<Compare>(`history/compare`);
  return response.data;
};

export const refreshStats = async (): Promise<{ success: boolean }> => {
  const response = await api.post<{ success: boolean }>(`history`);
  return response.data;
};

export const cleanupStats = async (): Promise<{ success: boolean }> => {
  const response = await api.post<{ success: boolean }>(`history/cleanup`);
  return response.data;
};
