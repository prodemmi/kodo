import { create } from "zustand";
import { persist } from "zustand/middleware";

interface AppState {
  investor: boolean;
  activeTab: number;
  setActiveTab: (tab: number) => void;
  setInvestor: (investor: boolean) => void;
}

export const useAppState = create<AppState>()(
  persist(
    (set) => ({
      activeTab: 0,
      investor: false,
      setActiveTab: (activeTab: number) => set({ activeTab }),
      setInvestor: (investor: boolean) => set({ investor }),
    }),
    {
      name: "__APP_STATE__",
      partialize: (state) => ({ activeTab: state.activeTab }),
    }
  )
);
