import React from "react";

import Header from "../Header";
import Selector from "./Selector";
import List from "./List";
import Input from "./Input";
import { useChatContext } from "@/contexts/ChatContext";
import ChatHeader from "./ChatHeader";
import { ErrorBoundary } from "@/components/ui/ErrorBoundary";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { useErrorToast } from "@/hooks/useErrorToast";
import { BlockError } from "@/components/ui/BlockError";
import { analytics } from "@/lib/segment";
import { ASK_AI_CHANGE_THREAD } from "@/lib/events";

export type SidebarMode = "timeline" | "chat";

const Chat = () => {
  const { editor } = useRogueEditorContext();
  const showErrorToast = useErrorToast();
  const {
    activeThreadID,
    createThreadMutation,
    createThreadMessage: _createThreadMessage,
    errorThreads,
    errorMsgs,
    loadingThreads,
    loadingCreateThread,
    loadingMsgs,
    loadingCreateThreadMessage,
    messages,
    setActiveThreadID,
    threads,
    uploadAttachment,
  } = useChatContext();

  const handleSelectThread = (threadId: string) => {
    analytics.track(ASK_AI_CHANGE_THREAD);
    setActiveThreadID(threadId);
  };

  const handleCreateThread = async () => {
    await createThreadMutation();
  };

  if (errorThreads) {
    console.error("Error loading threads", errorThreads);
    showErrorToast("Error loading message thread");
  }

  if (errorMsgs) {
    console.error("Error loading messages", errorMsgs);
    showErrorToast("Error loading messages");
  }

  if (errorThreads || errorMsgs) {
    console.error("Error loading sidebar", {
      errorThreads,
      errorMsgs,
    });
    return <div>Error</div>;
  }

  if (!threads || !activeThreadID) {
    return <div>Loading...</div>;
  }

  return (
    <>
      <Header>
        <ChatHeader
          createThread={handleCreateThread}
          loadingCreateThread={loadingCreateThread}
        />
      </Header>
      <ErrorBoundary
        fallback={
          <BlockError text="The chat couldn't be loaded due to an error." />
        }
      >
        <Selector
          threadId={activeThreadID}
          onSelectThread={handleSelectThread}
          loading={loadingThreads}
          threads={threads}
        />
        <List
          messages={(messages || []) as MessageFieldsFragment[]}
          loading={loadingMsgs}
        />
        <footer className="pt-4 pb-0 min-w-[22rem] pr-1">
          <Input
            createThreadMessage={(input: MessageInput) => {
              _createThreadMessage(
                input,
                editor?.currentContentAddress() || "",
              );
            }}
            createThread={createThreadMutation}
            loadingCreateThreadMessage={loadingCreateThreadMessage}
            loadingCreateThread={loadingCreateThread}
            uploadAttachment={uploadAttachment}
          />
        </footer>
      </ErrorBoundary>
    </>
  );
};

export default Chat;
