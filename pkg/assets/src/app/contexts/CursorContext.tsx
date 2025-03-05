import React, { useEffect, createContext, useContext, useState } from "react";
import { GetUsers } from "@/queries/user";
import { AuthorInfo } from "@/../rogueEditor";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { useLazyQuery } from "@apollo/client";

type CursorContextState = ReturnType<typeof useSetupCursor>;

type CursorContextProviderProps = {
  children: React.ReactNode;
};

const CursorContext = createContext<CursorContextState | undefined>(undefined);

export const useSetupCursor = () => {
  const { editor } = useRogueEditorContext();
  const [editing, _setEditing] = useState(false);
  const [cursors, _setCursors] = useState<AuthorInfo[]>([]);
  const [connectedUsers, _setConnectedUsers] = useState<User[]>([]);

  const [getUsers] = useLazyQuery(GetUsers);

  useEffect(() => {
    if (!editor) {
      return;
    }

    const setCursors = (records: Record<string, AuthorInfo>) => {
      _setCursors(Object.values(records));

      const userIds = Object.values(records)
        .map((authorInfo) => authorInfo.userID)
        .filter((userId) => userId !== "");

      if (userIds.length > 0) {
        getUsers({ variables: { ids: userIds } }).then(({ data }) => {
          if (data) {
            _setConnectedUsers(data.users as User[]);
          }
        });
      }
    };

    const setEditing = (editing: boolean) => {
      _setEditing(editing);
    };

    editor.subscribe<Record<string, AuthorInfo>>("cursors", setCursors);
    editor.subscribe<boolean>("editing", setEditing);

    return () => {
      editor.unsubscribe<Record<string, AuthorInfo>>("cursors", setCursors);
      editor.unsubscribe<boolean>("editing", setEditing);
      _setCursors([]);
      _setEditing(false);
      _setConnectedUsers([]);
    };
  }, [editor, getUsers]);

  return {
    editing,
    cursors,
    connectedUsers,
  };
};

export const CursorContextProvider = function ({
  children,
}: CursorContextProviderProps) {
  const state = useSetupCursor();

  return (
    <CursorContext.Provider value={state}>{children}</CursorContext.Provider>
  );
};

export const useCursorContext = () => {
  const context = useContext(CursorContext);
  if (context === undefined) {
    throw new Error(
      "useCursorContext must be used within a CursorContextProvider",
    );
  }
  return context;
};
