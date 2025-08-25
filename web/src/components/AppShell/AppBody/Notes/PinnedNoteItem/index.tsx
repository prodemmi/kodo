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
  useMantineTheme,
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
import { RoleGuard } from "../../../../Investor";

interface Props {
  note: Note;
}

export default function PinnedNoteItem({ note }: Props) {
  const selectNote = useNoteStore((s) => s.selectNote);
  const selectedNote = useNoteStore((s) => s.selectedNote);
  const { primaryColor } = useMantineTheme();

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
              color={isSelected ? `var(--mantine-color-${primaryColor}-4)` : "var(--mantine-color-gray-4)"}
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

          <RoleGuard.Investor>
            <ActionIcon
              variant="subtle"
              size="xs"
              mr="xs"
              style={{ pointerEvents: "none" }}
            >
              <IconPinFilled size={12} />
            </ActionIcon>
          </RoleGuard.Investor>

          <RoleGuard.Consumer>
            <Group gap="xs">
              <ActionIcon variant="subtle" size="xs" onClick={handleTogglePin}>
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
                  <MenuItem leftSection={<IconTrash size={12} color="red" />}>
                    Delete
                  </MenuItem>
                </MenuDropdown>
              </Menu>
            </Group>
          </RoleGuard.Consumer>
        </Group>
      </UnstyledButton>
    </Box>
  );
}
