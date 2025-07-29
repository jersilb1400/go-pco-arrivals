import React, { createContext, useContext, useEffect, useState } from 'react';
import type { ReactNode } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import type { User, AuthStatusResponse } from '../types/api';
import apiService from '../services/api';

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (rememberMe?: boolean) => void;
  logout: () => Promise<void>;
  refreshAuth: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const queryClient = useQueryClient();
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  const {
    data: authStatus,
    isLoading,
    error,
    refetch,
  } = useQuery<AuthStatusResponse>({
    queryKey: ['auth', 'status'],
    queryFn: () => apiService.getAuthStatus(),
    retry: false,
    refetchOnWindowFocus: false,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });

  // Debug logging
  useEffect(() => {
    console.log('Auth status:', authStatus);
    console.log('Auth error:', error);
    console.log('Auth loading:', isLoading);
  }, [authStatus, error, isLoading]);

  useEffect(() => {
    if (authStatus) {
      setIsAuthenticated(authStatus.is_authenticated);
    }
  }, [authStatus]);

  const login = (rememberMe: boolean = false) => {
    apiService.login(rememberMe);
  };

  const logout = async () => {
    try {
      await apiService.logout();
      setIsAuthenticated(false);
      queryClient.clear();
    } catch (error) {
      console.error('Logout failed:', error);
    }
  };

  const refreshAuth = async () => {
    try {
      await refetch();
    } catch (error) {
      console.error('Auth refresh failed:', error);
    }
  };

  const value: AuthContextType = {
    user: authStatus?.user || null,
    isAuthenticated,
    isLoading,
    login,
    logout,
    refreshAuth,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}; 