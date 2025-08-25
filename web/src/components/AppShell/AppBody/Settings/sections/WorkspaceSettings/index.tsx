import { useSettingsState } from "../../../../../../states/settings.state";
import {
  Container,
  Title,
  Stack,
  Select,
  Group,
  Box,
  ColorSwatch,
  Tooltip,
  Switch,
  Text,
  useMantineColorScheme,
} from "@mantine/core";
import { useEffect } from "react";

const colors = ["dark", "blue", "orange", "green", "red"];

export function WorkspaceSettings() {
  const { setColorScheme } = useMantineColorScheme();
  const { workspace_settings, setWorkspaceSettings } = useSettingsState();

  useEffect(() => {
    setColorScheme(workspace_settings.theme);
  }, [workspace_settings.theme]);

  // Save workspace settings to localStorage on change
  useEffect(() => {
    localStorage.setItem(
      "workspaceSettings",
      JSON.stringify(workspace_settings)
    );
  }, [workspace_settings]);

  return (
    <Container fluid p="xs">
      <Title size="h3">Workspace Configuration</Title>
      <Stack gap="md" p="xs">
        <Group justify="space-between">
          <Select
            flex={1}
            label="Theme"
            value={workspace_settings.theme}
            data={[
              { value: "auto", label: "Auto (System)" },
              { value: "light", label: "Light" },
              { value: "dark", label: "Dark" },
            ]}
            allowDeselect={false}
            onChange={(value) => setWorkspaceSettings({ theme: value as any })}
          />
          <Box flex={1}>
            <Text fw={500} size="sm" mb="xs">
              Primary Color
            </Text>
            <Group gap="xs">
              {colors.map((color) => (
                <Tooltip key={color} label={color}>
                  <ColorSwatch
                    color={`var(--mantine-color-${color}-5)`}
                    size={25}
                    style={{ cursor: "pointer" }}
                    onClick={() =>
                      setWorkspaceSettings({ primary_color: color })
                    }
                    component="button"
                    type="button"
                  >
                    {workspace_settings.primary_color === color && (
                      <Text c="white" size="xs">
                        ✓
                      </Text>
                    )}
                  </ColorSwatch>
                </Tooltip>
              ))}
            </Group>
          </Box>
        </Group>
        <Switch
          label="Show Code Lines"
          checked={workspace_settings.show_line_preview}
          onChange={(e) =>
            setWorkspaceSettings({ show_line_preview: e.currentTarget.checked })
          }
        />
      </Stack>
    </Container>
  );
}
