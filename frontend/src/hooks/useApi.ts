// Custom hooks for API operations
import { useState, useCallback } from 'react';
import { ApiService } from '../services/apiClient';
import type { ApiError } from '../types/api';

// Generic API hook type
interface UseApiState<T> {
  data: T | null;
  isLoading: boolean;
  error: ApiError | null;
}

interface UseApiReturn<T> extends UseApiState<T> {
  execute: (...args: any[]) => Promise<T>;
  reset: () => void;
}

// Generic API hook
export function useApi<T>(
  apiFunction: (...args: any[]) => Promise<T>
): UseApiReturn<T> {
  const [state, setState] = useState<UseApiState<T>>({
    data: null,
    isLoading: false,
    error: null
  });

  const execute = useCallback(async (...args: any[]): Promise<T> => {
    setState(prev => ({
      ...prev,
      isLoading: true,
      error: null
    }));

    try {
      const result = await apiFunction(...args);
      setState({
        data: result,
        isLoading: false,
        error: null
      });
      return result;
    } catch (error) {
      const apiError = error as ApiError;
      setState(prev => ({
        ...prev,
        isLoading: false,
        error: apiError
      }));
      throw error;
    }
  }, [apiFunction]);

  const reset = useCallback(() => {
    setState({
      data: null,
      isLoading: false,
      error: null
    });
  }, []);

  return {
    ...state,
    execute,
    reset
  };
}

// Specific API hooks
export const useHealthCheck = () => {
  return useApi(ApiService.healthCheck);
};

export const usePing = () => {
  return useApi(ApiService.ping);
};

export const useCreateUser = () => {
  return useApi(ApiService.createUser);
};

export const useValidateUser = () => {
  return useApi(ApiService.validateUser);
};

export const useCreateSession = () => {
  return useApi(ApiService.createSession);
};

export const useGetSession = () => {
  return useApi(ApiService.getSession);
};

export const useUpdateSession = () => {
  return useApi(ApiService.updateSession);
};

export const useDeleteSession = () => {
  return useApi(ApiService.deleteSession);
};

export const useSearchAddress = () => {
  return useApi(ApiService.searchAddress);
};

export const useGetPrefectures = () => {
  return useApi(ApiService.getPrefectures);
};

export const useGetOptions = () => {
  return useApi(ApiService.getOptions);
};

export const useGetPlans = () => {
  return useApi(ApiService.getPlans);
};

export const useCheckInventory = () => {
  return useApi(ApiService.checkInventory);
};

export const useCheckRegion = () => {
  return useApi(ApiService.checkRegion);
};

// Combined hook for address search with postal code formatting
export const useAddressSearch = () => {
  const { execute, ...rest } = useSearchAddress();

  const searchByPostalCode = useCallback(async (postalCode1: string, postalCode2: string) => {
    const fullPostalCode = `${postalCode1}${postalCode2}`;
    return execute(fullPostalCode);
  }, [execute]);

  return {
    ...rest,
    searchByPostalCode
  };
};

// Combined hook for option availability check
export const useOptionAvailability = () => {
  const inventoryHook = useCheckInventory();
  const regionHook = useCheckRegion();

  const checkAvailability = useCallback(async (
    optionTypes: string[],
    prefecture: string,
    city: string
  ) => {
    try {
      const [inventoryResult, regionResult] = await Promise.all([
        inventoryHook.execute(optionTypes),
        regionHook.execute(prefecture, city, optionTypes)
      ]);

      return {
        inventory: inventoryResult.inventory,
        restrictions: regionResult.restrictions
      };
    } catch (error) {
      throw error;
    }
  }, [inventoryHook.execute, regionHook.execute]);

  return {
    data: {
      inventory: inventoryHook.data?.inventory || {},
      restrictions: regionHook.data?.restrictions || {}
    },
    isLoading: inventoryHook.isLoading || regionHook.isLoading,
    error: inventoryHook.error || regionHook.error,
    checkAvailability,
    reset: () => {
      inventoryHook.reset();
      regionHook.reset();
    }
  };
};