import { useState, useEffect, useCallback, lazy, Suspense } from "react";
import {
  DndContext,
  closestCenter,
  useSensor,
  useSensors,
  PointerSensor,
  DragEndEvent,
  DragOverlay,
  DragStartEvent,
} from "@dnd-kit/core";
import {
  SortableContext,
  verticalListSortingStrategy,
  arrayMove,
} from "@dnd-kit/sortable";
import {
  Group,
  Text,
  Badge,
  Title,
  Container,
  Paper,
  Alert,
  Button,
} from "@mantine/core";
import { IconAlertCircle, IconHistory } from "@tabler/icons-react";
import { useQueryClient } from "@tanstack/react-query";
import { Item } from "../../../../types/item";
import { useItems, useUpdateItem } from "../../../../hooks/use-items";
import ItemDetailDrawer from "./sections/ItemDetailDrawer";
import SortableTask from "./sections/SortableTask";
import { useSettings } from "../../../../hooks/use-settings";

const DroppableColumn = lazy(() => import("./sections/DroppableColumn"));
const HistoryDrawer = lazy(() => import("./sections/HistoryDrawer"));

interface Column {
  title: string;
  tasks: Item[];
  color?: string;
}

export default function Board() {
  const queryClient = useQueryClient();
  const {
    data: items,
    isSuccess: isSuccessItems,
    error: itemsError,
  } = useItems();
  const {
    data: settings,
    isSuccess: isSuccessSettings,
    error: settingsError,
  } = useSettings();

  const [columns, setColumns] = useState<Record<string, Column>>({});
  const [columnOrder, setColumnOrder] = useState<string[]>([]);
  const [openItemHistory, setOpenItemHistory] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [drawerOpened, setDrawerOpened] = useState(false);
  const [selectedItem, setSelectedItem] = useState<Item | null>(null);
  const [activeItem, setActiveItem] = useState<Item | null>(null);

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 3,
      },
    })
  );

  const { mutate } = useUpdateItem();

  // Fetch items and organize into columns
  useEffect(() => {
    const fetchData = async () => {
      try {
        if (!settings || !items) {
          return;
        }

        const newColumns: Record<string, Column> = {};
        const newColumnOrder: string[] = [];

        // Create columns from settings
        settings.kanban_columns.forEach((column) => {
          newColumns[column.id] = {
            title: column.name,
            tasks: [],
            color: column.color,
          };
          newColumnOrder.push(column.id);
        });

        // Distribute items into columns
        items.forEach((item) => {
          if (newColumns[item.status]) {
            newColumns[item.status].tasks.push(item);
          }
        });

        setColumns(newColumns);
        setColumnOrder(newColumnOrder);
      } catch (err) {
        setError("Failed to fetch items");
      }
    };

    if (isSuccessItems && isSuccessSettings) {
      fetchData();
    }
  }, [items, settings, isSuccessItems, isSuccessSettings]);

  const handleDragStart = (event: DragStartEvent) => {
    const { active } = event;
    const activeId = active.id.toString();

    // Find the active item
    const activeColumnId = Object.keys(columns).find((colId) =>
      columns[colId].tasks.some((item) => item.id.toString() === activeId)
    );

    if (activeColumnId) {
      const item = columns[activeColumnId].tasks.find(
        (item) => item.id.toString() === activeId
      );
      setActiveItem(item || null);
    }
  };

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    setActiveItem(null);

    if (!over) return;

    const activeId = active.id.toString();
    const overId = over.id.toString();

    // Find which column contains the active item
    const activeColumnId = Object.keys(columns).find((colId) =>
      columns[colId].tasks.some((item) => item.id.toString() === activeId)
    );

    if (!activeColumnId) return;

    // Determine the target column
    let overColumnId: string | undefined;

    // Check if overId is a column ID directly
    if (columnOrder.includes(overId)) {
      overColumnId = overId;
    } else {
      // Check if overId is an item ID, find its column
      overColumnId = Object.keys(columns).find((colId) =>
        columns[colId].tasks.some((item) => item.id.toString() === overId)
      );
    }

    if (!overColumnId) return;

    if (activeColumnId === overColumnId) {
      // Reordering within the same column
      const tasks = columns[activeColumnId].tasks;
      const activeIndex = tasks.findIndex(
        (item) => item.id.toString() === activeId
      );
      const overIndex = tasks.findIndex(
        (item) => item.id.toString() === overId
      );

      if (activeIndex !== overIndex) {
        const newTasks = arrayMove(tasks, activeIndex, overIndex);
        setColumns({
          ...columns,
          [activeColumnId]: { ...columns[activeColumnId], tasks: newTasks },
        });
      }
      return;
    }

    // Moving between columns
    const sourceTasks = [...columns[activeColumnId].tasks];
    const destTasks = [...columns[overColumnId].tasks];

    const activeIndex = sourceTasks.findIndex(
      (item) => item.id.toString() === activeId
    );

    if (activeIndex === -1) return;

    const [movedTask] = sourceTasks.splice(activeIndex, 1);
    const updatedTask = {
      ...movedTask,
      status: overColumnId as Item["status"],
    };

    // Add to destination at the end
    destTasks.push(updatedTask);

    // Update state optimistically
    setColumns({
      ...columns,
      [activeColumnId]: { ...columns[activeColumnId], tasks: sourceTasks },
      [overColumnId]: { ...columns[overColumnId], tasks: destTasks },
    });

    // Call the API update with optimistic update handling
    updateItemStatus(Number(updatedTask.id), updatedTask.status);
  };

  function updateItemStatus(itemId: number, status: string) {
    console.log("Calling API - itemId:", itemId, "status:", status);

    mutate(
      { id: itemId, status },
      {
        onSuccess: async (newItem) => {
          // Cancel any outgoing refetches
          await queryClient.cancelQueries({ queryKey: ["items"] });

          // Snapshot the previous value
          const previousItems = queryClient.getQueryData(["items"]);

          // Optimistically update the cache
          queryClient.setQueryData(["items"], (old: Item[] | undefined) => {
            if (!old) return old;
            return old.map((item) =>
              item.id === newItem.id
                ? { ...item, status: newItem.status }
                : item
            );
          });

          return { previousItems };
        },
        onError: (_, __, context: any) => {
          // Rollback on error
          queryClient.setQueryData(["items"], context?.previousItems);

          // Also rollback local state
          if (items && settings) {
            const rollbackColumns: Record<string, Column> = {};
            settings.kanban_columns.forEach((column) => {
              rollbackColumns[column.id] = {
                title: column.name,
                tasks: [],
                color: column.color,
              };
            });

            items.forEach((item) => {
              if (rollbackColumns[item.status]) {
                rollbackColumns[item.status].tasks.push(item);
              }
            });

            setColumns(rollbackColumns);
          }
        },
        onSettled: () => {
          // Refetch to ensure consistency
          queryClient.invalidateQueries({ queryKey: ["items"] });
        },
      }
    );
  }

  const handleItemClick = (item: Item) => {
    setSelectedItem(item);
    setDrawerOpened(true);
  };

  const columnHeader = useCallback((column: Column) => {
    return (
      <Group justify="space-between">
        <Text fw={600} size="sm">
          {column.title}
        </Text>
        <Badge variant="light" size="sm" color="gray">
          {column.tasks.length}
        </Badge>
      </Group>
    );
  }, []);

  // Error states
  const combinedError = error || itemsError || settingsError;
  if (combinedError) {
    return (
      <Container size="xl" py="md">
        <Alert icon={<IconAlertCircle size={16} />} title="Error" color="red">
          {typeof combinedError === "string"
            ? combinedError
            : combinedError.message || "An error occurred"}
        </Alert>
      </Container>
    );
  }

  return (
    items &&
    settings && (
      <>
        <Container size="xl" py="lg">
          <Group justify="space-between" align="center" w="100%">
            <Title order={2} mb="md">
              Kanban Board
            </Title>
            <Button
              leftSection={<IconHistory size={16} />}
              onClick={() => setOpenItemHistory(true)}
            >
              History
            </Button>
          </Group>
          <DndContext
            sensors={sensors}
            collisionDetection={closestCenter}
            onDragStart={handleDragStart}
            onDragEnd={handleDragEnd}
          >
            <Group align="start" gap="lg">
              {columnOrder.map((columnId) => {
                const column = columns[columnId];
                if (!column) return null;

                return (
                  <Paper
                    shadow="sm"
                    p="md"
                    style={{
                      width: 300,
                      borderRadius: 8,
                    }}
                    withBorder
                    key={columnId}
                  >
                    <Suspense>
                      <DroppableColumn
                        color={column.color!}
                        columnId={columnId}
                        header={columnHeader(column)}
                      >
                        <SortableContext
                          items={column.tasks.map((item) => item.id)}
                          strategy={verticalListSortingStrategy}
                        >
                          <div style={{ minHeight: 300 }}>
                            {column.tasks.map((item) => (
                              <SortableTask
                                key={item.id}
                                item={item}
                                onItemClick={handleItemClick}
                              />
                            ))}
                            {column.tasks.length === 0 && (
                              <div
                                style={{
                                  height: 200,
                                  display: "flex",
                                  alignItems: "center",
                                  justifyContent: "center",
                                  color: "#adb5bd",
                                  fontSize: "14px",
                                  borderRadius: 8,
                                }}
                              >
                                Drop items here
                              </div>
                            )}
                          </div>
                        </SortableContext>
                      </DroppableColumn>
                    </Suspense>
                  </Paper>
                );
              })}
            </Group>

            <DragOverlay>
              {activeItem ? (
                <div
                  style={{
                    transform: "rotate(5deg)",
                    opacity: 0.9,
                    transition: "all 0.4s ease",
                  }}
                >
                  <SortableTask item={activeItem} onItemClick={() => {}} />
                </div>
              ) : null}
            </DragOverlay>
          </DndContext>
        </Container>
        <ItemDetailDrawer
          drawerOpened={drawerOpened}
          selectedItem={selectedItem}
          setDrawerOpened={setDrawerOpened}
        />
        <Suspense>
          <HistoryDrawer
            isOpen={openItemHistory}
            onClose={() => {
              setOpenItemHistory(false);
            }}
          />
        </Suspense>
      </>
    )
  );
}
