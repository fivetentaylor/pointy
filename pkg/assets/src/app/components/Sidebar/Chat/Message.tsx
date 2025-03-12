import React, { forwardRef } from "react";
import { useNavigate } from "react-router-dom";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { MessageFieldsFragment } from "@/__generated__/graphql";
import { FileIcon, FileTextIcon } from "lucide-react";
import MessageRenderer from "./MessageRenderer";
import { cn } from "@/lib/utils";
import { timeAgo } from "@/lib/utils";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { RevisoUserID } from "@/constants";
import { ContentTypeDisplayInfoMap } from "./ContentTypeDislpay";
import { Spinner } from "@/components/ui/spinner";

const debugMessages = false;

type AttachmentType = Extract<
  MessageFieldsFragment["attachments"][number],
  { __typename: string }
>;
type SuggestionType = Extract<AttachmentType, { __typename: "Suggestion" }>;
type SelectionType = Extract<AttachmentType, { __typename: "Selection" }>;
type RevisionType = Extract<AttachmentType, { __typename: "Revision" }>;
type AttachedRevisoDocumentType = Extract<
  AttachmentType,
  { __typename: "AttachedRevisoDocument" }
>;
type AttachmentContentType = Extract<
  AttachmentType,
  { __typename: "AttachmentContent" }
>;
type AttachmentErrorType = Extract<
  AttachmentType,
  { __typename: "AttachmentError" }
>;
type AttachmentFileType = Extract<
  AttachmentType,
  { __typename: "AttachmentFile" }
>;

type MessageProps = {
  message: MessageFieldsFragment;
  previousRevisoMessage: MessageFieldsFragment | null;
};

const Message = forwardRef<HTMLDivElement, MessageProps>(
  ({ message, previousRevisoMessage: previousMessage }, ref) => {
    if (message.hidden) {
      return (
        <div>
          <DebugMessage message={message} />
        </div>
      );
    }

    if (message.user.id === RevisoUserID) {
      return (
        <RevisoMessage
          ref={ref}
          message={message}
          previousRevisoMessage={previousMessage}
        />
      );
    }

    return (
      <UserMessage
        ref={ref}
        message={message}
        previousRevisoMessage={previousMessage}
      />
    );
  },
);

const RevisoMessage = forwardRef<HTMLDivElement, MessageProps>(
  ({ message, previousRevisoMessage: previousMessage }, ref) => {
    const { editor } = useRogueEditorContext();

    const handleBeforeClick = () => {
      if (!editor) return;
      editor.setAddressDescription(
        `Draft version ${timeAgo(message.createdAt)} ago`,
      );

      if (previousMessage && previousMessage.metadata.contentAddressAfter) {
        editor.setHistoryDiff(
          previousMessage.metadata.contentAddressAfter,
          message.metadata.contentAddressBefore,
        );
        return;
      }

      editor.setAddress(
        message.metadata.contentAddressBefore,
        "history",
        false,
      );
    };

    const handleAfterClick = () => {
      editor?.setAddressDescription(
        `Draft version ${timeAgo(message.metadata.contentAddressAfterTimestamp || message.createdAt)} ago`,
      );
      editor?.setHistoryDiff(
        message.metadata.contentAddressBefore,
        message.metadata.contentAddressAfter,
      );
    };

    const uniqueContentAddress =
      message.metadata.contentAddress !==
        message.metadata.contentAddressAfter &&
      message.metadata.contentAddress !== message.metadata.contentAddressBefore;

    let hasRevisions = false;
    const filteredAttachments = message.attachments
      .filter((attachment) => {
        if (attachment.__typename !== "Revision") {
          return true;
        } else if (hasRevisions) {
          return false;
        } else {
          hasRevisions = true;
          return true;
        }
      })
      .sort((a, b) => {
        if (a.__typename === "Revision" && b.__typename !== "Revision") {
          return -1;
        } else if (a.__typename !== "Revision" && b.__typename === "Revision") {
          return 1;
        } else {
          return 0;
        }
      });

    const showPreviousDraftVersion = hasRevisions; /*&&
      message.metadata.contentAddressBefore &&
      (!previousMessage ||
        message.metadata.contentAddressBefore !==
        previousMessage.metadata.contentAddressAfter)*/

    return (
      <div ref={ref} className="mb-4 first:mt-auto">
        <DebugMessage message={message} />
        {showPreviousDraftVersion && (
          <div className="flex items-end mb-4 mt-[-0.5rem]">
            <div className="flex-grow"></div>
            <button
              className="flex items-center align-middle flex-shrink text-muted-foreground text-sm"
              onClick={handleBeforeClick}
            >
              View draft before changes
            </button>
          </div>
        )}

        <div className="flex items-start">
          <div className="flex-shrink">
            <Avatar className="w-6 h-6">
              <AvatarFallback className="text-background bg-reviso">
                R
              </AvatarFallback>
            </Avatar>
          </div>
          <div className="flex-grow ml-4">
            {message.lifecycleStage === "REVISING" && (
              <div className="text-muted-foreground mb-2">Revising...</div>
            )}

            {message.lifecycleStage === "PENDING" &&
              message.lifecycleReason != "" && (
                <div className="text-muted-foreground mb-2">
                  {message.lifecycleReason}
                </div>
              )}

            {filteredAttachments.map((attachment, idx) => (
              <Attachement
                key={`attachment-${idx}`}
                attachment={attachment}
                message={message}
              />
            ))}

            {message.content.length !== 0 && (
              <div className="mb-2">
                <MessageRenderer content={message.content} />
              </div>
            )}

            {message.aiContent?.feedback && (
              <MessageRenderer content={message.aiContent.feedback} />
            )}
            {message.aiContent?.concludingMessage && (
              <MessageRenderer content={message.aiContent.concludingMessage} />
            )}
          </div>
        </div>

        {message.metadata.contentAddressAfter && uniqueContentAddress && (
          <div className="flex items-start mt-2 first:mt-auto">
            <button
              className="ml-[2.375rem] flex items-center align-middle flex-shrink text-muted-foreground text-sm"
              onClick={handleAfterClick}
            >
              View accepted changes
            </button>
            <div className="flex-grow"></div>
          </div>
        )}
      </div>
    );
  },
);

