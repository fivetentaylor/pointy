import React from "react";
import Timeline from "./Timeline";
import { useSidebarContext } from "@/contexts/SidebarContext";
import Chat from "./Chat";
import { Skeleton } from "../ui/skeleton";

type SidebarProps = {
  maxWidth: number;
  documentListShowing: boolean;
};

export type SidebarMode = "timeline" | "chat";

export const Sidebar = ({ maxWidth }: SidebarProps) => {
  const { sidebarMode } = useSidebarContext();

  return (
    <div className="h-[calc(100dvh-1rem)] ml-2 pt-4 pl-2">
      <div className="h-full flex">
        <div
          id="SidebarContainer"
          className="flex flex-col flex-grow h-full pl-4 pr-6"
          style={{
            minWidth: `${maxWidth}px`,
            maxWidth: `${maxWidth}px`,
          }}
        >
          {sidebarMode === "timeline" && <Timeline />}
          {sidebarMode === "chat" && <Chat />}
          {!sidebarMode && (
            <div className="h-full mt-2">
              <div className="flex w-full gap-2">
                <Skeleton className="h-9 w-28 bg-elevated" />
                <Skeleton className="h-9 w-20 bg-elevated" />
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
