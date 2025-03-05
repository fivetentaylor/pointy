import React, { useState, useEffect, useCallback } from "react";
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/avatar";
import { Document, User } from "@/__generated__/graphql";
import {
  CatIcon,
  CheckIcon,
  ClipboardIcon,
  DeleteIcon,
  EllipsisIcon,
  GraduationCap,
  MaximizeIcon,
  MinimizeIcon,
  ShareIcon,
  FileLineChartIcon,
} from "lucide-react";
import DocumentSaved from "@/icons/DocumentSaved";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { formatShortDate } from "@/lib/utils";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Spinner } from "@/components/ui/spinner";
import { Skeleton } from "../ui/skeleton";
import { Share } from "./Share";
import { SidebarTrigger } from "@/components/ui/sidebar";
import DocumentLoading from "@/icons/DocumentLoading";
import DocumentDisconnected from "@/icons/DocumentDisconnected";
import { WithTooltip } from "@/components/ui/FloatingTooltip";
import { analytics } from "@/lib/segment";
import { DOCUMENT_UPDATE_TITLE } from "@/lib/events";
import { ErrorBoundary } from "../ui/ErrorBoundary";
import { getInitials } from "@/lib/utils";
import { AuthorInfo } from "../../../rogueEditor";
import { useCursorContext } from "@/contexts/CursorContext";
import { useWsDisconnect } from "@/hooks/useWsDisconnect";
import { DeleteDraftAlert } from "../ui/DeleteDraftAlert";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { useIsMobile } from "@/hooks/use-mobile"; // Import the useIsMobile hook

export type TitleProps = {
  deleteDocument: () => void;
  deleteLoading: boolean;
  document: Pick<Document, "id" | "title" | "updatedAt" | "access"> & {
    ownedBy: { id: string };
  };
  documentState: DocumentState;
  lastEdit: Date | null;
  handleCopy: () => boolean;
  loading: boolean;
  maximized: boolean;
  me: User | null;
  toggleMaximize: () => void;
  updateTitle: (newTitle: string) => void;
};

export enum DocumentState {
  Loading,
  Disconnected,
  Saved,
}

const TitleLoading: React.FC<{ maximized: boolean; isMobile: boolean }> = ({
  maximized,
  isMobile,
}) => {
  return (
    <div
      className={cn(
        "flex items-end h-9 min-h-9 max-h-9 mr-[1.25rem] mb-4",
        maximized ? " pl-[calc(60%-38rem)]" : "",
      )}
    >
      <div className="flex items-center text-foreground text-base font-sans leading-normal flex-grow h-9">
        <Skeleton className="w-48 h-6" />
      </div>
      <div className="flex gap-2 items-center text-foreground text-base font-sans leading-normal text-right">
        <Button variant="icon" size="icon" disabled>
          <DocumentLoading />
        </Button>
        <div className="flex whitespace-nowrap overflow-hidden text-ellipsis text-muted-foreground">
          <Skeleton className="w-24 h-6" />
        </div>
        <Button variant="icon" size="icon" disabled>
          {maximized ? (
            <MinimizeIcon className="w-4 h-4" />
          ) : (
            <MaximizeIcon className="w-4 h-4" />
          )}
        </Button>
        {!isMobile && (
          <Button variant="icon" size="icon" disabled>
            <ClipboardIcon className="w-4 h-4" />
          </Button>
        )}
        <Button variant="icon" size="icon" disabled>
          <ShareIcon className="w-4 h-4" />
        </Button>
      </div>
    </div>
  );
};

