import React, {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
} from "react";
import { WavRecorder, WavStreamPlayer } from "@/lib/wavtools";
import { useErrorToast } from "@/hooks/useErrorToast";

type StreamingState = "idle" | "connecting" | "streaming";

interface VoiceModeContextValue {
  streamingState: StreamingState;
  connectConversation: (params: {
    documentId: string;
    threadId: string;
    authorId: string;
    refreshMessages: () => void;
  }) => Promise<void>;
  disconnectConversation: () => Promise<void>;
}

const VoiceModeContext = createContext<VoiceModeContextValue | undefined>(
  undefined,
);

const VoiceModeProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const showErrorToast = useErrorToast();
  const [streamingState, setStreamingState] = useState<StreamingState>("idle");
  const websocketRef = useRef<WebSocket | null>(null);
  const wavRecorderRef = useRef<WavRecorder>(
    new WavRecorder({ sampleRate: 24000 }),
  );
  const wavStreamPlayerRef = useRef<WavStreamPlayer>(
    new WavStreamPlayer({ sampleRate: 24000 }),
  );

  const connectConversation = useCallback(
    async ({
      documentId,
      threadId,
      authorId,
      refreshMessages,
    }: {
      documentId: string;
      threadId: string;
      authorId: string;
      refreshMessages: () => void;
    }) => {
      const serverUrl = `/api/v1/documents/${documentId}/threads/${threadId}/authors/${authorId}/stream`;
      const wavRecorder = wavRecorderRef.current;
      const wavStreamPlayer = wavStreamPlayerRef.current;

      await wavRecorder.begin();
      await wavStreamPlayer.connect();

      const websocket = new WebSocket(serverUrl);
      websocketRef.current = websocket;
      websocket.onopen = () => {
        console.log("Realtime WebSocket connection established", serverUrl);
        setStreamingState("connecting");
      };

      websocket.onerror = (error) => {
        console.error("WebSocket error:", error);
        disconnectConversation();
      };

      websocket.onclose = () => {
        console.log("Realtime WebSocket connection closed");
        if (streamingState !== "idle") {
          showErrorToast("Realtime server disconnected");
        }
        disconnectConversation();
      };

      websocket.onmessage = (event) => {
        if (typeof event.data === "string") {
          const payload = JSON.parse(event.data);
          if (payload.is_speaking) {
            wavStreamPlayer.interrupt();
            return;
          }

          if (payload.type === "response.audio.delta") {
            const arrayBuffer = base64ToArrayBuffer(payload.delta);
            const appendValues = new Int16Array(arrayBuffer);

            wavStreamPlayer.add16BitPCM(appendValues, payload.item_id);
            return;
          }

          if (payload.type === "connected") {
            setStreamingState("streaming");
            return;
          }

          if (payload.type === "new_message") {
            refreshMessages();
            return;
          }

          if (payload.type === "failure") {
            disconnectConversation();
            showErrorToast("Realtime server failure:" + payload.reason);
            return;
          }

          console.log("Unexpected text message from server:", event.data);
        }
      };

      await wavRecorder.record((data) => websocket.send(data.mono));
    },
    [],
  );

  const disconnectConversation = useCallback(async () => {
    if (!websocketRef.current) {
      return;
    }

    setStreamingState("idle");
    console.log("Realtime WebSocket disconnecting");

    const wavRecorder = wavRecorderRef.current;
    await wavRecorder.end();

    const wavStreamPlayer = wavStreamPlayerRef.current;
    wavStreamPlayer.interrupt();

    const websocket = websocketRef.current;
    websocket.close();
    websocketRef.current = null;
  }, []);

  useEffect(() => {
    return () => {
      disconnectConversation();
    };
  }, [disconnectConversation]);

  return (
    <VoiceModeContext.Provider
      value={{ streamingState, connectConversation, disconnectConversation }}
    >
      {children}
    </VoiceModeContext.Provider>
  );
};

const useVoiceMode = (): VoiceModeContextValue => {
  const context = useContext(VoiceModeContext);
  if (!context) {
    throw new Error("useVoiceMode must be used within a VoiceModeProvider");
  }
  return context;
};

function base64ToArrayBuffer(base64: string) {
  const binaryString = atob(base64);
  const len = binaryString.length;
  const bytes = new Uint8Array(len);
  for (let i = 0; i < len; i++) {
    bytes[i] = binaryString.charCodeAt(i);
  }
  return bytes.buffer;
}

export { VoiceModeProvider, useVoiceMode };
