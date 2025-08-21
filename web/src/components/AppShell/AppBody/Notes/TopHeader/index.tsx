import { Group, Button } from "@mantine/core";
import { IconFolderPlus, IconPlus } from "@tabler/icons-react";
import { RoleGuard } from "../../../../Investor";
import {
  useNewFolderModalStore,
  useNewNoteModalStore,
} from "../../../../../states/note.state";

export default function TopHeader() {
  const openNewNoteModal = useNewNoteModalStore((s) => s.openModal);
  const setIsFolderModalOpen = useNewFolderModalStore((s) => s.openModal);

  return (
    <Group justify="flex-end" w="100%" px="sm">
      <RoleGuard.Consumer>
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
          leftSection={<IconPlus size={16} />}
          onClick={openNewNoteModal}
        >
          New Note
        </Button>
      </RoleGuard.Consumer>
    </Group>
  );
}
