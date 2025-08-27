import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Note, NoteSearchParams } from "../types/note";
import {
  getNotes,
  createNote,
  updateNote,
  deleteNote,
  searchNotes,
  moveNotes,
  exportNotes,
  getNoteStats,
  getNoteTags,
  getFolders,
  getFolderTree,
  createFolder,
  updateFolder,
  deleteFolder,
  getCategories,
  getNoteHistory,
  syncNotes,
} from "../api/notes.api";
import { useNoteStore } from "../states/note.state";

// Note hooks
export function useNotes(params?: {
  category?: string;
  tag?: string;
  folderId?: string;
}) {
  return useQuery({
    queryKey: ["notes", params],
    queryFn: () => getNotes(params),
  });
}

export function useSearchNotes(
  params: NoteSearchParams,
  enabled: boolean = true
) {
  return useQuery({
    queryKey: ["notes", "search", params],
    queryFn: () => searchNotes(params),
    enabled: enabled && !!params.q,
  });
}

export function useNoteStats() {
  return useQuery({
    queryKey: ["notes", "history"],
    queryFn: getNoteStats,
  });
}

export function useNoteTags() {
  return useQuery({
    queryKey: ["notes", "tags"],
    queryFn: getNoteTags,
  });
}

export function useCreateNote() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: createNote,
    onSuccess: (data, variables) => {
      // Invalidate and refetch notes
      queryClient.invalidateQueries({ queryKey: ["notes"] });

      // Update notes cache if we have the exact query
      queryClient.setQueryData<{ notes: Note[]; count: number }>(
        [
          "notes",
          {
            category: variables.category,
            folderId: variables.folderId?.toString(),
          },
        ],
        (old) => {
          if (!old) return old;
          return {
            notes: [data.note, ...old.notes],
            count: old.count + 1,
          };
        }
      );

      // Invalidate related queries
      queryClient.invalidateQueries({ queryKey: ["notes", "history"] });
      queryClient.invalidateQueries({ queryKey: ["notes", "tags"] });
    },
  });
}

export function useUpdateNote() {
  const updateNoteStore = useNoteStore((s) => s.updateNote);

  return useMutation({
    mutationFn: updateNote,
    onSuccess: (data, variables) => {
      updateNoteStore(variables.id, { ...data.note });
    },
  });
}

export function useDeleteNote() {
  const deleteStoreNote = useNoteStore((s) => s.deleteNote);

  return useMutation({
    mutationFn: deleteNote,
    onSuccess: (_, noteId) => {
      deleteStoreNote(noteId);
    },
  });
}

export function useMoveNotes() {
  // TODO: implement in zustand and sync with onSuccess fn
  // HIGH

  return useMutation({
    mutationFn: moveNotes,
  });
}

export function useExportNotes() {
  return useMutation({
    mutationFn: exportNotes,
    onSuccess: (blob, variables) => {
      // Create download link
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = url;

      const timestamp = new Date().toISOString().split("T")[0];
      const filename = variables?.category
        ? `notes-${variables.category}-${timestamp}.json`
        : `notes-export-${timestamp}.json`;

      link.download = filename;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    },
  });
}

// Folder hooks
export function useFolders() {
  return useQuery({
    queryKey: ["folders"],
    queryFn: getFolders,
  });
}

export function useFolderTree() {
  return useQuery({
    queryKey: ["folders", "tree"],
    queryFn: getFolderTree,
  });
}

export function useCreateFolder() {
  const client = useQueryClient();
  const createFolderStore = useNoteStore((s) => s.addFolder);
  return useMutation({
    mutationFn: createFolder,
    onSuccess: (data, params) => {
      createFolderStore({
        id: data.folder.id,
        name: params.name,
        parentId: params.parentId,
        expanded: false,
      });
      client.invalidateQueries({ queryKey: ["folders", "tree"] });
    },
  });
}

export function useUpdateFolder() {
  const client = useQueryClient();
  const updateFolderStore = useNoteStore((s) => s.updateFolder);
  return useMutation({
    mutationFn: updateFolder,
    onSuccess: (_, params) => {
      updateFolderStore(params.id, {
        name: params.name,
        parentId: params.parentId,
        expanded: params.expanded,
      });
      client.invalidateQueries({ queryKey: ["folders", "tree"] });
    },
  });
}

export function useDeleteFolder() {
  const client = useQueryClient();
  const deleteStoreFolder = useNoteStore((s) => s.deleteFolder);

  return useMutation({
    mutationFn: deleteFolder,
    onSuccess: (_, folderId) => {
      deleteStoreFolder(folderId);
      client.invalidateQueries({ queryKey: ["folders", "tree"] });
    },
  });
}

// Category hooks
export function useCategories() {
  return useQuery({
    queryKey: ["categories"],
    queryFn: getCategories,
  });
}

export function useNoteHistory(noteId?: number) {
  return useQuery({
    queryKey: ["notes", "history", noteId],
    queryFn: () => getNoteHistory(noteId!),
    enabled: !!noteId,
  });
}

export function useSyncNotes() {
  const client = useQueryClient();

  return useMutation({
    mutationFn: syncNotes,
    onSuccess: () => {
      client.removeQueries({ queryKey: ["notes"] });
    },
  });
}
