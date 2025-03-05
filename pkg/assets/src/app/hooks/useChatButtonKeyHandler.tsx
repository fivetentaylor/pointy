import { useState, useCallback } from "react";

interface UseChatButtonKeyHandlerProps {
  isEnabled: boolean;
  sendMessage: (type: "chat" | "revise") => void;
}

export const useChatButtonKeyHandler = ({
  isEnabled,
  sendMessage,
}: UseChatButtonKeyHandlerProps) => {
  const [altKeyDown, setAltKeyDown] = useState(false);

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent | KeyboardEvent) => {
      if (!e.altKey && (e.shiftKey || e.key !== "Enter")) {
        return;
      }

      e.preventDefault();

      if (e.altKey) {
        setAltKeyDown(true);
      }

      if (e.key === "Enter" && isEnabled) {
        if (e.altKey) {
          sendMessage("chat");
        } else {
          sendMessage("revise");
        }
      }

      return false;
    },
    [isEnabled, sendMessage],
  );

  const handleKeyUp = useCallback((e: React.KeyboardEvent | KeyboardEvent) => {
    if (e.key === "Alt") {
      setAltKeyDown(false);
    }
  }, []);

  return { altKeyDown, handleKeyDown, handleKeyUp };
};
