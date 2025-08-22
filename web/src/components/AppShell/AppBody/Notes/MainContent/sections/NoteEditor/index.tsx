import { RichTextEditor } from "@mantine/tiptap";
import { useEffect } from "react";
import { useNoteStore } from "../../../../../../../states/note.state";
import { RoleGuard } from "../../../../../../Investor";
import { Editor } from "@tiptap/react";
import { Box } from "@mantine/core";

interface Props {
  editor: Editor;
}

export default function NoteEditor({ editor }: Props) {
  const selectedNote = useNoteStore((s) => s.selectedNote);
  const isEditingNote = useNoteStore((s) => s.isEditingNote);

  useEffect(() => {
    if (selectedNote) editor.commands.setContent(selectedNote.content);
  }, [selectedNote, editor]);

  useEffect(() => {
    if (editor) {
      editor.setEditable(isEditingNote);
    }
  }, [editor, isEditingNote]);

  return (
    <Box
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
              style={{ border: "none" }}
              bg="var(---mantine-color-dark-4)"
              m="0"
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

              <RichTextEditor.ControlsGroup>
                <RichTextEditor.TaskList />
                <RichTextEditor.TaskListLift />
                <RichTextEditor.TaskListSink />
              </RichTextEditor.ControlsGroup>
            </RichTextEditor.Toolbar>
          )}
        </RoleGuard.Consumer>

        <RichTextEditor.Content
          bg="transparent"
          style={{
            cursor: isEditingNote ? "text" : "auto",
            fontSize: "16px",
            lineHeight: "1.6",
            minHeight: "calc(100dvh - 348px)",
            overflow: "auto",
          }}
        />
      </RichTextEditor>
    </Box>
  );
}
