import { Settings } from "../types/settings";
import api from "../utils/api";

export const getSettings = async () => {
  const response = await api.get<Settings>("/settings", {
    headers: {
      "Content-Type": "application/json",
    },
  });
  if (response.status < 200 || response.status >= 300) {
    throw new Error("Failed to get settings");
  }
  return response.data;
};

export const updateSettings = async (
  settings: Partial<Settings>
): Promise<any> => {
  const response = await api.put("/settings/update", settings, {
    headers: {
      "Content-Type": "application/json",
    },
  });
  if (response.status < 200 || response.status >= 300) {
    throw new Error("Failed to update settings");
  }
  return response.data;
};
