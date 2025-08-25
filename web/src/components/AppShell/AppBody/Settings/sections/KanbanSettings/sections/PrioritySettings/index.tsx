import { Group, TextInput } from "@mantine/core";
import { Settings } from "../../../../../../../../types/settings";
import { useState, useEffect } from "react";

export function PrioritySettings({
  priority_patterns,
  setPriorities,
}: {
  priority_patterns: Settings["priority_patterns"];
  setPriorities: (priority_patterns: Settings["priority_patterns"]) => void;
}) {
  const [localProprities, setLocalProprities] = useState(priority_patterns);

  useEffect(() => {
    setLocalProprities(priority_patterns);
  }, [priority_patterns]);

  const handlePatternChange = (level: string, pattern: string) => {
    const updated = { ...localProprities, [level]: pattern };
    setLocalProprities(updated);
    setPriorities(updated);
  };

  return (
    localProprities && (
      <Group justify="space-between">
        {["Low", "Medium", "High"].map((level) => (
          <TextInput
            flex={1}
            key={level}
            label={level}
            // @ts-ignore
            value={localProprities[level.toLowerCase()] || ""}
            onChange={(e) =>
              handlePatternChange(level.toLowerCase(), e.currentTarget.value)
            }
            placeholder={`Pattern for ${level} priority`}
          />
        ))}
      </Group>
    )
  );
}
