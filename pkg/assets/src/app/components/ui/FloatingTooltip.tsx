import React, { Children, ReactElement, cloneElement, forwardRef } from "react";
import { useFloatingTooltip } from "@/hooks/useFloatingTooltip";
import { FloatingArrow, Placement } from "@floating-ui/react";

type FloatingTooltipProps = {
  isOpen: boolean;
  tooltipText: string | undefined;
  setFloating: any;
  floatingStyles: any;
  transitionStyles: any;
  arrowRef: any;
  context: any;
  arrowWidth: number;
  arrowHeight: number;
};

export const FloatingTooltip = function ({
  isOpen,
  setFloating,
  floatingStyles,
  transitionStyles,
  tooltipText,
  arrowRef,
  context,
  arrowWidth,
  arrowHeight,
}: FloatingTooltipProps) {
  return (
    isOpen && (
      <div ref={setFloating} style={floatingStyles} className="z-10">
        <div
          style={transitionStyles}
          className="bg-black dark:bg-white text-white dark:text-black h-9 border-1 rounded-md border-black mx-auto w-max min-w-[2rem] flex items-center justify-center px-3"
        >
          {tooltipText}
          <FloatingArrow
            ref={arrowRef}
            context={context}
            width={arrowWidth}
            height={arrowHeight}
            className="
    fill-black dark:fill-white
    [&>path:first-of-type]:black dark:[&>path:first-of-type]:white
    [&>path:last-of-type]:black dark:[&>path:last-of-type]:white
  "
          />
        </div>
      </div>
    )
  );
};

export const WithTooltip = forwardRef(function WithTooltip(
  {
    tooltipText,
    placement,
    children,
  }: {
    tooltipText?: string;
    placement?: Placement;
    children: ReactElement<any>;
  },
  forwardedRef: React.Ref<any>,
) {
  const {
    setReference,
    getReferenceProps,
    arrowRef,
    context,
    floatingStyles,
    transitionStyles,
    isOpen,
    setFloating,
  } = useFloatingTooltip({ placement });

  // Combined ref to handle both forwardedRef and setReference
  const combinedRef = (node: HTMLElement | null) => {
    setReference(node);
    if (typeof forwardedRef === "function") {
      forwardedRef(node);
    } else if (forwardedRef) {
      // @ts-expect-error: TypeScript does not recognize forwardedRef as a valid ref type
      forwardedRef.current = node;
    }
  };

  return (
    <>
      {cloneElement(children, {
        ref: combinedRef,
        ...getReferenceProps(),
        ...children.props,
        children: [
          ...Children.toArray(children.props.children),
          <FloatingTooltip
            key="tooltip"
            isOpen={isOpen}
            setFloating={setFloating}
            floatingStyles={floatingStyles}
            transitionStyles={transitionStyles}
            tooltipText={tooltipText}
            arrowRef={arrowRef}
            context={context}
            arrowWidth={0}
            arrowHeight={0}
          />,
        ],
      })}
    </>
  );
});
