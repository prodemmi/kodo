import { ActionIcon, AppShell, Button, ButtonGroup, Flex } from "@mantine/core";
import { useAppState } from "../../../states/app.state";
import { IconSettings } from "@tabler/icons-react";

type Props = {};

export default function AppHeader({}: Props) {
  const { activeTab, setActiveTab } = useAppState((state) => state);

  return (
    <AppShell.Header>
      <Flex justify="space-between" align="center" h="100%" p="xs">
        <ButtonGroup>
          <Button
            variant={activeTab === 0 ? "light" : "subtle"}
            onClick={() => setActiveTab(0)}
          >
            Board
          </Button>
          <Button
            variant={activeTab === 1 ? "light" : "subtle"}
            onClick={() => setActiveTab(1)}
          >
            Chat
          </Button>
          <Button
            variant={activeTab === 2 ? "light" : "subtle"}
            onClick={() => setActiveTab(2)}
          >
            History
          </Button>
        </ButtonGroup>

        <ActionIcon
          variant={activeTab === 3 ? "light" : "subtle"}
          onClick={() => setActiveTab(3)}
        >
          <IconSettings />
        </ActionIcon>
      </Flex>
    </AppShell.Header>
  );
}
