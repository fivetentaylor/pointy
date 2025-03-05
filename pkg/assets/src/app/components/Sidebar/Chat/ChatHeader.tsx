import React from "react";
import { Button } from "@/components/ui/button";
import { MessageCirclePlusIcon } from "lucide-react";
import { WithTooltip } from "@/components/ui/FloatingTooltip";
import { ASK_AI_CREATE_NEW_THREAD } from "@/lib/events";
import { analytics } from "@/lib/segment";

type HeaderProps = {
  createThread: () => void;
  loadingCreateThread: boolean;
};

const ChatHeader = ({ createThread, loadingCreateThread }: HeaderProps) => {
  const handleCreateNewThread = () => {
    analytics.track(ASK_AI_CREATE_NEW_THREAD);
    createThread();
  };

  return (
    <div className="flex gap-2 items-center text-foreground text-base font-sans leading-normal text-right">
      <WithTooltip tooltipText="New topic">
        <Button
          variant="icon"
          size="icon"
          disabled={loadingCreateThread}
          onClick={handleCreateNewThread}
        >
          <MessageCirclePlusIcon className="w-4 h-4" />
        </Button>
      </WithTooltip>
    </div>
  );
};

export default ChatHeader;
