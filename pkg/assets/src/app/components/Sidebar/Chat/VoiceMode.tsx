import React from "react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { AudioLines } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import { useVoiceMode } from "@/contexts/VoiceModeContext";
import { WithTooltip } from "@/components/ui/FloatingTooltip";

interface VoiceModeProps {
  documentId: string;
  threadId: string;
  authorId: string;
  refreshMessages: () => void;
}

const VoiceMode: React.FC<VoiceModeProps> = ({
  documentId,
  threadId,
  authorId,
  refreshMessages,
}) => {
  const { streamingState, connectConversation, disconnectConversation } =
    useVoiceMode();

  const hasAllIds = documentId !== "" && threadId !== "" && authorId !== "";
  const isStreaming = streamingState !== "idle";
  const isConnecting = streamingState === "connecting";

  const handleClick = () => {
    if (isStreaming) {
      disconnectConversation();
    } else {
      connectConversation({ documentId, threadId, authorId, refreshMessages });
    }
  };

  return (
    <WithTooltip tooltipText="Voice Mode. Chat with Reviso about your document.">
      <Button
        onClick={handleClick}
        variant="icon"
        className={cn(
          "px-3 mx-1",
          isStreaming && "bg-orange-500 animate-pulse hover:bg-red-500",
          isConnecting && "bg-blue-500 animate-pulse hover:bg-red-500",
        )}
        disabled={!hasAllIds}
      >
        {isConnecting && <Spinner className="w-4 h-4 max-w-4 min-w-4 mx-0" />}
        {!isConnecting && <AudioLines className={cn("w-4 h-4")} />}
      </Button>
    </WithTooltip>
  );
};

export default VoiceMode;
