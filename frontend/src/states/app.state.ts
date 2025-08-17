import { create } from "zustand";

interface AppState {
  activeTab: number;
  setActiveTab: (tab: number) => void;
}

export const useAppState = create<AppState>((set) => ({
  activeTab: 0,
  setActiveTab: (activeTab: number) => set({ activeTab }),
}));
