import { Message } from '@/types';

const WS_BASE = import.meta.env.DEV ? 'ws://localhost:8080/ws' : '/ws';

interface ReconnectCallbacks {
  onReconnecting?: (attempt: number, max: number) => void;
  onReconnected?: () => void;
  onReconnectFailed?: () => void;
}

export class WebSocketService {
  private static instance: WebSocketService;
  private ws: WebSocket | null = null;
  private handlers: Map<string, ((payload: any) => void)[]> = new Map();
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000;
  private wasReconnecting = false;

  private onReconnecting?: (attempt: number, max: number) => void;
  private onReconnected?: () => void;
  private onReconnectFailed?: () => void;

  static getInstance(): WebSocketService {
    if (!WebSocketService.instance) {
      WebSocketService.instance = new WebSocketService();
    }
    return WebSocketService.instance;
  }

  setCallbacks(cbs: ReconnectCallbacks) {
    this.onReconnecting = cbs.onReconnecting;
    this.onReconnected = cbs.onReconnected;
    this.onReconnectFailed = cbs.onReconnectFailed;
  }

  connect(token: string, roomId?: string): Promise<void> {
    return new Promise((resolve, reject) => {
      try {
        const url = `${WS_BASE}?token=${token}`;
        this.ws = new WebSocket(url);

        this.ws.onopen = () => {
          console.log('WebSocket connected');
          if (this.wasReconnecting) {
            this.onReconnected?.();
            this.wasReconnecting = false;
          }
          this.reconnectAttempts = 0;
          if (roomId) {
            this.emit('join_room', { room_id: roomId });
          }
          resolve();
        };

        this.ws.onmessage = (event) => {
          try {
            const message: Message = JSON.parse(event.data);
            this.handleMessage(message);
          } catch (e) {
            console.error('Failed to parse message:', e);
          }
        };

        this.ws.onerror = (error) => {
          console.error('WebSocket error:', error);
          reject(error);
        };

        this.ws.onclose = () => {
          console.log('WebSocket disconnected');
          this.attemptReconnect(token, roomId);
        };
      } catch (error) {
        reject(error);
      }
    });
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  emit(event: string, payload: any) {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      console.warn(`WebSocket not ready. Message "${event}" not sent.`);
      return;
    }

    const message: Message = { event, payload };
    this.ws.send(JSON.stringify(message));
  }

  on(event: string, handler: (payload: any) => void) {
    if (!this.handlers.has(event)) {
      this.handlers.set(event, []);
    }
    this.handlers.get(event)!.push(handler);
  }

  off(event: string, handler: (payload: any) => void) {
    const handlers = this.handlers.get(event);
    if (handlers) {
      const index = handlers.indexOf(handler);
      if (index > -1) {
        handlers.splice(index, 1);
      }
    }
  }

  private handleMessage(message: Message) {
    const handlers = this.handlers.get(message.event);
    if (handlers) {
      handlers.forEach((handler) => {
        try {
          handler(message.payload);
        } catch (e) {
          console.error(`Error in handler for event "${message.event}":`, e);
        }
      });
    }
  }

  private attemptReconnect(token: string, roomId?: string) {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      this.wasReconnecting = true;
      const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);
      console.log(`Attempting to reconnect in ${delay}ms... (attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
      this.onReconnecting?.(this.reconnectAttempts, this.maxReconnectAttempts);
      setTimeout(() => {
        this.connect(token, roomId).catch(() => {});
      }, delay);
    } else {
      console.error('Max reconnect attempts reached.');
      this.onReconnectFailed?.();
    }
  }

  sendAudio(data: ArrayBuffer) {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) return;
    this.ws.send(data);
  }

  isConnected(): boolean {
    return this.ws !== null && this.ws.readyState === WebSocket.OPEN;
  }
}
