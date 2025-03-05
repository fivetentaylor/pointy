import React, { useEffect, useRef, useState, useCallback } from "react";
import {
  Panel,
  PanelGroup,
  PanelResizeHandle,
  ImperativePanelHandle,
} from "react-resizable-panels";

import { RogueEditorContextProvider } from "@/contexts/RogueEditorContext";
import Editor from "@/components/Editor";
import { ChatContextProvider } from "@/contexts/ChatContext";
import { ErrorBoundary } from "@/components/ui/ErrorBoundary";
import { CursorContextProvider } from "@/contexts/CursorContext";
import { DocumentContextProvider } from "@/contexts/DocumentContext";
import { SidebarContextProvider } from "@/contexts/SidebarContext";
import { Sidebar } from "@/components/Sidebar";
import { TimelineProvider } from "@/components/Sidebar/Timeline/TimelineContext";
import { BlockError } from "@/components/ui/BlockError";
import { SidebarProvider, useSidebar } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/AppSidebar";
import { VoiceModeProvider } from "@/contexts/VoiceModeContext";

const CombinedProviders = ({ children }: { children: React.ReactNode }) => (
  <DocumentContextProvider>
    <SidebarProvider>
      <RogueEditorContextProvider>
        <SidebarContextProvider>
          <ChatContextProvider>
            <TimelineProvider>
              <CursorContextProvider>
                <VoiceModeProvider>{children}</VoiceModeProvider>
              </CursorContextProvider>
            </TimelineProvider>
          </ChatContextProvider>
        </SidebarContextProvider>
      </RogueEditorContextProvider>
    </SidebarProvider>
  </DocumentContextProvider>
);

const ChatDocHitAreaMargins = {
  coarse: 20,
  fine: 15,
};

const Drafts = () => {
  return (
    <CombinedProviders>
      <AppSidebar />
      <DraftsContent />
    </CombinedProviders>
  );
};

const DraftsContent = function () {
  const [maximized, setMaximized] = useState(false);
  const [showDocumentsPanel] = useState(true);
  const [sidebarMaxWidth, setSidebarMaxWidth] = useState(0);
  const chatPanelRef = useRef<ImperativePanelHandle>(null);
  const resizeTimeoutRef = useRef<number | null>(null);
  const { open: appSidebarOpen } = useSidebar();
  const animationRef = useRef<number>();

  const handleResize = useCallback(() => {
    if (resizeTimeoutRef.current !== null) {
      cancelAnimationFrame(resizeTimeoutRef.current);
    }

    resizeTimeoutRef.current = requestAnimationFrame(() => {
      const panelGroup = document.getElementById("DraftPanelGroup");
      if (panelGroup && chatPanelRef.current) {
        setSidebarMaxWidth(
          panelGroup.clientWidth * (chatPanelRef.current.getSize() / 100),
        );
      }
      resizeTimeoutRef.current = null;
    });
  }, []);

  useEffect(() => {
    window.addEventListener("resize", handleResize);
    return () => {
      window.removeEventListener("resize", handleResize);
      if (resizeTimeoutRef.current !== null) {
        cancelAnimationFrame(resizeTimeoutRef.current);
      }
    };
  }, [handleResize]);

  useEffect(() => {
    handleResize();
  }, []);

  useEffect(() => {
    let startTime: number;
    const ANIMATION_DURATION = 300;

    const updateDuringAnimation = (timestamp: number) => {
      if (!startTime) startTime = timestamp;
      const elapsed = timestamp - startTime;

      handleResize();

      if (elapsed < ANIMATION_DURATION) {
        animationRef.current = requestAnimationFrame(updateDuringAnimation);
      }
    };

    // Start the animation loop
    animationRef.current = requestAnimationFrame(updateDuringAnimation);

    // Cleanup
    return () => {
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
      }
    };
  }, [appSidebarOpen]);

  return (
    <PanelGroup
      id="DraftPanelGroup"
      /*autoSaveId={draftId}*/
      direction="horizontal"
    >
      <ErrorBoundary fallback={<div>Left Panel Error</div>}>
        {!maximized && (
          <ErrorBoundary
            fallback={
              <div className="max-w-96 px-4 mt-20">
                <BlockError text="The sidebar couldn't be loaded due to an error." />
              </div>
            }
          >
            <Panel
              id="ChatPanel"
              ref={chatPanelRef}
              order={2}
              defaultSize={30}
              minSize={25}
              className="hidden md:block"
              onResize={handleResize}
            >
              <Sidebar
                maxWidth={sidebarMaxWidth}
                documentListShowing={showDocumentsPanel}
              />
            </Panel>
            <PanelResizeHandle hitAreaMargins={ChatDocHitAreaMargins} />
          </ErrorBoundary>
        )}
      </ErrorBoundary>
      <Panel id="editor" order={3} defaultSize={50} minSize={25}>
        <ErrorBoundary fallback={<EditorPanelEditor />}>
          <Editor
            maximized={maximized}
            toggleMaximize={() => {
              setMaximized(!maximized);
              setTimeout(() => window.dispatchEvent(new Event("resize")), 50);
            }}
          />
        </ErrorBoundary>
      </Panel>
    </PanelGroup>
  );
};

const EditorPanelEditor = () => {
  return <BlockError text="Your draft couldn't be loaded due to an error." />;
};

export default Drafts;
