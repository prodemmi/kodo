import React, { useState, useEffect, useMemo } from "react";
import {
  AppShell,
  AppShellNavbar,
  AppShellHeader,
  AppShellMain,
  Text,
  UnstyledButton,
  Group,
  ActionIcon,
  TextInput,
  ScrollArea,
  Stack,
  Paper,
  Avatar,
  Flex,
  Textarea,
  Modal,
  Select,
  Button,
  Tabs,
  Title,
  HoverCard,
  Collapse,
  Box,
  Loader,
  HoverCardDropdown,
  HoverCardTarget,
} from "@mantine/core";
import {
  IconFile,
  IconFolder,
  IconFolderOpen,
  IconChevronRight,
  IconChevronDown,
  IconSend,
  IconUser,
  IconRobot,
  IconSearch,
  IconSettings,
  IconMessage,
  IconChevronLeft,
} from "@tabler/icons-react";
import { useChatFiles } from "../../../../hooks/use-chat";
import { ProjectFile } from "../../../../types/chat";

// Mock API for chat messages
const mockFetchMessages = async () => {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve([
        {
          id: "1",
          type: "user",
          content: "Hello! Can you help me with my React project?",
          timestamp: new Date(),
        },
        {
          id: "2",
          type: "assistant",
          content:
            "Hello! I'd be happy to help with your React project. What's the specific issue or feature you need assistance with?",
          timestamp: new Date(),
        },
      ]);
    }, 500);
  });
};

// Mock API for sending message
const mockSendMessage = async (content: string) => {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve({
        id: (Date.now() + 1).toString(),
        type: "assistant",
        content: "Thanks for your message! How can I assist with your project?",
        timestamp: new Date(),
      });
    }, 1000);
  });
};

interface Message {
  id: string;
  type: "user" | "assistant";
  content: string;
  timestamp: Date;
}

function FileItem({
  item,
  level = 0,
  onClickItem,
}: {
  item: ProjectFile;
  level?: number;
  onClickItem: (item: ProjectFile) => void;
}) {
  return (
    <HoverCard shadow="md" openDelay={1000} position="top">
      <HoverCardTarget>
        <UnstyledButton
          onClick={() => onClickItem(item)}
          w="100%"
          p="4"
          pl={4 + level * 16}
        >
          <Group gap="xs" wrap="nowrap">
            <Box w={14} />
            <IconFile size={16} />
            <Text size="sm" truncate>
              {item.name}
            </Text>
          </Group>
        </UnstyledButton>
      </HoverCardTarget>
      <HoverCardDropdown>
        <Text size="xs">{item.path}</Text>
      </HoverCardDropdown>
    </HoverCard>
  );
}

function DirectoryItem({
  item,
  isOpened,
  level = 0,
  searchMode = false,
  onClickItem,
}: {
  item: ProjectFile;
  isOpened: boolean;
  searchMode: boolean;
  level?: number;
  onClickItem: (item: ProjectFile) => void;
}) {
  const [files, setFiles] = useState<ProjectFile[]>([]);
  const [loading, setLoading] = useState(false);
  const [opened, setOpened] = useState(isOpened);

  const {
    data: directoryFiles,
    isLoading,
    isError,
  } = useChatFiles(item.path, null, !searchMode && loading);

  useEffect(() => {
    if (!isError && !isLoading && loading) {
      const timeout = setTimeout(() => {
        setFiles(directoryFiles || []);
        setOpened(true); // open after loading
        setLoading(false); // reset loading trigger
      }, 120);

      return () => clearTimeout(timeout); // cancel previous timeout if effect re-runs
    }
  }, [directoryFiles, isLoading, isError, loading]);

  return (
    <>
      <HoverCard shadow="md" openDelay={1000} position="top">
        <HoverCardTarget>
          <UnstyledButton
            onClick={() => {
              if (!opened && files.length === 0) {
                // first click on unopened folder → load
                setLoading(true);
              } else {
                // toggle open/close freely
                setOpened((o) => !o);
              }
            }}
            w="100%"
            p="4"
            pl={4 + level * 16}
          >
            <Group gap="xs" wrap="nowrap">
              {opened ? (
                <IconChevronDown size={14} />
              ) : (
                <IconChevronRight size={14} />
              )}
              {opened ? <IconFolderOpen size={16} /> : <IconFolder size={16} />}
              <Text size="sm" truncate>
                {item.name}
              </Text>
              {loading && <Loader size="8" ml="auto" mr="sm" />}
            </Group>
          </UnstyledButton>
        </HoverCardTarget>
        <HoverCardDropdown>
          <Text size="xs">{item.path}</Text>
        </HoverCardDropdown>
      </HoverCard>

      <Collapse in={opened}>
        {isLoading ? (
          <Text size="sm" c="dimmed" pl={4 + (level + 1) * 16}>
            Loading files...
          </Text>
        ) : (
          files?.map((child) =>
            child.type === "folder" ? (
              <DirectoryItem
                key={child.id}
                item={child}
                isOpened={false}
                level={level + 1}
                searchMode={searchMode}
                onClickItem={onClickItem}
              />
            ) : (
              <FileItem
                key={child.id}
                item={child}
                level={level + 1}
                onClickItem={onClickItem}
              />
            )
          )
        )}
      </Collapse>
    </>
  );
}

