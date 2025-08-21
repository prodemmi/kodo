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
