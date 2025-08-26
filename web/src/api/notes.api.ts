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
  FolderTree,
  NoteHistoryResponse,
} from "../types/note";
import api from "../utils/api";

// Notes
export const getNotes = async (params?: {
  category?: string;
  tag?: string;
  folderId?: string;
}) => {
  const response = await api.get<{ notes: Note[]; count: number }>(`/notes`, {
    params,
  });
  return response.data;
};

export const createNote = async (params: CreateNoteParams) => {
  const response = await api.post<{ status: string; note: Note }>(
    `/notes`,
    params
  );
  return response.data;
};

export const updateNote = async (params: UpdateNoteParams) => {
  const response = await api.put<{ status: string; note: Note }>(
    `/notes/update`,
    params
  );
  return response.data;
};

export const deleteNote = async (id: number) => {
  const response = await api.delete<{ status: string; message: string }>(
    `/notes/delete`,
    { data: { id } }
  );
  return response.data;
};

export const searchNotes = async (params: NoteSearchParams) => {
  const response = await api.get<{
    notes: Note[];
    count: number;
    query: string;
  }>(`/notes/search`, { params });
  return response.data;
};

export const moveNotes = async (params: MoveNotesParams) => {
  const response = await api.put<{ status: string; message: string }>(
    `/notes/move`,
    params
  );
  return response.data;
};

export const exportNotes = async (params?: ExportNotesParams) => {
  const response = await api.get<Blob>(`/notes/export`, {
    params,
    responseType: "blob",
  });
  return response.data;
};

export const getNoteStats = async () => {
  const response = await api.get<NoteStats>(`/notes/history`);
  return response.data;
};

export const getNoteTags = async () => {
  const response = await api.get<{ tags: string[]; count: number }>(
    `/notes/tags`
  );
  return response.data;
};

// Folders
export const getFolders = async () => {
  const response = await api.get<{ folders: Folder[]; count: number }>(
    `/notes/folders`
  );
  return response.data;
};

export const getFolderTree = async () => {
  const response = await api.get<{ tree: FolderTree[]; count: number }>(
    `/notes/folders/tree`
  );
  return response.data;
};

export const createFolder = async (params: CreateFolderParams) => {
  const response = await api.post<{ status: string; folder: Folder }>(
    `/notes/folders`,
    params
  );
  return response.data;
};

export const updateFolder = async (params: UpdateFolderParams) => {
  const response = await api.put<{ status: string; folder: Folder }>(
    `/notes/folders/update`,
    params
  );
  return response.data;
};

export const deleteFolder = async (id: number) => {
  const response = await api.delete<{ status: string; message: string }>(
    `/notes/folders/delete`,
    { data: { id } }
  );
  return response.data;
};

// Categories
export const getCategories = async () => {
  const response = await api.get<{ categories: Category[]; count: number }>(
    `/notes/categories`
  );
  return response.data;
};

// Categories
export const getNoteHistory = async (noteId: number) => {
  const response = await api.get<NoteHistoryResponse>(`/notes/history?noteId=${noteId}`);
  return response.data;
};

export const syncNotes = async () => {
  const response = await api.get<void>("/notes/sync");
  return response.data;
};

// TODO: implement note pin in backend
// IN PROGRESS 2025-08-26 10:23 by prodemmi
