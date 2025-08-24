import {
  ActionIcon,
  AppShell,
  Text,
  Button,
  ButtonGroup,
  Flex,
  Box,
  useMantineTheme,
} from "@mantine/core";
import { useAppState } from "../../../states/app.state";
import { IconSettings } from "@tabler/icons-react";
import { RoleGuard } from "../../Investor";
import { useCallback } from "react";

type Props = {};

export default function AppHeader({}: Props) {
  const { activeTab, setActiveTab } = useAppState((state) => state);
  const { primaryColor, colors } = useMantineTheme();

  const buttonBg = useCallback(
    (tab: number) => {
      if (activeTab === tab) return colors[primaryColor][5];
      return colors[primaryColor][4];
    },
    [activeTab, primaryColor]
  );

  return (
    <AppShell.Header>
      <Flex justify="space-between" align="center" h="100%" p="xs">
        <ButtonGroup>
          <Button onClick={() => setActiveTab(0)} bg={buttonBg(0)}>
            Kanban
          </Button>
          {/* <Button
             bg={
              activeTab === 02
                ? "var(--mantine-color-dark-6)"
                : "var(--mantine-color-dark-7)"
            }
            onClick={() => setActiveTab(1)}
          >
            Chat
          </Button> */}
          <Button onClick={() => setActiveTab(1)} bg={buttonBg(1)}>
            Notes
          </Button>
        </ButtonGroup>

        <RoleGuard.Consumer>
          <ActionIcon onClick={() => setActiveTab(2)} c={buttonBg(2)} variant="subtle">
            <IconSettings />
          </ActionIcon>
        </RoleGuard.Consumer>

        <RoleGuard.Investor>
          <Box bg="red" color="white" px="xs" py="3">
            <Text size="xs">View Only</Text>
          </Box>
        </RoleGuard.Investor>
      </Flex>
    </AppShell.Header>
  );
}
