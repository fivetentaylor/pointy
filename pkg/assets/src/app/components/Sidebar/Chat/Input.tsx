import React, { useState, useEffect, useMemo } from "react";

import { Button } from "@/components/ui/button";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { MessageInput, MsgLlm } from "@/__generated__/graphql";
import TextareaAutosize from "react-textarea-autosize";
import { useLocation } from "react-router-dom";
import { analytics } from "@/lib/segment";
import { ASK_AI_SEND_MESSAGE } from "@/lib/events";
import { PostHogFeature } from "posthog-js/react";
import { useWsDisconnect } from "@/hooks/useWsDisconnect";
import useSessionStorageState from "use-session-storage-state";
import { TipBox } from "@/components/ui/TipBox";
import Attachments, { AttachmentType } from "./Attachments";
import { useChatContext } from "@/contexts/ChatContext";
import VoiceMode from "./VoiceMode";
import { useDocumentContext } from "@/contexts/DocumentContext";

type InputProps = {
  documentId: string;
  threadId: string;
  authorId: string;
  selectedHtml: string;
  loading: boolean;
  showDisconnectedMessage?: boolean;
  onClearSelection: () => void;
  onSendReviseMessage: (message: string) => Promise<any>;
  onSendChatMessage: (message: string) => Promise<any>;
  onCreateThread: () => Promise<any>;
  uploadAttachment: (file: File) => Promise<any>;
  activeAttachments: any[];
  setActiveAttachments: (value: any) => void;
  refetchMessages: () => void;
};

