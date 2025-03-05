import React, { useEffect, useRef, memo, useMemo, useState } from "react";
import CommentInput from "./CommentInput";
import { ScrollArea } from "@/components/ui/scroll-area";
import { TimelineEntry, TimelineSpacer } from "./TimelineEntry";
import { TimelineEventFieldsFragment } from "@/__generated__/graphql";
import TimelineActiveWriters from "./TimelineActiveWriters";
import { useTimelineContext } from "./TimelineContext";
import CommentThread, { TimelineEventWithMessage } from "./CommentThread";
import { ChevronLeftIcon, ChevronRightIcon } from "lucide-react";
import { cn } from "@/lib/utils";
import { ErrorBoundary } from "@/components/ui/ErrorBoundary";
import Header from "../Header";
import TimelineHeader from "./TimelineHeader";
import { BlockError } from "@/components/ui/BlockError";
import { useWsDisconnect } from "@/hooks/useWsDisconnect";

const MemoizedTimelineEntry = memo(TimelineEntry, (prevProps, nextProps) => {
  const isUnchanged = prevProps.timelineEvent === nextProps.timelineEvent;
  return isUnchanged;
});

interface DateGroupData {
  sortKey: string; // Full ISO date for accurate sorting
  displayDate: string; // Formatted date for display
  events: TimelineEventFieldsFragment[];
}

const DateGroup = memo(
  ({
    displayDate,
    events,
  }: {
    displayDate: string;
    events: TimelineEventFieldsFragment[];
  }) => {
    return (
      <div className="first:mt-auto">
        <TimelineSpacer date={displayDate} />
        {events.map((event) => (
          <MemoizedTimelineEntry key={event.id} timelineEvent={event} />
        ))}
      </div>
    );
  },
  (prevProps, nextProps) =>
    prevProps.displayDate === nextProps.displayDate &&
    prevProps.events === nextProps.events,
);

DateGroup.displayName = "DateGroup";

const CommentThreadNavigator = ({
  activeComments,
  onNavigate,
}: {
  activeComments: string[];
  onNavigate: (id: string) => void;
}) => {
  const [activeCommentIndex, setActiveCommentIndex] = useState(0);

  if (activeComments.length <= 1) return null;

  return (
    <div className="mt-auto flex items-center gap-[0.375rem] bg-card rounded-sm border border-border p-[0.375rem] max-w-fit">
      <ChevronLeftIcon
        className="w-4 h-4 cursor-pointer"
        onClick={(event) => {
          event.stopPropagation();
          const newIndex = Math.max(0, activeCommentIndex - 1);
          setActiveCommentIndex(newIndex);
          onNavigate(activeComments[newIndex]);
        }}
      />
      <span>
        {activeCommentIndex + 1} of {activeComments.length} threads
      </span>
      <ChevronRightIcon
        className="w-4 h-4 cursor-pointer"
        onClick={(event) => {
          event.stopPropagation();
          const newIndex = Math.min(
            activeComments.length - 1,
            activeCommentIndex + 1,
          );
          setActiveCommentIndex(newIndex);
          onNavigate(activeComments[newIndex]);
        }}
      />
    </div>
  );
};

