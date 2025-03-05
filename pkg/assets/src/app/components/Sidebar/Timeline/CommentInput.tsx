import React, { useState, useEffect, useMemo, useRef } from "react";
import { XIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { useLocation } from "react-router-dom";
import { analytics } from "@/lib/segment";
import { CREATE_COMMENT } from "@/lib/events";
import { TimelineMessageInput } from "@/__generated__/graphql";
import { Editor } from "@tiptap/react";
import MentionEditor from "../../ui/MentionEditor";
import { useTimelineContext } from "./TimelineContext";
import { TipBox } from "@/components/ui/TipBox";
import useSessionStorageState from "use-session-storage-state";

export type CommentInputProps = {
  isReply: boolean;
  isResolved: boolean;
  onSendMessage: (message: string) => Promise<any>;
  onEditorLoaded?: (editor: Editor) => void;
  selectedHtml: string;
  loading: boolean;
  onClearSelection: () => void;
  isDisconnected: boolean;
};

export const CommentInput = ({
  isReply,
  isResolved,
  onSendMessage,
  onClearSelection,
  onEditorLoaded,
  selectedHtml,
  loading,
  isDisconnected,
}: CommentInputProps) => {
  const [message, setMessage, { removeItem: resetMessage }] =
    useSessionStorageState<string>("timeline-comment-message", {
      defaultValue: "",
    });
  const editorRef = useRef<Editor | null>(null);

  const isEnabled = useMemo(() => {
    return !loading && message?.trim().length > 0;
  }, [loading, message]);

  const handleEditorLoaded = (editor: Editor) => {
    editorRef.current = editor;
    if (onEditorLoaded) {
      onEditorLoaded(editor);
    }
  };

  const sendMessage = async () => {
    if (!isEnabled) {
      return;
    }

    analytics.track(CREATE_COMMENT, {
      message,
      hasSelection: !!(selectedHtml && selectedHtml.length > 0),
    });

    await onSendMessage(message);
    resetMessage();

    focusInput();
  };

  const focusInput = () => {
    const textarea = document.querySelector("textarea");
    if (textarea) {
      textarea.focus();
    }
  };

  let placeholder = isReply ? "Reply..." : "Leave a comment...";
  if (isResolved && isReply) {
    placeholder = "Reply to reopen thread...";
  }

  return (
    <div className="flex flex-col items-center relative rounded-md shadow border border-border focus-within:border-primary bg-card">
      {isDisconnected && (
        <TipBox>You&apos;re disconnected. Attempting to reconnect...</TipBox>
      )}
      {selectedHtml && (
        <div className="w-full">
          <div className="flex items-center rounded-md border border-border mx-2 mt-2 pl-4">
            <div
              className="flex-grow h-6 leading-6 overflow-hidden text-sm font-normal font-sans"
              dangerouslySetInnerHTML={{ __html: selectedHtml }}
            />
            <div className="flex-shrink pr-[0.375rem]">
              <Button
                variant="icon"
                size="icon"
                className="rounded-[1.625rem]"
                onClick={onClearSelection}
              >
                <XIcon className="w-4 h-4" />
              </Button>
            </div>
          </div>
        </div>
      )}
      <div className="flex items-center relative w-full">
        <div className="flex-grow w-full pt-3 pl-3 pr-4 mb-4">
          <MentionEditor
            disabled={loading}
            placeholder={placeholder}
            initialContent={message}
            onChange={(content) => {
              setMessage(content);
            }}
            onEnter={() => {
              sendMessage();
              if (editorRef.current) {
                editorRef.current.commands.clearContent();
              }
            }}
            onLoaded={handleEditorLoaded}
          />
        </div>
      </div>
      <div className="flex w-full mt-[-0.4rem] mb-2 px-2 items-end justify-end">
        <div className="flex-shrink">
          <Button
            className="bg-primary text-white hover:bg-primary/90 h-9 py-2 px-3"
            disabled={!isEnabled}
            onClick={() => {
              sendMessage();
              setMessage("");
              if (editorRef.current) {
                editorRef.current.commands.clearContent();
              }
            }}
          >
            <span className="mr-2">⏎</span>
            Comment
          </Button>
        </div>
      </div>
    </div>
  );
};

export type CommentInputContainerProps = {
  loadingCreateThreadMessage: boolean;
  createTimelineMessage: (input: any) => any;
  isDisconnected: boolean;
  onSendMessage: () => void;
};

function CommentInputContainer({
  loadingCreateThreadMessage,
  createTimelineMessage,
  isDisconnected,
  onSendMessage,
}: CommentInputContainerProps) {
  const location = useLocation();
  const { editor } = useRogueEditorContext();
  const [selectedHtml, setSelectedHtml] = useState("");
  const { activeReplyId, isActiveReplyResolved, updateMessageResolution } =
    useTimelineContext();
  const commentEditorRef = useRef<Editor | null>(null);

  useEffect(() => {
    if (!editor) {
      return;
    }

    const onSelectionChange: (value: string | null) => void = (value) => {
      if (!value) {
        setSelectedHtml("");
      }

      if (value) {
        const parser = new DOMParser();
        const doc = parser.parseFromString(value, "text/html");

        // Select all elements with the data-rid attribute and remove it
        doc.querySelectorAll("[data-rid]").forEach((element) => {
          element.removeAttribute("data-rid");
        });

        setSelectedHtml(doc.body.innerHTML);
      }
    };

    editor.subscribe<string | null>("selectedHtml", onSelectionChange);

    return () => {
      editor.unsubscribe("selectedHtml", onSelectionChange);
    };
  }, [editor]);

  useEffect(() => {
    handleRemoveSelectedHtml();
  }, [location.pathname]);

  useEffect(() => {
    if (commentEditorRef.current && activeReplyId) {
      commentEditorRef.current.commands.focus();
    }
  }, [activeReplyId]);

  const sendMessage = async (message: string) => {
    if (isActiveReplyResolved && activeReplyId) {
      // we need to reopen the thread
      updateMessageResolution(activeReplyId, editor?.authorId || "", {
        resolved: false,
      });
    }

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

    if (activeReplyId) {
      input.replyTo = activeReplyId;
    }

    const selection = editor?.aiMessageSelection();
    if (selection) {
      input.startID = selection.start;
      input.endID = selection.end;
      input.selectionMarkdown = selection.content;
    }

    editor?.clearCurrentRogueRange();

    createTimelineMessage(input);

    onSendMessage();
  };

  const handleRemoveSelectedHtml = () => {
    if (!editor) {
      return;
    }
    editor.clearCurrentRogueRange();
  };

  return (
    <CommentInput
      isReply={!!activeReplyId}
      isResolved={isActiveReplyResolved}
      loading={loadingCreateThreadMessage || isDisconnected}
      onSendMessage={sendMessage}
      selectedHtml={selectedHtml}
      onClearSelection={handleRemoveSelectedHtml}
      onEditorLoaded={(editor) => (commentEditorRef.current = editor)}
      isDisconnected={isDisconnected}
    />
  );
}

export default CommentInputContainer;
