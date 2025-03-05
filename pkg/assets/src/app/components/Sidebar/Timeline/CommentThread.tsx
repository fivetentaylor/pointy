import React, { useState, useEffect, useRef, forwardRef, useMemo } from "react";
import ReactMarkdown from "react-markdown";
import { CheckCircle2, MoreHorizontal } from "lucide-react";
import { Comment } from "./Comment";
import {
  TimelineEventFieldsFragment,
  TlMessageV1,
} from "@/__generated__/graphql";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { RidToRogueID, RogueEditor } from "../../../../rogueEditor";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { useTimelineContext } from "./TimelineContext";
import { TimelineAvatar, TimelineDate } from "./TimelineHelpers";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import MentionEditor from "@/components/ui/MentionEditor";
import { useParams } from "react-router-dom";
import { useCurrentUserContext } from "@/contexts/CurrentUserContext";
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogCancel,
  AlertDialogAction,
} from "@/components/ui/alert-dialog";
import { WithTooltip } from "@/components/ui/FloatingTooltip";
import { analytics } from "@/lib/segment";
import { REOPEN_COMMENT_THREAD, RESOLVE_COMMENT_THREAD } from "@/lib/events";

export type TimelineEventWithMessage = TimelineEventFieldsFragment & {
  event: TlMessageV1;
};

export type CommentThreadProps = {
  className?: string;
  timelineEvent: TimelineEventWithMessage;
  showExpando?: boolean;
};

type TimelineEventUser = TimelineEventFieldsFragment["user"];

const CommentHeader = function ({
  user,
  createdAt,
}: {
  user: TimelineEventUser;
  createdAt: string;
}) {
  return (
    <div className="flex items-center">
      <TimelineAvatar user={user} className="mt-0" />
      <span className="text-[0.875rem] leading-[1.125rem] font-medium">
        {user.name}
      </span>
      <span className="text-xs ml-2">
        <TimelineDate date={createdAt} />
      </span>
    </div>
  );
};

const CommentBody = ({
  content,
  selectionMarkdown,
}: {
  content: string;
  selectionMarkdown: string;
}) => {
  const hasSelection = !!selectionMarkdown;

  return (
    <div className="ml-6">
      <div className="flex items-center justify-end">
        {hasSelection && (
          <div className="flex-grow h-6 leading-6 overflow-hidden text-sm italic text-muted-foreground">
            <ReactMarkdown>{selectionMarkdown}</ReactMarkdown>
          </div>
        )}
      </div>
      <Comment content={content} />
    </div>
  );
};

const ReplyButton = ({
  onClick,
  isThread,
}: {
  onClick: () => void;
  isThread: boolean;
}) => {
  return (
    <div className="ml-6">
      <Button
        variant="link"
        className="p-0 text-xs text-primary m-0 mt-[-1rem] h-4"
        onClick={onClick}
      >
        {isThread ? "View thread & reply" : "Reply"}
      </Button>
    </div>
  );
};

const CommentEditor = ({
  availableWidth,
  content,
  onCancel,
  onSave,
  selectionMarkdown,
}: {
  availableWidth?: string;
  content: string;
  onCancel: () => void;
  onSave: (content: string) => void;
  selectionMarkdown?: string;
}) => {
  const [message, setMessage] = useState(content);

  const onSendMessage = (content: string) => {
    onSave(content);
  };

  return (
    <>
      <div className="flex items-center justify-end">
        {selectionMarkdown && (
          <div className="flex-grow h-6 leading-6 overflow-hidden text-sm italic text-muted-foreground mt-2 cursor-pointer">
            <ReactMarkdown>{selectionMarkdown}</ReactMarkdown>
          </div>
        )}
      </div>
      <div
        className="border border-ring rounded-md shadow-sm mt-2 mr-[-1.75rem] p-2 min-h-20 flex flex-col"
        style={availableWidth ? { width: availableWidth } : {}}
        onClick={(event) => {
          event.stopPropagation();
        }}
      >
        <MentionEditor
          autoFocus
          placeholder="Updated comment..."
          initialContent={content}
          onChange={(value) => {
            setMessage(value);
          }}
          onEnter={() => {
            onSendMessage(message);
          }}
          onLoaded={() => {}}
        />
        <div className="flex justify-end space-x-1 mt-2">
          <Button
            variant="ghost"
            className="h-9"
            onClick={(event) => {
              event.stopPropagation();
              onCancel();
            }}
          >
            Cancel
          </Button>
          <Button
            className="bg-primary hover:bg-primary/90 text-primary-foreground h-9"
            onClick={(event) => {
              event.stopPropagation();
              onSave(message);
            }}
          >
            Save
          </Button>
        </div>
      </div>
    </>
  );
};

