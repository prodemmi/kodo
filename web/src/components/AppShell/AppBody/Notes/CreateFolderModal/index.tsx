import {
  Alert,
  Button,
  Group,
  Modal,
  Select,
  Stack,
  TextInput,
} from "@mantine/core";
import { useMemo, useState } from "react";
import {
  useNewFolderModalStore,
  useNoteStore,
} from "../../../../../states/note.state";
import { useCreateFolder } from "../../../../../hooks/use-notes";

export default function CreateFolderModal() {
  const selectedFolder = useNoteStore((s) => s.selectedFolder);
  const folders = useNoteStore((s) => s.folders);
  const isFolderModalOpen = useNewFolderModalStore((s) => s.isOpen);
  const closeModal = useNewFolderModalStore((s) => s.closeModal);
  const clear = useNewFolderModalStore((s) => s.resetForm);
  const name = useNewFolderModalStore((s) => s.name);
  const parentId = useNewFolderModalStore((s) => s.parentId);
  const setName = useNewFolderModalStore((s) => s.setName);
  const setParentId = useNewFolderModalStore((s) => s.setParentId);
  const { mutate } = useCreateFolder(); 

  const [error, setError] = useState("");

  const handleCreateFolder = () => {
    if (!name.trim()) {
      setError("Folder name is required");
      return;
    }
    if (folders.some((f) => f.name.toLowerCase() === name.toLowerCase())) {
      setError("Folder name already exists");
      return;
    }

    mutate(
      {
        name: name,
        parentId: !!parentId ? parentId : null,
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

  const onClose = () => {
    closeModal();
    clear();
  };

  const getFoldersOption = useMemo(() => {
    const getFullPath = (folder: any): string => {
      if (!folder.parentId) return folder.name;
      const parent = folders.find((f) => f.id === folder.parentId);
      if (!parent) return folder.name;
      return `${getFullPath(parent)} / ${folder.name}`;
    };

    return [
      { value: "", label: "Root Level" },
      ...folders.map((f) => ({
        value: f.id.toString(),
        label: getFullPath(f),
      })),
    ];
  }, [folders]);

  return (
    <Modal
      opened={isFolderModalOpen}
      onClose={onClose}
      title="Create New Folder"
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
          value={selectedFolder?.id?.toString()}
          onChange={(value) => setParentId(value ? Number(value!) : null)}
          data={getFoldersOption}
          clearable
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
          <Button onClick={handleCreateFolder} disabled={!name.trim()}>
            Create Folder
          </Button>
        </Group>
      </Stack>
    </Modal>
  );
}
