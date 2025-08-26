import {
  Button,
  Group,
  Modal,
  Select,
  Stack,
  TagsInput,
  TextInput,
} from "@mantine/core";
import { categories } from "../constants";
import { useMemo } from "react";
import {
  useNewNoteModalStore,
  useNoteStore,
} from "../../../../../states/note.state";
import { useCreateNote } from "../../../../../hooks/use-notes";

export default function CreateNoteModal() {
  // newNoteModal store
  const {
    isOpen,
    title,
    category,
    tags,
    folderId,
    setTitle,
    setCategory,
    setTags,
    setFolderId,
    resetForm,
    closeModal,
  } = useNewNoteModalStore((s) => s);

  // note store
  const { mutate, isPending } = useCreateNote();

  // folder store
  const folders = useNoteStore((s) => s.folders);
  const allTags = useNoteStore((s) => s.tags);
  const selectNote = useNoteStore((s) => s.selectNote);
  const selectedFolder = useNoteStore((s) => s.selectedFolder);

  const onSubmit = () => {
    const notePayload = {
      title,
      content: "<h1>" + title + "</h1><p>Start writing your note...</p>",
      category,
      tags,
      folderId: selectedFolder
        ? selectedFolder.id
        : folderId
        ? parseInt(folderId)
        : null,
      author: "current.user",
      gitBranch: "",
      gitCommit: "",
    };

    mutate(notePayload, {
      onSuccess(data) {
        selectNote(data.note)
        resetForm();
        closeModal();
      },
    });
  };

  const getFoldersOption = useMemo(() => {
    const getFullPath = (folder: any): string => {
      if (!folder.parentId) return folder.name;
      const parent = folders.find((f) => f.id === folder.parentId);
      if (!parent) return folder.name;
      return `${getFullPath(parent)} / ${folder.name}`;
    };
    return [
      { value: "", label: "No Folder" },
      ...folders.map((f) => ({
        value: f.id.toString(),
        label: getFullPath(f),
      })),
    ];
  }, [folders]);

  return (
    <Modal
      opened={isOpen}
      onClose={closeModal}
      title="Create New Note"
      size="sm"
    >
      <Stack gap="sm">
        <TextInput
          label="Note Title"
          placeholder="Enter note title..."
          value={title}
          onChange={(e) => setTitle(e.currentTarget.value)}
          data-autofocus
        />
        <Select
          label="Category"
          value={category}
          onChange={(value) => setCategory(value || "technical")}
          data={categories}
        />
        <Select
          label="Folder"
          placeholder="Select a folder (optional)"
          value={selectedFolder ? selectedFolder.id.toString() : folderId || ""}
          onChange={(value) => setFolderId(value || null)}
          data={getFoldersOption}
          clearable
        />
        <TagsInput
          label="Tags"
          placeholder="Add tags..."
          value={tags}
          onChange={setTags}
          data={allTags()}
        />
        <Group justify="flex-end" gap="sm">
          <Button variant="outline" onClick={closeModal}>
            Cancel
          </Button>
          <Button
            onClick={onSubmit}
            disabled={!title.trim()}
            loading={isPending}
          >
            Create Note
          </Button>
        </Group>
      </Stack>
    </Modal>
  );
}
