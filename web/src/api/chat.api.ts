import { ProjectFile } from "../types/chat";

export const getFiles = async (
  dir: string | null,
  search: string | null
): Promise<ProjectFile[]> => {
  const params = new URLSearchParams();
  if (dir) params.set("dir", dir);
  if (search) params.set("search", search);

  const response = await fetch(
    `http://localhost:8080/api/chat/project-files?${params.toString()}`
  );

  if (!response.ok) throw new Error("Failed to load project-files");
  return await response.json();
};
