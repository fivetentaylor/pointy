import { Signal, computed, effect, signal } from "@preact/signals-react";

type Theme = "light" | "dark";
export type ThemePreference = Theme | "system";

const LOCAL_STORAGE_THEME_PREFERENCE_KEY = "THEME_PREFERENCE";

class ColorThemeService {
  mediaTheme: Signal<Theme> | null;
  themePreference: Signal<ThemePreference> = signal(
    (localStorage.getItem(
      LOCAL_STORAGE_THEME_PREFERENCE_KEY,
    ) as ThemePreference) || "system",
  );
  theme: Signal<Theme> = computed(() => {
    if (this.themePreference.value === "system") {
      return this.mediaTheme ? this.mediaTheme.value : "light";
    }
    return this.themePreference.value;
  });

  constructor() {
    const darkThemeMq = window.matchMedia("(prefers-color-scheme: dark)");
    effect(() => {
      this.onThemeChange(this.theme.value);
    });

    this.mediaTheme = signal(darkThemeMq.matches ? "dark" : "light");

    darkThemeMq.addEventListener("change", this.onMediaChange.bind(this));
  }

  onMediaChange(event: MediaQueryListEvent) {
    if (!this.mediaTheme) {
      return;
    }

    this.mediaTheme.value = event.matches ? "dark" : "light";
  }

  onThemeChange(theme: string) {
    const html = document.querySelector("html");
    if (!html) {
      return;
    }
    if (theme === "dark") {
      html.classList.add("dark");
      html.classList.remove("light");
      html.style.colorScheme = "dark";
    } else {
      html.classList.add("light");
      html.classList.remove("dark");
      html.style.colorScheme = "light";
    }
  }

  setPreferredTheme(themePreference: ThemePreference) {
    localStorage.setItem(LOCAL_STORAGE_THEME_PREFERENCE_KEY, themePreference);
    this.themePreference.value = themePreference;
  }
}

export const colorThemeService = new ColorThemeService();
