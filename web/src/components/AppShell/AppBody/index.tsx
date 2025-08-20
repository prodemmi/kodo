import { AppShell } from "@mantine/core";
import { useAppState } from "../../../states/app.state";
import Board from "./Board";
import Chat from "./Chat";
import History from "./History";
import Settings from "./Settings";
import Notes from "./Notes";

const tabs = [<Board />, /*<Chat />,*/ <Notes />, <History />, <Settings />];

export default function AppBody() {
  const activeTab = useAppState((state) => state.activeTab);

  return <AppShell.Main>{tabs[activeTab] || null}</AppShell.Main>;
}
