import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { getSettings, updateSettings } from "../api/settings.api";
import { Settings } from "../types/settings";

/**
 * Hook to update settings in the backend using React Query.
 * Automatically calls the mutation and updates Zustand on success.
 */
export function useSettings() {
  return useQuery<Settings, Error>({
    queryKey: ["settings"],
    queryFn: getSettings,
  });
}

/**
 * Hook to update settings in the backend using React Query.
 * Automatically calls the mutation and updates Zustand on success.
 */
export function useUpdateSettings() {
  const client = useQueryClient();
  const mutation = useMutation({
    mutationFn: updateSettings,
    onSuccess: () => {
      client.invalidateQueries({ queryKey: ["settings"] });
    },
  });

  // Returns a function to call mutation; accepts partial settings object
  return (settings: Partial<Settings>) => {
    mutation.mutate(settings);
  };
}
