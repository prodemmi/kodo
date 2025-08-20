import { ProjectFile } from "./chat";

export interface CodeSnippet {
  id: string;
  fileName: string;
  filePath: string;
  content: string;
  startLine?: number;
  endLine?: number;
  language?: string;
}

export interface Context {
  id: string;
  name: string;
  description?: string;
  files: ProjectFile[];
  snippets: CodeSnippet[];
  createdAt: Date;
  updatedAt: Date;
}
