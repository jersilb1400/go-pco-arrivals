import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { 
  Wifi, 
  WifiOff, 
  Bell, 
  X,
  CheckCircle,
  AlertCircle,
  Info
} from 'lucide-react';
import { useWebSocket } from '../contexts/WebSocketContext';
import type { Notification } from '../types/api';

interface RealTimeNotificationProps {
  className?: string;
}

interface ToastNotification {
  id: string;
  type: 'success' | 'error' | 'info' | 'warning';
  message: string;
  timestamp: Date;
}

const RealTimeNotification: React.FC<RealTimeNotificationProps> = ({ className = '' }) => {
  const { connectionStatus, onNotificationUpdate, onBillboardStateChange } = useWebSocket();
  const [toasts, setToasts] = useState<ToastNotification[]>([]);

  // Handle real-time notifications
  useEffect(() => {
    const unsubscribeNotifications = onNotificationUpdate((notification) => {
      if (notification.child_name) {
        addToast('info', `New pickup request: ${notification.child_name} (${notification.security_code})`);
      }
    });

    const unsubscribeBillboardState = onBillboardStateChange((state) => {
      if (state.is_active) {
        addToast('success', `Billboard launched for ${state.event_name}`);
      } else {
        addToast('info', 'Billboard cleared');
      }
    });

    return () => {
      unsubscribeNotifications();
      unsubscribeBillboardState();
    };
  }, [onNotificationUpdate, onBillboardStateChange]);

  const addToast = (type: ToastNotification['type'], message: string) => {
    const id = Date.now().toString();
    const newToast: ToastNotification = {
      id,
      type,
      message,
      timestamp: new Date(),
    };

    setToasts(prev => [...prev, newToast]);

    // Auto-remove toast after 5 seconds
    setTimeout(() => {
      removeToast(id);
    }, 5000);
  };

  const removeToast = (id: string) => {
    setToasts(prev => prev.filter(toast => toast.id !== id));
  };

  const getToastIcon = (type: ToastNotification['type']) => {
    switch (type) {
      case 'success':
        return <CheckCircle className="h-5 w-5 text-green-500" />;
      case 'error':
        return <AlertCircle className="h-5 w-5 text-red-500" />;
      case 'warning':
        return <AlertCircle className="h-5 w-5 text-yellow-500" />;
      case 'info':
        return <Info className="h-5 w-5 text-blue-500" />;
    }
  };

  const getToastStyles = (type: ToastNotification['type']) => {
    switch (type) {
      case 'success':
        return 'bg-green-50 border-green-200 text-green-800';
      case 'error':
        return 'bg-red-50 border-red-200 text-red-800';
      case 'warning':
        return 'bg-yellow-50 border-yellow-200 text-yellow-800';
      case 'info':
        return 'bg-blue-50 border-blue-200 text-blue-800';
    }
  };

  return (
    <div className={`fixed top-4 right-4 z-50 space-y-2 ${className}`}>
      {/* Connection Status */}
      <motion.div
        initial={{ opacity: 0, x: 100 }}
        animate={{ opacity: 1, x: 0 }}
        className={`flex items-center space-x-2 px-3 py-2 rounded-md shadow-md ${
          connectionStatus.isConnected 
            ? 'bg-green-50 border border-green-200 text-green-800' 
            : 'bg-red-50 border border-red-200 text-red-800'
        }`}
      >
        {connectionStatus.isConnected ? (
          <Wifi className="h-4 w-4" />
        ) : (
          <WifiOff className="h-4 w-4" />
        )}
        <span className="text-sm font-medium">
          {connectionStatus.isConnected ? 'Live Updates' : 'Offline'}
        </span>
        {connectionStatus.error && (
          <span className="text-xs opacity-75">({connectionStatus.error})</span>
        )}
      </motion.div>

      {/* Toast Notifications */}
      <AnimatePresence>
        {toasts.map((toast) => (
          <motion.div
            key={toast.id}
            initial={{ opacity: 0, x: 100, scale: 0.9 }}
            animate={{ opacity: 1, x: 0, scale: 1 }}
            exit={{ opacity: 0, x: 100, scale: 0.9 }}
            transition={{ duration: 0.2 }}
            className={`flex items-start space-x-3 p-4 rounded-md shadow-lg border max-w-sm ${getToastStyles(toast.type)}`}
          >
            <div className="flex-shrink-0 mt-0.5">
              {getToastIcon(toast.type)}
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium">{toast.message}</p>
              <p className="text-xs opacity-75 mt-1">
                {toast.timestamp.toLocaleTimeString()}
              </p>
            </div>
            <button
              onClick={() => removeToast(toast.id)}
              className="flex-shrink-0 text-gray-400 hover:text-gray-600 transition-colors"
            >
              <X className="h-4 w-4" />
            </button>
          </motion.div>
        ))}
      </AnimatePresence>
    </div>
  );
};

export default RealTimeNotification; 