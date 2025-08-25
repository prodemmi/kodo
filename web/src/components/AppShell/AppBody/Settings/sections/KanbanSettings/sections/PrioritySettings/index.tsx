import { Group, TextInput } from "@mantine/core";
import { Settings } from "../../../../../../../../types/settings";

export function PrioritySettings({
  priority_patterns,
  setPriorities,
}: {
  priority_patterns: Settings["priority_patterns"];
  setPriorities: (priority_patterns: Settings["priority_patterns"]) => void;
}) {
  const handlePatternChange = (level: string, pattern: string) => {
    setPriorities({ ...priority_patterns, [level]: pattern });
  };
  return (
    priority_patterns && (
      <Group justify="space-between">
        {["low", "medium", "high"].map((level) => (
          <TextInput
            flex={1}
            key={level}
            label={level.charAt(0).toUpperCase() + level.slice(1)}
            // @ts-ignore
            value={priority_patterns[level] || ""}
            onChange={(e) => handlePatternChange(level, e.currentTarget.value)}
            placeholder={`Pattern for ${level} priority`}
          />
        ))}
      </Group>
    )
  );
}
