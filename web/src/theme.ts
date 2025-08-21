import { createTheme } from "@mantine/core";

export const theme = createTheme({
  // Primary colors - Using dark for consistency
  primaryColor: "dark",
  primaryShade: { light: 6, dark: 8 },

  // Custom color palette - organized by relationship
  colors: {
    // Main dark palette - foundation colors
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

  // Background colors - using consistent dark theme
  white: "#21262d",
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
      h1: { fontSize: "2rem", lineHeight: "1.3" },
      h2: { fontSize: "1.5rem", lineHeight: "1.35" },
      h3: { fontSize: "1.25rem", lineHeight: "1.4" },
      h4: { fontSize: "1.125rem", lineHeight: "1.45" },
      h5: { fontSize: "1rem", lineHeight: "1.5" },
      h6: { fontSize: "0.875rem", lineHeight: "1.5" },
    },
  },

  // Component customizations - organized by component type
  components: {
    // Layout components
    AppShell: {
      styles: {
        root: {
          backgroundColor: "var(--mantine-color-dark-9)",
        },
        header: {
          backgroundColor: "var(--mantine-color-dark-8)",
          borderBottom: "1px solid var(--mantine-color-dark-6)",
        },
        navbar: {
          backgroundColor: "var(--mantine-color-dark-8)",
          borderRight: "1px solid var(--mantine-color-dark-6)",
        },
        aside: {
          backgroundColor: "var(--mantine-color-dark-8)",
          borderLeft: "1px solid var(--mantine-color-dark-6)",
        },
        footer: {
          backgroundColor: "var(--mantine-color-dark-8)",
          borderTop: "1px solid var(--mantine-color-dark-6)",
        },
        main: {
          backgroundColor: "var(--mantine-color-dark-9)",
        },
      },
    },

    // Surface components
    Card: {
      styles: {
        root: {
          backgroundColor: "var(--mantine-color-dark-7)",
          border: "1px solid var(--mantine-color-dark-6)",
          borderRadius: "var(--mantine-radius-md)",
          boxShadow: "var(--mantine-shadow-sm)",
          transition: "all 200ms ease",
          "&:hover": {
            backgroundColor: "var(--mantine-color-dark-6)",
            borderColor: "var(--mantine-color-dark-5)",
            transform: "translateY(-1px)",
            boxShadow: "var(--mantine-shadow-md)",
          },
        },
      },
    },

    Paper: {
      styles: {
        root: {
          backgroundColor: "var(--mantine-color-dark-7)",
          border: "1px solid var(--mantine-color-dark-6)",
        },
      },
    },

    // Interactive components
    Button: {
      styles: {
        root: {
          fontWeight: 500,
          transition: "all 200ms ease",
          backgroundColor: "var(--mantine-color-dark-7)",
          borderColor: "var(--mantine-color-dark-6)",
          color: "var(--mantine-color-dark-0)",
          "&:hover": {
            transform: "translateY(-1px)",
          },
        },
      },
    },

    ActionIcon: {
      styles: {
        root: {
          color: "var(--mantine-color-dark-2)",
          transition: "all 200ms ease",
          "&:hover": {
            backgroundColor: "var(--mantine-color-dark-6)",
            color: "var(--mantine-color-dark-0)",
            transform: "scale(1.05)",
          },
        },
      },
    },

    // Form components
    Input: {
      styles: {
        input: {
          backgroundColor: "var(--mantine-color-dark-8)",
          borderColor: "var(--mantine-color-dark-6)",
          color: "var(--mantine-color-dark-0)",
          transition: "all 200ms ease",
          "&::placeholder": {
            color: "var(--mantine-color-dark-3)",
          },
          "&:focus": {
            borderColor: "var(--mantine-color-blue-5)",
            backgroundColor: "var(--mantine-color-dark-7)",
          },
        },
      },
    },

    TagsInput: {
      styles: {
        input: {
          backgroundColor: "var(--mantine-color-dark-8)",
          borderColor: "var(--mantine-color-dark-6)",
          color: "var(--mantine-color-dark-0)",
          transition: "all 200ms ease",
          "&::placeholder": {
            color: "var(--mantine-color-dark-3)",
          },
          "&:focus": {
            borderColor: "var(--mantine-color-blue-5)",
            backgroundColor: "var(--mantine-color-dark-7)",
          },
        },
      },
    },

    Select: {
      styles: {
        input: {
          backgroundColor: "var(--mantine-color-dark-7)",
          borderColor: "var(--mantine-color-dark-6)",
          color: "var(--mantine-color-dark-0)",
          "&::placeholder": {
            color: "var(--mantine-color-dark-3)",
          },
          "&:focus": {
            borderColor: "var(--mantine-color-blue-5)",
            backgroundColor: "var(--mantine-color-dark-6)",
          },
        },
        dropdown: {
          backgroundColor: "var(--mantine-color-dark-6)",
          borderColor: "var(--mantine-color-dark-5)",
          boxShadow: "var(--mantine-shadow-lg)",
        },
        option: {
          "&:hover": {
            backgroundColor: "var(--mantine-color-dark-5)",
          },
          "&[data-selected]": {
            backgroundColor: "var(--mantine-color-blue-8)",
            color: "var(--mantine-color-blue-1)",
          },
        },
      },
    },

    // Data display components
    Table: {
      styles: {
        root: {
          backgroundColor: "var(--mantine-color-dark-7)",
          borderRadius: "var(--mantine-radius-md)",
          overflow: "hidden",
        },
        th: {
          backgroundColor: "var(--mantine-color-dark-6)",
          borderBottom: "1px solid var(--mantine-color-dark-5)",
          color: "var(--mantine-color-dark-0)",
          fontWeight: 600,
        },
        td: {
          borderBottom: "1px solid var(--mantine-color-dark-6)",
          color: "var(--mantine-color-dark-1)",
        },
        tr: {
          transition: "background-color 150ms ease",
          "&:hover": {
            backgroundColor: "var(--mantine-color-dark-6)",
          },
        },
      },
    },

    Badge: {
      styles: {
        root: {
          color: "var(--mantine-color-dark-0)",
          border: "1px solid var(--mantine-color-dark-5)",
        },
      },
    },

    // Overlay components
    Modal: {
      styles: {
        content: {
          backgroundColor: "var(--mantine-color-dark-7)",
          border: "1px solid var(--mantine-color-dark-5)",
          borderRadius: "var(--mantine-radius-lg)",
        },
        header: {
          backgroundColor: "var(--mantine-color-dark-7)",
          borderBottom: "1px solid var(--mantine-color-dark-6)",
        },
        body: {
          padding: "var(--mantine-spacing-lg)",
        },
      },
    },

    Menu: {
      styles: {
        dropdown: {
          backgroundColor: "var(--mantine-color-dark-6)",
          border: "1px solid var(--mantine-color-dark-5)",
          boxShadow: "var(--mantine-shadow-xl)",
          borderRadius: "var(--mantine-radius-md)",
        },
        item: {
          borderRadius: "var(--mantine-radius-sm)",
          transition: "background-color 150ms ease",
          "&:hover": {
            backgroundColor: "var(--mantine-color-dark-5)",
          },
        },
      },
    },

    // Navigation components
    Tabs: {
      styles: {
        root: {
          backgroundColor: "transparent",
        },
        tab: {
          color: "var(--mantine-color-dark-2)",
          fontWeight: 500,
          transition: "all 200ms ease",
          "&:hover": {
            backgroundColor: "var(--mantine-color-dark-6)",
            color: "var(--mantine-color-dark-0)",
          },
          "&[data-active]": {
            color: "var(--mantine-color-blue-4)",
            borderColor: "var(--mantine-color-blue-5)",
            fontWeight: 600,
          },
        },
        tabsList: {
          borderBottom: "1px solid var(--mantine-color-dark-6)",
        },
      },
    },

    // Feedback components
    Notification: {
      styles: {
        root: {
          backgroundColor: "var(--mantine-color-dark-6)",
          border: "1px solid var(--mantine-color-dark-5)",
          borderRadius: "var(--mantine-radius-md)",
          boxShadow: "var(--mantine-shadow-lg)",
        },
      },
    },

    Progress: {
      styles: {
        root: {
          backgroundColor: "var(--mantine-color-dark-6)",
          borderRadius: "var(--mantine-radius-sm)",
        },
        bar: {
          transition: "width 300ms ease",
        },
      },
    },

    // Search components
    Spotlight: {
      styles: {
        root: {
          backgroundColor: "var(--mantine-color-dark-6)",
          border: "1px solid var(--mantine-color-dark-5)",
          borderRadius: "var(--mantine-radius-lg)",
        },
        action: {
          borderRadius: "var(--mantine-radius-sm)",
          "&:hover": {
            backgroundColor: "var(--mantine-color-dark-5)",
          },
        },
      },
    },
  },

  // Design tokens
  defaultRadius: "md",

  spacing: {
    xs: "0.5rem",
    sm: "0.75rem",
    md: "1rem",
    lg: "1.5rem",
    xl: "2rem",
  },

  shadows: {
    xs: "0 1px 2px rgba(0, 0, 0, 0.05)",
    sm: "0 1px 3px rgba(0, 0, 0, 0.12), 0 1px 2px rgba(0, 0, 0, 0.24)",
    md: "0 4px 6px rgba(0, 0, 0, 0.07), 0 2px 4px rgba(0, 0, 0, 0.06)",
    lg: "0 10px 15px rgba(0, 0, 0, 0.1), 0 4px 6px rgba(0, 0, 0, 0.05)",
    xl: "0 20px 25px rgba(0, 0, 0, 0.15), 0 10px 10px rgba(0, 0, 0, 0.04)",
  },

  // Global styles for smooth interactions
  globalStyles: () => ({
    "*": {
      boxSizing: "border-box",
    },

    body: {
      backgroundColor: "var(--mantine-color-dark-9)",
      color: "var(--mantine-color-dark-0)",
      fontFamily: "var(--mantine-font-family)",
    },

    // Custom scrollbars
    "*::-webkit-scrollbar": {
      width: "6px",
      height: "6px",
    },
    "*::-webkit-scrollbar-track": {
      backgroundColor: "var(--mantine-color-dark-8)",
    },
    "*::-webkit-scrollbar-thumb": {
      backgroundColor: "var(--mantine-color-dark-5)",
      borderRadius: "3px",
      "&:hover": {
        backgroundColor: "var(--mantine-color-dark-4)",
      },
    },

    // Selection styling
    "::selection": {
      backgroundColor: "var(--mantine-color-blue-7)",
      color: "var(--mantine-color-blue-0)",
    },
  }),
});
