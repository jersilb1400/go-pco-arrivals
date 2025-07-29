// API Response Types
export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}

// User Types
export interface User {
  id: number;
  pco_user_id: string;
  name: string;
  email: string;
  avatar?: string;
  is_admin: boolean;
  is_active: boolean;
  last_login: string;
  last_activity: string;
  created_at: string;
}

// Session Types
export interface SessionData {
  user_id: number;
  pco_user_id: string;
  email: string;
  is_admin: boolean;
  is_remember_me: boolean;
  expires_at: string;
}

// Auth Types
export interface AuthStatusResponse {
  is_authenticated: boolean;
  user?: User;
  session?: SessionData;
  expires_at?: string;
}

export interface LoginRequest {
  remember_me: boolean;
}

export interface LogoutResponse {
  success: boolean;
  message: string;
}

// Check-in Types
export interface CheckInDisplay {
  id: string;
  person_name: string;
  check_in_time: string;
  location_name: string;
  notes?: string;
  time_ago: string;
}

// Billboard Types
export interface BillboardState {
  location_id: string;
  location_name: string;
  last_updated: string;
  total_check_ins: number;
  recent_check_ins: CheckInDisplay[];
  is_online: boolean;
}

export interface BillboardStateResponse {
  success: boolean;
  state?: BillboardState;
  error?: string;
}

export interface CheckInsResponse {
  success: boolean;
  check_ins?: CheckInDisplay[];
  total?: number;
  error?: string;
}

export interface StatsResponse {
  success: boolean;
  stats?: {
    total_check_ins: number;
    today_check_ins: number;
    weekly_check_ins: number;
    period_days: number;
    location_id: string;
  };
  error?: string;
}

export interface SyncResponse {
  success: boolean;
  message: string;
  error?: string;
}

// Location Types
export interface Location {
  id: number;
  pco_location_id: string;
  name: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface LocationsResponse {
  success: boolean;
  locations?: Location[];
  error?: string;
}

// System Status Types
export interface SystemStatus {
  success: boolean;
  system: {
    locations: {
      total: number;
    };
    check_ins: {
      today: number;
    };
    sessions: {
      active: number;
    };
    timestamp: string;
  };
}

// Real-time Update Types
export interface RealTimeUpdate {
  type: string;
  location_id: string;
  check_in?: CheckInDisplay;
  state?: BillboardState;
  timestamp: string;
}

// WebSocket Types
export interface WebSocketMessage {
  type: string;
  data: any;
  timestamp: string;
}

export interface WebSocketConnectionStatus {
  isConnected: boolean;
  isConnecting: boolean;
  error?: string;
} 

// Admin Types
export interface Event {
  id: string;
  name: string;
  date: string;
  location: string;
  time?: string;
  description?: string;
  is_active: boolean;
  created_at: string;
}

export interface Notification {
  id: string;
  message: string;
  type: string;
  created_at: string;
  child_name?: string;
  security_code?: string;
  location_name?: string;
  event_name?: string;
  parent_name?: string;
  status: string;
  expires_at?: string;
}

export interface SecurityCode {
  code: string;
  created_at: string;
  created_by: string;
}

export interface BillboardControl {
  event_id: string;
  event_name: string;
  location_id: string;
  location_name: string;
  security_codes: string[];
  is_active: boolean;
  last_updated: string;
}

export interface EventsResponse {
  success: boolean;
  events?: Event[];
  error?: string;
}

export interface NotificationsResponse {
  success: boolean;
  notifications?: Notification[];
  error?: string;
}

export interface SecurityCodesResponse {
  success: boolean;
  codes?: SecurityCode[];
  error?: string;
}

export interface BillboardControlResponse {
  success: boolean;
  control?: BillboardControl;
  error?: string;
} 

// Advanced Location Status Types
export interface LocationStatus {
  id: string;
  name: string;
  active_children: number;
  total_check_ins: number;
  today_check_ins: number;
  notifications: Notification[];
  recent_check_ins: CheckInDisplay[];
  last_updated: string;
}

export interface LocationStatusResponse {
  success: boolean;
  location?: LocationStatus;
  error?: string;
}

export interface LocationAnalytics {
  location_id: string;
  period_days: number;
  daily_stats: Array<{
    date: string;
    count: number;
  }>;
  peak_hours: Array<{
    hour: number;
    count: number;
  }>;
  avg_wait_time_mins: number;
  efficiency_rate: number;
  total_notifications: number;
  completed_pickups: number;
  expired_notifications: number;
  generated_at: string;
}

export interface LocationAnalyticsResponse {
  success: boolean;
  analytics?: LocationAnalytics;
  error?: string;
}

export interface LocationOverview {
  id: string;
  name: string;
  active_children: number;
  today_check_ins: number;
  total_check_ins: number;
  is_active: boolean;
  last_updated: string;
}

export interface LocationsOverviewResponse {
  success: boolean;
  locations?: LocationOverview[];
  summary?: {
    total_locations: number;
    active_locations: number;
    total_children: number;
    generated_at: string;
  };
  error?: string;
}

// Enhanced Location Status with Real-time Features
export interface EnhancedLocationStatus extends LocationStatus {
  efficiency_score: number;
  wait_time_trend: 'improving' | 'stable' | 'declining';
  capacity_utilization: number;
  staff_alert_level: 'low' | 'medium' | 'high';
  last_activity: string;
  is_online: boolean;
}

export interface LocationPerformanceMetrics {
  location_id: string;
  avg_pickup_time: number;
  completion_rate: number;
  peak_hour_efficiency: number;
  staff_response_time: number;
  parent_satisfaction_score: number;
  last_calculated: string;
}

export interface LocationAlert {
  id: string;
  location_id: string;
  type: 'capacity' | 'wait_time' | 'staff_shortage' | 'system_issue';
  severity: 'low' | 'medium' | 'high' | 'critical';
  message: string;
  created_at: string;
  resolved_at?: string;
  is_active: boolean;
}

export interface LocationAlertsResponse {
  success: boolean;
  alerts?: LocationAlert[];
  error?: string;
} 