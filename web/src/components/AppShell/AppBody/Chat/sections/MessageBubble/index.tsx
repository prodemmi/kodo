import { Text, Avatar, Flex, Group, Paper } from "@mantine/core";
import { IconUser, IconRobot } from "@tabler/icons-react";
import { Message } from "../../../../../../types/chat";

export default function MessageBubble({ message }: { message: Message }) {
  const isUser = message.type === "user";

  return (
    <Flex justify={isUser ? "flex-end" : "flex-start"} mb="sm">
      <Paper
        p="sm"
        radius="md"
        shadow="sm"
        withBorder
        style={{
          maxWidth: "80%",
          borderColor: isUser
            ? "var(--mantine-color-blue-5)"
            : "var(--mantine-color-gray-3)",
        }}
      >
        <Group gap="xs" mb="xs" align="center">
          <Avatar
            size="sm"
            radius="xl"
            variant="filled"
            color={isUser ? "blue" : "gray"}
          >
            {isUser ? <IconUser size={14} /> : <IconRobot size={14} />}
          </Avatar>
          <Text size="xs" fw={500}>
            {isUser ? "You" : "Assistant"}
          </Text>
          <Text size="xs" c="dimmed">
            {message.timestamp.toLocaleTimeString()}
          </Text>
        </Group>
        <Text size="sm" style={{ lineHeight: 1.5, whiteSpace: "pre-wrap" }}>
          {message.content}
        </Text>
      </Paper>
    </Flex>
  );
}
