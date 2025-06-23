// API client for backend communication
import axios, { AxiosInstance, AxiosResponse, AxiosError } from 'axios';
import { securityService } from './securityService';
import type {
  ApiResponse,
  ApiError,
  UserCreateRequest,
  UserValidateRequest,
  UserCreateResponse,
  UserValidateResponse,
  SessionCreateRequest,
  SessionCreateResponse,
  SessionGetResponse,
  AddressSearchRequest,  
  AddressSearchResponse,
  PrefecturesGetResponse,
  OptionsGetResponse,
  PlansGetResponse,
  InventoryCheckRequest,
  InventoryCheckResponse,
  RegionCheckRequest,
  RegionCheckResponse,
  HealthCheckResponse
} from '../types/api';

// API client configuration
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
const API_TIMEOUT = 30000; // 30 seconds

// Create axios instance
const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: API_TIMEOUT,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor
apiClient.interceptors.request.use(
  async (config) => {
    // Add request ID for tracking
    config.headers['X-Request-ID'] = `web-${Date.now()}-${Math.random().toString(36).substring(2, 9)}`;
    
    // Add CSRF token for non-GET requests
    if (config.method && !['get', 'head', 'options'].includes(config.method.toLowerCase())) {
      try {
        const csrfToken = await securityService.getCSRFToken();
        config.headers['X-CSRF-Token'] = csrfToken;
      } catch (error) {
        console.error('Failed to get CSRF token:', error);
        // Don't block the request, let the server handle the missing token
      }
    }
    
    // Add security headers
    config.headers['X-Requested-With'] = 'XMLHttpRequest';
    
    // Log request in development
    if (import.meta.env.DEV) {
      console.log(`🚀 API Request [${config.method?.toUpperCase()}] ${config.url}`, {
        data: config.data,
        params: config.params,
        headers: {
          'X-CSRF-Token': config.headers['X-CSRF-Token'] ? '***' : undefined,
          'X-Request-ID': config.headers['X-Request-ID']
        }
      });
    }
    
    return config;
  },
  (error) => {
    console.error('❌ Request Error:', error);
    return Promise.reject(error);
  }
);

// Response interceptor
apiClient.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    // Log response in development
    if (import.meta.env.DEV) {
      console.log(`✅ API Response [${response.status}] ${response.config.url}`, response.data);
    }
    
    return response;
  },
  async (error: AxiosError<ApiResponse>) => {
    // Handle CSRF token errors
    if (error.response?.status === 403 && 
        error.response?.data?.error?.code === 'CSRF_TOKEN_INVALID') {
      console.warn('CSRF token expired, clearing cache');
      securityService.clearCSRFToken();
      
      // Retry the request once with a new token
      if (error.config && !(error.config as any)._retry) {
        (error.config as any)._retry = true;
        
        try {
          const newToken = await securityService.getCSRFToken();
          error.config.headers = error.config.headers || {};
          error.config.headers['X-CSRF-Token'] = newToken;
          return apiClient.request(error.config);
        } catch (retryError) {
          console.error('Failed to retry request with new CSRF token:', retryError);
        }
      }
    }
    
    // Handle rate limiting
    if (error.response?.status === 429) {
      const retryAfter = error.response.headers['retry-after'];
      if (retryAfter) {
        console.warn(`Rate limited. Retry after ${retryAfter} seconds`);
      }
    }
    
    // Log error in development
    if (import.meta.env.DEV) {
      console.error(`❌ API Error [${error.response?.status}] ${error.config?.url}`, {
        error: error.response?.data?.error,
        message: error.message
      });
    }
    
    // Transform error for consistent handling
    const apiError: ApiError = {
      code: error.response?.data?.error?.code || 'NETWORK_ERROR',
      message: error.response?.data?.error?.message || error.message || 'ネットワークエラーが発生しました',
      ...(error.response?.data?.error?.details && { details: error.response.data.error.details })
    };
    
    return Promise.reject(apiError);
  }
);

// API service class
export class ApiService {
  // Health check endpoints
  static async healthCheck(): Promise<HealthCheckResponse> {
    const response = await apiClient.get<ApiResponse<HealthCheckResponse>>('/health');
    if (!response.data.success || !response.data.data) {
      throw new Error('Health check failed');
    }
    return response.data.data;
  }

  static async ping(): Promise<{ message: string }> {
    const response = await apiClient.get<ApiResponse<{ message: string }>>('/api/v1/ping');
    if (!response.data.success || !response.data.data) {
      throw new Error('Ping failed');
    }
    return response.data.data;
  }

