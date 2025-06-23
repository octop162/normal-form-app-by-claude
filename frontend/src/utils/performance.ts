// Performance optimization utilities
import { useCallback, useMemo, useRef, useEffect } from 'react';

// Debounce hook for real-time validation
export function useDebounce<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = React.useState<T>(value);

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    return () => {
      clearTimeout(handler);
    };
  }, [value, delay]);

  return debouncedValue;
}

// Throttle hook for API calls
export function useThrottle<T extends (...args: any[]) => any>(
  callback: T,
  delay: number
): T {
  const lastCallRef = useRef<number>(0);
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  return useCallback(
    (...args: Parameters<T>) => {
      const now = Date.now();
      const timeSinceLastCall = now - lastCallRef.current;

      if (timeSinceLastCall >= delay) {
        lastCallRef.current = now;
        return callback(...args);
      } else {
        if (timeoutRef.current) {
          clearTimeout(timeoutRef.current);
        }
        
        timeoutRef.current = setTimeout(() => {
          lastCallRef.current = Date.now();
          callback(...args);
        }, delay - timeSinceLastCall);
      }
    },
    [callback, delay]
  ) as T;
}

// Memoization for expensive calculations
export function useMemoizedValidator<T>(
  validator: (value: T) => string | null,
  dependencies: any[]
) {
  return useMemo(() => validator, dependencies);
}

// Lazy loading hook for components
export function useLazyComponent<T>(
  factory: () => Promise<{ default: T }>,
  deps: any[] = []
) {
  const [Component, setComponent] = React.useState<T | null>(null);
  const [loading, setLoading] = React.useState(false);
  const [error, setError] = React.useState<Error | null>(null);

  useEffect(() => {
    setLoading(true);
    setError(null);
    
    factory()
      .then((module) => {
        setComponent(module.default);
      })
      .catch((err) => {
        setError(err);
      })
      .finally(() => {
        setLoading(false);
      });
  }, deps);

  return { Component, loading, error };
}

// Virtual scrolling for large lists (if needed)
export function useVirtualScrolling<T>(
  items: T[],
  itemHeight: number,
  containerHeight: number
) {
  const [scrollTop, setScrollTop] = React.useState(0);
  
  const visibleStart = Math.floor(scrollTop / itemHeight);
  const visibleEnd = Math.min(
    visibleStart + Math.ceil(containerHeight / itemHeight) + 1,
    items.length
  );
  
  const visibleItems = items.slice(visibleStart, visibleEnd);
  const totalHeight = items.length * itemHeight;
  const offsetY = visibleStart * itemHeight;
  
  return {
    visibleItems,
    totalHeight,
    offsetY,
    setScrollTop,
  };
}

// Cache hook for API responses
export function useApiCache<T>(key: string, initialValue?: T) {
  const cache = useRef(new Map<string, { data: T; timestamp: number }>());
  const CACHE_DURATION = 5 * 60 * 1000; // 5 minutes

  const get = useCallback((cacheKey: string): T | null => {
    const cached = cache.current.get(cacheKey);
    if (!cached) return null;
    
    if (Date.now() - cached.timestamp > CACHE_DURATION) {
      cache.current.delete(cacheKey);
      return null;
    }
    
    return cached.data;
  }, []);

  const set = useCallback((cacheKey: string, data: T) => {
    cache.current.set(cacheKey, {
      data,
      timestamp: Date.now(),
    });
  }, []);

  const clear = useCallback((cacheKey?: string) => {
    if (cacheKey) {
      cache.current.delete(cacheKey);
    } else {
      cache.current.clear();
    }
  }, []);

  return { get, set, clear };
}

// Optimized form field component with memoization
export const OptimizedFormField = React.memo<{
  name: string;
  value: string;
  onChange: (name: string, value: string) => void;
  onBlur?: (name: string) => void;
  error?: string;
  label: string;
  type?: string;
  placeholder?: string;
  disabled?: boolean;
}>(({
  name,
  value,
  onChange,
  onBlur,
  error,
  label,
  type = 'text',
  placeholder,
  disabled = false,
}) => {
  const handleChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      onChange(name, e.target.value);
    },
    [name, onChange]
  );

  const handleBlur = useCallback(() => {
    onBlur?.(name);
  }, [name, onBlur]);

  return (
    <div className="form-field">
      <label htmlFor={name} className="form-label">
        {label}
      </label>
      <input
        id={name}
        name={name}
        type={type}
        value={value}
        onChange={handleChange}
        onBlur={handleBlur}
        placeholder={placeholder}
        disabled={disabled}
        className={`form-input ${error ? 'error' : ''}`}
        autoComplete="off"
      />
      {error && <span className="error-message">{error}</span>}
    </div>
  );
});

// Performance monitoring utilities
export class PerformanceMonitor {
  private static marks = new Map<string, number>();
  private static measures = new Map<string, number>();

  static mark(name: string): void {
    this.marks.set(name, performance.now());
    if (performance.mark) {
      performance.mark(name);
    }
  }

  static measure(name: string, startMark: string, endMark?: string): number {
    const startTime = this.marks.get(startMark);
    const endTime = endMark ? this.marks.get(endMark) : performance.now();
    
    if (startTime === undefined) {
      console.warn(`Start mark "${startMark}" not found`);
      return 0;
    }
    
    const duration = (endTime || performance.now()) - startTime;
    this.measures.set(name, duration);
    
    if (performance.measure && endMark) {
      try {
        performance.measure(name, startMark, endMark);
      } catch (error) {
        console.warn('Performance measure failed:', error);
      }
    }
    
    return duration;
  }

  static getMetrics(): Record<string, number> {
    return Object.fromEntries(this.measures);
  }

  static logMetrics(): void {
    const metrics = this.getMetrics();
    console.table(metrics);
  }

  static clearMetrics(): void {
    this.marks.clear();
    this.measures.clear();
    if (performance.clearMarks) {
      performance.clearMarks();
    }
    if (performance.clearMeasures) {
      performance.clearMeasures();
    }
  }
}

// Bundle size optimization - dynamic imports
export const lazyLoadComponent = (importFn: () => Promise<any>) => {
  return React.lazy(importFn);
};

// Memory leak prevention
export function useCleanup(cleanup: () => void, deps: any[] = []) {
  useEffect(() => {
    return cleanup;
  }, deps);
}

// Optimized event handlers
export function useStableCallback<T extends (...args: any[]) => any>(
  callback: T
): T {
  const callbackRef = useRef(callback);
  callbackRef.current = callback;

  return useCallback(
    (...args: Parameters<T>) => callbackRef.current(...args),
    []
  ) as T;
}

// Image optimization utilities
export function useImagePreloader(sources: string[]): boolean {
  const [loaded, setLoaded] = React.useState(false);

  useEffect(() => {
    if (sources.length === 0) {
      setLoaded(true);
      return;
    }

    let loadedCount = 0;
    const images: HTMLImageElement[] = [];

    sources.forEach((src) => {
      const img = new Image();
      img.onload = () => {
        loadedCount++;
        if (loadedCount === sources.length) {
          setLoaded(true);
        }
      };
      img.onerror = () => {
        loadedCount++;
        if (loadedCount === sources.length) {
          setLoaded(true);
        }
      };
      img.src = src;
      images.push(img);
    });

    return () => {
      images.forEach((img) => {
        img.onload = null;
        img.onerror = null;
      });
    };
  }, [sources]);

  return loaded;
}

// React import fix
import React from 'react';