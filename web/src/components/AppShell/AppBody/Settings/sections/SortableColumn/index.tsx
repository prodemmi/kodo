import { useEffect, useState, useCallback } from "react";
import debounce from "lodash.debounce";
import {
  Text,
  LoadingOverlay,
  Paper,
  Stack,
  Group,
  TextInput,
  ColorSwatch,
} from "@mantine/core";
import snakeCase from "lodash.snakecase";
import {
  useSettings,
  useUpdateSettings,
} from "../../../../../../hooks/use-settings";

export function SortableColumn({
  column,
  showPatternInput,
}: {
  column: any;
  showPatternInput: boolean;
}) {
  const { data: settings, isLoading, isSuccess } = useSettings();
  const updateSettings = useUpdateSettings();

  // --- Local states
  const [localName, setLocalName] = useState(column.name);
  const [localPattern, setLocalPattern] = useState(
    column.auto_assign_pattern || ""
  );

  // Sync when column prop changes (after remote update)
  useEffect(() => {
    setLocalName(column.name);
    setLocalPattern(column.auto_assign_pattern || "");
  }, [column.name, column.auto_assign_pattern]);

  const handleNameChange = (val: string) => {
    setLocalName(val);
  };

  const handleNameBlur = () => {
    const cols = settings?.kanban_columns.map((col) =>
      col.id === column.id
        ? { ...col, name: localName, id: snakeCase(localName) }
        : col
    );
    updateSettings({ kanban_columns: cols });
  };

  const handleColorChange = (color: string) => {
    const cols = settings?.kanban_columns.map((col) =>
      col.id === column.id ? { ...col, color } : col
    );
    updateSettings({ kanban_columns: cols });
  };

  const debouncedPatternUpdate = useCallback(
    debounce((pattern: string) => {
      const cols = settings?.kanban_columns.map((col) =>
        col.id === column.id ? { ...col, auto_assign_pattern: pattern } : col
      );
      updateSettings({ kanban_columns: cols });
    }, 400),
    [settings, column.id, updateSettings]
  );

  const handlePatternChange = (val: string) => {
    setLocalPattern(val);
    debouncedPatternUpdate(val);
  };

  const colors = ["dark", "blue", "orange", "green", "red"];

  if (!isSuccess || isLoading) return <LoadingOverlay />;

  return (
    <Paper shadow="xs" p="md" mt="xs" radius="md" withBorder>
      <Stack gap="sm">
        <Group justify="space-between">
          <Group gap="xs" style={{ flex: 1 }}>
            <TextInput
              value={localName}
              onChange={(e) => handleNameChange(e.currentTarget.value)}
              onBlur={handleNameBlur}
              placeholder="Column name"
              size="sm"
              radius="md"
              style={{ flex: 1 }}
            />
          </Group>
        </Group>

        {showPatternInput && (
          <TextInput
            value={localPattern}
            onChange={(e) => handlePatternChange(e.currentTarget.value)}
            placeholder="Auto-assign pattern (e.g., 'TODO:')"
            size="sm"
            radius="md"
          />
        )}

        <Group gap="xs">
          <Text size="sm" c="dimmed">
            Color:
          </Text>
          {colors.map((color) => (
            <ColorSwatch
              key={color}
              color={`var(--mantine-color-${color}-6)`}
              size={20}
              style={{ cursor: "pointer" }}
              onClick={() => handleColorChange(color)}
              component="button"
              type="button"
            >
              {column.color === color && (
                <Text c="white" size="xs">
                  ✓
                </Text>
              )}
            </ColorSwatch>
          ))}
        </Group>
      </Stack>
    </Paper>
  );
}
