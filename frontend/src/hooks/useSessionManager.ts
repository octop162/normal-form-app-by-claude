// Custom hook for enhanced session management
import { useCallback, useEffect, useRef, useState } from 'react';
import { useFormContext } from '../contexts/FormContext';
import { useCreateSession, useUpdateSession, useDeleteSession } from './useApi';
import { transformFormDataToApiRequest } from '../utils/apiTransformers';
import { SESSION_CONFIG } from '../utils/constants';
import type { ApiError } from '../types/api';

interface UseSessionManagerReturn {
  // Session operations
  saveSession: () => Promise<boolean>;
  loadSession: () => Promise<boolean>;
  clearSession: () => Promise<boolean>;
  extendSession: () => Promise<boolean>;
  
  // Session state
  isSessionLoading: boolean;
  sessionError: ApiError | null;
  sessionExists: boolean;
  
  // Timeout management
  timeoutWarning: {
    show: boolean;
    remainingMinutes: number;
  };
  
  // Auto-save management
  enableAutoSave: () => void;
  disableAutoSave: () => void;
  isAutoSaveEnabled: boolean;
}

export const useSessionManager = (): UseSessionManagerReturn => {
  const { 
    formData, 
    sessionId, 
    setSessionId, 
    getSessionTimeoutWarning, 
    extendSession: contextExtendSession,
    hasFormData 
  } = useFormContext();
  
  // API hooks
  const createSession = useCreateSession();
  const updateSession = useUpdateSession();
  const deleteSession = useDeleteSession();
  
  // Local state
  const [sessionError, setSessionError] = useState<ApiError | null>(null);
  const [sessionExists, setSessionExists] = useState<boolean>(false);
  const [timeoutWarning, setTimeoutWarning] = useState({ show: false, remainingMinutes: 0 });
  const [isAutoSaveEnabled, setIsAutoSaveEnabled] = useState<boolean>(true);
  
  // Refs for intervals and timeouts
  const autoSaveInterval = useRef<NodeJS.Timeout | null>(null);
  const timeoutCheckInterval = useRef<NodeJS.Timeout | null>(null);
  const lastSaveTime = useRef<Date>(new Date());

  // Check if session loading
  const isSessionLoading = createSession.isLoading || updateSession.isLoading || deleteSession.isLoading;

  // Save session to server
  const saveSession = useCallback(async (): Promise<boolean> => {
    try {
      setSessionError(null);
      
      // Don't save if no meaningful form data
      if (!hasFormData()) {
        return true;
      }

      const apiRequest = transformFormDataToApiRequest(formData);
      const sessionData = {
        user_data: apiRequest,
        expires_at: new Date(Date.now() + SESSION_CONFIG.TIMEOUT_HOURS * 60 * 60 * 1000).toISOString()
      };

      if (sessionId) {
        // Update existing session
        await updateSession.execute(sessionId, sessionData);
      } else {
        // Create new session
        const result = await createSession.execute(sessionData);
        setSessionId(result.session_id);
        setSessionExists(true);
      }

      lastSaveTime.current = new Date();
      
      // Update local storage timestamp
      localStorage.setItem(SESSION_CONFIG.STORAGE_KEYS.LAST_SAVED, new Date().toISOString());
      
      return true;
    } catch (error) {
      const apiError = error as ApiError;
      setSessionError(apiError);
      console.error('Session save failed:', apiError);
      return false;
    }
  }, [formData, sessionId, setSessionId, hasFormData, createSession, updateSession]);

  // Load session from server (currently not implemented in API)
  const loadSession = useCallback(async (): Promise<boolean> => {
    // This would be implemented when we have a session retrieval endpoint
    // For now, we rely on localStorage which is handled in FormContext
    return true;
  }, []);

  // Clear session from server and local storage
  const clearSession = useCallback(async (): Promise<boolean> => {
    try {
      setSessionError(null);
      
      if (sessionId) {
        await deleteSession.execute(sessionId);
      }
      
      // Clear local storage
      Object.values(SESSION_CONFIG.STORAGE_KEYS).forEach(key => {
        localStorage.removeItem(key);
      });
      
      setSessionId(undefined);
      setSessionExists(false);
      
      return true;
    } catch (error) {
      const apiError = error as ApiError;
      setSessionError(apiError);
      console.error('Session clear failed:', apiError);
      return false;
    }
  }, [sessionId, setSessionId, deleteSession]);

  // Extend session timeout
  const extendSession = useCallback(async (): Promise<boolean> => {
    try {
      // Update context (local storage)
      contextExtendSession();
      
      // Update server session if exists
      if (sessionId && hasFormData()) {
        const apiRequest = transformFormDataToApiRequest(formData);
        const sessionData = {
          user_data: apiRequest,
          expires_at: new Date(Date.now() + SESSION_CONFIG.TIMEOUT_HOURS * 60 * 60 * 1000).toISOString()
        };
        
        await updateSession.execute(sessionId, sessionData);
      }
      
      setTimeoutWarning({ show: false, remainingMinutes: 0 });
      return true;
    } catch (error) {
      const apiError = error as ApiError;
      setSessionError(apiError);
      console.error('Session extend failed:', apiError);
      return false;
    }
  }, [contextExtendSession, sessionId, formData, hasFormData, updateSession]);

  // Enable auto-save
  const enableAutoSave = useCallback(() => {
    setIsAutoSaveEnabled(true);
  }, []);

  // Disable auto-save
  const disableAutoSave = useCallback(() => {
    setIsAutoSaveEnabled(false);
  }, []);

  // Auto-save functionality
  useEffect(() => {
    if (!isAutoSaveEnabled) return;

    // Auto-save every 30 seconds if there's meaningful form data
    const interval = setInterval(async () => {
      if (hasFormData() && !isSessionLoading) {
        // Only save if it's been more than 10 seconds since last save
        const timeSinceLastSave = Date.now() - lastSaveTime.current.getTime();
        if (timeSinceLastSave > 10000) {
          await saveSession();
        }
      }
    }, 30000);

    autoSaveInterval.current = interval;

    return () => {
      if (autoSaveInterval.current) {
        clearInterval(autoSaveInterval.current);
      }
    };
  }, [isAutoSaveEnabled, hasFormData, isSessionLoading, saveSession]);

  // Session timeout monitoring
  useEffect(() => {
    const checkTimeout = () => {
      const warning = getSessionTimeoutWarning();
      setTimeoutWarning(warning);
    };

    // Check immediately
    checkTimeout();

    // Check every minute
    const interval = setInterval(checkTimeout, 60000);
    timeoutCheckInterval.current = interval;

    return () => {
      if (timeoutCheckInterval.current) {
        clearInterval(timeoutCheckInterval.current);
      }
    };
  }, [getSessionTimeoutWarning]);

  // Update session exists flag when sessionId changes
  useEffect(() => {
    setSessionExists(!!sessionId);
  }, [sessionId]);

  // Save session when form data changes (debounced)
  useEffect(() => {
    if (!isAutoSaveEnabled || !hasFormData()) return;

    const timeoutId = setTimeout(async () => {
      if (!isSessionLoading) {
        await saveSession();
      }
    }, 2000); // 2 second debounce

    return () => clearTimeout(timeoutId);
  }, [formData, isAutoSaveEnabled, hasFormData, isSessionLoading, saveSession]);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (autoSaveInterval.current) {
        clearInterval(autoSaveInterval.current);
      }
      if (timeoutCheckInterval.current) {
        clearInterval(timeoutCheckInterval.current);
      }
    };
  }, []);

  // Save session before page unload
  useEffect(() => {
    const handleBeforeUnload = async (event: BeforeUnloadEvent) => {
      if (hasFormData() && isAutoSaveEnabled) {
        // For immediate save on page unload, we use the synchronous approach
        // This is a limitation of the beforeunload event
        event.preventDefault();
        event.returnValue = '入力内容が保存されていない可能性があります。ページを離れますか？';
        
        // Attempt to save (though this may not complete)
        saveSession();
      }
    };

    window.addEventListener('beforeunload', handleBeforeUnload);

    return () => {
      window.removeEventListener('beforeunload', handleBeforeUnload);
    };
  }, [hasFormData, isAutoSaveEnabled, saveSession]);

  return {
    // Session operations
    saveSession,
    loadSession,
    clearSession,
    extendSession,
    
    // Session state
    isSessionLoading,
    sessionError,
    sessionExists,
    
    // Timeout management
    timeoutWarning,
    
    // Auto-save management
    enableAutoSave,
    disableAutoSave,
    isAutoSaveEnabled
  };
};