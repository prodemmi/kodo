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
import { useMemo } from "react";

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
  const primaryColor = useAppState((s) => s.primaryColor);

  const theme = useMemo(() => createTheme(primaryColor), [primaryColor]);
  console.log('primaryColor ===>', primaryColor);
  
  return (
    <PersistQueryClientProvider
      client={queryClient}
      persistOptions={{ persister }}
    >
      <MantineProvider theme={theme} forceColorScheme="dark">
        <Notifications />
        <AppShell>App</AppShell>
      </MantineProvider>
    </PersistQueryClientProvider>
  );
}
