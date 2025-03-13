import React, { useState, useEffect } from "react";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { useParams } from "react-router-dom";

import Title, { TitleProps, DocumentState } from "./Title";
import { useCurrentUserContext } from "@/contexts/CurrentUserContext";
import { useErrorToast } from "@/hooks/useErrorToast";
import { useDocumentContext } from "@/contexts/DocumentContext";

type TitleContainerProps = Pick<TitleProps, "maximized" | "toggleMaximize">;

function TitleContainer(props: TitleContainerProps) {
  const { editor } = useRogueEditorContext();
  const [lastEdit, setLastEdit] = useState<Date | null>(null);
  const [documentState, setDocumentState] = useState<DocumentState>(
    DocumentState.Loading,
  );
  const { draftId } = useParams();
  const { currentUser } = useCurrentUserContext();
  const showErrorToast = useErrorToast();
  const {
    docData,
    deleteDocument,
    savingDocument,
    deletingDocument,
    updateDocument,
  } = useDocumentContext();

  const updateTitle = async (newTitle: string) => {
    const originalTitle = docData?.title;
    let validatedTitle = newTitle;
    if (validatedTitle.length === 0) {
      validatedTitle = "Untitled";
    }
    document.title = `${validatedTitle} - Pointy`;

    const { errors } = await updateDocument({
      title: validatedTitle,
    });

    if (errors) {
      document.title = `${originalTitle} - Pointy`;
      showErrorToast("Title failed to save");
    }
  };

  useEffect(() => {
    if (!editor) {
      return;
    }

    const handleChanges = (newValue: Date) => {
      setLastEdit(newValue);
    };

    const updateConnectionState = (_: any) => {
      if (!editor.connected) {
        setDocumentState(DocumentState.Disconnected);
        return;
      }

      if (!editor.syncing) {
        setDocumentState(DocumentState.Saved);
        return;
      }

      setDocumentState(DocumentState.Loading);
    };

    const onEditorKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Enter" && docData?.title === "Untitled") {
        const html = editor.container?.innerHTML;
        if (html?.startsWith("<h1")) {
          const parser = new DOMParser();
          const doc = parser.parseFromString(html, "text/html");

          const headings = doc.querySelectorAll("h1");
          const title = headings[0].textContent;
          if (title) {
            updateTitle(title);
            // we don't need the keydown listener anymore
            editor.removeEventListener("keydown", onEditorKeyDown);
          }
        }
      }
    };

    editor.subscribe<Date>("lastEdit", handleChanges);
    editor.subscribe<boolean>("connected", updateConnectionState);
    editor.subscribe<boolean>("syncing", updateConnectionState);

    if (editor && docData?.title === "Untitled") {
      editor.addEventListener("keydown", onEditorKeyDown);
    }

    return () => {
      editor.unsubscribe<Date>("lastEdit", handleChanges);
      editor.unsubscribe<boolean>("connected", updateConnectionState);
      editor.unsubscribe<boolean>("syncing", updateConnectionState);
      if (editor && docData?.title === "Untitled") {
        editor.removeEventListener("keydown", onEditorKeyDown);
      }
    };
  }, [editor, docData?.title]);

  useEffect(() => {
    if (docData) {
      document.title = `${docData.title} - Pointy`;
      if (docData.updatedAt) {
        setLastEdit(new Date(docData.updatedAt));
      }
    }
  }, [docData]);

  const handleDeleteDocument = () => {
    draftId && deleteDocument(draftId);
  };

  const onCopy = (): boolean => {
    if (editor) {
      editor.copyDoc();
      return true;
    }
    return false;
  };

  if (!docData) {
    return <div>Loading...</div>;
  }

  return (
    <Title
      {...props}
      me={currentUser}
      document={docData}
      documentState={documentState}
      loading={savingDocument}
      updateTitle={updateTitle}
      deleteDocument={handleDeleteDocument}
      deleteLoading={deletingDocument}
      lastEdit={lastEdit}
      handleCopy={onCopy}
    />
  );
}

export default TitleContainer;
