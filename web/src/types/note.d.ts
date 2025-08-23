export interface Note {
  id: number;
  title: string;
  content: string;
  author: string;
  createdAt: Date;
  updatedAt: Date;
  tags: string[];
  category: string;
  folderId: number | null;
  gitBranch?: string;
  gitCommit?: string;
  pinned?: boolean;
}

export interface Folder {
  id: number;
  name: string;
  parentId: number | null;
  expanded: boolean;
}

export interface Category {
  value: string;
  label: string;
}

export type TagColor =
  | "blue"
  | "red"
  | "green"
  | "purple"
  | "orange"
  | "pink"
  | "yellow"
  | "cyan"
  | "teal"
  | "lime";

export interface CreateNoteParams {
  title: string;
  content: string;
  author: string;
  tags: string[];
  category: string;
  folderId: number | null;
}

export interface UpdateNoteParams {
  id: number;
  title: string;
  content: string;
  author: string;
  tags: string[];
  category: string;
  pinned: boolean;
  folderId: number | null;
}

export interface CreateFolderParams {
  name: string;
  parentId: number | null;
}

export interface UpdateFolderParams {
  id: number;
  name: string;
  parentId: number | null;
  expanded: boolean;
}

export interface NoteSearchParams {
  q: string | null;
  category: string | null;
  folderId: number | null;
}

export interface MoveNotesParams {
  noteIds: number[];
  targetFolderId: number | null;
}

export interface ExportNotesParams {
  category: string;
  folderId: string | null;
}

export interface NoteStats {
  total_notes: number;
  total_folders: number;
  by_category: Record<string, number>;
  by_author: Record<string, number>;
  by_month: Record<string, number>;
  recent_notes: Note[];
}

export interface FolderTree {
  id: number;
  name: string;
  parentId: number | null;
  expanded: boolean;
  noteCount: number;
  totalNoteCount: number;
  children: FolderStats[];
  hasChildren: boolean;
}

export interface NoteHistoryStats {
  total_entries: number;
  by_action: Record<string, number>;
  by_author: Record<string, number>;
  by_branch: Record<string, number>;
  by_day: Record<string, number>;
  most_active_notes: MostActiveNote[];
  recent_activity: RecentActivity[];
}

export interface MostActiveNote {
  note_id: number;
  note_title: string;
  action_count: number;
  last_action: string; // ISO date string
  authors: string[];
}

export interface RecentActivity {
  id: number;
  note_id: number;
  action: string;
  author: string;
  timestamp: string;
  git_branch: string;
  git_commit: string;
  changes: {
    [key: string]: {
      from: number;
      to: number;
    };
  };
  old_value: NoteSnapshot;
  new_value: NoteSnapshot;
  message: string;
}

export interface NoteSnapshot {
  id: number;
  title: string;
  content: string;
  author: string;
  category: string;
  tags: string[];
  folderId: number | null;
  gitBranch: string;
  gitCommit: string;
  createdAt: string;
  updatedAt: string;
}

export interface NoteHistoryResponse {
  count: number;
  filter: {
    note_id: number;
    limit: number;
  };
  history: NoteHistory[];
}

export interface NoteHistory {
  id: number;
  note_id: number;
  action: string;
  author: string;
  timestamp: string;
  git_branch: string;
  git_commit: string;
  changes: Record<string, { from: number | string; to: number | string }>;
  old_value: NoteSnapshot;
  new_value: NoteSnapshot;
  message: string;
}

export interface NoteSnapshot {
  author: string;
  category: string;
  content: string;
  createdAt: string;
  folderId: number | null;
  gitBranch: string;
  gitCommit: string;
  id: number;
  tags: string[];
  title: string;
  updatedAt: string;
}
