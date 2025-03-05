import React from "react";
import RogueEditor from "./RogueEditor";
import TitleContainer from "./TitleContainer";
import { cn } from "@/lib/utils";
import { FormatBarContainer } from "./FormatBar";
import { BubbleMenuWrapper } from "./BubbleMenu/BubbleMenuWrapper";
import { RevisionToast } from "./RevisionToast";
import { HistoryToast } from "./HistoryToast";
import { ErrorBoundary } from "../ui/ErrorBoundary";
import { BlockError } from "../ui/BlockError";
import { Stats } from "./Stats";

type EditorProps = {
  maximized: boolean;
  toggleMaximize: () => void;
};

const Editor = ({ maximized, toggleMaximize }: EditorProps) => {
  return (
    <div
      className={cn(
        "h-[calc(100dvh-1rem)] m-2 bg-card rounded-md shadow border border-border pt-4",

        maximized ? "p-4" : "pl-4 sm:pl-10",
      )}
    >
      <TitleContainer maximized={maximized} toggleMaximize={toggleMaximize} />

      <ErrorBoundary
        fallback={
          <BlockError text="Your draft couldn't be loaded due to an error." />
        }
      >
        <div
          className={cn(
            "overflow-auto h-[calc(100dvh-4.625rem)] relative pt-[0.625rem] mt-[-0.625rem]",
            maximized ? "px-[calc(60%-38rem)]" : "pr-10",
          )}
        >
          <Stats />
          <RevisionToast />
          <HistoryToast />
          <RogueEditor />
          <FormatBarContainer />
          <BubbleMenuWrapper
            toggleMaximize={() => {
              if (maximized) {
                toggleMaximize();
              }
            }}
          />
        </div>
      </ErrorBoundary>
    </div>
  );
};

export default Editor;
