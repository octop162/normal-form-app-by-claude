// Custom hook for session recovery management
import { useCallback, useEffect, useState } from 'react';
import { useFormContext } from '../contexts/FormContext';
import { SESSION_CONFIG } from '../utils/constants';
import type { UserFormData } from '../types/form';

interface UseSessionRecoveryReturn {
  // Recovery state
  hasRecoverableSession: boolean;
  showRecoveryDialog: boolean;
  
  // Recovery actions
  recoverSession: () => void;
  discardSession: () => void;
  closeRecoveryDialog: () => void;
  
  // Recovery data
  recoveryData: UserFormData | null;
  recoveryTimestamp: Date | null;
}

export const useSessionRecovery = (): UseSessionRecoveryReturn => {
  const { formData, updateFormData, hasFormData, resetForm } = useFormContext();
  const [showRecoveryDialog, setShowRecoveryDialog] = useState(false);
  const [recoveryData, setRecoveryData] = useState<UserFormData | null>(null);
  const [recoveryTimestamp, setRecoveryTimestamp] = useState<Date | null>(null);

  // Check if there's recoverable session data
  const hasRecoverableSession = useCallback((): boolean => {
    try {
      const savedData = localStorage.getItem(SESSION_CONFIG.STORAGE_KEYS.FORM_DATA);
      const lastSaved = localStorage.getItem(SESSION_CONFIG.STORAGE_KEYS.LAST_SAVED);
      
      if (!savedData || !lastSaved) {
        return false;
      }

      const lastSavedTime = new Date(lastSaved);
      const now = new Date();
      const hoursDiff = (now.getTime() - lastSavedTime.getTime()) / (1000 * 60 * 60);

      // Check if session is still valid (within timeout period)
      if (hoursDiff >= SESSION_CONFIG.TIMEOUT_HOURS) {
        return false;
      }

      const parsedData = JSON.parse(savedData) as UserFormData;
      
      // Check if the saved data has meaningful content
      const hasMeaningfulData = Object.values(parsedData).some(value => {
        if (Array.isArray(value)) {
          return value.length > 0;
        }
        return value && value.toString().trim() !== '';
      });

      return hasMeaningfulData;
    } catch (error) {
      console.error('Error checking recoverable session:', error);
      return false;
    }
  }, []);

  // Load recovery data from storage
  const loadRecoveryData = useCallback((): void => {
    try {
      const savedData = localStorage.getItem(SESSION_CONFIG.STORAGE_KEYS.FORM_DATA);
      const lastSaved = localStorage.getItem(SESSION_CONFIG.STORAGE_KEYS.LAST_SAVED);
      
      if (savedData && lastSaved) {
        const parsedData = JSON.parse(savedData) as UserFormData;
        const timestamp = new Date(lastSaved);
        
        setRecoveryData(parsedData);
        setRecoveryTimestamp(timestamp);
      }
    } catch (error) {
      console.error('Error loading recovery data:', error);
      setRecoveryData(null);
      setRecoveryTimestamp(null);
    }
  }, []);

  // Show recovery dialog
  const showRecovery = useCallback((): void => {
    if (hasRecoverableSession()) {
      loadRecoveryData();
      setShowRecoveryDialog(true);
    }
  }, [hasRecoverableSession, loadRecoveryData]);

  // Recover session data
  const recoverSession = useCallback((): void => {
    if (recoveryData) {
      updateFormData(recoveryData);
      setShowRecoveryDialog(false);
      
      // Clear recovery data since it's now loaded
      setRecoveryData(null);
      setRecoveryTimestamp(null);
    }
  }, [recoveryData, updateFormData]);

  // Discard saved session data
  const discardSession = useCallback((): void => {
    // Clear from localStorage
    Object.values(SESSION_CONFIG.STORAGE_KEYS).forEach(key => {
      localStorage.removeItem(key);
    });
    
    // Reset form
    resetForm();
    
    // Clear recovery state
    setRecoveryData(null);
    setRecoveryTimestamp(null);
    setShowRecoveryDialog(false);
  }, [resetForm]);

  // Close recovery dialog without action
  const closeRecoveryDialog = useCallback((): void => {
    setShowRecoveryDialog(false);
  }, []);

  // Check for recoverable session on mount
  useEffect(() => {
    // Only show recovery dialog if current form is empty
    if (!hasFormData() && hasRecoverableSession()) {
      // Small delay to ensure page is fully loaded
      setTimeout(() => {
        showRecovery();
      }, 1000);
    }
  }, [hasFormData, hasRecoverableSession, showRecovery]);

  // Auto-cleanup expired sessions
  useEffect(() => {
    const cleanupExpiredSessions = () => {
      const lastSaved = localStorage.getItem(SESSION_CONFIG.STORAGE_KEYS.LAST_SAVED);
      if (lastSaved) {
        const lastSavedTime = new Date(lastSaved);
        const now = new Date();
        const hoursDiff = (now.getTime() - lastSavedTime.getTime()) / (1000 * 60 * 60);

        if (hoursDiff >= SESSION_CONFIG.TIMEOUT_HOURS) {
          // Session expired, clean up
          Object.values(SESSION_CONFIG.STORAGE_KEYS).forEach(key => {
            localStorage.removeItem(key);
          });
          
          setRecoveryData(null);
          setRecoveryTimestamp(null);
        }
      }
    };

    // Check on mount
    cleanupExpiredSessions();

    // Check periodically (every 5 minutes)
    const interval = setInterval(cleanupExpiredSessions, 5 * 60 * 1000);

    return () => clearInterval(interval);
  }, []);

  return {
    // Recovery state
    hasRecoverableSession: hasRecoverableSession(),
    showRecoveryDialog,
    
    // Recovery actions
    recoverSession,
    discardSession,
    closeRecoveryDialog,
    
    // Recovery data
    recoveryData,
    recoveryTimestamp
  };
};