import { useNoteStore } from "../../../../../states/note.state";
import { Link } from "@mantine/tiptap";
import CodeBlockLowlight from "@tiptap/extension-code-block-lowlight";
import Superscript from "@tiptap/extension-superscript";
import TextAlign from "@tiptap/extension-text-align";
import Underline from "@tiptap/extension-underline";
import { useEditor } from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";
import SubScript from "@tiptap/extension-subscript";
import Highlight from "@tiptap/extension-highlight";
import { createLowlight } from "lowlight";
import { useMemo } from "react";
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
      Link,
      Superscript,
      SubScript,
      Highlight,
      TextAlign.configure({ types: ["heading", "paragraph"] }),
      CodeBlockLowlight.configure({
        lowlight,
      }),
    ],
    content: selectedNote?.content || "",
    editable: false,
  });

  return (
    <div
      style={{
        flex: 1,
        display: "flex",
        flexDirection: "column",
        overflow: "hidden",
      }}
    >
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
    </div>
  );
}
