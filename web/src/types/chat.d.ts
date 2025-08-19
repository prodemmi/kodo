export interface ProjectFile {
  id: string;
  name: string;
  path: string;
  type: "file" | "folder";
  children?: ProjectFile[];
}
