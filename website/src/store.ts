import { create } from "zustand";

export type Locale = "en" | "zh" | "ja" | "ko" | "es" | "fr" | "de";
export type CatalogFilter = "all" | "components" | "platforms";

type UiState = {
  locale: Locale;
  filter: CatalogFilter;
  mobileMenu: boolean;
  setLocale: (locale: Locale) => void;
  setFilter: (filter: CatalogFilter) => void;
  toggleMobileMenu: () => void;
};

export const useUiStore = create<UiState>((set) => ({
  locale: "zh",
  filter: "all",
  mobileMenu: false,
  setLocale: (locale) => set({ locale }),
  setFilter: (filter) => set({ filter }),
  toggleMobileMenu: () => set((state) => ({ mobileMenu: !state.mobileMenu })),
}));
