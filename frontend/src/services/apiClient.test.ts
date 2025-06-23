import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import axios from 'axios';
import type { ApiResponse, ApiError } from '../types';
import * as apiClientModule from './apiClient';

// Mock axios
vi.mock('axios');
const mockedAxios = vi.mocked(axios);

// Mock securityService
vi.mock('./securityService', () => ({
  securityService: {
    getCSRFToken: vi.fn().mockResolvedValue('mock-csrf-token'),
    clearCSRFToken: vi.fn(),
  },
}));

describe('apiClient', () => {
  const mockAxiosInstance = {
    request: vi.fn(),
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
    patch: vi.fn(),
    defaults: {},
    interceptors: {
      request: {
        use: vi.fn(),
      },
      response: {
        use: vi.fn(),
      },
    },
  };

  beforeEach(() => {
    vi.clearAllMocks();
    mockedAxios.create.mockReturnValue(mockAxiosInstance as any);
    
    // Clear console mocks
    vi.mocked(console.log).mockClear();
    vi.mocked(console.error).mockClear();
  });

  afterEach(() => {
    vi.resetModules();
  });

  describe('Instance Creation', () => {
    it('should create axios instance with correct base configuration', () => {
      // Re-import to trigger module initialization
      require('./apiClient');

      expect(mockedAxios.create).toHaveBeenCalledWith({
        baseURL: 'http://localhost:8080',
        timeout: 10000,
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        },
      });
    });

    it('should setup request and response interceptors', () => {
      require('./apiClient');

      expect(mockAxiosInstance.interceptors.request.use).toHaveBeenCalled();
      expect(mockAxiosInstance.interceptors.response.use).toHaveBeenCalled();
    });
  });

  describe('Request Interceptor', () => {
    it('should add request ID to headers', async () => {
      const config = {
        method: 'get',
        url: '/test',
        headers: {},
      };

      // Get the request interceptor function
      const requestInterceptor = mockAxiosInstance.interceptors.request.use.mock.calls[0][0];
      const modifiedConfig = await requestInterceptor(config);

      expect(modifiedConfig.headers['X-Request-ID']).toMatch(/^web-\d+-[a-z0-9]+$/);
    });

    it('should add CSRF token for non-GET requests', async () => {
      const config = {
        method: 'post',
        url: '/test',
        headers: {},
      };

      const requestInterceptor = mockAxiosInstance.interceptors.request.use.mock.calls[0][0];
      const modifiedConfig = await requestInterceptor(config);

      expect(modifiedConfig.headers['X-CSRF-Token']).toBe('mock-csrf-token');
      expect(modifiedConfig.headers['X-Requested-With']).toBe('XMLHttpRequest');
    });

    it('should not add CSRF token for GET requests', async () => {
      const config = {
        method: 'get',
        url: '/test',
        headers: {},
      };

      const requestInterceptor = mockAxiosInstance.interceptors.request.use.mock.calls[0][0];
      const modifiedConfig = await requestInterceptor(config);

      expect(modifiedConfig.headers['X-CSRF-Token']).toBeUndefined();
    });

    it('should log request in development mode', async () => {
      // Mock development environment
      vi.stubGlobal('import.meta', { env: { DEV: true } });

      const config = {
        method: 'post',
        url: '/test',
        data: { test: 'data' },
        headers: {},
      };

      const requestInterceptor = mockAxiosInstance.interceptors.request.use.mock.calls[0][0];
      await requestInterceptor(config);

      expect(console.log).toHaveBeenCalledWith(
        expect.stringContaining('🚀 API Request [POST] /test'),
        expect.any(Object)
      );
    });
  });

  describe('Response Interceptor', () => {
    it('should log successful response in development mode', () => {
      vi.stubGlobal('import.meta', { env: { DEV: true } });

      const response = {
        status: 200,
        config: { method: 'get', url: '/test' },
        data: { success: true },
      };

      const successInterceptor = mockAxiosInstance.interceptors.response.use.mock.calls[0][0];
      const result = successInterceptor(response);

      expect(console.log).toHaveBeenCalledWith(
        expect.stringContaining('✅ API Response [200] /test'),
        expect.any(Object)
      );
      expect(result).toBe(response);
    });

    it('should handle CSRF token errors with retry', async () => {
      const error = {
        response: {
          status: 403,
          data: {
            error: {
              code: 'CSRF_TOKEN_INVALID',
            },
          },
        },
        config: {
          headers: {},
          _retry: undefined,
        },
      };

      const errorInterceptor = mockAxiosInstance.interceptors.response.use.mock.calls[0][1];
      
      // Mock the retry request
      mockAxiosInstance.request.mockResolvedValue({ data: 'retry success' });

      const result = await errorInterceptor(error);

      expect(result.data).toBe('retry success');
      expect(error.config._retry).toBe(true);
      expect(error.config.headers['X-CSRF-Token']).toBe('mock-csrf-token');
    });

    it('should handle rate limiting errors', async () => {
      const error = {
        response: {
          status: 429,
          headers: {
            'retry-after': '60',
          },
          data: {
            error: {
              code: 'RATE_LIMIT_EXCEEDED',
              message: 'Too many requests',
            },
          },
        },
        config: { url: '/test' },
        message: 'Request failed',
      };

      const errorInterceptor = mockAxiosInstance.interceptors.response.use.mock.calls[0][1];

      try {
        await errorInterceptor(error);
      } catch (apiError) {
        const err = apiError as ApiError;
        expect(err.code).toBe('RATE_LIMIT_EXCEEDED');
        expect(err.message).toBe('Too many requests');
      }

      expect(console.warn).toHaveBeenCalledWith(
        expect.stringContaining('Rate limited. Retry after 60 seconds')
      );
    });

    it('should transform network errors into ApiError format', async () => {
      const error = {
        response: undefined,
        config: { url: '/test' },
        message: 'Network Error',
      };

      const errorInterceptor = mockAxiosInstance.interceptors.response.use.mock.calls[0][1];

      try {
        await errorInterceptor(error);
      } catch (apiError) {
        const err = apiError as ApiError;
        expect(err.code).toBe('NETWORK_ERROR');
        expect(err.message).toBe('ネットワークエラーが発生しました');
      }
    });

    it('should log errors in development mode', async () => {
      vi.stubGlobal('import.meta', { env: { DEV: true } });

      const error = {
        response: {
          status: 500,
          data: {
            error: {
              code: 'INTERNAL_ERROR',
              message: 'Server error',
            },
          },
        },
        config: { url: '/test' },
        message: 'Request failed',
      };

      const errorInterceptor = mockAxiosInstance.interceptors.response.use.mock.calls[0][1];

      try {
        await errorInterceptor(error);
      } catch (apiError) {
        // Error should be caught
      }

      expect(console.error).toHaveBeenCalledWith(
        expect.stringContaining('❌ API Error [500] /test'),
        expect.any(Object)
      );
    });
  });

  describe('Configuration', () => {
    it('should use correct timeout value', () => {
      require('./apiClient');

      expect(mockedAxios.create).toHaveBeenCalledWith(
        expect.objectContaining({
          timeout: 10000,
        })
      );
    });

    it('should set correct content type headers', () => {
      require('./apiClient');

      expect(mockedAxios.create).toHaveBeenCalledWith(
        expect.objectContaining({
          headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
          },
        })
      );
    });

    it('should use localhost base URL in test environment', () => {
      require('./apiClient');

      expect(mockedAxios.create).toHaveBeenCalledWith(
        expect.objectContaining({
          baseURL: 'http://localhost:8080',
        })
      );
    });
  });
});