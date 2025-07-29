import React, { useState, useEffect, useMemo } from 'react';
import { useQuery } from '@tanstack/react-query';
import { motion, AnimatePresence } from 'framer-motion';
import { 
  Users, 
  Clock, 
  RefreshCw, 
  MapPin,
  AlertCircle,
  Calendar,
  Building,
  Shield,
  TrendingUp,
  TrendingDown,
  Activity,
  BarChart3,
  Filter,
  Search,
  Settings,
  Eye,
  EyeOff,
  Zap,
  Target,
  Timer,
  CheckCircle,
  XCircle,
  AlertTriangle,
  Info
} from 'lucide-react';
import apiService from '../services/api';
import { useApp } from '../contexts/AppContext';
import { useWebSocket } from '../contexts/WebSocketContext';
import type { 
  Notification, 
  LocationOverview, 
  LocationStatus as ApiLocationStatus, 
  LocationAnalytics,
  LocationsOverviewResponse 
} from '../types/api';

interface LocationStatus {
  location_id: string;
  location_name: string;
  child_count: number;
  children: Notification[];
}

const LocationStatusPage: React.FC = () => {
  const { state: appState } = useApp();
  const { 
    connectionStatus,
    onNotificationUpdate,
    onBillboardStateChange,
    subscribeToNotifications,
    subscribeToBillboardState
  } = useWebSocket();
  
  // State management
  const [lastUpdated, setLastUpdated] = useState<Date>(new Date());
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [selectedLocation, setSelectedLocation] = useState<string | null>(null);
  const [viewMode, setViewMode] = useState<'overview' | 'detailed' | 'analytics'>('overview');
  const [searchTerm, setSearchTerm] = useState('');
  const [filterStatus, setFilterStatus] = useState<'all' | 'active' | 'inactive'>('all');
  const [showOfflineLocations, setShowOfflineLocations] = useState(true);
  const [analyticsPeriod, setAnalyticsPeriod] = useState<number>(30);

  // Get global billboard state (active event)
  const { data: billboardControl } = useQuery({
    queryKey: ['billboard-control'],
    queryFn: () => apiService.getBillboardControl(),
    refetchInterval: 30000, // Poll every 30 seconds as backup
  });

  // Get locations overview
  const { data: locationsOverview, refetch: refetchLocations } = useQuery({
    queryKey: ['locations-overview'],
    queryFn: () => apiService.getLocationsOverview(),
    refetchInterval: 60000, // Poll every 60 seconds as backup
  });

  // Get detailed location status (when a location is selected)
  const { data: locationStatus } = useQuery({
    queryKey: ['location-status', selectedLocation],
    queryFn: () => selectedLocation ? apiService.getLocationStatus(selectedLocation) : null,
    enabled: !!selectedLocation && viewMode === 'detailed',
    refetchInterval: 30000,
  });

  // Get location analytics (when a location is selected)
  const { data: locationAnalytics } = useQuery({
    queryKey: ['location-analytics', selectedLocation, analyticsPeriod],
    queryFn: () => selectedLocation ? apiService.getLocationAnalytics(selectedLocation, analyticsPeriod) : null,
    enabled: !!selectedLocation && viewMode === 'analytics',
  });

  // Subscribe to real-time updates
  useEffect(() => {
    subscribeToNotifications();
    subscribeToBillboardState();
  }, [subscribeToNotifications, subscribeToBillboardState]);

  // Handle real-time updates
  useEffect(() => {
    const unsubscribeNotifications = onNotificationUpdate((notification) => {
      console.log('Real-time notification update:', notification);
      refetchLocations();
      setLastUpdated(new Date());
    });

    const unsubscribeBillboardState = onBillboardStateChange((state) => {
      console.log('Real-time billboard state change:', state);
      setLastUpdated(new Date());
    });

    return () => {
      unsubscribeNotifications();
      unsubscribeBillboardState();
    };
  }, [onNotificationUpdate, onBillboardStateChange, refetchLocations]);

  // Filter and search locations
  const filteredLocations = useMemo(() => {
    if (!locationsOverview?.locations) return [];
    
    return locationsOverview.locations.filter(location => {
      const matchesSearch = location.name.toLowerCase().includes(searchTerm.toLowerCase());
      const matchesStatus = filterStatus === 'all' || 
        (filterStatus === 'active' && location.active_children > 0) ||
        (filterStatus === 'inactive' && location.active_children === 0);
      const matchesOnline = showOfflineLocations || location.is_active;
      
      return matchesSearch && matchesStatus && matchesOnline;
    });
  }, [locationsOverview?.locations, searchTerm, filterStatus, showOfflineLocations]);

  // Group notifications by location (for detailed view)
  const locationStatuses = useMemo(() => {
    if (!locationStatus?.location?.notifications) return [];
    
    const notifications = locationStatus.location.notifications;
    const locationMap = new Map<string, LocationStatus>();

    notifications.forEach(notification => {
      const locationId = notification.location_name || 'Unknown Location';
      
      if (!locationMap.has(locationId)) {
        locationMap.set(locationId, {
          location_id: locationId,
          location_name: locationId,
          child_count: 0,
          children: [],
        });
      }

      const location = locationMap.get(locationId)!;
      location.child_count++;
      location.children.push(notification);
    });

    return Array.from(locationMap.values()).sort((a, b) => b.child_count - a.child_count);
  }, [locationStatus?.location?.notifications]);

  const handleRefresh = async () => {
    setIsRefreshing(true);
    try {
      await refetchLocations();
      setLastUpdated(new Date());
    } catch (error) {
      console.error('Failed to refresh:', error);
    } finally {
      setIsRefreshing(false);
    }
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

  const getEfficiencyColor = (rate: number) => {
    if (rate >= 90) return 'text-green-600';
    if (rate >= 70) return 'text-yellow-600';
    return 'text-red-600';
  };

  const getEfficiencyIcon = (rate: number) => {
    if (rate >= 90) return <CheckCircle className="h-4 w-4 text-green-600" />;
    if (rate >= 70) return <AlertTriangle className="h-4 w-4 text-yellow-600" />;
    return <XCircle className="h-4 w-4 text-red-600" />;
  };

  // Check if billboard is active
  const isBillboardActive = billboardControl?.control?.is_active;
  const activeEvent = billboardControl?.control;

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
      {/* Enhanced Header */}
      <div className="bg-white border-b border-gray-200 px-6 py-4 shadow-sm">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Location Status Dashboard</h1>
            {activeEvent && (
              <p className="mt-1 text-lg text-gray-600">
                {activeEvent.event_name} • {appState.selectedDate}
              </p>
            )}
            <p className="mt-1 text-sm text-gray-500">
              {filteredLocations.length} location{filteredLocations.length !== 1 ? 's' : ''} • 
              {locationsOverview?.summary?.total_children || 0} children awaiting pickup
            </p>
          </div>
          
          <div className="flex items-center space-x-4">
            {/* Connection Status */}
            <div className="flex items-center space-x-2">
              <div className={`w-3 h-3 rounded-full ${connectionStatus.isConnected ? 'bg-green-500 animate-pulse' : 'bg-red-500'}`}></div>
              <span className="text-sm text-gray-500">
                {connectionStatus.isConnected ? 'Live Updates' : 'Offline'}
              </span>
            </div>

            {/* Last Updated */}
            <div className="text-sm text-gray-500">
              Last updated: {lastUpdated.toLocaleTimeString()}
            </div>

            {/* Refresh Button */}
            <button
              onClick={handleRefresh}
              disabled={isRefreshing}
              className="flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              <RefreshCw className={`h-4 w-4 mr-2 ${isRefreshing ? 'animate-spin' : ''}`} />
              {isRefreshing ? 'Refreshing...' : 'Refresh'}
            </button>
          </div>
        </div>

        {/* View Mode Tabs */}
        <div className="mt-4 flex items-center space-x-4">
          <div className="flex space-x-1 bg-gray-100 rounded-lg p-1">
            {[
              { key: 'overview', label: 'Overview', icon: <Building className="h-4 w-4" /> },
              { key: 'detailed', label: 'Detailed', icon: <Eye className="h-4 w-4" /> },
              { key: 'analytics', label: 'Analytics', icon: <BarChart3 className="h-4 w-4" /> }
            ].map(({ key, label, icon }) => (
              <button
                key={key}
                onClick={() => setViewMode(key as any)}
                className={`flex items-center px-3 py-2 rounded-md text-sm font-medium transition-colors ${
                  viewMode === key 
                    ? 'bg-white text-blue-600 shadow-sm' 
                    : 'text-gray-600 hover:text-gray-900'
                }`}
              >
                {icon}
                <span className="ml-2">{label}</span>
              </button>
            ))}
          </div>

          {/* Filters */}
          <div className="flex items-center space-x-4">
            {/* Search */}
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
              <input
                type="text"
                placeholder="Search locations..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-10 pr-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>

            {/* Status Filter */}
            <select
              value={filterStatus}
              onChange={(e) => setFilterStatus(e.target.value as any)}
              className="px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="all">All Locations</option>
              <option value="active">Active Only</option>
              <option value="inactive">Inactive Only</option>
            </select>

            {/* Online/Offline Toggle */}
            <button
              onClick={() => setShowOfflineLocations(!showOfflineLocations)}
              className={`flex items-center px-3 py-2 rounded-md text-sm font-medium transition-colors ${
                showOfflineLocations 
                  ? 'bg-green-100 text-green-700' 
                  : 'bg-gray-100 text-gray-700'
              }`}
            >
              {showOfflineLocations ? <Eye className="h-4 w-4 mr-1" /> : <EyeOff className="h-4 w-4 mr-1" />}
              {showOfflineLocations ? 'Show All' : 'Online Only'}
            </button>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="p-6">
        {!isBillboardActive ? (
          // No active billboard
          <div className="flex items-center justify-center h-64">
            <div className="text-center">
              <AlertCircle className="mx-auto h-12 w-12 text-gray-400" />
              <h3 className="mt-2 text-lg font-medium text-gray-900">
                No Active Billboard
              </h3>
              <p className="mt-1 text-sm text-gray-500">
                Please launch a billboard from the admin panel to view location status.
              </p>
            </div>
          </div>
        ) : viewMode === 'overview' ? (
          // Overview Mode
          <div className="space-y-6">
            {/* Summary Cards */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                className="bg-white rounded-lg shadow-sm border border-gray-200 p-6"
              >
                <div className="flex items-center">
                  <div className="flex-shrink-0">
                    <Building className="h-8 w-8 text-blue-600" />
                  </div>
                  <div className="ml-4">
                    <p className="text-sm font-medium text-gray-500">Total Locations</p>
                    <p className="text-2xl font-bold text-gray-900">
                      {locationsOverview?.summary?.total_locations || 0}
                    </p>
                  </div>
                </div>
              </motion.div>

              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.1 }}
                className="bg-white rounded-lg shadow-sm border border-gray-200 p-6"
              >
                <div className="flex items-center">
                  <div className="flex-shrink-0">
                    <Users className="h-8 w-8 text-green-600" />
                  </div>
                  <div className="ml-4">
                    <p className="text-sm font-medium text-gray-500">Active Children</p>
                    <p className="text-2xl font-bold text-gray-900">
                      {locationsOverview?.summary?.total_children || 0}
                    </p>
                  </div>
                </div>
              </motion.div>

              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.2 }}
                className="bg-white rounded-lg shadow-sm border border-gray-200 p-6"
              >
                <div className="flex items-center">
                  <div className="flex-shrink-0">
                    <Activity className="h-8 w-8 text-purple-600" />
                  </div>
                  <div className="ml-4">
                    <p className="text-sm font-medium text-gray-500">Active Locations</p>
                    <p className="text-2xl font-bold text-gray-900">
                      {filteredLocations.filter(l => l.active_children > 0).length}
                    </p>
                  </div>
                </div>
              </motion.div>

              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.3 }}
                className="bg-white rounded-lg shadow-sm border border-gray-200 p-6"
              >
                <div className="flex items-center">
                  <div className="flex-shrink-0">
                    <Zap className="h-8 w-8 text-orange-600" />
                  </div>
                  <div className="ml-4">
                    <p className="text-sm font-medium text-gray-500">Today's Check-ins</p>
                    <p className="text-2xl font-bold text-gray-900">
                      {filteredLocations.reduce((sum, l) => sum + l.today_check_ins, 0)}
                    </p>
                  </div>
                </div>
              </motion.div>
            </div>

            {/* Location Cards Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              <AnimatePresence>
                {filteredLocations.map((location, index) => (
                  <motion.div
                    key={location.id}
                    initial={{ opacity: 0, scale: 0.9 }}
                    animate={{ opacity: 1, scale: 1 }}
                    exit={{ opacity: 0, scale: 0.9 }}
                    transition={{ duration: 0.3, delay: index * 0.1 }}
                    className="bg-white rounded-lg shadow-sm border border-gray-200 hover:shadow-md transition-shadow cursor-pointer"
                    onClick={() => {
                      setSelectedLocation(location.id);
                      setViewMode('detailed');
                    }}
                  >
                    <div className="p-6">
                      <div className="flex items-center justify-between mb-4">
                        <div className="flex items-center space-x-3">
                          <div className="flex-shrink-0">
                            <div className={`h-10 w-10 rounded-full flex items-center justify-center ${
                              location.active_children > 0 ? 'bg-red-100' : 'bg-green-100'
                            }`}>
                              <MapPin className={`h-5 w-5 ${
                                location.active_children > 0 ? 'text-red-600' : 'text-green-600'
                              }`} />
                            </div>
                          </div>
                          <div>
                            <h3 className="text-lg font-semibold text-gray-900">
                              {location.name}
                            </h3>
                            <p className="text-sm text-gray-500">
                              {location.active_children} children awaiting pickup
                            </p>
                          </div>
                        </div>
                        <div className="flex items-center space-x-2">
                          {location.active_children > 0 && (
                            <div className="px-2 py-1 bg-red-100 text-red-800 rounded-full text-xs font-medium">
                              {location.active_children}
                            </div>
                          )}
                        </div>
                      </div>

                      <div className="grid grid-cols-2 gap-4 text-sm">
                        <div>
                          <p className="text-gray-500">Today's Check-ins</p>
                          <p className="font-semibold text-gray-900">{location.today_check_ins}</p>
                        </div>
                        <div>
                          <p className="text-gray-500">Total Check-ins</p>
                          <p className="font-semibold text-gray-900">{location.total_check_ins}</p>
                        </div>
                      </div>

                      <div className="mt-4 pt-4 border-t border-gray-100">
                        <div className="flex items-center justify-between text-xs text-gray-500">
                          <span>Last updated: {formatTimeAgo(location.last_updated)}</span>
                          <div className={`w-2 h-2 rounded-full ${
                            location.is_active ? 'bg-green-500' : 'bg-gray-400'
                          }`}></div>
                        </div>
                      </div>
                    </div>
                  </motion.div>
                ))}
              </AnimatePresence>
            </div>
          </div>
        ) : viewMode === 'detailed' && selectedLocation ? (
          // Detailed View
          <div className="space-y-6">
            {/* Location Header */}
            <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
              <div className="flex items-center justify-between">
                <div>
                  <h2 className="text-2xl font-bold text-gray-900">
                    {locationStatus?.location?.name || 'Location Details'}
                  </h2>
                  <p className="text-gray-600">
                    Detailed status and recent activity
                  </p>
                </div>
                <button
                  onClick={() => setViewMode('overview')}
                  className="px-4 py-2 text-gray-600 hover:text-gray-900"
                >
                  ← Back to Overview
                </button>
              </div>
            </div>

            {/* Location Metrics */}
            {locationStatus?.location && (
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
                  <div className="flex items-center">
                    <Users className="h-8 w-8 text-blue-600" />
                    <div className="ml-4">
                      <p className="text-sm font-medium text-gray-500">Active Children</p>
                      <p className="text-2xl font-bold text-gray-900">
                        {locationStatus.location.active_children}
                      </p>
                    </div>
                  </div>
                </div>

                <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
                  <div className="flex items-center">
                    <Clock className="h-8 w-8 text-green-600" />
                    <div className="ml-4">
                      <p className="text-sm font-medium text-gray-500">Today's Check-ins</p>
                      <p className="text-2xl font-bold text-gray-900">
                        {locationStatus.location.today_check_ins}
                      </p>
                    </div>
                  </div>
                </div>

                <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
                  <div className="flex items-center">
                    <Target className="h-8 w-8 text-purple-600" />
                    <div className="ml-4">
                      <p className="text-sm font-medium text-gray-500">Total Check-ins</p>
                      <p className="text-2xl font-bold text-gray-900">
                        {locationStatus.location.total_check_ins}
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* Children List */}
            {locationStatuses.length > 0 ? (
              <div className="bg-white rounded-lg shadow-sm border border-gray-200">
                <div className="px-6 py-4 border-b border-gray-200">
                  <h3 className="text-lg font-semibold text-gray-900">
                    Children Awaiting Pickup
                  </h3>
                </div>
                <div className="p-6">
                  <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                    {locationStatuses[0]?.children.map((child, childIndex) => (
                      <motion.div
                        key={child.id}
                        initial={{ opacity: 0, scale: 0.9 }}
                        animate={{ opacity: 1, scale: 1 }}
                        exit={{ opacity: 0, scale: 0.9 }}
                        transition={{ duration: 0.2, delay: childIndex * 0.05 }}
                        className="bg-gray-50 border border-gray-200 rounded-lg p-4"
                      >
                        <div className="flex items-start justify-between">
                          <div className="flex-1">
                            <h4 className="text-lg font-semibold text-gray-900 mb-2">
                              {child.child_name || 'Unknown Child'}
                            </h4>
                            
                            {child.security_code && (
                              <div className="mb-2">
                                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                                  <Shield className="h-3 w-3 mr-1" />
                                  Code: {child.security_code}
                                </span>
                              </div>
                            )}
                            
                            {child.parent_name && (
                              <p className="text-sm text-gray-600 mb-1">
                                Parent: {child.parent_name}
                              </p>
                            )}
                            
                            <div className="flex items-center text-sm text-gray-500">
                              <Clock className="h-3 w-3 mr-1" />
                              {formatTimeAgo(child.created_at)}
                            </div>
                          </div>
                          
                          <div className="ml-4 p-2 rounded-full bg-orange-100">
                            <Users className="h-5 w-5 text-orange-600" />
                          </div>
                        </div>
                      </motion.div>
                    ))}
                  </div>
                </div>
              </div>
            ) : (
              <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-12 text-center">
                <Building className="mx-auto h-12 w-12 text-gray-400" />
                <h3 className="mt-2 text-lg font-medium text-gray-900">
                  No Children Awaiting Pickup
                </h3>
                <p className="mt-1 text-sm text-gray-500">
                  All children have been picked up or no pickup requests are active.
                </p>
              </div>
            )}
          </div>
        ) : viewMode === 'analytics' && selectedLocation ? (
          // Analytics View
          <div className="space-y-6">
            {/* Analytics Header */}
            <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
              <div className="flex items-center justify-between">
                <div>
                  <h2 className="text-2xl font-bold text-gray-900">
                    {locationStatus?.location?.name || 'Location'} Analytics
                  </h2>
                  <p className="text-gray-600">
                    Performance metrics and trends
                  </p>
                </div>
                <div className="flex items-center space-x-4">
                  <select
                    value={analyticsPeriod}
                    onChange={(e) => setAnalyticsPeriod(Number(e.target.value))}
                    className="px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  >
                    <option value={7}>Last 7 days</option>
                    <option value={30}>Last 30 days</option>
                    <option value={90}>Last 90 days</option>
                  </select>
                  <button
                    onClick={() => setViewMode('overview')}
                    className="px-4 py-2 text-gray-600 hover:text-gray-900"
                  >
                    ← Back to Overview
                  </button>
                </div>
              </div>
            </div>

            {/* Analytics Metrics */}
            {locationAnalytics?.analytics && (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
                  <div className="flex items-center">
                    <Timer className="h-8 w-8 text-blue-600" />
                    <div className="ml-4">
                      <p className="text-sm font-medium text-gray-500">Avg Wait Time</p>
                      <p className="text-2xl font-bold text-gray-900">
                        {Math.round(locationAnalytics.analytics.avg_wait_time_mins)}m
                      </p>
                    </div>
                  </div>
                </div>

                <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
                  <div className="flex items-center">
                    <Target className="h-8 w-8 text-green-600" />
                    <div className="ml-4">
                      <p className="text-sm font-medium text-gray-500">Efficiency Rate</p>
                      <div className="flex items-center">
                        <p className={`text-2xl font-bold ${getEfficiencyColor(locationAnalytics.analytics.efficiency_rate)}`}>
                          {Math.round(locationAnalytics.analytics.efficiency_rate)}%
                        </p>
                        <span className="ml-2">
                          {getEfficiencyIcon(locationAnalytics.analytics.efficiency_rate)}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>

                <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
                  <div className="flex items-center">
                    <CheckCircle className="h-8 w-8 text-purple-600" />
                    <div className="ml-4">
                      <p className="text-sm font-medium text-gray-500">Completed Pickups</p>
                      <p className="text-2xl font-bold text-gray-900">
                        {locationAnalytics.analytics.completed_pickups}
                      </p>
                    </div>
                  </div>
                </div>

                <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
                  <div className="flex items-center">
                    <BarChart3 className="h-8 w-8 text-orange-600" />
                    <div className="ml-4">
                      <p className="text-sm font-medium text-gray-500">Total Notifications</p>
                      <p className="text-2xl font-bold text-gray-900">
                        {locationAnalytics.analytics.total_notifications}
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* Analytics Charts Placeholder */}
            <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">Performance Trends</h3>
              <div className="h-64 bg-gray-50 rounded-lg flex items-center justify-center">
                <div className="text-center">
                  <BarChart3 className="mx-auto h-12 w-12 text-gray-400" />
                  <p className="mt-2 text-sm text-gray-500">Analytics charts coming soon</p>
                </div>
              </div>
            </div>
          </div>
        ) : (
          // No locations with children
          <div className="flex items-center justify-center h-64">
            <div className="text-center">
              <Building className="mx-auto h-12 w-12 text-gray-400" />
              <h3 className="mt-2 text-lg font-medium text-gray-900">
                No Children Awaiting Pickup
              </h3>
              <p className="mt-1 text-sm text-gray-500">
                All children have been picked up or no pickup requests are active.
              </p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default LocationStatusPage; 