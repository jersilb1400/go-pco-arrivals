import React, { useEffect, useState, useCallback, useRef } from 'react';
import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { motion, AnimatePresence, useAnimation } from 'framer-motion';
import '../components/Billboard.css';
import { 
  Users, 
  Clock, 
  RefreshCw, 
  Wifi,
  WifiOff,
  AlertCircle,
  Maximize2,
  Minimize2,
  Home,
  Shield,
  Bell,
  CheckCircle,
  Play,
  Pause,
  Settings,
  RotateCcw
} from 'lucide-react';
import apiService from '../services/api';
import { useWebSocket } from '../contexts/WebSocketContext';
import { useApp } from '../contexts/AppContext';
import type { CheckInDisplay, BillboardState, Notification } from '../types/api';

const BillboardPage: React.FC = () => {
  const { locationId } = useParams<{ locationId: string }>();
  const { 
    connect, 
    disconnect, 
    onCheckInUpdate, 
    onNotificationUpdate,
    onBillboardStateChange,
    connectionStatus,
    subscribeToLocation,
    subscribeToNotifications,
    subscribeToBillboardState
  } = useWebSocket();
  const { state: appState } = useApp();
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [lastUpdated, setLastUpdated] = useState<Date>(new Date());
  const [autoRefresh, setAutoRefresh] = useState(true);
  const [refreshInterval, setRefreshInterval] = useState(30); // seconds
  const [showSettings, setShowSettings] = useState(false);
  const [isPaused, setIsPaused] = useState(false);
  const [notificationCount, setNotificationCount] = useState(0);
  const autoRefreshRef = useRef<number | null>(null);
  const controls = useAnimation();

  // Get global billboard state (active event)
  const { data: billboardControl, refetch: refetchBillboardControl } = useQuery({
    queryKey: ['billboard-control'],
    queryFn: () => apiService.getBillboardControl(),
    refetchInterval: autoRefresh && !isPaused ? refreshInterval * 1000 : false,
  });

  // Get active pickup notifications (reduced polling since we have WebSocket)
  const { data: notificationsResponse, refetch: refetchNotifications } = useQuery({
    queryKey: ['notifications', 'active'],
    queryFn: () => apiService.getActiveNotifications(),
    refetchInterval: autoRefresh && !isPaused ? refreshInterval * 1000 : false,
    enabled: !!billboardControl?.control?.is_active,
  });

  // Get location-specific billboard state
  const { data: billboardState, refetch: refetchState } = useQuery({
    queryKey: ['billboard', 'state', locationId],
    queryFn: () => apiService.getBillboardState(locationId!),
    enabled: !!locationId,
    refetchInterval: autoRefresh && !isPaused ? refreshInterval * 1000 : false,
  });

  // Connect to WebSocket for real-time updates
  useEffect(() => {
    if (locationId) {
      connect(locationId);
      subscribeToLocation(locationId);
    }

    return () => {
      disconnect();
    };
  }, [locationId, connect, disconnect, subscribeToLocation]);

  // Subscribe to real-time updates
  useEffect(() => {
    subscribeToNotifications();
    subscribeToBillboardState();
  }, [subscribeToNotifications, subscribeToBillboardState]);

  // Handle real-time updates
  useEffect(() => {
    // Listen for check-in updates
    const unsubscribeCheckIns = onCheckInUpdate((update) => {
      if (update.location_id === locationId) {
        console.log('Real-time check-in update:', update);
        refetchState();
        setLastUpdated(new Date());
        // Trigger animation for new updates
        controls.start({ scale: [1, 1.05, 1], transition: { duration: 0.3 } });
      }
    });

    // Listen for notification updates
    const unsubscribeNotifications = onNotificationUpdate((notification) => {
      console.log('Real-time notification update:', notification);
      refetchNotifications();
      setLastUpdated(new Date());
      // Trigger animation for new notifications
      controls.start({ scale: [1, 1.05, 1], transition: { duration: 0.3 } });
    });

    // Listen for billboard state changes
    const unsubscribeBillboardState = onBillboardStateChange((state) => {
      console.log('Real-time billboard state change:', state);
      refetchBillboardControl();
      setLastUpdated(new Date());
    });

    return () => {
      unsubscribeCheckIns();
      unsubscribeNotifications();
      unsubscribeBillboardState();
    };
  }, [
    locationId,
    onCheckInUpdate,
    onNotificationUpdate,
    onBillboardStateChange,
    refetchState,
    refetchNotifications,
    refetchBillboardControl,
    controls
  ]);

  // Auto-refresh logic
  useEffect(() => {
    if (autoRefresh && !isPaused && refreshInterval > 0) {
      autoRefreshRef.current = window.setInterval(() => {
        handleRefresh();
      }, refreshInterval * 1000);
    }

    return () => {
      if (autoRefreshRef.current) {
        window.clearInterval(autoRefreshRef.current);
      }
    };
  }, [autoRefresh, isPaused, refreshInterval]);

  // Update notification count
  useEffect(() => {
    const notifications = notificationsResponse?.notifications || [];
    setNotificationCount(notifications.length);
  }, [notificationsResponse]);

  const handleRefresh = useCallback(async () => {
    try {
      await Promise.all([
        refetchNotifications(),
        refetchState(),
        refetchBillboardControl(),
      ]);
      setLastUpdated(new Date());
      // Trigger refresh animation
      controls.start({ rotate: [0, 360], transition: { duration: 0.5 } });
    } catch (error) {
      console.error('Failed to refresh:', error);
    }
  }, [refetchNotifications, refetchState, refetchBillboardControl, controls]);

  const toggleFullscreen = () => {
    if (!document.fullscreenElement) {
      document.documentElement.requestFullscreen();
      setIsFullscreen(true);
    } else {
      document.exitFullscreen();
      setIsFullscreen(false);
    }
  };

  const togglePause = () => {
    setIsPaused(!isPaused);
  };

  const formatTimeAgo = (timestamp: string) => {
    const now = new Date();
    const time = new Date(timestamp);
    const diffMs = now.getTime() - time.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    
    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins}m ago`;
    const diffHours = Math.floor(diffMins / 60);
    return `${diffHours}h ${diffMins % 60}m ago`;
  };

  // Group notifications by security code
  const groupedNotifications = (notificationsResponse?.notifications || []).reduce((groups, notification) => {
    const code = notification.security_code || 'No Code';
    if (!groups[code]) {
      groups[code] = [];
    }
    groups[code].push(notification);
    return groups;
  }, {} as Record<string, Notification[]>);

  // Check if billboard is active
  const isBillboardActive = billboardControl?.control?.is_active;
  const activeEvent = billboardControl?.control;

  if (!locationId) {
    return (
      <div className="flex items-center justify-center h-64">
        <AlertCircle className="h-8 w-8 text-gray-400" />
        <span className="ml-2 text-gray-500">Location ID is required</span>
      </div>
    );
  }

  return (
    <div className={`min-h-screen ${isFullscreen ? 'bg-gradient-to-br from-gray-900 via-black to-gray-900 text-white' : 'bg-gradient-to-br from-gray-50 to-gray-100'}`}>
      {/* Enhanced Header */}
      <motion.div 
        className={`${isFullscreen ? 'bg-black/80 backdrop-blur-sm border-gray-700' : 'bg-white/95 backdrop-blur-sm border-gray-200'} border-b px-6 py-4 shadow-lg`}
        initial={{ y: -50, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ duration: 0.5 }}
      >
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <motion.div
              animate={controls}
              className={`p-3 rounded-full ${isFullscreen ? 'bg-green-600/20' : 'bg-green-100'}`}
            >
              <Bell className={`h-6 w-6 ${isFullscreen ? 'text-green-400' : 'text-green-600'}`} />
            </motion.div>
            <div>
              <h1 className={`${isFullscreen ? 'text-4xl' : 'text-2xl'} font-bold ${isFullscreen ? 'text-white' : 'text-gray-900'}`}>
                {isBillboardActive ? 'Child Pickup Requests' : 'Billboard Display'}
              </h1>
              {activeEvent && (
                <p className={`mt-1 ${isFullscreen ? 'text-xl' : 'text-lg'} ${isFullscreen ? 'text-gray-300' : 'text-gray-600'}`}>
                  {activeEvent.event_name} • {activeEvent.location_name}
                </p>
              )}
              <p className={`mt-1 text-sm ${isFullscreen ? 'text-gray-400' : 'text-gray-500'}`}>
                {notificationCount} children ready for pickup
              </p>
            </div>
          </div>
          
          <div className="flex items-center space-x-3">
            {/* Connection Status */}
            <motion.div 
              className="flex items-center space-x-2"
              whileHover={{ scale: 1.05 }}
            >
              {connectionStatus.isConnected ? (
                <Wifi className={`h-4 w-4 ${isFullscreen ? 'text-green-400' : 'text-green-600'}`} />
              ) : (
                <WifiOff className={`h-4 w-4 ${isFullscreen ? 'text-red-400' : 'text-red-600'}`} />
              )}
              <span className={`text-sm font-medium ${isFullscreen ? 'text-gray-300' : 'text-gray-500'}`}>
                {connectionStatus.isConnected ? 'Live' : 'Offline'}
              </span>
            </motion.div>

            {/* Last Updated */}
            <div className={`text-sm ${isFullscreen ? 'text-gray-400' : 'text-gray-500'}`}>
              Last updated: {lastUpdated.toLocaleTimeString()}
            </div>

            {/* Auto-refresh Status */}
            {autoRefresh && (
              <motion.div 
                className="flex items-center space-x-1"
                animate={{ opacity: isPaused ? 0.5 : 1 }}
              >
                <RotateCcw className={`h-4 w-4 ${isFullscreen ? 'text-blue-400' : 'text-blue-600'}`} />
                <span className={`text-sm ${isFullscreen ? 'text-gray-300' : 'text-gray-500'}`}>
                  {refreshInterval}s
                </span>
              </motion.div>
            )}

            {/* Pause/Play Button */}
            <motion.button
              onClick={togglePause}
              className={`p-2 rounded-md ${isFullscreen ? 'bg-gray-800 text-white hover:bg-gray-700' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'}`}
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              title={isPaused ? 'Resume Auto-refresh' : 'Pause Auto-refresh'}
            >
              {isPaused ? <Play className="h-4 w-4" /> : <Pause className="h-4 w-4" />}
            </motion.button>

            {/* Refresh Button */}
            <motion.button
              onClick={handleRefresh}
              className={`p-2 rounded-md ${isFullscreen ? 'bg-gray-800 text-white hover:bg-gray-700' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'}`}
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              title="Manual Refresh"
            >
              <RefreshCw className="h-4 w-4" />
            </motion.button>

            {/* Settings Button */}
            <motion.button
              onClick={() => setShowSettings(!showSettings)}
              className={`p-2 rounded-md ${isFullscreen ? 'bg-gray-800 text-white hover:bg-gray-700' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'}`}
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              title="Settings"
            >
              <Settings className="h-4 w-4" />
            </motion.button>

            {/* Fullscreen Toggle */}
            <motion.button
              onClick={toggleFullscreen}
              className={`p-2 rounded-md ${isFullscreen ? 'bg-gray-800 text-white hover:bg-gray-700' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'}`}
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              title={isFullscreen ? 'Exit Fullscreen' : 'Enter Fullscreen'}
            >
              {isFullscreen ? <Minimize2 className="h-4 w-4" /> : <Maximize2 className="h-4 w-4" />}
            </motion.button>

            {/* Back to Admin */}
            <motion.a
              href="/admin"
              className={`p-2 rounded-md ${isFullscreen ? 'bg-gray-800 text-white hover:bg-gray-700' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'}`}
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              title="Back to Admin"
            >
              <Shield className="h-4 w-4" />
            </motion.a>
          </div>
        </div>

        {/* Settings Panel */}
        <AnimatePresence>
          {showSettings && (
            <motion.div
              initial={{ height: 0, opacity: 0 }}
              animate={{ height: 'auto', opacity: 1 }}
              exit={{ height: 0, opacity: 0 }}
              transition={{ duration: 0.3 }}
              className={`mt-4 p-4 rounded-lg ${isFullscreen ? 'bg-gray-800/50' : 'bg-gray-50'}`}
            >
              <div className="flex items-center space-x-6">
                <div className="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    id="autoRefresh"
                    checked={autoRefresh}
                    onChange={(e) => setAutoRefresh(e.target.checked)}
                    className="rounded"
                  />
                  <label htmlFor="autoRefresh" className={`text-sm ${isFullscreen ? 'text-gray-300' : 'text-gray-600'}`}>
                    Auto-refresh
                  </label>
                </div>
                <div className="flex items-center space-x-2">
                  <label className={`text-sm ${isFullscreen ? 'text-gray-300' : 'text-gray-600'}`}>
                    Interval:
                  </label>
                  <select
                    value={refreshInterval}
                    onChange={(e) => setRefreshInterval(Number(e.target.value))}
                    className={`px-2 py-1 rounded text-sm ${isFullscreen ? 'bg-gray-700 text-white' : 'bg-white text-gray-900'}`}
                  >
                    <option value={15}>15s</option>
                    <option value={30}>30s</option>
                    <option value={60}>1m</option>
                    <option value={120}>2m</option>
                  </select>
                </div>
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </motion.div>

      {/* Enhanced Main Content */}
      <div className={`p-6 ${isFullscreen ? 'bg-transparent' : ''}`}>
        {!isBillboardActive ? (
          // No active billboard
          <motion.div 
            className="flex items-center justify-center h-64"
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.5 }}
          >
            <div className="text-center">
              <motion.div
                animate={{ rotate: [0, 10, -10, 0] }}
                transition={{ duration: 2, repeat: Infinity, repeatDelay: 3 }}
              >
                <AlertCircle className={`mx-auto h-16 w-16 ${isFullscreen ? 'text-gray-400' : 'text-gray-400'}`} />
              </motion.div>
              <h3 className={`mt-4 text-xl font-medium ${isFullscreen ? 'text-white' : 'text-gray-900'}`}>
                No Active Billboard
              </h3>
              <p className={`mt-2 text-sm ${isFullscreen ? 'text-gray-400' : 'text-gray-500'}`}>
                Please launch a billboard from the admin panel to display pickup notifications.
              </p>
            </div>
          </motion.div>
        ) : Object.keys(groupedNotifications).length === 0 ? (
          // No notifications
          <motion.div 
            className="flex items-center justify-center h-64"
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.5 }}
          >
            <div className="text-center">
              <motion.div
                animate={{ y: [0, -10, 0] }}
                transition={{ duration: 2, repeat: Infinity }}
              >
                <Users className={`mx-auto h-16 w-16 ${isFullscreen ? 'text-gray-400' : 'text-gray-400'}`} />
              </motion.div>
              <h3 className={`mt-4 text-xl font-medium ${isFullscreen ? 'text-white' : 'text-gray-900'}`}>
                No Pickup Requests
              </h3>
              <p className={`mt-2 text-sm ${isFullscreen ? 'text-gray-400' : 'text-gray-500'}`}>
                Children will appear here when they are ready for pickup.
              </p>
            </div>
          </motion.div>
        ) : (
          // Enhanced notifications display
          <div className="space-y-8">
            <AnimatePresence>
              {Object.entries(groupedNotifications).map(([securityCode, notifications], groupIndex) => (
                <motion.div
                  key={securityCode}
                  initial={{ opacity: 0, y: 50, scale: 0.9 }}
                  animate={{ opacity: 1, y: 0, scale: 1 }}
                  exit={{ opacity: 0, y: -50, scale: 0.9 }}
                  transition={{ 
                    duration: 0.5, 
                    delay: groupIndex * 0.1,
                    type: "spring",
                    stiffness: 100
                  }}
                  className={`${isFullscreen ? 'bg-gradient-to-r from-gray-900/80 to-gray-800/80 backdrop-blur-sm border-gray-600' : 'bg-white border-gray-200'} border rounded-xl p-6 shadow-xl`}
                >
                  {/* Enhanced Security Code Header */}
                  <div className="mb-6">
                    <motion.div
                      className="flex items-center space-x-3 mb-3"
                      initial={{ x: -20, opacity: 0 }}
                      animate={{ x: 0, opacity: 1 }}
                      transition={{ delay: 0.2 }}
                    >
                      <div className={`p-3 rounded-full ${isFullscreen ? 'bg-blue-600/20' : 'bg-blue-100'}`}>
                        <Shield className={`h-6 w-6 ${isFullscreen ? 'text-blue-400' : 'text-blue-600'}`} />
                      </div>
                      <h2 className={`${isFullscreen ? 'text-4xl' : 'text-3xl'} font-bold ${isFullscreen ? 'text-white' : 'text-gray-900'}`}>
                        Security Code: {securityCode}
                      </h2>
                    </motion.div>
                    <motion.p 
                      className={`${isFullscreen ? 'text-xl' : 'text-lg'} ${isFullscreen ? 'text-gray-300' : 'text-gray-600'}`}
                      initial={{ opacity: 0 }}
                      animate={{ opacity: 1 }}
                      transition={{ delay: 0.3 }}
                    >
                      {notifications.length} child{notifications.length !== 1 ? 'ren' : ''} ready for pickup
                    </motion.p>
                  </div>

                  {/* Enhanced Children Grid */}
                  <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                    {notifications.map((notification, index) => (
                      <motion.div
                        key={notification.id}
                        initial={{ opacity: 0, scale: 0.8, y: 20 }}
                        animate={{ opacity: 1, scale: 1, y: 0 }}
                        exit={{ opacity: 0, scale: 0.8, y: -20 }}
                        transition={{ 
                          duration: 0.4, 
                          delay: index * 0.1,
                          type: "spring",
                          stiffness: 200
                        }}
                        whileHover={{ 
                          scale: 1.02,
                          y: -5,
                          transition: { duration: 0.2 }
                        }}
                        className={`${isFullscreen ? 'bg-gradient-to-br from-gray-800/80 to-gray-700/80 backdrop-blur-sm border-gray-600' : 'bg-gradient-to-br from-gray-50 to-white border-gray-200'} border rounded-xl p-5 shadow-lg hover:shadow-xl transition-all duration-300`}
                      >
                        <div className="flex items-start justify-between">
                          <div className="flex-1">
                            <motion.h3 
                              className={`${isFullscreen ? 'text-2xl' : 'text-xl'} font-bold ${isFullscreen ? 'text-white' : 'text-gray-900'} mb-3`}
                              initial={{ opacity: 0 }}
                              animate={{ opacity: 1 }}
                              transition={{ delay: 0.1 }}
                            >
                              {notification.child_name || 'Unknown Child'}
                            </motion.h3>
                            <motion.div
                              className="space-y-2"
                              initial={{ opacity: 0 }}
                              animate={{ opacity: 1 }}
                              transition={{ delay: 0.2 }}
                            >
                              <p className={`${isFullscreen ? 'text-lg' : 'text-base'} ${isFullscreen ? 'text-gray-300' : 'text-gray-600'}`}>
                                📍 {notification.location_name || 'Unknown Location'}
                              </p>
                              {notification.parent_name && (
                                <p className={`${isFullscreen ? 'text-lg' : 'text-base'} ${isFullscreen ? 'text-gray-300' : 'text-gray-600'}`}>
                                  👤 {notification.parent_name}
                                </p>
                              )}
                              <p className={`text-sm ${isFullscreen ? 'text-gray-400' : 'text-gray-500'} flex items-center`}>
                                <Clock className="h-3 w-3 mr-1" />
                                {formatTimeAgo(notification.created_at)}
                              </p>
                            </motion.div>
                          </div>
                          <motion.div 
                            className={`ml-4 p-3 rounded-full ${isFullscreen ? 'bg-green-600/20' : 'bg-green-100'}`}
                            animate={{ scale: [1, 1.1, 1] }}
                            transition={{ duration: 2, repeat: Infinity, repeatDelay: 3 }}
                          >
                            <CheckCircle className={`h-6 w-6 ${isFullscreen ? 'text-green-400' : 'text-green-600'}`} />
                          </motion.div>
                        </div>
                      </motion.div>
                    ))}
                  </div>
                </motion.div>
              ))}
            </AnimatePresence>
          </div>
        )}
      </div>
    </div>
  );
};

export default BillboardPage; 