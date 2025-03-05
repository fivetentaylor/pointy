import React, { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { CheckCheckIcon, EyeIcon, Undo2Icon, XIcon } from "lucide-react";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { WithTooltip } from "../ui/FloatingTooltip";
import { cn } from "@/lib/utils";

export const HistoryToast = function () {
  const { editor } = useRogueEditorContext();
  const [showDiffHighlights, setShowDiffHighlights] = useState(false);
  const [addressDescription, setAddressDescription] = useState<string | null>(
    null,
  );

  useEffect(() => {
    if (!editor) return;

    editor.subscribe<boolean>("showDiffHighlights", setShowDiffHighlights);
    setShowDiffHighlights(editor.showDiffHighlights);
    editor.subscribe<string>("addressDescription", setAddressDescription);
    setAddressDescription(editor.addressDescription);

    return () => {
      editor.unsubscribe("showDiffHighlights", setShowDiffHighlights);
      editor.unsubscribe("addressDescription", setAddressDescription);
    };
  }, [editor]);

  const handleDone = () => {
    if (!editor?.address && editor?.editorMode === "history") {
      return;
    }

    editor?.resetAddress();
  };

  const handleRevert = () => {
    if (!editor?.address && editor?.editorMode === "history") {
      return;
    }

    editor?.rewindAll();
    editor?.resetAddress();
  };

  if (!editor || !editor.address || editor.editorMode !== "history") {
    return null;
  }

  const handleToggleHighlights = () => {
    if (editor) {
      editor.showDiffHighlights = !editor?.showDiffHighlights;
      editor.renderRogue();
    }
  };

  return (
    <div className="sticky top-0 left-0 w-full z-[9999] pointer-events-auto">
      <div className="absolute top-0 left-0 right-0">
        <div className="flex items-start justify-center fade-in h-[1.75rem]">
          <div className="bg-reviso text-white rounded shadow-lg flex items-center justify-between h-9">
            <div className="flex items-center">
              <span className="flex-grow px-4">
                {addressDescription || "Draft"}
              </span>
              <WithTooltip
                tooltipText={showDiffHighlights ? "Hide Diff" : "Show Diff"}
              >
                <Button
                  className={cn(
                    "hover:bg-primary/80 hover:text-inherit",
                    showDiffHighlights ? "bg-primary" : "bg-reviso",
                  )}
                  variant="ghost"
                  size="sm"
                  onClick={handleToggleHighlights}
                >
                  <EyeIcon className="w-4 h-4" />
                </Button>
              </WithTooltip>
              <WithTooltip tooltipText="Revert back to this draft">
                <Button
                  className="hover:bg-primary/80 hover:text-inherit"
                  variant="ghost"
                  size="sm"
                  onClick={handleRevert}
                >
                  <Undo2Icon className="w-4 h-4" />
                </Button>
              </WithTooltip>
            </div>
            <WithTooltip tooltipText="Back to current draft">
              <Button
                className="hover:bg-primary/80 hover:text-inherit"
                variant="ghost"
                size="sm"
                onClick={handleDone}
              >
                <XIcon className="w-4 h-4" />
              </Button>
            </WithTooltip>
          </div>
        </div>
      </div>
    </div>
  );
};
