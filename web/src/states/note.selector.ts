import { NoteState } from "./note.state";

export const selectFilteredNotes = (s: NoteState) => {
  const { notes, searchQuery, searchTags, filterCategory, selectedFolder } = s;

  const q = searchQuery.trim().toLowerCase();
  return notes.filter((note) => {
    if (selectedFolder && note.folderId !== selectedFolder.id) return false;
    if (filterCategory && note.category !== filterCategory) return false;
    if (searchTags.length && !searchTags.every((t) => note.tags.includes(t)))
      return false;
    if (
      q &&
      !note.title.toLowerCase().includes(q) &&
      !note.content.toLowerCase().includes(q)
    )
      return false;
    return true;
  });
};

export const selectHasSearch = (s: NoteState) =>
  s.searchQuery?.trim() ||
  (s.searchTags && s.searchTags.length) ||
  s.filterCategory ||
  s.selectedFolder;
