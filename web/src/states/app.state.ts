import { create } from "zustand";
import { persist } from "zustand/middleware";

interface AppState {
  investor: boolean;
  activeTab: number;
  theme: "light" | "dark" | "auto";
  primaryColor: string;
  setActiveTab: (tab: number) => void;
  setInvestor: (investor: boolean) => void;
  setPrimaryColor: (primaryColor: string) => void;
  setTheme: (theme: "light" | "dark" | "auto") => void;
}

export const useAppState = create<AppState>()(
  persist(
    (set) => ({
      activeTab: 0,
      investor: false,
      theme: "dark",
      primaryColor: "dark",
      setActiveTab: (activeTab: number) => set({ activeTab }),
      setInvestor: (investor: boolean) => set({ investor }),
      setPrimaryColor: (primaryColor: string) => set({ primaryColor }),
      setTheme: (theme: "light" | "dark" | "auto") => set({ theme }),
    }),
    {
      name: "__APP_STATE__",
      partialize: (state) => ({
        theme: state.theme,
        activeTab: state.activeTab,
        primaryColor: state.primaryColor,
      }),
    }
  )
);
