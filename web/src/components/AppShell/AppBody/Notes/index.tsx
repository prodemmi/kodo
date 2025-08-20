import React, { useState, useEffect, useRef } from "react";
import {
  Container,
  Title,
  Button,
  Group,
  Stack,
  Card,
  Text,
  Badge,
  ActionIcon,
  Modal,
  TextInput,
  Select,
  Divider,
  Avatar,
  Tooltip,
  Alert,
  Loader,
  Menu,
  ScrollArea,
  Box,
  Textarea,
  MultiSelect,
  Drawer,
  UnstyledButton,
  Collapse,
  Flex,
  Paper,
  TagsInput,
} from "@mantine/core";
import { RichTextEditor, Link } from "@mantine/tiptap";
import { useEditor } from "@tiptap/react";
import Highlight from "@tiptap/extension-highlight";
import StarterKit from "@tiptap/starter-kit";
import Underline from "@tiptap/extension-underline";
import TextAlign from "@tiptap/extension-text-align";
import Superscript from "@tiptap/extension-superscript";
import SubScript from "@tiptap/extension-subscript";
import CodeBlockLowlight from "@tiptap/extension-code-block-lowlight";
import { createLowlight } from "lowlight";
import {
  IconPlus,
  IconEdit,
  IconTrash,
  IconSearch,
  IconFilter,
  IconFileText,
  IconGitBranch,
  IconDotsVertical,
  IconFolder,
  IconFolderPlus,
  IconChevronDown,
  IconChevronRight,
  IconTags,
  IconX,
  IconCheck,
  IconAlertTriangle,
} from "@tabler/icons-react";
import debounce from "lodash.debounce";
import "@mantine/tiptap/styles.css";
import { RoleGuard } from "../../../Investor";

// Mock data with folders
const mockFolders = [
  { id: 1, name: "Architecture", parentId: null, expanded: true },
  { id: 2, name: "Meetings", parentId: null, expanded: true },
  { id: 3, name: "Ideas", parentId: null, expanded: false },
  { id: 4, name: "Sprint Planning", parentId: 2, expanded: true },
];

const mockNotes = [
  {
    id: 1,
    title: "Architecture Review Notes",
    content:
      '<h2>Database Schema Changes</h2><p>Need to refactor the user authentication system to support OAuth2. Current implementation is becoming difficult to maintain.</p><ul><li>Update user table structure</li><li>Add OAuth provider mapping</li><li>Migrate existing passwords</li></ul><pre><code>type User struct {\n    ID       int    `json:"id"`\n    Username string `json:"username"`\n    Email    string `json:"email"`\n}</code></pre>',
    author: "john.doe",
    createdAt: new Date("2024-08-15T10:30:00"),
    updatedAt: new Date("2024-08-16T14:20:00"),
    tags: ["architecture", "auth", "database"],
    category: "technical",
    folderId: 1,
    gitBranch: "feature/oauth2",
    gitCommit: "a1b2c3d",
  },
  {
    id: 2,
    title: "Meeting Notes - Sprint Planning",
    content:
      "<h2>Sprint 23 Planning</h2><p>Discussed the upcoming features and technical debt items.</p><h3>Key Decisions:</h3><p>Focus on performance improvements this sprint. The TODO items analysis shows we have 23 performance-related tasks.</p><blockquote><p>Important: All team members should prioritize code reviews this week.</p></blockquote>",
    author: "jane.smith",
    createdAt: new Date("2024-08-14T09:00:00"),
    updatedAt: new Date("2024-08-14T09:00:00"),
    tags: ["meeting", "planning", "sprint"],
    category: "meeting",
    folderId: 4,
    gitBranch: "main",
    gitCommit: "e4f5g6h",
  },
  {
    id: 3,
    title: "Performance Optimization Ideas",
    content:
      "<h1>Performance Ideas</h1><p>Collection of ideas for improving application performance:</p><ul><li><strong>Database indexing</strong> - Add indexes on frequently queried columns</li><li><strong>Caching layer</strong> - Implement Redis for session storage</li><li><strong>Code splitting</strong> - Lazy load components in frontend</li></ul>",
    author: "bob.wilson",
    createdAt: new Date("2024-08-13T16:45:00"),
    updatedAt: new Date("2024-08-13T16:45:00"),
    tags: ["performance", "optimization", "database", "caching"],
    category: "idea",
    folderId: 3,
    gitBranch: "main",
    gitCommit: "x9y8z7w",
  },
];