const CommentDeleteAlert = ({
  open,
  onConfirm,
  onCancel,
  onOpenChange,
}: {
  open: boolean;
  onConfirm: () => void;
  onCancel: () => void;
  onOpenChange: (open: boolean) => void;
}) => {
  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Delete comment?</AlertDialogTitle>
          <AlertDialogDescription>
            Please confirm you want to delete your comment.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel onClick={onCancel}>Cancel</AlertDialogCancel>
          <AlertDialogAction
            onClick={onConfirm}
            className="bg-destructive hover:bg-destructive/80"
          >
            Delete
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
};

type SingleCommentProps = {
  className?: string;
  timelineEvent: TimelineEventWithMessage;
  commentRef?: React.RefObject<HTMLDivElement>;
  isResolved?: boolean;
  isRootComment?: boolean;
  isThread?: boolean;
  onClickComment: () => void;
  onClickReply?: () => void;
  onUpdateResolved: (resolved: boolean) => void;
};

const SingleComment = ({
  className,
  timelineEvent,
  commentRef,
  onClickComment,
  onClickReply,
  onUpdateResolved,
  isResolved = false,
  isRootComment = false,
  isThread = false,
}: SingleCommentProps) => {
  const {
    editTimelineMessage,
    deleteTimelineMessage,
    updateMessageResolution,
  } = useTimelineContext();
  const { draftId } = useParams();
  const { currentUser } = useCurrentUserContext();
  const { editor } = useRogueEditorContext();
  const [availableWidth, setAvailableWidth] = useState<string>("");
  const [showDeleteAlert, setShowDeleteAlert] = useState(false);
  const commentContainerRef = useRef<HTMLDivElement>(null);

  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const [isEditing, setIsEditing] = useState(false);

  const hasSelection = !!timelineEvent.event.selectionStartId;
  const isOwner = currentUser?.id === timelineEvent.user.id;

  useEffect(() => {
    if (!commentContainerRef || !commentContainerRef.current) {
      return;
    }
    const updateAvailableWidth = () => {
      if (commentContainerRef && commentContainerRef.current) {
        const { width } = commentContainerRef.current.getBoundingClientRect();
        setAvailableWidth(`${width}px`);
      }
    };

    updateAvailableWidth();

    const resizeObserver = new ResizeObserver((entries) => {
      for (const entry of entries) {
        if (entry.contentBoxSize) {
          updateAvailableWidth();
        }
      }
    });

    if (commentContainerRef.current) {
      resizeObserver.observe(commentContainerRef.current);
    }

    return () => {
      if (commentContainerRef.current) {
        resizeObserver.unobserve(commentContainerRef.current);
      }
    };
  }, [commentContainerRef, commentContainerRef?.current]);

  const onEditMessage = (content: string) => {
    setIsEditing(false);
    editTimelineMessage({
      variables: {
        documentId: draftId!,
        messageId: timelineEvent.event.eventId,
        input: { content },
      },
    });
  };

  const onDeleteMessage = () => {
    deleteTimelineMessage({
      variables: {
        documentId: draftId!,
        messageId: timelineEvent.event.eventId,
      },
    });
  };

  const onClickResolveButton = () => {
    analytics.track(RESOLVE_COMMENT_THREAD);
    onUpdateResolved(true);
    updateMessageResolution(
      timelineEvent.event.eventId,
      editor?.authorId || "",
      { resolved: true },
    );
  };

  const onClickReopenButton = () => {
    analytics.track(REOPEN_COMMENT_THREAD);
    onUpdateResolved(false);
    updateMessageResolution(
      timelineEvent.event.eventId,
      editor?.authorId || "",
      { resolved: false },
    );
  };

  const showDropdown = isOwner || (isRootComment && isResolved);

  return (
    <>
      <CommentDeleteAlert
        open={showDeleteAlert}
        onOpenChange={setShowDeleteAlert}
        onConfirm={onDeleteMessage}
        onCancel={() => setShowDeleteAlert(false)}
      />
      <div
        ref={commentContainerRef}
        className={cn(
          "flex flex-col justify-center relative group",
          !isEditing && hasSelection && "cursor-pointer",
          className,
        )}
        onClick={onClickComment}
      >
        <div ref={commentRef} className="flex justify-between items-start">
          <div className="flex-grow">
            <CommentHeader
              user={timelineEvent.user}
              createdAt={timelineEvent.createdAt}
            />
            {isEditing ? (
              <CommentEditor
                availableWidth={availableWidth}
                content={timelineEvent.event.content}
                selectionMarkdown={
                  hasSelection
                    ? timelineEvent.event.selectionMarkdown
                    : undefined
                }
                onCancel={() => {
                  setIsEditing(false);
                }}
                onSave={onEditMessage}
              />
            ) : (
              <CommentBody
                content={timelineEvent.event.content}
                selectionMarkdown={timelineEvent.event.selectionMarkdown}
              />
            )}
          </div>
          <div
            className={cn(
              "flex items-center transition-all duration-200 ease-in-out text-muted-foreground",
              isDropdownOpen
                ? "opacity-100 visible"
                : "opacity-0 invisible group-hover:opacity-100 group-hover:visible",
              isEditing ? "opacity-0 invisible" : "",
            )}
          >
            {/* future emoji button
          <Button variant="ghost" size="sm" className="h-7 w-7 p-0">
            <Smile className="h-4 w-4" />
          </Button>
          */}
            {!isResolved && isRootComment && (
              <WithTooltip tooltipText="Resolve thread">
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-7 w-7 p-0"
                  onClick={onClickResolveButton}
                >
                  <CheckCircle2 className="h-4 w-4" />
                </Button>
              </WithTooltip>
            )}
            {showDropdown && (
              <DropdownMenu onOpenChange={setIsDropdownOpen} modal={false}>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" size="sm" className="h-7 w-7 p-0">
                    <MoreHorizontal className="h-4 w-4" />
                    <span className="sr-only">Open menu</span>
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent
                  className="z-[51]"
                  align="end"
                  onCloseAutoFocus={(event) => {
                    event.stopPropagation();
                    event.preventDefault();
                  }}
                >
                  {isResolved && isRootComment && (
                    <DropdownMenuItem
                      onClick={(event) => {
                        event.stopPropagation();
                        onClickReopenButton();
                      }}
                    >
                      Reopen
                    </DropdownMenuItem>
                  )}
                  {isOwner && (
                    <>
                      <DropdownMenuItem
                        onClick={(event) => {
                          event.stopPropagation();
                          setIsEditing(true);
                        }}
                      >
                        Edit
                      </DropdownMenuItem>
                      <DropdownMenuItem
                        onClick={(event) => {
                          event.stopPropagation();
                          setShowDeleteAlert(true);
                        }}
                      >
                        Delete
                      </DropdownMenuItem>
                    </>
                  )}
                </DropdownMenuContent>
              </DropdownMenu>
            )}
          </div>
        </div>
      </div>
      {onClickReply && !isEditing && (
        <ReplyButton onClick={onClickReply} isThread={isThread} />
      )}
    </>
  );
};

