import React, { useEffect, useRef, useState } from "react";
import {
  CheckCircle2Icon,
  ClipboardPasteIcon,
  MilestoneIcon,
  PencilIcon,
  TrashIcon,
} from "lucide-react";
import { arrayToSentence, cn, timeAgo } from "@/lib/utils";
import {
  TimelineEventFieldsFragment,
  TlAttributeChangeV1,
} from "@/__generated__/graphql";
import { TimelineEventWithMessage, TLMessageEntry } from "./CommentThread";
import MessageRenderer from "../Chat/MessageRenderer";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { TimelineAvatar, TimelineDate } from "./TimelineHelpers";
import { WithTooltip } from "@/components/ui/FloatingTooltip";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useTimelineContext } from "./TimelineContext";
import { useDocumentContext } from "@/contexts/DocumentContext";
import {
  TIMELINE_DELETE_FLAGGED_VERSION,
  TIMELINE_FLAG_VERSION,
  TIMELINE_SHOW_UPDATE_VERSION,
  TIMELINE_UPDATE_COMMENT_SUMMARY,
  TIMELINE_UPDATE_VERSION_SUMMARY,
} from "@/lib/events";
import { analytics } from "@/lib/segment";

export const TimelineSpacer: React.FC<{ date: Date | string }> = ({ date }) => {
  const formattedDate =
    typeof date === "string"
      ? date
      : date
          .toLocaleDateString("en-US", {
            weekday: "short",
            month: "short",
            day: "numeric",
            year: undefined,
          })
          .replace(",", "");

  return (
    <div className="flex items-center my-4">
      <div className="flex-grow border-t border-gray-300"></div>
      <div className="mx-4 text-sm text-muted-foreground">{formattedDate}</div>
      <div className="flex-grow border-t border-gray-300"></div>
    </div>
  );
};

export const TimelineEntry: React.FC<{
  timelineEvent: TimelineEventFieldsFragment;
}> = ({ timelineEvent }) => {
  if (timelineEvent.event?.__typename === "TLUpdateV1") {
    return <TLUpdate key={timelineEvent.id} timelineEvent={timelineEvent} />;
  } else if (timelineEvent.event?.__typename === "TLMessageV1") {
    return (
      <TLMessageEntry
        timelineEvent={timelineEvent as TimelineEventWithMessage}
      />
    );
  } else if (timelineEvent.event?.__typename === "TLMessageResolutionV1") {
    return (
      <TLMessageResolution
        timelineEvent={timelineEvent as TimelineEventWithMessage}
      />
    );
  } else if (timelineEvent.event?.__typename === "TLMarkerV1") {
    return (
      <TimelineMarker
        key={timelineEvent.id}
        title={timelineEvent.event.title}
      />
    );
  } else if (timelineEvent.event?.__typename === "TLJoinV1") {
    return <TLJoin key={timelineEvent.id} timelineEvent={timelineEvent} />;
  } else if (timelineEvent.event?.__typename === "TLAccessChangeV1") {
    return (
      <TLAccessChange key={timelineEvent.id} timelineEvent={timelineEvent} />
    );
  } else if (timelineEvent.event?.__typename === "TLAttributeChangeV1") {
    return (
      <TLAttributeChange key={timelineEvent.id} timelineEvent={timelineEvent} />
    );
  } else if (timelineEvent.event?.__typename === "TLPasteV1") {
    return <TLPaste key={timelineEvent.id} timelineEvent={timelineEvent} />;
  } else {
    console.error("Unexpected event type", timelineEvent.event);
    return null;
  }
};

const TimelineMarker: React.FC<{ title: string }> = ({ title }) => {
  return (
    <div className="flex items-center my-4">
      <div className="flex-grow border-t border-gray-300"></div>
      <div className="mx-4 text-sm text-muted-foreground">{title}</div>
      <div className="flex-grow border-t border-gray-300"></div>
    </div>
  );
};

const TLJoin = ({
  timelineEvent,
}: {
  timelineEvent: TimelineEventFieldsFragment;
}) => {
  if (timelineEvent.event?.__typename !== "TLJoinV1") {
    console.error("Unexpected event type", timelineEvent.event?.__typename);
    return null;
  }

  return (
    <div className="flex items-start mb-4 first:mt-auto px-4">
      <div className="w-4 h-4 flex items-center justify-start mr-2 ml-0.5">
        <TimelineAvatar user={timelineEvent.user} />
      </div>
      <div className="flex-grow mt-[0.125rem]">
        <div className="text-sm font-normal">
          <span>{timelineEvent.user.name} </span>
          <span className="">
            {timelineEvent.event.action === "create" ? (
              <>
                {"created "}
                <span className="">the draft</span>
              </>
            ) : (
              " opened the draft for the first time"
            )}
          </span>
          <TimelineDate date={timelineEvent.createdAt} />
        </div>
      </div>
    </div>
  );
};

