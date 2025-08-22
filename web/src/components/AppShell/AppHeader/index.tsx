import {
  ActionIcon,
  AppShell,
  Text,
  Button,
  ButtonGroup,
  Flex,
  Box,
} from "@mantine/core";
import { useAppState } from "../../../states/app.state";
import { IconSettings } from "@tabler/icons-react";
import { RoleGuard } from "../../Investor";
type Props = {};

export default function AppHeader({}: Props) {
  const { activeTab, setActiveTab } = useAppState((state) => state);

  return (
    <AppShell.Header>
      <Flex justify="space-between" align="center" h="100%" p="xs">
        <ButtonGroup>
          <Button
            bg={
              activeTab === 0
                ? "var(--mantine-color-dark-6)"
                : "var(--mantine-color-dark-7)"
            }
            onClick={() => setActiveTab(0)}
          >
            Board
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
          <Button
            bg={
              activeTab === 1
                ? "var(--mantine-color-dark-6)"
                : "var(--mantine-color-dark-7)"
            }
            onClick={() => setActiveTab(1)}
          >
            Notes
          </Button>
        </ButtonGroup>

        <RoleGuard.Consumer>
          <ActionIcon
            bg={
              activeTab === 2
                ? "var(--mantine-color-dark-6)"
                : "var(--mantine-color-dark-7)"
            }
            onClick={() => setActiveTab(2)}
          >
            <IconSettings />
          </ActionIcon>
        </RoleGuard.Consumer>

        <RoleGuard.Investor>
          <Box bg="red" c="white" px="xs" py="3">
            <Text size="xs">View Only</Text>
          </Box>
        </RoleGuard.Investor>
      </Flex>
    </AppShell.Header>
  );
}
