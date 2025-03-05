import React, { useRef, useState } from "react";
import { Button } from "@/components/ui/button";
import { ErrorBoundary } from "@/components/ui/ErrorBoundary";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { useChatContext } from "@/contexts/ChatContext";
import { Editor } from "@tiptap/react";
import MentionEditor from "@/components/ui/MentionEditor";
import { analytics } from "@/lib/segment";
import { AI_CLICK_ASK_REVISO } from "@/lib/events";
import { useSidebarContext } from "@/contexts/SidebarContext";

export const AskAIAnything = ({
  onSendMessage,
}: {
  container: HTMLElement | null;
  onSendMessage: (message: string, type: "chat" | "revise") => void;
}) => {
  const [message, setMessage] = useState("");
  const editorRef = useRef<Editor | null>(null);
  const isEnabled = message.length > 0;

  const sendMessage = (type: "chat" | "revise") => {
    onSendMessage(message, type);
    setMessage("");
    if (editorRef.current) {
      editorRef.current.commands.clearContent();
    }
  };

  return (
    <div className="relative flex items-center p-2 min-h-6">
      <div className="flex-grow mt-[-0.25rem] min-w-[19.3125rem] max-w-[16.3125rem]">
        <MentionEditor
          autoFocus
          mentionsEnabled={false}
          placeholder="Edit instructions"
          onChange={(value) => {
            setMessage(value);
          }}
          onEnter={() => {
            sendMessage("revise");
          }}
          onLoaded={(editor) => {
            editorRef.current = editor;
          }}
        />
      </div>
      <Button
        className="bg-primary text-white hover:bg-primary/90 h-8 py-2 px-2 text-xs leading-none mt-auto"
        disabled={!isEnabled}
        onClick={() => sendMessage("revise")}
      >
        <span className="mr-2">⏎</span>
        Edit
      </Button>
    </div>
  );
};

const AskAIAnythingWrapper = ({
  container,
  onSendMessage,
}: {
  container: HTMLElement | null;
  onSendMessage: () => void;
}) => {
  const { createThreadMessage } = useChatContext();
  const { setSidebarMode } = useSidebarContext();
  const { editor } = useRogueEditorContext();

  const handleSendMessage = (message: string, type: "chat" | "revise") => {
    const input: MessageInput = {
      content: message,
      authorId: editor?.authorId || "",
      allowDraftEdits: type === "revise",
      contentAddress: editor?.currentContentAddress() || "",
    };

    const selection = editor?.aiMessageSelection();
    if (selection) {
      input.selection = selection;
    }

    analytics.track(AI_CLICK_ASK_REVISO, {
      prompt: message,
    });

    editor?.clearSelection();

    createThreadMessage(input, editor?.currentContentAddress() || "");
    onSendMessage();
    setSidebarMode("chat");
  };

  return (
    <ErrorBoundary fallback={<div>Error</div>}>
      <AskAIAnything container={container} onSendMessage={handleSendMessage} />
    </ErrorBoundary>
  );
};

export default AskAIAnythingWrapper;