const categories = [
  { value: "technical", label: "Technical" },
  { value: "meeting", label: "Meeting" },
  { value: "idea", label: "Idea" },
  { value: "documentation", label: "Documentation" },
  { value: "bug-analysis", label: "Bug Analysis" },
  { value: "review", label: "Review" },
];

const tagColors: any = {
  architecture: "blue",
  auth: "red",
  database: "green",
  meeting: "purple",
  planning: "orange",
  sprint: "pink",
  performance: "yellow",
  security: "red",
  feature: "blue",
  bug: "red",
  refactor: "cyan",
  optimization: "teal",
  caching: "lime",
};

export default function Notes() {
  const [notes, setNotes] = useState<any>(mockNotes);
  const [folders, setFolders] = useState(mockFolders);
  const [selectedNote, setSelectedNote] = useState<any>(null);
  const [isEditing, setIsEditing] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [searchTags, setSearchTags] = useState<any[]>([]);
  const [filterCategory, setFilterCategory] = useState("");
  const [selectedFolder, setSelectedFolder] = useState<any>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isFolderModalOpen, setIsFolderModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [noteToDelete, setNoteToDelete] = useState<any>(null);
  const [loading, setLoading] = useState(false);
  const [newNoteTitle, setNewNoteTitle] = useState("");
  const [newNoteCategory, setNewNoteCategory] = useState<any>("technical");
  const [newNoteTags, setNewNoteTags] = useState<any[]>([]);
  const [newFolderName, setNewFolderName] = useState("");
  const [newFolderParent, setNewFolderParent] = useState("");
  const [editingTags, setEditingTags] = useState(false);
  const [tempTags, setTempTags] = useState<any[]>([]);
  const [error, setError] = useState("");

  const lowlight = createLowlight();

  const editor = useEditor({
    extensions: [
      StarterKit,
      Underline,
      Link,
      Superscript,
      SubScript,
      Highlight,
      TextAlign.configure({ types: ["heading", "paragraph"] }),
      CodeBlockLowlight.configure({
        lowlight,
      }),
    ],
    content: selectedNote?.content || "",
    editable: isEditing,
    onUpdate: ({ editor }) => {
      if (selectedNote && isEditing) {
        setSelectedNote({
          ...selectedNote,
          content: editor.getHTML(),
          updatedAt: new Date(),
        });
      }
    },
  });

  useEffect(() => {
    if (editor && selectedNote) {
      editor.commands.setContent(selectedNote.content);
      editor.setEditable(isEditing);
    }
  }, [selectedNote, isEditing, editor]);

  // Debounced search
  const debouncedSearch = debounce((value) => {
    setSearchQuery(value);
  }, 300);

  // Get all available tags from existing notes
  const allTags: any = [...new Set(notes.flatMap((note: any) => note.tags))];

  // Filter notes based on search, tags, category, and folder
  const filteredNotes = notes.filter((note: any) => {
    const matchesSearch =
      note.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
      note.content.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesTags =
      searchTags.length === 0 ||
      searchTags.every((tag) => note.tags.includes(tag));
    const matchesCategory = !filterCategory || note.category === filterCategory;
    const matchesFolder = !selectedFolder || note.folderId === selectedFolder;
    return matchesSearch && matchesTags && matchesCategory && matchesFolder;
  });

  const handleCreateNote = async () => {
    if (!newNoteTitle.trim()) {
      setError("Note title is required");
      return;
    }

    const newNote = {
      id: Math.max(...notes.map((n: any) => n.id), 0) + 1,
      title: newNoteTitle,
      content: "<h1>" + newNoteTitle + "</h1><p>Start writing your note...</p>",
      author: "current.user",
      createdAt: new Date(),
      updatedAt: new Date(),
      tags: newNoteTags,
      category: newNoteCategory,
      folderId: selectedFolder,
      gitBranch: "main",
      gitCommit: "current",
    };

    setNotes([...notes, newNote]);
    setSelectedNote(newNote);
    setIsEditing(true);
    setIsModalOpen(false);
    setNewNoteTitle("");
    setNewNoteTags([]);
    setError("");
  };

  const handleCreateFolder = () => {
    if (!newFolderName.trim()) {
      setError("Folder name is required");
      return;
    }
    if (
      folders.some((f) => f.name.toLowerCase() === newFolderName.toLowerCase())
    ) {
      setError("Folder name already exists");
      return;
    }

    const newFolder = {
      id: Math.max(...folders.map((f) => f.id), 0) + 1,
      name: newFolderName,
      parentId: newFolderParent ? parseInt(newFolderParent) : null,
      expanded: true,
    };

    setFolders([...folders, newFolder]);
    setIsFolderModalOpen(false);
    setNewFolderName("");
    setNewFolderParent("");
    setError("");
  };

  const handleSaveNote = () => {
    if (selectedNote) {
      setLoading(true);
      setTimeout(() => {
        setNotes(
          notes.map((note: any) =>
            note.id === selectedNote.id ? selectedNote : note
          )
        );
        setIsEditing(false);
        setLoading(false);
      }, 500);
    }
  };

  const handleDeleteNote = () => {
    if (noteToDelete) {
      setNotes(notes.filter((note: any) => note.id !== noteToDelete.id));
      if (selectedNote?.id === noteToDelete.id) {
        setSelectedNote(null);
      }
      setIsDeleteModalOpen(false);
      setNoteToDelete(null);
    }
  };

  const handleTagsEdit = () => {
    setTempTags([...selectedNote.tags]);
    setEditingTags(true);
  };

  const handleTagsSave = () => {
    if (tempTags.length === 0) {
      setError("At least one tag is required");
      return;
    }
    setSelectedNote({
      ...selectedNote,
      tags: tempTags,
      updatedAt: new Date(),
    });
    setEditingTags(false);
    setError("");
  };

  const toggleFolder = (folderId: any) => {
    setFolders(
      folders.map((folder) =>
        folder.id === folderId
          ? { ...folder, expanded: !folder.expanded }
          : folder
      )
    );
  };

  const getFolderHierarchy = () => {
    const folderMap = new Map();
    folders.forEach((folder) =>
      folderMap.set(folder.id, { ...folder, children: [] })
    );

    const rootFolders: any = [];
    folders.forEach((folder) => {
      if (folder.parentId) {
        const parent = folderMap.get(folder.parentId);
        if (parent) {
          parent.children.push(folderMap.get(folder.id));
        }
      } else {
        rootFolders.push(folderMap.get(folder.id));
      }
    });

    return rootFolders;
  };

  const renderFolder = (folder: any, level = 0) => {
    const folderNotes = notes.filter(
      (note: any) => note.folderId === folder.id
    );
    const hasNotes = folderNotes.length > 0;
    const hasChildren = (folder.children?.length ?? 0) > 0;
    const isSelected = selectedFolder === folder.id;

    return (
      <div key={folder.id}>
        <UnstyledButton
          p="xs"
          style={{
            width: "100%",
            paddingLeft: `${12 + level * 20}px`,
          }}
          onClick={() => setSelectedFolder(isSelected ? null : folder.id)}
        >
          <Group gap="xs">
            <IconFolder size={16} color={isSelected ? "#339af0" : "#868e96"} />

            <Text size="sm" fw={isSelected ? 600 : 400}>
              {folder.name}
            </Text>

            {hasChildren && (
              <ActionIcon
                variant="transparent"
                size="xs"
                onClick={(e) => {
                  e.stopPropagation();
                  toggleFolder(folder.id);
                }}
              >
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
          {folder.children?.map((child: any) => renderFolder(child, level + 1))}
        </Collapse>
      </div>
    );
  };

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

  const folderOptions = folders
    .filter((f) => !f.parentId)
    .map((f) => ({ value: f.id.toString(), label: f.name }));

  return (
    <Container
      fluid
      style={{
        height: "calc(100vh - 86px)",
        display: "flex",
        flexDirection: "column",
        overflow: "hidden",
      }}
    >
      {/* Top Header */}
      <Box
        p="xs"
        pt="0"
        style={{
          borderBottom: "1px solid var(--mantine-color-gray-8)",
        }}
      >
        <Group justify="space-between">
          <Title order={2} size="h3">
            <IconFileText
              size={24}
              style={{ marginRight: "8px", verticalAlign: "middle" }}
            />
            Notes
          </Title>
          <Group gap="sm">
            <Group gap="sm">
              <TextInput
                flex={1}
                placeholder="Search notes..."
                value={searchQuery}
                onChange={(e) => debouncedSearch(e.currentTarget.value)}
                leftSection={<IconSearch size={16} />}
              />
              <TagsInput
                flex={1}
                placeholder="Search by tags..."
                value={searchTags}
                onChange={(value) => setSearchTags(value as any[])}
                data={allTags}
                leftSection={<IconTags size={16} />}
                clearable
              />
              <Select
                flex={1}
                placeholder="Filter by category"
                value={filterCategory}
                onChange={(value) => setFilterCategory(value!)}
                data={[{ value: "", label: "All Categories" }, ...categories]}
                leftSection={<IconFilter size={16} />}
                clearable
              />
            </Group>
            <RoleGuard.Consumer>
              <Button
                variant="light"
                size="sm"
                leftSection={<IconFolderPlus size={16} />}
                onClick={() => setIsFolderModalOpen(true)}
              >
                New Folder
              </Button>
              <Button
                size="sm"
                leftSection={<IconPlus size={16} />}
                onClick={() => setIsModalOpen(true)}
              >
                New Note
              </Button>
            </RoleGuard.Consumer>
          </Group>
        </Group>
      </Box>

      <Box style={{ display: "flex", flex: 1, overflow: "hidden" }}>
        {/* Sidebar */}
        <div
          style={{
            width: "350px",
            borderRight: "1px solid var(--mantine-color-gray-8)",
            display: "flex",
            flexDirection: "column",
          }}
        >
          {/* Folders */}
          <Box
            p="xs"
            style={{
              borderBottom: "1px solid var(--mantine-color-gray-8)",
            }}
          >
            <Group justify="space-between" mb="sm">
              <Text size="sm" fw={600} c="dimmed">
                FOLDERS
              </Text>
              <RoleGuard.Consumer>
                <ActionIcon
                  variant="subtle"
                  size="sm"
                  onClick={() => setIsFolderModalOpen(true)}
                >
                  <IconFolderPlus size={14} />
                </ActionIcon>
              </RoleGuard.Consumer>
            </Group>

            <UnstyledButton
              style={{
                width: "100%",
                padding: "8px 12px",
                borderRadius: "6px",
                marginBottom: "4px",
              }}
              onClick={() => setSelectedFolder(null)}
            >
              <Group gap="xs">
                <IconFileText
                  size={16}
                  color={selectedFolder === null ? "#339af0" : "#868e96"}
                />
                <Text size="sm" fw={selectedFolder === null ? 600 : 400}>
                  All Notes
                </Text>
                <Badge size="xs" variant="light" color="gray">
                  {notes.length}
                </Badge>
              </Group>
            </UnstyledButton>

            {getFolderHierarchy().map((folder: any) => renderFolder(folder))}
          </Box>

          {/* Notes List */}
          <ScrollArea style={{ flex: 1 }}>
            <Box p="xs" pr="sm">
              {filteredNotes.map((note: any) => (
                <Card
                  key={note.id}
                  padding="xs"
                  mb="xs"
                  style={{
                    cursor: "pointer",
                    border:
                      selectedNote?.id === note.id
                        ? "2px solid #339af0"
                        : "1px solid var(--mantine-color-gray-8)",
                  }}
                  onClick={() => {
                    setSelectedNote(note);
                    setIsEditing(false);
                  }}
                >
                  <Group justify="space-between" mb="xs">
                    <Text fw={500} size="sm" style={{ flex: 1 }} truncate>
                      {note.title}
                    </Text>
                    <RoleGuard.Consumer>
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
                            leftSection={<IconEdit size={14} />}
                            onClick={(e) => {
                              e.stopPropagation();
                              setSelectedNote(note);
                              setIsEditing(true);
                            }}
                          >
                            Edit
                          </Menu.Item>
                          <Menu.Item
                            leftSection={<IconTrash size={14} />}
                            color="red"
                            onClick={(e) => {
                              e.stopPropagation();
                              setNoteToDelete(note);
                              setIsDeleteModalOpen(true);
                            }}
                          >
                            Delete
                          </Menu.Item>
                        </Menu.Dropdown>
                      </Menu>
                    </RoleGuard.Consumer>
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
                    <Group gap="xs">
                      <Avatar size={16} color="blue">
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

              {filteredNotes.length === 0 && (
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
                      onClick={() => setIsModalOpen(true)}
                    >
                      Create a note
                    </Button>
                  </Stack>
                </Card>
              )}
            </Box>
          </ScrollArea>
        </div>

        {/* Main Content */}
        <div
          style={{
            flex: 1,
            display: "flex",
            flexDirection: "column",
            overflow: "hidden",
          }}
        >
          {selectedNote ? (
            <>
              {/* Note Header */}
              <Box
                p="xs"
                style={{
                  borderBottom: "1px solid var(--mantine-color-gray-8)",
                  minHeight: "auto",
                }}
              >
                <Group justify="space-between" mb="sm">
                  <div style={{ flex: 1 }}>
                    <Title order={1} size="h2" mb="xs">
                      {selectedNote.title}
                    </Title>
                    <Group gap="sm">
                      <Group gap="xs">
                        <Avatar size={24} color="blue">
                          {selectedNote.author.charAt(0).toUpperCase()}
                        </Avatar>
                        <div>
                          <Text size="sm" fw={500}>
                            {selectedNote.author}
                          </Text>
                          <Text size="xs" c="dimmed">
                            Created {formatDate(selectedNote.createdAt)}
                            {selectedNote.updatedAt > selectedNote.createdAt &&
                              ` • Updated ${formatDate(
                                selectedNote.updatedAt
                              )}`}
                          </Text>
                        </div>
                      </Group>
                      <Divider orientation="vertical" />
                      <Group gap="xs">
                        <IconGitBranch size={16} color="#868e96" />
                        <Text size="sm" c="dimmed">
                          {selectedNote.gitBranch}
                        </Text>
                        <Text size="xs" c="dimmed">
                          ({selectedNote.gitCommit})
                        </Text>
                      </Group>
                    </Group>
                  </div>

                  <RoleGuard.Consumer>
                    <Group gap="sm">
                      {!isEditing ? (
                        <Button
                          leftSection={<IconEdit size={16} />}
                          onClick={() => setIsEditing(true)}
                        >
                          Edit
                        </Button>
                      ) : (
                        <Group gap="sm">
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={() => {
                              setIsEditing(false);
                              editor?.commands.setContent(selectedNote.content);
                              setError("");
                            }}
                          >
                            Cancel
                          </Button>
                          <Button
                            size="sm"
                            onClick={handleSaveNote}
                            loading={loading}
                          >
                            Save
                          </Button>
                        </Group>
                      )}
                    </Group>
                  </RoleGuard.Consumer>
                </Group>

                {/* Tags and Category */}
                <Group gap="sm" mb="sm">
                  <Badge
                    color={getCategoryColor(selectedNote.category)}
                    variant="light"
                  >
                    {selectedNote.category}
                  </Badge>

                  {!editingTags ? (
                    <Group gap="xs">
                      {selectedNote.tags.map((tag: any) => (
                        <Badge
                          key={tag}
                          color={tagColors[tag] || "gray"}
                          variant="outline"
                          size="sm"
                        >
                          {tag}
                        </Badge>
                      ))}
                      <RoleGuard.Consumer>
                        <ActionIcon
                          variant="subtle"
                          size="sm"
                          onClick={handleTagsEdit}
                        >
                          <IconEdit size={12} />
                        </ActionIcon>
                      </RoleGuard.Consumer>
                    </Group>
                  ) : (
                    <Group gap="xs" style={{ flex: 1 }}>
                      <TagsInput
                        value={tempTags}
                        onChange={setTempTags}
                        data={allTags}
                        placeholder="Add tags..."
                        size="xs"
                        style={{ flex: 1 }}
                        error={error}
                      />
                      <ActionIcon
                        variant="subtle"
                        color="green"
                        size="sm"
                        onClick={handleTagsSave}
                      >
                        <IconCheck size={12} />
                      </ActionIcon>
                      <ActionIcon
                        variant="subtle"
                        color="red"
                        size="sm"
                        onClick={() => {
                          setEditingTags(false);
                          setError("");
                        }}
                      >
                        <IconX size={12} />
                      </ActionIcon>
                    </Group>
                  )}
                </Group>
                {error && (
                  <Alert color="red" variant="light">
                    {error}
                  </Alert>
                )}
              </Box>

              {/* Editor/Content */}
              <div
                style={{
                  flex: 1,
                  display: "flex",
                  flexDirection: "column",
                  overflow: "hidden",
                }}
              >
                <RichTextEditor
                  editor={editor}
                  style={{
                    border: "none",
                    display: "flex",
                    flexDirection: "column",
                    height: "100%",
                  }}
                >
                  <RoleGuard.Consumer>
                    {isEditing && (
                      <RichTextEditor.Toolbar
                        sticky
                        stickyOffset={0}
                        style={{
                          borderBottom: "1px solid var(--mantine-color-gray-8)",
                          zIndex: 99,
                          padding: "8px 16px",
                        }}
                      >
                        <RichTextEditor.ControlsGroup>
                          <RichTextEditor.Bold />
                          <RichTextEditor.Italic />
                          <RichTextEditor.Underline />
                          <RichTextEditor.Strikethrough />
                          <RichTextEditor.ClearFormatting />
                          <RichTextEditor.Highlight />
                          <RichTextEditor.Code />
                        </RichTextEditor.ControlsGroup>

                        <RichTextEditor.ControlsGroup>
                          <RichTextEditor.H1 />
                          <RichTextEditor.H2 />
                          <RichTextEditor.H3 />
                          <RichTextEditor.H4 />
                        </RichTextEditor.ControlsGroup>

                        <RichTextEditor.ControlsGroup>
                          <RichTextEditor.Blockquote />
                          <RichTextEditor.Hr />
                          <RichTextEditor.BulletList />
                          <RichTextEditor.OrderedList />
                          <RichTextEditor.Subscript />
                          <RichTextEditor.Superscript />
                        </RichTextEditor.ControlsGroup>

                        <RichTextEditor.ControlsGroup>
                          <RichTextEditor.Link />
                          <RichTextEditor.Unlink />
                        </RichTextEditor.ControlsGroup>

                        <RichTextEditor.ControlsGroup>
                          <RichTextEditor.AlignLeft />
                          <RichTextEditor.AlignCenter />
                          <RichTextEditor.AlignJustify />
                          <RichTextEditor.AlignRight />
                        </RichTextEditor.ControlsGroup>

                        <RichTextEditor.ControlsGroup>
                          <RichTextEditor.Undo />
                          <RichTextEditor.Redo />
                        </RichTextEditor.ControlsGroup>

                        <RichTextEditor.ControlsGroup>
                          <RichTextEditor.CodeBlock />
                        </RichTextEditor.ControlsGroup>
                      </RichTextEditor.Toolbar>
                    )}
                  </RoleGuard.Consumer>

                  <RichTextEditor.Content
                    style={{
                      flex: 1,
                      padding: "24px",
                      fontSize: "16px",
                      lineHeight: "1.6",
                      overflow: "auto",
                    }}
                  />
                </RichTextEditor>
              </div>
            </>
          ) : (
            /* Welcome State */
            <div
              style={{
                flex: 1,
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
              }}
            >
              <Stack align="center" gap="xs">
                <IconFileText size={46} color="#ced4da" />
                <div style={{ textAlign: "center" }}>
                  <RoleGuard.Consumer>
                    <Title order={2} c="dimmed" mb="sm">
                      Select a note to view or edit
                    </Title>
                  </RoleGuard.Consumer>

                  <RoleGuard.Investor>
                    <Title order={2} c="dimmed" mb="sm">
                      Select a note to view
                    </Title>
                  </RoleGuard.Investor>

                  <RoleGuard.Consumer>
                    <Text c="dimmed" size="sm">
                      Choose a note from the sidebar or create a new one to get
                      started
                    </Text>
                  </RoleGuard.Consumer>
                </div>
                <RoleGuard.Consumer>
                  <Button
                    size="sm"
                    leftSection={<IconPlus size={20} />}
                    onClick={() => setIsModalOpen(true)}
                  >
                    Create New Note
                  </Button>
                </RoleGuard.Consumer>
              </Stack>
            </div>
          )}
        </div>
      </Box>

      {/* Create Note Modal */}
      <Modal
        opened={isModalOpen}
        onClose={() => {
          setIsModalOpen(false);
          setError("");
        }}
        title="Create New Note"
        size="sm"
      >
        <Stack gap="sm">
          <TextInput
            label="Note Title"
            placeholder="Enter note title..."
            value={newNoteTitle}
            onChange={(e) => setNewNoteTitle(e.currentTarget.value)}
            data-autofocus
            error={error}
          />
          <Select
            label="Category"
            value={newNoteCategory}
            onChange={setNewNoteCategory}
            data={categories}
          />
          <Select
            label="Folder"
            placeholder="Select a folder (optional)"
            value={selectedFolder?.toString() || ""}
            onChange={(value) =>
              setSelectedFolder(value ? parseInt(value) : null)
            }
            data={[{ value: "", label: "No Folder" }, ...folderOptions]}
            clearable
          />
          <TagsInput
            label="Tags"
            placeholder="Add tags..."
            value={newNoteTags}
            onChange={setNewNoteTags}
            data={allTags}
          />
          {error && (
            <Alert color="red" variant="light">
              {error}
            </Alert>
          )}
          <Group justify="flex-end" gap="sm">
            <Button
              variant="outline"
              onClick={() => {
                setIsModalOpen(false);
                setError("");
              }}
            >
              Cancel
            </Button>
            <Button onClick={handleCreateNote} disabled={!newNoteTitle.trim()}>
              Create Note
            </Button>
          </Group>
        </Stack>
      </Modal>

      {/* Create Folder Modal */}
      <Modal
        opened={isFolderModalOpen}
        onClose={() => {
          setIsFolderModalOpen(false);
          setError("");
        }}
        title="Create New Folder"
        size="sm"
      >
        <Stack gap="sm">
          <TextInput
            label="Folder Name"
            placeholder="Enter folder name..."
            value={newFolderName}
            onChange={(e) => setNewFolderName(e.currentTarget.value)}
            data-autofocus
            error={error}
          />
          <Select
            label="Parent Folder"
            placeholder="Select parent folder (optional)"
            value={newFolderParent}
            onChange={(value) => setNewFolderParent(value!)}
            data={[
              { value: "", label: "Root Level" },
              ...folders
                .filter((f) => !f.parentId)
                .map((f) => ({
                  value: f.id.toString(),
                  label: f.name,
                })),
            ]}
            clearable
          />
          {error && (
            <Alert color="red" variant="light">
              {error}
            </Alert>
          )}
          <Group justify="flex-end" gap="sm">
            <Button
              variant="outline"
              onClick={() => {
                setIsFolderModalOpen(false);
                setError("");
              }}
            >
              Cancel
            </Button>
            <Button
              onClick={handleCreateFolder}
              disabled={!newFolderName.trim()}
            >
              Create Folder
            </Button>
          </Group>
        </Stack>
      </Modal>

      {/* Delete Confirmation Modal */}
      <Modal
        opened={isDeleteModalOpen}
        onClose={() => {
          setIsDeleteModalOpen(false);
          setError("");
        }}
        title="Delete Note"
        size="sm"
      >
        <Stack gap="sm">
          <Group gap="sm">
            <IconAlertTriangle size={20} color="#fa5252" />
            <div>
              <Text fw={500}>Are you sure you want to delete this note?</Text>
              <Text size="sm" c="dimmed">
                "{noteToDelete?.title}"
              </Text>
            </div>
          </Group>
          <Alert color="red" variant="light">
            This action cannot be undone. The note will be permanently deleted.
          </Alert>
          <Group justify="flex-end" gap="sm">
            <Button
              variant="outline"
              onClick={() => {
                setIsDeleteModalOpen(false);
                setError("");
              }}
            >
              Cancel
            </Button>
            <Button color="red" onClick={handleDeleteNote}>
              Delete Note
            </Button>
          </Group>
        </Stack>
      </Modal>
    </Container>
  );
}
