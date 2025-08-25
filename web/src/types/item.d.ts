export type ItemStatus = "todo" | "in_progress" | "done";
export type ItemPriority = "low" | "medium" | "high";

export type ItemType =
  | "REFACTOR"
  | "OPTIMIZE"
  | "CLEANUP"
  | "DEPRECATED"
  | "BUG"
  | "FIXME"
  | "TODO"
  | "FEATURE"
  | "ENHANCE"
  | "DOC"
  | "TEST"
  | "EXAMPLE"
  | "SECURITY"
  | "COMPLIANCE"
  | "DEBT"
  | "ARCHITECTURE"
  | "CONFIG"
  | "DEPLOY"
  | "MONITOR"
  | "NOTE"
  | "QUESTION"
  | "IDEA"
  | "REVIEW";

export interface StatusHistory {
  status: ItemStatus;
  timestamp: string; // ISO string
  user: string;
}

export interface Item {
  id: number;
  type: ItemType;
  title: string;
  description: string;
  file: string;
  line: number;
  status: ItemStatus;
  priority: ItemPriority;

  // Track status changes over time
  history?: StatusHistory[];

  // Computed / convenience fields
  created_at: string; // ISO string
  updated_at: string; // ISO string
  current_user?: string;

  is_done?: boolean;
  done_at?: string;
  done_by?: string;

  is_in_progress?: boolean;
  in_progress_at?: string;
  in_progress_by?: string;

  full_title?: string;
}

export interface ItemContext {
  file?: string;
  itemLine?: number; // Changed from todoLine to match backend
  lines?: CodeLine[];
  codeElements?: CodeElement[]; // NEW: Code elements found after TODO
  error?: string;
}

export interface CodeLine {
  number: number;
  content: string;
  isTodo: boolean;
}

export interface CodeElement {
  type: string; // function, variable, struct, class, interface, etc.
  name: string; // element name
  line: number; // line number where element is defined
  language: string; // programming language
  signature: string; // full signature/declaration
  context: string; // additional context info (e.g., "function body starts")
  visibility: string; // public, private, protected, etc. (optional)
}

export interface OpenFileParams {
  filename: string;
  line: number;
}

export interface OpenFileResponse {
  status: string;
  opened: boolean;
}

export interface UpdateItemParams {
  id: number;
  status: string;
}

export interface UpdateItemResponse {
  status: string;
}