const Title = ({
  deleteDocument,
  deleteLoading,
  document: doc,
  documentState,
  handleCopy,
  lastEdit,
  loading,
  maximized,
  me,
  toggleMaximize,
  updateTitle,
}: TitleProps) => {
  const { editor, editorMode } = useRogueEditorContext();
  const documentId = doc?.id || "";
  const canDelete = doc?.ownedBy?.id === me?.id;
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [newTitle, setNewTitle] = useState(doc?.title || "Untitled");
  const { isDisconnected } = useWsDisconnect();
  const [hasCopied, setHasCopied] = useState(false);
  const showStats = editorMode === "xray";
  const isMobile = useIsMobile(); // Use the useIsMobile hook

  const handleEdit = () => {
    setNewTitle(doc?.title || "Untitled");
    setIsEditing(true);
  };

  const handleKeyDown = useCallback(
    (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        setNewTitle(doc?.title || "Untitled");
        setIsEditing(false);
      }
    },
    [setIsEditing, setNewTitle, doc?.title],
  );

  const handleSubmit = (
    e:
      | React.FormEvent<HTMLFormElement>
      | React.FocusEvent<HTMLInputElement, Element>
      | React.MouseEvent,
  ) => {
    e.preventDefault();

    analytics.track(DOCUMENT_UPDATE_TITLE, {
      title: newTitle,
      documentId,
    });

    updateTitle(newTitle);
    setIsEditing(false);
  };

  const handleMaximize = () => {
    toggleMaximize();
  };

  const handleToggleStats = () => {
    if (editor) {
      editor.toggleXRayMode();
    }
  };

  useEffect(() => {
    if (isEditing) {
      window.addEventListener("keydown", handleKeyDown);
    } else {
      window.removeEventListener("keydown", handleKeyDown);
    }

    return () => {
      window.removeEventListener("keydown", handleKeyDown);
    };
  }, [isEditing, handleKeyDown]);

  let documentStateText = "Syncing...";
  switch (documentState) {
    case DocumentState.Loading:
      documentStateText = "Syncing...";
      break;
    case DocumentState.Disconnected:
      documentStateText = "Disconnected";
      break;
    case DocumentState.Saved:
      documentStateText = "All changes saved";
      break;
  }

  const onCopy = () => {
    if (hasCopied) {
      return;
    }
    const success = handleCopy();
    if (success) {
      setHasCopied(true);
      setTimeout(() => {
        setHasCopied(false);
      }, 2500);
    }
  };

  if (loading) {
    return <TitleLoading maximized={maximized} isMobile={isMobile} />;
  }

  return (
    <ErrorBoundary fallback={<div>Title Error</div>}>
      <DeleteDraftAlert
        showDeleteDialog={showDeleteDialog}
        setShowDeleteDialog={(open) => {
          setShowDeleteDialog(open);
          if (!open) {
            setTimeout(() => (document.body.style.pointerEvents = ""), 0);
          }
        }}
        onClickDelete={deleteDocument}
        draftTitle={doc?.title}
      />
      <div
        className={cn(
          "flex items-end h-9 min-h-9 max-h-9 mr-[1.25rem] mb-4 max-w-full",
          maximized ? " pl-[calc(60%-38rem)]" : "",
        )}
      >
        {/* Title section */}
        <div className="flex items-center text-foreground text-base font-sans leading-normal min-w-0 flex-grow">
          {isEditing ? (
            <>
              <form onSubmit={(e) => handleSubmit(e)}>
                <input
                  className="font-medium outline-none bg-card w-full"
                  type="text"
                  value={newTitle}
                  onChange={(e) => setNewTitle(e.target.value)}
                  onBlur={(e) => handleSubmit(e)}
                  autoFocus
                  disabled={loading}
                />
              </form>
              <Button
                className="ml-2"
                variant="icon"
                size="icon"
                onClick={(e) => handleSubmit(e)}
              >
                <CheckIcon className="w-4 h-4" />
              </Button>
            </>
          ) : (
            <div className="flex min-w-0 flex-grow items-center">
              <div
                onClick={() => handleEdit()}
                className="font-medium cursor-pointer whitespace-nowrap overflow-hidden text-ellipsis"
              >
                {doc?.title || "Untitled"}
              </div>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button
                    variant="icon"
                    size="icon"
                    className="ml-2 flex-shrink-0"
                    disabled={isDisconnected}
                  >
                    {deleteLoading ? (
                      <Spinner />
                    ) : (
                      <EllipsisIcon className="w-4 h-4" />
                    )}
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent className="w-20" align="start">
                  <DropdownMenuGroup>
                    <DropdownMenuItem
                      onClick={() => setShowDeleteDialog(true)}
                      disabled={deleteLoading || !canDelete}
                    >
                      {deleteLoading ? (
                        <Spinner />
                      ) : (
                        <>
                          <DeleteIcon className="w-4 h-4 mr-2" />
                          <span>Delete</span>
                        </>
                      )}
                    </DropdownMenuItem>
                    {me && me.isAdmin && (
                      <a
                        href={"/admin/documents/" + documentId}
                        target="_blank"
                        rel="noreferrer"
                      >
                        <DropdownMenuItem>
                          <CatIcon className="w-4 h-4 mr-2 text-reviso" />
                          <span>Admin</span>
                        </DropdownMenuItem>
                      </a>
                    )}
                  </DropdownMenuGroup>
                </DropdownMenuContent>
              </DropdownMenu>
              {doc?.access == "admin" && (
                <WithTooltip
                  tooltipText="You are an admin of this document"
                  placement="bottom"
                >
                  <Button
                    variant="icon"
                    size="icon"
                    className="cursor-default hover:bg-none text-muted-foreground"
                  >
                    <GraduationCap className="w-4 h-4 text-reviso" />
                  </Button>
                </WithTooltip>
              )}
              {!isMobile && (
                <Button
                  variant="icon"
                  size="icon"
                  className={cn("cursor-default text-muted-foreground", {
                    "hover:bg-none": !showStats,
                    "bg-muted": showStats,
                  })}
                  onClick={handleToggleStats}
                >
                  <FileLineChartIcon
                    className={cn("w-4 h-4", { "text-reviso": showStats })}
                  />
                </Button>
              )}
            </div>
          )}
          {me && me.isAdmin && process.env.IMAGE_TAG && (
            <div className="ml-2 text-muted text-xs">
              {process.env.IMAGE_TAG.substring(0, 7)}
            </div>
          )}
        </div>

        <div className="flex gap-2 items-center text-foreground text-base font-sans leading-normal text-right flex-shrink-0 ml-2">
          {/* Conditionally render SidebarTrigger on mobile */}
          {isMobile && <SidebarTrigger />}

          {/* Existing buttons and icons */}
          {!maximized && <FacePile me={me} />}

          {/* Hide timestamp, changes saved, and focus button on mobile */}
          {!isMobile && (
            <>
              <WithTooltip tooltipText={documentStateText} placement="bottom">
                <Button
                  variant="icon"
                  size="icon"
                  className="cursor-default hover:bg-none text-muted-foreground"
                >
                  {documentState === DocumentState.Loading && (
                    <DocumentLoading />
                  )}
                  {documentState === DocumentState.Disconnected && (
                    <DocumentDisconnected className="text-destructive" />
                  )}
                  {documentState === DocumentState.Saved && <DocumentSaved />}
                </Button>
              </WithTooltip>
              <div className="flex whitespace-nowrap overflow-hidden text-ellipsis text-muted-foreground">
                Edited {lastEdit ? formatShortDate(lastEdit) : "Never"}
              </div>
              <WithTooltip tooltipText={maximized ? "Minimize" : "Focus"}>
                <Button variant="icon" size="icon" onClick={handleMaximize}>
                  {maximized ? (
                    <MinimizeIcon className="w-4 h-4" />
                  ) : (
                    <MaximizeIcon className="w-4 h-4" />
                  )}
                </Button>
              </WithTooltip>
            </>
          )}

          {/* Hide copy button on mobile */}
          {!isMobile && (
            <WithTooltip
              tooltipText={
                hasCopied ? "Copied to clipboard" : "Copy to clipboard"
              }
            >
              <Button variant="icon" size="icon" onClick={onCopy}>
                {hasCopied ? (
                  <CheckIcon className="w-4 h-4" />
                ) : (
                  <ClipboardIcon className="w-4 h-4" />
                )}
              </Button>
            </WithTooltip>
          )}

          {!isDisconnected && <Share documentId={documentId} />}
          {isDisconnected && (
            <Button variant="icon" size="icon" disabled>
              <ShareIcon className="w-4 h-4" />
            </Button>
          )}
        </div>
      </div>
    </ErrorBoundary>
  );
};

