import {
  Container,
  Box,
  Paper,
  Text,
  ActionIcon,
  Group,
  Menu,
  LoadingOverlay,
  ScrollArea,
  Textarea,
} from "@mantine/core";
import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  DragEndEvent,
  DragOverlay,
} from "@dnd-kit/core";
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  rectSortingStrategy,
  useSortable,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import {
  IconGripVertical,
  IconTrash,
  IconEdit,
  IconPalette,
  IconPin,
  IconPinFilled,
} from "@tabler/icons-react";
import { useState, useEffect } from "react";
import { useNotes, useFolders } from "../../../../../hooks/use-notes";
import { useNoteStore } from "../../../../../states/note.state";

interface StickyNote {
  id: string;
  title: string;
  content: string;
  color: string;
  position: { x: number; y: number };
  isPinned: boolean;
  folderId?: string;
  createdAt: string;
  updatedAt: string;
}

function StickyNoteCard({
  note,
  onUpdate,
  onDelete,
}: {
  note: StickyNote;
  onUpdate: (id: string, updates: Partial<StickyNote>) => void;
  onDelete: (id: string) => void;
}) {
  const [isEditing, setIsEditing] = useState(false);
  const [editContent, setEditContent] = useState(note.content);
  const [editTitle, setEditTitle] = useState(note.title);

  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: note.id,
  });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition: `${transition}, box-shadow 0.2s ease`,
    opacity: isDragging ? 0.7 : 1,
  };

  const handleSave = () => {
    onUpdate(note.id, {
      title: editTitle,
      content: editContent,
      updatedAt: new Date().toISOString(),
    });
    setIsEditing(false);
  };

  const handleCancel = () => {
    setEditTitle(note.title);
    setEditContent(note.content);
    setIsEditing(false);
  };

  const noteColors = [
    "var(--mantine-color-blue-2)",
    "var(--mantine-color-green-2)",
    "var(--mantine-color-yellow-2)",
    "var(--mantine-color-red-2)",
    "var(--mantine-color-purple-2)",
    "var(--mantine-color-orange-2)",
    "var(--mantine-color-teal-2)",
    "var(--mantine-color-pink-2)",
  ];

  return (
    <Paper
      ref={setNodeRef}
      style={style}
      shadow="sm"
      p="md"
      radius="md"
      bg={note.color}
      w={{ base: 260, sm: 280, md: 300 }}
      h={{ base: 300, sm: 320, md: 340 }}
      pos="relative"
    >
      {/* Drag Handle */}
      <Box
        {...attributes}
        {...listeners}
        pos="absolute"
        style={{
          cursor: "grab",
          touchAction: "none",
        }}
      >
        <IconGripVertical size={16} />
      </Box>

      {/* Actions Menu */}
      <Group
        pos="absolute"
        className="sticky-actions"
        style={{
          opacity: 0,
          transition: "opacity 0.2s ease",
        }}
      >
        <ActionIcon
          size="sm"
          variant="subtle"
          color="gray"
          onClick={() => onUpdate(note.id, { isPinned: !note.isPinned })}
        >
          {note.isPinned ? <IconPinFilled size={14} /> : <IconPin size={14} />}
        </ActionIcon>

        <Menu shadow="md" width={200} withinPortal>
          <Menu.Target>
            <ActionIcon size="sm" variant="subtle" color="gray">
              <IconPalette size={14} />
            </ActionIcon>
          </Menu.Target>
          <Menu.Dropdown>
            <Menu.Label>Change Color</Menu.Label>
            <Box p="xs">
              <Group gap="xs">
                {noteColors.map((color) => (
                  <Box
                    key={color}
                    w={24}
                    h={24}
                    bg={color}
                    style={{
                      borderRadius: "50%",
                      cursor: "pointer",
                      border:
                        color === note.color
                          ? `2px solid var(--mantine-color-dark-7)`
                          : `1px solid var(--mantine-color-dark-3)`,
                      transition: "transform 0.2s ease",
                      "&:hover": {
                        transform: "scale(1.1)",
                      },
                    }}
                    onClick={() => onUpdate(note.id, { color })}
                  />
                ))}
              </Group>
            </Box>
          </Menu.Dropdown>
        </Menu>

        <ActionIcon
          size="sm"
          variant="subtle"
          onClick={() => setIsEditing(!isEditing)}
        >
          <IconEdit size={14} />
        </ActionIcon>

        <ActionIcon
          size="sm"
          variant="subtle"
          onClick={() => onDelete(note.id)}
        >
          <IconTrash size={14} color="red" />
        </ActionIcon>
      </Group>

      {/* Content */}
      <Box mt={"xl"} h="calc(100% - 32px)">
        {isEditing ? (
          <Box h="100%">
            <Textarea
              placeholder="Note title..."
              value={editTitle}
              onChange={(e) => setEditTitle(e.target.value)}
              size="sm"
              mb="xs"
              autosize
              minRows={1}
              maxRows={2}
              styles={{
                input: {
                  fontWeight: 600,
                  fontSize: "md",
                  border: `1px solid var(--mantine-color-gray-3)`,
                },
              }}
            />
            <Textarea
              placeholder="Write your note..."
              value={editContent}
              onChange={(e) => setEditContent(e.target.value)}
              size="sm"
              h="calc(100% - 80px)"
              styles={{
                input: {
                  height: "100%",
                  resize: "none",
                  fontSize: "sm",
                  border: `1px solid var(--mantine-color-gray-3)}`,
                },
              }}
            />
            <Group justify="flex-end" gap={"xs"} mt="xs">
              <ActionIcon
                size="sm"
                variant="filled"
                color="green"
                onClick={handleSave}
              >
                ✓
              </ActionIcon>
              <ActionIcon
                size="sm"
                variant="filled"
                color="gray"
                onClick={handleCancel}
              >
                ✗
              </ActionIcon>
            </Group>
          </Box>
        ) : (
          <Box
            h="100%"
            onClick={() => setIsEditing(true)}
            style={{ cursor: "text" }}
          >
            <Text fw={600} mb="xs" lineClamp={2} size="md">
              {note.title || "Untitled Note"}
            </Text>
            <ScrollArea h="calc(100% - 40px)" type="auto">
              <Text
                size="sm"
                style={{ whiteSpace: "pre-wrap", lineHeight: 1.5 }}
              >
                {note.content || "Click to add content..."}
              </Text>
            </ScrollArea>
          </Box>
        )}
      </Box>
    </Paper>
  );
}

