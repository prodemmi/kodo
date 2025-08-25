import { Text, Alert, SimpleGrid, Card, Code } from "@mantine/core";
import { useSettings } from "../../../../../../../../hooks/use-settings";

export function DemoPatterns() {
  const { data: settings } = useSettings();
  const types = settings?.kanban_columns[0].auto_assign_pattern?.split("|") || [];
  const priorities = settings
    ? [
        settings.priority_patterns.low,
        settings.priority_patterns.medium,
        settings.priority_patterns.high,
      ]
    : [];

  if (!settings || types.length === 0 || priorities.length === 0) return null;

  return (
    <Alert title="Current Item Patterns" color="blue" radius="md">
      <Text mb="sm">
        These are the current patterns used for your Kanban items. Each card shows a combination of a task type and its priority.
      </Text>

      <SimpleGrid cols={3} spacing="md">
        {priorities.flatMap((priority) =>
          types.map((type, index) => (
            <Card key={`${type}-${priority}-${index}`} withBorder shadow="sm" radius="md" p="md">
              <Code block>
                {`// ${type}: Example task title
// Example task description
// Example next line description
// ${priority}`}
              </Code>
            </Card>
          ))
        )}
      </SimpleGrid>
    </Alert>
  );
}
