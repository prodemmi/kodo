import React, { useState, useEffect, useMemo } from "react";
import {
  AppShell,
  AppShellNavbar,
  AppShellHeader,
  AppShellMain,
  Text,
  Group,
  ActionIcon,
  TextInput,
  ScrollArea,
  Stack,
  Paper,
  Flex,
  Textarea,
  Modal,
  Select,
  Button,
  Tabs,
  Title,
  Box,
  MultiSelect,
} from "@mantine/core";
import {
  IconChevronRight,
  IconSend,
  IconRobot,
  IconSearch,
  IconSettings,
  IconMessage,
  IconChevronLeft,
  IconContainer,
} from "@tabler/icons-react";
import { useChatFiles } from "../../../../hooks/use-chat";
import { Message, ProjectFile } from "../../../../types/chat";
import { FileItem } from "./sections/FileItem";
import ContextManager from "./sections/ContextManager";
import { useContextStore } from "../../../../states/context.state";
import DirectoryItem from "./sections/DirectoryItem";
import MessageBubble from "./sections/MessageBubble";

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
  const [isContextManagerOpened, setIsContextManagerOpened] = useState(false);
  const contexts = useContextStore((s) => s.contexts);
  const selectedContexts = useContextStore((s) => s.selectedContexts);
  const toggleContexts = useContextStore((s) => s.toggleContexts);

  const contextOptionsData = useMemo(() => {
    return contexts.map((ctx) => ({
      value: ctx.id,
      label: ctx.name,
    }));
  }, [contexts]);

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
              <MultiSelect
                size="xs"
                placeholder={
                  selectedContexts.length === 0 ? "Pick a Context" : undefined
                }
                data={contextOptionsData}
                value={selectedContexts} // Changed from defaultValue to value
                onChange={(val) => toggleContexts(val)}
                searchable
              />
            </Group>
            <Group gap="xs">
              <Button
                size="sm"
                variant="subtle"
                leftSection={<IconContainer size={16} />}
                onClick={() => setIsContextManagerOpened(true)}
              >
                Context Manager
              </Button>
              <Button
                size="sm"
                variant="subtle"
                leftSection={<IconMessage size={16} />}
                onClick={() => setMessages([])}
              >
                New Chat
              </Button>
              <ActionIcon
                variant="subtle"
                onClick={() => setSettingsOpened(true)}
                size="lg"
              >
                <IconSettings size={18} />
              </ActionIcon>
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
                  files.map((item) => (
                    <Group justify="space-between" m="0" p="0">
                      {item.type === "folder" ? (
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
                          onAddToContext={() => alert("added to context")}
                          onCopyPath={() => alert("added to context")}
                          onDelete={() => alert("added to context")}
                          onDownload={() => alert("added to context")}
                          onPreview={() => alert("added to context")}
                          onSelectCode={() => alert("added to context")}
                        />
                      )}
                    </Group>
                  ))
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
      <ContextManager
        opened={isContextManagerOpened}
        onClose={() => setIsContextManagerOpened(false)}
      />
    </>
  );
}
