import React, {
  ReactElement,
  cloneElement,
  useState,
  useEffect,
  useRef,
} from "react";
import ScrubSliderWrapper from "./ScrubSlider";
import { Button } from "@/components/ui/button";
import { useSignals } from "@preact/signals-react/runtime";
import { Signal, signal } from "@preact/signals-react";
import {
  ChevronDownIcon,
  ChevronUpIcon,
  CodeIcon,
  SquareCodeIcon,
  Heading1Icon,
  Heading2Icon,
  Heading3Icon,
  ListIcon,
  ListOrderedIcon,
  QuoteIcon,
  TypeIcon,
  BoldIcon,
  ItalicIcon,
  StrikethroughIcon,
  UnderlineIcon,
  Undo2Icon,
  Redo2Icon,
  HistoryIcon,
  AudioLines,
} from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { analytics } from "@/lib/segment";
import {
  DOCUMENT_APPLY_BLOCK_FORMAT,
  DOCUMENT_TOGGLE_FORMAT,
} from "@/lib/events";
import { cn } from "@/lib/utils";
import { Spinner } from "@/components/ui/spinner";
import { useVoiceMode } from "@/contexts/VoiceModeContext";
import { useDocumentContext } from "@/contexts/DocumentContext";
import { useChatContext } from "@/contexts/ChatContext";
import { WithTooltip } from "../ui/FloatingTooltip";

export type FormatBarProps = {
  hideBar: any;
  spanFormat: any;
  lineFormat: any;
  toggleSpanFormat: any;
  setLineFormat: any;
};

interface ButtonGroupProps {
  hideBar: Signal<boolean>;
  children: React.ReactNode;
}

function ButtonGroup({ hideBar, children }: ButtonGroupProps) {
  return (
    <div
      className={`flex bg-card items-center print:hidden ${hideBar.value ? "animate-slideOutY" : "animate-slideInY"} rounded-lg border border-border shadow-[0_2px_4px_-2px_hsla(220,43%,11%,0.1),0_4px_6px_-1px_hsla(0,0%,0%,0.1)]`}
    >
      {children}
    </div>
  );
}

export function FormatBar({
  hideBar,
  spanFormat,
  lineFormat,
  toggleSpanFormat,
  setLineFormat,
}: FormatBarProps) {
  useSignals();

  return (
    <ButtonGroup hideBar={hideBar}>
      <FormatSelect lineFormat={lineFormat} setLineFormat={setLineFormat} />
      <FormatButton
        icon={<BoldIcon />}
        isActive={spanFormat.value["b"] === "true"}
        onClick={() => {
          toggleSpanFormat("b");
        }}
      />
      <FormatButton
        icon={<ItalicIcon />}
        isActive={spanFormat.value["i"] === "true"}
        onClick={() => {
          toggleSpanFormat("i");
        }}
      />
      <FormatButton
        icon={<UnderlineIcon />}
        isActive={spanFormat.value["u"] === "true"}
        onClick={() => {
          toggleSpanFormat("u");
        }}
      />
      <FormatButton
        icon={<StrikethroughIcon />}
        isActive={spanFormat.value["s"] === "true"}
        onClick={() => {
          toggleSpanFormat("s");
        }}
      />
      <FormatButton
        icon={<CodeIcon />}
        isActive={spanFormat.value["c"] === "true"}
        onClick={() => {
          toggleSpanFormat("c");
        }}
      />
    </ButtonGroup>
  );
}

export function HistoryButtons({
  hideBar,
  canUndo,
  canRedo,
  onUndo,
  onRedo,
  onScrub,
}: {
  hideBar: any;
  canUndo: Signal<boolean>;
  canRedo: Signal<boolean>;
  onUndo: () => void;
  onRedo: () => void;
  onScrub: () => void;
}) {
  useSignals();
  return (
    <ButtonGroup hideBar={hideBar}>
      <FormatButton
        icon={<Undo2Icon />}
        isActive={false}
        onClick={onUndo}
        disabled={!canUndo.value}
      />
      <FormatButton
        icon={<Redo2Icon />}
        isActive={false}
        onClick={onRedo}
        disabled={!canRedo.value}
      />
      <FormatButton
        icon={<HistoryIcon />}
        isActive={false}
        onClick={onScrub}
        disabled={false}
      />
    </ButtonGroup>
  );
}