export const TLMessageEntry = ({
  timelineEvent,
}: {
  timelineEvent: TimelineEventWithMessage;
}) => {
  const { editor } = useRogueEditorContext();
  const { setActiveReplyId } = useTimelineContext();

  const hasSelection = !!timelineEvent.event.selectionStartId;

  const showHighlights = function (editor: RogueEditor, event: TlMessageV1) {
    editor.showCommentHighlights([
      {
        eventId: event.eventId,
        address: event.contentAddress,
        startId: RidToRogueID(event.selectionStartId),
        endId: RidToRogueID(event.selectionEndId),
      },
    ]);
  };

  useEffect(() => {
    if (!editor) return;

    const updateConnectionState = (connected: boolean) => {
      if (connected && hasSelection) {
        showHighlights(editor, timelineEvent.event);
      }
    };

    editor.subscribe<boolean>("connected", updateConnectionState);

    updateConnectionState(editor.connected);

    return () => {
      editor.unsubscribe("connected", updateConnectionState);
    };
  }, [editor]);

  const updateHighlightsForResolution = (resolved: boolean) => {
    if (!editor) {
      return;
    }
    if (!resolved) {
      showHighlights(editor, timelineEvent.event);
    } else {
      editor.hideCommentHighlights([timelineEvent.event.eventId]);
      editor.hideActiveCommentHighlight(timelineEvent.event.eventId);
    }
  };

  return (
    <div
      className="mb-4 bg-card border border-border p-4 shadow-sm rounded-sm"
      onMouseOver={() => {
        if (editor && hasSelection) {
          editor.showActiveCommentHighlight({
            eventId: timelineEvent.event.eventId,
            startId: RidToRogueID(timelineEvent.event.selectionStartId),
            endId: RidToRogueID(timelineEvent.event.selectionEndId),
            address: timelineEvent.event.contentAddress,
          });
        }
      }}
      onMouseLeave={() => {
        if (editor) {
          editor.hideActiveCommentHighlight(timelineEvent.event.eventId);
        }
      }}
    >
      <SingleComment
        timelineEvent={timelineEvent}
        isResolved={false}
        isRootComment={
          timelineEvent.replyTo === "" || timelineEvent.replyTo === null
        }
        isThread={
          !!timelineEvent.replyTo || timelineEvent.event.replies.length > 0
        }
        onClickComment={() => {
          if (hasSelection && editor) {
            editor.activateComment({
              eventId: timelineEvent.event.eventId,
            });
          }
        }}
        onClickReply={() => {
          if (timelineEvent.replyTo) {
            setActiveReplyId(timelineEvent.replyTo);
          } else {
            setActiveReplyId(timelineEvent.id);
          }
        }}
        onUpdateResolved={updateHighlightsForResolution}
      />
    </div>
  );
};

