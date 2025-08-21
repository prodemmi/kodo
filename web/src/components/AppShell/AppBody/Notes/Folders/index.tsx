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

export default function Folders() {
  const storeNotes = useNoteStore((s) => s.notes);
  const folders = useNoteStore((s) => s.folders);
  const selectFolder = useNoteStore((s) => s.selectFolder);
  const selectedFolder = useNoteStore((s) => s.selectedFolder);
  const setIsFolderModalOpen = useNewFolderModalStore((s) => s.openModal);
  const isSmall = useMediaQuery("(max-width: 920px)");

  const getFolderHierarchy = () => {
    const folderMap = new Map();
    folders.forEach((folder) =>
      folderMap.set(folder.id, { ...folder, children: [] })
    );

    const rootFolders: Folder[] = [];
    folders.forEach((folder) => {
      if (folder.parentId) {
        const parent = folderMap.get(folder.parentId);
        if (parent) {
          parent.children.push(folderMap.get(folder.id));
        }
      } else {
        rootFolders.push(folderMap.get(folder.id));
      }
    });

    return rootFolders;
  };

  return (
    <ScrollArea h={isSmall ? undefined : "100%"} w={isSmall ? "100%" : "40%"}>
      <Box px="xs" h={isSmall ? undefined : "100%"}>
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
          p="6"
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

        {getFolderHierarchy().map((folder: Folder) => (
          <FolderItem key={folder.id} folder={folder} level={0} />
        ))}
      </Box>
    </ScrollArea>
  );
}
