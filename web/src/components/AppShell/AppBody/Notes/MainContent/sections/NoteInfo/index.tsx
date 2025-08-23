import {
  Text,
  Group,
  Avatar,
  Divider,
  Button,
  Badge,
  ActionIcon,
  TagsInput,
  Alert,
  Stack,
  Menu,
  MenuTarget,
  MenuDropdown,
  MenuItem,
} from "@mantine/core";
import {
  IconGitBranch,
  IconEdit,
  IconCheck,
  IconX,
  IconHistory,
} from "@tabler/icons-react";
import { useState } from "react";
import { Editor } from "@tiptap/react";
import {
  useNoteHistoryModalStore,
  useNoteStore,
} from "../../../../../../../states/note.state";
import { RoleGuard } from "../../../../../../Investor";
import { categories, tagColors } from "../../../constants";
import { useUpdateNote } from "../../../../../../../hooks/use-notes";
import NoteTitle from "./sections/NoteTitle";

interface Props {
  editor: Editor;
}

export default function NoteInfo({ editor }: Props) {
  const allTags = useNoteStore((s) => s.tags);
  const openHistoryForNote = useNoteHistoryModalStore((s) => s.openForNote);
  const isEditingTags = useNoteStore((s) => s.isEditingTags);
  const setIsEditingTags = useNoteStore((s) => s.setIsEditingTags);
  const setIsEditingNote = useNoteStore((s) => s.setIsEditingNote);
  const selectedNote = useNoteStore((s) => s.selectedNote);
  const tempTags = useNoteStore((s) => s.tempTags);
  const setTempTags = useNoteStore((s) => s.setTempTags);
  const isEditingNote = useNoteStore((s) => s.isEditingNote);
  const updateNote = useNoteStore((s) => s.updateNote);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { mutate } = useUpdateNote();

  const handleTagsEdit = () => {
    if (selectedNote) setTempTags([...selectedNote.tags]);
    setIsEditingTags(true);
  };

  const handleTagsSave = () => {
    if (selectedNote)
      mutate(
        {
          ...selectedNote,
          tags: tempTags,
          id: selectedNote?.id!,
        },
        {
          onSuccess() {
            updateNote(selectedNote?.id!, {
              ...selectedNote,
              tags: tempTags,
              updatedAt: new Date(),
            });
            setIsEditingTags(false);
            setError(null);
          },
        }
      );
  };

  const updateCategory = (value: string) => {
    if (selectedNote)
      mutate(
        {
          ...selectedNote,
          category: value,
          id: selectedNote?.id!,
        },
        {
          onSuccess() {
            updateNote(selectedNote?.id!, {
              ...selectedNote,
              category: value,
              updatedAt: new Date(),
            });
            setError(null);
          },
        }
      );
  };

  const getCategoryColor = (category: any) => {
    const colors: any = {
      technical: "blue",
      meeting: "purple",
      idea: "green",
      documentation: "orange",
      "bug-analysis": "red",
      review: "cyan",
    };
    return colors[category] || "gray";
  };

  const formatDate = (date: any) => {
    return new Intl.DateTimeFormat("en-US", {
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    }).format(new Date(date));
  };

  const onSubmit = () => {
    if (selectedNote)
      mutate(
        {
          ...selectedNote,
          content: editor?.getHTML() || "",
          id: selectedNote?.id!,
        },
        {
          onSuccess() {
            setIsEditingNote(false);
            setError(null);
          },
        }
      );
  };

  return (
    selectedNote && (
      <Stack
        px="md"
        pb="xs"
        gap="sm"
        style={{
          minHeight: "auto",
        }}
      >
        <Group justify="space-between">
          <NoteTitle />
          <RoleGuard.Investor>
            <Button
              leftSection={<IconHistory size={16} />}
              onClick={() => openHistoryForNote(selectedNote)}
            >
              History
            </Button>
          </RoleGuard.Investor>
          <RoleGuard.Consumer>
            <Group gap="sm">
              {!isEditingNote ? (
                <Group>
                  <Button
                    leftSection={<IconHistory size={16} />}
                    onClick={() => openHistoryForNote(selectedNote)}
                  >
                    History
                  </Button>
                  <Button
                    leftSection={<IconEdit size={16} />}
                    onClick={() => setIsEditingNote(true)}
                  >
                    Edit
                  </Button>
                </Group>
              ) : (
                <Group gap="sm">
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={() => {
                      setIsEditingNote(false);
                      editor?.commands.setContent(selectedNote.content);
                      setError(null);
                    }}
                  >
                    Cancel
                  </Button>
                  <Button size="sm" onClick={onSubmit} loading={loading}>
                    Save
                  </Button>
                </Group>
              )}
            </Group>
          </RoleGuard.Consumer>
        </Group>

        <Group justify="space-between">
          <Group gap="sm" align="flex-start">
            <Group gap="xs">
              <Avatar size={24} color="blue">
                {selectedNote.author.charAt(0).toUpperCase()}
              </Avatar>
              <div>
                <Text size="sm" fw={500}>
                  {selectedNote.author}
                </Text>
                <Text size="xs" c="dimmed">
                  Created {formatDate(selectedNote.createdAt)}
                  {selectedNote.updatedAt > selectedNote.createdAt &&
                    ` • Updated ${formatDate(selectedNote.updatedAt)}`}
                </Text>
              </div>
            </Group>
            <Divider orientation="vertical" />
            <Group gap="xs">
              <IconGitBranch size={16} color="#868e96" />
              <Text size="sm" c="dimmed">
                {selectedNote.gitBranch}
              </Text>
              <Text size="xs" c="dimmed">
                ({selectedNote.gitCommit})
              </Text>
            </Group>
          </Group>
        </Group>

        {/* Tags and Category */}
        <Group gap="sm" h="32px">
          <RoleGuard.Consumer>
            <Menu>
              <MenuTarget>
                <Badge color={getCategoryColor(selectedNote.category)} pr={0}>
                  <Group align="center" gap="0" justify="space-between">
                    <Text size="xs">{selectedNote.category}</Text>
                    <ActionIcon m={0}>
                      <IconEdit size={12} />
                    </ActionIcon>
                  </Group>
                </Badge>
              </MenuTarget>
              <MenuDropdown>
                {categories.map((cat) => (
                  <MenuItem
                    value={cat.value}
                    leftSection={
                      cat.value === selectedNote.category && (
                        <IconCheck size={12} />
                      )
                    }
                    onClick={() => updateCategory(cat.value)}
                  >
                    {cat.label}
                  </MenuItem>
                ))}
              </MenuDropdown>
            </Menu>
          </RoleGuard.Consumer>

          <RoleGuard.Investor>
            <Badge color={getCategoryColor(selectedNote.category)}>
              <Text size="xs">{selectedNote.category}</Text>
            </Badge>
          </RoleGuard.Investor>

          {selectedNote.tags && selectedNote.tags.length && (
            <Divider orientation="vertical" />
          )}

          {!isEditingTags ? (
            <Group gap="xs">
              {selectedNote.tags.map((tag: any) => (
                <Badge
                  key={tag}
                  color={tagColors[tag] || "gray"}
                  variant="outline"
                  size="sm"
                >
                  {tag}
                </Badge>
              ))}
              <RoleGuard.Consumer>
                {!selectedNote.tags ||
                  (selectedNote.tags.length === 0 && (
                    <Text size="xs">Add tags</Text>
                  ))}
                <ActionIcon variant="subtle" size="sm" onClick={handleTagsEdit}>
                  <IconEdit size={12} />
                </ActionIcon>
              </RoleGuard.Consumer>
            </Group>
          ) : (
            <Group gap="xs" style={{ flex: 1 }} align="center">
              <TagsInput
                value={tempTags}
                onChange={setTempTags}
                data={allTags()}
                placeholder="Add tags..."
                size="xs"
                style={{ flex: 1 }}
                styles={{
                  pill: { size: "sm", variant: "outline" },
                  input: {
                    padding: "0",
                    paddingTop: "4px",
                    border: "none",
                    backgroundColor: "transparent",
                  },
                }}
              />
              <ActionIcon
                variant="filled"
                color="green"
                size="sm"
                onClick={handleTagsSave}
              >
                <IconCheck size={12} />
              </ActionIcon>
              <ActionIcon
                variant="filled"
                color="red"
                size="sm"
                onClick={() => {
                  setTempTags([]);
                  setIsEditingTags(false);
                  setError(null);
                }}
              >
                <IconX size={12} />
              </ActionIcon>
            </Group>
          )}
        </Group>
        {error && (
          <Alert color="red" variant="light">
            {error}
          </Alert>
        )}
      </Stack>
    )
  );
}