const FormatButton = function ({
  disabled = false,
  icon,
  isActive,
  onClick,
}: {
  isActive: boolean;
  onClick: () => void;
  disabled?: boolean;
  icon: ReactElement<any>;
}) {
  useSignals();
  return (
    <Button
      className={`border-none p-0 w-7 sm:w-9 h-[calc(2.25rem-2px)] rounded-none ${
        isActive ? "bg-accent text-forground" : "hover:bg-elevated"
      }`}
      size="sm"
      variant="ghost"
      disabled={disabled}
      onClick={() => {
        if (onClick) {
          onClick();
          return;
        }
      }}
    >
      {icon && cloneElement(icon, { ...icon.props, className: "w-4 h-4" })}
    </Button>
  );
};

export function VoiceButtons({
  hideBar,
  documentId,
  threadId,
  authorId,
  refreshMessages,
}: {
  hideBar: any;
  documentId: string;
  threadId: string;
  authorId: string;
  refreshMessages: () => void;
}) {
  const { streamingState, connectConversation, disconnectConversation } =
    useVoiceMode();

  const hasAllIds = documentId !== "" && threadId !== "" && authorId !== "";
  const isStreaming = streamingState !== "idle";
  const isConnecting = streamingState === "connecting";

  const handleClick = () => {
    if (isStreaming) {
      disconnectConversation();
    } else {
      connectConversation({ documentId, threadId, authorId, refreshMessages });
    }
  };

  return (
    <ButtonGroup hideBar={hideBar}>
      <WithTooltip tooltipText="Voice Mode. Chat with Reviso about your document.">
        <Button
          onClick={handleClick}
          variant="ghost"
          className={cn(
            "border-none p-0 w-7 sm:w-9 h-[calc(2.25rem-2px)] rounded-none",
            isStreaming && "bg-orange-500 animate-pulse hover:bg-red-500",
            isConnecting && "bg-blue-500 animate-pulse hover:bg-red-500",
          )}
          disabled={!hasAllIds}
        >
          {isConnecting && <Spinner className="w-4 h-4 max-w-4 min-w-4 mx-0" />}
          {!isConnecting && <AudioLines className={cn("w-4 h-4")} />}
        </Button>
      </WithTooltip>
    </ButtonGroup>
  );
}

const getIconForValue = (value: string) => {
  switch (value) {
    case "header:1":
      return <Heading1Icon className="w-4 h-4" />;
    case "header:2":
      return <Heading2Icon className="w-4 h-4" />;
    case "header:3":
      return <Heading3Icon className="w-4 h-4" />;
    case "list:bullet":
      return <ListIcon className="w-4 h-4" />;
    case "list:ordered":
      return <ListOrderedIcon className="w-4 h-4" />;
    case "blockquote":
      return <QuoteIcon className="w-4 h-4" />;
    case "code-block":
      return <SquareCodeIcon className="w-4 h-4" />;
    default:
      return <TypeIcon className="w-4 h-4" />;
  }
};

const getFormat = (format: Record<string, any> | undefined) => {
  if (!format) {
    return "text";
  }
  if (format.h) {
    return `header:${format.h}`;
  } else if (format.ol) {
    return `list:ordered`;
  } else if (format.ul) {
    return `list:bullet`;
  } else if (format.bq) {
    return "blockquote";
  } else if (typeof format["cb"] === "string") {
    return "code-block";
  } else {
    return "text";
  }
};

