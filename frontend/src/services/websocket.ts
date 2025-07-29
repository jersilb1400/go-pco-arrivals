import type { RealTimeUpdate, WebSocketMessage, WebSocketConnectionStatus, Notification, BillboardControl } from '../types/api';

interface WebSocketEventMap {
  'new_check_in': RealTimeUpdate;
  'state_update': RealTimeUpdate;
  'notification_update': Notification;
  'billboard_state_change': BillboardControl;
  'security_code_added': { code: string; event_id: string };
  'security_code_removed': { code: string; event_id: string };
  'billboard_launched': BillboardControl;
  'billboard_cleared': { event_id: string };
  'connection_status': WebSocketConnectionStatus;
}

class WebSocketService {
  private ws: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000;
  private listeners: Map<string, Set<(data: any) => void>> = new Map();
  private connectionStatus: WebSocketConnectionStatus = {
    isConnected: false,
    isConnecting: false,
  };
  private connectionStatusListeners: Set<(status: WebSocketConnectionStatus) => void> = new Set();
  private heartbeatInterval: number | null = null;
  private reconnectTimeout: number | null = null;
  private isConnecting = false;
  private shouldReconnect = true;

  private baseURL: string;

  constructor(baseURL: string = 'ws://localhost:3000') {
    this.baseURL = baseURL;
  }

  connect(locationId?: string): Promise<void> {
    return new Promise((resolve, reject) => {
      // Prevent multiple simultaneous connection attempts
      if (this.isConnecting) {
        console.log('WebSocket connection already in progress, skipping...');
        resolve();
        return;
      }

      if (this.ws?.readyState === WebSocket.OPEN) {
        console.log('WebSocket already connected');
        resolve();
        return;
      }

      this.isConnecting = true;
      this.connectionStatus.isConnecting = true;
      this.connectionStatus.error = undefined;
      this.notifyConnectionStatusChange();

      const wsUrl = locationId 
        ? `${this.baseURL}/ws/billboard/${locationId}`
        : `${this.baseURL}/ws`;

      console.log(`Connecting to WebSocket: ${wsUrl}`);
      
      try {
        this.ws = new WebSocket(wsUrl);

        this.ws.onopen = () => {
          console.log('WebSocket connected successfully');
          this.isConnecting = false;
          this.connectionStatus.isConnected = true;
          this.connectionStatus.isConnecting = false;
          this.connectionStatus.error = undefined;
          this.reconnectAttempts = 0;
          this.notifyConnectionStatusChange();
          this.startHeartbeat();
          resolve();
        };

        this.ws.onmessage = (event) => {
          try {
            const message: WebSocketMessage = JSON.parse(event.data);
            console.log('WebSocket message received:', message.type, message.data);
            this.handleMessage(message);
          } catch (error) {
            console.error('Failed to parse WebSocket message:', error);
          }
        };

        this.ws.onclose = (event) => {
          console.log('WebSocket disconnected:', event.code, event.reason);
          this.isConnecting = false;
          this.connectionStatus.isConnected = false;
          this.connectionStatus.isConnecting = false;
          this.notifyConnectionStatusChange();
          this.stopHeartbeat();

          // Only attempt reconnection if we should and it wasn't a clean close
          if (this.shouldReconnect && !event.wasClean && this.reconnectAttempts < this.maxReconnectAttempts) {
            this.scheduleReconnect();
          }
        };

        this.ws.onerror = (error) => {
          console.error('WebSocket error:', error);
          this.isConnecting = false;
          this.connectionStatus.error = 'Connection failed';
          this.connectionStatus.isConnecting = false;
          this.notifyConnectionStatusChange();
          reject(error);
        };
      } catch (error) {
        this.isConnecting = false;
        this.connectionStatus.error = 'Failed to create WebSocket connection';
        this.connectionStatus.isConnecting = false;
        this.notifyConnectionStatusChange();
        reject(error);
      }
    });
  }

