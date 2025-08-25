import { useEffect } from "react";
import AppHeader from "./AppHeader";
import { AppShell } from "@mantine/core";
import AppBody from "./AppBody";
import { useInvestor } from "../../hooks/use-config";
import { useAppState } from "../../states/app.state";

type Props = {};

export default function ({}: Props) {
  const { data: investor, isError, isLoading } = useInvestor();
  const { setInvestor } = useAppState();

  useEffect(() => {
    if (investor && !isError && !isLoading) {
      setInvestor(investor.investor);
    } else {
      setInvestor(false);
    }
  }, [investor, isError, isLoading]);

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
