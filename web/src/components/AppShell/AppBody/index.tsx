import { AppShell } from "@mantine/core";
import { useAppState } from "../../../states/app.state";
import Board from "./Board";
import Settings from "./Settings";
import Notes from "./Notes";

const tabs = [<Board />, /*<Chat />,*/ <Notes />, <Settings />];

export default function AppBody() {
  const activeTab = useAppState((state) => state.activeTab);

  return (
    <AppShell.Main
      pt="calc(var(--app-shell-header-offset, 0rem))"
      px="0"
      pb="0"
    >
      {tabs[activeTab] || null}
    </AppShell.Main>
  );
}