export const FormatSelect = function ({
  lineFormat,
  setLineFormat,
}: {
  lineFormat: any;
  setLineFormat: any;
}) {
  useSignals();

  const [isOpen, setIsOpen] = useState(false);

  return (
    <DropdownMenu onOpenChange={setIsOpen}>
      <DropdownMenuTrigger className="border-none rounded hover:h-[calc(100%-2px)] flex items-center text-foreground gap-2 px-2 hover:bg-elevated/90">
        {getIconForValue(getFormat(lineFormat.value))}
        {isOpen ? (
          <ChevronUpIcon className="w-4 h-4" />
        ) : (
          <ChevronDownIcon className="w-4 h-4" />
        )}
      </DropdownMenuTrigger>
      <DropdownMenuContent
        className="ml-[5.5rem] mb-3"
        onCloseAutoFocus={(e: any) => e.preventDefault()}
      >
        <DropdownMenuItem
          className="flex gap-2"
          onClick={() => setLineFormat("text", "true")}
        >
          <TypeIcon className="w-4 h-4" />
          <span>Text</span>
        </DropdownMenuItem>
        <DropdownMenuItem
          className="flex gap-2"
          onClick={() => {
            setLineFormat("h", "1");
          }}
        >
          <Heading1Icon className="w-4 h-4" />
          <span>Heading 1</span>
        </DropdownMenuItem>
        <DropdownMenuItem
          className="flex gap-2"
          onClick={() => setLineFormat("h", "2")}
        >
          <Heading2Icon className="w-4 h-4" />
          <span>Heading 2</span>
        </DropdownMenuItem>
        <DropdownMenuItem
          className="flex gap-2"
          onClick={() => setLineFormat("h", "3")}
        >
          <Heading3Icon className="w-4 h-4" />
          <span>Heading 3</span>
        </DropdownMenuItem>
        <DropdownMenuItem
          className="flex gap-2"
          onClick={() => setLineFormat("ol", "0")}
        >
          <ListOrderedIcon className="w-4 h-4" />
          <span>Numbered List</span>
        </DropdownMenuItem>
        <DropdownMenuItem
          className="flex gap-2"
          onClick={() => setLineFormat("ul", "0")}
        >
          <ListIcon className="w-4 h-4" />
          <span>Bulleted List</span>
        </DropdownMenuItem>
        <DropdownMenuItem
          className="flex gap-2"
          onClick={() => setLineFormat("bq", "true")}
        >
          <QuoteIcon className="w-4 h-4" />
          <span>Quote</span>
        </DropdownMenuItem>
        <DropdownMenuItem
          className="flex gap-2"
          onClick={() => setLineFormat("cb", "true")}
        >
          <SquareCodeIcon className="w-4 h-4" />
          <span>Code</span>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export function FormatBarContainer() {
  const [scrubMode, setScrubMode] = useState(false);

  const { draftId: documentId } = useDocumentContext();
  const { activeThreadID: threadId, refetchMessages } = useChatContext();
  const { editor } = useRogueEditorContext();

  const spanFormat = signal({});
  const lineFormat = signal({});
  const canUndo = signal<boolean>(false);
  const canRedo = signal<boolean>(false);

  const divRef = useRef<HTMLDivElement>(null);
  useEffect(() => {
    if (divRef.current) {
      divRef.current.focus();
    }
  }, [scrubMode]);

  const handleBlur = (event: React.FocusEvent<HTMLDivElement>) => {
    if (!divRef.current?.contains(event.relatedTarget as Node)) {
      setScrubMode(false);
    }
  };

  if (!editor) {
    return;
  }

  const toggleSpanFormat = (format: string) => {
    analytics.track(DOCUMENT_TOGGLE_FORMAT, { format });
    editor.toggleSpanFormat(format);
    spanFormat.value = editor.curSpanFormat;
  };

  const setLineFormat = (format: string, value: string) => {
    analytics.track(DOCUMENT_APPLY_BLOCK_FORMAT, { format: value });
    editor.format(format, value);
    lineFormat.value = editor.curLineFormat;
  };

  editor.subscribe("curSpanFormat", (value: any) => {
    spanFormat.value = value;
  });

  editor.subscribe("curLineFormat", (value: any) => {
    lineFormat.value = value;
  });

  editor.subscribe("canUndo", (value: any) => {
    canUndo.value = value;
  });

  editor.subscribe("canRedo", (value: any) => {
    canRedo.value = value;
  });

  return (
    <div
      className="flex justify-center sticky bottom-6 left-0 right-0 gap-2 sm:gap-4 mx-1 sm:mx-0 sm:bottom-3 sm:ml-0"
      ref={divRef}
      tabIndex={0}
      onBlur={handleBlur}
    >
      {scrubMode && <ScrubSliderWrapper />}

      {!scrubMode && (
        <FormatBar
          hideBar={false}
          spanFormat={spanFormat}
          lineFormat={lineFormat}
          toggleSpanFormat={toggleSpanFormat}
          setLineFormat={setLineFormat}
        />
      )}

      {!scrubMode && (
        <HistoryButtons
          hideBar={false}
          canUndo={canUndo}
          canRedo={canRedo}
          onUndo={() => {
            editor.undo();
          }}
          onRedo={() => {
            editor.redo();
          }}
          onScrub={() => {
            setScrubMode(true);
          }}
        />
      )}

      {editor.authorId && documentId && threadId && !scrubMode && (
        <VoiceButtons
          hideBar={false}
          documentId={documentId}
          authorId={editor.authorId}
          threadId={threadId}
          refreshMessages={refetchMessages}
        />
      )}
    </div>
  );
}
