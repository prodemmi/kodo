import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { History, TrendData, RecentChanges, Compare } from "../types/stat";
import {
  loadHistory,
  loadTrends,
  loadChanges,
  loadComparison,
  refreshStats,
  cleanupStats,
} from "../api/stats.api";

export function useHistory(enabled: boolean) {
  return useQuery<History, Error>({
    queryKey: ["stats", "history"],
    queryFn: loadHistory,
    enabled,
  });
}

export function useTrends(enabled: boolean) {
  return useQuery<TrendData, Error>({
    queryKey: ["stats", "trends"],
    queryFn: loadTrends,
    enabled,
  });
}

export function useChanges(enabled: boolean) {
  return useQuery<RecentChanges, Error>({
    queryKey: ["stats", "changes"],
    queryFn: loadChanges,
    enabled,
  });
}

export function useComparison(enabled: boolean) {
  return useQuery<Compare, Error>({
    queryKey: ["stats", "comparison"],
    queryFn: loadComparison,
    enabled,
  });
}

export function useRefreshStats() {
  const queryClient = useQueryClient();
  return useMutation<{ success: boolean }, Error>({
    mutationFn: refreshStats,

    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["stats"] });
    },
  });
}

export function useCleanupStats() {
  const queryClient = useQueryClient();
  return useMutation<{ success: boolean }, Error>({
    mutationFn: cleanupStats,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["stats"] });
    },
  });
}
