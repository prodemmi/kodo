import { Container, Box, Stack, Flex, Group, Divider } from "@mantine/core";

import "@mantine/tiptap/styles.css";
import CreateNoteModal from "./CreateNoteModal";
import CreateFolderModal from "./CreateFolderModal";
import MainContent from "./MainContent";
import NoteList from "./NotesList";
import TopHeader from "./TopHeader";
import DeleteConfirmationModal from "./DeleteConfirmationModal";
import Folders from "./Folders";
import { useMediaQuery } from "@mantine/hooks";
import { useFolders, useNotes } from "../../../../hooks/use-notes";
import { useNoteStore } from "../../../../states/note.state";
import { useEffect } from "react";

export default function Notes() {
  const isSmall = useMediaQuery("(max-width: 920px)");
  const {
    data: notesData,
    isError: notesError,
    isLoading: notesLoading,
  } = useNotes();
  const {
    data: foldersData,
    isError: foldersError,
    isLoading: foldersLoading,
  } = useFolders();

  const setNotes = useNoteStore((s) => s.setNotes);
  const setFolders = useNoteStore((s) => s.setFolders);

  useEffect(() => {
    if (!notesError && !notesLoading && notesData && notesData.count > 0) {
      setNotes(notesData.notes);
    }
  }, [notesData, notesError, notesLoading]);

  useEffect(() => {
    if (
      !foldersError &&
      !foldersLoading &&
      foldersData &&
      foldersData.count > 0
    ) {
      setFolders(foldersData.folders);
    }
  }, [foldersData, foldersError, foldersLoading]);

  return (
    <Container
      fluid
      p="0"
      m="0"
      style={{
        height: "calc(100vh - 86px)",
        display: "flex",
        flexDirection: "column",
        overflow: "hidden",
      }}
    >
      {/* Top Header */}
      <Box
        style={{
          display: "flex",
          flex: 1,
          overflow: "hidden",
        }}
      >
        {/* Sidebar */}
        <Flex
          gap="xs"
          align="flex-start"
          style={{
            width: isSmall ? "100%" : "65dvw",
            flexDirection: isSmall ? "column" : "row",
          }}
        >
          {/* Folders */}
          <Folders />

          <Divider
            orientation={isSmall ? "horizontal" : "vertical"}
            w={isSmall ? "100%" : undefined}
          />

          {/* Notes List */}
          <NoteList />

          {!isSmall && <Divider orientation="vertical" />}
        </Flex>

        {/* Main Content */}
        <Stack w="100%">
          <TopHeader />
          <MainContent />
        </Stack>
      </Box>

      {/* Create Note Modal */}
      <CreateNoteModal />

      {/* Create Folder Modal */}
      <CreateFolderModal />

      {/* Delete Confirmation Modal */}
      <DeleteConfirmationModal />
    </Container>
  );
}
