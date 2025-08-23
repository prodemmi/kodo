import { useState, useContext, createContext, useEffect } from "react";
import {
  Container,
  Title,
  Tabs,
  Stack,
  Group,
  Button,
  Select,
  Switch,
  NumberInput,
  Text,
  Paper,
  TextInput,
  ActionIcon,
  ColorSwatch,
  Modal,
  Textarea,
  Badge,
  Card,
  MultiSelect,
  Divider,
  Code,
  Alert,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import {
  IconPlus,
  IconTrash,
  IconFolder,
  IconNote,
  IconSearch,
  IconBrandGithub,
  IconCheck,
  IconAlertCircle,
} from "@tabler/icons-react";
import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  DragEndEvent,
} from "@dnd-kit/core";
import {
  SortableContext,
  sortableKeyboardCoordinates,
  verticalListSortingStrategy,
  useSortable,
  arrayMove,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { useAppState } from "../../../../states/app.state";

interface UserSettings {
  // Personal productivity
  preferredIde: string;
  autoSaveInterval: number;
  theme: string;
  defaultAssignee: string;
  personalKeywords: string[];
  showLinePreview: boolean;
}

interface CodeScanSettings {
  // Comment parsing
  commentKeywords: string[];
  fileExtensions: string[];
  excludeDirectories: string[];
  excludeFiles: string[];
  caseSensitive: boolean;

  // Auto-categorization patterns
  selectedPatterns: string[];

  // Notes template
  notesTemplate: string;
}

interface WorkspaceSettings {
  // Project configuration
  projectName: string;
  primaryColor: string;

  // Provider sync
  syncEnabled: boolean;
  syncProvider: string;
  autoCreateIssues: boolean;
  syncBranches: boolean;
  issueLabels: string[];

  // Git integration
  autoCommit: boolean;
  commitTemplate: string;
  branchNaming: string;

  // Team collaboration
  teamMembers: string[];
}

interface KanbanColumn {
  id: string;
  name: string;
  color: string;
  autoAssignPattern?: string;
}

interface UserContextType {
  userId: string;
  username: string;
  settings: UserSettings;
  updateSettings: (settings: Partial<UserSettings>) => void;
}

interface WorkspaceContextType {
  workspaceId: string;
  name: string;
  codeScanSettings: CodeScanSettings;
  workspaceSettings: WorkspaceSettings;
  updateCodeScanSettings: (settings: Partial<CodeScanSettings>) => void;
  updateWorkspaceSettings: (settings: Partial<WorkspaceSettings>) => void;
  kanbanColumns: KanbanColumn[];
  setKanbanColumns: (columns: KanbanColumn[]) => void;
}

// Create contexts with realistic defaults
const UserContext = createContext<UserContextType>({
  userId: "user-1",
  username: "prodemmi",
  settings: {
    preferredIde: "vscode",
    autoSaveInterval: 5,
    theme: "auto",
    defaultAssignee: "prodemmi",
    personalKeywords: ["MYFIX", "PERSONAL"],
    showLinePreview: true,
  },
  updateSettings: () => {},
});

const WorkspaceContext = createContext<WorkspaceContextType>({
  workspaceId: "ws-1",
  name: "My Project",
  codeScanSettings: {
    commentKeywords: ["TODO", "FIXME", "BUG", "FEAT", "NOTE"],
    fileExtensions: [
      ".js",
      ".ts",
      ".jsx",
      ".tsx",
      ".py",
      ".java",
      ".cpp",
      ".c",
    ],
    excludeDirectories: ["node_modules", ".git", "dist", "build"],
    excludeFiles: ["*.min.js", "*.bundle.js", "*.d.ts"],
    caseSensitive: false,
    selectedPatterns: ["basic-todo", "status-tracking", "feature-request"],
    notesTemplate: "# {title}\n\n## Overview\n\n## Tasks\n- [ ] \n\n## Notes\n",
  },
  workspaceSettings: {
    projectName: "My Awesome Project",
    primaryColor: "blue",
    syncEnabled: false,
    syncProvider: "github",
    autoCreateIssues: true,
    syncBranches: true,
    issueLabels: [
      "ENHANCEMENT",
      "BUG",
      "TODO",
      "FEATURE|FEAT",
      "DOC|DOCUMENTATION",
      "HELP",
    ],
    autoCommit: false,
    commitTemplate: "feat: {description}",
    branchNaming: "feature/{ticket}-{description}",
    teamMembers: ["prodemmi"],
  },
  updateCodeScanSettings: () => {},
  updateWorkspaceSettings: () => {},
  kanbanColumns: [
    { id: "col-1", name: "Backlog", color: "gray", autoAssignPattern: "TODO:" },
    {
      id: "col-2",
      name: "In Progress",
      color: "blue",
      autoAssignPattern: "IN PROGRESS",
    },
    {
      id: "col-3",
      name: "Features",
      color: "green",
      autoAssignPattern: "FEAT:",
    },
    { id: "col-4", name: "Done", color: "teal", autoAssignPattern: "DONE" },
  ],
  setKanbanColumns: () => {},
});

// Pattern examples for code scanning
const commentPatterns = [
  {
    id: "basic-todo",
    name: "Basic TODO",
    description: "Simple todo comments",
    example: "// TODO: implement in zustand and sync with onSuccess fn",
    keywords: ["TODO:", "TO DO:"],
  },
  {
    id: "status-tracking",
    name: "Status Tracking",
    description: "Comments with status and timestamps",
    example: "// IN PROGRESS from 2025-08-22 20:52 by prodemmi",
    keywords: ["IN PROGRESS", "DONE", "BLOCKED"],
  },
  {
    id: "feature-request",
    name: "Feature Requests",
    description: "Feature implementation comments",
    example:
      "// FEAT: store collapse\n// The collapse should store in localStorage (using zustand persist)",
    keywords: ["FEAT:", "FEATURE:"],
  },
  {
    id: "bug-fix",
    name: "Bug Fixes",
    description: "Bug and fix comments",
    example:
      "// FIXME: memory leak in component unmount\n// BUG: validation not working properly",
    keywords: ["FIXME:", "BUG:", "FIX:"],
  },
  {
    id: "notes",
    name: "Notes & Documentation",
    description: "General notes and documentation",
    example: "// NOTE: This function needs refactoring for performance",
    keywords: ["NOTE:", "DOCS:", "INFO:"],
  },
];

// Enhanced sortable column component with proper drag handling
interface SortableColumnProps {
  column: KanbanColumn;
  onNameChange: (id: string, name: string) => void;
  onColorChange: (id: string, color: string) => void;
  onPatternChange: (id: string, pattern: string) => void;
  onDelete: (id: string) => void;
  colors: string[];
}

function SortableColumn({
  column,
  onNameChange,
  onColorChange,
  onPatternChange,
  onDelete,
  colors,
}: SortableColumnProps) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: column.id,
  });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  return (
    <Paper
      shadow="xs"
      p="md"
      mt="xs"
      radius="md"
      withBorder
      ref={setNodeRef}
      style={style}
    >
      <Stack gap="sm">
        <Group justify="space-between">
          <Group gap="xs" style={{ flex: 1 }}>
            <div
              {...attributes}
              {...listeners}
              style={{ cursor: "grab", padding: "4px" }}
            >
              <Text size="xs" c="dimmed">
                ⋮⋮
              </Text>
            </div>
            <TextInput
              value={column.name}
              onChange={(event) =>
                onNameChange(column.id, event.currentTarget.value)
              }
              placeholder="Column name"
              size="sm"
              radius="md"
              style={{ flex: 1 }}
            />
          </Group>
          <ActionIcon
            color="red"
            variant="subtle"
            onClick={() => onDelete(column.id)}
          >
            <IconTrash size={16} />
          </ActionIcon>
        </Group>

        <TextInput
          value={column.autoAssignPattern || ""}
          onChange={(event) =>
            onPatternChange(column.id, event.currentTarget.value)
          }
          placeholder="Auto-assign pattern (e.g., 'TODO:', 'IN PROGRESS')"
          size="sm"
          radius="md"
          label="Auto-assign when comment contains:"
        />

        <Group gap="xs">
          <Text size="sm" c="dimmed">
            Color:
          </Text>
          {colors.map((color) => (
            <ColorSwatch
              key={color}
              color={`var(--mantine-color-${color}-6)`}
              size={20}
              style={{ cursor: "pointer" }}
              onClick={() => onColorChange(column.id, color)}
              component="button"
              type="button"
            >
              {column.color === color && (
                <Text c="white" size="xs">
                  ✓
                </Text>
              )}
            </ColorSwatch>
          ))}
        </Group>
      </Stack>
    </Paper>
  );
}

