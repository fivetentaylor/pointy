import React, { useEffect, useState, useRef } from "react";
import { useParams } from "react-router-dom";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { RogueEditor as RogueEditorElement } from "../../../rogueEditor";
import { cn } from "@/lib/utils";
import { Skeleton } from "../ui/skeleton";

const RogueEditor = () => {
  const [loading, setLoading] = useState(true);
  const [htmlContent, setHtmlContent] = useState("");
  const [address, setAddress] = useState<string | null>(null);
  const rogueEditorRef = useRef<HTMLDivElement | null>(null);
  const {
    editor,
    setCurrentEditor,
    uploadImage,
    getImage,
    listDocumentImages,
    getImageSignedUrl,
  } = useRogueEditorContext();
  const { draftId } = useParams();

  const onAddressChange = (value: string | null) => {
    setAddress(value);
  };

  useEffect(() => {
    if (editor) {
      editor.subscribe<string | null>("address", onAddressChange);
    }

    return () => {
      if (editor) {
        editor.unsubscribe<string | null>("address", onAddressChange);
      }
    };
  }, [editor, draftId]);

  useEffect(() => {
    setLoading(true);

    fetch(`/api/v1/documents/${draftId}/editor.html`)
      .then((response) => response.text())
      .then((data) => {
        setLoading(false);
        setHtmlContent(data);
      })
      .catch((error) => {
        console.error("Error fetching HTML:", error);
        setLoading(false);
      });
  }, [draftId]);

  useEffect(() => {
    if (htmlContent && rogueEditorRef.current) {
      const rogueEditorElement = rogueEditorRef.current.querySelector(
        "rogue-editor",
      ) as RogueEditorElement;
      if (rogueEditorElement) {
        // Directly assign the functions to the RogueEditor
        rogueEditorElement.uploadImage = uploadImage;
        rogueEditorElement.getImage = getImage;
        rogueEditorElement.listDocumentImages = listDocumentImages;
        rogueEditorElement.getImageSignedUrl = getImageSignedUrl;

        setCurrentEditor(rogueEditorElement);
        rogueEditorElement.resetAddress();
      }
    }
  }, [
    htmlContent,
    setCurrentEditor,
    uploadImage,
    getImage,
    listDocumentImages,
    getImageSignedUrl,
  ]);

  if (loading) {
    return (
      <div className="pt-[1.8125rem] pb-[30dvh] min-h-[90dvh]">
        <Skeleton className="h-[2.375rem] w-[24.635rem] bg-elevated mb-[1.3125rem]" />
        <Skeleton className="h-6 w-full bg-elevated mb-[1.3125rem]" />
        <Skeleton className="h-6 w-full bg-elevated mb-[1.3125rem]" />
        <Skeleton className="h-6 w-[11.56rem] bg-elevated mb-[1.3125rem]" />
      </div>
    );
  }

  return (
    <>
      <style>
        {`
          rogue-editor {
            position: relative;
          }
        `}
      </style>
      <div
        dangerouslySetInnerHTML={{ __html: htmlContent }}
        ref={rogueEditorRef}
        className={cn(
          "pt-[1.8125rem] pb-[30dvh] min-h-[90dvh]",
          address !== null ? "readonly" : "",
        )}
      />
    </>
  );
};

export default RogueEditor;
