import React, { useEffect, useRef, useState } from "react";
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/avatar";
import { useCurrentUserContext } from "@/contexts/CurrentUserContext";
import { getInitials } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { ErrorBoundary } from "@/components/ui/ErrorBoundary";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { Editor } from "@tiptap/react";
import MentionEditor from "@/components/ui/MentionEditor";
import { useTimelineContext } from "@/components/Sidebar/Timeline/TimelineContext";
import { useSidebarContext } from "@/contexts/SidebarContext";

export const CommentEditor = ({
  author,
  container,
  onSendMessage,
  onCommentUpdate,
}: {
  author: User | null;
  container: HTMLElement | null;
  onSendMessage: (message: string) => void;
  onCommentUpdate: (message: string) => void;
}) => {
  const [message, setMessage] = useState("");
  const editorRef = useRef<Editor | null>(null);

  useEffect(() => {
    onCommentUpdate(message.trim());
  }, [message]);

  return (
    <div className="relative flex items-start p-2">
      <Avatar className="w-6 h-6 mr-2 flex-shrink-0 mt-1">
        <AvatarImage alt="Profile icon" src={author?.picture || undefined} />
        <AvatarFallback className="text-background bg-primary">
          {getInitials(author?.name || "")}
        </AvatarFallback>
      </Avatar>
      <div className="flex flex-col flex-grow">
        <div className="flex items-center min-h-6">
          <div className="flex-grow mt-[-0.25rem] min-w-[16.3125rem] max-w-[16.3125rem]">
            <MentionEditor
              autoFocus
              mentionContainer={container}
              placeholder="Leave comment..."
              onChange={(value) => {
                setMessage(value);
              }}
              onEnter={() => {
                onSendMessage(message);
                setMessage("");
                if (editorRef.current) {
                  editorRef.current.commands.clearContent();
                }
              }}
              onLoaded={(editor) => {
                editorRef.current = editor;
              }}
            />
          </div>
          <Button
            className="bg-primary text-white hover:bg-primary/90 h-8 py-2 px-2 text-xs leading-none mt-auto"
            onClick={() => {
              onSendMessage(message);
              setMessage("");
              if (editorRef.current) {
                editorRef.current.commands.clearContent();
              }
            }}
            disabled={message.length === 0}
          >
            <span className="mr-2">⏎</span>
            Comment
          </Button>
        </div>
      </div>
    </div>
  );
};

const CommentEditorWrapper = ({
  container,
  onCreateTimelineMessage,
  onCommentUpdate,
}: {
  container: HTMLElement | null;
  onCreateTimelineMessage: () => void;
  onCommentUpdate: (message: string) => void;
}) => {
  const { currentUser } = useCurrentUserContext();
  const { createTimelineMessage } = useTimelineContext();
  const { setSidebarMode } = useSidebarContext();
  const { editor } = useRogueEditorContext();

  const handleSendMessage = (message: string) => {
    const address = editor?.currentContentAddress();
    if (!address) {
      console.error("No current address");
      return;
    }

    const input: TimelineMessageInput = {
      content: message,
      contentAddress: address,
      authorId: editor?.authorId || "",
    };

    const selection = editor?.aiMessageSelection();
    if (selection) {
      input.startID = selection.start;
      input.endID = selection.end;
      input.selectionMarkdown = selection.content;
    }

    createTimelineMessage(input);
    onCreateTimelineMessage();
    setSidebarMode("timeline");
  };

  return (
    <ErrorBoundary fallback={<div>Error</div>}>
      <CommentEditor
        author={currentUser}
        container={container}
        onSendMessage={handleSendMessage}
        onCommentUpdate={onCommentUpdate}
      />
    </ErrorBoundary>
  );
};

export default CommentEditorWrapper;