const Input = ({
  documentId,
  threadId,
  authorId,
  onSendReviseMessage,
  onSendChatMessage,
  onCreateThread,
  onClearSelection,
  selectedHtml,
  loading,
  showDisconnectedMessage = false,
  uploadAttachment,
  activeAttachments,
  setActiveAttachments,
  refetchMessages,
}: InputProps) => {
  const [message, setMessage, { removeItem: resetMessage }] =
    useSessionStorageState<string>("ask-ai-message", {
      defaultValue: "",
    });
  const [altKeyDown, setAltKeyDown] = useState(false);
  const [ctrlKeyDown, setCtrltKeyDown] = useState(false);
  const [activeSelection, setActiveSelection] = useState("");

  const isEnabled = useMemo(() => {
    return !loading && message.length > 0;
  }, [loading, message]);

  useEffect(() => {
    setActiveSelection("");
  }, [selectedHtml]);

  const sendMessage = async (type: "chat" | "revise") => {
    if (!isEnabled) {
      return;
    }

    analytics.track(ASK_AI_SEND_MESSAGE, {
      message,
      hasSelection: !!(selectedHtml && selectedHtml.length > 0),
      type,
    });

    await (type === "chat"
      ? onSendChatMessage(message)
      : onSendReviseMessage(message));
    resetMessage();
    focusInput();
  };

  const newTopic = async () => {
    await onCreateThread();
    focusInput();
  };

  const focusInput = () => {
    const textarea = document.querySelector("textarea");
    if (textarea) {
      textarea.focus();
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (!e.altKey && !e.ctrlKey && (e.shiftKey || e.key !== "Enter")) {
      return;
    }

    if (e.altKey) {
      setAltKeyDown(true);
    }

    if (e.ctrlKey) {
      setCtrltKeyDown(true);
    }

    if (e.key === "Enter") {
      if (e.altKey) {
        sendMessage("revise");
      } else if (e.ctrlKey) {
        newTopic();
      } else {
        sendMessage("chat");
      }
    }

    return false;
  };

  const handleKeyUp = (
    e: KeyboardEvent | React.KeyboardEvent<HTMLTextAreaElement>,
  ) => {
    if (e.key === "Alt") {
      setAltKeyDown(false);
    }
    if (e.key === "Control") {
      setCtrltKeyDown(false);
    }
  };

  useEffect(() => {
    window.addEventListener("keyup", handleKeyUp);
    return () => {
      window.removeEventListener("keyup", handleKeyUp);
    };
  }, []);

  return (
    <div className="flex flex-col items-center relative rounded-md shadow border border-border focus-within:border-primary bg-card">
      <PostHogFeature flag="show-ai-outage-banner" match={true}>
        <TipBox>
          Our AI Provider is experiencing an outage. Chat may be unreliable.
        </TipBox>
      </PostHogFeature>
      {showDisconnectedMessage && (
        <TipBox>You&apos;re disconnected. Attempting to reconnect...</TipBox>
      )}
      <div className="flex flex-wrap min-h-6 w-full px-3 mt-3 items-center gap-2">
        <Attachments
          selectedHtml={selectedHtml}
          uploadAttachment={uploadAttachment}
          activeAttachments={activeAttachments}
          setActiveAttachments={setActiveAttachments}
          setActiveSelection={setActiveSelection}
          onClearSelection={onClearSelection}
        />
      </div>
      {activeSelection && (
        <div className="flex w-full mt-2 px-3">
          <div
            className="flex-grow overflow-y-auto overflow-x-hidden text-sm font-normal border border-border rounded-md p-2 max-h-20 selected-text-preview"
            dangerouslySetInnerHTML={{ __html: activeSelection }}
          />
        </div>
      )}
      <div className="flex items-center relative w-full">
        <div className="flex-grow text-sm font-normal font-sans w-full pt-1 px-4">
          <TextareaAutosize
            rows={1}
            maxRows={4}
            className="w-full pt-[0.6875rem] mb-[0.6875rem] leading-[1.3125rem] bg-transparent outline-none resize-none placeholder:text-muted-foreground placeholder:font-sans placeholder:text-sm"
            placeholder="Ask anything..."
            value={message}
            onKeyDown={handleKeyDown}
            onKeyUp={handleKeyUp}
            onChange={(e) => setMessage(e.target.value)}
            disabled={loading}
          ></TextareaAutosize>
        </div>
      </div>
      <div className="flex w-full mt-[-0.4rem] mb-2 px-2 justify-end items-center">
        <div className="flex items-center flex-grow">
          <VoiceMode
            documentId={documentId}
            threadId={threadId}
            authorId={authorId}
            refreshMessages={refetchMessages}
          />
        </div>

        <div className="flex-shrink">
          {!ctrlKeyDown && (
            <Button
              variant="ghost"
              size="sm"
              className="mr-2"
              disabled={!isEnabled}
              onClick={() => sendMessage("revise")}
            >
              <span className="mr-2">⌥↵</span>
              Write
            </Button>
          )}
        </div>
        {!altKeyDown && !ctrlKeyDown && (
          <div className="flex-shrink">
            <Button
              className="bg-primary text-white hover:bg-primary/90 h-9 py-2 px-3"
              disabled={!isEnabled}
              onClick={() => sendMessage("chat")}
            >
              <span className="mr-2">⏎</span>
              Chat
            </Button>
          </div>
        )}
      </div>
    </div>
  );
};

type InputContainerProps = {
  loadingCreateThreadMessage: boolean;
  loadingCreateThread: boolean;
  createThreadMessage: (input: any) => any;
  createThread: () => any;
  uploadAttachment: (file: File) => any;
};

function InputContainer({
  loadingCreateThreadMessage,
  createThreadMessage,
  loadingCreateThread,
  createThread,
  uploadAttachment,
}: InputContainerProps) {
  const location = useLocation();
  const { editor } = useRogueEditorContext();
  const { draftId } = useDocumentContext();
  const { activeThreadID, refetchMessages } = useChatContext();
  const [selectedHtml, setSelectedHtml] = useState("");
  const [activeAttachments, setActiveAttachments] = useSessionStorageState<
    AttachmentType[]
  >(`ask-ai-attachments-${activeThreadID}`, {
    defaultValue: [],
  });
  const llm: MsgLlm = "CLAUDE";

  useEffect(() => {
    if (!editor) {
      return;
    }

    if (!editor.curRogueRange) {
      setSelectedHtml("");
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

  const { isDisconnected } = useWsDisconnect();

  const sendMessage = async (message: string, type: "chat" | "revise") => {
    const input: MessageInput = {
      content: message,
      authorId: editor?.authorId || "",
      allowDraftEdits: type === "revise",
      contentAddress: editor?.currentContentAddress() || "",
      llm,
    };

    const selection = editor?.aiMessageSelection();
    if (selection) {
      input.selection = selection;
    }

    if (activeAttachments.length > 0) {
      input.attachments = activeAttachments.map((a) => {
        return {
          id: a.id,
          type: a.type === "draft" ? "DRAFT" : "FILE",
          name: a.name,
          contentType: a.contentType,
        };
      });
    }

    editor?.clearSelection();

    return createThreadMessage(input);
  };

  const handleRemoveSelectedHtml = () => {
    if (!editor) {
      return;
    }
    editor.clearSelection();
  };

  const handleReload = () => {
    refetchMessages();
  };

  return (
    <Input
      documentId={draftId || ""}
      threadId={activeThreadID || ""}
      authorId={editor?.authorId || ""}
      loading={
        loadingCreateThread ||
        loadingCreateThreadMessage ||
        !editor ||
        !editor.authorId ||
        isDisconnected
      }
      refetchMessages={handleReload}
      showDisconnectedMessage={isDisconnected}
      onSendChatMessage={(message) => sendMessage(message, "chat")}
      onSendReviseMessage={(message) => sendMessage(message, "revise")}
      onCreateThread={createThread}
      selectedHtml={selectedHtml}
      onClearSelection={handleRemoveSelectedHtml}
      uploadAttachment={uploadAttachment}
      activeAttachments={activeAttachments}
      setActiveAttachments={setActiveAttachments}
    />
  );
}

export default InputContainer;
