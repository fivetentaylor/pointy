import React, { useState, useEffect, useRef } from "react";
import { Button } from "@/components/ui/button";
import { SendIcon } from "lucide-react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { cn } from "@/lib/utils";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { useCurrentUserContext } from "@/contexts/CurrentUserContext";
import { getInitials } from "@/lib/utils";
import MentionEditor from "@/components/ui/MentionEditor";
import { Editor } from "@tiptap/react";
import { useTimelineContext } from "./TimelineContext";

export const ReplyInput = function ({
  className,
  showReplyInput,
  setShowReplyInput,
  parentEventId,
  isChildReply = false,
  onMessageUpdated,
}: {
  parentEventId: string;
  className?: string;
  showReplyInput: boolean;
  isChildReply?: boolean;
  setShowReplyInput: (showReplyInput: boolean) => void;
  onMessageUpdated: (message: string) => void;
}) {
  const { currentUser } = useCurrentUserContext();

  const [message, setMessage] = useState("");

  const { editor } = useRogueEditorContext();
  const { createTimelineMessage } = useTimelineContext();
  const editorRef = useRef<Editor>();

  useEffect(() => {
    onMessageUpdated(message);
  }, [message]);

  const handleSendMessage = () => {
    const address = editor?.currentContentAddress();
    if (!address) {
      console.error("No current address");
      return;
    }

    const input: TimelineMessageInput = {
      replyTo: parentEventId,
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
    setMessage("");
    setTimeout(() => {
      setShowReplyInput(false);
    }, 0);
  };

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        event.target instanceof Element &&
        !event.target.closest(".reply-input")
      ) {
        setShowReplyInput(false);
      }
    };

    window.addEventListener("mousedown", handleClickOutside);

    return () => {
      window.removeEventListener("mousedown", handleClickOutside);
    };
  }, [setShowReplyInput]);

  return showReplyInput ? (
    <div
      className={cn(
        "reply-input",
        "relative rounded-md shadow border border-border focus-within:border-primary mt-3 min-h-14",
        isChildReply ? "ml-[-2.25rem]" : "",
        className,
      )}
    >
      <div className="flex items-start p-2">
        <Avatar className="w-6 h-6 mr-2">
          <AvatarImage
            alt="Profile icon"
            src={currentUser?.picture || undefined}
          />
          <AvatarFallback className="text-background bg-primary">
            {getInitials(currentUser?.name || "")}
          </AvatarFallback>
        </Avatar>
        <div className="flex-grow">
          <MentionEditor
            placeholder="Reply..."
            onChange={(value) => {
              setMessage(value);
            }}
            onEnter={() => {
              handleSendMessage();
              if (editorRef.current) {
                editorRef.current.commands.clearContent();
              }
            }}
            onLoaded={(editor) => {
              editorRef.current = editor;
            }}
          />
        </div>
      </div>
      <div className="absolute right-2 bottom-0">
        <Button
          variant="icon"
          size="icon"
          className="rounded-[26px]"
          onClick={handleSendMessage}
          disabled={message.length === 0}
        >
          <SendIcon
            className={`w-4 h-4 ${
              message.length === 0 ? "text-muted-foreground" : "text-primary"
            }`}
          />
        </Button>
      </div>
    </div>
  ) : (
    <Button
      variant="link"
      className="p-0 text-xs text-foreground m-0 mt-[-1rem] h-4"
      onClick={() => setShowReplyInput(true)}
    >
      Reply
    </Button>
  );
};