export const CommentThread = forwardRef<HTMLDivElement, CommentThreadProps>(
  ({ className, timelineEvent }, ref) => {
    const commentRef = useRef<HTMLDivElement>(null);
    const { editor } = useRogueEditorContext();
    const { isActiveReplyResolved } = useTimelineContext();
    const localResolution = useRef(isActiveReplyResolved);

    const commentsToHighlight = useMemo(() => {
      const comments: TlMessageV1[] = [];
      if (
        timelineEvent.event.__typename === "TLMessageV1" &&
        !!timelineEvent.event.selectionStartId
      ) {
        comments.push(timelineEvent.event);
      }
      comments.push(
        ...timelineEvent.event.replies
          .filter(
            (reply): reply is TimelineEventWithMessage =>
              reply.event.__typename === "TLMessageV1" &&
              !!reply.event.selectionStartId,
          )
          .map((reply) => reply.event),
      );

      return comments;
    }, [timelineEvent]);

    const updateHighlightsForResolution = (resolved: boolean) => {
      if (!editor) {
        return;
      }
      if (!resolved) {
        editor.hideCommentHighlights(
          commentsToHighlight.map((comment) => comment.eventId),
        );
        editor.hideActiveCommentHighlight(timelineEvent.event.eventId);
      }
      localResolution.current = resolved;
    };

    const hasSelection = commentsToHighlight.length > 0;
    const messages = [timelineEvent, ...timelineEvent.event.replies];

    useEffect(() => {
      const cleanup = () => {
        if (editor && localResolution.current) {
          commentsToHighlight.forEach((comment) => {
            editor.hideActiveCommentHighlight(comment.eventId);
          });
          editor.hideCommentHighlights(
            commentsToHighlight.map((comment) => comment.eventId),
          );
        }
      };

      return cleanup;
    }, []);

    return (
      <div
        ref={ref}
        className={cn(
          "mb-4 bg-card border border-border p-4 shadow-sm rounded-sm",
          className,
        )}
        onMouseOver={() => {
          if (editor) {
            commentsToHighlight.forEach((comment) => {
              editor.showActiveCommentHighlight({
                eventId: comment.eventId,
                startId: RidToRogueID(comment.selectionStartId),
                endId: RidToRogueID(comment.selectionEndId),
                address: comment.contentAddress,
              });
            });
          }
        }}
        onMouseLeave={() => {
          if (editor) {
            commentsToHighlight.forEach((comment) => {
              editor.hideActiveCommentHighlight(comment.eventId);
            });
          }
        }}
      >
        {messages.map((event, idx) => {
          if (event.event.__typename === "TLMessageV1") {
            return (
              <SingleComment
                key={event.id}
                className={idx > 0 ? "mt-2" : ""}
                timelineEvent={event as TimelineEventWithMessage}
                commentRef={commentRef}
                isRootComment={idx === 0}
                isResolved={isActiveReplyResolved}
                onClickComment={() => {
                  if (hasSelection && editor) {
                    editor.activateComment({
                      eventId: event.id,
                    });
                  }
                }}
                onUpdateResolved={updateHighlightsForResolution}
              />
            );
          }
          if (event.event.__typename === "TLMessageResolutionV1") {
            const previousEvent = messages[idx - 1];
            const resolutionEvent = event.event as TlMessageResolutionV1;
            return (
              <div
                key={event.id}
                className={`flex items-start mb-2 ${
                  previousEvent?.event.__typename === "TLMessageV1"
                    ? "mt-4"
                    : ""
                }`}
              >
                <div className="flex-grow mt-[0.125rem]">
                  <div className="text-xs font-normal">
                    <span>{event.user.name} </span>
                    <span className="text-muted-foreground">
                      {`${resolutionEvent.resolved ? "resolved" : "reopened"} the thread`}
                    </span>
                    <TimelineDate date={event.createdAt} />
                  </div>
                </div>
              </div>
            );
          }
        })}
      </div>
    );
  },
);

CommentThread.displayName = "CommentThread";
export default CommentThread;
