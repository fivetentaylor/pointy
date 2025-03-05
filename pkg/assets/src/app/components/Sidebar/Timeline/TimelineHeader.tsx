import React, { useState } from "react";
import {
  ListFilterIcon,
  HistoryIcon,
  MessageSquareIcon,
  TextIcon,
  CheckIcon,
  BellIcon,
} from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownMenuShortcut,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import TimelineNotificationSettingsDialog from "./TimelineNotificationSettings";
import { useTimelineContext } from "./TimelineContext";
import { analytics } from "@/lib/segment";
import { NOTIFICATIONS_SETTINGS_OPEN } from "@/lib/events";

const TimelineHeader = () => {
  const [settingsOpen, setSettingsOpen] = useState(false);
  const { timelineFilter, setTimelineFilter } = useTimelineContext();

  return (
    <div className="flex gap-1 items-center text-foreground text-base font-sans leading-normal text-right">
      <TimelineNotificationSettingsDialog
        showOpen={settingsOpen}
        onClose={() => setSettingsOpen(false)}
      />

      <Button
        variant="icon"
        size="icon"
        className="flex items-center focus:outline-none"
        onClick={() => {
          analytics.track(NOTIFICATIONS_SETTINGS_OPEN);
          setSettingsOpen(true);
        }}
      >
        <BellIcon className="w-4 h-4" />
      </Button>

      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button
            variant="icon"
            size="icon"
            className="flex items-center focus:outline-none"
          >
            <ListFilterIcon className="w-4 h-4" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItem
            onClick={() => {
              setTimelineFilter("ALL");
            }}
          >
            <HistoryIcon className="mr-2 h-4 w-4" />
            <span>All</span>
            <DropdownMenuShortcut className="ml-2">
              {timelineFilter === "ALL" && (
                <CheckIcon className="w-4 h-4 text-muted-foreground" />
              )}
            </DropdownMenuShortcut>
          </DropdownMenuItem>
          <DropdownMenuItem
            onClick={() => {
              setTimelineFilter("COMMENTS");
            }}
          >
            <MessageSquareIcon className="mr-2 h-4 w-4" />
            <span>Comments</span>
            <DropdownMenuShortcut className="ml-2">
              {timelineFilter === "COMMENTS" && (
                <CheckIcon className="w-4 h-4 text-muted-foreground" />
              )}
            </DropdownMenuShortcut>
          </DropdownMenuItem>
          <DropdownMenuItem
            onClick={() => {
              setTimelineFilter("EDITS");
            }}
          >
            <TextIcon className="mr-2 h-4 w-4" />
            <span>Edits</span>
            <DropdownMenuShortcut className="ml-2">
              {timelineFilter === "EDITS" && (
                <CheckIcon className="w-4 h-4 text-muted-foreground" />
              )}
            </DropdownMenuShortcut>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
};

export default TimelineHeader;