const UpdateCommentResolutionSummaryDialog = ({
  isOpen,
  setIsOpen,
  initialSummary,
  handleUpdateChangeSummary,
}: {
  isOpen: boolean;
  setIsOpen: (isOpen: boolean) => void;
  initialSummary: string;
  handleUpdateChangeSummary: (summary: string) => void;
}) => {
  const [summary, setSummary] = useState(initialSummary);

  useEffect(() => {
    setSummary(initialSummary);
  }, [isOpen]);

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Summarize thread</DialogTitle>
          <DialogDescription>
            Describe the topic of the comment thread.
          </DialogDescription>
          <div>
            <div className="my-4">
              <Label className="block mb-2" htmlFor="summary">
                Summary
              </Label>
              <div className="flex items-center">
                <textarea
                  className="mt-1 border border-border rounded-md px-3.5 py-3 min-h-[8.125rem] w-full bg-transparent outline-none resize-none"
                  onChange={(e) => setSummary(e.target.value)}
                  onKeyDown={(e) => {
                    if (e.key === "Enter") {
                      handleUpdateChangeSummary(summary);
                    }
                  }}
                  value={summary}
                />
              </div>
            </div>
            <div className="flex justify-end">
              <Button
                size="sm"
                className="bg-primary hover:bg-primary/90"
                disabled={!summary.trim()}
                onClick={() => {
                  handleUpdateChangeSummary(summary);
                }}
              >
                Save
              </Button>
            </div>
          </div>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
};

export const TLMessageResolution: React.FC<{
  timelineEvent: TimelineEventFieldsFragment;
}> = ({ timelineEvent }) => {
  const { setActiveReplyId } = useTimelineContext();
  const [isUpdatingSummary, setIsUpdatingSummary] = useState(false);
  const resolutionEvent = timelineEvent.event as TlMessageResolutionV1;
  const { editMessageResolutionSummary } = useTimelineContext();

  const handleUpdateChangeSummary = (summary: string) => {
    analytics.track(TIMELINE_UPDATE_COMMENT_SUMMARY);
    setIsUpdatingSummary(false);
    editMessageResolutionSummary(timelineEvent.id, summary);
  };

  return (
    <>
      <UpdateCommentResolutionSummaryDialog
        isOpen={isUpdatingSummary}
        setIsOpen={setIsUpdatingSummary}
        initialSummary={resolutionEvent.resolutionSummary}
        handleUpdateChangeSummary={handleUpdateChangeSummary}
      />
      <div className="flex flex-col mb-4 first:mt-auto pl-1 pt-2 relative group hover:bg-elevated">
        <div className="absolute top-0 right-0 w-[33px] h-[34px] p-[1px] mr-2 flex items-center justify-start bg-card transform -translate-y-[50%] rounded-md opacity-0 group-hover:opacity-100 transition-opacity shadow-sm">
          <WithTooltip tooltipText="Edit summary">
            <Button
              variant="ghost"
              size="icon"
              className="p-2 w-8 h-8 min-w-8 text-muted-foreground"
              onClick={() => setIsUpdatingSummary(true)}
            >
              <PencilIcon className="w-4 h-4 min-w-4 text-muted-foreground" />
            </Button>
          </WithTooltip>
        </div>
        <div className="flex items-start">
          <div className="w-4 h-4 flex items-center justify-start mr-2 ml-0.5 mt-1">
            <CheckCircle2Icon className="w-4 h-4 text-muted-foreground" />
          </div>
          <div className="flex-grow mt-[0.125rem]">
            <div className="text-sm font-medium">
              <span>{timelineEvent.user.name} </span>
              <span>resolved a thread</span>
              <TimelineDate date={timelineEvent.createdAt} />
            </div>
          </div>
        </div>
        <div className="flex flex-col ml-[1.625rem]">
          {resolutionEvent.resolutionSummary && (
            <div className="text-sm font-normal">
              <span>{resolutionEvent.resolutionSummary}</span>
            </div>
          )}
          <div>
            <Button
              variant="link"
              className="p-0 text-sm text-primary m-0 mt-[-1rem] h-4"
              onClick={() => setActiveReplyId(timelineEvent.replyTo)}
            >
              View thread
            </Button>
          </div>
        </div>
      </div>
    </>
  );
};

