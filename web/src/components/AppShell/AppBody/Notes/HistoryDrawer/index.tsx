import {
  Text,
  Card,
  Drawer,
  LoadingOverlay,
  Stack,
  Title,
  Badge,
  Group,
  Timeline,
  ActionIcon,
  Tooltip,
  Paper,
  ScrollArea,
  Box,
  Divider,
  ThemeIcon,
  useMantineTheme,
} from "@mantine/core";
import {
  IconJson,
  IconClock,
  IconUser,
  IconGitBranch,
  IconGitCommit,
  IconEdit,
  IconPlus,
  IconTrash,
  IconHistory,
} from "@tabler/icons-react";
import { useNoteHistoryModalStore } from "../../../../../states/note.state";
import { useNoteHistory } from "../../../../../hooks/use-notes";
import { NoteHistory } from "../../../../../types/note";
import { useMemo } from "react";
import { RoleGuard } from "../../../../Investor";

export default function HistoryDrawer() {
  const note = useNoteHistoryModalStore((s) => s.note);
  const isOpen = useNoteHistoryModalStore((s) => s.isOpen);
  const closeModal = useNoteHistoryModalStore((s) => s.closeModal);
  const noteId = note?.id!;
  const { data: history, isLoading } = useNoteHistory(noteId);
  const { primaryColor } = useMantineTheme();

  const formatDate = (
    dateStr: string
  ): { date: string; time: string; relative: string } => {
    try {
      const date = new Date(dateStr);
      if (!(date instanceof Date) || isNaN(date.getTime())) {
        return { date: "Invalid Date", time: "", relative: "" };
      }

      const now = new Date();
      const diffMs = now.getTime() - date.getTime();
      const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
      const diffDays = Math.floor(diffHours / 24);

      let relative = "";
      if (diffDays === 0) {
        if (diffHours === 0) relative = "Just now";
        else if (diffHours === 1) relative = "1 hour ago";
        else relative = `${diffHours} hours ago`;
      } else if (diffDays === 1) {
        relative = "Yesterday";
      } else if (diffDays < 7) {
        relative = `${diffDays} days ago`;
      } else {
        relative = date.toLocaleDateString();
      }

      return {
        date: date.toLocaleDateString(),
        time: date.toLocaleTimeString([], {
          hour: "2-digit",
          minute: "2-digit",
        }),
        relative,
      };
    } catch {
      return { date: "Invalid Date", time: "", relative: "" };
    }
  };

  const getActionIcon = (action: string) => {
    switch (action?.toLowerCase()) {
      case "create":
      case "created":
        return <IconPlus size={16} />;
      case "update":
      case "updated":
      case "edit":
      case "edited":
        return <IconEdit size={16} />;
      case "delete":
      case "deleted":
        return <IconTrash size={16} />;
      default:
        return <IconHistory size={16} />;
    }
  };

  const getActionColor = (action: string) => {
    switch (action?.toLowerCase()) {
      case "create":
      case "created":
        return "green";
      case "update":
      case "updated":
      case "edit":
      case "edited":
        return "blue";
      case "delete":
      case "deleted":
        return "red";
      default:
        return "gray";
    }
  };

  const showJson = () => {
    window.open(
      `http://localhost:8080/api/notes/history?noteId=${noteId}`,
      "_blank"
    );
  };

  const headerContent = useMemo(
    () => (
      <Group justify="space-between" align="center" w="100%" pr="sm">
        <Group gap="sm">
          <ThemeIcon variant="light" size="lg">
            <IconHistory size={20} />
          </ThemeIcon>
          <Box>
            <Title size="h3">History</Title>
            <Text size="sm" c="dimmed">
              {note?.title || "Untitled Note"}
            </Text>
          </Box>
        </Group>
        <RoleGuard.Consumer>
          <Text
            size="xs"
            style={{ cursor: "pointer" }}
            td="underline"
            onClick={showJson}
            c={primaryColor}
          >
            Show Json
          </Text>
        </RoleGuard.Consumer>
      </Group>
    ),
    [note?.title]
  );

  if (!note || !noteId) {
    return null;
  }

  return (
    <Drawer
      opened={isOpen}
      onClose={closeModal}
      title={headerContent}
      size="xl"
      styles={{
        header: { background: "var(--mantine-color-dark-8)" },
        title: { width: "100%" },
        body: { padding: 0 },
      }}
    >
      <Box pos="relative" h="100%">
        <LoadingOverlay visible={isLoading} />

        <ScrollArea h="100%" px="md">
          {!isLoading && !history?.history?.length && (
            <Paper p="xl" withBorder radius="md" mt="md">
              <Stack align="center" gap="md">
                <ThemeIcon size={60} variant="light" color="gray">
                  <IconHistory size={30} />
                </ThemeIcon>
                <Text size="lg" fw={500} ta="center">
                  No History Available
                </Text>
                <Text c="dimmed" ta="center">
                  This note doesn't have any recorded history yet.
                </Text>
              </Stack>
            </Paper>
          )}

          {!isLoading && history && history.count > 0 && (
            <Box py="md">
              {/* Summary Stats */}
              <Paper p="md" withBorder radius="md" mb="lg">
                <Group justify="space-between">
                  <Group gap="lg">
                    <Box ta="center">
                      <Text size="xl" fw={700}>
                        {history.history.length}
                      </Text>
                      <Text size="sm" c="dimmed">
                        Changes
                      </Text>
                    </Box>
                    <Divider orientation="vertical" />
                    <Box ta="center">
                      <Text size="xl" fw={700} c="green">
                        {new Set(history.history.map((h) => h.author)).size}
                      </Text>
                      <Text size="sm" c="dimmed">
                        Contributors
                      </Text>
                    </Box>
                  </Group>
                  <Badge size="lg" variant="light">
                    Last updated{" "}
                    {formatDate(history.history[0]?.timestamp).relative}
                  </Badge>
                </Group>
              </Paper>

              {/* Timeline */}
              <Timeline active={-1} bulletSize={24} lineWidth={2}>
                {history.history.map((entry: NoteHistory, index: number) => {
                  const dateInfo = formatDate(entry.timestamp);
                  const isLatest = index === 0;

                  return (
                    <Timeline.Item
                      key={entry.id}
                      bullet={
                        <ThemeIcon
                          size={20}
                          color={getActionColor(entry.action)}
                          variant={isLatest ? "filled" : "light"}
                        >
                          {getActionIcon(entry.action)}
                        </ThemeIcon>
                      }
                    >
                      <Card
                        p="md"
                        radius="md"
                        withBorder
                        shadow={isLatest ? "md" : "xs"}
                        style={{
                          borderColor: isLatest
                            ? "var(--mantine-color-blue-4)"
                            : undefined,
                          borderWidth: isLatest ? 2 : undefined,
                        }}
                      >
                        {/* Header */}
                        <Group justify="space-between" mb="sm">
                          <Group gap="xs">
                            <Badge
                              color={getActionColor(entry.action)}
                              variant={isLatest ? "filled" : "light"}
                              size="sm"
                            >
                              {entry.action || "Unknown Action"}
                            </Badge>
                            {isLatest && (
                              <Badge variant="dot" size="sm">
                                Latest
                              </Badge>
                            )}
                          </Group>
                          <Group gap="xs">
                            <IconClock size={14} />
                            <Text size="sm" c="dimmed">
                              {dateInfo.relative}
                            </Text>
                          </Group>
                        </Group>

                        {/* Author and timestamp */}
                        <Group gap="md" mb="sm">
                          <Group gap="xs">
                            <IconUser size={14} />
                            <Text size="sm" fw={500}>
                              {entry.author || "Unknown User"}
                            </Text>
                          </Group>
                          <Text size="xs" c="dimmed">
                            {dateInfo.date} at {dateInfo.time}
                          </Text>
                        </Group>

                        {/* Message */}
                        {entry.message && (
                          <Text
                            size="sm"
                            mb="sm"
                            style={{ fontStyle: "italic" }}
                          >
                            "{entry.message}"
                          </Text>
                        )}

                        {/* Changes */}
                        {entry.changes &&
                          Object.keys(entry.changes).length > 0 && (
                            <Box mb="sm">
                              <Text size="sm" fw={500} mb="xs">
                                Modified Fields:
                              </Text>
                              <Group gap="xs">
                                {Object.keys(entry.changes).map((field) => (
                                  <Badge
                                    key={field}
                                    size="sm"
                                    variant="outline"
                                    color="gray"
                                  >
                                    {field}
                                  </Badge>
                                ))}
                              </Group>
                            </Box>
                          )}

                        {/* Git info */}
                        {(entry.git_branch || entry.git_commit) && (
                          <Group
                            gap="md"
                            mt="sm"
                            pt="sm"
                            style={{
                              borderTop:
                                "1px solid var(--mantine-color-gray-3)",
                            }}
                          >
                            {entry.git_branch && (
                              <Group gap={4}>
                                <IconGitBranch size={12} />
                                <Text size="xs" c="dimmed">
                                  {entry.git_branch}
                                </Text>
                              </Group>
                            )}
                            {entry.git_commit && (
                              <Group gap={4}>
                                <IconGitCommit size={12} />
                                <Text size="xs" c="dimmed" ff="monospace">
                                  {entry.git_commit.substring(0, 7)}
                                </Text>
                              </Group>
                            )}
                          </Group>
                        )}
                      </Card>
                    </Timeline.Item>
                  );
                })}
              </Timeline>
            </Box>
          )}
        </ScrollArea>
      </Box>
    </Drawer>
  );
}