const UserMessage = forwardRef<HTMLDivElement, MessageProps>(
  ({ message }, ref) => {
    return (
      <div ref={ref} className="flex flex-col justify-center first:mt-auto">
        <DebugMessage message={message} />
        <div className="flex items-center justify-end">
          <div className="flex flex-col items-end justify-end max-w-[-moz-fill-available] max-w-[-webkit-fill-available]">
            {[...message.attachments]
              .sort((a, b) => {
                if (
                  a.__typename === "AttachmentContent" &&
                  b.__typename !== "AttachmentContent"
                ) {
                  return 1;
                } else if (
                  a.__typename !== "AttachmentContent" &&
                  b.__typename === "AttachmentContent"
                ) {
                  return -1;
                }
                return 0;
              })
              .map((attachment, idx) => (
                <Attachement
                  key={idx}
                  attachment={attachment}
                  message={message}
                />
              ))}
          </div>
        </div>
        <div className="flex items-center justify-end mb-4">
          <div className="ml-4 rounded-[26px] bg-elevated p-4">
            {message.lifecycleStage === "PENDING" && (
              <Spinner className="w-5 h-5 ml-3" />
            )}
            <MessageRenderer content={message.content} />
          </div>
        </div>
      </div>
    );
  },
);

const DebugMessage = ({ message }: { message: MessageFieldsFragment }) => {
  if (!message || !debugMessages) {
    return null;
  }

  return (
    <div className="text-xs border-t mb-2 pt-2">
      <div className="flex items-center">
        <span className="text-muted-foreground">{message.id}</span>
        <div className="flex-grow"></div>
        <span className="text-muted-foreground">
          {message.lifecycleStage} | {message.metadata.revisionStatus}
        </span>
      </div>
      <div className="flex items-center text-muted-foreground">
        <span>{message.metadata.contentAddressBefore || "NULL"}</span>
        <div className="flex-grow text-center">|</div>
        <span>{message.metadata.contentAddress || "NULL"}</span>
        <div className="flex-grow text-center">|</div>
        <span>{message.metadata.contentAddressAfter || "NULL"}</span>
      </div>

      {message.hidden && (
        <div className="flex items-center text-muted-foreground">
          <span>Hidden: </span>
          <span>{message.content}</span>
        </div>
      )}
    </div>
  );
};

const Attachement = ({
  message,
  attachment,
}: {
  message: MessageFieldsFragment;
  attachment: MessageFieldsFragment["attachments"][number];
}) => {
  if (attachment.__typename === "Selection") {
    return <Selection selection={attachment} />;
  }
  if (
    attachment.__typename === "Revision" ||
    attachment.__typename === "Suggestion"
  ) {
    return <Delta delta={attachment} message={message} />;
  }
  if (attachment.__typename === "AttachmentContent") {
    return <AttachmentContent attachment={attachment} />;
  }
  if (attachment.__typename === "AttachmentError") {
    return <AttachmentError attachment={attachment} />;
  }
  if (attachment.__typename === "AttachmentFile") {
    return <AttachedFile attachment={attachment} />;
  }
  if (attachment.__typename === "AttachedRevisoDocument") {
    return <AttachedRevisoDocument attachment={attachment} />;
  }

  return null;
};

