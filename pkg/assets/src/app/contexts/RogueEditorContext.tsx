"use client";

import { RogueEditor } from "@/../rogueEditor";
import React, { useContext, useEffect, useState, useCallback } from "react";

import { useLazyQuery, useMutation } from "@apollo/client";
import {
  UploadImage,
  GetImage,
  ListDocumentImages,
  GetImageSignedUrl,
} from "@/queries/images";

type RogueEditorContextState = {
  editor: RogueEditor | undefined;
  setCurrentEditor: (value: RogueEditor) => void;
  address: string | null;
  editorMode: "diff" | "history" | "paste" | "edit" | "xray";
  uploadImage: (file: File, docId: string) => Promise<Image | undefined>;
  getImage: (docId: string, imageId: string) => Promise<Image | undefined>;
  listDocumentImages: (docId: string) => Promise<Image[] | undefined>;
  getImageSignedUrl: (
    docId: string,
    imageId: string,
  ) => Promise<string | undefined>;
};

type RogueEditorContextProviderProps = {
  children: React.ReactNode;
};

const RogueEditorContext = React.createContext<
  RogueEditorContextState | undefined
>(undefined);

export const RogueEditorContextProvider: React.FC<
  RogueEditorContextProviderProps
> = ({ children }) => {
  const [editor, setCurrentEditor] = useState<RogueEditor>();
  const [address, _setAddress] = useState<string | null>(null);
  const [editorMode, _setEditorMode] = useState<
    "diff" | "history" | "paste" | "edit" | "xray"
  >("diff");

  const [uploadImageMutation] = useMutation(UploadImage);
  const [getImageQuery] = useLazyQuery(GetImage);
  const [listDocumentImagesQuery] = useLazyQuery(ListDocumentImages);
  const [getImageSignedUrlQuery] = useLazyQuery(GetImageSignedUrl);

  const uploadImage = useCallback(
    async (file: File, docId: string): Promise<Image | undefined> => {
      try {
        const { data } = await uploadImageMutation({
          variables: { file, docId },
        });
        return data?.uploadImage as Image; // Assuming the mutation returns the uploaded image ID
      } catch (error) {
        console.error("Error uploading image:", error);
        throw error;
      }
    },
    [uploadImageMutation],
  );

  const getImage = useCallback(
    async (docId: string, imageId: string): Promise<Image | undefined> => {
      try {
        const { data } = await getImageQuery({
          variables: { docId, imageId },
          fetchPolicy: "network-only",
        });
        console.log("GetImage DATA", data);
        return data?.getImage;
      } catch (error) {
        console.error("Error getting image:", error);
        throw error;
      }
    },
    [getImageQuery],
  );

  const listDocumentImages = useCallback(
    async (docId: string): Promise<Image[] | undefined> => {
      try {
        const { data } = await listDocumentImagesQuery({
          variables: { docId },
        });
        return data?.listDocumentImages as Image[];
      } catch (error) {
        console.error("Error listing document images:", error);
        return [];
      }
    },
    [listDocumentImagesQuery],
  );

  const getImageSignedUrl = useCallback(
    async (docId: string, imageId: string): Promise<string | undefined> => {
      try {
        const { data } = await getImageSignedUrlQuery({
          variables: { docId, imageId },
        });
        return data?.getImageSignedUrl.url;
      } catch (error) {
        console.error("Error getting image signed URL:", error);
        throw error;
      }
    },
    [getImageSignedUrlQuery],
  );

  useEffect(() => {
    if (editor) {
      _setAddress(editor.address);
      _setEditorMode(editor.editorMode);
      editor.subscribe<string | null>("address", _setAddress);
      editor.subscribe<"diff" | "history" | "paste" | "edit">(
        "editorMode",
        _setEditorMode,
      );
    }

    return () => {
      if (editor) {
        editor.unsubscribe("address", _setAddress);
        editor.unsubscribe("editorMode", _setEditorMode);
      }
    };
  }, [editor]);

  const value = {
    editor,
    setCurrentEditor,
    address,
    editorMode,
    uploadImage,
    getImage,
    listDocumentImages,
    getImageSignedUrl,
  };

  return (
    <RogueEditorContext.Provider value={value}>
      {children}
    </RogueEditorContext.Provider>
  );
};

export const useRogueEditorContext = () => {
  const context = useContext(RogueEditorContext);
  if (context === undefined) {
    throw new Error(
      "useRogueEditor must be used within a RogueEditorContextProvider",
    );
  }
  return context;
};
