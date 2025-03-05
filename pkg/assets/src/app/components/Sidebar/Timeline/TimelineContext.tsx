import { EditMessageResolutionSummaryDocument } from "@/__generated__/graphql";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { useErrorToast } from "@/hooks/useErrorToast";
import {
  CreateTimelineMessage,
  DeleteTimelineMessage,
  EditTimelineMessage,
  EditTimelineUpdateSummary,
  ForceTimelineUpdateSummary,
  GetDocumentTimeline,
  TimelineEventDeleted,
  TimelineEventInserted,
  TimelineEventUpdated,
  UpdateMessageResolution,
} from "@/queries/timeline";
import { useMutation, useQuery, useSubscription } from "@apollo/client";
import React, { createContext, useState, useContext, useEffect } from "react";
import { useParams } from "react-router-dom";

type TimelineContextState = ReturnType<typeof useSetupTimeline>;

type TimelieContextProviderProps = {
  children: React.ReactNode;
};

const TimelineContext = createContext<TimelineContextState | undefined>(
  undefined,
);

const useSetupTimeline = () => {
  const { draftId } = useParams();
  const [activeReplyId, setActiveReplyId] = useState<string | null>(null);
  const [isActiveReplyResolved, setIsActiveReplyResolved] =
    useState<boolean>(false);
  const [hasReplyInput, setHasReplyInput] = useState<boolean>(false);
  const [activeComments, setActiveComments] = useState<string[]>([]);
  const showErrorToast = useErrorToast();
  const [timelineFilter, _setTimelineFilter] =
    useState<TimelineEventFilter>("ALL");

  const { editor } = useRogueEditorContext();

  useEffect(() => {
    if (!editor) return;
    const onActiveCommentsChange = (comments: string[]) => {
      setActiveComments(comments);
    };
    editor.subscribe("activeComments", onActiveCommentsChange);

    return () => {
      editor.unsubscribe("activeComments", onActiveCommentsChange);
    };
  }, [editor]);

  const {
    data: timelineData,
    loading: loadingTimeline,
    error: errorTimeline,
    refetch: refetchTimeline,
    subscribeToMore: subscribeToMoreTimeline,
  } = useQuery(GetDocumentTimeline, {
    variables: {
      documentId: draftId || "",
      filter: timelineFilter,
    },
    skip: !draftId,
    fetchPolicy: "network-only",
  });

  const timelineEvents =
    timelineData?.getDocumentTimeline || ([] as TimelineEventFieldsFragment[]);

  useEffect(() => {
    const timelineEvent = (
      timelineEvents as TimelineEventFieldsFragment[]
    ).find((event) => event.id === activeReplyId);

    if (timelineEvent) {
      const event = timelineEvent.event as TlMessageV1;
      const resolutions: TlMessageResolutionV1[] = [];
      const messageReplies: TlMessageV1[] = [];

      if (event.replies) {
        event.replies.forEach((reply) => {
          if (reply.event.__typename === "TLMessageResolutionV1") {
            const event: TlMessageResolutionV1 =
              reply.event as TlMessageResolutionV1;
            resolutions.push(event);
          } else {
            messageReplies.push(reply.event as TlMessageV1);
          }
        });
      }

      const isResolved =
        resolutions.length > 0 && resolutions[resolutions.length - 1]?.resolved;

      setIsActiveReplyResolved(isResolved);
    }
  }, [activeReplyId, timelineEvents]);

  const setTimelineFilter = (filter: TimelineEventFilter) => {
    if (filter === timelineFilter) {
      return;
    }
    _setTimelineFilter(filter);
    refetchTimeline();
  };

  useEffect(() => {
    const unsubscribe = subscribeToMoreTimeline({
      document: TimelineEventInserted,
      variables: {
        documentId: draftId || "",
        filter: timelineFilter,
      },
      updateQuery: (prev, { subscriptionData }) => {
        if (!subscriptionData.data) return prev;
        const newEvent = subscriptionData.data.timelineEventInserted;
        return {
          getDocumentTimeline: [...prev.getDocumentTimeline, newEvent],
        };
      },
    });

    return () => {
      unsubscribe();
    };
  }, [draftId, subscribeToMoreTimeline, timelineFilter]);

  const { error: timelineDeletedSubscriptionError } = useSubscription(
    TimelineEventDeleted,
    {
      variables: {
        documentId: draftId || "",
      },
      onData: (options) => {
        const cache = options.client.cache;
        const deletedEvent = options.data.data?.timelineEventDeleted;
        if (deletedEvent) {
          cache.evict({
            id: cache.identify({
              __typename: "TimelineEvent",
              id: deletedEvent.id,
            }),
          });
          cache.gc();
        }
      },
    },
  );

  const { error: timelineUpdatedSubscriptionError } = useSubscription(
    TimelineEventUpdated,
    {
      variables: {
        documentId: draftId || "",
      },
    },
  );

  const [
    createTimelineMessageMutation,
    { loading: loadingCreateTimelineMessage },
  ] = useMutation(CreateTimelineMessage, {
    onError: (error) => {
      showErrorToast(`Failed to create message: ${error.message}`);
    },
    variables: {
      documentId: draftId || "",
      input: {
        replyTo: null,
        content: "",
        contentAddress: "",
        authorId: "",
        startID: "",
        endID: "",
      },
    },
  });

  const createTimelineMessage = async (input: TimelineMessageInput) => {
    return createTimelineMessageMutation({
      variables: {
        input,
        documentId: draftId || "",
      },
    });
  };

  const [editTimelineMessage] = useMutation(EditTimelineMessage, {
    onError: (error) => {
      showErrorToast(`Failed to edit message: ${error.message}`);
    },
    variables: {
      documentId: draftId || "",
      messageId: "",
      input: {
        content: "",
      },
    },
  });

  const [updateMessageResolutionMutation] = useMutation(
    UpdateMessageResolution,
    {
      onError: (error) => {
        showErrorToast(`Failed to update message resolution: ${error.message}`);
      },
      variables: {
        documentId: draftId || "",
        messageId: "",
        input: {
          authorID: "",
          resolved: false,
        },
      },
    },
  );

  const updateMessageResolution = async (
    messageId: string,
    authorId: string,
    { resolved }: { resolved: boolean },
  ) => {
    return updateMessageResolutionMutation({
      variables: {
        documentId: draftId || "",
        messageId,
        input: {
          authorID: authorId,
          resolved,
        },
      },
    });
  };

  const [editTimelineUpdateSummaryMutation] = useMutation(
    EditTimelineUpdateSummary,
    {
      onError: (error) => {
        showErrorToast(`Failed to edit update summary: ${error.message}`);
      },
      variables: {
        documentId: draftId || "",
        updateId: "",
        summary: "",
      },
    },
  );

  const editTimelineUpdateSummary = async (
    updateId: string,
    summary: string,
  ) => {
    return editTimelineUpdateSummaryMutation({
      variables: {
        documentId: draftId || "",
        updateId,
        summary,
      },
    });
  };

  const [editMessageResolutionSummaryMutation] = useMutation(
    EditMessageResolutionSummaryDocument,
    {
      onError: (error) => {
        showErrorToast(
          `Failed to edit message resolution summary: ${error.message}`,
        );
      },
    },
  );

  const editMessageResolutionSummary = async (
    messageId: string,
    summary: string,
  ) => {
    return editMessageResolutionSummaryMutation({
      variables: {
        documentId: draftId || "",
        messageId,
        summary,
      },
    });
  };

  const [
    forceTimelineUpdateSummary,
    { loading: loadingForceTimelineUpdateSummary },
  ] = useMutation(ForceTimelineUpdateSummary, {
    onError: (error) => {
      showErrorToast(`Failed to create message: ${error.message}`);
    },
    variables: {
      documentId: draftId || "",
      userId: "",
    },
    refetchQueries: [
      {
        query: GetDocumentTimeline,
        variables: {
          documentId: draftId || "",
          filter: timelineFilter,
        },
      },
    ],
  });

  const [deleteTimelineMessage] = useMutation(DeleteTimelineMessage, {
    onError: (error) => {
      showErrorToast(`Failed to delete message: ${error.message}`);
    },
    refetchQueries: [
      {
        query: GetDocumentTimeline,
        variables: {
          documentId: draftId || "",
          filter: timelineFilter,
        },
      },
    ],
  });

  return {
    activeComments,
    activeReplyId,
    isActiveReplyResolved,
    setActiveReplyId,
    hasReplyInput,
    setHasReplyInput,
    createTimelineMessage,
    deleteTimelineMessage,
    editTimelineMessage,
    editTimelineUpdateSummary,
    editMessageResolutionSummary,
    errorTimeline,
    loadingCreateTimelineMessage,
    forceTimelineUpdateSummary,
    updateMessageResolution,
    loadingForceTimelineUpdateSummary,
    loadingTimeline,
    setTimelineFilter,
    setIsActiveReplyResolved,
    timelineEvents,
    timelineUpdatedSubscriptionError,
    timelineDeletedSubscriptionError,
    timelineFilter,
  };
};

// Create a provider component
export const TimelineProvider = ({ children }: TimelieContextProviderProps) => {
  const state = useSetupTimeline();

  return (
    <TimelineContext.Provider value={state}>
      {children}
    </TimelineContext.Provider>
  );
};

// Custom hook to use the TimelineContext
export const useTimelineContext = () => {
  const context = useContext(TimelineContext);
  if (context === undefined) {
    throw new Error(
      "useTimelineContext must be used within a TimelineProvider",
    );
  }
  return context;
};
