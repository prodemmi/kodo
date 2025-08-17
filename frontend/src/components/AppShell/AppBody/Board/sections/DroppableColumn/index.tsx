import { useDroppable } from "@dnd-kit/core";
import { Box } from "@mantine/core";
import { ReactNode } from "react";

type Props = {
  columnId: string;
  children: ReactNode;
};

// BUG: Hi
// Create a DroppableColumn component
export default function DroppableColumn({ columnId, children }: Props) {
  const { setNodeRef, isOver } = useDroppable({
    id: columnId,
  });

  return (
    <Box
      ref={setNodeRef}
      style={{
        minHeight: isOver ? 350 : 300,
        // Removed background color changes and transitions
      }}
    >
      {children}
    </Box>
  );
}
