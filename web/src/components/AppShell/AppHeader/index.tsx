import {
  AppShell,
  Text,
  Button,
  ButtonGroup,
  Flex,
  Box,
  useMantineTheme,
  Group,
} from "@mantine/core";
import { useAppState } from "../../../states/app.state";
import { IconRefresh, IconSettings } from "@tabler/icons-react";
import { RoleGuard } from "../../Investor";
import { useCallback } from "react";
import packageJson from "../../../../package.json";
import { useSettings } from "../../../hooks/use-settings";
import { useSyncNotes } from "../../../hooks/use-notes";

type Props = {};

export default function AppHeader({}: Props) {
  const { data: settings } = useSettings();
  const { mutate: syncNotes, isPending: isLoadingSyncNotes } = useSyncNotes();
  const { activeTab, setActiveTab } = useAppState((state) => state);
  const { primaryColor, colors } = useMantineTheme();

  const buttonBg = useCallback(
    (tab: number) => {
      if (activeTab === tab) return colors[primaryColor][6];
      return colors[primaryColor][8];
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

          <RoleGuard.Consumer>
            <Button
              onClick={() => setActiveTab(2)}
              bg={buttonBg(2)}
              rightSection={<IconSettings size={18} />}
            >
              Settings
            </Button>
          </RoleGuard.Consumer>
        </ButtonGroup>

        <Group align="center" justify="flex-end" gap="sm">
          {settings &&
            settings.code_scan_settings.sync_enabled &&
            activeTab === 1 && (
              <Button
                onClick={() => syncNotes()}
                loading={isLoadingSyncNotes}
                rightSection={<IconRefresh size={18} />}
              >
                Sync Notes
              </Button>
            )}

          <RoleGuard.Investor>
            <Box bg="red" color="white" px="xs" py="3">
              <Text size="xs">View Only</Text>
            </Box>
          </RoleGuard.Investor>

          <Box px="xs" py="3">
            <Text size="xs">V{packageJson.version}</Text>
          </Box>
        </Group>
      </Flex>
    </AppShell.Header>
  );
}
