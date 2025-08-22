import { useEffect, useState, useRef, useCallback, useMemo } from "react";
import { Group, TextInput, ActionIcon } from "@mantine/core";
import { IconEdit } from "@tabler/icons-react";
import { useNoteStore } from "../../../../../../../../../states/note.state";
import { useUpdateNote } from "../../../../../../../../../hooks/use-notes";

export default function NoteTitle() {
  const selectedNote = useNoteStore((s) => s.selectedNote);
  const [value, setValue] = useState(selectedNote?.title || "");
  const [focus, setFocus] = useState(false);
  const [error, setError] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);
  const measureRef = useRef<HTMLSpanElement>(null);
  const { mutate, isPending } = useUpdateNote();

  // Memoized original title for comparison
  const originalTitle = useMemo(
    () => selectedNote?.title || "",
    [selectedNote?.title]
  );

  useEffect(() => {
    setValue(originalTitle);
  }, [originalTitle]);

  useEffect(() => {
    if (!value.trim()) {
      setError(true);
    } else {
      setError(false);
    }
  }, [value]);

  // Calculate width based on text content
  const calculateWidth = useMemo(() => {
    if (!measureRef.current) return 200;

    // Create a temporary span to measure text width
    const canvas = document.createElement("canvas");
    const context = canvas.getContext("2d");
    if (context) {
      context.font =
        '500 14px -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif';
      const textWidth = context.measureText(value || "").width;
      const padding = 54;
      return textWidth + padding;
    }

    return 200;
  }, [value]);

  const handleFocus = useCallback(() => {
    setFocus(true);
    setIsEditing(true);
    setTimeout(() => inputRef.current?.select(), 0);
  }, []);

  const handleBlur = useCallback(() => {
    if (error) {
      return;
    }
    setFocus(false);
    setIsEditing(false);

    // Update only if there are changes
    if (selectedNote && value.trim() !== originalTitle) {
      mutate({
        ...selectedNote,
        title: value.trim(),
        id: selectedNote.id,
      });
    }
  }, [selectedNote, value, error, originalTitle, mutate]);

  const handleSave = useCallback(() => {
    if (error) {
      return;
    }
    if (!selectedNote || value.trim() === originalTitle) return;

    mutate({
      ...selectedNote,
      title: value.trim(),
      id: selectedNote.id,
    });

    setFocus(false);
    setIsEditing(false);
  }, [selectedNote, value, originalTitle, mutate]);

  const handleCancel = useCallback(() => {
    setValue(originalTitle);
    setFocus(false);
    setIsEditing(false);
  }, [originalTitle]);

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
      e.stopPropagation();

      if (e.key === "Enter") {
        e.preventDefault();
        handleSave();
      } else if (e.key === "Escape") {
        e.preventDefault();
        handleCancel();
      }
    },
    [handleSave, handleCancel]
  );

  const handleEditClick = useCallback(
    (e: React.MouseEvent) => {
      e.stopPropagation();
      handleFocus();
    },
    [handleFocus]
  );

  if (!selectedNote) return null;

  const isEmpty = !value.trim();

  return (
    <>
      {/* Hidden span for text measurement */}
      <span
        ref={measureRef}
        style={{
          visibility: "hidden",
          position: "absolute",
          whiteSpace: "nowrap",
          font: '500 14px -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
        }}
      >
        {value || "Untitled Note"}
      </span>

      <Group
        gap="xs"
        wrap="nowrap"
        align="center"
        style={{ position: "relative", width: "max-content" }}
      >
        <TextInput
          ref={inputRef}
          value={value}
          onBlur={handleBlur}
          onKeyDown={handleKeyDown}
          onChange={(e) => setValue(e.currentTarget.value)}
          size="md"
          mr="xl"
          ml="0"
          placeholder="Write Note Title"
          disabled={isPending}
          c={isEmpty ? "var(--mantine-color-dimmed)" : undefined}
          fw={focus || !isEmpty ? 500 : 400}
          styles={{
            input: {
              border: !error && !focus ? "none" : undefined,
              borderColor: error
                ? "red"
                : focus
                ? "var(--mantine-color-blue-6)"
                : "var(--mantine-color-dark-3)",
              backgroundColor: "transparent",
              pointerEvents: focus ? "auto" : "none",
              transition: "all 0.2s ease",
              minWidth: value.trim() ? 80 : 160,
              width: calculateWidth,
              cursor: !focus ? "pointer" : "text",
              padding: focus ? undefined : 0
            },
          }}
        />

        {/* Edit icon when not focused */}
        {!focus && !isPending && (
          <ActionIcon
            size="sm"
            variant="subtle"
            onClick={handleEditClick}
            style={{
              position: "absolute",
              right: 4,
              opacity: 0.6,
            }}
          >
            <IconEdit size={14} />
          </ActionIcon>
        )}
      </Group>
    </>
  );
}
