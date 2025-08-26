import { useNoteStore } from "../../../../../states/note.state";
import { getTaskListExtension, Link } from "@mantine/tiptap";
import CodeBlockLowlight from "@tiptap/extension-code-block-lowlight";
import Superscript from "@tiptap/extension-superscript";
import TextAlign from "@tiptap/extension-text-align";
import Underline from "@tiptap/extension-underline";
import { useEditor } from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";
import SubScript from "@tiptap/extension-subscript";
import Highlight from "@tiptap/extension-highlight";
import TaskItem from "@tiptap/extension-task-item";
import TipTapTaskList from "@tiptap/extension-task-list";
import { createLowlight } from "lowlight";
import { useMemo } from "react";
import { Stack } from "@mantine/core";
import NoteInfo from "./sections/NoteInfo";
import WelcomeState from "./sections/WelcomeState";
import NoteEditor from "./sections/NoteEditor";

export default function MainContent() {
  const selectedNote = useNoteStore((s) => s.selectedNote);
  const lowlight = useMemo(() => createLowlight(), []);

  const editor = useEditor({
    autofocus: true,
    extensions: [
      StarterKit,
      Underline,
      Link.configure({
        HTMLAttributes: {
          class: "editor-link",
        },
      }),
      Superscript,
      SubScript,
      Highlight,
      getTaskListExtension(TipTapTaskList),
      TaskItem.configure({
        nested: true,
        HTMLAttributes: {
          class: "test-item",
        },
      }),
      TextAlign.configure({ types: ["heading", "paragraph"] }),
      CodeBlockLowlight.configure({
        lowlight,
      }),
    ],
    content: selectedNote?.content || "",
    editable: true, // Enable editing to test mentions
  });

  return (
    <Stack gap="xs" p="xs" flex={1}>
      {selectedNote ? (
        <>
          {/* Editor Header */}
          <NoteInfo editor={editor} />

          {/* Editor/Content */}
          <NoteEditor editor={editor} />
        </>
      ) : (
        /* Welcome State */
        <WelcomeState />
      )}
    </Stack>
  );
}
