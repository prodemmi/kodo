import { Group, Button } from "@mantine/core";
import { IconFolderPlus, IconPlus } from "@tabler/icons-react";
import { RoleGuard } from "../../../../Investor";
import {
  useNewFolderModalStore,
  useNewNoteModalStore,
  useNoteStore,
} from "../../../../../states/note.state";
import NoteTitle from "../MainContent/sections/NoteInfo/sections/NoteTitle";

export default function TopHeader() {
  const selectedNote = useNoteStore((s) => s.selectedNote);
  const openNewNoteModal = useNewNoteModalStore((s) => s.openModal);
  const setIsFolderModalOpen = useNewFolderModalStore((s) => s.openModal);

  return (
    <Group justify={selectedNote ? "space-between" : "flex-end"} w="100%" px="sm">
      <NoteTitle />
      <RoleGuard.Consumer>
        <Group>
          <Button
            variant="light"
            size="xs"
            leftSection={<IconFolderPlus size={16} />}
            onClick={() => setIsFolderModalOpen()}
          >
            New Folder
          </Button>
          <Button
            size="xs"
            bg="green"
            leftSection={<IconPlus size={16} />}
            onClick={openNewNoteModal}
          >
            New Note
          </Button>
        </Group>
      </RoleGuard.Consumer>
    </Group>
  );
}
