import {
  UnstyledButton,
  Group,
  ActionIcon,
  Badge,
  Text,
  Collapse,
} from "@mantine/core";
import {
  IconFolder,
  IconChevronDown,
  IconChevronRight,
} from "@tabler/icons-react";
import { useNoteStore } from "../../../../../../states/note.state";
import { Folder } from "../../../../../../types/note";

interface Props {
  folder: Folder;
  level: number;
}

export default function FolderItem({ folder, level = 0 }: Props) {
  const storeNotes = useNoteStore((s) => s.notes);
  const folders = useNoteStore((s) => s.folders);
  const selectedFolder = useNoteStore((s) => s.selectedFolder);
  const selectFolder = useNoteStore((s) => s.selectFolder);
  const toggleFolder = useNoteStore((s) => s.toggleFolder);

  const folderNotes = storeNotes.filter(
    (note: any) => note.folderId === folder.id
  );
  const hasNotes = folderNotes.length > 0;
  const folderChildren = folders.filter((f) => f.parentId === folder.id);
  const hasChildren = folderChildren?.length > 0;
  const isSelected = selectedFolder?.id === folder.id;

  return (
    <div key={folder.id}>
      <UnstyledButton
        p="6"
        style={{
          width: "100%",
          paddingLeft: `${13 + level * 20}px`,
        }}
        onClick={(e) => {
          e.stopPropagation();
          selectFolder(folder);
          toggleFolder(folder.id);
        }}
      >
        <Group gap="xs">
          <IconFolder size={16} color={isSelected ? "#339af0" : "#868e96"} />

          <Text size="sm" fw={isSelected ? 600 : 400}>
            {folder.name}
          </Text>

          {hasChildren && (
            <ActionIcon variant="transparent" size="xs">
              {folder.expanded ? (
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
      </UnstyledButton>

      <Collapse in={folder.expanded} pl="md">
        {folderChildren?.map((child: Folder) => (
          <FolderItem key={child.id} folder={child} level={level + 1} />
        ))}
      </Collapse>
    </div>
  );
}
