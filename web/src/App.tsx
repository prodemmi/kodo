import "./global.css";
import "@mantine/core/styles.css";
import "@mantine/notifications/styles.css";
import { MantineProvider } from "@mantine/core";
import { createTheme } from "./theme";
import AppShell from "./components/AppShell";
import { Notifications } from "@mantine/notifications";
import { useEffect, useMemo, useState } from "react";
import { useSettingsState } from "./states/settings.state";

export default function App() {
  const [isDark, setIsDark] = useState(true);
  const workspace_settings = useSettingsState((s) => s.workspace_settings);
  console.log('workspace_settings ===>', workspace_settings);
  

  useEffect(() => {
    const html = document.documentElement;
    const observer = new MutationObserver((mutations) => {
      for (const mutation of mutations) {
        if (mutation.type === "attributes") {
          setIsDark(html.getAttribute(mutation.attributeName!) === "dark");
        }
      }
    });
    observer.observe(html, {
      attributes: true,
      attributeFilter: ["data-mantine-color-scheme"],
    });
    return () => observer.disconnect();
  }, []);

  const theme = useMemo(
    () => createTheme(workspace_settings.primary_color, isDark),
    [workspace_settings, isDark]
  );

  return (
    <MantineProvider theme={theme} defaultColorScheme="dark">
      <Notifications />
      <AppShell />
    </MantineProvider>
  );
}
