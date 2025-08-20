import { create } from "zustand";

interface AppState {
  investor: boolean;
  activeTab: number;
  setActiveTab: (tab: number) => void;
}

export const useAppState = create<AppState>((set) => ({
  activeTab: 0,
  investor: true,
  setActiveTab: (activeTab: number) => set({ activeTab }),
}));
