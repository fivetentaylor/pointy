import React, { useEffect, useMemo, useRef, useState } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { useSignal } from "@preact/signals-react/runtime";
import { ExternalLinkIcon, Link2Icon, Link2OffIcon } from "lucide-react";

type LinkEditorProps = {
  onBlur: () => void;
  currentLink: string;
  setCurrentLink: (link: string) => void;
  removeLink: () => void;
};

function validateLink(link: string) {
  try {
    new URL(link);
    return true;
  } catch {
    return false;
  }
}

function addProtocol(link: string) {
  return link.match(/^https?:\/\//) ? link : `https://${link}`;
}

export const LinkEditor = function ({
  onBlur,
  currentLink,
  setCurrentLink,
  removeLink,
}: LinkEditorProps) {
  const [linkInput, setLinkInput] = useState("");
  const [hasMadeChanges, setHasMadeChanges] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const hasValidLink = useMemo(
    () =>
      hasMadeChanges
        ? validateLink(addProtocol(linkInput))
        : validateLink(currentLink),
    [linkInput, currentLink, hasMadeChanges],
  );

  const onClose = (newLink: string) => {
    setLinkInput(newLink);
    if (newLink) {
      setCurrentLink(addProtocol(linkInput));
    } else {
      removeLink();
    }
    setHasMadeChanges(false);
    onBlur();
  };

  return (
    <>
      <div ref={containerRef} className="relative">
        <Link2Icon size={20} className="fixed stroke-gray-400 top-2 left-2" />
        <Input
          ref={inputRef}
          className="pl-8 h-9 border-0 rounded-none focus:border-0 focus:ring-0 focus-visible:ring-0 focus-visible:ring-offset-0"
          placeholder="Enter URL..."
          type="text"
          value={hasMadeChanges ? linkInput : currentLink}
          onChange={(e) => {
            const newValue = e.target.value;
            if (newValue?.length === 0) {
              removeLink();
            }
            setLinkInput(newValue);
            setHasMadeChanges(true);
          }}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              e.preventDefault();
              onClose(linkInput);
            }
          }}
        />
      </div>
      <div className="w-[1px] bg-border py-1 h-7" />
      <Button
        className={`border-none dark:bg-gray-800 dark:hover:bg-white dark:hover:text-gray-800 px-2 h-9 
        ${linkInput || currentLink ? "" : "opacity-50 cursor-not-allowed"}`}
        variant="outline"
        onClick={() => {
          onClose("");
        }}
      >
        <Link2OffIcon className="h-4 w-4" />
      </Button>
      <Button
        className={`border-none dark:bg-gray-800 dark:hover:bg-white dark:hover:text-gray-800 px-2 h-9 ${
          hasValidLink ? "" : "opacity-50 cursor-not-allowed"
        }`}
        variant="outline"
        onMouseDown={() => {
          if (hasValidLink) {
            window.open(
              hasMadeChanges ? addProtocol(linkInput) : currentLink,
              "_blank",
            );
          }
        }}
      >
        <ExternalLinkIcon size="sm" className="h-4 w-4" />
      </Button>
    </>
  );
};

export const LinkEditorContainer = function (
  props: Pick<LinkEditorProps, "onBlur">,
) {
  const { editor } = useRogueEditorContext();
  const currentLink = useSignal("");

  useEffect(() => {
    if (!editor) {
      return;
    }
    const onFormatChange = (format: any) => {
      currentLink.value = format.a;
    };

    editor.subscribe("curSpanFormat", onFormatChange);
    onFormatChange(editor.curSpanFormat);

    return () => {
      editor.unsubscribe("curSpanFormat", onFormatChange);
    };
  }, [editor]);

  const setCurrentLink = (link: string) => {
    editor?.format("a", link);
  };

  const removeLink = () => {
    editor?.format("a", "");
  };

  return (
    <LinkEditor
      {...props}
      currentLink={currentLink.value}
      setCurrentLink={setCurrentLink}
      removeLink={removeLink}
    />
  );
};

export default LinkEditorContainer;
