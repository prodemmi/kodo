import { History, TrendData, RecentChanges, Compare } from "../types/stat";

export const loadHistory = async (): Promise<History> => {
  const response = await fetch("http://localhost:8080/api/stats/history");
  if (!response.ok) throw new Error("Failed to load history");
  return await response.json();
};

export const loadTrends = async (): Promise<TrendData> => {
  const response = await fetch("http://localhost:8080/api/stats/trends");
  if (!response.ok) throw new Error("Failed to load trends");
  return await response.json();
};

export const loadChanges = async (): Promise<RecentChanges> => {
  const response = await fetch("http://localhost:8080/api/stats/changes");
  if (!response.ok) throw new Error("Failed to load changes");
  return await response.json();
};

export const loadComparison = async (): Promise<Compare> => {
  const response = await fetch("http://localhost:8080/api/stats/compare");
  if (!response.ok) throw new Error("Failed to load comparison");
  return await response.json();
};

export const refreshStats = async (): Promise<{ success: boolean }> => {
  const response = await fetch("http://localhost:8080/api/stats", {
    method: "POST",
  });
  if (!response.ok) throw new Error("Failed to refresh stats");
  return await response.json();
};

export const cleanupStats = async (): Promise<{ success: boolean }> => {
  const response = await fetch("http://localhost:8080/api/stats/cleanup", {
    method: "POST",
  });
  if (!response.ok) throw new Error("Failed to cleanup stats");
  return await response.json();
};
