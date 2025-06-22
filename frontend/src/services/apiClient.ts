// API client for backend communication
import axios, { AxiosInstance, AxiosResponse, AxiosError } from 'axios';
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
  (config) => {
    // Add request ID for tracking
    config.headers['X-Request-ID'] = `web-${Date.now()}-${Math.random().toString(36).substring(2, 9)}`;
    
    // Log request in development
    if (import.meta.env.DEV) {
      console.log(`🚀 API Request [${config.method?.toUpperCase()}] ${config.url}`, {
        data: config.data,
        params: config.params
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
  (error: AxiosError<ApiResponse>) => {
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
      details: error.response?.data?.error?.details
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