import React, { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

import { ChevronLeftIcon, ShareIcon } from "lucide-react";
import { ShareModalRoot } from "./ShareModalRoot";
import { EmailView, isValidEmail } from "./EmailView";
import { AccessView } from "./AccessView";
import { useErrorToast } from "@/hooks/useErrorToast";
import { analytics } from "@/lib/segment";
import {
  SHARE_ADD_EMAIL,
  SHARE_CLICK_EDIT_ACCESS,
  SHARE_OPEN,
  SHARE_REMOVE_EDITOR,
  SHARE_SEND_INVITE,
  SHARE_UPDATE_LINK,
} from "@/lib/events";
import { WithTooltip } from "@/components/ui/FloatingTooltip";
import { Button } from "@/components/ui/button";
import { useDocumentContext } from "@/contexts/DocumentContext";
import { VisuallyHidden } from "@radix-ui/react-visually-hidden";
import DocumentLoading from "@/icons/DocumentLoading";
export enum ShareState {
  Primary,
  Email,
  AccessList,
  ShareACopy,
  DownloadACopy,
}

const defaultOptions = {
  position: "center",
  initialState: ShareState.Primary,
  dismissOnShare: false,
};

export const Share = function ({
  documentId,
  options,
}: {
  documentId: string;
  options?: {
    position?: "center" | "topright";
    initialState?: ShareState;
    dismissOnShare?: boolean;
  };
}) {
  const { position, initialState, dismissOnShare } = {
    ...defaultOptions,
    ...options,
  };
  const {
    docLoading,
    refetchDoc,
    refetchSharedLinks,
    sharedLinks,
    editors,
    isDocPublic,
    shareDocument,
    updateShareLink,
    unshareDocument,
    savingDocument,
    updateDocumentVisibility,
  } = useDocumentContext();

  const [open, setOpen] = useState(false);
  const [_, setError] = useState("");
  const [email, setEmail] = useState("");
  const [currentState, setCurrentState] = useState<ShareState>(initialState);
  const [emails, setEmails] = useState<string[]>([]);
  const showErrorToast = useErrorToast();

  const onOpenChange = (open: boolean) => {
    setOpen(open);

    analytics.track(SHARE_OPEN, { documentId: documentId });

    if (open) {
      setCurrentState(initialState);
      refetchDoc();
      refetchSharedLinks();
    }
  };

  const updateCurrentState = (state: ShareState) => {
    setEmail("");
    setCurrentState(state);
  };

  const handleShareDocument = async (emails: string[], message?: string) => {
    const { data, errors } = await shareDocument(emails, message);
    if (errors) {
      setOpen(false);
      setCurrentState(initialState);
      showErrorToast("Failed to share document");
    }
    return { data, errors };
  };

  const handleUpdateDocumentVisibility = async (isPublic: boolean) => {
    const { errors } = await updateDocumentVisibility(isPublic);
    if (errors) {
      setOpen(false);
      setCurrentState(initialState);
      showErrorToast("Failed to update document visibility");
    }
  };

  const handleSubmit = async (message?: string) => {
    const validEmails = [email, ...emails].filter((email) =>
      isValidEmail(email),
    );

    analytics.track(SHARE_SEND_INVITE, {
      documentId: documentId,
      emails: validEmails,
    });

    await handleShareDocument(validEmails, message);

    setEmail("");
    setEmails([]);
    updateCurrentState(initialState);
    if (dismissOnShare) {
      setOpen(false);
    }
  };

  const handleUpdateLink = async (inviteLink: string, isActive: boolean) => {
    analytics.track(SHARE_UPDATE_LINK, {
      inviteLink,
      isActive,
      documentId: documentId,
    });

    const { errors } = await updateShareLink(inviteLink, isActive);

    if (errors) {
      setError(errors[0].message);
      setOpen(false);
      showErrorToast("Failed to update link");
      return;
    }

    setCurrentState(ShareState.Primary);
  };

  const handleRemoveEditor = async (editorId: string) => {
    analytics.track(SHARE_REMOVE_EDITOR, {
      documentId: documentId,
      editorId,
    });

    setCurrentState(ShareState.Primary);
    const { errors } = await unshareDocument(editorId);

    if (errors) {
      setError(errors[0].message);
      setOpen(false);
      showErrorToast("Failed to remove editor");
      return;
    }

    setCurrentState(ShareState.Primary);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <VisuallyHidden asChild>
        <DialogDescription>Share document modal</DialogDescription>
      </VisuallyHidden>
      <DialogTrigger asChild>
        <WithTooltip tooltipText="Share">
          <Button
            variant="icon"
            size="icon"
            onClick={() => {
              setOpen(true);
            }}
          >
            <ShareIcon className="w-4 h-4" />
          </Button>
        </WithTooltip>
      </DialogTrigger>
      <DialogContent
        className={`max-w-[28.75rem] translate-y-[-14.375rem] translate-x-[-50%] 
          ${
            position === "topright"
              ? "left-auto top-[21rem] right-[-9rem] "
              : ""
          }
          `}
      >
        <DialogHeader>
          <DialogTitle>
            {currentState === ShareState.Primary ? (
              "Share this document"
            ) : (
              <ChevronLeftIcon
                className="cursor-pointer"
                onClick={() => {
                  updateCurrentState(ShareState.Primary);
                  setEmails([]);
                }}
              />
            )}
          </DialogTitle>
        </DialogHeader>
        {currentState === ShareState.Primary && (
          <ShareModalRoot
            docId={documentId}
            isDocumentUpdating={savingDocument}
            isDocPublic={isDocPublic === undefined ? false : isDocPublic}
            email={email}
            sharedLinks={sharedLinks}
            editors={editors}
            onEmailChange={(email) => {
              const lastChar = email.charAt(email.length - 1);
              if (lastChar === " " || lastChar === ",") {
                updateCurrentState(ShareState.Email);
              } else {
                setEmail(email);
                setEmails([email]);
              }
            }}
            onAccessChange={handleUpdateDocumentVisibility}
            onEnter={() => {
              analytics.track(SHARE_ADD_EMAIL, {
                email,
                documentId: documentId,
              });
              updateCurrentState(ShareState.Email);
            }}
            onEditAccess={() => {
              analytics.track(SHARE_CLICK_EDIT_ACCESS, {
                documentId: documentId,
              });
              updateCurrentState(ShareState.AccessList);
            }}
          />
        )}
        {currentState === ShareState.Email && (
          <EmailView
            email={email}
            setEmail={setEmail}
            emails={emails}
            setEmails={setEmails}
            loadingshareDocument={savingDocument || docLoading}
            onSubmit={handleSubmit}
          />
        )}
        {currentState === ShareState.AccessList && (
          <AccessView
            sharedLinks={sharedLinks}
            editors={editors}
            onRemoveLink={(link) => handleUpdateLink(link.inviteLink, false)}
            onRemoveEditor={(editor) => {
              handleRemoveEditor(editor.id);
            }}
            email={email}
            onEmailChange={(email) => {
              const lastChar = email.charAt(email.length - 1);
              if (lastChar === " " || lastChar === ",") {
                updateCurrentState(ShareState.Email);
              } else {
                setEmail(email);
                setEmails([email]);
              }
            }}
            onEnter={() => {
              updateCurrentState(ShareState.Email);
            }}
          />
        )}
      </DialogContent>
    </Dialog>
  );
};
