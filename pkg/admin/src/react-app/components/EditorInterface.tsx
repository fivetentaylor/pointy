import React, {
  ReactElement,
  cloneElement,
  useState,
} from "react";
import {
  BoldIcon,
  ItalicIcon,
  StrikethroughIcon,
  UnderlineIcon,
  UndoIcon,
  RedoIcon,
} from "lucide-react";
import { Button } from "../ui/button";
import { useSignals } from "@preact/signals-react/runtime";

export const EditorInterface = ({
  hideBar,
  spanFormat,
  lineFormat,
  toggleSpanFormat,
  setLineFormat,
  undo,
  redo,
}: {
  hideBar: any;
  spanFormat: any;
  lineFormat: any;
  toggleSpanFormat: any;
  setLineFormat: any;
  undo: any;
  redo: any;
}) => {
  useSignals();

  return (
    <div
      className={`flex bg-background h-9 w-44 items-center mx-auto print:hidden ${hideBar.value ? "animate-slideOutY" : "animate-slideInY"} rounded-lg border border-border shadow-[0_2px_4px_-2px_hsla(220,43%,11%,0.1),0_4px_6px_-1px_hsla(0,0%,0%,0.1)]`}
    >
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
        icon={<UndoIcon />}
        isActive={false}
        onClick={() => {
          undo();
        }}
      />
      <FormatButton
        icon={<RedoIcon />}
        isActive={false}
        onClick={() => {
          redo();
        }}
      />
    </div>
  );
};

const FormatButton = function({
  icon,
  isActive,
  onClick,
}: {
  isActive: boolean;
  onClick: () => void;
  icon: ReactElement<any>;
}) {
  return (
    <Button
      className={`border-none p-0 w-9 h-[calc(2.25rem-2px)] ${isActive ? "bg-accent text-forground" : "hover:bg-elevated"
        }`}
      size="sm"
      variant="ghost"
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

import {
  ChevronDownIcon,
  ChevronUpIcon,
  SquareCodeIcon,
  Heading1Icon,
  Heading2Icon,
  Heading3Icon,
  ListIcon,
  ListOrderedIcon,
  QuoteIcon,
  TypeIcon,
} from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu";

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

export const FormatSelect = function({
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
        className="ml-8"
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
