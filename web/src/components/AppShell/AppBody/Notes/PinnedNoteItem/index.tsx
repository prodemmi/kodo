import {
  UnstyledButton,
  Group,
  Text,
  ActionIcon,
  Menu,
  MenuTarget,
  MenuDropdown,
  MenuItem,
  Box,
} from "@mantine/core";
import {
  IconFileText,
  IconPin,
  IconPinFilled,
  IconDotsVertical,
  IconTrash,
  IconEdit,
} from "@tabler/icons-react";
import { useNoteStore } from "../../../../../states/note.state";
import { Note } from "../../../../../types/note";
import { useUpdateNote } from "../../../../../hooks/use-notes";

interface Props {
  note: Note;
}

export default function PinnedNoteItem({ note }: Props) {
  const selectNote = useNoteStore((s) => s.selectNote);
  const selectedNote = useNoteStore((s) => s.selectedNote);

  const isSelected = selectedNote?.id === note.id;

  const { mutate: updateNote } = useUpdateNote();

  const togglePinNote = (note: Note) => {
    updateNote({
      ...note,
      id: note.id,
      pinned: !note.pinned,
    });
  };

  const handleTogglePin = (e: React.MouseEvent) => {
    e.stopPropagation();
    togglePinNote(note);
  };

  const handleSelectNote = () => {
    selectNote(note);
  };

  return (
    <Box mb="xs">
      <UnstyledButton style={{ width: "100%" }} onClick={handleSelectNote}>
        <Group justify="space-between" w="100%">
          <Group gap="xs">
            <IconFileText
              size={16}
              color={isSelected ? "#339af0" : "#868e96"}
            />
            <Text
              size="sm"
              fw={400}
              style={{
                maxWidth: 150,
                overflow: "hidden",
                textOverflow: "ellipsis",
                whiteSpace: "nowrap",
              }}
            >
              {note.title}
            </Text>
          </Group>

          <Group gap="xs">
            <ActionIcon
              variant="subtle"
              size="xs"
              onClick={handleTogglePin}
              color="blue"
            >
              <IconPinFilled size={12} />
            </ActionIcon>

            <Menu position="right-start">
              <MenuTarget>
                <ActionIcon variant="subtle" size="xs">
                  <IconDotsVertical size={12} />
                </ActionIcon>
              </MenuTarget>
              <MenuDropdown>
                <MenuItem
                  onClick={handleTogglePin}
                  leftSection={<IconPin size={12} />}
                >
                  Unpin Note
                </MenuItem>
                <MenuItem leftSection={<IconEdit size={12} />}>Edit</MenuItem>
                <MenuItem color="red" leftSection={<IconTrash size={12} />}>
                  Delete
                </MenuItem>
              </MenuDropdown>
            </Menu>
          </Group>
        </Group>
      </UnstyledButton>
    </Box>
  );
}
