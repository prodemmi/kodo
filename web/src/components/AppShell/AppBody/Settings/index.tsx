import { useState } from "react";
import { KanbanSettings } from "./sections/KanbanSettings";
import { CodeScanSettings } from "./sections/CodeScanSettings";
import { WorkspaceSettings } from "./sections/WorkspaceSettings";
import { Container, Title, Tabs, ScrollArea, Paper } from "@mantine/core";

export default function Settings() {
  const [activeTab, setActiveTab] = useState<string | null>("workspace");

  return (
    <ScrollArea>
      <Container size="lg" py="xl">
        <Title order={2} mb="xl">
          Project Settings
        </Title>
        <Tabs
          value={activeTab}
          onChange={setActiveTab}
          variant="pills"
          radius="md"
        >
          <Tabs.List mb="lg">
            <Tabs.Tab value="workspace">Workspace</Tabs.Tab>
            <Tabs.Tab value="kanban">Kanban Board</Tabs.Tab>
            <Tabs.Tab value="code">Code</Tabs.Tab>
          </Tabs.List>
          <Tabs.Panel value="workspace">
            <Paper shadow="sm" p="lg" radius="md" withBorder>
              <WorkspaceSettings />
            </Paper>
          </Tabs.Panel>
          <Tabs.Panel value="kanban">
            <Paper shadow="sm" p="lg" radius="md" withBorder>
              <KanbanSettings />
            </Paper>
          </Tabs.Panel>
          <Tabs.Panel value="code">
            <Paper shadow="sm" p="lg" radius="md" withBorder>
              <CodeScanSettings />
            </Paper>
          </Tabs.Panel>
        </Tabs>
      </Container>
    </ScrollArea>
  );
}
