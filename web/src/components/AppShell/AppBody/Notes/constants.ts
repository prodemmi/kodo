import { Category, Folder, Note } from "../../../../types/note";

export const mockFolders: Folder[] = [
  { id: 1, name: "Architecture", parentId: null, expanded: true },
  { id: 2, name: "Meetings", parentId: null, expanded: true },
  { id: 3, name: "Ideas", parentId: null, expanded: false },
  { id: 4, name: "Sprint Planning", parentId: 2, expanded: true },
];

export const mockNotes: Note[] = [
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

export const categories: Category[] = [
  { value: "technical", label: "Technical" },
  { value: "feature", label: "Feature" },
  { value: "bug", label: "Bug" },
  { value: "idea", label: "Idea" },
  { value: "improvement", label: "Improvement" },
  { value: "documentation", label: "Documentation" },
  { value: "review", label: "Review" },
  { value: "meeting", label: "Meeting" },
];

export const tagColors: any = {
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