export default function NotesStickyView() {
  const [stickyNotes, setStickyNotes] = useState<StickyNote[]>([]);
  const [activeId, setActiveId] = useState<string | null>(null);

  const {
    data: notesData,
    isError: notesError,
    isLoading: notesLoading,
  } = useNotes();

  const {
    data: foldersData,
    isError: foldersError,
    isLoading: foldersLoading,
  } = useFolders();

  const notes = useNoteStore((s) => s.notes);
  const folders = useNoteStore((s) => s.folders);
  const setNotes = useNoteStore((s) => s.setNotes);
  const setFolders = useNoteStore((s) => s.setFolders);

  const noteColors = [
    "var(--mantine-color-blue-2)",
    "var(--mantine-color-green-2)",
    "var(--mantine-color-yellow-2)",
    "var(--mantine-color-red-2)",
    "var(--mantine-color-purple-2)",
    "var(--mantine-color-orange-2)",
    "var(--mantine-color-teal-2)",
    "var(--mantine-color-pink-2)",
  ];

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8,
      },
    }),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  );

  // Convert notes to sticky notes format
  useEffect(() => {
    if (notes && notes.length > 0) {
      const convertedNotes: StickyNote[] = notes.map((note, index) => ({
        id: note.id,
        title: note.title || "Untitled",
        content: note.content || "",
        color: noteColors[index % noteColors.length],
        position: {
          x: (index % 4) * 320 + 20,
          y: Math.floor(index / 4) * 360 + 20,
        },
        isPinned: false,
        folderId: note.folderId,
        createdAt: note.createdAt || new Date().toISOString(),
        updatedAt: note.updatedAt || new Date().toISOString(),
      }));
      setStickyNotes(convertedNotes);
    }
  }, [notes]);

  // Load initial data
  useEffect(() => {
    if (
      !foldersError &&
      !foldersLoading &&
      foldersData &&
      foldersData.count > 0
    ) {
      setFolders(foldersData.folders);
    }
    if (!notesError && !notesLoading && notesData && notesData.count > 0) {
      setNotes(notesData.notes);
    }
  }, [
    notesData,
    notesError,
    notesLoading,
    foldersData,
    foldersError,
    foldersLoading,
    setNotes,
    setFolders,
  ]);

  const handleDragStart = (event: any) => {
    setActiveId(event.active.id);
  };

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;

    if (active.id !== over?.id) {
      setStickyNotes((items) => {
        const oldIndex = items.findIndex((item) => item.id === active.id);
        const newIndex = items.findIndex((item) => item.id === over?.id);

        return arrayMove(items, oldIndex, newIndex);
      });
    }

    setActiveId(null);
  };

  const handleUpdateNote = (id: string, updates: Partial<StickyNote>) => {
    setStickyNotes((prev) =>
      prev.map((note) => (note.id === id ? { ...note, ...updates } : note))
    );
  };

  const handleDeleteNote = (id: string) => {
    setStickyNotes((prev) => prev.filter((note) => note.id !== id));
  };

  const activeNote = activeId
    ? stickyNotes.find((note) => note.id === activeId)
    : null;

  return (
    notes &&
    folders && (
      <Container
        fluid
        p="lg"
        style={{
          height: "calc(100dvh - 60px)",
          backgroundColor: "var(--mantine-color-gray-0)",
          position: "relative",
          overflow: "auto",
        }}
      >
        <DndContext
          sensors={sensors}
          collisionDetection={closestCenter}
          onDragStart={handleDragStart}
          onDragEnd={handleDragEnd}
        >
          <SortableContext items={stickyNotes} strategy={rectSortingStrategy}>
            <Box
              style={{
                display: "grid",
                gridTemplateColumns: "repeat(auto-fill, minmax(280px, 1fr))",
                gap: "lg",
                justifyContent: "start",
                minHeight: "100%",
                padding: "md",
              }}
            >
              {stickyNotes.map((note) => (
                <StickyNoteCard
                  key={note.id}
                  note={note}
                  onUpdate={handleUpdateNote}
                  onDelete={handleDeleteNote}
                />
              ))}
            </Box>
          </SortableContext>

          <DragOverlay>
            {activeNote && (
              <StickyNoteCard
                note={activeNote}
                onUpdate={() => {}}
                onDelete={() => {}}
              />
            )}
          </DragOverlay>
        </DndContext>

        {stickyNotes.length === 0 && (
          <Box
            style={{
              position: "absolute",
              top: "50%",
              left: "50%",
              transform: "translate(-50%, -50%)",
              textAlign: "center",
            }}
          >
            <Text size="xl" fw={500} c="dimmed" mb="md">
              No notes to display
            </Text>
            <Text size="sm" c="dimmed">
              Create your first note to see it as a sticky note!
            </Text>
          </Box>
        )}
      </Container>
    )
  );
}
