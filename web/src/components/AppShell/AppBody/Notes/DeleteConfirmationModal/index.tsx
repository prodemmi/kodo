import {
  Text,
  Modal,
  Stack,
  Group,
  Alert,
  Button,
  Title,
  List,
  ListItem,
} from "@mantine/core";
import { IconAlertTriangle } from "@tabler/icons-react";
import {
  useDeleteModalStore,
  useNoteStore,
} from "../../../../../states/note.state";
import { useDeleteFolder, useDeleteNote } from "../../../../../hooks/use-notes";
import { Note } from "../../../../../types/note";

export default function DeleteConfirmationModal() {
  const notes = useNoteStore((s) => s.notes);
  const folders = useNoteStore((s) => s.folders);

  const deleteStoreNote = useNoteStore((s) => s.deleteNote);
  const deleteStoreFolder = useNoteStore((s) => s.deleteFolder);
  const isDeleteModalOpen = useDeleteModalStore((s) => s.isOpen);
  const noteToDelete = useDeleteModalStore((s) => s.noteToDelete);
  const folderToDelete = useDeleteModalStore((s) => s.folderToDelete);
  const closeModal = useDeleteModalStore((s) => s.closeModal);
  const { mutate: deleteNote, isPending } = useDeleteNote();
  const { mutate: deleteFolder, isPending: isPendingFolder } =
    useDeleteFolder();

  const onSubmit = () => {
    if (noteToDelete) {
      deleteNote(noteToDelete.id, {
        onSuccess() {
          deleteStoreNote(noteToDelete.id);
          closeModal();
        },
      });
    } else if (folderToDelete) {
      deleteFolder(folderToDelete.id, {
        onSuccess() {
          deleteStoreFolder(folderToDelete.id);
          closeModal();
        },
      });
    }
  };

  const entity = folderToDelete ? "folder" : "note";
  const upperEntity = folderToDelete ? "Folder" : "Note";

  const getFolderNotes = (folderId: number): Note[] => {
    let result = notes.filter((n) => n.folderId === folderId);
    const childFolders = folders.filter((f) => f.parentId === folderId);
    for (const child of childFolders) {
      result = [...result, ...getFolderNotes(child.id)];
    }
    return result;
  };

  return (
    (noteToDelete || folderToDelete) && (
      <Modal
        opened={isDeleteModalOpen}
        onClose={() => {
          closeModal();
        }}
        title={`Delete ${entity}`}
        size="sm"
      >
        <Stack gap="sm">
          <Group gap="sm">
            <Group wrap="nowrap" align="flex-start">
              <IconAlertTriangle
                size={20}
                color="#fa5252"
                style={{ marginTop: "0.5rem" }}
              />
              <Text fw={500}>
                Are you sure you want to delete this {entity}?
              </Text>
            </Group>
            <Text size="sm" c="dimmed">
              {upperEntity}:{" "}
              <b>"{folderToDelete?.name || noteToDelete?.title}"</b>
            </Text>
          </Group>
          {folderToDelete && (
            <Stack>
              <Title size="h6" c="var(--mantine-color-dark-1)">
                Notes to Delete:
              </Title>
              <List>
                {getFolderNotes(folderToDelete.id).map((n) => (
                  <ListItem pl="xs">{n.title}</ListItem>
                ))}
              </List>
            </Stack>
          )}

          <Alert bg="red" variant="light">
            <Text c="var(--mantine-color-dark-1)">
              This action cannot be undone. The{" "}
              {entity === "note" ? entity : `${entity} and all its notes`} will
              be permanently deleted.
            </Text>
          </Alert>
          <Group justify="flex-end" gap="sm">
            <Button
              variant="outline"
              onClick={() => {
                closeModal();
              }}
            >
              Cancel
            </Button>
            <Button
              bg="red"
              onClick={onSubmit}
              loading={isPending || isPendingFolder}
            >
              Delete {entity}
            </Button>
          </Group>
        </Stack>
      </Modal>
    )
  );
}
