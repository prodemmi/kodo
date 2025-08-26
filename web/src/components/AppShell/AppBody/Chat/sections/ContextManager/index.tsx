import { useState } from "react";
import {
  Box,
  Button,
  Card,
  Group,
  Text,
  TextInput,
  ActionIcon,
  Badge,
  Stack,
  Collapse,
  Textarea,
  Modal,
  Menu,
  MenuItem,
  MenuTarget,
  MenuDropdown,
  Code,
  Tooltip,
} from "@mantine/core";
import {
  IconPlus,
  IconTrash,
  IconEdit,
  IconFile,
  IconCode,
  IconChevronDown,
  IconChevronRight,
  IconDots,
  IconCopy,
  IconX,
  IconBrain,
  IconCodeDots,
  IconFileText,
  IconEye,
} from "@tabler/icons-react";
import { Context } from "../../../../../../types/context";
import { useContextStore } from "../../../../../../states/context.state";

type Props = {
  opened: boolean;
  onClose: () => void;
};

export default function ContextManager({ opened, onClose }: Props) {
  const {
    contexts,
    addContext,
    removeFile,
    removeContext,
    removeSnippet,
    updateContextName,
    duplicateContext,
  } = useContextStore((state) => state);

  const [expandedContexts, setExpandedContexts] = useState<Set<string>>(
    new Set(["1"])
  );
  const [newContextModalOpen, setNewContextModalOpen] = useState(false);
  const [editingContext, setEditingContext] = useState<string | null>(null);
  const [newContextName, setNewContextName] = useState("");
  const [newContextDescription, setNewContextDescription] = useState("");

  const createContext = () => {
    if (!newContextName.trim()) return;

    const newContext: Context = {
      id: Date.now().toString(),
      name: newContextName,
      description: newContextDescription || undefined,
      files: [],
      snippets: [],
      createdAt: new Date(),
      updatedAt: new Date(),
    };

    addContext(newContext);
    setNewContextName("");
    setNewContextDescription("");
    setNewContextModalOpen(false);
    setExpandedContexts((prev) => new Set([...prev, newContext.id]));
  };

  const toggleExpanded = (contextId: string) => {
    setExpandedContexts((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(contextId)) {
        newSet.delete(contextId);
      } else {
        newSet.add(contextId);
      }
      return newSet;
    });
  };

  return (
    <>
      <Modal opened={opened} onClose={onClose} size="xl">
        <Box p="md">
          <Group justify="space-between" mb="lg">
            <Group gap="sm">
              <IconBrain size={24} />
              <Text size="xl" fw={600}>
                AI Contexts
              </Text>
              <Badge variant="light" size="sm">
                {contexts.length}
              </Badge>
            </Group>

            <Button
              leftSection={<IconPlus size={16} />}
              onClick={() => setNewContextModalOpen(true)}
              size="sm"
            >
              New Context
            </Button>
          </Group>

          <Stack gap="sm">
            {contexts.map((context) => {
              const isExpanded = expandedContexts.has(context.id);
              const totalItems = context.files.length + context.snippets.length;

              return (
                <Card key={context.id} padding="sm" withBorder>
                  <Group
                    justify="space-between"
                    style={{ cursor: "pointer" }}
                    wrap="nowrap"
                  >
                    <Group
                      gap="sm"
                      style={{ flex: 1 }}
                      onClick={() => toggleExpanded(context.id)}
                      className="cursor-pointer"
                    >
                      {isExpanded ? (
                        <IconChevronDown size={16} />
                      ) : (
                        <IconChevronRight size={16} />
                      )}

                      <Box style={{ flex: 1 }}>
                        {editingContext === context.id ? (
                          <TextInput
                            value={context.name}
                            onChange={(e) =>
                              updateContextName(context.id, e.target.value)
                            }
                            onBlur={() => setEditingContext(null)}
                            onKeyDown={(e) => {
                              if (e.key === "Enter") setEditingContext(null);
                              if (e.key === "Escape") setEditingContext(null);
                            }}
                            size="sm"
                            autoFocus
                            onClick={(e) => e.stopPropagation()}
                          />
                        ) : (
                          <Group gap="sm">
                            <Text fw={500}>{context.name}</Text>
                            <Badge variant="light" size="xs">
                              {totalItems} items
                            </Badge>
                          </Group>
                        )}

                        {context.description && !editingContext && (
                          <Text size="xs" c="dimmed" mt={2}>
                            {context.description}
                          </Text>
                        )}
                      </Box>
                    </Group>

                    <Group gap={4}>
                      <Tooltip label="Add to Context">
                        <ActionIcon size={20} variant="subtle">
                          <IconFileText size={12} />
                        </ActionIcon>
                      </Tooltip>

                      <Tooltip label="Select Code">
                        <ActionIcon size={20} variant="subtle" color="green">
                          <IconCode size={12} />
                        </ActionIcon>
                      </Tooltip>

                      <Tooltip label="Preview">
                        <ActionIcon size={20} variant="subtle" color="orange">
                          <IconEye size={12} />
                        </ActionIcon>
                      </Tooltip>

                      <Menu position="bottom-end" withinPortal>
                        <MenuTarget>
                          <ActionIcon size={20} variant="subtle" color="gray">
                            <IconDots size={12} />
                          </ActionIcon>
                        </MenuTarget>

                        <MenuDropdown>
                          <MenuItem
                            leftSection={<IconEdit size={14} />}
                            onClick={() => setEditingContext(context.id)}
                          >
                            Rename
                          </MenuItem>
                          <MenuItem
                            leftSection={<IconCopy size={14} />}
                            onClick={() => duplicateContext(context.id)}
                          >
                            Duplicate
                          </MenuItem>
                          <Menu.Divider />
                          <MenuItem
                            leftSection={<IconTrash size={14} color="red" />}
                            onClick={() => removeContext(context.id)}
                          >
                            Delete
                          </MenuItem>
                        </MenuDropdown>
                      </Menu>
                    </Group>
                  </Group>

                  <Collapse in={isExpanded}>
                    <Box mt="sm">
                      {context.files.length > 0 && (
                        <Box mb="sm" pl="xs">
                          <Stack gap={4}>
                            {context.files.map((file) => (
                              <Group key={file.path} gap="sm" pl="md">
                                <IconFile size={12} style={{ opacity: 0.6 }} />
                                <Text size="xs" style={{ flex: 1 }} truncate>
                                  {file.name}
                                </Text>
                                <Code c="dimmed">{file.path}</Code>
                                <ActionIcon
                                  size={16}
                                  variant="subtle"
                                  color="red"
                                  onClick={() =>
                                    removeFile(context.id, file.id)
                                  }
                                >
                                  <IconX size={10} />
                                </ActionIcon>
                              </Group>
                            ))}
                          </Stack>
                        </Box>
                      )}

                      {context.snippets.length > 0 && (
                        <Box>
                          <Group gap="xs" mb="xs">
                            <IconCodeDots size={14} />
                            <Text size="sm" fw={500}>
                              Code Snippets ({context.snippets.length})
                            </Text>
                          </Group>

                          <Stack gap={4}>
                            {context.snippets.map((snippet) => (
                              <Group key={snippet.id} gap="sm" pl="md">
                                <IconCode size={12} style={{ opacity: 0.6 }} />
                                <Box style={{ flex: 1 }}>
                                  <Text size="xs" truncate>
                                    {snippet.fileName}
                                  </Text>
                                  <Text size="xs" c="dimmed">
                                    Lines {snippet.startLine}-{snippet.endLine}
                                  </Text>
                                </Box>
                                <Code c="dimmed">{snippet.language}</Code>
                                <ActionIcon
                                  size={16}
                                  variant="subtle"
                                  color="red"
                                  onClick={() =>
                                    removeSnippet(context.id, snippet.id)
                                  }
                                >
                                  <IconX size={10} />
                                </ActionIcon>
                              </Group>
                            ))}
                          </Stack>
                        </Box>
                      )}

                      {totalItems === 0 && (
                        <Text size="sm" c="dimmed" ta="center" py="md">
                          No files or code snippets added yet
                        </Text>
                      )}
                    </Box>
                  </Collapse>
                </Card>
              );
            })}

            {contexts.length === 0 && (
              <Card padding="xl" withBorder>
                <Stack align="center" gap="md">
                  <IconBrain size={48} style={{ opacity: 0.5 }} />
                  <Text size="lg" fw={500}>
                    No Contexts Created
                  </Text>
                  <Text size="sm" c="dimmed" ta="center">
                    Create your first context to start organizing code for AI
                    assistance
                  </Text>
                  <Button
                    leftSection={<IconPlus size={16} />}
                    onClick={() => setNewContextModalOpen(true)}
                  >
                    Create First Context
                  </Button>
                </Stack>
              </Card>
            )}
          </Stack>
        </Box>
      </Modal>
      <Modal
        opened={newContextModalOpen}
        onClose={() => {
          setNewContextModalOpen(false);
          setNewContextName("");
          setNewContextDescription("");
        }}
        title="Create New Context"
        size="md"
      >
        <Stack gap="md">
          <TextInput
            label="Context Name"
            placeholder="e.g., Authentication System, Payment Flow"
            value={newContextName}
            onChange={(e) => setNewContextName(e.target.value)}
            required
          />

          <Textarea
            label="Description (Optional)"
            placeholder="Describe what this context contains and its purpose..."
            value={newContextDescription}
            onChange={(e) => setNewContextDescription(e.target.value)}
            rows={3}
          />

          <Group justify="flex-end" gap="sm">
            <Button
              variant="light"
              onClick={() => {
                setNewContextModalOpen(false);
                setNewContextName("");
                setNewContextDescription("");
              }}
            >
              Cancel
            </Button>
            <Button onClick={createContext} disabled={!newContextName.trim()}>
              Create Context
            </Button>
          </Group>
        </Stack>
      </Modal>
    </>
  );
}
