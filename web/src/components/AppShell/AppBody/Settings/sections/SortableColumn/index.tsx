import { CSS } from "@dnd-kit/utilities";
import { useSortable } from "@dnd-kit/sortable";
import { IconGripVertical } from "@tabler/icons-react";
import {
  Paper,
  Stack,
  Group,
  TextInput,
  ActionIcon,
  ColorSwatch,
  Text,
  LoadingOverlay,
} from "@mantine/core";
import {
  useSettings,
  useUpdateSettings,
} from "../../../../../../hooks/use-settings";
import { useCallback } from "react";

export function SortableColumn({
  column,
  showPatternInput,
}: {
  column: any;
  showPatternInput: boolean;
}) {
  const { data: settings, isLoading, isSuccess } = useSettings();
  const updateSettings = useUpdateSettings();
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: column.id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  const handleNameChange = (name: string) => {
    const cols = settings?.kanban_columns.map((col) =>
      col.id === column.id ? { ...col, name } : col
    );
    updateSettings({ kanban_columns: cols });
  };

  const handleColorChange = (color: string) => {
    const cols = settings?.kanban_columns.map((col) =>
      col.id === column.id ? { ...col, color } : col
    );
    updateSettings({ kanban_columns: cols });
  };

  const debouncedPatternChange = useCallback(
    (pattern: string) => {
      const cols = settings?.kanban_columns.map((col) =>
        col.id === column.id ? { ...col, auto_assign_pattern: pattern } : col
      );
      updateSettings({ kanban_columns: cols });
    },
    [settings, column.id, updateSettings]
  );

  const colors = ["dark", "blue", "orange", "green", "red"];

  if (!isSuccess || isLoading) return <LoadingOverlay />;

  return (
    <Paper
      shadow="xs"
      p="md"
      mt="xs"
      radius="md"
      withBorder
      ref={setNodeRef}
      style={style}
    >
      <Stack gap="sm">
        <Group justify="space-between">
          <Group gap="xs" style={{ flex: 1 }}>
            <ActionIcon
              {...attributes}
              {...listeners}
              style={{ cursor: "grab", padding: "4px" }}
              size="md"
              variant="transparent"
            >
              <IconGripVertical />
            </ActionIcon>
            <TextInput
              value={column.name}
              onChange={(e) => handleNameChange(e.currentTarget.value)}
              placeholder="Column name"
              size="sm"
              radius="md"
              style={{ flex: 1 }}
            />
          </Group>
        </Group>
        {showPatternInput && (
          <TextInput
            value={column.auto_assign_pattern || ""}
            onChange={(e) => debouncedPatternChange(e.currentTarget.value)}
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
