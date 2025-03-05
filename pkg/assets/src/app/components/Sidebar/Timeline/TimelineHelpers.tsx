import React, { useState } from "react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { cn, formatLocalTime, timeSpecific } from "@/lib/utils";
import { getInitials } from "@/lib/utils";
import { WithTooltip } from "@/components/ui/FloatingTooltip";

type TimelineEventUser = TimelineEventFieldsFragment["user"];

interface TimelineDateProps {
  date: string;
}

export const TimelineDate: React.FC<TimelineDateProps> = ({ date }) => {
  return (
    <WithTooltip tooltipText={timeSpecific(date)}>
      <span className="text-muted-foreground">
        {" • "}
        {formatLocalTime(date)}
      </span>
    </WithTooltip>
  );
};

interface TimelineAvatarProps {
  user: TimelineEventUser;
  className?: string;
}

export const TimelineAvatar: React.FC<TimelineAvatarProps> = ({
  user,
  className,
}) => {
  return (
    <Avatar className={cn("w-4 h-4 mt-1 mr-2", className)}>
      <AvatarImage alt="Profile icon" src={user.picture || undefined} />
      <AvatarFallback className="text-background bg-primary">
        {getInitials(user.name || "")}
      </AvatarFallback>
    </Avatar>
  );
};
