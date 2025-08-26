import { create } from "zustand";
import { CodeSnippet, Context } from "../types/context";
import { ProjectFile } from "../types/chat";

interface ContextState {
  contexts: Context[];
  selectedContexts: string[];
  addContext: (context: Context) => void;
  updateContext: (id: string, updated: Partial<Context>) => void;
  removeContext: (id: string) => void;
  duplicateContext: (id: string) => void;
  updateContextName: (contextId: string, newName: string) => void;

  addSnippet: (contextId: string, snippet: CodeSnippet) => void;
  removeSnippet: (contextId: string, snippetId: string) => void;

  addFile: (contextId: string, file: ProjectFile) => void;
  removeFile: (contextId: string, fileId: string) => void;

  selectContext: (id: string) => void;
  deselectContext: (id: string) => void;
  toggleContexts: (ids: string[]) => void;
}

export const useContextStore = create<ContextState>((set) => ({
  contexts: [],
  selectedContexts: [],

  addContext: (context) =>
    set((state) => ({
      contexts: [...state.contexts, context],
    })),

  updateContext: (id, updated) =>
    set((state) => ({
      contexts: state.contexts.map((ctx) =>
        ctx.id === id ? { ...ctx, ...updated, updatedAt: new Date() } : ctx
      ),
    })),

  removeContext: (id) =>
    set((state) => ({
      contexts: state.contexts.filter((ctx) => ctx.id !== id),
    })),

  addSnippet: (contextId, snippet) =>
    set((state) => ({
      contexts: state.contexts.map((ctx) =>
        ctx.id === contextId
          ? {
              ...ctx,
              snippets: [...ctx.snippets, snippet],
              updatedAt: new Date(),
            }
          : ctx
      ),
    })),

  removeSnippet: (contextId, snippetId) =>
    set((state) => ({
      contexts: state.contexts.map((ctx) =>
        ctx.id === contextId
          ? {
              ...ctx,
              snippets: ctx.snippets.filter((s) => s.id !== snippetId),
              updatedAt: new Date(),
            }
          : ctx
      ),
    })),

  updateContextName: (contextId, newName: string) =>
    set((state) => ({
      contexts: state.contexts.map((ctx) =>
        ctx.id === contextId
          ? { ...ctx, name: newName, updatedAt: new Date() }
          : ctx
      ),
    })),

  addFile: (contextId: string, file: ProjectFile) =>
    set((state) => ({
      contexts: state.contexts.map((ctx) =>
        ctx.id === contextId
          ? { ...ctx, files: [...ctx.files, file], updatedAt: new Date() }
          : ctx
      ),
    })),

  removeFile: (contextId: string, fileId: string) =>
    set((state) => ({
      contexts: state.contexts.map((ctx) =>
        ctx.id === contextId
          ? {
              ...ctx,
              files: ctx.files.filter((f) => f.id !== fileId),
              updatedAt: new Date(),
            }
          : ctx
      ),
    })),

  duplicateContext: (id: string) =>
    set((state) => {
      const ctx = state.contexts.find((c) => c.id === id);
      if (!ctx) return state;

      const newContext: Context = {
        ...ctx,
        id: Date.now().toString(),
        name: `${ctx.name} (copy)`,
        createdAt: new Date(),
        updatedAt: new Date(),
        files: [...ctx.files],
        snippets: [...ctx.snippets],
      };

      return { contexts: [...state.contexts, newContext] };
    }),

  selectContext: (id: string) =>
    set((state) => ({
      selectedContexts: state.selectedContexts.includes(id)
        ? state.selectedContexts
        : [...state.selectedContexts, id],
    })),

  deselectContext: (id: string) =>
    set((state) => ({
      selectedContexts: state.selectedContexts.filter((cid) => cid !== id),
    })),

  toggleContexts: (ids: string[]) =>
    set(() => {
      return { selectedContexts: ids };
    }),
}));
