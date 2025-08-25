import {
  Text,
  ScrollArea,
  Box,
  Card,
  Group,
  Menu,
  ActionIcon,
  Badge,
  Avatar,
  Tooltip,
  Stack,
  Button,
  TextInput,
  TagsInput,
  Select,
  Popover,
  PopoverTarget,
  PopoverDropdown,
  Title,
  useMantineTheme,
} from "@mantine/core";
import {
  IconDotsVertical,
  IconTrash,
  IconGitBranch,
  IconFileText,
  IconPlus,
  IconFilter,
  IconSearch,
  IconTags,
  IconPin,
  IconPinFilled,
} from "@tabler/icons-react";
import { RoleGuard } from "../../../../Investor";
import { categories, tagColors } from "../constants";
import {
  useNewNoteModalStore,
  useDeleteModalStore,
  useNoteStore,
} from "../../../../../states/note.state";
import { useMediaQuery } from "@mantine/hooks";
import { useMemo, useRef } from "react";
import debounce from "lodash.debounce";
import { selectHasSearch } from "../../../../../states/note.selector";
import { useUpdateNote } from "../../../../../hooks/use-notes";
import { Note } from "../../../../../types/note";

export default function NoteList() {
  const notes = useNoteStore((s) => s.notes);
  const selectedNote = useNoteStore((s) => s.selectedNote);
  const setIsEditingNote = useNoteStore((s) => s.setIsEditingNote);
  const isPinned = useNoteStore((s) => s.isPinned);
  const selectNote = useNoteStore((s) => s.selectNote);
  const selectedFolder = useNoteStore((s) => s.selectedFolder);
  const openForNote = useDeleteModalStore((s) => s.openForNote);
  const openNewNoteModal = useNewNoteModalStore((s) => s.openModal);
  const isSmall = useMediaQuery("(max-width: 920px)");

  const allTags = useNoteStore.getState().tags();
  const searchQuery = useNoteStore((s) => s.searchQuery);
  const setSearchQuery = useNoteStore((s) => s.setSearchQuery);
  const searchTags = useNoteStore((s) => s.searchTags);
  const setSearchTags = useNoteStore((s) => s.setSearchTags);
  const filterCategory = useNoteStore((s) => s.filterCategory);
  const setFilterCategory = useNoteStore((s) => s.setFilterCategory);
  const clearSearch = useNoteStore((s) => s.clearSearch);
  const hasSearch = useNoteStore(selectHasSearch);
  const searchInputRef = useRef<HTMLInputElement | null>(null);
  const { primaryColor, colors } = useMantineTheme();

  const { mutate: updateNote } = useUpdateNote();

  const togglePinNote = (note: Note) => {
    updateNote({
      ...note,
      id: note.id,
      pinned: !note.pinned,
    });
  };

  const notesInDirectory = useMemo(() => {
    if (selectedFolder) {
      return notes.filter((n) => n.folderId === selectedFolder.id);
    }

    return notes;
  }, [notes, selectedFolder]);

  const filteredNotes = useMemo(() => {
    let shallowNotes = [...notes];
    const q = searchQuery.trim().toLowerCase();

    if (selectedFolder) {
      shallowNotes = shallowNotes.filter(
        (n) => n.folderId === selectedFolder.id
      );
    }

    return shallowNotes.filter((note) => {
      if (filterCategory && note.category !== filterCategory) return false;
      if (searchTags.length && !searchTags.every((t) => note.tags.includes(t)))
        return false;
      if (q) {
        const textMatch =
          note.title.toLowerCase().includes(q) ||
          note.content.toLowerCase().includes(q) ||
          note.gitBranch?.toLowerCase().includes(q) ||
          note.gitCommit?.toLowerCase().includes(q) ||
          note.author.toLowerCase().includes(q);

        if (!textMatch) return false;
      }

      return true;
    });
  }, [notes, searchQuery, searchTags, filterCategory, selectedFolder]);

  const getCategoryColor = (category: any) => {
    const colors: any = {
      technical: "blue",
      meeting: "purple",
      idea: "green",
      documentation: "orange",
      "bug-analysis": "red",
      review: "cyan",
    };
    return colors[category] || "gray";
  };

  const formatDate = (date: any) => {
    return new Intl.DateTimeFormat("en-US", {
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    }).format(new Date(date));
  };

  const debouncedSearch = useMemo(
    () =>
      debounce((value: string) => {
        setSearchQuery(value);
      }, 300),
    [setSearchQuery]
  );

  const clearSearches = () => {
    clearSearch();
    if (searchInputRef.current) {
      searchInputRef.current.value = "";
    }
  };

  return (
    <Stack
      w="60%"
      p="sm"
      pt="md"
      pr="xs"
      gap="xs"
      style={{ height: "calc(100dvh - 44px)" }}
    >
      <Group justify="space-between" w="100%">
        <Title size="h6">{selectedFolder?.name || "All Notes"}</Title>
        {notesInDirectory && notesInDirectory.length > 0 && (
          <Popover position="right-start">
            <PopoverTarget>
              <ActionIcon size="sm" c={hasSearch ? primaryColor : undefined}>
                <IconFilter size={14} color="white" />
              </ActionIcon>
            </PopoverTarget>
            <PopoverDropdown>
              <Stack gap="xs" align="flex-end" maw="230px">
                <TextInput
                  ref={searchInputRef}
                  w="100%"
                  size="xs"
                  placeholder="Search notes..."
                  defaultValue={searchQuery}
                  onChange={(e) => debouncedSearch(e.currentTarget.value)}
                  leftSection={<IconSearch size={16} />}
                />

                <TagsInput
                  w="100%"
                  size="xs"
                  placeholder="Search by tags..."
                  allowDuplicates={false}
                  value={searchTags}
                  onChange={(val) => {
                    const filtered = val.filter((v) => allTags.includes(v));
                    setSearchTags(filtered);
                  }}
                  comboboxProps={{ withinPortal: false }}
                  data={allTags}
                  leftSection={<IconTags size={16} />}
                  acceptValueOnBlur={false}
                />

                <Select
                  w="100%"
                  size="xs"
                  placeholder="Filter by category"
                  value={filterCategory}
                  onChange={(value) => setFilterCategory(value || "")}
                  comboboxProps={{ withinPortal: false }}
                  data={[{ value: "", label: "All Categories" }, ...categories]}
                  leftSection={<IconFilter size={16} />}
                  clearable
                />
                {hasSearch && (
                  <Button onClick={clearSearches}>Clear Search</Button>
                )}
              </Stack>
            </PopoverDropdown>
          </Popover>
        )}
      </Group>

      <ScrollArea style={{ flex: 1 }} h={isSmall ? undefined : "100%"}>
        <Box>
          {filteredNotes.map((note: any) => (
            <Card
              key={note.id}
              padding="xs"
              mb="xs"
              style={{
                cursor: "pointer",
                border:
                  selectedNote?.id === note.id
                    ? `1px solid ${colors[primaryColor][5]}`
                    : `1px solid ${colors.dark[5]}`,
              }}
              onClick={() => {
                selectNote(note);
                setIsEditingNote(false);
              }}
            >
              <Group justify="space-between" mb="xs">
                <Text fw={500} size="sm" style={{ flex: 1 }} truncate>
                  {note.title}
                </Text>
                <Group gap="xs">
                  <RoleGuard.Consumer>
                    <ActionIcon
                      variant="subtle"
                      size="sm"
                      onClick={(e) => {
                        e.stopPropagation();
                        togglePinNote(note);
                      }}
                    >
                      {isPinned(note.id) ? (
                        <IconPinFilled size={14} />
                      ) : (
                        <IconPin size={14} />
                      )}
                    </ActionIcon>
                    <Menu>
                      <Menu.Target>
                        <ActionIcon
                          variant="subtle"
                          size="sm"
                          onClick={(e) => e.stopPropagation()}
                        >
                          <IconDotsVertical size={14} />
                        </ActionIcon>
                      </Menu.Target>
                      <Menu.Dropdown>
                        <Menu.Item
                          leftSection={<IconTrash size={14} color="red" />}
                          onClick={(e) => {
                            e.stopPropagation();
                            openForNote(note);
                          }}
                        >
                          Delete
                        </Menu.Item>
                      </Menu.Dropdown>
                    </Menu>
                  </RoleGuard.Consumer>
                </Group>
              </Group>

              <Text
                size="xs"
                c="dimmed"
                mb="xs"
                style={{
                  display: "-webkit-box",
                  WebkitLineClamp: 2,
                  WebkitBoxOrient: "vertical",
                  overflow: "hidden",
                  whiteSpace: "pre",
                }}
              >
                {note.content.replace(/<[^>]*>/g, "")}
              </Text>

              <Group gap="xs" mb="xs">
                <Badge
                  size="xs"
                  color={getCategoryColor(note.category)}
                  variant="light"
                >
                  {note.category}
                </Badge>
                {note.tags.slice(0, 2).map((tag: any) => (
                  <Badge
                    key={tag}
                    size="xs"
                    color={tagColors[tag] || "gray"}
                    variant="outline"
                  >
                    {tag}
                  </Badge>
                ))}
                {note.tags.length > 2 && (
                  <Badge size="xs" variant="outline" color="gray">
                    +{note.tags.length - 2}
                  </Badge>
                )}
              </Group>

              <Group justify="space-between" align="center">
                <Group gap="xs" align="center">
                  <Avatar size={18}>
                    {note.author.charAt(0).toUpperCase()}
                  </Avatar>
                  <Text size="xs" c="dimmed">
                    {note.author}
                  </Text>
                </Group>
                <Group gap="xs">
                  <Tooltip label={`Branch: ${note.gitBranch}`}>
                    <IconGitBranch size={12} color="#868e96" />
                  </Tooltip>
                  <Text size="xs" c="dimmed">
                    {formatDate(note.updatedAt)}
                  </Text>
                </Group>
              </Group>
            </Card>
          ))}

          {filteredNotes.length === 0 ? (
            <Card padding="lg">
              <Stack align="center" gap="sm">
                <IconFileText size={48} color="#ced4da" />
                <Text c="dimmed" ta="center">
                  {searchQuery || filterCategory || searchTags.length > 0
                    ? "No notes match your filters"
                    : selectedFolder
                    ? "No notes in this folder"
                    : "No notes yet"}
                </Text>
                <Button
                  variant="light"
                  leftSection={<IconPlus size={16} />}
                  onClick={openNewNoteModal}
                >
                  Create a note
                </Button>
              </Stack>
            </Card>
          ) : (
            <RoleGuard.Consumer>
              <Button
                size="sm"
                w="100%"
                leftSection={<IconPlus size={20} />}
                onClick={openNewNoteModal}
              >
                Create New Note
              </Button>
            </RoleGuard.Consumer>
          )}
        </Box>
      </ScrollArea>
    </Stack>
  );
}
