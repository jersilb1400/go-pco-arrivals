import type {
  AuthStatusResponse,
  BillboardStateResponse,
  CheckInsResponse,
  StatsResponse,
  SyncResponse,
  LocationsResponse,
  SystemStatus,
  User,
  Event,
  Notification,
  SecurityCode,
  BillboardControl,
  EventsResponse,
  NotificationsResponse,
  SecurityCodesResponse,
  BillboardControlResponse,
  LocationStatusResponse,
  LocationAnalyticsResponse,
  LocationsOverviewResponse,
} from '../types/api';

const API_BASE_URL = '';

class ApiService {
  private baseURL: string;

  constructor(baseURL: string = API_BASE_URL) {
    this.baseURL = baseURL;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    
    const config: RequestInit = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      credentials: 'include', // Include cookies for session management
      ...options,
    };

    try {
      const response = await fetch(url, config);
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      
      return await response.json();
    } catch (error) {
      console.error('API request failed:', error);
      throw error;
    }
  }

  // Auth endpoints
  async getAuthStatus(): Promise<AuthStatusResponse> {
    return this.request<AuthStatusResponse>('/auth/status');
  }

  async login(rememberMe: boolean = false): Promise<void> {
    const url = `/auth/login?remember_me=${rememberMe}`;
    window.location.href = `${this.baseURL}${url}`;
  }

  async logout(): Promise<void> {
    await this.request('/auth/logout', { method: 'POST' });
  }

  async refreshToken(): Promise<void> {
    await this.request('/auth/refresh', { method: 'POST' });
  }

  async getUserProfile(): Promise<User> {
    return this.request<User>('/auth/profile');
  }

  async updateUserProfile(profile: Partial<User>): Promise<User> {
    return this.request<User>('/auth/profile', {
      method: 'PUT',
      body: JSON.stringify(profile),
    });
  }

  // Billboard endpoints
  async getBillboardState(locationId: string): Promise<BillboardStateResponse> {
    return this.request<BillboardStateResponse>(`/billboard/state/${locationId}`);
  }

  async getRecentCheckIns(
    locationId: string,
    limit: number = 10
  ): Promise<CheckInsResponse> {
    return this.request<CheckInsResponse>(
      `/billboard/check-ins/${locationId}?limit=${limit}`
    );
  }

  async getCheckInStats(
    locationId: string,
    days: number = 7
  ): Promise<StatsResponse> {
    return this.request<StatsResponse>(
      `/api/billboard/stats/${locationId}?days=${days}`
    );
  }

  async syncPCOCheckIns(locationId: string): Promise<SyncResponse> {
    return this.request<SyncResponse>(`/billboard/sync/${locationId}`, {
      method: 'POST',
    });
  }

  async getLocations(): Promise<LocationsResponse> {
    return this.request<LocationsResponse>('/billboard/locations');
  }

  async addLocation(pcoLocationId: string, name: string): Promise<void> {
    await this.request('/billboard/locations', {
      method: 'POST',
      body: JSON.stringify({ pco_location_id: pcoLocationId, name }),
    });
  }

  async getLocationBillboard(locationId: string): Promise<BillboardStateResponse> {
    return this.request<BillboardStateResponse>(`/billboard/location/${locationId}`);
  }

  async cleanupOldData(): Promise<void> {
    await this.request('/billboard/cleanup', { method: 'POST' });
  }

  async getSystemStatus(): Promise<SystemStatus> {
    return this.request<SystemStatus>('/billboard/status');
  }

  // PCO API endpoints
  async getPCOLocations(): Promise<{ locations: any[] }> {
    return this.request<{ locations: any[] }>('/api/locations');
  }

  async getPCOCheckIns(locationId?: string, since?: string): Promise<{ check_ins: any[] }> {
    let url = '/api/check-ins';
    const params = new URLSearchParams();
    if (locationId) params.append('location_id', locationId);
    if (since) params.append('since', since);
    if (params.toString()) url += `?${params.toString()}`;
    return this.request<{ check_ins: any[] }>(url);
  }

  async getPCOCheckInsByLocation(locationId: string, since?: string): Promise<{ check_ins: any[], location_id: string }> {
    let url = `/api/check-ins/location/${locationId}`;
    if (since) url += `?since=${since}`;
    return this.request<{ check_ins: any[], location_id: string }>(url);
  }

  // Health endpoints
  async getHealth(): Promise<{ status: string }> {
    return this.request<{ status: string }>('/health');
  }

  async getDetailedHealth(): Promise<any> {
    return this.request('/health/detailed');
  }

  // Admin endpoints
  async getEventsByDate(date: string): Promise<EventsResponse> {
    return this.request<EventsResponse>(`/api/events?date=${date}`);
  }

  async getActiveNotifications(): Promise<NotificationsResponse> {
    return this.request<NotificationsResponse>('/api/notifications/active');
  }

  async getSecurityCodes(): Promise<SecurityCodesResponse> {
    return this.request<SecurityCodesResponse>('/api/security-codes');
  }

  async addSecurityCode(code: string): Promise<{ success: boolean; message: string }> {
    return this.request<{ success: boolean; message: string }>('/api/security-codes', {
      method: 'POST',
      body: JSON.stringify({ code }),
    });
  }

  async removeSecurityCode(code: string): Promise<{ success: boolean; message: string }> {
    return this.request<{ success: boolean; message: string }>(`/api/security-codes/${encodeURIComponent(code)}`, {
      method: 'DELETE',
    });
  }

  async getBillboardControl(): Promise<BillboardControlResponse> {
    return this.request<BillboardControlResponse>('/api/billboard/control');
  }

  async launchBillboard(eventId: string, locationId: string, securityCodes: string[]): Promise<{ success: boolean; message: string }> {
    return this.request<{ success: boolean; message: string }>('/api/billboard/launch', {
      method: 'POST',
      body: JSON.stringify({
        event_id: eventId,
        location_id: locationId,
        security_codes: securityCodes,
      }),
    });
  }

  async clearBillboard(): Promise<{ success: boolean; message: string }> {
    return this.request<{ success: boolean; message: string }>('/api/billboard/clear', {
      method: 'POST',
    });
  }

  // Advanced Location Status API endpoints
  async getLocationStatus(locationId: string): Promise<LocationStatusResponse> {
    return this.request<LocationStatusResponse>(`/api/locations/${locationId}/status`);
  }

  async getLocationAnalytics(locationId: string, days: number = 30): Promise<LocationAnalyticsResponse> {
    return this.request<LocationAnalyticsResponse>(`/api/locations/${locationId}/analytics?days=${days}`);
  }

  async getLocationsOverview(): Promise<LocationsOverviewResponse> {
    return this.request<LocationsOverviewResponse>('/api/locations/overview');
  }

  async getLocationCheckInStats(locationId: string, days: number = 7): Promise<StatsResponse> {
    return this.request<StatsResponse>(`/api/billboard/stats/${locationId}?days=${days}`);
  }
}

// Create the API service instance
const apiService = new ApiService();

export { apiService };
export default apiService; 