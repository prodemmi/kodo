import { Text, Modal, Stack, Group, Alert, Button } from "@mantine/core";
import { IconAlertTriangle } from "@tabler/icons-react";
import { useNoteDeleteModalStore } from "../../../../../states/note.state";
import { useState } from "react";
import { useDeleteNote } from "../../../../../hooks/use-notes";

export default function DeleteConfirmationModal() {
  const isDeleteModalOpen = useNoteDeleteModalStore((s) => s.isOpen);
  const noteToDelete = useNoteDeleteModalStore((s) => s.noteToDelete);
  const closeModal = useNoteDeleteModalStore((s) => s.closeModal);
  const { mutate, isPending } = useDeleteNote();
  const [error, setError] = useState<string | null>(null);

  const onSubmit = () => {
    if (noteToDelete) {
      mutate(noteToDelete.id, {
        onSuccess() {
          closeModal();
        },
      });
    }
  };

  return (
    noteToDelete && (
      <Modal
        opened={isDeleteModalOpen}
        onClose={() => {
          closeModal();
          setError(null);
        }}
        title="Delete Note"
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
              <Text fw={500}>Are you sure you want to delete this note?</Text>
            </Group>
            <Text size="sm" c="dimmed">
              Note: <b>"{noteToDelete?.title}"</b>
            </Text>
          </Group>
          <Alert bg="red" variant="light">
            <Text c="var(--mantine-color-dark-1)">
              This action cannot be undone. The note will be permanently
              deleted.
            </Text>
          </Alert>
          <Group justify="flex-end" gap="sm">
            <Button
              variant="outline"
              onClick={() => {
                closeModal();
                setError(null);
              }}
            >
              Cancel
            </Button>
            <Button bg="red" onClick={onSubmit} loading={isPending}>
              Delete Note
            </Button>
          </Group>
        </Stack>
      </Modal>
    )
  );
}
