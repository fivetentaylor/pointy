import React, { ReactNode } from "react";
import { CatIcon, FileClockIcon, SparklesIcon } from "lucide-react";
import { PostHogFeature } from "posthog-js/react";
import { SidebarMode } from "./Chat";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { useSidebarContext } from "@/contexts/SidebarContext";
import { useRogueEditorContext } from "@/contexts/RogueEditorContext";
import { Skeleton } from "@/components/ui/skeleton";
import { useParams } from "react-router-dom";
import { analytics } from "@/lib/segment";
import { SIDEBAR_SELECT_CHAT, SIDEBAR_SELECT_TIMELINE } from "@/lib/events";
import { SidebarTrigger } from "../ui/sidebar";
import { useCurrentUserContext } from "@/contexts/CurrentUserContext";

type HeaderProps = {
  documentId: string;
  children: ReactNode;
  onUpdateSidebarMode: (mode: SidebarMode) => void;
  currentSidebarMode: SidebarMode;
  me: User | null;
};

const LoadingHeader = () => {
  return (
    <div className="flex items-center h-9 my-2">
      <div className="flex text-foreground text-base font-sans leading-normal flex-grow">
        <div className="flex items-center font-medium gap-2">
          <Skeleton className="w-4 h-4" />
        </div>
      </div>
    </div>
  );
};

const Header = ({
  children,
  documentId,
  onUpdateSidebarMode,
  currentSidebarMode,
  me,
}: HeaderProps) => {
  const handleButtonClick = (mode: SidebarMode) => {
    analytics.track(
      mode === "timeline" ? SIDEBAR_SELECT_TIMELINE : SIDEBAR_SELECT_CHAT,
    );
    onUpdateSidebarMode(mode);
  };

  return (
    <div className="flex items-center h-9 my-2">
      <div className="flex text-foreground text-base font-sans leading-normal flex-grow">
        <div className="flex items-center font-medium gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => handleButtonClick("timeline")}
            className={cn(
              "shadow-sm px-4 py-2",
              currentSidebarMode === "timeline"
                ? "bg-card border-foreground/40"
                : "",
            )}
          >
            <FileClockIcon className="w-4 h-4 mr-2" />
            Timeline
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => handleButtonClick("chat")}
            className={cn(
              "shadow-sm px-4 py-2",
              currentSidebarMode === "chat"
                ? "bg-card border-foreground/40"
                : "",
            )}
          >
            <SparklesIcon className="w-4 h-4 mr-2" />
            Ask AI
          </Button>
          {me && me.isAdmin && (
            <PostHogFeature flag="show-admin-link" match={true}>
              <div className="ml-2">
                <a
                  href={"/admin/documents/" + documentId + "/dags"}
                  target="_blank"
                  rel="noreferrer"
                  className="text-muted-foreground hover:text-reviso"
                >
                  <CatIcon className="w-4 h-4" />
                </a>
              </div>
            </PostHogFeature>
          )}
        </div>
      </div>

      {children}
    </div>
  );
};

const HeaderWrapper = ({ children }: { children: ReactNode }) => {
  const { currentUser } = useCurrentUserContext();
  const { editor } = useRogueEditorContext();
  const { draftId } = useParams();
  const {
    sidebarMode,
    sidebarLoading,
    setSidebarMode: _setSidebarMode,
  } = useSidebarContext();

  const setSidebarMode = (mode: SidebarMode) => {
    _setSidebarMode(mode);
    if (editor && mode === "chat") {
      editor.hideCommentHighlights();
    }
  };

  if (!sidebarMode || sidebarLoading || !draftId) {
    return <LoadingHeader></LoadingHeader>;
  }

  return (
    <Header
      documentId={draftId}
      onUpdateSidebarMode={setSidebarMode}
      currentSidebarMode={sidebarMode}
      me={currentUser}
    >
      {children}
    </Header>
  );
};

export default HeaderWrapper;
