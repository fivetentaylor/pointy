import { wsEventEmitter } from "@/lib/wsEmitter";
import { useEffect, useState } from "react";

export const useWsDisconnect = () => {
  const [isDisconnected, setIsDisconnected] = useState(
    wsEventEmitter.currentStatus.open === false,
  );

  useEffect(() => {
    const onWsStatusChange = (status: { open: boolean }) => {
      setIsDisconnected(status.open === false);
    };

    wsEventEmitter.on("wsStatus", onWsStatusChange);

    return () => {
      wsEventEmitter.off("wsStatus", onWsStatusChange);
    };
  }, []);

  return { isDisconnected };
};