const AttributeChangeText = function ({
  event,
}: {
  event: TlAttributeChangeV1;
}) {
  if (event.attribute === "is_public") {
    return (
      <>
        <span className="">updated draft public visibility to</span>
        <WithTooltip
          tooltipText={"Previous " + event.attribute + ": " + event.oldValue}
        >
          <span className="font-normal">{event.newValue}</span>
        </WithTooltip>
      </>
    );
  }
  return (
    <>
      <span className="">changed {event.attribute} to </span>
      <WithTooltip
        tooltipText={"Previous " + event.attribute + ": " + event.oldValue}
      >
        <span className="font-normal">{event.newValue}</span>
      </WithTooltip>
    </>
  );
};

const TLAttributeChange = ({
  timelineEvent,
}: {
  timelineEvent: TimelineEventFieldsFragment;
}) => {
  if (timelineEvent.event?.__typename !== "TLAttributeChangeV1") {
    console.error("Unexpected event type", timelineEvent.event?.__typename);
    return null;
  }

  return (
    <div className="flex items-start mb-4 first:mt-auto px-4">
      <div className="w-4 h-4 flex items-center justify-start mr-2 ml-0.5 mt-1">
        <PencilIcon className="w-4 h-4 text-muted-foreground" />
      </div>
      <div className="flex-grow mt-[0.125rem]">
        <div className="text-sm font-normal">
          <span>{timelineEvent.user.name} </span>
          <AttributeChangeText event={timelineEvent.event} />
          <TimelineDate date={timelineEvent.createdAt} />
        </div>
      </div>
    </div>
  );
};

const TLPaste = ({
  timelineEvent,
}: {
  timelineEvent: TimelineEventFieldsFragment;
}) => {
  if (timelineEvent.event?.__typename !== "TLPasteV1") {
    console.error("Unexpected event type", timelineEvent.event?.__typename);
    return null;
  }

  const { editor } = useRogueEditorContext();

  const handleShowChanges = () => {
    if (timelineEvent.event?.__typename !== "TLPasteV1") {
      console.error("Unexpected event type", timelineEvent.event?.__typename);
      return null;
    }

    if (editor) {
      const description =
        timelineEvent.user.name + " pasted " + timeAgo(timelineEvent.createdAt);

      editor.setAddressDescription(description);
      editor.setHistoryDiff(
        timelineEvent.event.contentAddressBefore,
        timelineEvent.event.contentAddressAfter,
      );
    }
  };

  return (
    <div
      className="flex items-start p-4 first:mt-auto hover:bg-elevated cursor-pointer"
      onClick={handleShowChanges}
    >
      <div className="w-4 h-4 flex items-center justify-start mr-2 ml-0.5 mt-1">
        <ClipboardPasteIcon className="w-4 h-4 text-muted-foreground" />
      </div>
      <div className="flex-grow mt-[0.125rem]">
        <div className="text-sm font-normal">
          <span>{timelineEvent.user.name} </span>
          <span className="">pasted text into the draft</span>
          <TimelineDate date={timelineEvent.createdAt} />
        </div>
      </div>
    </div>
  );
};

const TLAccessChange = ({
  timelineEvent,
}: {
  timelineEvent: TimelineEventFieldsFragment;
}) => {
  if (timelineEvent.event?.__typename !== "TLAccessChangeV1") {
    console.error("Unexpected event type", timelineEvent.event?.__typename);
    return null;
  }

  return (
    <div className="flex items-start mb-4 first:mt-auto px-4">
      <div className="w-4 h-4 flex items-center justify-start mr-2 ml-0.5">
        <TimelineAvatar user={timelineEvent.user} />
      </div>
      <div className="flex-grow mt-[0.125rem]">
        <div className="text-sm font-normal">
          <span>{timelineEvent.user.name} </span>
          <span className="">
            {timelineEvent.event.action === "REMOVE_ACTION" ? (
              <>
                {"revoked access to "}
                <span className="">
                  {arrayToSentence(timelineEvent.event.userIdentifiers)}
                </span>
              </>
            ) : (
              <>
                {"invited "}
                <span className="">
                  {arrayToSentence(timelineEvent.event.userIdentifiers)}
                </span>
              </>
            )}
          </span>
          <TimelineDate date={timelineEvent.createdAt} />
        </div>
      </div>
    </div>
  );
};

