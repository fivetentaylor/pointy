// src/utils/eventEmitter.ts
import { EventEmitter } from "events";

interface WebSocketStatusEvent {
  open: boolean;
}

class WebSocketEventEmitter extends EventEmitter {
  public currentStatus: WebSocketStatusEvent = { open: true };

  emit(event: "wsStatus", data: WebSocketStatusEvent): boolean {
    this.currentStatus = data;
    return super.emit(event, data);
  }

  on(event: "wsStatus", listener: (data: WebSocketStatusEvent) => void): this {
    return super.on(event, listener);
  }

  off(event: "wsStatus", listener: (data: WebSocketStatusEvent) => void): this {
    return super.off(event, listener);
  }
}

export const wsEventEmitter = new WebSocketEventEmitter();
