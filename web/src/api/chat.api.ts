import { ProjectFile } from "../types/chat";
import api from "../utils/api";

export const getFiles = async (
  dir: string | null,
  search: string | null
): Promise<ProjectFile[]> => {
  try {
    const params: Record<string, string> = {};
    if (dir) params.dir = dir;
    if (search) params.search = search;

    const response = await api.get<ProjectFile[]>("/chat/project-files", {
      params,
    });

    return response.data;
  } catch (error: any) {
    // Axios interceptor already shows notification
    throw new Error(error.response?.data?.message || "Failed to load project files");
  }
};
