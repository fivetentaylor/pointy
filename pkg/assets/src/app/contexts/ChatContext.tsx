import React, {
  useEffect,
  createContext,
  useContext,
  useMemo,
  useState,
  useRef,
} from "react";
import {
  CreateAIThread,
  CreateAIThreadMessage,
  GetAIThreadMessages,
  GetAIThreads,
  MessageUpserted,
  ThreadUpserted,
  UpdateMessageRevisionStatus,
} from "@/queries/messaging";
import {
  FetchResult,
  useMutation,
  useQuery,
  useSubscription,
} from "@apollo/client";
import { useParams } from "react-router-dom";
import { useErrorToast } from "@/hooks/useErrorToast";
import { useDebounce } from "use-debounce";
import { UploadAttachment } from "@/queries/attachments";

type ChatContextState = ReturnType<typeof useSetupChat>;

type ChatContextProviderProps = {
  children: React.ReactNode;
};

const ChatContext = createContext<ChatContextState | undefined>(undefined);

export const useSetupChat = () => {
  const { draftId } = useParams();
  const [debouncedDraftId] = useDebounce(draftId, 100);
  const showErrorToast = useErrorToast();
  const [activeThreadID, setActiveThreadID] = useState<string | null>(null);
  const isSwitchingDocs = useRef(false);

  const {
    data: threadsData,
    loading: loadingThreads,
    error: errorThreads,
  } = useQuery(GetAIThreads, {
    variables: {
      documentId: debouncedDraftId || "",
    },
    skip: !debouncedDraftId,
  });

  const threads = threadsData?.getAskAiThreads;

  useEffect(() => {
    if (draftId) {
      isSwitchingDocs.current = true;
      setActiveThreadID(null);
    }
  }, [draftId]);

  // Effect to set activeThreadID after threads are fetched
  useEffect(() => {
    if (threads?.[0]?.id) {
      setActiveThreadID(threads[0].id);
    } else if (isSwitchingDocs.current && (!threads || threads.length === 0)) {
      setActiveThreadID(null);
    }
    isSwitchingDocs.current = false;
  }, [threads]);

  const {
    data: msgsData,
    loading: loadingMsgs,
    error: errorMsgs,
    refetch: refetchMessages,
  } = useQuery(GetAIThreadMessages, {
    variables: {
      documentId: debouncedDraftId || "",
      threadId: activeThreadID || "",
    },
    skip: !debouncedDraftId || !activeThreadID || isSwitchingDocs.current,
  });

  const messages = msgsData?.getAskAiThreadMessages;

  const { error: subscriptionError } = useSubscription(MessageUpserted, {
    variables: {
      documentId: debouncedDraftId || "",
      channelId: activeThreadID || "",
    },
    skip: !debouncedDraftId || !activeThreadID || isSwitchingDocs.current,
    onError: (error) => {
      showErrorToast(`Failed to subscribe to messages: ${error.message}`);
      console.error(error);
    },
  });

  const { error: threadSubscriptionError } = useSubscription(ThreadUpserted, {
    variables: {
      documentId: debouncedDraftId || "",
    },
    skip: !debouncedDraftId || isSwitchingDocs.current,
    onError: (error) => {
      showErrorToast(`Failed to subscribe to threads: ${error.message}`);
      console.error(error);
    },
  });

  const [createThreadMutation, { loading: loadingCreateThread }] = useMutation(
    CreateAIThread,
    {
      variables: {
        documentId: debouncedDraftId || "",
      },
      refetchQueries: [
        {
          query: GetAIThreads,
          variables: {
            documentId: debouncedDraftId || "",
          },
        },
      ],
    },
  );

  const [_markRevisionStatus] = useMutation(UpdateMessageRevisionStatus, {
    onError: (error) => {
      showErrorToast(`Failed to update message: ${error.message}`);
    },
  });

  const markRevisionStatus = async function (
    status: MessageRevisionStatus,
    contendAddress: string,
  ) {
    if (!messages) {
      return;
    }

    let lastRevisionMessage;
    for (let i = messages.length - 1; i >= 0; i--) {
      const currentMessage = messages[i] as MessageFieldsFragment;
      if (
        currentMessage.attachments.some(
          (attachment) => attachment.__typename === "Revision",
        )
      ) {
        lastRevisionMessage = currentMessage;
        break;
      }
    }

    console.log("MARK REVISION STATUS", lastRevisionMessage);

    if (lastRevisionMessage) {
      return _markRevisionStatus({
        variables: {
          containerId: lastRevisionMessage.containerId,
          messageId: lastRevisionMessage.id,
          status,
          contentAddress: contendAddress,
        },
      });
    }
  };

  const [createThreadMessageMutation, { loading: loadingCreateThreadMessage }] =
    useMutation(CreateAIThreadMessage, {
      onError: (error) => {
        showErrorToast(`Failed to create message: ${error.message}`);
      },
      variables: {
        documentId: debouncedDraftId || "",
        threadId: activeThreadID || "",
        input: {
          content: "",
          authorId: "",
          allowDraftEdits: true,
          contentAddress: "",
        },
      },
      refetchQueries: [
        {
          query: GetAIThreadMessages,
          variables: {
            documentId: debouncedDraftId || "",
            threadId: activeThreadID || "",
          },
        },
      ],
    });

  const createThreadMessage = async (
    input: MessageInput,
    currentContentAddress: string,
  ) => {
    const lastMessage =
      (messages?.[messages.length - 1] as MessageFieldsFragment) || undefined;
    if (
      lastMessage &&
      lastMessage.metadata.contentAddress &&
      lastMessage.metadata.revisionStatus === "UNSPECIFIED"
    ) {
      markRevisionStatus("ACCEPTED", currentContentAddress);
    }
    input.contentAddress = currentContentAddress;
    return createThreadMessageMutation({
      variables: {
        input,
        documentId: debouncedDraftId || "",
        threadId: activeThreadID || "",
      },
    });
  };

  const [
    _uploadAttachment,
    { loading: loadingAttachment, error: attachmentError },
  ] = useMutation(UploadAttachment, {
    onError: (error) => {
      showErrorToast(`Failed to upload attachment: ${error.message}`);
    },
  });

  const uploadAttachment = async (
    file: File,
  ): Promise<FetchResult<UploadAttachmentMutation>> => {
    const data = _uploadAttachment({
      variables: {
        file,
        docId: debouncedDraftId || "",
      },
    });

    return data;
  };

  const isRevising = useMemo(() => {
    const lastMessage =
      (messages?.[messages.length - 1] as MessageFieldsFragment) || undefined;
    return lastMessage && lastMessage.lifecycleStage === "REVISING";
  }, [messages]);

  return {
    activeThreadID,
    attachmentError,
    createThreadMessage,
    createThreadMutation,
    errorMsgs,
    errorThreads,
    isRevising,
    loadingAttachment,
    loadingCreateThread,
    loadingCreateThreadMessage,
    loadingMsgs,
    loadingThreads,
    markRevisionStatus,
    messages,
    refetchMessages,
    setActiveThreadID,
    subscriptionError,
    threadSubscriptionError,
    threads,
    uploadAttachment,
  };
};

export const ChatContextProvider = function ({
  children,
}: ChatContextProviderProps) {
  const state = useSetupChat();

  return <ChatContext.Provider value={state}>{children}</ChatContext.Provider>;
};

export const useChatContext = () => {
  const context = useContext(ChatContext);
  if (context === undefined) {
    throw new Error("useChatContext must be used within a ChatContextProvider");
  }
  return context;
};