const AttachmentContent = ({
  attachment,
}: {
  attachment: AttachmentContentType;
}) => {
  return (
    <>{attachment.text && <MessageRenderer content={attachment.text} />}</>
  );
};

const Selection = ({ selection }: { selection: SelectionType }) => {
  return (
    <div className="mb-2 text-muted-foreground w-full overflow-hidden whitespace-nowrap text-ellipsis text-right pl-9">
      {selection.content}
    </div>
  );
};

const AttachmentError = ({
  attachment,
}: {
  attachment: AttachmentErrorType;
}) => {
  return (
    <div className="mb-2 w-full mt-0.5">
      <span className="text-muted-foreground">{attachment.title}</span>
      <div className="mt-1">{attachment.text}</div>
    </div>
  );
};

const AttachedRevisoDocument = ({
  attachment,
}: {
  attachment: AttachedRevisoDocumentType;
}) => {
  const navigate = useNavigate();

  const handleClick = () => {
    const params = new URLSearchParams(location.search);
    const sbParam = params.get("sb");
    let target = `/drafts/${attachment.id}`;
    if (sbParam) {
      target += `?sb=${sbParam}`;
    }
    navigate(target);
  };

  return (
    <div
      className={cn(
        "w-[250px] p-2 bg-card rounded-md shadow border border-border justify-start items-start gap-2 inline-flex mb-2",
        "hover:bg-muted hover:cursor-pointer",
      )}
      onClick={handleClick}
    >
      <div
        className={cn(
          "w-9 h-9 rounded-md justify-center items-center flex flex-shrink-0",
          "bg-reviso",
        )}
      >
        <FileIcon className="w-4 h-4 text-white" />
      </div>
      <div className="h-9 flex-1 min-w-0">
        <div className="mt-0.5 text-foreground text-xs font-medium font-['Inter'] leading-none truncate overflow-hidden whitespace-nowrap">
          {attachment.title}
        </div>
        <div className="text-muted-foreground text-xs font-normal font-['Inter'] leading-none mt-2">
          Your Draft
        </div>
      </div>
    </div>
  );
};

const AttachedFile = ({ attachment }: { attachment: AttachmentFileType }) => {
  const handleClick = () => {
    if (attachment.contentType !== "text/url") {
      return;
    }

    window.open(attachment.filename, "_blank");
  };

  const identifier = attachment.filename;
  let extraClasses = "";

  const displayInfo = ContentTypeDisplayInfoMap[
    attachment.contentType || ""
  ] || {
    label: "",
    color: "#000000",
    icon: <FileTextIcon className="w-4 h-4 text-white" />,
  };
  const subtitle = displayInfo.label;
  const color = displayInfo.color;
  const icon = displayInfo.icon;

  if (attachment.contentType === "text/url") {
    extraClasses = "hover:bg-muted hover:cursor-pointer";
  }

  return (
    <div
      className={cn(
        "w-[250px] p-2 bg-card rounded-md shadow border border-border justify-start items-start gap-2 inline-flex mb-2",
        extraClasses,
      )}
      onClick={handleClick}
    >
      <div
        className={cn(
          "w-9 h-9 rounded-md justify-center items-center flex flex-shrink-0",
          color,
        )}
      >
        {icon}
      </div>
      <div className="h-9 flex-1 min-w-0">
        <div className="mt-0.5 text-foreground text-xs font-medium font-['Inter'] leading-none truncate overflow-hidden whitespace-nowrap">
          {identifier}
        </div>
        <div className="text-muted-foreground text-xs font-normal font-['Inter'] leading-none mt-2">
          {subtitle}
        </div>
      </div>
    </div>
  );
};

const Delta = ({
  message,
}: {
  message: MessageFieldsFragment;
  delta: RevisionType | SuggestionType;
}) => {
  const { editor } = useRogueEditorContext();

  const handleClick = () => {
    editor?.setAddressDescription(
      `Revised version ${timeAgo(message.createdAt)} ago`,
    );
    editor?.setHistoryDiff(
      message.metadata.contentAddressBefore,
      message.metadata.contentAddress,
    );
  };

  return (
    <div className={cn("mb-2 mt-0.5 text-muted-foreground w-full")}>
      {message.lifecycleStage === "REVISED" && <>Revised</>}
      {message.lifecycleStage === "COMPLETED" &&
        message.metadata.contentAddress && (
          <button onClick={handleClick}>View revised version</button>
        )}
    </div>
  );
};

Message.displayName = "Message";
RevisoMessage.displayName = "PointyMessage";
UserMessage.displayName = "UserMessage";
Selection.displayName = "Selection";
Delta.displayName = "Delta";

export default Message;
