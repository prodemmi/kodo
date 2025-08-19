import {
  Stack,
  Text,
  Drawer,
  Loader,
  Badge,
  Group,
  Card,
  Collapse,
  Button,
} from "@mantine/core";
import {
  CodeHighlight,
  CodeHighlightAdapterProvider,
  createShikiAdapter,
} from "@mantine/code-highlight";
import { Item } from "../../../../../../types/item";
import { useItemContext, useOpenFile } from "../../../../../../hooks/use-items";
import { useMemo, useEffect, useState, useCallback } from "react";

// Import styles for CodeHighlight
import "@mantine/code-highlight/styles.css";
import {
  IconChevronDown,
  IconChevronRight,
  IconCode,
  IconFileText,
} from "@tabler/icons-react";

// Shiki requires async code to load the highlighter
async function loadShiki() {
  const { createHighlighter } = await import("shiki");
  const shiki = await createHighlighter({
    langs: [
      "tsx",
      "scss",
      "html",
      "bash",
      "json",
      "go",
      "javascript",
      "typescript",
      "jsx",
      "python",
      "css",
    ],
    themes: ["github-light", "github-dark"],
  });

  return shiki;
}

const shikiAdapter = createShikiAdapter(loadShiki);

type Props = {
  selectedItem: Item | null;
  drawerOpened: boolean;
  setDrawerOpened: (drawerOpened: boolean) => void;
};

export default function ItemDetailDrawer({
  drawerOpened,
  selectedItem,
  setDrawerOpened,
}: Props) {
  const { data, isLoading } = useItemContext(
    selectedItem?.file!,
    selectedItem?.line!
  );

  const { mutate } = useOpenFile();

  const [showCodeContext, setShowCodeContext] = useState(true);

  const language = useMemo(() => {
    const splitted = data && data.file ? data.file.split(".") : [];
    const ext = splitted[splitted.length - 1];
    const languageMap: { [key: string]: string } = {
      js: "javascript",
      ts: "typescript",
      tsx: "tsx",
      jsx: "jsx",
      py: "python",
      css: "css",
      html: "html",
      go: "go",
    };
    return languageMap[ext] || "plaintext";
  }, [data]);

  const plainCode = useMemo(() => {
    if (!data?.lines) return "";
    return data.lines.map((line) => line.content).join("\n");
  }, [data]);

  const targetLineNumber = useMemo(() => {
    if (!selectedItem?.line || !data?.lines) return null;
    const lineIndex = data.lines.findIndex(
      (line) => line.number === selectedItem.line
    );
    return lineIndex !== -1 ? lineIndex + 1 : null;
  }, [selectedItem?.line, data?.lines]);

  const gotoFile = useCallback(
    (item: Item) => {
      mutate({ filename: item.file, line: item.line });
    },
    [mutate]
  );

  // Apply highlighting using DOM manipulation after render
  useEffect(() => {
    if (!targetLineNumber) return;

    const timer = setTimeout(() => {
      const codeElements = document.querySelectorAll(
        "code.mantine-CodeHighlight-code .line"
      );
      console.log('codeElements ===>', codeElements);
      

      codeElements.forEach((element) => {
        element.classList.remove("highlighted-line");
      });

      if (codeElements[targetLineNumber - 1]) {
        codeElements[targetLineNumber - 1].classList.add("highlighted-line");
      }
    }, 100);

    return () => clearTimeout(timer);
  }, [targetLineNumber, plainCode]);

  // Inject CSS for highlighted lines
  useEffect(() => {
    const style = document.createElement("style");
    style.textContent = `
      .highlighted-line {
        background-color: rgba(255, 0, 0, 0.1) !important;
        border-left: 3px solid #ff4444 !important;
      }
      .highlighted-line td {
        background-color: rgba(255, 0, 0, 0.1) !important;
      }
    `;
    document.head.appendChild(style);

    return () => {
      document.head.removeChild(style);
    };
  }, []);

  return (
    <Drawer
      opened={drawerOpened}
      onClose={() => setDrawerOpened(false)}
      title="Item Details"
      position="right"
      size="xl"
    >
      {selectedItem && (
        <Stack gap="md">
          {/* TODO Item Information */}
          <Card withBorder p="md">
            <Stack gap="xs">
              <Group justify="space-between">
                <Text fw={600} size="lg">
                  {selectedItem.title}
                </Text>
                <Badge color="blue" variant="light">
                  {selectedItem.status}
                </Badge>
              </Group>
              <Text
                size="sm"
                c="dimmed"
                styles={{ root: { whiteSpace: "break-spaces" } }}
              >
                {selectedItem.description}
              </Text>
              <Group gap="sm">
                <Badge color="orange" variant="outline" size="sm">
                  Priority: {selectedItem.priority}
                </Badge>
                <Badge color="gray" variant="outline" size="sm">
                  Type: {selectedItem.type}
                </Badge>
              </Group>
              <Group align="center" justify="space-between">
                <Text size="xs" c="dimmed" style={{ fontFamily: "monospace" }}>
                  📁 {selectedItem.file}:{selectedItem.line}
                </Text>
                <Button
                  size="compact-sm"
                  variant="light"
                  leftSection={<IconCode size={16} />}
                  onClick={() => gotoFile(selectedItem)}
                >
                  Go to file
                </Button>
              </Group>
            </Stack>
          </Card>

          {isLoading && <Loader size="sm" />}

          {data && !data.error && (
            <>
              {/* Code Context */}
              <Card withBorder p="md">
                <Group justify="space-between" mb="sm">
                  <Group gap="sm">
                    <IconFileText size={16} color="#495057" />
                    <Text size="sm" fw={600}>
                      {data.file}
                    </Text>
                  </Group>
                  <Button
                    variant="subtle"
                    size="xs"
                    rightSection={
                      showCodeContext ? (
                        <IconChevronDown size={12} />
                      ) : (
                        <IconChevronRight size={12} />
                      )
                    }
                    onClick={() => setShowCodeContext(!showCodeContext)}
                  >
                    {showCodeContext ? "Hide" : "Show"} Code
                  </Button>
                </Group>

                <Collapse in={showCodeContext}>
                  {data.lines && (
                    <CodeHighlightAdapterProvider adapter={shikiAdapter}>
                      <CodeHighlight
                        code={plainCode}
                        language={language}
                        defaultExpanded
                        withExpandButton={false}
                        withCopyButton
                      />
                    </CodeHighlightAdapterProvider>
                  )}
                </Collapse>
              </Card>
            </>
          )}

          {data?.error && (
            <Card withBorder p="md" style={{ borderColor: "#fa5252" }}>
              <Text c="red" size="sm" fw={500}>
                ❌ Error: {data.error}
              </Text>
            </Card>
          )}
        </Stack>
      )}
    </Drawer>
  );
}