  // User endpoints
  static async createUser(userData: UserCreateRequest): Promise<UserCreateResponse> {
    const response = await apiClient.post<ApiResponse<UserCreateResponse>>('/api/v1/users', userData);
    if (!response.data.success || !response.data.data) {
      throw response.data.error || new Error('User creation failed');
    }
    return response.data.data;
  }

  static async validateUser(userData: UserValidateRequest): Promise<UserValidateResponse> {
    const response = await apiClient.post<ApiResponse<UserValidateResponse>>('/api/v1/users/validate', userData);
    if (!response.data.success || !response.data.data) {
      throw response.data.error || new Error('User validation failed');
    }
    return response.data.data;
  }

  // Session management endpoints
  static async createSession(sessionData: SessionCreateRequest): Promise<SessionCreateResponse> {
    const response = await apiClient.post<ApiResponse<SessionCreateResponse>>('/api/v1/sessions', sessionData);
    if (!response.data.success || !response.data.data) {
      throw response.data.error || new Error('Session creation failed');
    }
    return response.data.data;
  }

  static async getSession(sessionId: string): Promise<SessionGetResponse> {
    const response = await apiClient.get<ApiResponse<SessionGetResponse>>(`/api/v1/sessions/${sessionId}`);
    if (!response.data.success || !response.data.data) {
      throw response.data.error || new Error('Session retrieval failed');
    }
    return response.data.data;
  }

  static async updateSession(sessionId: string, sessionData: SessionCreateRequest): Promise<void> {
    const response = await apiClient.put<ApiResponse<void>>(`/api/v1/sessions/${sessionId}`, sessionData);
    if (!response.data.success) {
      throw response.data.error || new Error('Session update failed');
    }
  }

  static async deleteSession(sessionId: string): Promise<void> {
    const response = await apiClient.delete<ApiResponse<void>>(`/api/v1/sessions/${sessionId}`);
    if (!response.data.success) {
      throw response.data.error || new Error('Session deletion failed');
    }
  }

  // Address and prefecture endpoints
  static async searchAddress(postalCode: string): Promise<AddressSearchResponse> {
    const response = await apiClient.get<ApiResponse<AddressSearchResponse>>('/api/v1/address/search', {
      params: { postal_code: postalCode }
    });
    if (!response.data.success || !response.data.data) {
      throw response.data.error || new Error('Address search failed');
    }
    return response.data.data;
  }

  static async getPrefectures(): Promise<PrefecturesGetResponse> {
    const response = await apiClient.get<ApiResponse<PrefecturesGetResponse>>('/api/v1/prefectures');
    if (!response.data.success || !response.data.data) {
      throw response.data.error || new Error('Prefecture retrieval failed');
    }
    return response.data.data;
  }

  // Option and plan endpoints
  static async getOptions(): Promise<OptionsGetResponse> {
    const response = await apiClient.get<ApiResponse<OptionsGetResponse>>('/api/v1/options');
    if (!response.data.success || !response.data.data) {
      throw response.data.error || new Error('Options retrieval failed');
    }
    return response.data.data;
  }

  static async getPlans(): Promise<PlansGetResponse> {
    const response = await apiClient.get<ApiResponse<PlansGetResponse>>('/api/v1/plans');
    if (!response.data.success || !response.data.data) {
      throw response.data.error || new Error('Plans retrieval failed');
    }
    return response.data.data;
  }

  // Inventory and region check endpoints
  static async checkInventory(optionTypes: string[]): Promise<InventoryCheckResponse> {
    const response = await apiClient.post<ApiResponse<InventoryCheckResponse>>('/api/v1/options/check-inventory', {
      option_types: optionTypes
    });
    if (!response.data.success || !response.data.data) {
      throw response.data.error || new Error('Inventory check failed');
    }
    return response.data.data;
  }

  static async checkRegion(prefecture: string, city: string, optionTypes: string[]): Promise<RegionCheckResponse> {
    const response = await apiClient.post<ApiResponse<RegionCheckResponse>>('/api/v1/region/check', {
      prefecture,
      city,
      option_types: optionTypes
    });
    if (!response.data.success || !response.data.data) {
      throw response.data.error || new Error('Region check failed');
    }
    return response.data.data;
  }
}

// Export axios instance for direct use if needed
export { apiClient };

// Export default API service
export default ApiService;