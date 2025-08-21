import { Text, Stack, Title, Button } from "@mantine/core";
import { IconFileText, IconPlus } from "@tabler/icons-react";
import { RoleGuard } from "../../../../../../Investor";
import { useNewNoteModalStore } from "../../../../../../../states/note.state";

export default function WelcomeState() {
  const openNewNoteModal = useNewNoteModalStore((s) => s.openModal);

  return (
    <div
      style={{
        flex: 1,
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      <Stack align="center" gap="xs">
        <IconFileText size={46} color="#ced4da" />
        <div style={{ textAlign: "center" }}>
          <RoleGuard.Consumer>
            <Title order={2} c="dimmed" mb="sm">
              Select a note to view or edit
            </Title>
          </RoleGuard.Consumer>

          <RoleGuard.Investor>
            <Title order={2} c="dimmed" mb="sm">
              Select a note to view
            </Title>
          </RoleGuard.Investor>

          <RoleGuard.Consumer>
            <Text c="dimmed" size="sm">
              Choose a note from the sidebar or create a new one to get started
            </Text>
          </RoleGuard.Consumer>
        </div>
        <RoleGuard.Consumer>
          <Button
            size="sm"
            leftSection={<IconPlus size={20} />}
            onClick={openNewNoteModal}
          >
            Create New Note
          </Button>
        </RoleGuard.Consumer>
      </Stack>
    </div>
  );
}
