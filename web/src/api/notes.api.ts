import { 
  Note, 
  Folder, 
  Category, 
  CreateNoteParams, 
  UpdateNoteParams, 
  CreateFolderParams, 
  UpdateFolderParams,
  NoteSearchParams,
  MoveNotesParams,
  ExportNotesParams,
  NoteStats,
  FolderTree
} from "../types/note";

const API_BASE = "http://localhost:8080/api";

// Note operations
export const getNotes = async (params?: {
  category?: string;
  tag?: string;
  folderId?: string;
}): Promise<{ notes: Note[]; count: number }> => {
  const searchParams = new URLSearchParams();
  if (params?.category) searchParams.append("category", params.category);
  if (params?.tag) searchParams.append("tag", params.tag);
  if (params?.folderId) searchParams.append("folderId", params.folderId);

  const response = await fetch(`${API_BASE}/notes?${searchParams.toString()}`);
  if (!response.ok) {
    throw new Error(`Failed to fetch notes: ${response.statusText}`);
  }
  return await response.json();
};

export const createNote = async (params: CreateNoteParams): Promise<{ status: string; note: Note }> => {
  const response = await fetch(`${API_BASE}/notes`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(params),
  });
  
  if (!response.ok) {
    throw new Error(`Failed to create note: ${response.statusText}`);
  }
  return await response.json();
};

export const updateNote = async (params: UpdateNoteParams): Promise<{ status: string; note: Note }> => {
  const response = await fetch(`${API_BASE}/notes/update`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(params),
  });
  
  if (!response.ok) {
    throw new Error(`Failed to update note: ${response.statusText}`);
  }
  return await response.json();
};

export const deleteNote = async (id: number): Promise<{ status: string; message: string }> => {
  const response = await fetch(`${API_BASE}/notes/delete`, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ id }),
  });
  
  if (!response.ok) {
    throw new Error(`Failed to delete note: ${response.statusText}`);
  }
  return await response.json();
};

export const searchNotes = async (params: NoteSearchParams): Promise<{ notes: Note[]; count: number; query: string }> => {
  const searchParams = new URLSearchParams();
  if (params.q) searchParams.append("q", params.q);
  if (params.category) searchParams.append("category", params.category);
  if (params.folderId) searchParams.append("folderId", params.folderId.toString());

  const response = await fetch(`${API_BASE}/notes/search?${searchParams.toString()}`);
  if (!response.ok) {
    throw new Error(`Failed to search notes: ${response.statusText}`);
  }
  return await response.json();
};

export const moveNotes = async (params: MoveNotesParams): Promise<{ status: string; message: string }> => {
  const response = await fetch(`${API_BASE}/notes/move`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(params),
  });
  
  if (!response.ok) {
    throw new Error(`Failed to move notes: ${response.statusText}`);
  }
  return await response.json();
};

export const exportNotes = async (params?: ExportNotesParams): Promise<Blob> => {
  const searchParams = new URLSearchParams();
  if (params?.category) searchParams.append("category", params.category);
  if (params?.folderId) searchParams.append("folderId", params.folderId.toString());

  const response = await fetch(`${API_BASE}/notes/export?${searchParams.toString()}`);
  if (!response.ok) {
    throw new Error(`Failed to export notes: ${response.statusText}`);
  }
  return await response.blob();
};

export const getNoteStats = async (): Promise<NoteStats> => {
  const response = await fetch(`${API_BASE}/notes/stats`);
  if (!response.ok) {
    throw new Error(`Failed to fetch note stats: ${response.statusText}`);
  }
  return await response.json();
};

export const getNoteTags = async (): Promise<{ tags: string[]; count: number }> => {
  const response = await fetch(`${API_BASE}/notes/tags`);
  if (!response.ok) {
    throw new Error(`Failed to fetch note tags: ${response.statusText}`);
  }
  return await response.json();
};

// Folder operations
export const getFolders = async (): Promise<{ folders: Folder[]; count: number }> => {
  const response = await fetch(`${API_BASE}/folders`);
  if (!response.ok) {
    throw new Error(`Failed to fetch folders: ${response.statusText}`);
  }
  return await response.json();
};

export const getFolderTree = async (): Promise<{ tree: FolderTree[]; count: number }> => {
  const response = await fetch(`${API_BASE}/folders/tree`);
  if (!response.ok) {
    throw new Error(`Failed to fetch folder tree: ${response.statusText}`);
  }
  return await response.json();
};

export const createFolder = async (params: CreateFolderParams): Promise<{ status: string; folder: Folder }> => {
  const response = await fetch(`${API_BASE}/folders`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(params),
  });
  
  if (!response.ok) {
    throw new Error(`Failed to create folder: ${response.statusText}`);
  }
  return await response.json();
};

export const updateFolder = async (params: UpdateFolderParams): Promise<{ status: string; folder: Folder }> => {
  const response = await fetch(`${API_BASE}/folders/update`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(params),
  });
  
  if (!response.ok) {
    throw new Error(`Failed to update folder: ${response.statusText}`);
  }
  return await response.json();
};

export const deleteFolder = async (id: number): Promise<{ status: string; message: string }> => {
  const response = await fetch(`${API_BASE}/folders/delete`, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ id }),
  });
  
  if (!response.ok) {
    throw new Error(`Failed to delete folder: ${response.statusText}`);
  }
  return await response.json();
};

// Categories
export const getCategories = async (): Promise<{ categories: Category[]; count: number }> => {
  const response = await fetch(`${API_BASE}/categories`);
  if (!response.ok) {
    throw new Error(`Failed to fetch categories: ${response.statusText}`);
  }
  return await response.json();
};