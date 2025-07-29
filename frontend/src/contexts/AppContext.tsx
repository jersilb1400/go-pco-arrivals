import React, { createContext, useContext, useEffect, useState } from 'react';
import type { ReactNode } from 'react';
import type { Event } from '../types/api';

interface AppState {
  selectedEvent: Event | null;
  selectedDate: string;
  selectedLocationId: string | null;
  billboardActive: boolean;
}

interface AppContextType {
  state: AppState;
  setSelectedEvent: (event: Event | null) => void;
  setSelectedDate: (date: string) => void;
  setSelectedLocationId: (locationId: string | null) => void;
  setBillboardActive: (active: boolean) => void;
  clearSelection: () => void;
}

const AppContext = createContext<AppContextType | undefined>(undefined);

export const useApp = () => {
  const context = useContext(AppContext);
  if (context === undefined) {
    throw new Error('useApp must be used within an AppProvider');
  }
  return context;
};

interface AppProviderProps {
  children: ReactNode;
}

// Storage keys for persistence
const STORAGE_KEYS = {
  SELECTED_EVENT: 'pco_billboard_selected_event',
  SELECTED_DATE: 'pco_billboard_selected_date',
  SELECTED_LOCATION: 'pco_billboard_selected_location',
  BILLBOARD_ACTIVE: 'pco_billboard_active',
} as const;

// Helper functions for localStorage
const getStoredValue = <T,>(key: string, defaultValue: T): T => {
  try {
    const item = localStorage.getItem(key);
    return item ? JSON.parse(item) : defaultValue;
  } catch (error) {
    console.warn(`Failed to parse stored value for ${key}:`, error);
    return defaultValue;
  }
};

const setStoredValue = <T,>(key: string, value: T): void => {
  try {
    localStorage.setItem(key, JSON.stringify(value));
  } catch (error) {
    console.warn(`Failed to store value for ${key}:`, error);
  }
};

export const AppProvider: React.FC<AppProviderProps> = ({ children }) => {
  // Initialize state from localStorage
  const [state, setState] = useState<AppState>(() => ({
    selectedEvent: getStoredValue(STORAGE_KEYS.SELECTED_EVENT, null),
    selectedDate: getStoredValue(STORAGE_KEYS.SELECTED_DATE, new Date().toISOString().split('T')[0]),
    selectedLocationId: getStoredValue(STORAGE_KEYS.SELECTED_LOCATION, null),
    billboardActive: getStoredValue(STORAGE_KEYS.BILLBOARD_ACTIVE, false),
  }));

  // Update localStorage when state changes
  useEffect(() => {
    setStoredValue(STORAGE_KEYS.SELECTED_EVENT, state.selectedEvent);
  }, [state.selectedEvent]);

  useEffect(() => {
    setStoredValue(STORAGE_KEYS.SELECTED_DATE, state.selectedDate);
  }, [state.selectedDate]);

  useEffect(() => {
    setStoredValue(STORAGE_KEYS.SELECTED_LOCATION, state.selectedLocationId);
  }, [state.selectedLocationId]);

  useEffect(() => {
    setStoredValue(STORAGE_KEYS.BILLBOARD_ACTIVE, state.billboardActive);
  }, [state.billboardActive]);

  const setSelectedEvent = (event: Event | null) => {
    setState(prev => ({ ...prev, selectedEvent: event }));
  };

  const setSelectedDate = (date: string) => {
    setState(prev => ({ ...prev, selectedDate: date }));
  };

  const setSelectedLocationId = (locationId: string | null) => {
    setState(prev => ({ ...prev, selectedLocationId: locationId }));
  };

  const setBillboardActive = (active: boolean) => {
    setState(prev => ({ ...prev, billboardActive: active }));
  };

  const clearSelection = () => {
    setState(prev => ({
      ...prev,
      selectedEvent: null,
      selectedLocationId: null,
    }));
  };

  const value: AppContextType = {
    state,
    setSelectedEvent,
    setSelectedDate,
    setSelectedLocationId,
    setBillboardActive,
    clearSelection,
  };

  return <AppContext.Provider value={value}>{children}</AppContext.Provider>;
}; 