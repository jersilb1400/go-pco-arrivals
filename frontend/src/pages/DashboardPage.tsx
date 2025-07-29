import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { motion } from 'framer-motion';
import { 
  Users, 
  Clock, 
  TrendingUp, 
  Activity,
  Plus,
  ArrowRight,
  BarChart3,
  Settings,
  Bell,
  Shield,
  Database,
  Wifi,
  WifiOff,
  RefreshCw,
  Calendar,
  MapPin,
  Eye,
  Edit,
  Trash2,
  Download,
  Upload,
  Zap,
  AlertTriangle,
  CheckCircle,
  XCircle
} from 'lucide-react';
import apiService from '../services/api';
import type { Location, SystemStatus, StatsResponse } from '../types/api';

interface DashboardStats {
  totalLocations: number;
  todayCheckIns: number;
  activeSessions: number;
  totalUsers: number;
  systemHealth: 'healthy' | 'warning' | 'error';
  lastSync: string;
  uptime: string;
}

type ChangeType = 'positive' | 'negative' | 'neutral';

const DashboardPage: React.FC = () => {
  const [activeTab, setActiveTab] = useState('overview');
  const [isRefreshing, setIsRefreshing] = useState(false);

  const { data: systemStatus, refetch: refetchSystemStatus } = useQuery({
    queryKey: ['system', 'status'],
    queryFn: apiService.getSystemStatus,
    refetchInterval: 30000, // Refresh every 30 seconds
  });

  const { data: locationsResponse, refetch: refetchLocations } = useQuery({
    queryKey: ['locations'],
    queryFn: apiService.getLocations,
  });

  const { data: pcoLocationsResponse } = useQuery({
    queryKey: ['pco-locations'],
    queryFn: apiService.getPCOLocations,
    enabled: false, // Don't fetch automatically
  });

  const { data: checkInStats } = useQuery({
    queryKey: ['check-in-stats'],
    queryFn: () => apiService.getCheckInStats('all', 7), // Last 7 days
    refetchInterval: 60000, // Refresh every minute
  });

  const locations = locationsResponse?.locations || [];

  const handleRefresh = async () => {
    setIsRefreshing(true);
    await Promise.all([
      refetchSystemStatus(),
      refetchLocations(),
    ]);
    setIsRefreshing(false);
  };

  const stats = [
    {
      name: 'Total Locations',
      value: systemStatus?.system.locations.total || 0,
      icon: MapPin,
      color: 'text-blue-600',
      bgColor: 'bg-blue-100',
      change: '+2 this week',
      changeType: 'positive' as ChangeType,
    },
    {
      name: 'Today\'s Check-ins',
      value: systemStatus?.system.check_ins.today || 0,
      icon: Clock,
      color: 'text-green-600',
      bgColor: 'bg-green-100',
      change: '+15% vs yesterday',
      changeType: 'positive' as ChangeType,
    },
    {
      name: 'Active Sessions',
      value: systemStatus?.system.sessions.active || 0,
      icon: Activity,
      color: 'text-purple-600',
      bgColor: 'bg-purple-100',
      change: '3 users online',
      changeType: 'neutral' as ChangeType,
    },
    {
      name: 'System Health',
      value: 'Healthy',
      icon: CheckCircle,
      color: 'text-green-600',
      bgColor: 'bg-green-100',
      change: 'All systems operational',
      changeType: 'positive' as ChangeType,
    },
  ];

  const quickActions = [
    {
      name: 'Add Location',
      description: 'Configure a new billboard location',
      icon: Plus,
      href: '/locations',
      color: 'text-blue-600',
      bgColor: 'bg-blue-100',
    },
    {
      name: 'Manage Locations',
      description: 'View and edit your locations',
      icon: MapPin,
      href: '/locations',
      color: 'text-green-600',
      bgColor: 'bg-green-100',
    },
    {
      name: 'View Analytics',
      description: 'Check-in statistics and trends',
      icon: BarChart3,
      href: '/analytics',
      color: 'text-purple-600',
      bgColor: 'bg-purple-100',
    },
    {
      name: 'System Settings',
      description: 'Configure system preferences',
      icon: Settings,
      href: '/settings',
      color: 'text-orange-600',
      bgColor: 'bg-orange-100',
    },
    {
      name: 'User Management',
      description: 'Manage user accounts and permissions',
      icon: Users,
      href: '/users',
      color: 'text-indigo-600',
      bgColor: 'bg-indigo-100',
    },
    {
      name: 'Notifications',
      description: 'Configure notification settings',
      icon: Bell,
      href: '/notifications',
      color: 'text-red-600',
      bgColor: 'bg-red-100',
    },
  ];

  const systemAlerts = [
    {
      type: 'info' as const,
      message: 'System running smoothly',
      icon: CheckCircle,
      color: 'text-green-600',
      bgColor: 'bg-green-50',
    },
    {
      type: 'warning' as const,
      message: '2 locations need attention',
      icon: AlertTriangle,
      color: 'text-yellow-600',
      bgColor: 'bg-yellow-50',
    },
  ];

  const recentActivity = [
    {
      type: 'check-in',
      message: 'New check-in at Main Campus',
      time: '2 minutes ago',
      icon: Users,
      color: 'text-blue-600',
    },
    {
      type: 'location',
      message: 'Location "Youth Room" added',
      time: '15 minutes ago',
      icon: MapPin,
      color: 'text-green-600',
    },
    {
      type: 'user',
      message: 'New user registered',
      time: '1 hour ago',
      icon: Users,
      color: 'text-purple-600',
    },
  ];

  const tabs = [
    { id: 'overview', name: 'Overview', icon: BarChart3 },
    { id: 'locations', name: 'Locations', icon: MapPin },
    { id: 'analytics', name: 'Analytics', icon: TrendingUp },
    { id: 'system', name: 'System', icon: Settings },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Admin Dashboard</h1>
          <p className="mt-1 text-sm text-gray-500">
            Welcome back! Here's what's happening with your PCO Arrivals Billboard system.
          </p>
        </div>
        <div className="flex items-center space-x-3">
          <button
            onClick={handleRefresh}
            disabled={isRefreshing}
            className="btn-secondary flex items-center space-x-2"
          >
            <RefreshCw className={`h-4 w-4 ${isRefreshing ? 'animate-spin' : ''}`} />
            <span>Refresh</span>
          </button>
        </div>
      </div>

      {/* System Alerts */}
      {systemAlerts.length > 0 && (
        <div className="space-y-3">
          {systemAlerts.map((alert, index) => {
            const Icon = alert.icon;
            return (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: -10 }}
                animate={{ opacity: 1, y: 0 }}
                className={`p-4 rounded-lg border ${alert.bgColor} border-l-4 ${alert.color.replace('text-', 'border-l-')}`}
              >
                <div className="flex items-center">
                  <Icon className={`h-5 w-5 ${alert.color} mr-3`} />
                  <p className="text-sm font-medium text-gray-900">{alert.message}</p>
                </div>
              </motion.div>
            );
          })}
        </div>
      )}

      {/* Stats Grid */}
      <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat, index) => {
          const Icon = stat.icon;
          return (
            <motion.div
              key={stat.name}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.1 }}
              className="card p-6"
            >
              <div className="flex items-center">
                <div className={`flex-shrink-0 ${stat.bgColor} rounded-md p-3`}>
                  <Icon className={`h-6 w-6 ${stat.color}`} />
                </div>
                <div className="ml-5 w-0 flex-1">
                  <dl>
                    <dt className="text-sm font-medium text-gray-500 truncate">
                      {stat.name}
                    </dt>
                    <dd className="text-lg font-medium text-gray-900">
                      {stat.value}
                    </dd>
                    <dd className={`text-xs ${
                      stat.changeType === 'positive' ? 'text-green-600' :
                      stat.changeType === 'negative' ? 'text-red-600' :
                      'text-gray-500'
                    }`}>
                      {stat.change}
                    </dd>
                  </dl>
                </div>
              </div>
            </motion.div>
          );
        })}
      </div>

      {/* Quick Actions */}
      <div className="card p-6">
        <h2 className="text-lg font-medium text-gray-900 mb-4">Quick Actions</h2>
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {quickActions.map((action, index) => {
            const Icon = action.icon;
            return (
              <motion.div
                key={action.name}
                initial={{ opacity: 0, scale: 0.95 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ delay: index * 0.05 }}
              >
                <Link
                  to={action.href}
                  className="group relative rounded-lg border border-gray-300 bg-white p-6 hover:border-primary-500 hover:shadow-md transition-all duration-200"
                >
                  <div className="flex items-center">
                    <div className="flex-shrink-0">
                      <Icon className={`h-6 w-6 ${action.color}`} />
                    </div>
                    <div className="ml-4">
                      <h3 className="text-sm font-medium text-gray-900 group-hover:text-primary-600">
                        {action.name}
                      </h3>
                      <p className="text-xs text-gray-500">
                        {action.description}
                      </p>
                    </div>
                  </div>
                  <div className="absolute top-4 right-4 opacity-0 group-hover:opacity-100 transition-opacity">
                    <ArrowRight className="h-4 w-4 text-gray-400" />
                  </div>
                </Link>
              </motion.div>
            );
          })}
        </div>
      </div>

      {/* Main Content Tabs */}
      <div className="card">
        <div className="border-b border-gray-200">
          <nav className="-mb-px flex space-x-8">
            {tabs.map((tab) => {
              const Icon = tab.icon;
              return (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className={`py-4 px-1 border-b-2 font-medium text-sm flex items-center space-x-2 ${
                    activeTab === tab.id
                      ? 'border-primary-500 text-primary-600'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  }`}
                >
                  <Icon className="h-4 w-4" />
                  <span>{tab.name}</span>
                </button>
              );
            })}
          </nav>
        </div>

        <div className="p-6">
          {activeTab === 'overview' && (
            <div className="space-y-6">
              {/* Recent Activity */}
              <div>
                <h3 className="text-lg font-medium text-gray-900 mb-4">Recent Activity</h3>
                <div className="space-y-3">
                  {recentActivity.map((activity, index) => {
                    const Icon = activity.icon;
                    return (
                      <motion.div
                        key={index}
                        initial={{ opacity: 0, x: -20 }}
                        animate={{ opacity: 1, x: 0 }}
                        transition={{ delay: index * 0.1 }}
                        className="flex items-center space-x-3 p-3 bg-gray-50 rounded-lg"
                      >
                        <Icon className={`h-4 w-4 ${activity.color}`} />
                        <div className="flex-1">
                          <p className="text-sm text-gray-900">{activity.message}</p>
                          <p className="text-xs text-gray-500">{activity.time}</p>
                        </div>
                      </motion.div>
                    );
                  })}
                </div>
              </div>

              {/* System Status */}
              <div>
                <h3 className="text-lg font-medium text-gray-900 mb-4">System Status</h3>
                <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
                  <div className="text-center p-4 bg-green-50 rounded-lg">
                    <CheckCircle className="h-8 w-8 text-green-600 mx-auto mb-2" />
                    <p className="text-sm font-medium text-green-900">Database</p>
                    <p className="text-xs text-green-600">Connected</p>
                  </div>
                  <div className="text-center p-4 bg-green-50 rounded-lg">
                    <Wifi className="h-8 w-8 text-green-600 mx-auto mb-2" />
                    <p className="text-sm font-medium text-green-900">WebSocket</p>
                    <p className="text-xs text-green-600">Active</p>
                  </div>
                  <div className="text-center p-4 bg-green-50 rounded-lg">
                    <Zap className="h-8 w-8 text-green-600 mx-auto mb-2" />
                    <p className="text-sm font-medium text-green-900">PCO API</p>
                    <p className="text-xs text-green-600">Connected</p>
                  </div>
                </div>
              </div>
            </div>
          )}

          {activeTab === 'locations' && (
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <h3 className="text-lg font-medium text-gray-900">Recent Locations</h3>
                <Link to="/locations" className="btn-primary text-sm">
                  View All
                </Link>
              </div>
              <div className="space-y-4">
                {locations.slice(0, 5).map((location) => (
                  <motion.div
                    key={location.id}
                    initial={{ opacity: 0, y: 10 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
                  >
                    <div className="flex items-center">
                      <div className="flex-shrink-0">
                        <div className="h-8 w-8 rounded-full bg-primary-100 flex items-center justify-center">
                          <MapPin className="h-4 w-4 text-primary-600" />
                        </div>
                      </div>
                      <div className="ml-4">
                        <h3 className="text-sm font-medium text-gray-900">
                          {location.name}
                        </h3>
                        <p className="text-xs text-gray-500">
                          ID: {location.pco_location_id}
                        </p>
                      </div>
                    </div>
                    <div className="flex items-center space-x-2">
                      <Link
                        to={`/billboard/${location.pco_location_id}`}
                        className="btn-secondary text-xs"
                      >
                        <Eye className="h-3 w-3 mr-1" />
                        View
                      </Link>
                      <Link
                        to={`/locations/${location.id}/edit`}
                        className="btn-secondary text-xs"
                      >
                        <Edit className="h-3 w-3 mr-1" />
                        Edit
                      </Link>
                    </div>
                  </motion.div>
                ))}
              </div>
            </div>
          )}

          {activeTab === 'analytics' && (
            <div className="space-y-6">
              <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
                <div className="text-center p-4 bg-blue-50 rounded-lg">
                  <Users className="h-8 w-8 text-blue-600 mx-auto mb-2" />
                                     <p className="text-2xl font-bold text-blue-900">
                     {checkInStats?.stats?.total_check_ins || 0}
                   </p>
                   <p className="text-sm text-blue-600">Total Check-ins</p>
                 </div>
                 <div className="text-center p-4 bg-green-50 rounded-lg">
                   <TrendingUp className="h-8 w-8 text-green-600 mx-auto mb-2" />
                   <p className="text-2xl font-bold text-green-900">
                     {checkInStats?.stats?.today_check_ins || 0}
                   </p>
                   <p className="text-sm text-green-600">Today's Check-ins</p>
                 </div>
                 <div className="text-center p-4 bg-purple-50 rounded-lg">
                   <MapPin className="h-8 w-8 text-purple-600 mx-auto mb-2" />
                   <p className="text-2xl font-bold text-purple-900">
                     {locations.length}
                   </p>
                   <p className="text-sm text-purple-600">Active Locations</p>
                 </div>
                 <div className="text-center p-4 bg-orange-50 rounded-lg">
                   <Calendar className="h-8 w-8 text-orange-600 mx-auto mb-2" />
                   <p className="text-2xl font-bold text-orange-900">
                     {checkInStats?.stats?.weekly_check_ins || 0}
                   </p>
                   <p className="text-sm text-orange-600">Weekly Check-ins</p>
                </div>
              </div>
              
              <div className="text-center py-8">
                <BarChart3 className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                <p className="text-gray-500">Detailed analytics coming soon...</p>
              </div>
            </div>
          )}

          {activeTab === 'system' && (
            <div className="space-y-6">
              <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                <div className="p-4 border border-gray-200 rounded-lg">
                  <h4 className="font-medium text-gray-900 mb-2">System Information</h4>
                  <div className="space-y-2 text-sm">
                    <div className="flex justify-between">
                      <span className="text-gray-500">Version:</span>
                      <span className="text-gray-900">1.0.0</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-500">Uptime:</span>
                      <span className="text-gray-900">2 days, 14 hours</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-500">Last Sync:</span>
                      <span className="text-gray-900">2 minutes ago</span>
                    </div>
                  </div>
                </div>
                
                <div className="p-4 border border-gray-200 rounded-lg">
                  <h4 className="font-medium text-gray-900 mb-2">Database Status</h4>
                  <div className="space-y-2 text-sm">
                    <div className="flex justify-between">
                      <span className="text-gray-500">Status:</span>
                      <span className="text-green-600 flex items-center">
                        <CheckCircle className="h-3 w-3 mr-1" />
                        Connected
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-500">Type:</span>
                      <span className="text-gray-900">SQLite</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-500">Size:</span>
                      <span className="text-gray-900">2.4 MB</span>
                    </div>
                  </div>
                </div>
              </div>

              <div className="flex space-x-3">
                <button className="btn-secondary">
                  <Download className="h-4 w-4 mr-2" />
                  Export Data
                </button>
                <button className="btn-secondary">
                  <Upload className="h-4 w-4 mr-2" />
                  Import Data
                </button>
                <button className="btn-secondary">
                  <Database className="h-4 w-4 mr-2" />
                  Backup Database
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default DashboardPage; 