  private scheduleReconnect(): void {
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout);
    }

    this.reconnectAttempts++;
    const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);
    
    console.log(`Scheduling WebSocket reconnection attempt ${this.reconnectAttempts} in ${delay}ms`);
    
    this.reconnectTimeout = window.setTimeout(() => {
      if (this.shouldReconnect) {
        this.connect().catch(error => {
          console.error('WebSocket reconnection failed:', error);
        });
      }
    }, delay);
  }

  private startHeartbeat(): void {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval);
    }

    this.heartbeatInterval = window.setInterval(() => {
      if (this.ws?.readyState === WebSocket.OPEN) {
        this.ws.send(JSON.stringify({
          type: 'heartbeat',
          data: { timestamp: Date.now() },
          timestamp: new Date().toISOString()
        }));
      }
    }, 30000); // Send heartbeat every 30 seconds
  }

  private stopHeartbeat(): void {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval);
      this.heartbeatInterval = null;
    }
  }

  disconnect(): void {
    console.log('Disconnecting WebSocket...');
    this.shouldReconnect = false;
    this.isConnecting = false;
    
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout);
      this.reconnectTimeout = null;
    }
    
    this.stopHeartbeat();
    
    if (this.ws) {
      this.ws.close(1000, 'Client disconnecting');
      this.ws = null;
    }
    
    this.connectionStatus.isConnected = false;
    this.connectionStatus.isConnecting = false;
    this.notifyConnectionStatusChange();
  }

  private handleMessage(message: WebSocketMessage): void {
    const listeners = this.listeners.get(message.type);
    if (listeners) {
      listeners.forEach((listener) => {
        try {
          listener(message.data);
        } catch (error) {
          console.error('Error in WebSocket message listener:', error);
        }
      });
    }
  }

  private notifyConnectionStatusChange(): void {
    this.connectionStatusListeners.forEach((listener) => {
      try {
        listener({ ...this.connectionStatus });
      } catch (error) {
        console.error('Error in connection status listener:', error);
      }
    });
  }

  on<T extends keyof WebSocketEventMap>(eventType: T, callback: (data: WebSocketEventMap[T]) => void): () => void {
    if (!this.listeners.has(eventType)) {
      this.listeners.set(eventType, new Set());
    }
    
    this.listeners.get(eventType)!.add(callback);

    // Return unsubscribe function
    return () => {
      const listeners = this.listeners.get(eventType);
      if (listeners) {
        listeners.delete(callback);
        if (listeners.size === 0) {
          this.listeners.delete(eventType);
        }
      }
    };
  }

  send(message: WebSocketMessage): void {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message));
    } else {
      console.warn('WebSocket is not connected, cannot send message:', message.type);
    }
  }

  // Send specific message types
  subscribeToLocation(locationId: string): void {
    this.send({
      type: 'subscribe_location',
      data: { location_id: locationId },
      timestamp: new Date().toISOString(),
    });
  }

  subscribeToNotifications(): void {
    this.send({
      type: 'subscribe_notifications',
      data: {},
      timestamp: new Date().toISOString(),
    });
  }

  subscribeToBillboardState(): void {
    this.send({
      type: 'subscribe_billboard_state',
      data: {},
      timestamp: new Date().toISOString(),
    });
  }

  getConnectionStatus(): WebSocketConnectionStatus {
    return { ...this.connectionStatus };
  }

  // Specific event handlers for billboard updates
  onCheckInUpdate(callback: (checkIn: RealTimeUpdate) => void): () => void {
    return this.on('new_check_in', callback);
  }

  onStateUpdate(callback: (state: RealTimeUpdate) => void): () => void {
    return this.on('state_update', callback);
  }

  onNotificationUpdate(callback: (notification: Notification) => void): () => void {
    return this.on('notification_update', callback);
  }

  onBillboardStateChange(callback: (state: BillboardControl) => void): () => void {
    return this.on('billboard_state_change', callback);
  }

  onSecurityCodeAdded(callback: (data: { code: string; event_id: string }) => void): () => void {
    return this.on('security_code_added', callback);
  }

  onSecurityCodeRemoved(callback: (data: { code: string; event_id: string }) => void): () => void {
    return this.on('security_code_removed', callback);
  }

  onBillboardLaunched(callback: (state: BillboardControl) => void): () => void {
    return this.on('billboard_launched', callback);
  }

  onBillboardCleared(callback: (data: { event_id: string }) => void): () => void {
    return this.on('billboard_cleared', callback);
  }

  onConnectionStatusChange(callback: (status: WebSocketConnectionStatus) => void): () => void {
    this.connectionStatusListeners.add(callback);
    
    // Return unsubscribe function
    return () => {
      this.connectionStatusListeners.delete(callback);
    };
  }
}

export const websocketService = new WebSocketService();
export default websocketService; 