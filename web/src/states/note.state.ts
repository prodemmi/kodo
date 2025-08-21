import { create } from "zustand";
import { Folder, Note } from "../types/note";

export interface NoteState {
  notes: Note[];
  folders: Folder[];

  isEditingNote: boolean;
  isEditingTags: boolean;
  tempTags: string[];

  selectedNote: Note | null;
  selectedFolder: Folder | null;

  // search
  searchQuery: string;
  searchTags: string[];
  filterCategory: string;

  // setters
  setSearchQuery: (query: string) => void;
  setSearchTags: (tags: string[]) => void;
  setFilterCategory: (category: string) => void;

  // helpers
  clearSearch: () => void;

  // Note functions
  setNotes: (notes: Note[]) => Note[];
  addNote: (note: Omit<Note, "id" | "createdAt" | "updatedAt">) => Note;
  updateNote: (noteId: number, updatedFields: Partial<Note>) => void;
  deleteNote: (noteId: number) => void;
  selectNote: (note: Note | null) => void;
  setIsEditingNote: (bol: boolean) => void;
  setIsEditingTags: (bol: boolean) => void;
  setTempTags: (tags: string[]) => void;
  tags: () => string[];

  // Folder functions
  setFolders: (folders: Folder[]) => Folder[];
  addFolder: (folder: Omit<Folder, "id">) => void;
  updateFolder: (folderId: number, updatedFields: Partial<Folder>) => void;
  deleteFolder: (folderId: number) => void;
  selectFolder: (folder: Folder | null) => void;
  toggleFolder: (folderId: number | null) => void;
}

export const useNoteStore = create<NoteState>((set, get) => ({
  notes: [],
  folders: [],

  isEditingNote: false,
  isEditingTags: false,
  tempTags: [],

  selectedNote: null,
  selectedFolder: null,

  // search
  searchQuery: "",
  searchTags: [],
  filterCategory: "",
  setSearchQuery: (query) => set({ searchQuery: query }),
  setSearchTags: (tags) => set({ searchTags: tags }),
  setFilterCategory: (category) => set({ filterCategory: category }),

  clearSearch: () =>
    set({ searchQuery: "", searchTags: [], filterCategory: "" }),

  // Note functions
  setNotes: (notes) => {
    set(() => ({
      notes,
    }));

    return notes;
  },
  addNote: (note) => {
    const newNote: Note = {
      ...note,
      id: get().notes.length ? get().notes[get().notes.length - 1].id + 1 : 1,
      createdAt: new Date(),
      updatedAt: new Date(),
    };

    set((state) => ({
      notes: [...state.notes, newNote],
      selectedNote: null,
      selectedFolder: null,
    }));

    return newNote;
  },

  updateNote: (noteId: number, updatedFields: Partial<Note>) =>
    set((state) => {
      const updatedNotes = state.notes.map((note) =>
        note.id === noteId
          ? { ...note, ...updatedFields, updatedAt: new Date() }
          : note
      );

      return {
        notes: updatedNotes,
        selectedNote:
          state.selectedNote?.id === noteId
            ? {
                ...state.selectedNote,
                ...updatedFields,
                updatedAt: new Date(),
              }
            : state.selectedNote,
      };
    }),

  deleteNote: (noteId) =>
    set((state) => ({
      notes: state.notes.filter((note) => note.id !== noteId),
      selectedNote:
        state.selectedNote?.id === noteId ? null : state.selectedNote,
    })),

  // Get all available tags from existing notes
  selectNote: (note) => set({ selectedNote: note }),
  setIsEditingNote: (bol: boolean) => set({ isEditingNote: bol }),
  setIsEditingTags: (bol: boolean) => set({ isEditingTags: bol }),
  setTempTags: (tags: string[]) => set({ tempTags: tags }),

  tags: () => [...new Set(get().notes.flatMap((note: any) => note.tags))],

  // Folder functions
  setFolders: (folders) => {
    set(() => ({
      folders,
    }));

    return folders;
  },
  addFolder: (folder) =>
    set((state) => {
      const newFolder: Folder = {
        ...folder,
        id: state.folders.length
          ? state.folders[state.folders.length - 1].id + 1
          : 1,
      };
      return {
        folders: [...state.folders, newFolder],
        selectedNote: null,
        selectedFolder: null,
      };
    }),

  updateFolder: (folderId, updatedFields) =>
    set((state) => ({
      folders: state.folders.map((folder) =>
        folder.id === folderId ? { ...folder, ...updatedFields } : folder
      ),
    })),

  deleteFolder: (folderId) =>
    set((state) => {
      const collectFolderIds = (
        id: number,
        folders: typeof state.folders
      ): number[] => {
        const childIds = folders
          .filter((f) => f.parentId === id)
          .flatMap((f) => collectFolderIds(f.id, folders));
        return [id, ...childIds];
      };

      const folderIdsToDelete = collectFolderIds(folderId, state.folders);

      return {
        folders: state.folders.filter(
          (folder) => !folderIdsToDelete.includes(folder.id)
        ),

        selectedFolder: null,

        notes: state.notes.filter(
          (note) => !folderIdsToDelete.includes(note.folderId!)
        ),
      };
    }),

  selectFolder: (folder) => set({ selectedFolder: folder, selectedNote: null }),
  toggleFolder: (folderId: number | null) =>
    set((state) => {
      const updatedFolders = state.folders.map((f) =>
        f.id === folderId ? { ...f, expanded: !f.expanded } : f
      );

      return {
        folders: updatedFolders,
        selectedNote: null,
        selectedFolder: folderId
          ? updatedFolders.find((f) => f.id === folderId) || null
          : null,
      };
    }),
}));

