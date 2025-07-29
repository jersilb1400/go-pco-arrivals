import React, { useState, useEffect, useRef } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useAuth } from '../contexts/AuthContext';
import { useApp } from '../contexts/AppContext';
import { useWebSocket } from '../contexts/WebSocketContext';
import api from '../services/api';

interface Event {
  id: string;
  name: string;
  date: string;
  location: string;
  time?: string;
}

interface Notification {
  id: string;
  message: string;
  type: string;
  created_at: string;
}

const AdminPage: React.FC = () => {
  const { user } = useAuth();
  const { state, setSelectedEvent: setGlobalSelectedEvent, setSelectedDate: setGlobalSelectedDate } = useApp();
  const { 
    connectionStatus, 
    onNotificationUpdate, 
    onBillboardStateChange, 
    onSecurityCodeAdded, 
    onSecurityCodeRemoved,
    onBillboardLaunched,
    onBillboardCleared,
    subscribeToNotifications,
    subscribeToBillboardState
  } = useWebSocket();
  const queryClient = useQueryClient();
  const [selectedDate, setSelectedDate] = useState(state.selectedDate);
  const [selectedEvent, setSelectedEvent] = useState(state.selectedEvent?.id || '');
  const [newSecurityCode, setNewSecurityCode] = useState('');
  const [snackbar, setSnackbar] = useState<{ open: boolean; message: string; type: 'success' | 'error' | 'info' }>({ open: false, message: '', type: 'info' });
  const pollingRef = useRef<number | null>(null);

  // Fetch events by date
  const { data: eventsData, isLoading: eventsLoading } = useQuery({
    queryKey: ['events', selectedDate],
    queryFn: () => api.getEventsByDate(selectedDate),
    enabled: !!selectedDate,
    refetchOnWindowFocus: false,
    staleTime: 30000,
  });

  // Fetch active notifications (reduced polling since we have WebSocket)
  const { data: notificationsData, isLoading: notificationsLoading } = useQuery({
    queryKey: ['notifications', 'active'],
    queryFn: () => api.getActiveNotifications(),
    refetchInterval: 60000, // Poll every 60 seconds as backup
    refetchOnWindowFocus: false,
    staleTime: 30000,
  });

  // Fetch security codes
  const { data: securityCodesData, isLoading: securityCodesLoading } = useQuery({
    queryKey: ['security-codes'],
    queryFn: () => api.getSecurityCodes(),
    refetchOnWindowFocus: false,
    staleTime: 60000,
  });

  // Fetch billboard control state (reduced polling since we have WebSocket)
  const { data: billboardControlData, isLoading: billboardControlLoading } = useQuery({
    queryKey: ['billboard-control'],
    queryFn: () => api.getBillboardControl(),
    refetchInterval: 30000, // Poll every 30 seconds as backup
    refetchOnWindowFocus: false,
    staleTime: 15000,
  });

  // Mutations
  const addSecurityCodeMutation = useMutation({
    mutationFn: (code: string) => api.addSecurityCode(code),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['security-codes'] });
      setSnackbar({ open: true, message: 'Security code added successfully', type: 'success' });
    },
    onError: (error) => {
      setSnackbar({ open: true, message: 'Failed to add security code', type: 'error' });
    },
  });

  const removeSecurityCodeMutation = useMutation({
    mutationFn: (code: string) => api.removeSecurityCode(code),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['security-codes'] });
      setSnackbar({ open: true, message: 'Security code removed successfully', type: 'info' });
    },
    onError: (error) => {
      setSnackbar({ open: true, message: 'Failed to remove security code', type: 'error' });
    },
  });

  const launchBillboardMutation = useMutation({
    mutationFn: ({ eventId, locationId, securityCodes }: { eventId: string; locationId: string; securityCodes: string[] }) =>
      api.launchBillboard(eventId, locationId, securityCodes),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['billboard-control'] });
      setSnackbar({ open: true, message: 'Billboard launched successfully', type: 'success' });
    },
    onError: (error) => {
      setSnackbar({ open: true, message: 'Failed to launch billboard', type: 'error' });
    },
  });

  const clearBillboardMutation = useMutation({
    mutationFn: () => api.clearBillboard(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['billboard-control'] });
      setSnackbar({ open: true, message: 'Billboard cleared successfully', type: 'success' });
    },
    onError: (error) => {
      setSnackbar({ open: true, message: 'Failed to clear billboard', type: 'error' });
    },
  });

  // Extract data from API responses
  const events = eventsData?.events || [];
  const notifications = notificationsData?.notifications || [];
  const securityCodes = securityCodesData?.codes?.map((code: any) => code.code) || [];
  const billboardControl = billboardControlData?.control;

  // WebSocket real-time updates
  useEffect(() => {
    // Subscribe to real-time updates
    subscribeToNotifications();
    subscribeToBillboardState();

    // Listen for notification updates
    const unsubscribeNotifications = onNotificationUpdate((notification) => {
      console.log('Real-time notification update:', notification);
      queryClient.invalidateQueries({ queryKey: ['notifications', 'active'] });
      
      // Show toast for new notifications
      if (notification.child_name) {
        setSnackbar({ 
          open: true, 
          message: `New pickup request: ${notification.child_name} (${notification.security_code})`, 
          type: 'info' 
        });
      }
    });

    // Listen for billboard state changes
    const unsubscribeBillboardState = onBillboardStateChange((state) => {
      console.log('Real-time billboard state change:', state);
      queryClient.invalidateQueries({ queryKey: ['billboard-control'] });
    });

    // Listen for security code updates
    const unsubscribeSecurityCodeAdded = onSecurityCodeAdded((data) => {
      console.log('Real-time security code added:', data);
      queryClient.invalidateQueries({ queryKey: ['security-codes'] });
      setSnackbar({ open: true, message: `Security code ${data.code} added`, type: 'success' });
    });

    const unsubscribeSecurityCodeRemoved = onSecurityCodeRemoved((data) => {
      console.log('Real-time security code removed:', data);
      queryClient.invalidateQueries({ queryKey: ['security-codes'] });
      setSnackbar({ open: true, message: `Security code ${data.code} removed`, type: 'info' });
    });

    // Listen for billboard launch/clear events
    const unsubscribeBillboardLaunched = onBillboardLaunched((state) => {
      console.log('Real-time billboard launched:', state);
      queryClient.invalidateQueries({ queryKey: ['billboard-control'] });
      setSnackbar({ open: true, message: 'Billboard launched successfully', type: 'success' });
    });

    const unsubscribeBillboardCleared = onBillboardCleared((data) => {
      console.log('Real-time billboard cleared:', data);
      queryClient.invalidateQueries({ queryKey: ['billboard-control'] });
      setSnackbar({ open: true, message: 'Billboard cleared successfully', type: 'success' });
    });

    return () => {
      unsubscribeNotifications();
      unsubscribeBillboardState();
      unsubscribeSecurityCodeAdded();
      unsubscribeSecurityCodeRemoved();
      unsubscribeBillboardLaunched();
      unsubscribeBillboardCleared();
    };
  }, [
    onNotificationUpdate,
    onBillboardStateChange,
    onSecurityCodeAdded,
    onSecurityCodeRemoved,
    onBillboardLaunched,
    onBillboardCleared,
    subscribeToNotifications,
    subscribeToBillboardState,
    queryClient
  ]);

  // Handlers
  const handleAddSecurityCode = () => {
    if (!newSecurityCode.trim()) {
      setSnackbar({ open: true, message: 'Enter a security code', type: 'error' });
      return;
    }
    if (securityCodes.includes(newSecurityCode.trim())) {
      setSnackbar({ open: true, message: 'Code already exists', type: 'error' });
      return;
    }
    addSecurityCodeMutation.mutate(newSecurityCode.trim());
    setNewSecurityCode('');
  };

  const handleRemoveSecurityCode = (code: string) => {
    removeSecurityCodeMutation.mutate(code);
  };

  const handleLaunchBillboard = () => {
    if (!selectedEvent) {
      setSnackbar({ open: true, message: 'Please select an event', type: 'error' });
      return;
    }
    
    const selectedEventObj = events.find(e => e.id === selectedEvent);
    if (!selectedEventObj) {
      setSnackbar({ open: true, message: 'Selected event not found', type: 'error' });
      return;
    }

    launchBillboardMutation.mutate({
      eventId: selectedEvent,
      locationId: selectedEventObj.location,
      securityCodes: securityCodes,
    });
  };

  const handleClearBillboard = () => {
    clearBillboardMutation.mutate();
  };

  const closeSnackbar = () => setSnackbar({ ...snackbar, open: false });

  // Event selection handlers
  const handleDateChange = (date: string) => {
    setSelectedDate(date);
    setGlobalSelectedDate(date);
  };

  const handleEventChange = (eventId: string) => {
    setSelectedEvent(eventId);
    const eventObj = events.find(e => e.id === eventId);
    setGlobalSelectedEvent(eventObj || null);
  };

  // Find selected event details
  const selectedEventObj = events.find(e => e.id === selectedEvent);

  if (!user?.is_admin) {
    return (
      <div className="p-6">
        <div className="bg-red-50 border border-red-200 rounded-md p-4">
          <h3 className="text-red-800 font-medium">Access Denied</h3>
          <p className="text-red-600 mt-1">You don't have permission to access the admin panel.</p>
        </div>
      </div>
    );
  }

  const isLoading = eventsLoading || notificationsLoading || securityCodesLoading || billboardControlLoading;

  if (isLoading) {
    return (
      <div className="p-6">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-1/4 mb-4"></div>
          <div className="space-y-3">
            <div className="h-4 bg-gray-200 rounded"></div>
            <div className="h-4 bg-gray-200 rounded w-5/6"></div>
            <div className="h-4 bg-gray-200 rounded w-4/6"></div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Admin Panel</h1>
          <p className="text-gray-600">Manage events, security codes, and billboard display</p>
        </div>
        
        {/* Connection Status */}
        <div className="flex items-center space-x-2">
          <div className={`w-2 h-2 rounded-full ${connectionStatus.isConnected ? 'bg-green-500' : 'bg-red-500'}`}></div>
          <span className="text-sm text-gray-500">
            {connectionStatus.isConnected ? 'Live Updates' : 'Offline'}
          </span>
        </div>
      </div>

      {/* Event Selection */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Event Selection</h2>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label htmlFor="date" className="block text-sm font-medium text-gray-700 mb-2">
              Select Date
            </label>
            <input
              type="date"
              id="date"
              value={selectedDate}
              onChange={(e) => handleDateChange(e.target.value)}
              className="border border-gray-300 rounded-md px-3 py-2 w-full"
            />
          </div>
          
          <div>
            <div className="flex items-center justify-between mb-2">
              <label htmlFor="event" className="block text-sm font-medium text-gray-700">Select Event</label>
              {selectedEvent && (
                <button
                  onClick={() => handleEventChange('')}
                  className="text-xs text-gray-500 hover:text-gray-700 underline"
                >
                  Clear Selection
                </button>
              )}
            </div>
            <select
              id="event"
              value={selectedEvent}
              onChange={e => handleEventChange(e.target.value)}
              className="border border-gray-300 rounded-md px-3 py-2 w-full"
            >
              <option value="">Select an event...</option>
              {events.map((event) => (
                <option key={event.id} value={event.id}>
                  {event.name} - {event.time || 'All day'}
                </option>
              ))}
            </select>
          </div>
        </div>
      </div>

      {/* Security Code Management */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Security Codes</h2>
        
        <div className="flex space-x-2 mb-4">
          <input
            type="text"
            value={newSecurityCode}
            onChange={(e) => setNewSecurityCode(e.target.value)}
            placeholder="Enter security code"
            className="flex-1 border border-gray-300 rounded-md px-3 py-2"
            onKeyPress={(e) => e.key === 'Enter' && handleAddSecurityCode()}
          />
          <button
            onClick={handleAddSecurityCode}
            disabled={addSecurityCodeMutation.isPending}
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
          >
            {addSecurityCodeMutation.isPending ? 'Adding...' : 'Add'}
          </button>
        </div>
        
        <div className="flex flex-wrap gap-2">
          {securityCodes.map((code) => (
            <div
              key={code}
              className="flex items-center space-x-2 bg-gray-100 px-3 py-1 rounded-full"
            >
              <span className="text-sm font-medium">{code}</span>
              <button
                onClick={() => handleRemoveSecurityCode(code)}
                disabled={removeSecurityCodeMutation.isPending}
                className="text-red-600 hover:text-red-800 text-sm"
              >
                ×
              </button>
            </div>
          ))}
        </div>
      </div>

      {/* Billboard Controls */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Billboard Controls</h2>
        
        <div className="flex space-x-4">
          <button
            onClick={handleLaunchBillboard}
            disabled={!selectedEvent || launchBillboardMutation.isPending}
            className="px-6 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 disabled:opacity-50"
          >
            {launchBillboardMutation.isPending ? 'Launching...' : 'Launch Billboard'}
          </button>
          
          <button
            onClick={handleClearBillboard}
            disabled={clearBillboardMutation.isPending}
            className="px-6 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 disabled:opacity-50"
          >
            {clearBillboardMutation.isPending ? 'Clearing...' : 'Clear Billboard'}
          </button>
        </div>
        
        {billboardControl && (
          <div className="mt-4 p-4 bg-green-50 border border-green-200 rounded-md">
            <h3 className="font-medium text-green-800">Active Billboard</h3>
            <p className="text-green-700 text-sm">
              Event: {billboardControl.event_name} • Location: {billboardControl.location_name}
            </p>
          </div>
        )}
      </div>

      {/* Active Notifications */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Active Pickup Notifications</h2>
        
        {notifications.length === 0 ? (
          <p className="text-gray-500">No active pickup notifications</p>
        ) : (
          <div className="space-y-3">
            {notifications.map((notification) => (
              <div
                key={notification.id}
                className="flex items-center justify-between p-3 bg-gray-50 rounded-md"
              >
                <div>
                  <p className="font-medium">{notification.child_name}</p>
                  <p className="text-sm text-gray-600">
                    Code: {notification.security_code} • Location: {notification.location_name}
                  </p>
                </div>
                <span className="text-sm text-gray-500">
                  {new Date(notification.created_at).toLocaleTimeString()}
                </span>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Snackbar */}
      {snackbar.open && (
        <div className={`fixed bottom-4 right-4 p-4 rounded-md shadow-lg ${
          snackbar.type === 'success' ? 'bg-green-500 text-white' :
          snackbar.type === 'error' ? 'bg-red-500 text-white' :
          'bg-blue-500 text-white'
        }`}>
          <div className="flex items-center space-x-2">
            <span>{snackbar.message}</span>
            <button
              onClick={closeSnackbar}
              className="text-white hover:text-gray-200"
            >
              ×
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default AdminPage; 