const FacePile = ({ me }: { me: User | null }) => {
  const { cursors, connectedUsers } = useCursorContext();
  if (!me) {
    return null;
  }

  const cursorsWithoutMe = cursors.filter((cursor) => cursor.userID !== me.id);

  const cursorCount = cursorsWithoutMe.length;

  if (cursorCount === 0) {
    return null;
  }

  const user = (userId: string): User => {
    return connectedUsers.find((user) => user.id === userId) || ({} as User);
  };

  return (
    <div className="flex">
      {cursorsWithoutMe.map((cursor) => (
        <UserButton
          key={cursor.userID}
          cursor={cursor}
          user={user(cursor.userID)}
        />
      ))}
    </div>
  );
};

const UserButton = ({ cursor, user }: { cursor: AuthorInfo; user: User }) => {
  const handleClick = () => {
    const rogueHighlight = document.querySelector(
      "rogue-highlight#cursor-" + cursor.authorID,
    );
    if (!rogueHighlight) {
      return;
    }

    const rogueCaret = rogueHighlight.querySelector("rogue-caret") as any;
    if (rogueCaret) {
      rogueCaret.scrollIntoView({
        behavior: "smooth",
        block: "center",
        inline: "center",
      });

      rogueCaret.boop();
    }
  };
  return (
    <Avatar
      className="w-6 h-6 -ml-2 border cursor-pointer hover:border-none"
      style={{ borderColor: cursor.color }}
      onClick={handleClick}
    >
      <AvatarImage alt="Profile icon" src={user.picture || undefined} />
      <AvatarFallback className="text-background bg-primary">
        {getInitials(user.name || "")}
      </AvatarFallback>
    </Avatar>
  );
};

export default Title;
