import { create } from "zustand";
import { persist } from "zustand/middleware";

interface AppState {
  investor: boolean;
  activeTab: number;
  primaryColor: string;
  setActiveTab: (tab: number) => void;
  setInvestor: (investor: boolean) => void;
  setPrimaryColor: (primaryColor: string) => void;
}

export const useAppState = create<AppState>()(
  persist(
    (set) => ({
      activeTab: 0,
      investor: false,
      primaryColor: "dark",
      setActiveTab: (activeTab: number) => set({ activeTab }),
      setInvestor: (investor: boolean) => set({ investor }),
      setPrimaryColor: (primaryColor: string) => set({ primaryColor }),
    }),
    {
      name: "__APP_STATE__",
      partialize: (state) => ({ activeTab: state.activeTab, primaryColor: state.primaryColor }),
    }
  )
);
