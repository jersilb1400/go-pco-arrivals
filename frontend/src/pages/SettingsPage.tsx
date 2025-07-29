import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { motion } from 'framer-motion';
import { 
  User, 
  Settings, 
  Shield, 
  Bell,
  Save,
  Eye,
  EyeOff,
  CheckCircle,
  AlertCircle
} from 'lucide-react';
import { useAuth } from '../contexts/AuthContext';
import apiService from '../services/api';
import type { User as UserType } from '../types/api';

interface ProfileForm {
  name: string;
  email: string;
}

interface NotificationSettings {
  emailNotifications: boolean;
  browserNotifications: boolean;
  realTimeUpdates: boolean;
}

const SettingsPage: React.FC = () => {
  const { user, refreshAuth } = useAuth();
  const queryClient = useQueryClient();
  const [showPassword, setShowPassword] = useState(false);
  const [activeTab, setActiveTab] = useState('profile');

  const [profileForm, setProfileForm] = useState<ProfileForm>({
    name: user?.name || '',
    email: user?.email || '',
  });

  const [notificationSettings, setNotificationSettings] = useState<NotificationSettings>({
    emailNotifications: true,
    browserNotifications: true,
    realTimeUpdates: true,
  });

  const { data: userProfile, isLoading } = useQuery({
    queryKey: ['user', 'profile'],
    queryFn: apiService.getUserProfile,
    enabled: !!user,
  });

  const updateProfileMutation = useMutation({
    mutationFn: (data: ProfileForm) => apiService.updateUserProfile(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user', 'profile'] });
      refreshAuth();
    },
    onError: (error) => {
      console.error('Failed to update profile:', error);
    },
  });

  const handleProfileSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    updateProfileMutation.mutate(profileForm);
  };

  const handleProfileInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setProfileForm(prev => ({ ...prev, [name]: value }));
  };

  const handleNotificationChange = (setting: keyof NotificationSettings) => {
    setNotificationSettings(prev => ({
      ...prev,
      [setting]: !prev[setting],
    }));
  };

  const tabs = [
    { id: 'profile', name: 'Profile', icon: User },
    { id: 'notifications', name: 'Notifications', icon: Bell },
    { id: 'security', name: 'Security', icon: Shield },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Settings</h1>
        <p className="mt-1 text-sm text-gray-500">
          Manage your account settings and preferences
        </p>
      </div>

      {/* Tabs */}
      <div className="border-b border-gray-200">
        <nav className="-mb-px flex space-x-8">
          {tabs.map((tab) => {
            const Icon = tab.icon;
            return (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`py-2 px-1 border-b-2 font-medium text-sm flex items-center space-x-2 ${
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

      {/* Tab Content */}
      <div className="space-y-6">
        {/* Profile Tab */}
        {activeTab === 'profile' && (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="space-y-6"
          >
            <div className="card p-6">
              <h2 className="text-lg font-medium text-gray-900 mb-6">Profile Information</h2>
              
              {isLoading ? (
                <div className="text-center py-8">
                  <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 mx-auto"></div>
                  <p className="mt-2 text-sm text-gray-500">Loading profile...</p>
                </div>
              ) : (
                <form onSubmit={handleProfileSubmit} className="space-y-6">
                  <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
                    <div>
                      <label htmlFor="name" className="block text-sm font-medium text-gray-700">
                        Full Name
                      </label>
                      <input
                        type="text"
                        id="name"
                        name="name"
                        value={profileForm.name}
                        onChange={handleProfileInputChange}
                        className="input mt-1"
                        required
                      />
                    </div>

                    <div>
                      <label htmlFor="email" className="block text-sm font-medium text-gray-700">
                        Email Address
                      </label>
                      <input
                        type="email"
                        id="email"
                        name="email"
                        value={profileForm.email}
                        onChange={handleProfileInputChange}
                        className="input mt-1"
                        required
                      />
                    </div>
                  </div>

                                     <div className="flex items-center justify-between">
                     <div className="text-sm text-gray-500">
                       Last updated: {userProfile?.last_activity ? 
                         new Date(userProfile.last_activity).toLocaleDateString() : 
                         'Never'
                       }
                     </div>
                    <button
                      type="submit"
                      disabled={updateProfileMutation.isPending}
                      className="btn-primary"
                    >
                      <Save className="h-4 w-4 mr-2" />
                      {updateProfileMutation.isPending ? 'Saving...' : 'Save Changes'}
                    </button>
                  </div>
                </form>
              )}
            </div>

            {/* Account Info */}
            <div className="card p-6">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Account Information</h3>
              <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                <div>
                  <p className="text-sm font-medium text-gray-500">PCO User ID</p>
                  <p className="text-sm text-gray-900">{userProfile?.pco_user_id || 'N/A'}</p>
                </div>
                <div>
                  <p className="text-sm font-medium text-gray-500">Account Type</p>
                  <p className="text-sm text-gray-900">
                    {userProfile?.is_admin ? 'Administrator' : 'User'}
                  </p>
                </div>
                <div>
                  <p className="text-sm font-medium text-gray-500">Member Since</p>
                  <p className="text-sm text-gray-900">
                    {userProfile?.created_at ? 
                      new Date(userProfile.created_at).toLocaleDateString() : 
                      'N/A'
                    }
                  </p>
                </div>
                <div>
                  <p className="text-sm font-medium text-gray-500">Last Login</p>
                  <p className="text-sm text-gray-900">
                    {userProfile?.last_login ? 
                      new Date(userProfile.last_login).toLocaleString() : 
                      'N/A'
                    }
                  </p>
                </div>
              </div>
            </div>
          </motion.div>
        )}

        {/* Notifications Tab */}
        {activeTab === 'notifications' && (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="space-y-6"
          >
            <div className="card p-6">
              <h2 className="text-lg font-medium text-gray-900 mb-6">Notification Preferences</h2>
              
              <div className="space-y-6">
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="text-sm font-medium text-gray-900">Email Notifications</h3>
                    <p className="text-sm text-gray-500">
                      Receive email updates about system events
                    </p>
                  </div>
                  <button
                    onClick={() => handleNotificationChange('emailNotifications')}
                    className={`relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 ${
                      notificationSettings.emailNotifications ? 'bg-primary-600' : 'bg-gray-200'
                    }`}
                  >
                    <span
                      className={`pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out ${
                        notificationSettings.emailNotifications ? 'translate-x-5' : 'translate-x-0'
                      }`}
                    />
                  </button>
                </div>

                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="text-sm font-medium text-gray-900">Browser Notifications</h3>
                    <p className="text-sm text-gray-500">
                      Show desktop notifications for real-time updates
                    </p>
                  </div>
                  <button
                    onClick={() => handleNotificationChange('browserNotifications')}
                    className={`relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 ${
                      notificationSettings.browserNotifications ? 'bg-primary-600' : 'bg-gray-200'
                    }`}
                  >
                    <span
                      className={`pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out ${
                        notificationSettings.browserNotifications ? 'translate-x-5' : 'translate-x-0'
                      }`}
                    />
                  </button>
                </div>

                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="text-sm font-medium text-gray-900">Real-time Updates</h3>
                    <p className="text-sm text-gray-500">
                      Enable WebSocket connections for live updates
                    </p>
                  </div>
                  <button
                    onClick={() => handleNotificationChange('realTimeUpdates')}
                    className={`relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 ${
                      notificationSettings.realTimeUpdates ? 'bg-primary-600' : 'bg-gray-200'
                    }`}
                  >
                    <span
                      className={`pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out ${
                        notificationSettings.realTimeUpdates ? 'translate-x-5' : 'translate-x-0'
                      }`}
                    />
                  </button>
                </div>
              </div>
            </div>
          </motion.div>
        )}

        {/* Security Tab */}
        {activeTab === 'security' && (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="space-y-6"
          >
            <div className="card p-6">
              <h2 className="text-lg font-medium text-gray-900 mb-6">Security Settings</h2>
              
              <div className="space-y-6">
                <div>
                  <h3 className="text-sm font-medium text-gray-900 mb-2">Session Management</h3>
                  <p className="text-sm text-gray-500 mb-4">
                    Your current session is managed through Planning Center Online OAuth.
                  </p>
                  <div className="bg-gray-50 p-4 rounded-lg">
                    <div className="flex items-center space-x-2">
                      <CheckCircle className="h-5 w-5 text-green-600" />
                      <span className="text-sm text-gray-700">OAuth authentication active</span>
                    </div>
                  </div>
                </div>

                <div>
                  <h3 className="text-sm font-medium text-gray-900 mb-2">Account Security</h3>
                  <p className="text-sm text-gray-500 mb-4">
                    Your account security is managed by Planning Center Online.
                  </p>
                  <div className="space-y-3">
                    <div className="flex items-center space-x-2">
                      <CheckCircle className="h-4 w-4 text-green-600" />
                      <span className="text-sm text-gray-700">Two-factor authentication (if enabled in PCO)</span>
                    </div>
                    <div className="flex items-center space-x-2">
                      <CheckCircle className="h-4 w-4 text-green-600" />
                      <span className="text-sm text-gray-700">Secure OAuth 2.0 authentication</span>
                    </div>
                    <div className="flex items-center space-x-2">
                      <CheckCircle className="h-4 w-4 text-green-600" />
                      <span className="text-sm text-gray-700">Automatic session refresh</span>
                    </div>
                  </div>
                </div>

                <div>
                  <h3 className="text-sm font-medium text-gray-900 mb-2">Data Privacy</h3>
                  <p className="text-sm text-gray-500 mb-4">
                    This application only stores check-in data and user session information.
                  </p>
                  <div className="bg-blue-50 p-4 rounded-lg">
                    <div className="flex items-start space-x-2">
                      <AlertCircle className="h-5 w-5 text-blue-600 mt-0.5" />
                      <div>
                        <p className="text-sm text-blue-800">
                          Your personal data is managed by Planning Center Online. 
                          This application only processes check-in information for display purposes.
                        </p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </motion.div>
        )}
      </div>
    </div>
  );
};

export default SettingsPage; 