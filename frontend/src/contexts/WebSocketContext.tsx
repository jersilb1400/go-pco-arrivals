import React, { createContext, useContext, useState, useEffect } from 'react';
import type { ReactNode } from 'react';
import { websocketService } from '../services/websocket';
import type { WebSocketConnectionStatus, RealTimeUpdate, Notification, BillboardControl } from '../types/api';

interface WebSocketContextType {
  connectionStatus: WebSocketConnectionStatus;
  connect: (locationId?: string) => Promise<void>;
  disconnect: () => void;
  onCheckInUpdate: (callback: (checkIn: RealTimeUpdate) => void) => () => void;
  onStateUpdate: (callback: (state: RealTimeUpdate) => void) => () => void;
  onNotificationUpdate: (callback: (notification: Notification) => void) => () => void;
  onBillboardStateChange: (callback: (state: BillboardControl) => void) => () => void;
  onSecurityCodeAdded: (callback: (data: { code: string; event_id: string }) => void) => () => void;
  onSecurityCodeRemoved: (callback: (data: { code: string; event_id: string }) => void) => () => void;
  onBillboardLaunched: (callback: (state: BillboardControl) => void) => () => void;
  onBillboardCleared: (callback: (data: { event_id: string }) => void) => () => void;
  subscribeToLocation: (locationId: string) => void;
  subscribeToNotifications: () => void;
  subscribeToBillboardState: () => void;
}

const WebSocketContext = createContext<WebSocketContextType | undefined>(undefined);

export const useWebSocket = () => {
  const context = useContext(WebSocketContext);
  if (context === undefined) {
    throw new Error('useWebSocket must be used within a WebSocketProvider');
  }
  return context;
};

interface WebSocketProviderProps {
  children: ReactNode;
}

export const WebSocketProvider: React.FC<WebSocketProviderProps> = ({ children }) => {
  const [connectionStatus, setConnectionStatus] = useState<WebSocketConnectionStatus>({
    isConnected: false,
    isConnecting: false,
  });

  useEffect(() => {
    const unsubscribe = websocketService.onConnectionStatusChange((status) => {
      setConnectionStatus(status);
    });

    return unsubscribe;
  }, []);

  const connect = async (locationId?: string) => {
    try {
      await websocketService.connect(locationId);
    } catch (error) {
      console.error('WebSocket connection failed:', error);
      // Don't throw the error to prevent uncaught promise rejections
    }
  };

  const disconnect = () => {
    websocketService.disconnect();
  };

  const onCheckInUpdate = (callback: (checkIn: RealTimeUpdate) => void) => {
    return websocketService.onCheckInUpdate(callback);
  };

  const onStateUpdate = (callback: (state: RealTimeUpdate) => void) => {
    return websocketService.onStateUpdate(callback);
  };

  const onNotificationUpdate = (callback: (notification: Notification) => void) => {
    return websocketService.onNotificationUpdate(callback);
  };

  const onBillboardStateChange = (callback: (state: BillboardControl) => void) => {
    return websocketService.onBillboardStateChange(callback);
  };

  const onSecurityCodeAdded = (callback: (data: { code: string; event_id: string }) => void) => {
    return websocketService.onSecurityCodeAdded(callback);
  };

  const onSecurityCodeRemoved = (callback: (data: { code: string; event_id: string }) => void) => {
    return websocketService.onSecurityCodeRemoved(callback);
  };

  const onBillboardLaunched = (callback: (state: BillboardControl) => void) => {
    return websocketService.onBillboardLaunched(callback);
  };

  const onBillboardCleared = (callback: (data: { event_id: string }) => void) => {
    return websocketService.onBillboardCleared(callback);
  };

  const subscribeToLocation = (locationId: string) => {
    websocketService.subscribeToLocation(locationId);
  };

  const subscribeToNotifications = () => {
    websocketService.subscribeToNotifications();
  };

  const subscribeToBillboardState = () => {
    websocketService.subscribeToBillboardState();
  };

  const value: WebSocketContextType = {
    connectionStatus,
    connect,
    disconnect,
    onCheckInUpdate,
    onStateUpdate,
    onNotificationUpdate,
    onBillboardStateChange,
    onSecurityCodeAdded,
    onSecurityCodeRemoved,
    onBillboardLaunched,
    onBillboardCleared,
    subscribeToLocation,
    subscribeToNotifications,
    subscribeToBillboardState,
  };

  return <WebSocketContext.Provider value={value}>{children}</WebSocketContext.Provider>;
}; 