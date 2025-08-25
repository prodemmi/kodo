import { create } from "zustand";
import { persist } from "zustand/middleware";
import { WorkspaceSettings } from "../types/settings";

export interface SettingsState {
  workspace_settings: WorkspaceSettings;

  setWorkspaceSettings: (settings: Partial<WorkspaceSettings>) => void;
}

export const useSettingsState = create<SettingsState>()(
  persist(
    (set, get) => ({
      workspace_settings: {
        theme: "dark",
        primary_color: "dark",
        show_line_preview: true,
      },
      setWorkspaceSettings: (settings: Partial<WorkspaceSettings>) =>
        set({
          workspace_settings: { ...get().workspace_settings, ...settings },
        }),
    }),
    {
      name: "__APP_SETTINGS_STATE__",
      partialize: (state) => ({
        workspace_settings: state.workspace_settings,
      }),
    }
  )
);
