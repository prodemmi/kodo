import {
  Box,
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
} from "@mantine/core";
import { IconGitBranch, IconEdit, IconCheck, IconX } from "@tabler/icons-react";
import { useState } from "react";
import { Editor } from "@tiptap/react";
import { useNoteStore } from "../../../../../../../states/note.state";
import { RoleGuard } from "../../../../../../Investor";
import { tagColors } from "../../../constants";
import NoteTitle from "./sections/NoteTitle";
import { useUpdateNote } from "../../../../../../../hooks/use-notes";

interface Props {
  editor: Editor;
}

export default function NoteInfo({ editor }: Props) {
  const allTags = useNoteStore((s) => s.tags);
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
            updateNote(selectedNote?.id!, {
              ...selectedNote,
              content: editor?.getHTML() || "",
              updatedAt: new Date(),
            });
            setIsEditingTags(false);
            setError(null);
          },
        }
      );
  };

  return (
    selectedNote && (
      <Stack
        px="md"
        pb="md"
        gap="xs"
        style={{
          borderBottom: "1px solid var(--mantine-color-gray-8)",
          minHeight: "auto",
        }}
      >
        <Group justify="space-between">
          <div style={{ flex: 1 }}>
            <NoteTitle />
            <Group gap="sm">
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
          </div>

          <RoleGuard.Consumer>
            <Group gap="sm">
              {!isEditingNote ? (
                <Button
                  leftSection={<IconEdit size={16} />}
                  onClick={() => setIsEditingNote(true)}
                >
                  Edit
                </Button>
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

        {/* Tags and Category */}
        <Group gap="sm" h="32px">
          <Badge
            color={getCategoryColor(selectedNote.category)}
            variant="light"
          >
            {selectedNote.category}
          </Badge>

          <Divider orientation="vertical" />

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
