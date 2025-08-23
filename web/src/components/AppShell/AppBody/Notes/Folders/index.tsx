import {
  Text,
  Box,
  Group,
  ActionIcon,
  UnstyledButton,
  Badge,
  ScrollArea,
  LoadingOverlay,
  Title,
  Divider,
} from "@mantine/core";
import {
  IconFolderPlus,
  IconFileText,
  IconPinFilled,
} from "@tabler/icons-react";
import { RoleGuard } from "../../../../Investor";
import FolderItem from "./FolderItem";
import {
  useNewFolderModalStore,
  useNoteStore,
} from "../../../../../states/note.state";
import { Folder, Note } from "../../../../../types/note";
import { useMediaQuery } from "@mantine/hooks";
import { useFolderTree } from "../../../../../hooks/use-notes";
import { useEffect, useMemo } from "react";
import PinnedNoteItem from "../PinnedNoteItem";

export default function Folders() {
  const storeNotes = useNoteStore((s) => s.notes);
  const folderTree = useNoteStore((s) => s.folderTree);
  const setFolderTree = useNoteStore((s) => s.setFolderTree);
  const selectFolder = useNoteStore((s) => s.selectFolder);
  const selectedFolder = useNoteStore((s) => s.selectedFolder);
  const setIsFolderModalOpen = useNewFolderModalStore((s) => s.openModal);
  const isSmall = useMediaQuery("(max-width: 920px)");
  const { data: remoteFolderTree, isError, isPending, isLoading } = useFolderTree();

  // Get pinned and unpinned notes
  const pinnedNotes = useMemo(
    () => storeNotes.filter((note) => note.pinned),
    [storeNotes]
  );

  useEffect(() => {
    if (remoteFolderTree && !isError && !isLoading) {
      setFolderTree(remoteFolderTree.tree);
    }
  }, [remoteFolderTree, isError, isLoading, setFolderTree]);

  if (isLoading) return <LoadingOverlay />;

  return (
    <ScrollArea h={isSmall ? undefined : "100%"} w="40%" p="sm" pt="md">
      <Box h={isSmall ? undefined : "100%"}>
        <Group justify="space-between" mb="sm">
          <Title size="h6">Folders</Title>
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

        {/* Pinned Notes Section */}
        {pinnedNotes.length > 0 && (
          <>
            <Group gap="xs" mb="xs">
              <IconPinFilled size={16} color="#339af0" />
              <Text size="sm" fw={500} c="blue">
                Pinned Notes
              </Text>
              <Badge size="xs" variant="light" color="blue">
                {pinnedNotes.length}
              </Badge>
            </Group>

            {pinnedNotes.map((note: Note) => (
              <PinnedNoteItem key={note.id} note={note} />
            ))}

            <Divider my="md" />
          </>
        )}

        {/* All Notes Button */}
        <UnstyledButton
          style={{ width: "100%" }}
          onClick={() => selectFolder(null)}
        >
          <Group gap="xs">
            <IconFileText
              size={16}
              color={selectedFolder === null ? "#339af0" : "#868e96"}
            />
            <Text size="sm" fw={400}>
              All Notes
            </Text>
            <Badge size="xs" variant="light" color="gray">
              {storeNotes.length}
            </Badge>
          </Group>
        </UnstyledButton>

        {/* Folder Tree */}
        {folderTree &&
          folderTree.length > 0 &&
          folderTree.map((folder: Folder) => (
            <FolderItem key={folder.id} folder={folder} level={0} />
          ))}
      </Box>
    </ScrollArea>
  );
}
