import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { History, TrendData, RecentChanges, Compare } from "../types/stat";
import {
  loadHistory,
  loadTrends,
  loadChanges,
  loadComparison,
  refreshStats,
  cleanupStats,
} from "../api/history.api";

export function useHistory(enabled: boolean) {
  return useQuery<History, Error>({
    queryKey: ["history", "history"],
    queryFn: loadHistory,
    enabled,
  });
}

export function useTrends(enabled: boolean) {
  return useQuery<TrendData, Error>({
    queryKey: ["history", "trends"],
    queryFn: loadTrends,
    enabled,
  });
}

export function useChanges(enabled: boolean) {
  return useQuery<RecentChanges, Error>({
    queryKey: ["history", "changes"],
    queryFn: loadChanges,
    enabled,
  });
}

export function useComparison(enabled: boolean) {
  return useQuery<Compare, Error>({
    queryKey: ["history", "comparison"],
    queryFn: loadComparison,
    enabled,
  });
}

export function useRefreshStats() {
  const queryClient = useQueryClient();
  return useMutation<{ success: boolean }, Error>({
    mutationFn: refreshStats,

    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["history"] });
    },
  });
}

export function useCleanupStats() {
  const queryClient = useQueryClient();
  return useMutation<{ success: boolean }, Error>({
    mutationFn: cleanupStats,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["history"] });
    },
  });
}