function MessageBubble({ message }: { message: Message }) {
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

export default function ProjectChat() {
  const [sidebarOpened, setSidebarOpened] = useState(true);
  const [messages, setMessages] = useState<Message[]>([]);
  const [inputValue, setInputValue] = useState("");
  const [searchValue, setSearchValue] = useState("");
  const [settingsOpened, setSettingsOpened] = useState(false);
  const [aiProvider, setAiProvider] = useState<string | null>("default");
  const [model, setModel] = useState<string | null>("default-model");
  const [files, setFiles] = useState<ProjectFile[]>([]);
  const [openedDirs, setOpenedDirs] = useState<Record<string, boolean>>({});
  const {
    data: selectedFiles,
    isLoading,
    error,
  } = useChatFiles(null, searchValue, true);
  const [navbarSize, setNavbarSize] = useState(250);

  useEffect(() => {
    if (selectedFiles) {
      setFiles(selectedFiles);
    }
  }, [selectedFiles]);

  const toggleDir = (item: ProjectFile) => {
    setOpenedDirs((prev) => ({
      ...prev,
      [item.id]: !prev[item.id],
    }));
  };

  const handleSendMessage = async () => {
    if (!inputValue.trim()) return;

    const newMessage = {
      id: Date.now().toString(),
      type: "user" as const,
      content: inputValue,
      timestamp: new Date(),
    };

    setMessages((prev) => [...prev, newMessage]);
    setInputValue("");

    const assistantResponse = await mockSendMessage(inputValue);
    setMessages((prev) => [...prev, assistantResponse as Message]);
  };

  const handleKeyPress = (event: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (event.key === "Enter" && !event.shiftKey) {
      event.preventDefault();
      handleSendMessage();
    }
  };

  const onMoveResizeHandler = (e: any) => {
    e.preventDefault();
    const startX = e.clientX;
    const startWidth = navbarSize;

    const onMouseMove = (e: MouseEvent) => {
      const newWidth = startWidth + (e.clientX - startX);
      setNavbarSize(Math.max(250, Math.min(newWidth, 550)));
    };

    const onMouseUp = () => {
      window.removeEventListener("mousemove", onMouseMove);
      window.removeEventListener("mouseup", onMouseUp);
    };

    window.addEventListener("mousemove", onMouseMove);
    window.addEventListener("mouseup", onMouseUp);
  };

  const chatLeft = useMemo(() => {
    if (sidebarOpened) {
      return `calc((100% + ${navbarSize}px) / 2 - 20%)`;
    }
    return "calc(100vw / 2 - 20vw)";
  }, [sidebarOpened, navbarSize]);

  return (
    <>
      {/* Settings Modal */}
      <Modal
        opened={settingsOpened}
        onClose={() => setSettingsOpened(false)}
        title={<Title order={4}>AI Assistant Settings</Title>}
        centered
        size="lg"
      >
        <Tabs defaultValue="ai" variant="outline">
          <Tabs.List>
            <Tabs.Tab value="ai" leftSection={<IconRobot size={16} />}>
              AI Configuration
            </Tabs.Tab>
            <Tabs.Tab value="general" leftSection={<IconSettings size={16} />}>
              General
            </Tabs.Tab>
          </Tabs.List>

          <Tabs.Panel value="ai" pt="md">
            <Stack gap="md">
              <Select
                label="AI Provider"
                placeholder="Select AI provider"
                data={[
                  { value: "default", label: "Default AI" },
                  { value: "openai", label: "OpenAI" },
                  { value: "anthropic", label: "Anthropic" },
                ]}
                value={aiProvider}
                onChange={setAiProvider}
              />
              <Select
                label="Model"
                placeholder="Select model"
                data={[
                  { value: "default-model", label: "Default Model" },
                  { value: "advanced", label: "Advanced Model" },
                  { value: "experimental", label: "Experimental Model" },
                ]}
                value={model}
                onChange={setModel}
                disabled={!aiProvider || aiProvider === "default"}
              />
              <Button onClick={() => setSettingsOpened(false)}>
                Save Settings
              </Button>
            </Stack>
          </Tabs.Panel>

          <Tabs.Panel value="general" pt="md">
            <Stack gap="md">
              <TextInput label="API Key" placeholder="Enter your API key" />
              <Select
                label="Theme"
                placeholder="Select theme"
                data={[
                  { value: "light", label: "Light" },
                  { value: "dark", label: "Dark" },
                ]}
              />
              <Button onClick={() => setSettingsOpened(false)}>
                Save Settings
              </Button>
            </Stack>
          </Tabs.Panel>
        </Tabs>
      </Modal>

      <AppShell
        layout="alt"
        navbar={{
          width: navbarSize,
          breakpoint: "sm",
          collapsed: { mobile: !sidebarOpened, desktop: !sidebarOpened },
        }}
        header={{ height: 46 }}
      >
        <AppShellHeader p="md">
          <Flex align="center" justify="space-between" h="100%">
            <Group gap="xs">
              {!sidebarOpened && (
                <ActionIcon
                  variant="subtle"
                  size="sm"
                  onClick={() => setSidebarOpened(true)}
                >
                  <IconChevronRight size={14} />
                </ActionIcon>
              )}
              <Text fw={600} size="lg">
                Project AI Assistant
              </Text>
            </Group>
            <Group gap="xs">
              <ActionIcon
                variant="subtle"
                onClick={() => setSettingsOpened(true)}
                size="lg"
              >
                <IconSettings size={18} />
              </ActionIcon>
              <Button
                variant="subtle"
                size="compact-md"
                leftSection={<IconMessage size={16} />}
                onClick={() => setMessages([])}
              >
                New Chat
              </Button>
            </Group>
          </Flex>
        </AppShellHeader>

        <AppShellNavbar p="xs">
          <Flex direction="column" h="100%" pos="relative">
            <Group mb="sm" justify="space-between">
              <Text fw={600} size="md">
                Project Files
              </Text>
              <Group gap={4}>
                <ActionIcon
                  variant="subtle"
                  size="sm"
                  onClick={() => setSidebarOpened(false)}
                >
                  <IconChevronLeft size={14} />
                </ActionIcon>
              </Group>
            </Group>

            <TextInput
              placeholder="Search files..."
              size="xs"
              leftSection={<IconSearch size={14} />}
              value={searchValue}
              onChange={(event) => setSearchValue(event.currentTarget.value)}
              mb="sm"
            />

            <ScrollArea style={{ flex: 1 }}>
              <Stack gap={2}>
                {isLoading && <Text>Loading...</Text>}
                {error && <Text color="red">Error loading files</Text>}
                {!isLoading && files && files.length > 0 ? (
                  files.map((item) =>
                    item.type === "folder" ? (
                      <DirectoryItem
                        key={item.id}
                        item={item}
                        isOpened={!!openedDirs[item.id]}
                        searchMode={!!searchValue}
                        level={0}
                        onClickItem={toggleDir}
                      />
                    ) : (
                      <FileItem
                        key={item.id}
                        item={item}
                        level={0}
                        onClickItem={toggleDir}
                      />
                    )
                  )
                ) : (
                  <Text size="sm" c="dimmed">
                    No files found
                  </Text>
                )}
              </Stack>
            </ScrollArea>
          </Flex>

          <Box
            pos="absolute"
            top={0}
            bottom={0}
            right={0}
            w={2}
            h="100%"
            className="resizer"
            onMouseDown={onMoveResizeHandler}
          />
        </AppShellNavbar>

        <AppShellMain p="0">
          <Flex direction="column" h="calc(100vh - 60px)" w="full">
            <ScrollArea style={{ flex: 1 }} p="md">
              <Stack gap="xs">
                {messages.length > 0 ? (
                  messages.map((message) => (
                    <MessageBubble key={message.id} message={message} />
                  ))
                ) : (
                  <Text size="sm" c="dimmed" ta="center" mt="xl">
                    Start a new conversation
                  </Text>
                )}
              </Stack>
            </ScrollArea>

            <Paper
              p="md"
              shadow="sm"
              style={{ transition: "all 0.2s ease" }}
              withBorder
              radius="lg"
              pos="fixed"
              bottom="12px"
              w="40%"
              left={chatLeft}
            >
              <Group align="flex-end" gap="xs">
                <Textarea
                  placeholder="Ask about your project..."
                  value={inputValue}
                  onChange={(event) => setInputValue(event.currentTarget.value)}
                  onKeyPress={handleKeyPress}
                  minRows={1}
                  maxRows={6}
                  autosize
                  style={{ flex: 1 }}
                  variant="filled"
                  styles={{ input: { pt: "120px" } }}
                />
                <ActionIcon
                  size="lg"
                  variant="filled"
                  color="blue"
                  onClick={handleSendMessage}
                  disabled={!inputValue.trim()}
                >
                  <IconSend size={16} />
                </ActionIcon>
              </Group>
            </Paper>
          </Flex>
        </AppShellMain>
      </AppShell>
    </>
  );
}
