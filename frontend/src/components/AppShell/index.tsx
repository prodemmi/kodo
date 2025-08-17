import { ReactNode } from "react";
import AppHeader from "./AppHeader";
import { AppShell } from "@mantine/core";
import AppBody from "./AppBody";

type Props = {
  children: ReactNode;
};

export default function ({}: Props) {
  return (
    <AppShell
      padding="md"
      header={{ height: 52 }}
      styles={{ main: { overflow: "hidden" } }}
    >
      <AppHeader />
      <AppBody />
    </AppShell>
  );
}
