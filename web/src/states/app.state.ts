import { create } from "zustand";
import { persist } from "zustand/middleware";

interface AppState {
  investor: boolean;
  activeTab: number;
  settingsActiveTab: string;
  setActiveTab: (tab: number) => void;
  setSettingsActiveTab: (tab: string) => void;
  setInvestor: (investor: boolean) => void;
}

export const useAppState = create<AppState>()(
  persist(
    (set) => ({
      activeTab: 0,
      settingsActiveTab: "workspace",
      investor: false,
      theme: "dark",
      primaryColor: "dark",
      setActiveTab: (activeTab: number) => set({ activeTab }),
      setSettingsActiveTab: (settingsActiveTab: string) =>
        set({ settingsActiveTab }),
      setInvestor: (investor: boolean) => set({ investor }),
    }),
    {
      name: "__APP_STATE__",
      partialize: (state) => ({
        activeTab: state.activeTab,
        settingsActiveTab: state.settingsActiveTab,
      }),
    }
  )
);
