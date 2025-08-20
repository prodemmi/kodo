export interface ProjectFile {
  id: string;
  name: string;
  path: string;
  size: number;
  type: "file" | "folder";
  children?: ProjectFile[];
}

export interface Message {
  id: string;
  type: "user" | "assistant";
  content: string;
  timestamp: Date;
}
