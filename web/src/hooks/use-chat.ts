import { useQuery } from "@tanstack/react-query";
import { getFiles } from "../api/chat.api";
import { ProjectFile } from "../types/chat";

export function useChatFiles(
  dir: string | null,
  search: string | null,
  enabled: boolean
) {
  return useQuery<ProjectFile[], Error>({
    queryKey: ["chat", "files", dir, search],
    queryFn: () => getFiles(dir, search),
    enabled,
  });
}
