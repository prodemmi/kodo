import { useSortable } from "@dnd-kit/sortable";
import { Item } from "../../../../../../types/item";
import {
  Text,
  Badge,
  Card,
  Group,
  Select,
  ActionIcon,
  Stack,
  Divider,
  Box,
} from "@mantine/core";
import { CSS } from "@dnd-kit/utilities";
import { IconCode, IconEye, IconGripVertical } from "@tabler/icons-react";
import { useOpenFile, useUpdateItem } from "../../../../../../hooks/use-items";
import { useCallback } from "react";
import { RoleGuard } from "../../../../../Investor";
import { useElementSize, useMergedRef } from "@mantine/hooks";

export default function SortableTask({
  item,
  onItemClick,
}: {
  item: Item;
  onItemClick: (item: Item) => void;
}) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
    isSorting,
  } = useSortable({
    id: item.id,
    transition: {
      duration: 200,
      easing: "cubic-bezier(0.25, 1, 0.5, 1)",
    },
  });
  const { ref, width } = useElementSize();
  const cardRef = useMergedRef(setNodeRef, ref);
  const { mutate } = useUpdateItem();

  const style = {
    transform: CSS.Transform.toString(transform),
    transition: isSorting ? transition : undefined,
    zIndex: isDragging ? 999 : "auto",
  };

  const getTypeColor = (type: string) => {
    switch (type) {
      // Code Quality & Maintenance
      case "REFACTOR":
        return "blue";
      case "OPTIMIZE":
        return "teal";
      case "CLEANUP":
        return "gray";
      case "DEPRECATED":
        return "orange";

      // Bug Fixes & Issues
      case "BUG":
        return "red";
      case "FIXME":
        return "pink";
      case "HACK":
        return "yellow";

      // Features & Enhancements
      case "TODO":
        return "indigo";
      case "FEATURE":
        return "green";
      case "ENHANCE":
        return "cyan";

      // Documentation & Testing
      case "DOC":
        return "violet";
      case "TEST":
        return "lime";
      case "EXAMPLE":
        return "grape";

      // Security & Compliance
      case "SECURITY":
        return "darkred";
      case "COMPLIANCE":
        return "brown";

      // Technical Debt & Architecture
      case "DEBT":
        return "darkorange";
      case "ARCHITECTURE":
        return "navy";

      // Operations & Infrastructure
      case "CONFIG":
        return "darkcyan";
      case "DEPLOY":
        return "darkgreen";
      case "MONITOR":
        return "slateblue";

      // General & Miscellaneous
      case "NOTE":
        return "gray";
      case "QUESTION":
        return "purple";
      case "IDEA":
        return "gold";
      case "REVIEW":
        return "blueviolet";

      default:
        return "black"; // fallback
    }
  };

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case "high":
        return "red";
      case "medium":
        return "yellow";
      case "low":
        return "blue";
      default:
        return "gray";
    }
  };

  const statusOptions = [
    { value: "todo", label: "To Do" },
    { value: "in-progress", label: "In Progress" },
    { value: "done", label: "Done" },
  ];

  async function updateItemStatus(itemId: number, status: string) {
    mutate({ id: itemId, status });
  }

  const { mutate: mutateOpenFile, isPending: isLoadingGoToFile } =
    useOpenFile();

  const goToFile = useCallback(
    (item: Item) => {
      mutateOpenFile({ filename: item.file, line: item.line });
    },
    [mutateOpenFile]
  );

  return (
    <Card
      p="xs"
      mb="sm"
      withBorder
      shadow="none"
      ref={cardRef}
      style={{
        ...style,
        cursor: isDragging ? "grabbing" : "grab",
        opacity: isDragging ? 0.6 : 1,
        transform: isDragging
          ? `${CSS.Transform.toString(transform)} rotate(5deg) scale(1.02)`
          : CSS.Transform.toString(transform),
        transition: "all 0.4 ease",
      }}
    >
      <Stack gap="xs">
        <Group>
          <Group justify="space-between" align="flex-start" w="100%">
            <Group align="flex-start" gap="xs">
              <RoleGuard.Consumer>
                <IconGripVertical
                  size={18}
                  {...attributes}
                  {...listeners}
                  style={{
                    outline: "unset",
                    marginLeft: -4,
                    cursor: isDragging ? "grabbing" : "grab",
                  }}
                />
              </RoleGuard.Consumer>
              <Text
                size="sm"
                onClick={() => onItemClick(item)}
                truncate
                maw={width - 100}
              >
                {item.title}
              </Text>
            </Group>

            <Group gap="4" mt="2">
              <ActionIcon size="xs" variant="subtle" >
                <IconEye onClick={() => onItemClick(item)} />
              </ActionIcon>
              <ActionIcon
                size="xs"
                variant="subtle"
                loading={isLoadingGoToFile}
              >
                <IconCode onClick={() => goToFile(item)} />
              </ActionIcon>
            </Group>
          </Group>
        </Group>
        {item.description && (
          <Box
            bdrs="sm"
            p="xs"
            onClick={() => onItemClick(item)}
          >
            <Text
              size="xs"
              c="dimmed"
              lineClamp={2}
              truncate
              styles={{ root: { whiteSpace: "break-spaces" } }}
            >
              {item.description}
            </Text>
          </Box>
        )}

        <RoleGuard.Consumer>
          <Select
            onChange={(value) => updateItemStatus(item.id, value!)}
            value={item.status}
            data={statusOptions}
            size="xs"
          />
        </RoleGuard.Consumer>

        <Group gap="4" align="center" py="4">
          <Badge color={getTypeColor(item.type)} size="xs" variant="dot">
            {item.type}
          </Badge>
          <Divider orientation="vertical" />
          <Badge color={getPriorityColor(item.priority)} size="xs">
            {item.priority}
          </Badge>
        </Group>
      </Stack>
    </Card>
  );
}
