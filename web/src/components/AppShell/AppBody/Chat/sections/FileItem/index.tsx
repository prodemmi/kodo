import {
  ActionIcon,
  Box,
  Group,
  HoverCard,
  Menu,
  MenuItem,
  MenuDropdown,
  MenuTarget,
  Text,
  UnstyledButton,
} from "@mantine/core";
import { useInViewport } from "@mantine/hooks";
import {
  IconFile,
  IconDots,
  IconFileText,
  IconCode,
  IconCopy,
  IconEye,
  IconDownload,
  IconTrash,
  IconPlus,
} from "@tabler/icons-react";
import { ProjectFile } from "../../../../../../types/chat";

function FileItem({
  item,
  level = 0,
  onClickItem,
  onAddToContext,
  onSelectCode,
  onCopyPath,
  onPreview,
  onDownload,
  onDelete,
}: {
  item: ProjectFile;
  level?: number;
  onClickItem: (item: ProjectFile) => void;
  onAddToContext?: (item: ProjectFile) => void;
  onSelectCode?: (item: ProjectFile) => void;
  onCopyPath?: (item: ProjectFile) => void;
  onPreview?: (item: ProjectFile) => void;
  onDownload?: (item: ProjectFile) => void;
  onDelete?: (item: ProjectFile) => void;
}) {
  const { ref, inViewport } = useInViewport();

  return (
    <Box ref={ref} mih={"28px"} w="100%" pr="12">
      {inViewport && (
        <HoverCard openDelay={800} closeDelay={200}>
          <HoverCard.Target>
            <Group
              gap={4}
              wrap="nowrap"
              py={2}
              pl={level * 20}
              justify="space-between"
              style={{
                cursor: "pointer",
                borderRadius: 4,
              }}
              w="100%"
              onClick={() => onClickItem(item)}
            >
              <Group w="100%">
                <IconFile size={14} style={{ flexShrink: 0 }} />
                <Text size="xs" truncate style={{ flex: 1 }}>
                  {item.name}
                </Text>
              </Group>

              <Menu position="right-start" withinPortal>
                <MenuTarget>
                  <ActionIcon
                    size={16}
                    variant="subtle"
                    color="gray"
                    onClick={(e) => e.stopPropagation()}
                    style={{ opacity: 0.6 }}
                  >
                    <IconPlus size={10} />
                  </ActionIcon>
                </MenuTarget>

                <MenuDropdown>
                  {onAddToContext && (
                    <MenuItem
                      leftSection={<IconFileText size={14} />}
                      onClick={() => onAddToContext(item)}
                    >
                      Add File to Context
                    </MenuItem>
                  )}

                  {onSelectCode && (
                    <MenuItem
                      leftSection={<IconCode size={14} />}
                      onClick={() => onSelectCode(item)}
                    >
                      Select Code to Add
                    </MenuItem>
                  )}

                  {onCopyPath && (
                    <MenuItem
                      leftSection={<IconCopy size={14} />}
                      onClick={() => onCopyPath(item)}
                    >
                      Copy Path
                    </MenuItem>
                  )}

                  {onPreview && (
                    <MenuItem
                      leftSection={<IconEye size={14} />}
                      onClick={() => onPreview(item)}
                    >
                      Quick Preview
                    </MenuItem>
                  )}

                  {onDownload && (
                    <MenuItem
                      leftSection={<IconDownload size={14} />}
                      onClick={() => onDownload(item)}
                    >
                      Download
                    </MenuItem>
                  )}

                  {onDelete && (
                    <>
                      <Menu.Divider />
                      <MenuItem
                        leftSection={<IconTrash size={14} color="red" />}
                        onClick={() => onDelete(item)}
                      >
                        Delete
                      </MenuItem>
                    </>
                  )}
                </MenuDropdown>
              </Menu>
            </Group>
          </HoverCard.Target>

          <HoverCard.Dropdown>
            <Text size="xs" c="dimmed">
              {item.path}
            </Text>
            {item.size && (
              <Text size="xs" c="dimmed">
                {(item.size / 1024).toFixed(1)} KB
              </Text>
            )}
          </HoverCard.Dropdown>
        </HoverCard>
      )}
    </Box>
  );
}

// Even more compact version
function FileItemCompact({
  item,
  level = 0,
  onClickItem,
  onAddToContext,
  onSelectCode,
  onPreview,
  onCopyPath,
  onDownload,
  onDelete,
}: {
  item: ProjectFile;
  level?: number;
  onClickItem: (item: ProjectFile) => void;
  onAddToContext?: (item: ProjectFile) => void;
  onSelectCode?: (item: ProjectFile) => void;
  onPreview?: (item: ProjectFile) => void;
  onCopyPath?: (item: ProjectFile) => void;
  onDownload?: (item: ProjectFile) => void;
  onDelete?: (item: ProjectFile) => void;
}) {
  const { ref, inViewport } = useInViewport();

  return (
    <div ref={ref} style={{ minHeight: 24 }}>
      {inViewport && (
        <Group
          gap={4}
          wrap="nowrap"
          pl={4 + level * 10}
          pr={4}
          py={1}
          className="hover:bg-gray-50 cursor-pointer rounded-sm"
          onClick={() => onClickItem(item)}
        >
          <IconFile size={12} style={{ flexShrink: 0, opacity: 0.7 }} />
          <Text size="xs" truncate style={{ flex: 1 }}>
            {item.name}
          </Text>

          <Group
            gap={2}
            style={{ opacity: 0 }}
            className="group-hover:opacity-100"
          >
            {onAddToContext && (
              <ActionIcon
                size={16}
                variant="subtle"
                color="gray"
                onClick={(e) => {
                  e.stopPropagation();
                  onAddToContext(item);
                }}
                title="Add to Context"
              >
                <IconFileText size={10} />
              </ActionIcon>
            )}

            {onSelectCode && (
              <ActionIcon
                size={16}
                variant="subtle"
                color="gray"
                onClick={(e) => {
                  e.stopPropagation();
                  onSelectCode(item);
                }}
                title="Select Code"
              >
                <IconCode size={10} />
              </ActionIcon>
            )}

            {onPreview && (
              <ActionIcon
                size={16}
                variant="subtle"
                color="gray"
                onClick={(e) => {
                  e.stopPropagation();
                  onPreview(item);
                }}
                title="Preview"
              >
                <IconEye size={10} />
              </ActionIcon>
            )}

            <Menu position="right-start" withinPortal>
              <MenuTarget>
                <ActionIcon
                  size={16}
                  variant="subtle"
                  color="gray"
                  onClick={(e) => e.stopPropagation()}
                  title="More options"
                >
                  <IconPlus size={10} />
                </ActionIcon>
              </MenuTarget>

              <MenuDropdown>
                {onCopyPath && (
                  <MenuItem
                    leftSection={<IconCopy size={14} />}
                    onClick={() => onCopyPath(item)}
                  >
                    Copy Path
                  </MenuItem>
                )}

                {onDownload && (
                  <MenuItem
                    leftSection={<IconDownload size={14} />}
                    onClick={() => onDownload(item)}
                  >
                    Download
                  </MenuItem>
                )}

                {onDelete && (
                  <>
                    <Menu.Divider />
                    <MenuItem
                      leftSection={<IconTrash size={14} color="red" />}
                      onClick={() => onDelete(item)}
                    >
                      Delete
                    </MenuItem>
                  </>
                )}
              </MenuDropdown>
            </Menu>
          </Group>
        </Group>
      )}
    </div>
  );
}

export { FileItem, FileItemCompact };
