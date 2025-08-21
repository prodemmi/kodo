import {
  Text,
  Box,
  Group,
  ActionIcon,
  UnstyledButton,
  Badge,
  ScrollArea,
} from "@mantine/core";
import { IconFolderPlus, IconFileText } from "@tabler/icons-react";
import { RoleGuard } from "../../../../Investor";
import FolderItem from "./FolderItem";
import {
  useNewFolderModalStore,
  useNoteStore,
} from "../../../../../states/note.state";
import { Folder } from "../../../../../types/note";
import { useMediaQuery } from "@mantine/hooks";
import { useFolderTree } from "../../../../../hooks/use-notes";

export default function Folders() {
  const storeNotes = useNoteStore((s) => s.notes);
  const selectFolder = useNoteStore((s) => s.selectFolder);
  const selectedFolder = useNoteStore((s) => s.selectedFolder);
  const setIsFolderModalOpen = useNewFolderModalStore((s) => s.openModal);
  const isSmall = useMediaQuery("(max-width: 920px)");
  const { data: folders } = useFolderTree();

  return (
    <ScrollArea h={isSmall ? undefined : "100%"} miw={isSmall ? "100%" : "40%"} py="xs">
      <Box px="2" h={isSmall ? undefined : "100%"}>
        <Group justify="space-between" mb="sm">
          <Text size="sm" fw={600} c="dimmed">
            FOLDERS
          </Text>
          <RoleGuard.Consumer>
            <ActionIcon
              variant="subtle"
              size="sm"
              onClick={() => setIsFolderModalOpen()}
            >
              <IconFolderPlus size={14} />
            </ActionIcon>
          </RoleGuard.Consumer>
        </Group>

        <UnstyledButton
          p="2"
          style={{
            width: "100%",
          }}
          onClick={() => selectFolder(null)}
        >
          <Group gap="xs">
            <IconFileText
              size={16}
              color={selectedFolder === null ? "#339af0" : "#868e96"}
            />
            <Text size="sm" fw={selectedFolder === null ? 600 : 400}>
              All Notes
            </Text>
            <Badge size="xs" variant="light" color="gray">
              {storeNotes.length}
            </Badge>
          </Group>
        </UnstyledButton>

        {folders &&
          folders.count > 0 &&
          folders.tree.map((folder: Folder) => (
            <FolderItem key={folder.id} folder={folder} level={0} />
          ))}
      </Box>
    </ScrollArea>
  );
}