interface TimelineProps {
  isDisconnected: boolean;
}
const Timeline: React.FC<TimelineProps> = ({ isDisconnected }) => {
  const {
    activeReplyId,
    activeComments,
    timelineEvents,
    loadingTimeline,
    createTimelineMessage,
    setActiveReplyId,
  } = useTimelineContext();

  const containerRef = useRef<HTMLDivElement | null>(null);
  const commentThreadRef = useRef<HTMLDivElement | null>(null);

  const activeReply = useMemo(() => {
    if (!activeReplyId) return null;
    return (timelineEvents as TimelineEventFieldsFragment[]).find(
      (event) => event?.id === activeReplyId,
    );
  }, [timelineEvents, activeReplyId]);

  const handleSelectCommentId = (commentId: string) => {
    const parentEvent = (timelineEvents as TimelineEventFieldsFragment[]).find(
      (event) => {
        if (event?.id === commentId) {
          console.log("parent event found", event.id);
          return true;
        } else if (
          event.event.__typename === "TLMessageV1" &&
          event.event.replies.some((reply) => reply.id === commentId)
        ) {
          return true;
        }
        return false;
      },
    );

    if (parentEvent) {
      setActiveReplyId(parentEvent.id);
    }
  };

  useEffect(() => {
    if (activeComments.length > 0) {
      handleSelectCommentId(activeComments[0]);
    } else {
      setActiveReplyId(null);
    }
  }, [activeComments]);

  const scrollToBottom = () => {
    if (containerRef.current) {
      const timelineContainer =
        containerRef.current.querySelector(".TimelineContainer");
      const scrollingContainer = containerRef.current.querySelector("&>div");
      scrollingContainer?.scrollTo({ top: timelineContainer?.scrollHeight });
    }
  };

  useEffect(() => {
    if (!loadingTimeline) {
      scrollToBottom();
    }
  }, [loadingTimeline]);

  useEffect(() => {
    if (containerRef.current) {
      const scrollingContainer = containerRef.current.querySelector("&>div");
      if (!scrollingContainer) return;

      // Type assertion to HTMLElement
      const scrollableElement = scrollingContainer as HTMLElement;

      // Calculate how close to the bottom we are
      const distanceToBottom =
        scrollableElement.scrollHeight -
        (scrollableElement.scrollTop + scrollableElement.clientHeight);

      if (distanceToBottom < 300) {
        scrollToBottom();
      }
    }
  }, [timelineEvents]);

  const sortedGroupedEvents = useMemo(() => {
    const grouped = groupEventsByDate(
      timelineEvents as TimelineEventFieldsFragment[],
    );
    return Object.values(grouped).sort((a, b) =>
      a.sortKey.localeCompare(b.sortKey),
    );
  }, [timelineEvents]);

  if (loadingTimeline || !timelineEvents) {
    return null;
  }

  return (
    <>
      <Header>
        <TimelineHeader />
      </Header>
      <ErrorBoundary
        fallback={
          <BlockError text="Your timeline couldn't be loaded due to an error." />
        }
      >
        <ScrollArea
          ref={containerRef}
          className="TimelineContainerScrollArea flex-1 overflow-y-auto w-auto mr-[-0.52rem] pr-1"
        >
          {activeReply && (
            <div
              className="absolute top-0 left-0 w-full h-full z-[49] flex flex-col bg-background/90 pt-8 pr-3"
              onClick={(event) => {
                // exit if click is anything but the thread
                if (!commentThreadRef.current?.contains(event.target as Node)) {
                  setActiveReplyId(null);
                }
              }}
            >
              <CommentThreadNavigator
                activeComments={activeComments}
                onNavigate={(commentId) => {
                  handleSelectCommentId(commentId);
                }}
              />
              <CommentThread
                ref={commentThreadRef}
                className={cn(
                  "mb-0",
                  activeComments.length > 1 ? "mt-2" : "mt-auto",
                )}
                timelineEvent={activeReply as TimelineEventWithMessage}
              />
            </div>
          )}
          <div
            className={cn(
              "TimelineContainer flex-grow flex flex-col h-full pr-2",
              isDisconnected ? "opacity-50 pointer-events-none" : "",
            )}
          >
            {sortedGroupedEvents.map(({ sortKey, displayDate, events }) => (
              <DateGroup
                key={sortKey}
                displayDate={displayDate}
                events={events}
              />
            ))}
            <TimelineActiveWriters onMount={scrollToBottom} />
          </div>
        </ScrollArea>
        <footer className="py-4 pb-0 pr-1">
          <CommentInput
            createTimelineMessage={createTimelineMessage}
            loadingCreateThreadMessage={false}
            isDisconnected={isDisconnected}
            onSendMessage={() => {
              scrollToBottom();
            }}
          />
        </footer>
      </ErrorBoundary>
    </>
  );
};

const groupEventsByDate = (
  events: TimelineEventFieldsFragment[],
): Record<string, DateGroupData> => {
  const addEventToAccumulator = (
    acc: Record<string, DateGroupData>,
    event: TimelineEventFieldsFragment,
  ) => {
    const createdAt = new Date(event.createdAt);
    const sortKey = createdAt.toISOString().split("T")[0]; // YYYY-MM-DD format
    const displayDate = createdAt.toLocaleDateString("en-US", {
      weekday: "short",
      month: "short",
      day: "numeric",
    });

    if (!acc[sortKey]) {
      acc[sortKey] = {
        sortKey,
        displayDate,
        events: [],
      };
    }
    acc[sortKey].events.push(event);
  };

  const reducer = (
    acc: Record<string, DateGroupData>,
    event: TimelineEventFieldsFragment,
  ) => {
    if (event.event.__typename === "TLMessageV1") {
      if (event.event.replies?.length > 0) {
        let isResolved = false;
        let lastResolution;
        const resolutions: TimelineEventFieldsFragment[] = [];
        const replies: TimelineEventFieldsFragment[] = [];

        event.event.replies.forEach((reply) => {
          if (reply.event.__typename === "TLMessageResolutionV1") {
            resolutions.push(reply as TimelineEventFieldsFragment);
          } else {
            replies.push(reply as TimelineEventFieldsFragment);
          }
        });

        if (resolutions.length > 0) {
          lastResolution = resolutions[resolutions.length - 1];
          isResolved = (lastResolution.event as TlMessageResolutionV1).resolved;
        }

        if (isResolved && lastResolution) {
          addEventToAccumulator(acc, lastResolution);
        } else {
          addEventToAccumulator(acc, event);
          replies.forEach((reply) => {
            addEventToAccumulator(acc, reply);
          });
        }
      } else {
        addEventToAccumulator(acc, event);
      }
    } else {
      addEventToAccumulator(acc, event);
    }

    return acc;
  };

  const eventHash = events.reduce(reducer, {} as Record<string, DateGroupData>);

  // Sort events within each date group
  Object.values(eventHash).forEach((dateEvents) => {
    dateEvents.events.sort(
      (a, b) =>
        new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime(),
    );
  });

  return eventHash;
};

const TimelineWrapper = (props: Omit<TimelineProps, "isDisconnected">) => {
  const { isDisconnected } = useWsDisconnect();
  return <Timeline isDisconnected={isDisconnected} {...props} />;
};

export default TimelineWrapper;