const FlagDialog = ({
  flagName,
  isFlaggingVersion,
  setIsFlaggingVersion,
  handleFlagVersion,
}: {
  flagName: string;
  isFlaggingVersion: boolean;
  handleFlagVersion: (name: string) => void;
  setIsFlaggingVersion: (isFlaggingVersion: boolean) => void;
}) => {
  const [versionName, setVersionName] = useState(flagName);

  useEffect(() => {
    setVersionName(flagName);
  }, [isFlaggingVersion]);

  return (
    <Dialog open={isFlaggingVersion} onOpenChange={setIsFlaggingVersion}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Flag version</DialogTitle>
          <DialogDescription>
            Give this a name to refer to in the timeline.
          </DialogDescription>
          <div>
            <div className="mt-4 mb-8">
              <Label className="block mb-2" htmlFor="name">
                Name
              </Label>
              <Input
                className="w-full"
                onChange={(e) => setVersionName(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === "Enter") {
                    handleFlagVersion(versionName);
                  }
                }}
                value={versionName}
              />
            </div>
            <div className="flex justify-end">
              <Button
                size="sm"
                className="bg-primary hover:bg-primary/90"
                disabled={!versionName}
                onClick={() => {
                  handleFlagVersion(versionName);
                }}
              >
                Save
              </Button>
            </div>
          </div>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
};

const UpdateChangeSummaryDialog = ({
  isOpen,
  setIsOpen,
  initialSummary,
  handleUpdateChangeSummary,
}: {
  isOpen: boolean;
  setIsOpen: (isOpen: boolean) => void;
  initialSummary: string;
  handleUpdateChangeSummary: (summary: string) => void;
}) => {
  const [summary, setSummary] = useState(initialSummary);

  useEffect(() => {
    setSummary(initialSummary);
  }, [isOpen]);

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Summarize edits</DialogTitle>
          <DialogDescription>
            Describe what edits were made during this session.
          </DialogDescription>
          <div>
            <div className="my-4">
              <Label className="block mb-2" htmlFor="summary">
                Summary
              </Label>
              <div className="text-muted-foreground text-sm">
                We recommend starting the summary with a verb
              </div>
              <div className="flex items-center">
                <textarea
                  className="mt-1 border border-border rounded-md px-3.5 py-3 min-h-[8.125rem] w-full bg-transparent outline-none resize-none"
                  onChange={(e) => setSummary(e.target.value)}
                  onKeyDown={(e) => {
                    if (e.key === "Enter") {
                      handleUpdateChangeSummary(summary);
                    }
                  }}
                  value={summary}
                />
              </div>
            </div>
            <div className="flex justify-end">
              <Button
                size="sm"
                className="bg-primary hover:bg-primary/90"
                disabled={!summary.trim()}
                onClick={() => {
                  handleUpdateChangeSummary(summary);
                }}
              >
                Save
              </Button>
            </div>
          </div>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
};

