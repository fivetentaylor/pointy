import React, { useEffect, useState } from "react";
import "./animations.css"; // Import the custom CSS for animations
import { CheckCheckIcon, EyeIcon, Undo2Icon } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { WithTooltip } from "../ui/FloatingTooltip";
import { RevisoUserID } from "@/constants";
import { useChatContext } from "@/contexts/ChatContext";
import {
  DOCUMENT_REVISION_ACCEPT,
  DOCUMENT_REVISION_DECLINE,
} from "@/lib/events";
import { analytics } from "@/lib/segment";
import { cn } from "@/lib/utils";
import { Spinner } from "../ui/spinner";

export const RevisionToast = function () {
  const { editor } = useRogueEditorContext();
  const [showDiffHighlights, setShowDiffHighlights] = useState(false);
  const { messages, markRevisionStatus, isRevising } = useChatContext();

  useEffect(() => {
    if (!editor) return;
    const onShowDiffHighlights = (value: boolean) => {
      editor?.clearSelection();
      setShowDiffHighlights(value);
    };

    editor.subscribe<boolean>("showDiffHighlights", onShowDiffHighlights);
    setShowDiffHighlights(editor.showDiffHighlights);

    return () => {
      editor.unsubscribe("showDiffHighlights", setShowDiffHighlights);
    };
  }, [editor]);

  useEffect(() => {
    if (!editor || !showDiffHighlights) return;
    // find the first insert or rel in the editor
    const firstInsert = editor.querySelector("ins");
    const firstDel = editor.querySelector("del");
    // figure out which element is higher in the dom hierarchy
    if (firstInsert && firstDel) {
      // Compare their position in the document
      const comparison = firstInsert.compareDocumentPosition(firstDel);

      // If firstDel comes before firstInsert
      if (comparison === Node.DOCUMENT_POSITION_PRECEDING) {
        firstDel.scrollIntoView({ behavior: "smooth", block: "center" });
      } else {
        firstInsert.scrollIntoView({ behavior: "smooth", block: "center" });
      }
    } else if (firstInsert) {
      firstInsert.scrollIntoView({ behavior: "smooth", block: "center" });
    } else if (firstDel) {
      firstDel.scrollIntoView({ behavior: "smooth", block: "center" });
    }
  }, [showDiffHighlights]);

  useEffect(() => {
    const resetEditor = () => {
      editor?.resetAddress();
    };
    if (messages && messages.length > 0) {
      const targetMessage = findFirstMessageOfLastResponder(
        messages as MessageFieldsFragment[],
      );
      if (!targetMessage) {
        return;
      }

      // console.log("target revision message", targetMessage);

      if (
        targetMessage.lifecycleStage !== "COMPLETED" &&
        targetMessage.lifecycleStage !== "REVISED"
      ) {
        resetEditor();
        return;
      }

      if (
        targetMessage.attachments.length == 0 ||
        targetMessage.metadata.revisionStatus !== "UNSPECIFIED"
      ) {
        resetEditor();
        return;
      }

      if (
        !targetMessage.attachments.some(
          (attachment) => attachment.__typename === "Revision",
        )
      ) {
        resetEditor();
        return;
      }

      editor?.setAddress(targetMessage.metadata.contentAddressBefore, "diff");
    }
  }, [messages]);

  const handleAcceptChanges = () => {
    if (!editor?.address && editor?.editorMode === "diff") {
      console.log("No address to accept");
      return;
    }
    editor?.resetAddress();

    if (messages && messages.length > 0) {
      markRevisionStatus("ACCEPTED", editor?.currentContentAddress() || "");
      analytics.track(DOCUMENT_REVISION_ACCEPT, {
        messageId: (messages[messages.length - 1] as MessageFieldsFragment).id,
      });
    }
  };

  const handleUndoChanges = () => {
    if (!editor?.address && editor?.editorMode === "diff") {
      console.log("No address to accept");
      return;
    }

    editor?.rewindAll();
    editor?.resetAddress();

    if (messages && messages.length > 0) {
      markRevisionStatus("DECLINED", editor?.currentContentAddress() || "");
      analytics.track(DOCUMENT_REVISION_DECLINE, {
        messageId: (messages[messages.length - 1] as MessageFieldsFragment).id,
      });
    }
  };

  const handleToggleHighlights = () => {
    if (editor) {
      editor.showDiffHighlights = !editor?.showDiffHighlights;
      editor.renderRogue();
    }
  };

  if (isRevising) {
    return (
      <Container>
        <div
          className="h-9 border rounded flex items-center justify-center px-2 bg-background"
          style={{
            boxShadow:
              "0 2px 4px -2px rgba(16, 24, 40, 0.1), 0 4px 6px -1px rgba(0, 0, 0, 0.1)",
          }}
        >
          <Spinner className="w-4 h-4 inline ml-0 mr-2 text-reviso" />
          <span className="text-muted-foreground">Revising...</span>
        </div>
      </Container>
    );
  }

  if (editor?.address && editor?.editorMode === "diff") {
    return (
      <Container>
        <div className="bg-reviso text-white rounded shadow-lg flex items-center justify-between h-9 pointer-events-auto">
          <div className="flex items-center">
            <span className="flex-grow px-4">Accept changes?</span>
            <WithTooltip
              tooltipText={showDiffHighlights ? "Hide Diff" : "Show Diff"}
            >
              <Button
                className={cn(
                  "hover:bg-reviso hover:text-white",
                  showDiffHighlights
                    ? "bg-primary hover:bg-primary/90"
                    : "bg-reviso",
                )}
                variant="ghost"
                size="sm"
                onClick={handleToggleHighlights}
              >
                <EyeIcon className="w-4 h-4" />
              </Button>
            </WithTooltip>
            <WithTooltip tooltipText="Undo AI Changes">
              <Button
                className="hover:bg-primary/90 hover:text-white"
                variant="ghost"
                size="sm"
                onClick={handleUndoChanges}
              >
                <Undo2Icon className="w-4 h-4" />
              </Button>
            </WithTooltip>
          </div>
          <WithTooltip tooltipText="Accept AI Changes">
            <Button
              className="hover:bg-primary/90 hover:text-white"
              variant="ghost"
              size="sm"
              onClick={handleAcceptChanges}
            >
              <CheckCheckIcon className="w-4 h-4" />
            </Button>
          </WithTooltip>
        </div>
      </Container>
    );
  }
};

const Container = ({ children }: { children: React.ReactNode }) => {
  return (
    <div className="sticky top-0 left-0 w-full z-[9999] pointer-events-none">
      <div className="absolute top-0 left-0 right-0">
        <div className="flex items-start justify-center fade-in h-[1.75rem]">
          {children}
        </div>
      </div>
    </div>
  );
};

function findFirstMessageOfLastResponder(
  messages: MessageFieldsFragment[],
): MessageFieldsFragment | undefined {
  if (messages.length === 0) return undefined;

  const lastMessage = messages[messages.length - 1];
  const lastResponderIsReviso = lastMessage.user.id === RevisoUserID;

  for (let i = messages.length - 1; i >= 0; i--) {
    const currentMessage = messages[i];
    const isReviso = currentMessage.user.id === RevisoUserID;

    if (isReviso !== lastResponderIsReviso) {
      return messages[i + 1];
    }

    if (i === 0) {
      return messages[0];
    }

    if (
      currentMessage.attachments.some(
        (attachment) => attachment.__typename === "Revision",
      )
    ) {
      return messages[i];
    }
  }

  return undefined;
}
