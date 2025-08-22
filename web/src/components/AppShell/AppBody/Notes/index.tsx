import {
  Container,
  Box,
  Stack,
  Flex,
  Group,
  Divider,
  LoadingOverlay,
} from "@mantine/core";

import "@mantine/tiptap/styles.css";
import CreateNoteModal from "./CreateNoteModal";
import CreateFolderModal from "./CreateFolderModal";
import MainContent from "./MainContent";
import NoteList from "./NotesList";
import DeleteConfirmationModal from "./DeleteConfirmationModal";
import Folders from "./Folders";
import { useMediaQuery } from "@mantine/hooks";
import { useFolders, useNotes } from "../../../../hooks/use-notes";
import { useNoteStore } from "../../../../states/note.state";
import { useEffect } from "react";
import HistoryDrawer from "./HistoryModal";

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
    if (
      !foldersError &&
      !foldersLoading &&
      foldersData &&
      foldersData.count > 0
    ) {
      setFolders(foldersData.folders);
    }
    if (!notesError && !notesLoading && notesData && notesData.count > 0) {
      setNotes(notesData.notes);
    }
  }, [
    notesData,
    notesError,
    notesLoading,
    foldersData,
    foldersError,
    foldersLoading,
  ]);

  if (notesLoading || foldersLoading) {
    return (
      <LoadingOverlay zIndex={1000} overlayProps={{ radius: "sm", blur: 2 }} />
    );
  }

  return (
    <Container
      fluid
      p="0"
      m="0"
      style={{
        height: "calc(100dvh - 52px)",
        display: "flex",
        flexDirection: "column",
      }}
    >
      {/* Top Header */}
      <Box
        style={{
          display: "flex",
          flex: 1,
          overflow: "hidden",
        }}
        h="100dvh"
      >
        {/* Sidebar */}
        <Flex
          gap="xs"
          align="flex-start"
          h="100%"
          style={{
            width: isSmall ? "auto" : "45dvw",
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
        <MainContent />
      </Box>

      {/* Create Note Modal */}
      <CreateNoteModal />

      {/* Create Folder Modal */}
      <CreateFolderModal />

      {/* Delete Confirmation Modal */}
      <DeleteConfirmationModal />

      {/* Note History Modal */}
      <HistoryDrawer />
    </Container>
  );
}
