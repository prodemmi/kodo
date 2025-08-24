import "./global.css";
import "@mantine/core/styles.css";
import "@mantine/notifications/styles.css";
import { MantineProvider } from "@mantine/core";
import { createTheme } from "./theme";
import AppShell from "./components/AppShell";
import { QueryClient } from "@tanstack/react-query";
import { PersistQueryClientProvider } from "@tanstack/react-query-persist-client";
import { createAsyncStoragePersister } from "@tanstack/query-async-storage-persister";
import { Notifications } from "@mantine/notifications";
import { useAppState } from "./states/app.state";
import { useEffect, useState } from "react";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      gcTime: 1000 * 60 * 60 * 24, // 24 hours
    },
  },
});

const persister = createAsyncStoragePersister({
  storage: window.localStorage,
});

export default function App() {
  const [isDark, setIsDark] = useState(true);
  const primaryColor = useAppState((s) => s.primaryColor);

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

  const theme = createTheme(primaryColor, isDark);

  return (
    <PersistQueryClientProvider
      client={queryClient}
      persistOptions={{ persister }}
    >
      <MantineProvider theme={theme} defaultColorScheme="dark">
        <Notifications />
        <AppShell>App</AppShell>
      </MantineProvider>
    </PersistQueryClientProvider>
  );
}
