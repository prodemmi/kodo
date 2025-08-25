import {
  Container,
  Title,
  Stack,
  TextInput,
  Button,
  Group,
  Alert,
  Badge,
  Modal,
  Text,
  LoadingOverlay,
  Code,
} from "@mantine/core";
import { IconArrowRight, IconCheck, IconPlus } from "@tabler/icons-react";
import { useState } from "react";
import { SortableColumn } from "../SortableColumn";
import { PrioritySettings } from "./sections/PrioritySettings";
import { DndContext } from "@dnd-kit/core";
import {
  SortableContext,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import {
  useSettings,
  useUpdateSettings,
} from "../../../../../../hooks/use-settings";
import { DemoPatterns } from "./sections/DemoPatterns";
// import snakeCase from "lodash.snakecase";
// import { KanbanColumn } from "../../../../../../types/settings";

export function KanbanSettings() {
  const { data: settings, isLoading, isSuccess } = useSettings();
  const updateSettings = useUpdateSettings();
  // const [newColumnName, setNewColumnName] = useState("");
  const [deleteColumnId, setDeleteColumnId] = useState<string | null>(null);

  // Only allow deleting non-main columns (index > 2)
  const handleDeleteColumn = (id: string, index: number) => {
    if (index > 2) setDeleteColumnId(id);
  };

  const confirmDeleteColumn = () => {
    updateSettings({
      kanban_columns: settings?.kanban_columns.filter(
        (col) => col.id !== deleteColumnId
      ),
    });
    setDeleteColumnId(null);
  };

  // const handleAddColumn = () => {
  //   if (newColumnName.trim()) {
  //     const newCol: KanbanColumn = {
  //       id: snakeCase(newColumnName),
  //       name: newColumnName,
  //       color: "blue",
  //       auto_assign_pattern: "",
  //     };
  //     updateSettings({
  //       kanban_columns: [...(settings?.kanban_columns ?? []), newCol],
  //     });
  //     setNewColumnName("");
  //   }
  // };

  if (!isSuccess || isLoading) return <LoadingOverlay />;

  return (
    settings && (
      <>
        <Container fluid p="xs">
          {/* <Title fw={600} size="h3">
            Kanban Board Setup
          </Title> */}

          <Stack gap="md" p="xs">
            {/* <Group align="flex-end" px="xs">
              <TextInput
                label="Column Name"
                value={newColumnName}
                onChange={(e) => setNewColumnName(e.currentTarget.value)}
                style={{ flex: 1 }}
              />
              <Button
                leftSection={<IconPlus size={16} />}
                onClick={handleAddColumn}
              >
                Add Column
              </Button>
            </Group> */}

            <Title fw={600} size="h3">
              Kanban Columns
            </Title>

            <Stack gap="md" px="xs">
              <Alert variant="light" color="blue">
                <Group gap="xs" align="center">
                  <Text fw={500}>Workflow direction:</Text>
                  {settings?.kanban_columns.map((column, index) => (
                    <Group key={column.id} gap="xs" align="center">
                      <Text c={column.color} fw={500}>
                        {column.name}
                      </Text>
                      {index < settings.kanban_columns.length - 1 && (
                        <IconArrowRight size={14} />
                      )}
                      {index === settings.kanban_columns.length - 1 && (
                        <Badge
                          color="green"
                          size="sm"
                          radius="sm"
                          leftSection={<IconCheck size={12} />}
                        >
                          Done
                        </Badge>
                      )}
                    </Group>
                  ))}
                </Group>
              </Alert>
              <DndContext
                onDragEnd={(event) => {
                  /* implement drag end handler */
                }}
              >
                <SortableContext
                  items={settings.kanban_columns.map((col) => col.id)}
                  strategy={verticalListSortingStrategy}
                >
                  <Stack>
                    {settings.kanban_columns.map((column, idx) => (
                      <SortableColumn
                        key={column.id}
                        column={column}
                        // onDelete={() => handleDeleteColumn(column.id, idx)}
                        showPatternInput={idx === 0}
                      />
                    ))}
                  </Stack>
                </SortableContext>
              </DndContext>
            </Stack>

            <Title fw={600} size="h3">
              Priority Patterns
            </Title>

            <PrioritySettings
              priority_patterns={settings.priority_patterns}
              setPriorities={(priority_patterns) =>
                updateSettings({ priority_patterns })
              }
            />
            <DemoPatterns />
          </Stack>
        </Container>

        {/* Delete Confirmation Modal */}
        <Modal
          opened={!!deleteColumnId}
          onClose={() => setDeleteColumnId(null)}
          title="Delete Column"
          centered
        >
          <Stack>
            <Text>
              Are you sure you want to delete the column "
              {
                settings.kanban_columns.find((col) => col.id === deleteColumnId)
                  ?.name
              }
              "?
            </Text>
            <Group justify="flex-end">
              <Button variant="outline" onClick={() => setDeleteColumnId(null)}>
                Cancel
              </Button>
              <Button color="red" onClick={confirmDeleteColumn}>
                Delete
              </Button>
            </Group>
          </Stack>
        </Modal>
      </>
    )
  );
}
