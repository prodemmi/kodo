import { useRef, useState } from "react";
import { TextInput, ActionIcon, Group, Notification } from "@mantine/core";
import { IconEdit, IconCheck, IconX } from "@tabler/icons-react";
import { useNoteStore } from "../../../../../../../../../states/note.state";
import { RoleGuard } from "../../../../../../../../Investor";
import { useUpdateNote } from "../../../../../../../../../hooks/use-notes";
import { useTextWidth } from "../../../../../../../../../hooks/use-text-width";
import { getInputWidth } from "../../../../../../../../../utils/get-text-input";

export default function NoteTitle() {
  const selectedNote = useNoteStore((s) => s.selectedNote);
  const updateNote = useNoteStore((s) => s.updateNote);
  const [isEditing, setIsEditing] = useState(false);
  const [value, setValue] = useState(selectedNote?.title || "");
  const { mutate } = useUpdateNote();

  const handleSave = () => {
    if (!selectedNote) return;
    mutate(
      {
        ...selectedNote,
        title: value,
        id: selectedNote.id,
      },
      {
        onSuccess() {
          setIsEditing(false);
        },
      }
    );
  };

  const handleCancel = () => {
    setValue(selectedNote?.title || "");
    setIsEditing(false);
  };
  const inputWidth = getInputWidth(value);

  if (!selectedNote) return null;

  return (
    <Group gap="xs" align="center">
      {isEditing ? (
        <>
          <TextInput
            value={value}
            w={inputWidth}
            onChange={(e) => setValue(e.currentTarget.value)}
            size="md"
            autoFocus={isEditing}
            mb="xs"
            style={{ pointerEvents: isEditing ? "auto" : "none" }}
          />
          <ActionIcon
            variant="filled"
            color="green"
            size="sm"
            onClick={handleSave}
            disabled={!value}
          >
            <IconCheck size={14} />
          </ActionIcon>
          <ActionIcon
            variant="filled"
            color="red"
            size="sm"
            onClick={handleCancel}
          >
            <IconX size={14} />
          </ActionIcon>
        </>
      ) : (
        <>
          <TextInput
            defaultValue={value}
            size="md"
            w={inputWidth}
            variant="unstyled"
            autoFocus={false}
            mb="xs"
            px="md"
            styles={{
              input: {
                pointerEvents: isEditing ? "auto" : "none",
                border: "none",
                backgroundColor: "transparent",
              },
            }}
          />
          <RoleGuard.Consumer>
            <ActionIcon
              variant="filled"
              size="sm"
              onClick={() => setIsEditing(true)}
            >
              <IconEdit size={14} />
            </ActionIcon>
          </RoleGuard.Consumer>
        </>
      )}
    </Group>
  );
}
