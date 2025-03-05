import React, { CSSProperties, forwardRef, useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Edit2Icon,
  HelpCircleIcon,
  LinkIcon,
  HistoryIcon,
  MessageCircleIcon,
} from "lucide-react";
import { AI_CLICK_EXPLAIN } from "@/lib/events";
import { AskReviso } from "./AskReviso";
import LinkEditor from "./LinkEditor";
import CommentEditor from "./CommentEditor";
import AskAIAnything from "./AskAIAnything";
import Scrub from "./Scrub";
import { useWsDisconnect } from "@/hooks/useWsDisconnect";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { analytics } from "@/lib/segment";
import { useSidebarContext } from "@/contexts/SidebarContext";
import { useChatContext } from "@/contexts/ChatContext";

export const BubbleMenu = forwardRef<
  HTMLDivElement,
  {
    style?: CSSProperties;
    hide: boolean;
    onDismiss: () => void;
    hasLink: boolean;
    toggleMaximize: () => void;
    setHasCommentInput: (hasCommentInput: boolean) => void;
  }
>(function BubbleMenu(
  { style, onDismiss, hide, hasLink, toggleMaximize, setHasCommentInput },
  floatingRef,
) {
  const [menuView, setMenuView] = useState<
    "default" | "link" | "comment" | "askAIAnything" | "scrub"
  >(hasLink ? "link" : "default");
  const [container, setContainer] = useState<HTMLDivElement | null>(null);
  const { editor } = useRogueEditorContext();
  const { createThreadMessage } = useChatContext();
  const { setSidebarMode } = useSidebarContext();
  const { isDisconnected } = useWsDisconnect();

  useEffect(() => {
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.metaKey && event.shiftKey && event.key === "c") {
        setMenuView("comment");
      }

      if (event.metaKey && event.key === "k") {
        setMenuView("link");
      }

      if (event.metaKey && event.ctrlKey && event.key === "y") {
        setMenuView("askAIAnything");
      }
    };
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, []);

  const handleExplain = () => {
    const input: MessageInput = {
      content: "Please explain this to me.",
      authorId: editor?.authorId || "",
      allowDraftEdits: false,
      contentAddress: editor?.currentContentAddress() || "",
    };

    const selection = editor?.aiMessageSelection();
    if (selection) {
      input.selection = selection;
    }

    analytics.track(AI_CLICK_EXPLAIN, {});

    editor?.clearSelection();

    createThreadMessage(input, editor?.currentContentAddress() || "");
    setSidebarMode("chat");
    setMenuView("default");
    toggleMaximize();
    onDismiss();
  };

  return (
    <div
      ref={(node: HTMLDivElement) => {
        if (typeof floatingRef === "function") {
          floatingRef(node);
        } else if (floatingRef) {
          floatingRef.current = node;
        }
        setContainer(node);
      }}
      style={style}
      className={`mt-2 rounded-lg bg-card z-[900]
      border border-border
      shadow-[0_2px_4px_-2px_hsla(220,43%,11%,0.1),0_4px_6px_-1px_hsla(0,0%,0%,0.1)] ${hide ? "animate-fadeOut" : "animate-fadeIn"}`}
    >
      <div className="flex items-center">
        {(() => {
          switch (menuView) {
            case "link":
              return (
                <LinkEditor
                  onBlur={() => {
                    setMenuView("default");
                    onDismiss();
                  }}
                />
              );
            case "comment":
              return (
                <CommentEditor
                  container={container}
                  onCommentUpdate={(message: string) => {
                    setHasCommentInput(message.length > 0);
                  }}
                  onCreateTimelineMessage={() => {
                    setHasCommentInput(false);
                    setMenuView("default");
                    toggleMaximize();
                    onDismiss();
                  }}
                />
              );
            case "askAIAnything":
              return (
                <AskAIAnything
                  container={container}
                  onSendMessage={() => {
                    setMenuView("default");
                    toggleMaximize();
                    onDismiss();
                  }}
                />
              );
            case "scrub":
              return <Scrub container={container} />;
            default:
              return (
                <>
                  <Button
                    className={`border-none p-0 h-[calc(2.25rem-2px)] px-2 ${
                      hasLink ? "bg-accent text-forground" : "hover:bg-elevated"
                    }`}
                    size="sm"
                    variant="ghost"
                    onClick={() => {
                      setMenuView("askAIAnything");
                    }}
                  >
                    <Edit2Icon className="w-4 h-4 mr-2" />
                    Edit
                  </Button>
                  <Button
                    className={`border-none p-0 h-[calc(2.25rem-2px)] px-2 ${
                      hasLink ? "bg-accent text-forground" : "hover:bg-elevated"
                    }`}
                    size="sm"
                    variant="ghost"
                    onClick={handleExplain}
                  >
                    <HelpCircleIcon className="w-4 h-4 mr-2" />
                    Explain
                  </Button>

                  <AskReviso
                    container={container}
                    toggleMaximize={toggleMaximize}
                    isDisconnected={isDisconnected}
                  />

                  <Button
                    className={`border-none p-0 h-[calc(2.25rem-2px)] px-2 ${
                      hasLink ? "bg-accent text-forground" : "hover:bg-elevated"
                    }`}
                    size="sm"
                    variant="ghost"
                    onClick={() => {
                      setMenuView("comment");
                    }}
                    disabled={isDisconnected}
                  >
                    <MessageCircleIcon className="w-4 h-4 mr-2" />
                    Comment
                  </Button>
                  <Button
                    className={`border-none p-0 h-[calc(2.25rem-2px)] px-2 ${
                      hasLink ? "bg-accent text-forground" : "hover:bg-elevated"
                    }`}
                    size="sm"
                    variant="ghost"
                    onClick={() => {
                      setMenuView("link");
                    }}
                  >
                    <LinkIcon className="w-4 h-4 mr-2" />
                    Link
                  </Button>
                  <Button
                    className={`border-none p-0 h-[calc(2.25rem-2px)] px-2 ${
                      hasLink ? "bg-accent text-forground" : "hover:bg-elevated"
                    }`}
                    size="sm"
                    variant="ghost"
                    onClick={() => {
                      setMenuView("scrub");
                    }}
                  >
                    <HistoryIcon className="w-4 h-4 mr-2" />
                    Scrub
                  </Button>
                </>
              );
          }
        })()}
      </div>
    </div>
  );
});
