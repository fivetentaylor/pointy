import React, { useEffect, useRef, useState } from "react";

import {
  useFloating,
  useDismiss,
  useInteractions,
  shift,
  hide,
  inline,
  autoUpdate,
  offset,
  autoPlacement,
} from "@floating-ui/react";
import { BubbleMenu } from "./BubbleMenu";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { CommentDismissAlert } from "@/components/ui/CommentDismissAlert";

export const BubbleMenuWrapper = function ({
  toggleMaximize,
}: {
  toggleMaximize: () => void;
}) {
  const [isOpen, setIsOpen] = useState(false);
  const isMouseDownRef = useRef<boolean>(false);
  const [hasActiveLink, setHasActiveLink] = useState(false);
  const { editor } = useRogueEditorContext();
  const currentSelectionRef = useRef<Selection | null>(null);
  const [showDismissAlert, setShowDismissAlert] = useState(false);
  const [hasCommentInput, setHasCommentInput] = useState(false);

  const { refs, floatingStyles, context, middlewareData } = useFloating({
    open: isOpen,
    onOpenChange: (open) => {
      if (isOpen && !open && hasCommentInput) {
        setShowDismissAlert(true);
      } else {
        setIsOpen(open);
      }
    },
    middleware: [
      inline(),
      autoPlacement({
        allowedPlacements: ["top-start", "top-end"],
      }),
      offset((state) => {
        return state.placement.match("top") ? 14 : -7;
      }),
      /*
      flip({
        fallbackPlacements: ["bottom-start", "bottom-end"],
      }),
      */
      shift({ crossAxis: true }),
      hide({ padding: { top: -15, bottom: -14 } }),
    ],
    whileElementsMounted: (reference, floating, update) => {
      return autoUpdate(reference, floating, update);
    },
  });

  const dismiss = useDismiss(context);

  const { getFloatingProps } = useInteractions([dismiss]);

  useEffect(() => {
    if (!editor || !editor.subscribe) {
      return;
    }

    setIsOpen(false);

    const onFormatChange = (format: any) => {
      if (format && format.a) {
        setHasActiveLink(true);
      } else {
        setHasActiveLink(false);
      }
    };

    const onMouseUp = () => {
      isMouseDownRef.current = false;
      if (currentSelectionRef.current) {
        showMenuIfPossible(currentSelectionRef.current);
      }
    };

    const onMouseDown = () => {
      isMouseDownRef.current = true;
    };

    const showMenuIfPossible = (selection: Selection) => {
      if (!isMouseDownRef.current && selection && !selection.isCollapsed) {
        const domRange =
          typeof selection.rangeCount === "number" && selection.rangeCount > 0
            ? selection.getRangeAt(0)
            : null;
        if (domRange) {
          refs.setReference({
            getBoundingClientRect: () => domRange.getBoundingClientRect(),
            getClientRects: () => domRange.getClientRects(),
          });
          setIsOpen(true);
          return true;
        }
      }
      return false;
    };

    const onSelectionChange = (value: string | null) => {
      const selection = window.getSelection();

      const isPartialRewind =
        editor.editorMode === "scrub" && editor.scrubMode === "partial";

      if (
        !isPartialRewind &&
        (!value ||
          !selection ||
          !editor.container ||
          !editor.container.contains(selection.anchorNode) ||
          selection.isCollapsed ||
          (selection.anchorOffset === selection.focusOffset &&
            selection.toString().trim() === "") ||
          editor.address)
      ) {
        setIsOpen(false);
        currentSelectionRef.current = null;
        return;
      }

      if (!showMenuIfPossible(selection)) {
        // If we couldn't show the menu immediately, store the selection
        // and wait for mouseup
        currentSelectionRef.current = selection;
      }
    };

    editor.subscribe("curSpanFormat", onFormatChange);
    editor.addEventListener("mousedown", onMouseDown, { capture: true });
    // Listen for mouseup on window to catch all cases
    window.addEventListener("mouseup", onMouseUp, { capture: true });
    editor.subscribe<string | null>("selectedHtml", onSelectionChange);

    return () => {
      editor.unsubscribe("curSpanFormat", onFormatChange);
      editor.removeEventListener("mousedown", onMouseDown);
      window.removeEventListener("mouseup", onMouseUp);
      editor.unsubscribe("selectedHtml", onSelectionChange);
    };
  }, [editor]);

  return (
    <>
      {isOpen && (
        <BubbleMenu
          ref={(node: HTMLElement | null) => {
            refs.setFloating(node);
          }}
          style={floatingStyles}
          {...getFloatingProps()}
          hide={middlewareData.hide?.referenceHidden || false}
          onDismiss={() => {
            setIsOpen(false);
            editor?.clearSelection();
          }}
          hasLink={hasActiveLink}
          toggleMaximize={toggleMaximize}
          setHasCommentInput={setHasCommentInput}
        ></BubbleMenu>
      )}
      <CommentDismissAlert
        open={showDismissAlert}
        setOpen={setShowDismissAlert}
        onConfirm={() => {
          setIsOpen(false);
          setShowDismissAlert(false);
          setHasCommentInput(false);
        }}
        onCancel={() => {
          setShowDismissAlert(false);
        }}
      />
    </>
  );
};