function Settings() {
  const { settings: userSettings, updateSettings: updateUserSettings } =
    useContext(UserContext);
  const {
    codeScanSettings,
    workspaceSettings,
    updateCodeScanSettings,
    updateWorkspaceSettings,
    kanbanColumns,
    setKanbanColumns,
  } = useContext(WorkspaceContext);

  const [activeTab, setActiveTab] = useState<string | null>("workspace");
  const [deleteColumnId, setDeleteColumnId] = useState<string | null>(null);
  const primaryColor = useAppState((s) => s.primaryColor);
  const setPrimaryColor = useAppState((s) => s.setPrimaryColor);

  // User settings form
  const userForm = useForm<UserSettings>({
    initialValues: userSettings,
  });

  // Code scanning form
  const codeScanForm = useForm<CodeScanSettings>({
    initialValues: codeScanSettings,
  });

  // Workspace form
  const workspaceForm = useForm<WorkspaceSettings>({
    initialValues: workspaceSettings,
  });

  // Column form
  const columnForm = useForm<{ name: string; pattern: string }>({
    initialValues: { name: "", pattern: "" },
    validate: {
      name: (value) =>
        value.trim().length === 0 ? "Column name required" : null,
    },
  });

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

  // Handle submissions
  const handleUserSubmit = (values: UserSettings) => {
    updateUserSettings(values);
    console.log("User settings saved:", values);
  };

  const handleCodeScanSubmit = (values: CodeScanSettings) => {
    updateCodeScanSettings(values);
    console.log("Code scan settings saved:", values);
  };

  const handleWorkspaceSubmit = (values: WorkspaceSettings) => {
    updateWorkspaceSettings(values);
    console.log("Workspace settings saved:", values);
  };

  const handleAddColumn = (values: { name: string; pattern: string }) => {
    const newColumn: KanbanColumn = {
      id: `col-${Date.now()}`,
      name: values.name,
      color: "blue",
      autoAssignPattern: values.pattern || undefined,
    };
    setKanbanColumns([...kanbanColumns, newColumn]);
    columnForm.reset();
  };

  const handleColumnNameChange = (id: string, name: string) => {
    setKanbanColumns(
      kanbanColumns.map((col) => (col.id === id ? { ...col, name } : col))
    );
  };

  const handleColumnColorChange = (id: string, color: string) => {
    setKanbanColumns(
      kanbanColumns.map((col) => (col.id === id ? { ...col, color } : col))
    );
  };

  const handleColumnPatternChange = (id: string, pattern: string) => {
    setKanbanColumns(
      kanbanColumns.map((col) =>
        col.id === id
          ? { ...col, autoAssignPattern: pattern || undefined }
          : col
      )
    );
  };

  const handleDeleteColumn = (id: string) => {
    setDeleteColumnId(id);
  };

  const confirmDeleteColumn = () => {
    if (deleteColumnId) {
      setKanbanColumns(
        kanbanColumns.filter((col) => col.id !== deleteColumnId)
      );
      setDeleteColumnId(null);
    }
  };

  // Fixed drag end handler
  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;

    if (!over || active.id === over.id) {
      return;
    }

    setKanbanColumns((columns) => {
      const oldIndex = columns.findIndex((col) => col.id === active.id);
      const newIndex = columns.findIndex((col) => col.id === over.id);

      return arrayMove(columns, oldIndex, newIndex);
    });
  };

  const colors = ["dark", "blue", "orange", "green", "red", "gray"];

  return (
    <Container size="lg" py="xl">
      <Title order={2} mb="xl">
        Project Settings
      </Title>

      <Tabs
        value={activeTab}
        onChange={setActiveTab}
        variant="pills"
        radius="md"
      >
        <Tabs.List mb="lg">
          <Tabs.Tab value="workspace" leftSection={<IconFolder size={16} />}>
            Workspace & Sync
          </Tabs.Tab>
          <Tabs.Tab value="kanban" leftSection={<IconNote size={16} />}>
            Kanban Board
          </Tabs.Tab>
          <Tabs.Tab value="scanning" leftSection={<IconSearch size={16} />}>
            Code Scanning
          </Tabs.Tab>
        </Tabs.List>

        {/* Code Scanning Tab */}
        <Tabs.Panel value="scanning">
          <Paper shadow="sm" p="lg" radius="md" withBorder>
            <Stack gap="lg">
              <Text fw={600} size="lg">
                Code Comment Scanning
              </Text>

              <MultiSelect
                label="Comment Keywords"
                description="Keywords that create kanban items from comments"
                data={[
                  "TODO",
                  "FIXME",
                  "BUG",
                  "FEAT",
                  "NOTE",
                  "HACK",
                  "REFACTOR",
                  "OPTIMIZE",
                  "IN PROGRESS",
                  "DONE",
                  "BLOCKED",
                ]}
                placeholder="Select keywords..."
                {...codeScanForm.getInputProps("commentKeywords")}
              />

              <MultiSelect
                label="File Extensions"
                description="File types to scan for comments"
                data={[
                  ".js",
                  ".ts",
                  ".jsx",
                  ".tsx",
                  ".py",
                  ".java",
                  ".cpp",
                  ".c",
                  ".php",
                  ".rb",
                  ".go",
                  ".rs",
                ]}
                placeholder="Select file types..."
                {...codeScanForm.getInputProps("fileExtensions")}
              />

              <Group grow>
                <div>
                  <Text fw={500} size="sm" mb="xs">
                    Exclude Directories
                  </Text>
                  <Textarea
                    description="Directories to ignore (one per line)"
                    placeholder="node_modules&#10;.git&#10;dist&#10;build"
                    autosize
                    minRows={2}
                    maxRows={6}
                    value={codeScanForm.values.excludeDirectories.join("\n")}
                    onChange={(e) =>
                      codeScanForm.setFieldValue(
                        "excludeDirectories",
                        e.target.value.split("\n").filter(Boolean)
                      )
                    }
                  />
                </div>

                <div>
                  <Text fw={500} size="sm" mb="xs">
                    Exclude Files
                  </Text>
                  <Textarea
                    description="File patterns to ignore (one per line)"
                    placeholder="*.min.js&#10;*.bundle.js&#10;*.d.ts"
                    autosize
                    minRows={2}
                    maxRows={6}
                    value={codeScanForm.values.excludeFiles.join("\n")}
                    onChange={(e) =>
                      codeScanForm.setFieldValue(
                        "excludeFiles",
                        e.target.value.split("\n").filter(Boolean)
                      )
                    }
                  />
                </div>
              </Group>

              <Switch
                label="Case Sensitive Matching"
                description="Match keywords exactly as typed"
                {...codeScanForm.getInputProps("caseSensitive", {
                  type: "checkbox",
                })}
              />

              <Divider label="Comment Patterns" />

              <Text size="sm" c="dimmed">
                Select which comment patterns to recognize:
              </Text>

              <Stack gap="md">
                {commentPatterns.map((pattern) => (
                  <Card
                    key={pattern.id}
                    padding="md"
                    withBorder
                    style={{
                      cursor: "pointer",
                      backgroundColor:
                        codeScanForm.values.selectedPatterns.includes(
                          pattern.id
                        )
                          ? "var(--mantine-color-blue-0)"
                          : undefined,
                      borderColor:
                        codeScanForm.values.selectedPatterns.includes(
                          pattern.id
                        )
                          ? "var(--mantine-color-blue-6)"
                          : undefined,
                    }}
                    onClick={() => {
                      const current = codeScanForm.values.selectedPatterns;
                      const updated = current.includes(pattern.id)
                        ? current.filter((id) => id !== pattern.id)
                        : [...current, pattern.id];
                      codeScanForm.setFieldValue("selectedPatterns", updated);
                    }}
                  >
                    <Group justify="space-between" align="flex-start">
                      <Stack gap="xs" style={{ flex: 1 }}>
                        <Group gap="sm">
                          <Text fw={500}>{pattern.name}</Text>
                          {codeScanForm.values.selectedPatterns.includes(
                            pattern.id
                          ) && (
                            <Badge
                              
                              size="sm"
                              leftSection={<IconCheck size={12} />}
                            >
                              Active
                            </Badge>
                          )}
                        </Group>
                        <Text size="sm" c="dimmed">
                          {pattern.description}
                        </Text>
                        <Code block>{pattern.example}</Code>
                        <Group gap="xs">
                          <Text size="xs" c="dimmed">
                            Keywords:
                          </Text>
                          {pattern.keywords.map((keyword) => (
                            <Badge key={keyword} size="xs" variant="outline">
                              {keyword}
                            </Badge>
                          ))}
                        </Group>
                      </Stack>
                    </Group>
                  </Card>
                ))}
              </Stack>

              <Divider label="Notes Template" />

              <Textarea
                label="Notes Template"
                description="Default template for new project notes"
                minRows={6}
                placeholder="# {title}&#10;&#10;## Overview&#10;&#10;## Tasks&#10;- [ ] &#10;&#10;## Notes"
                {...codeScanForm.getInputProps("notesTemplate")}
              />

              <Group justify="flex-end">
                <Button
                  onClick={() => handleCodeScanSubmit(codeScanForm.values)}
                >
                  Save Scanning Settings
                </Button>
              </Group>
            </Stack>
          </Paper>
        </Tabs.Panel>

        {/* Workspace & Sync Tab */}
        <Tabs.Panel value="workspace">
          <Paper shadow="sm" p="lg" radius="md" withBorder>
            <Stack gap="lg">
              <Text fw={600} size="lg">
                Workspace Configuration
              </Text>

              <Group grow>
                <TextInput
                  label="Project Name"
                  description="Display name for this project"
                  placeholder="My Awesome Project"
                  {...workspaceForm.getInputProps("projectName")}
                />

                <div>
                  <Text fw={500} size="sm" mb="xs">
                    Primary Color
                  </Text>
                  <Group>
                    {colors.map((color) => (
                      <ColorSwatch
                        key={color}
                        color={`var(--mantine-color-${color}-6)`}
                        bd={primaryColor === color ? "2px solid white" : "none"}
                        size={25}
                        style={{ cursor: "pointer" }}
                        onClick={() => setPrimaryColor(color)}
                        component="button"
                        type="button"
                      >
                        {primaryColor === color && (
                          <Text c="white" size="xs">
                            ✓
                          </Text>
                        )}
                      </ColorSwatch>
                    ))}
                  </Group>
                </div>
              </Group>

              <Divider label="Provider Sync" />

              <Alert
                icon={<IconBrandGithub size={16} />}
                title="Sync with Version Control"
                
                variant="light"
              >
                <Text c="var(--mantine-color-dark-0)">
                  Connect your project with GitHub, GitLab, or other providers
                  to automatically create issues from code comments and sync
                  task status with branches.
                </Text>
              </Alert>

              <Group grow align="flex-end">
                <Select
                  label="Sync Provider"
                  description="Choose your version control provider"
                  data={[
                    { value: "github", label: "GitHub" },
                    { value: "gitlab", label: "GitLab" },
                    { value: "bitbucket", label: "Bitbucket" },
                    { value: "azure-devops", label: "Azure DevOps" },
                  ]}
                  {...workspaceForm.getInputProps("syncProvider")}
                />

                <Switch
                  label="Enable Sync"
                  description="Sync tasks with your repository"
                  {...workspaceForm.getInputProps("syncEnabled", {
                    type: "checkbox",
                  })}
                />
              </Group>

              {workspaceForm.values.syncEnabled && (
                <Stack gap="md">
                  <Switch
                    label="Auto-create Issues"
                    description="Create repository issues from TODO comments"
                    {...workspaceForm.getInputProps("autoCreateIssues", {
                      type: "checkbox",
                    })}
                  />

                  <MultiSelect
                    label="Issue Labels"
                    description="Default labels for auto-created issues"
                    data={[
                      "ENHANCEMENT",
                      "BUG",
                      "TODO",
                      "FEATURE|FEAT",
                      "DOC|DOCUMENTATION",
                      "HELP",
                    ]}
                    placeholder="Select default labels..."
                    {...workspaceForm.getInputProps("issueLabels")}
                  />
                </Stack>
              )}

              <Divider label="Git Integration" />

              <Switch
                label="Auto Commit"
                description="Auto-commit when tasks are moved to Done"
                {...workspaceForm.getInputProps("autoCommit", {
                  type: "checkbox",
                })}
              />

              <Group grow>
                <TextInput
                  label="Commit Template"
                  description="Template for automatic commits"
                  placeholder="feat: {description}"
                  {...workspaceForm.getInputProps("commitTemplate")}
                />

                <TextInput
                  label="Branch Naming"
                  description="Template for branch names"
                  placeholder="feature/{ticket}-{description}"
                  {...workspaceForm.getInputProps("branchNaming")}
                />
              </Group>

              <MultiSelect
                label="Team Members"
                description="People who can be assigned tasks"
                data={["prodemmi", "john.doe", "jane.smith"]}
                placeholder="Add team members..."
                {...workspaceForm.getInputProps("teamMembers")}
              />

              <Group justify="flex-end">
                <Button
                  onClick={() => handleWorkspaceSubmit(workspaceForm.values)}
                >
                  Save Workspace Settings
                </Button>
              </Group>
            </Stack>

            <Stack gap="lg">
              <Text fw={600} size="lg">
                Personal Settings
              </Text>

              <Group grow>
                <Select
                  label="Preferred IDE"
                  description="Your main development environment"
                  data={[
                    { value: "vscode", label: "VS Code" },
                    { value: "intellij", label: "IntelliJ IDEA" },
                    { value: "sublime", label: "Sublime Text" },
                    { value: "atom", label: "Atom" },
                    { value: "vim", label: "Vim/Neovim" },
                  ]}
                  {...userForm.getInputProps("preferredIde")}
                />

                <Select
                  label="Theme"
                  description="Interface appearance"
                  data={[
                    { value: "auto", label: "Auto (System)" },
                    { value: "light", label: "Light" },
                    { value: "dark", label: "Dark" },
                  ]}
                  {...userForm.getInputProps("theme")}
                />
              </Group>

              <NumberInput
                label="Auto-save Interval"
                description="Minutes between automatic saves"
                min={1}
                max={30}
                suffix=" min"
                {...userForm.getInputProps("autoSaveInterval")}
              />

              <Switch
                label="Show Line Preview"
                description="Show code context in kanban cards"
                {...userForm.getInputProps("showLinePreview", {
                  type: "checkbox",
                })}
              />

              <Group justify="flex-end">
                <Button onClick={() => handleUserSubmit(userForm.values)}>
                  Save User Settings
                </Button>
              </Group>
            </Stack>
          </Paper>
        </Tabs.Panel>

        {/* Kanban Board Tab */}
        <Tabs.Panel value="kanban">
          <Paper shadow="sm" p="lg" radius="md" withBorder>
            <Stack gap="lg">
              <Text fw={600} size="lg">
                Kanban Board Setup
              </Text>

              <Group align="flex-end">
                <TextInput
                  label="Column Name"
                  placeholder="e.g., Code Review"
                  style={{ flex: 1 }}
                  {...columnForm.getInputProps("name")}
                />
                <TextInput
                  label="Auto-assign Pattern"
                  placeholder="e.g., REVIEW:"
                  style={{ flex: 1 }}
                  {...columnForm.getInputProps("pattern")}
                />
                <Button
                  leftSection={<IconPlus size={16} />}
                  onClick={() => {
                    if (!columnForm.validate().hasErrors) {
                      handleAddColumn(columnForm.values);
                    }
                  }}
                >
                  Add Column
                </Button>
              </Group>

              <Alert
                icon={<IconAlertCircle size={16} />}
                
                variant="light"
              >
                Drag the grip (⋮⋮) to reorder columns. Auto-assign patterns
                automatically move items containing those keywords.
              </Alert>

              <DndContext
                sensors={sensors}
                collisionDetection={closestCenter}
                onDragEnd={handleDragEnd}
              >
                <SortableContext
                  items={kanbanColumns.map((col) => col.id)}
                  strategy={verticalListSortingStrategy}
                >
                  <Stack>
                    {kanbanColumns.map((column) => (
                      <SortableColumn
                        key={column.id}
                        column={column}
                        onNameChange={handleColumnNameChange}
                        onColorChange={handleColumnColorChange}
                        onPatternChange={handleColumnPatternChange}
                        onDelete={handleDeleteColumn}
                        colors={colors}
                      />
                    ))}
                  </Stack>
                </SortableContext>
              </DndContext>
            </Stack>
          </Paper>
        </Tabs.Panel>
      </Tabs>

      {/* Delete Confirmation Modal */}
      <Modal
        opened={!!deleteColumnId}
        onClose={() => setDeleteColumnId(null)}
        title="Delete Column"
        radius="md"
        centered
      >
        <Stack>
          <Text>
            Are you sure you want to delete the column "
            {kanbanColumns.find((col) => col.id === deleteColumnId)?.name}"?
          </Text>
          <Group justify="flex-end">
            <Button variant="outline" onClick={() => setDeleteColumnId(null)}>
              Cancel
            </Button>
            <Button color="red" onClick={confirmDeleteColumn}>
              Delete
            </Button>
          </Group>
        </Stack>
      </Modal>
    </Container>
  );
}

export default Settings;
