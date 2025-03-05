import React, { useCallback, useState } from "react";
import {
  ChevronDownIcon,
  ChevronUpIcon,
  SparklesIcon,
  SpellCheckIcon,
} from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useChatContext } from "@/contexts/ChatContext";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { analytics } from "@/lib/segment";
import { AI_CLICK_ASK_REVISO } from "@/lib/events";
import { cn } from "@/lib/utils";

const PROMPT_MAP = {
  add_clarity: "Please rewrite the selected section to add clarity.",
  shorten:
    "Please write a concise version of the selected section while maintaining its core message",
  expand:
    "Please expand the selected section while maintaining my existing tone of voice.",
  use_stronger_language:
    "Please rewrite the selected section to use stronger language.",
  soften_tone: "Please rewrite the selected section in a softer tone.",
  simplify_language:
    "Please simplify the language and structure of the selected section.",
  fix_mistakes:
    "Please rewrite the selected section with fixes to any grammar, punctuation, or spelling mistakes. Only fix mistakes, do not change the content. If there's no mistakes, please let me know.",
};

export const AskReviso = function ({
  container,
  toggleMaximize,
  isDisconnected,
}: {
  container?: HTMLDivElement | null;
  toggleMaximize: () => void;
  isDisconnected: boolean;
}) {
  const [isOpen, setIsOpen] = useState(false);
  const { createThreadMessage } = useChatContext();
  const { editor } = useRogueEditorContext();

  const sendMessage = useCallback(
    async (message: string) => {
      toggleMaximize();

      const input: MessageInput = {
        content: message,
        authorId: editor?.authorId || "",
        allowDraftEdits: true,
      };

      const selection = editor?.aiMessageSelection();
      if (selection) {
        input.selection = selection;
      }

      analytics.track(AI_CLICK_ASK_REVISO, {
        prompt: message,
      });

      return createThreadMessage(input, editor?.currentContentAddress() || "");
    },
    [createThreadMessage, editor, toggleMaximize],
  );

  return (
    <DropdownMenu onOpenChange={setIsOpen}>
      <DropdownMenuTrigger
        className={cn(
          `border-none rounded w-28 h-9 flex items-center text-foreground gap-2 px-2 focus-visible:ring-0 focus-visible:ring-offset-0`,
          isDisconnected
            ? "cursor-not-allowed opacity-50"
            : "hover:bg-elevated/90 ",
        )}
        disabled={isDisconnected}
      >
        <SpellCheckIcon className="w-4 h-4 min-w-4" />
        <span>Improve</span>
        {isOpen ? (
          <ChevronUpIcon className="w-4 h-4" />
        ) : (
          <ChevronDownIcon className="w-4 h-4" />
        )}
      </DropdownMenuTrigger>
      <DropdownMenuContent container={container || undefined} className="ml-16">
        <DropdownMenuItem
          onClick={() => {
            sendMessage(PROMPT_MAP["fix_mistakes"]);
          }}
        >
          Fix mistakes
        </DropdownMenuItem>
        <DropdownMenuItem
          onClick={() => {
            sendMessage(PROMPT_MAP["add_clarity"]);
          }}
        >
          Add clarity
        </DropdownMenuItem>
        <DropdownMenuItem
          onClick={() => {
            sendMessage(PROMPT_MAP["shorten"]);
          }}
        >
          Shorten
        </DropdownMenuItem>
        <DropdownMenuItem
          onClick={() => {
            sendMessage(PROMPT_MAP["expand"]);
          }}
        >
          Expand
        </DropdownMenuItem>
        <DropdownMenuItem
          onClick={() => {
            sendMessage(PROMPT_MAP["use_stronger_language"]);
          }}
        >
          Use stronger language
        </DropdownMenuItem>
        <DropdownMenuItem
          onClick={() => {
            sendMessage(PROMPT_MAP["soften_tone"]);
          }}
        >
          Soften tone
        </DropdownMenuItem>
        <DropdownMenuItem
          onClick={() => {
            sendMessage(PROMPT_MAP["simplify_language"]);
          }}
        >
          Simplify language
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};
