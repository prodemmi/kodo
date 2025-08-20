import { useEffect, useState } from "react";
import { useChatFiles } from "../../../../../../hooks/use-chat";
import { ProjectFile } from "../../../../../../types/chat";
import {
  Text,
  Collapse,
  Group,
  HoverCard,
  HoverCardDropdown,
  HoverCardTarget,
  Loader,
  UnstyledButton,
} from "@mantine/core";
import {
  IconChevronDown,
  IconChevronRight,
  IconFolder,
  IconFolderOpen,
} from "@tabler/icons-react";
import { FileItem } from "../FileItem";

export default function DirectoryItem({
  item,
  isOpened,
  level = 0,
  searchMode = false,
  onClickItem,
}: {
  item: ProjectFile;
  isOpened: boolean;
  searchMode: boolean;
  level?: number;
  onClickItem: (item: ProjectFile) => void;
}) {
  const [files, setFiles] = useState<ProjectFile[]>([]);
  const [loading, setLoading] = useState(false);
  const [opened, setOpened] = useState(isOpened);

  const {
    data: directoryFiles,
    isLoading,
    isError,
  } = useChatFiles(item.path, null, !searchMode && loading);

  useEffect(() => {
    if (!isError && !isLoading && loading) {
      const timeout = setTimeout(() => {
        setFiles(directoryFiles || []);
        setOpened(true); // open after loading
        setLoading(false); // reset loading trigger
      }, 120);

      return () => clearTimeout(timeout); // cancel previous timeout if effect re-runs
    }
  }, [directoryFiles, isLoading, isError, loading]);

  return (
    <>
      <HoverCard shadow="md" openDelay={1000} position="top">
        <HoverCardTarget>
          <UnstyledButton
            onClick={() => {
              if (!opened && files.length === 0) {
                // first click on unopened folder → load
                setLoading(true);
              } else {
                // toggle open/close freely
                setOpened((o) => !o);
              }
            }}
            w="100%"
            pr={4}
            py={2}
            pl={4 + level * 16}
          >
            <Group gap="xs" wrap="nowrap">
              {opened ? (
                <IconChevronDown size={14} />
              ) : (
                <IconChevronRight size={14} />
              )}
              {opened ? <IconFolderOpen size={16} /> : <IconFolder size={16} />}
              <Text size="sm" truncate>
                {item.name}
              </Text>
              {loading && <Loader size="8" ml="auto" mr="sm" />}
            </Group>
          </UnstyledButton>
        </HoverCardTarget>
        <HoverCardDropdown>
          <Text size="xs">{item.path}</Text>
        </HoverCardDropdown>
      </HoverCard>

      <Collapse in={opened}>
        {isLoading ? (
          <Text size="sm" c="dimmed" pl={4 + (level + 1) * 16}>
            Loading files...
          </Text>
        ) : (
          files?.map((child) =>
            child.type === "folder" ? (
              <DirectoryItem
                key={child.id}
                item={child}
                isOpened={false}
                level={level + 1}
                searchMode={searchMode}
                onClickItem={onClickItem}
              />
            ) : (
              <FileItem
                key={child.id}
                item={child}
                level={level + 1}
                onClickItem={onClickItem}
              />
            )
          )
        )}
      </Collapse>
    </>
  );
}
