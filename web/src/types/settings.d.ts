export type KanbanColumn = {
  id: string;
  name: string;
  color: string;
  auto_assign_pattern?: string;
};

export type PriorityPatterns = {
  low: string;
  medium: string;
  high: string;
};

export type GithubAuth = {
  token: string;
};

export type WorkspaceSettings = {
  theme: "auto" | "light" | "dark";
  primary_color: string;
  show_line_preview: boolean;
};

export type CodeScanSettings = {
  exclude_directories: string[];
  exclude_files: string[];
  sync_enabled: boolean;
};

export type Settings = {
  kanban_columns: KanbanColumn[];
  priority_patterns: PriorityPatterns;
  github_auth: GithubAuth;
  code_scan_settings: CodeScanSettings;
};
