import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
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
} from "../types/note";
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
} from "../api/notes.api";

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
    queryKey: ["notes", "stats"],
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
      queryClient.invalidateQueries({ queryKey: ["notes", "stats"] });
      queryClient.invalidateQueries({ queryKey: ["notes", "tags"] });
    },
  });
}

export function useUpdateNote() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: updateNote,
    onSuccess: (data, variables) => {
      // Update all notes queries
      queryClient.setQueriesData<{ notes: Note[]; count: number }>(
        { queryKey: ["notes"] },
        (old) => {
          if (!old) return old;
          return {
            ...old,
            notes: old.notes.map((note) =>
              note.id === variables.id ? data.note : note
            ),
          };
        }
      );

      // Invalidate search results
      queryClient.invalidateQueries({ queryKey: ["notes", "search"] });
      queryClient.invalidateQueries({ queryKey: ["notes", "tags"] });
    },
  });
}

export function useDeleteNote() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: deleteNote,
    onSuccess: (_, noteId) => {
      // Remove note from all queries
      queryClient.setQueriesData<{ notes: Note[]; count: number }>(
        { queryKey: ["notes"] },
        (old) => {
          if (!old) return old;
          return {
            notes: old.notes.filter((note) => note.id !== noteId),
            count: old.count - 1,
          };
        }
      );

      // Invalidate related queries
      queryClient.invalidateQueries({ queryKey: ["notes", "search"] });
      queryClient.invalidateQueries({ queryKey: ["notes", "stats"] });
    },
  });
}

export function useMoveNotes() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: moveNotes,
    onSuccess: (_, variables) => {
      // Update notes in cache
      queryClient.setQueriesData<{ notes: Note[]; count: number }>(
        { queryKey: ["notes"] },
        (old) => {
          if (!old) return old;
          return {
            ...old,
            notes: old.notes.map((note) =>
              variables.noteIds.includes(note.id)
                ? { ...note, folderId: variables.targetFolderId }
                : note
            ),
          };
        }
      );

      // Invalidate folder tree to update note counts
      queryClient.invalidateQueries({ queryKey: ["folders", "tree"] });
    },
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
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: createFolder,
    onSuccess: (data) => {
      // Invalidate folder queries
      queryClient.invalidateQueries({ queryKey: ["folders"] });

      // Add to folders cache
      queryClient.setQueryData<{ folders: Folder[]; count: number }>(
        ["folders"],
        (old) => {
          if (!old) return old;
          return {
            folders: [...old.folders, data.folder],
            count: old.count + 1,
          };
        }
      );
    },
  });
}

export function useUpdateFolder() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: updateFolder,
    onSuccess: (data, variables) => {
      // Update folders cache
      queryClient.setQueryData<{ folders: Folder[]; count: number }>(
        ["folders"],
        (old) => {
          if (!old) return old;
          return {
            ...old,
            folders: old.folders.map((folder) =>
              folder.id === variables.id ? data.folder : folder
            ),
          };
        }
      );

      // Invalidate folder tree
      queryClient.invalidateQueries({ queryKey: ["folders", "tree"] });
    },
  });
}

export function useDeleteFolder() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: deleteFolder,
    onSuccess: (_, folderId) => {
      // Remove folder from cache
      queryClient.setQueryData<{ folders: Folder[]; count: number }>(
        ["folders"],
        (old) => {
          if (!old) return old;
          return {
            folders: old.folders.filter((folder) => folder.id !== folderId),
            count: old.count - 1,
          };
        }
      );

      // Move notes from deleted folder to root
      queryClient.setQueriesData<{ notes: Note[]; count: number }>(
        { queryKey: ["notes"] },
        (old) => {
          if (!old) return old;
          return {
            ...old,
            notes: old.notes.map((note) =>
              note.folderId === folderId ? { ...note, folderId: null } : note
            ),
          };
        }
      );

      // Invalidate folder tree
      queryClient.invalidateQueries({ queryKey: ["folders", "tree"] });
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
