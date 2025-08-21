import { RichTextEditor } from "@mantine/tiptap";
import { useEffect } from "react";
import { useNoteStore } from "../../../../../../../states/note.state";
import { RoleGuard } from "../../../../../../Investor";
import { Editor } from "@tiptap/react";

interface Props {
  editor: Editor;
}

export default function NoteEditor({ editor }: Props) {
  const selectedNote = useNoteStore((s) => s.selectedNote);
  const isEditingNote = useNoteStore((s) => s.isEditingNote);

  useEffect(() => {
    if (selectedNote) editor.commands.setContent(selectedNote.content);
  }, [selectedNote, editor]);

  return (
    <div
      style={{
        flex: 1,
        display: "flex",
        flexDirection: "column",
        overflow: "hidden",
      }}
    >
      <RichTextEditor
        editor={editor}
        key={isEditingNote ? "1" : "0"}
        style={{
          border: "none",
          display: "flex",
          flexDirection: "column",
          height: "100%",
        }}
      >
        <RoleGuard.Consumer>
          {isEditingNote && (
            <RichTextEditor.Toolbar
              sticky
              stickyOffset={0}
              style={{
                borderBottom: "1px solid var(--mantine-color-gray-8)",
                zIndex: 99,
                padding: "8px 16px",
              }}
            >
              <RichTextEditor.ControlsGroup>
                <RichTextEditor.Bold />
                <RichTextEditor.Italic />
                <RichTextEditor.Underline />
                <RichTextEditor.Strikethrough />
                <RichTextEditor.ClearFormatting />
                <RichTextEditor.Highlight />
                <RichTextEditor.Code />
              </RichTextEditor.ControlsGroup>

              <RichTextEditor.ControlsGroup>
                <RichTextEditor.H1 />
                <RichTextEditor.H2 />
                <RichTextEditor.H3 />
                <RichTextEditor.H4 />
              </RichTextEditor.ControlsGroup>

              <RichTextEditor.ControlsGroup>
                <RichTextEditor.Blockquote />
                <RichTextEditor.Hr />
                <RichTextEditor.BulletList />
                <RichTextEditor.OrderedList />
                <RichTextEditor.Subscript />
                <RichTextEditor.Superscript />
              </RichTextEditor.ControlsGroup>

              <RichTextEditor.ControlsGroup>
                <RichTextEditor.Link />
                <RichTextEditor.Unlink />
              </RichTextEditor.ControlsGroup>

              <RichTextEditor.ControlsGroup>
                <RichTextEditor.AlignLeft />
                <RichTextEditor.AlignCenter />
                <RichTextEditor.AlignJustify />
                <RichTextEditor.AlignRight />
              </RichTextEditor.ControlsGroup>

              <RichTextEditor.ControlsGroup>
                <RichTextEditor.Undo />
                <RichTextEditor.Redo />
              </RichTextEditor.ControlsGroup>

              <RichTextEditor.ControlsGroup>
                <RichTextEditor.CodeBlock />
              </RichTextEditor.ControlsGroup>
            </RichTextEditor.Toolbar>
          )}
        </RoleGuard.Consumer>

        <RichTextEditor.Content
          bg="transparent"
          style={{
            flex: 1,
            fontSize: "16px",
            lineHeight: "1.6",
            overflow: "auto",
          }}
        />
      </RichTextEditor>
    </div>
  );
}
