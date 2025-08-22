import {
  UnstyledButton,
  Group,
  ActionIcon,
  Badge,
  Text,
  Collapse,
  Menu,
  MenuTarget,
  MenuDropdown,
  MenuItem,
  Box,
} from "@mantine/core";
import {
  IconFolder,
  IconChevronDown,
  IconChevronRight,
  IconTrash,
  IconEdit,
  IconDotsVertical,
} from "@tabler/icons-react";
import {
  useDeleteModalStore,
  useNewFolderModalStore,
  useNoteStore,
} from "../../../../../../states/note.state";
import { Folder } from "../../../../../../types/note";
import { useState } from "react";

interface Props {
  folder: Folder;
  level: number;
}

export default function FolderItem({ folder, level = 0 }: Props) {
  const storeNotes = useNoteStore((s) => s.notes);
  const folders = useNoteStore((s) => s.folders);
  const selectedFolder = useNoteStore((s) => s.selectedFolder);
  const selectFolder = useNoteStore((s) => s.selectFolder);
  const openForFolder = useDeleteModalStore((s) => s.openForFolder);
  const openEditModal = useNewFolderModalStore((s) => s.openEditModal);
  const [collapsed, setCollapsed] = useState(false);

  const folderNotes = storeNotes.filter(
    (note: any) => note.folderId === folder.id
  );
  const hasNotes = folderNotes.length > 0;
  const folderChildren = folders.filter((f) => f.parentId === folder.id);
  const hasChildren = folderChildren?.length > 0;
  const isSelected = selectedFolder?.id === folder.id;

  const openDeleteModal = () => {
    if (folder) openForFolder(folder);
  };

  return (
    folder && (
      <Box key={folder.id}>
        <UnstyledButton
          style={{
            width: "100%",
          }}
        >
          <Group justify="space-between" w="100%">
            <Group
              gap="xs"
              onClick={(e) => {
                e.stopPropagation();
                selectFolder(folder);
                setCollapsed((o) => !o);
              }}
            >
              <IconFolder
                size={16}
                color={isSelected ? "#339af0" : "#868e96"}
              />

              <Text size="sm" fw={400}>
                {folder.name}
              </Text>

              {hasChildren && (
                <ActionIcon variant="transparent" size="xs">
                  {collapsed ? (
                    <IconChevronDown size={12} />
                  ) : (
                    <IconChevronRight size={12} />
                  )}
                </ActionIcon>
              )}

              {hasNotes && (
                <Badge size="xs" variant="light" color="gray">
                  {folderNotes.length}
                </Badge>
              )}
            </Group>
            <Menu position="right-start">
              <MenuTarget>
                <ActionIcon>
                  <IconDotsVertical size={14} />
                </ActionIcon>
              </MenuTarget>
              <MenuDropdown>
                <MenuItem
                  onClick={() => openEditModal(folder.id)}
                  leftSection={<IconEdit size={12} />}
                >
                  Edit
                </MenuItem>
                <MenuItem
                  color="red"
                  onClick={openDeleteModal}
                  leftSection={<IconTrash size={12} />}
                >
                  Delete
                </MenuItem>
              </MenuDropdown>
            </Menu>
          </Group>
        </UnstyledButton>

        <Collapse in={collapsed} pl="md">
          {folderChildren?.map((child: Folder) => (
            <FolderItem key={child.id} folder={child} level={level + 1} />
          ))}
        </Collapse>
      </Box>
    )
  );
}