const TLUpdate: React.FC<{ timelineEvent: TimelineEventFieldsFragment }> = ({
  timelineEvent,
}) => {
  const { editor } = useRogueEditorContext();
  const [isShowingChanges, setIsShowingChanges] = useState(false);
  const selectedAddresses = useRef<(string | null)[] | null>(null); // [baseAddress, address]
  const [isFlaggingVersion, setIsFlaggingVersion] = useState(false);
  const [isUpdatingChangeSummary, setIsUpdatingChangeSummary] = useState(false);
  const { editTimelineUpdateSummary } = useTimelineContext();
  const { createFlaggedVersion, editFlaggedVersion, deleteFlaggedVersion } =
    useDocumentContext();

  useEffect(() => {
    if (editor) {
      if (
        selectedAddresses.current &&
        (selectedAddresses.current[0] !== editor.baseAddress ||
          selectedAddresses.current[1] !== editor.address)
      ) {
        selectedAddresses.current = null;
        setIsShowingChanges(false);
      }
    }
  }, [editor, editor?.address]);

  if (timelineEvent.event?.__typename !== "TLUpdateV1") {
    console.error("Unexpected event type", timelineEvent.event?.__typename);
    return null;
  }

  const event = timelineEvent.event as TlUpdateV1;
  const existingFlag = timelineEvent.event.flaggedVersionName;
  const hasFlag = !!(existingFlag && existingFlag !== "");

  const handleFlagVersion = (name: string) => {
    if (event.flaggedVersionID) {
      editFlaggedVersion(event.flaggedVersionID, {
        name,
        updateID: timelineEvent.id,
      });
    } else {
      analytics.track(TIMELINE_FLAG_VERSION);
      createFlaggedVersion({
        name,
        updateID: timelineEvent.id,
      });
    }
    setIsFlaggingVersion(false);
  };

  const handleDeleteVersion = () => {
    if (event.flaggedVersionID) {
      analytics.track(TIMELINE_DELETE_FLAGGED_VERSION);
      deleteFlaggedVersion(event.flaggedVersionID, timelineEvent.id);
    }
  };

  const handleShowChanges = () => {
    if (timelineEvent.event?.__typename !== "TLUpdateV1") {
      console.error("Unexpected event type", timelineEvent.event?.__typename);
      return null;
    }

    if (editor) {
      analytics.track(TIMELINE_SHOW_UPDATE_VERSION);
      const description =
        timelineEvent.user.name + " edited " + timeAgo(timelineEvent.createdAt);

      if (isShowingChanges && editor.addressDescription === description) {
        setIsShowingChanges(false);
        editor.showDiffHighlights = false;
        editor.resetAddress();
        return;
      }
      setIsShowingChanges(true);

      editor.setAddressDescription(description);
      editor.setHistoryDiff(
        timelineEvent.event.startingContentAddress,
        timelineEvent.event.endingContentAddress,
      );

      selectedAddresses.current = [editor.baseAddress, editor.address];
    }
  };

  const handleUpdateChangeSummary = (newSummary: string) => {
    analytics.track(TIMELINE_UPDATE_VERSION_SUMMARY);
    editTimelineUpdateSummary(timelineEvent.id, newSummary);
    setIsUpdatingChangeSummary(false);
  };

  const content = timelineEvent.event.content;
  const firstWord = content?.split(" ")[0].toLowerCase();
  const restOfContent = content?.slice(firstWord.length).trim();

  return (
    <>
      <FlagDialog
        flagName={hasFlag ? existingFlag : "Rough draft"}
        isFlaggingVersion={isFlaggingVersion}
        setIsFlaggingVersion={setIsFlaggingVersion}
        handleFlagVersion={handleFlagVersion}
      />
      <UpdateChangeSummaryDialog
        isOpen={isUpdatingChangeSummary}
        setIsOpen={setIsUpdatingChangeSummary}
        initialSummary={timelineEvent.event.content || ""}
        handleUpdateChangeSummary={handleUpdateChangeSummary}
      />
      <div
        className={cn(
          "group flex items-start mt-[-1rem] px-4 first:mt-auto relative hover:bg-elevated cursor-pointer",
          hasFlag ? "pt-4 pb-2 mb-2 rounded-t-md" : "py-4 rounded-md",
        )}
        onClick={() => handleShowChanges()}
      >
        <div className="absolute top-0 right-0 w-[66px] h-[34px] p-[1px] mr-2 flex items-center justify-start bg-card transform -translate-y-[50%] rounded-md opacity-0 group-hover:opacity-100 transition-opacity shadow-sm">
          <WithTooltip tooltipText="Flag this version">
            <Button
              variant="ghost"
              size="icon"
              className="p-2 w-8 h-8 min-w-8 text-muted-foreground"
              onClick={(event) => {
                event.preventDefault();
                event.stopPropagation();
                setIsFlaggingVersion(true);
              }}
              disabled={hasFlag}
            >
              <MilestoneIcon className="w-4 h-4 min-w-4 text-muted-foreground" />
            </Button>
          </WithTooltip>
          <WithTooltip tooltipText="Edit change summary">
            <Button
              variant="ghost"
              size="icon"
              className="p-2 w-8 h-8 min-w-8 text-muted-foreground"
              onClick={(event) => {
                event.preventDefault();
                event.stopPropagation();
                setIsUpdatingChangeSummary(true);
              }}
            >
              <PencilIcon className="w-4 h-4 min-w-4 text-muted-foreground" />
            </Button>
          </WithTooltip>
        </div>
        <div className="w-4 h-4 flex items-center justify-start mr-2 mt-1 ml-0.5">
          <PencilIcon className="w-4 h-4 text-muted-foreground" />
        </div>
        <div className="flex-grow mt-[0.125rem]">
          {timelineEvent.event.content &&
          /^[a-zA-Z]/.test(timelineEvent.event.content) ? (
            <div className="space-x-1 text-sm font-normal">
              <span>{timelineEvent.user.name}</span>
              <span className="">{firstWord}</span>
              <span>
                <MessageRenderer
                  content={restOfContent}
                  variant="timeline-update"
                />
              </span>
              <TimelineDate date={timelineEvent.createdAt} />
            </div>
          ) : (
            <>
              <div className="text-sm font-normal">
                <span>{timelineEvent.user.name}</span>
                <span className="">{timelineEvent.event.title}</span>
                <TimelineDate date={timelineEvent.createdAt} />
              </div>
              {timelineEvent.event.content && (
                <div className="ml-2">
                  <MessageRenderer
                    content={timelineEvent.event.content}
                    variant="timeline"
                  />
                </div>
              )}
            </>
          )}
        </div>
      </div>
      {hasFlag && (
        <FlagVersion
          timelineEvent={timelineEvent.event as TlUpdateV1}
          onClickEdit={() => setIsFlaggingVersion(true)}
          onClickDelete={handleDeleteVersion}
        />
      )}
    </>
  );
};

