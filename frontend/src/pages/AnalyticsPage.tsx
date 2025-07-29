import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { motion } from 'framer-motion';
import { 
  BarChart3, 
  TrendingUp, 
  Users, 
  MapPin, 
  Calendar,
  Clock,
  Download,
  Filter,
  RefreshCw,
  Eye,
  EyeOff
} from 'lucide-react';
import apiService from '../services/api';
import type { StatsResponse } from '../types/api';

const AnalyticsPage: React.FC = () => {
  const [selectedPeriod, setSelectedPeriod] = useState(7);
  const [selectedLocation, setSelectedLocation] = useState('all');
  const [isRefreshing, setIsRefreshing] = useState(false);

  const { data: checkInStats, refetch: refetchStats } = useQuery({
    queryKey: ['check-in-stats', selectedLocation, selectedPeriod],
    queryFn: () => apiService.getCheckInStats(selectedLocation, selectedPeriod),
    refetchInterval: 300000, // Refresh every 5 minutes
  });

  const { data: locationsResponse } = useQuery({
    queryKey: ['locations'],
    queryFn: apiService.getLocations,
  });

  const locations = locationsResponse?.locations || [];

  const handleRefresh = async () => {
    setIsRefreshing(true);
    await refetchStats();
    setIsRefreshing(false);
  };

  const periods = [
    { value: 1, label: 'Last 24 hours' },
    { value: 7, label: 'Last 7 days' },
    { value: 30, label: 'Last 30 days' },
    { value: 90, label: 'Last 90 days' },
  ];

  const analyticsCards = [
    {
      title: 'Total Check-ins',
      value: checkInStats?.stats?.total_check_ins || 0,
      change: '+12%',
      changeType: 'positive' as const,
      icon: Users,
      color: 'text-blue-600',
      bgColor: 'bg-blue-100',
    },
    {
      title: 'Today\'s Check-ins',
      value: checkInStats?.stats?.today_check_ins || 0,
      change: '+5%',
      changeType: 'positive' as const,
      icon: Clock,
      color: 'text-green-600',
      bgColor: 'bg-green-100',
    },
    {
      title: 'Weekly Check-ins',
      value: checkInStats?.stats?.weekly_check_ins || 0,
      change: '+8%',
      changeType: 'positive' as const,
      icon: Calendar,
      color: 'text-purple-600',
      bgColor: 'bg-purple-100',
    },
    {
      title: 'Active Locations',
      value: locations.length,
      change: 'No change',
      changeType: 'neutral' as const,
      icon: MapPin,
      color: 'text-orange-600',
      bgColor: 'bg-orange-100',
    },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Analytics</h1>
          <p className="mt-1 text-sm text-gray-500">
            Comprehensive check-in statistics and trends
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
          <button className="btn-secondary flex items-center space-x-2">
            <Download className="h-4 w-4" />
            <span>Export</span>
          </button>
        </div>
      </div>

      {/* Filters */}
      <div className="card p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-medium text-gray-900">Filters</h2>
          <Filter className="h-5 w-5 text-gray-400" />
        </div>
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Time Period
            </label>
            <select
              value={selectedPeriod}
              onChange={(e) => setSelectedPeriod(Number(e.target.value))}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            >
              {periods.map((period) => (
                <option key={period.value} value={period.value}>
                  {period.label}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Location
            </label>
            <select
              value={selectedLocation}
              onChange={(e) => setSelectedLocation(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            >
              <option value="all">All Locations</option>
              {locations.map((location) => (
                <option key={location.id} value={location.pco_location_id}>
                  {location.name}
                </option>
              ))}
            </select>
          </div>
        </div>
      </div>

      {/* Analytics Cards */}
      <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
        {analyticsCards.map((card, index) => {
          const Icon = card.icon;
          return (
            <motion.div
              key={card.title}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.1 }}
              className="card p-6"
            >
              <div className="flex items-center">
                <div className={`flex-shrink-0 ${card.bgColor} rounded-md p-3`}>
                  <Icon className={`h-6 w-6 ${card.color}`} />
                </div>
                <div className="ml-5 w-0 flex-1">
                  <dl>
                    <dt className="text-sm font-medium text-gray-500 truncate">
                      {card.title}
                    </dt>
                    <dd className="text-lg font-medium text-gray-900">
                      {card.value.toLocaleString()}
                    </dd>
                    <dd className={`text-xs ${
                      card.changeType === 'positive' ? 'text-green-600' :
                      card.changeType === 'negative' ? 'text-red-600' :
                      'text-gray-500'
                    }`}>
                      {card.change}
                    </dd>
                  </dl>
                </div>
              </div>
            </motion.div>
          );
        })}
      </div>

      {/* Charts Section */}
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        {/* Check-ins Over Time */}
        <div className="card p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-medium text-gray-900">Check-ins Over Time</h3>
            <div className="flex items-center space-x-2">
              <button className="p-1 text-gray-400 hover:text-gray-600">
                <Eye className="h-4 w-4" />
              </button>
              <button className="p-1 text-gray-400 hover:text-gray-600">
                <Download className="h-4 w-4" />
              </button>
            </div>
          </div>
          <div className="h-64 flex items-center justify-center bg-gray-50 rounded-lg">
            <div className="text-center">
              <BarChart3 className="h-12 w-12 text-gray-400 mx-auto mb-4" />
              <p className="text-gray-500">Chart visualization coming soon...</p>
            </div>
          </div>
        </div>

        {/* Location Distribution */}
        <div className="card p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-medium text-gray-900">Location Distribution</h3>
            <div className="flex items-center space-x-2">
              <button className="p-1 text-gray-400 hover:text-gray-600">
                <Eye className="h-4 w-4" />
              </button>
              <button className="p-1 text-gray-400 hover:text-gray-600">
                <Download className="h-4 w-4" />
              </button>
            </div>
          </div>
          <div className="h-64 flex items-center justify-center bg-gray-50 rounded-lg">
            <div className="text-center">
              <TrendingUp className="h-12 w-12 text-gray-400 mx-auto mb-4" />
              <p className="text-gray-500">Chart visualization coming soon...</p>
            </div>
          </div>
        </div>
      </div>

      {/* Detailed Statistics */}
      <div className="card p-6">
        <h3 className="text-lg font-medium text-gray-900 mb-4">Detailed Statistics</h3>
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Metric
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Value
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Change
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Trend
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              <tr>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                  Total Check-ins
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {checkInStats?.stats?.total_check_ins?.toLocaleString() || '0'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-green-600">
                  +12%
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <TrendingUp className="h-4 w-4 text-green-600" />
                </td>
              </tr>
              <tr>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                  Today's Check-ins
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {checkInStats?.stats?.today_check_ins?.toLocaleString() || '0'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-green-600">
                  +5%
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <TrendingUp className="h-4 w-4 text-green-600" />
                </td>
              </tr>
              <tr>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                  Weekly Check-ins
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {checkInStats?.stats?.weekly_check_ins?.toLocaleString() || '0'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-green-600">
                  +8%
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <TrendingUp className="h-4 w-4 text-green-600" />
                </td>
              </tr>
              <tr>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                  Active Locations
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {locations.length}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  No change
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="h-4 w-4" />
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      {/* Insights */}
      <div className="card p-6">
        <h3 className="text-lg font-medium text-gray-900 mb-4">Insights</h3>
        <div className="space-y-4">
          <div className="p-4 bg-blue-50 rounded-lg">
            <div className="flex">
              <div className="flex-shrink-0">
                <TrendingUp className="h-5 w-5 text-blue-400" />
              </div>
              <div className="ml-3">
                <h4 className="text-sm font-medium text-blue-800">Growth Trend</h4>
                <p className="text-sm text-blue-700 mt-1">
                  Check-ins have increased by 12% compared to the previous period.
                </p>
              </div>
            </div>
          </div>
          <div className="p-4 bg-green-50 rounded-lg">
            <div className="flex">
              <div className="flex-shrink-0">
                <Users className="h-5 w-5 text-green-400" />
              </div>
              <div className="ml-3">
                <h4 className="text-sm font-medium text-green-800">Peak Hours</h4>
                <p className="text-sm text-green-700 mt-1">
                  Most check-ins occur between 9:00 AM and 11:00 AM on Sundays.
                </p>
              </div>
            </div>
          </div>
          <div className="p-4 bg-purple-50 rounded-lg">
            <div className="flex">
              <div className="flex-shrink-0">
                <MapPin className="h-5 w-5 text-purple-400" />
              </div>
              <div className="ml-3">
                <h4 className="text-sm font-medium text-purple-800">Popular Locations</h4>
                <p className="text-sm text-purple-700 mt-1">
                  Main Campus and Youth Room are the most active locations.
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AnalyticsPage; 