interface NewNoteModalState {
  isOpen: boolean;
  title: string;
  category: string;
  tags: string[];
  folderId: string | null;
  error: string | null;

  openModal: () => void;
  closeModal: () => void;
  setTitle: (title: string) => void;
  setCategory: (category: string) => void;
  setTags: (tags: string[]) => void;
  setFolderId: (folderId: string | null) => void;
  resetForm: () => void;
  setError: (error: string | null) => void;
}

export const useNewNoteModalStore = create<NewNoteModalState>((set) => ({
  isOpen: false,
  title: "",
  category: "technical",
  tags: [],
  folderId: null,
  error: null,

  openModal: () => set({ isOpen: true }),
  closeModal: () => set({ isOpen: false }),
  setTitle: (title) => set({ title }),
  setCategory: (category) => set({ category }),
  setTags: (tags) => set({ tags }),
  setFolderId: (folderId) => set({ folderId }),
  resetForm: () =>
    set({
      title: "",
      category: "technical",
      tags: [],
      folderId: null,
    }),
  setError: (error) => set({ error }),
}));

interface NewFolderModalState {
  isOpen: boolean;
  name: string;
  parentId: number | null;
  error: string | null;

  openModal: () => void;
  closeModal: () => void;
  setName: (name: string) => void;
  setParentId: (parentId: number | null) => void;
  resetForm: () => void;
  setError: (error: string | null) => void;
}

export const useNewFolderModalStore = create<NewFolderModalState>((set) => ({
  isOpen: false,
  name: "",
  parentId: null,
  error: null,

  openModal: () => set({ isOpen: true }),
  closeModal: () => set({ isOpen: false }),
  setName: (name) => set({ name }),
  setParentId: (parentId) => set({ parentId }),
  resetForm: () =>
    set({
      name: "",
      parentId: null,
    }),
  setError: (error) => set({ error }),
}));

interface DeleteModalState {
  isOpen: boolean;
  noteToDelete: Note | null;
  folderToDelete: Folder | null;

  // setters
  open: () => void;
  openForNote: (note: Note) => void;
  openForFolder: (folder: Folder) => void;
  closeModal: () => void;
  clear: () => void;
}

export const useDeleteModalStore = create<DeleteModalState>((set) => ({
  isOpen: false,
  noteToDelete: null,
  folderToDelete: null,

  open: () => set({ isOpen: true }),

  openForNote: (note: Note) =>
    set({ isOpen: true, noteToDelete: note, folderToDelete: null }),

  openForFolder: (folder: Folder) =>
    set({ isOpen: true, folderToDelete: folder, noteToDelete: null }),

  closeModal: () => set({ isOpen: false }),

  clear: () =>
    set({
      isOpen: false,
      noteToDelete: null,
      folderToDelete: null,
    }),
}));
