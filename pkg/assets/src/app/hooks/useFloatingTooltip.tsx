import React, { useRef, useState } from "react";
import type { Placement } from "@floating-ui/react";
import {
  useFloating,
  arrow,
  offset,
  shift,
  flip,
  useTransitionStyles,
  autoUpdate,
  useHover,
  useInteractions,
} from "@floating-ui/react";

const ARROW_WIDTH = 12;
const ARROW_HEIGHT = 10;

export const useFloatingTooltip = function ({
  placement,
}: {
  placement?: Placement;
} = {}) {
  const [isOpen, setIsOpen] = useState(false);
  const arrowRef = useRef(null);

  const { refs, floatingStyles, context, middlewareData } = useFloating({
    placement: placement || "top",
    open: isOpen,
    onOpenChange: setIsOpen,
    middleware: [
      offset(ARROW_HEIGHT),
      flip({ padding: 5 }),
      shift({ padding: 5 }),
      arrow({ element: arrowRef }),
    ],
    whileElementsMounted: autoUpdate,
  });

  const arrowX = middlewareData.arrow?.x ?? 0;
  const arrowY = middlewareData.arrow?.y ?? 0;
  const transformX = arrowX + ARROW_WIDTH / 2;
  const transformY = arrowY + ARROW_HEIGHT;

  const { styles } = useTransitionStyles(context, {
    initial: {
      transform: "scale(0)",
    },
    common: ({ side }) => ({
      transformOrigin: {
        top: `${transformX}px calc(100% + ${ARROW_HEIGHT}px)`,
        bottom: `${transformX}px ${-ARROW_HEIGHT}px`,
        left: `calc(100% + ${ARROW_HEIGHT}px) ${transformY}px`,
        right: `${-ARROW_HEIGHT}px ${transformY}px`,
      }[side],
    }),
  });

  const hover = useHover(context, {
    delay: {
      open: 500,
      close: 0,
    },
  });
  const { getReferenceProps } = useInteractions([hover]);

  return {
    setReference: refs.setReference,
    setFloating: refs.setFloating,
    getReferenceProps,
    arrowRef,
    context,
    transitionStyles: styles,
    floatingStyles,
    isOpen,
    arrowHeight: ARROW_HEIGHT,
    arrowWidth: ARROW_WIDTH,
  };
};
