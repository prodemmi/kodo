import { createTheme as theme, Button, rem } from "@mantine/core";

export const createTheme = (primaryColor: string) => {
  return theme({
    // Primary colors - Using different shades for light/dark
    primaryColor,
    primaryShade: { light: 6, dark: 8 },

    // Custom color palette - organized by relationship
    colors: {
      // Main dark palette - foundation colors for dark mode
      dark: [
        "#f0f6ff", // lightest - for light text on dark
        "#c9d1d9", // very light gray - primary text
        "#b1bac4", // light gray - secondary text
        "#8b949e", // medium gray - muted text
        "#6e7681", // text secondary - disabled text
        "#484f58", // border color - subtle borders
        "#30363d", // surface color - elevated surfaces
        "#21262d", // card background - main content areas
        "#161b22", // page background - main app background
        "#0d1117", // darkest - deepest backgrounds
      ],

      // Light palette - foundation colors for light mode
      light: [
        "#fafbfc", // lightest - subtle backgrounds
        "#f8f9fa", // page background - main app background
        "#ffffff", // card background - main content areas
        "#f6f8fa", // surface color - elevated surfaces
        "#d0d7de", // border color - subtle borders
        "#8c959f", // text secondary - disabled text
        "#656d76", // medium gray - muted text
        "#32383f", // dark gray - secondary text
        "#24292f", // very dark gray - primary text
        "#0d1117", // darkest - for dark text on light
      ],

      // Status colors - grouped for workflow states
      blue: [
        "#dbeafe", // lightest
        "#bfdbfe",
        "#93c5fd",
        "#60a5fa",
        "#3b82f6", // medium - good for icons
        "#2563eb", // primary shade light
        "#1d4ed8", // good for buttons
        "#1e40af", // primary shade dark
        "#1e3a8a", // darker
        "#172554", // darkest
      ],

      orange: [
        "#fed7aa", // lightest
        "#fdba74",
        "#fb923c",
        "#f97316", // good for warnings
        "#ea580c", // medium - progress indicators
        "#dc2626",
        "#c2410c", // darker
        "#9a3412",
        "#7c2d12",
        "#431407", // darkest
      ],

      green: [
        "#dcfce7", // lightest
        "#bbf7d0",
        "#86efac",
        "#4ade80", // bright success
        "#22c55e", // medium success
        "#16a34a", // good for completed states
        "#15803d", // darker success
        "#166534",
        "#14532d",
        "#052e16", // darkest
      ],

      red: [
        "#fecaca", // lightest
        "#fca5a5",
        "#f87171",
        "#ef4444", // good for errors
        "#dc2626", // medium error
        "#b91c1c", // darker error
        "#991b1b",
        "#7f1d1d",
        "#6b2c2c",
        "#450a0a", // darkest
      ],

      // Neutral grays - for subtle variations
      gray: [
        "#f8fafc",
        "#f1f5f9",
        "#e2e8f0",
        "#cbd5e1",
        "#94a3b8",
        "#64748b",
        "#475569",
        "#334155",
        "#1e293b",
        "#0f172a",
      ],
    },

    // Background colors - using consistent theme colors
    white: "#ffffff",
    black: "#0d1117",

    // Typography - GitHub font stack
    fontFamily:
      '-apple-system, BlinkMacSystemFont, "Segoe UI", "Noto Sans", Helvetica, Arial, sans-serif',
    fontFamilyMonospace:
      'ui-monospace, SFMono-Regular, "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Consolas, "Courier New", monospace',

    headings: {
      fontFamily:
        '-apple-system, BlinkMacSystemFont, "Segoe UI", "Noto Sans", Helvetica, Arial, sans-serif',
      fontWeight: "600",
      sizes: {
        h1: { fontSize: rem(32), lineHeight: "1.3" },
        h2: { fontSize: rem(24), lineHeight: "1.35" },
        h3: { fontSize: rem(20), lineHeight: "1.4" },
        h4: { fontSize: rem(18), lineHeight: "1.45" },
        h5: { fontSize: rem(16), lineHeight: "1.5" },
        h6: { fontSize: rem(14), lineHeight: "1.5" },
      },
    },

    // Design tokens
    defaultRadius: "md",

    spacing: {
      xs: rem(8),
      sm: rem(12),
      md: rem(16),
      lg: rem(24),
      xl: rem(32),
    },

    shadows: {
      xs: "0 1px 2px rgba(0, 0, 0, 0.05)",
      sm: "0 1px 3px rgba(0, 0, 0, 0.12), 0 1px 2px rgba(0, 0, 0, 0.24)",
      md: "0 4px 6px rgba(0, 0, 0, 0.07), 0 2px 4px rgba(0, 0, 0, 0.06)",
      lg: "0 10px 15px rgba(0, 0, 0, 0.1), 0 4px 6px rgba(0, 0, 0, 0.05)",
      xl: "0 20px 25px rgba(0, 0, 0, 0.15), 0 10px 10px rgba(0, 0, 0, 0.04)",
    },

    // Component extensions using colorSchema variable
    components: {
      Button: Button.extend({
        defaultProps: {
          size: "xs",
        },
      }),
    },
  });
};