const FlagVersion: React.FC<{
  timelineEvent: TlUpdateV1;
  onClickEdit: () => void;
  onClickDelete: () => void;
}> = ({ timelineEvent, onClickEdit, onClickDelete }) => {
  const { editor } = useRogueEditorContext();
  const [showingVersion, setShowingVersion] = useState(false);

  return (
    <>
      <div
        className="flex items-start mt-[-0.25rem] mb-1 py-1 pl-3 first:mt-auto  hover:bg-elevated group rounded-b-md group relative z-48 cursor-pointer"
        onClick={() => {
          if (!editor) {
            return;
          }
          const description = `${timelineEvent.flaggedVersionName}`;

          if (showingVersion && editor.addressDescription === description) {
            editor.resetAddress();
            editor?.setAddressDescription("");
            setShowingVersion(false);
          } else {
            editor?.setAddress(
              timelineEvent.endingContentAddress,
              "history",
              false,
            );
            editor?.setAddressDescription(description);
            setShowingVersion(true);
          }
        }}
      >
        <div className="absolute top-0 right-0 w-[66px] h-[34px] mr-2 flex items-center justify-start bg-card transform -translate-y-[75%] rounded-md opacity-0 group-hover:opacity-100 transition-opacity border-1 border-border p-[1px] shadow-sm">
          <WithTooltip tooltipText="Change version name">
            <Button
              variant="ghost"
              size="icon"
              className="p-2 w-8 h-8 min-w-8 text-muted-foreground"
              onClick={(event) => {
                event.stopPropagation();
                event.preventDefault();
                onClickEdit();
              }}
            >
              <PencilIcon className="w-4 h-4 min-w-4 text-muted-foreground" />
            </Button>
          </WithTooltip>
          <WithTooltip tooltipText="Delete flag">
            <Button
              variant="ghost"
              size="icon"
              className="p-2 w-8 h-8 min-w-8 text-muted-foreground"
              onClick={(event) => {
                event.stopPropagation();
                event.preventDefault();
                onClickDelete();
              }}
            >
              <TrashIcon className="w-4 h-4 min-w-4 text-muted-foreground" />
            </Button>
          </WithTooltip>
        </div>
        <div className="border-l-2 border-b-2 border-border rounded-bl-3xl h-6 w-2 mr-1 ml-2 mt-[-0.75rem] z-50" />
        <div className="w-4 h-4 flex items-center justify-start mr-2 ml-2 mt-1">
          <MilestoneIcon className="w-4 h-4 text-muted-foreground" />
        </div>
        <div className="flex-grow mt-[0.125rem]">
          <div className="text-sm font-normal">
            <span>{timelineEvent.flaggedByUser?.name} </span>
            <span className="">flagged version </span>
            <span className="text-primary">
              {timelineEvent.flaggedVersionName}
            </span>
            {timelineEvent.flaggedVersionCreatedAt && (
              <TimelineDate date={timelineEvent.flaggedVersionCreatedAt} />
            )}
          </div>
        </div>
      </div>
    </>
  );
};
