import { useState, useCallback, useMemo } from "react";
import { Box, Group, ActionIcon, Collapse } from "@mantine/core";
import { useDroppable } from "@dnd-kit/core";
import { IconChevronDown } from "@tabler/icons-react";

interface Props {
  columnId: string;
  header: React.ReactNode;
  children: React.ReactNode;
}

// BUG: Hi
// Create a DroppableColumn component
// IN PROGRESS from 2025-08-20 12:06 by prodemmi
export default function DroppableColumn({ columnId, header, children }: Props) {
  const [collapsed, setCollapsed] = useState(false);

  const { setNodeRef, isOver } = useDroppable({
    id: columnId,
    disabled: collapsed,
  });

  // Memoize toggle function to prevent unnecessary re-renders
  const toggleCollapsed = useCallback(() => {
    setCollapsed((prev) => !prev);
  }, []);

  // Memoize styles for performance
  const boxStyles = useMemo(
    () => ({
      minHeight: collapsed ? undefined : isOver ? 350 : undefined,
      transition: "min-height 250ms cubic-bezier(0.4, 0, 0.2, 1)",
      overflow: "hidden",
    }),
    [collapsed, isOver]
  );

  const groupStyles = useMemo(
    () => ({
      marginBottom: collapsed ? undefined : "md",
      transition: "margin-bottom 250ms cubic-bezier(0.4, 0, 0.2, 1)",
      width: "100%",
    }),
    [collapsed]
  );

  // Memoize action icon styles
  const actionIconStyles = useMemo(
    () => ({
      transition: "transform 250ms cubic-bezier(0.4, 0, 0.2, 1)",
      transform: collapsed ? "rotate(0deg)" : "rotate(180deg)",
    }),
    [collapsed]
  );

  return (
    <Box ref={setNodeRef} style={boxStyles}>
      <Group
        style={groupStyles}
        justify="space-between"
        mb={collapsed ? undefined : "lg"}
      >
        {header}
        <ActionIcon
          size="xs"
          variant="transparent"
          onClick={toggleCollapsed}
          style={actionIconStyles}
          aria-label={collapsed ? "Expand column" : "Collapse column"}
        >
          <IconChevronDown />
        </ActionIcon>
      </Group>

      <Collapse
        in={!collapsed}
        transitionDuration={250}
        transitionTimingFunction="cubic-bezier(0.4, 0, 0.2, 1)"
      >
        {children}
      </Collapse>
    </Box>
  );
}
