import {
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
import packageJson from "../../../../package.json";

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

          <RoleGuard.Consumer>
            <Button
              onClick={() => setActiveTab(2)}
              bg={buttonBg(2)}
              rightSection={<IconSettings />}
            >
              Settings
            </Button>
          </RoleGuard.Consumer>
        </ButtonGroup>

        <RoleGuard.Investor>
          <Box bg="red" color="white" px="xs" py="3">
            <Text size="xs">View Only</Text>
          </Box>
        </RoleGuard.Investor>

        <Box px="xs" py="3">
          <Text size="xs">V{packageJson.version}</Text>
        </Box>
      </Flex>
    </AppShell.Header>
  );
}
