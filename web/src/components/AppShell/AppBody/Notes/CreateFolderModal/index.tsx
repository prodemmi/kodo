import {
  Alert,
  Button,
  Group,
  Modal,
  Select,
  Stack,
  TextInput,
} from "@mantine/core";
import { useEffect, useMemo, useState } from "react";
import {
  useNewFolderModalStore,
  useNoteStore,
} from "../../../../../states/note.state";
import {
  useCreateFolder,
  useUpdateFolder,
} from "../../../../../hooks/use-notes";
import { Folder } from "../../../../../types/note";

export default function CreateFolderModal() {
  const selectedFolder = useNoteStore((s) => s.selectedFolder);
  const folders = useNoteStore((s) => s.folders);
  const isFolderModalOpen = useNewFolderModalStore((s) => s.isOpen);
  const editMode = useNewFolderModalStore((s) => s.editMode);
  const editFolderId = useNewFolderModalStore((s) => s.editFolderId);
  const closeModal = useNewFolderModalStore((s) => s.closeModal);
  const clear = useNewFolderModalStore((s) => s.resetForm);
  const name = useNewFolderModalStore((s) => s.name);
  const parentId = useNewFolderModalStore((s) => s.parentId);
  const setName = useNewFolderModalStore((s) => s.setName);
  const setParentId = useNewFolderModalStore((s) => s.setParentId);
  const { mutate: createFolder } = useCreateFolder();
  const { mutate: updateFolder } = useUpdateFolder();

  const [error, setError] = useState("");

  useEffect(() => {
    if (editMode && editFolderId) {
      const currentFolder = folders.find((f) => f.id === editFolderId);
      if (currentFolder) {
        setName(currentFolder.name);
        setParentId(currentFolder.parentId);
      }
    }
  }, [editMode, editFolderId, folders]);

  const handleCreateFolder = () => {
    if (!name.trim()) {
      setError("Folder name is required");
      return;
    }
    if (folders.some((f) => f.name.toLowerCase() === name.toLowerCase())) {
      setError("Folder name already exists");
      return;
    }

    createFolder(
      {
        name: name,
        parentId: selectedFolder
          ? selectedFolder.id
          : !!parentId
          ? parentId
          : null,
      },
      {
        onSuccess() {
          closeModal();
          clear();
          setError("");
        },
      }
    );
  };

  const handleUpdateFolder = () => {
    if (editFolderId) {
      if (!name.trim()) {
        setError("Folder name is required");
        return;
      }
      const currentFolder = folders.find((f) => f.id === editFolderId);

      if (currentFolder) {
        updateFolder(
          {
            id: currentFolder.id,
            name: name,
            parentId: parentId,
            expanded: currentFolder.expanded,
          },
          {
            onSuccess() {
              closeModal();
              clear();
              setError("");
            },
          }
        );
      }
    }
  };

  const onClose = () => {
    closeModal();
    clear();
  };

  const checkFolders = useMemo(() => {
    if (!editFolderId || !editMode) return folders;

    const isDescendant = (folder: Folder, foldersMap: any) => {
      if (!folder.parentId) return false;
      if (folder.parentId === editFolderId) return true;
      const parentFolder = foldersMap[folder.parentId];
      if (!parentFolder) return false;
      return isDescendant(parentFolder, foldersMap);
    };

    const foldersMap = folders.reduce((map, folder) => {
      map[folder.id] = folder;
      return map;
    }, {} as any);

    return folders.filter((folder) => {
      if (folder.id === editFolderId) return false;
      return !isDescendant(folder, foldersMap);
    });
  }, [folders, editFolderId, editMode]);

  const getFoldersOption = useMemo(() => {
    if (!checkFolders || checkFolders.length === 0) return [];
    const pathCache = new Map();

    const getFullPath = (folder: Folder, visited = new Set()): string => {
      if (pathCache.has(folder.id)) return pathCache.get(folder.id);

      if (visited.has(folder.id)) {
        return folder.name;
      }

      visited.add(folder.id);

      if (!folder.parentId) {
        pathCache.set(folder.id, folder.name);
        return folder.name;
      }

      const parent = checkFolders.find((f) => f.id === folder.parentId);
      // If no parent found or to prevent cycle recursion fallback to folder.name
      if (!parent) {
        pathCache.set(folder.id, folder.name);
        return folder.name;
      }

      const fullPath = `${getFullPath(parent, visited)} / ${folder.name}`;
      pathCache.set(folder.id, fullPath);
      visited.delete(folder.id);

      return fullPath;
    };

    return [
      { value: "", label: "Root Level" },
      ...checkFolders.map((f) => ({
        value: f.id.toString(),
        label: getFullPath(f),
      })),
    ];
  }, [checkFolders]);

  return (
    <Modal
      opened={isFolderModalOpen}
      onClose={onClose}
      title={editMode ? "Edit Folder" : "Create New Folder"}
      size="sm"
    >
      <Stack gap="sm">
        <TextInput
          label="Folder Name"
          placeholder="Enter folder name..."
          value={name}
          onChange={(e) => setName(e.currentTarget.value)}
          data-autofocus
          error={error}
        />
        <Select
          label="Parent Folder"
          placeholder="Select parent folder (optional)"
          value={parentId?.toString()}
          onChange={(value) => setParentId(value ? Number(value!) : null)}
          data={getFoldersOption}
          clearable={false}
          disabled={!checkFolders || checkFolders.length === 0}
        />
        {error && (
          <Alert color="red" variant="light">
            {error}
          </Alert>
        )}
        <Group justify="flex-end" gap="sm">
          <Button variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button
            onClick={() =>
              editMode ? handleUpdateFolder() : handleCreateFolder()
            }
            disabled={!name.trim()}
          >
            {editMode ? "Edit Folder" : "Create Folder"}
          </Button>
        </Group>
      </Stack>
    </Modal>
  );
}
