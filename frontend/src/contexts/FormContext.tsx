// Form state management using React Context
import React, { createContext, useContext, useReducer, useCallback, useEffect } from 'react';
import type { 
  FormContextType, 
  UserFormData, 
  FormStep, 
  FormValidationErrors,
  SessionData
} from '../types/form';

// Initial form data
const initialFormData: UserFormData = {
  lastName: '',
  firstName: '',
  lastNameKana: '',
  firstNameKana: '',
  phone1: '',
  phone2: '',
  phone3: '',
  postalCode1: '',
  postalCode2: '',
  prefecture: '',
  city: '',
  town: '',
  chome: '',
  banchi: '',
  go: '',
  building: '',
  room: '',
  email: '',
  emailConfirm: '',
  planType: '',
  optionTypes: []
};

// Form state interface
interface FormState {
  formData: UserFormData;
  currentStep: FormStep;
  isLoading: boolean;
  errors: FormValidationErrors;
  sessionId?: string;
}

// Initial state
const initialState: FormState = {
  formData: initialFormData,
  currentStep: 'input',
  isLoading: false,
  errors: {},
  sessionId: undefined
};

// Action types
type FormAction =
  | { type: 'UPDATE_FORM_DATA'; payload: Partial<UserFormData> }
  | { type: 'RESET_FORM' }
  | { type: 'SET_CURRENT_STEP'; payload: FormStep }
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'SET_ERRORS'; payload: FormValidationErrors }
  | { type: 'SET_SESSION_ID'; payload: string | undefined }
  | { type: 'LOAD_FROM_SESSION'; payload: { formData: UserFormData; sessionId: string } };

// Reducer function
const formReducer = (state: FormState, action: FormAction): FormState => {
  switch (action.type) {
    case 'UPDATE_FORM_DATA':
      return {
        ...state,
        formData: {
          ...state.formData,
          ...action.payload
        }
      };
    
    case 'RESET_FORM':
      return {
        ...initialState,
        sessionId: undefined
      };
    
    case 'SET_CURRENT_STEP':
      return {
        ...state,
        currentStep: action.payload
      };
    
    case 'SET_LOADING':
      return {
        ...state,
        isLoading: action.payload
      };
    
    case 'SET_ERRORS':
      return {
        ...state,
        errors: action.payload
      };
    
    case 'SET_SESSION_ID':
      return {
        ...state,
        sessionId: action.payload
      };
    
    case 'LOAD_FROM_SESSION':
      return {
        ...state,
        formData: action.payload.formData,
        sessionId: action.payload.sessionId
      };
    
    default:
      return state;
  }
};

// Create context
const FormContext = createContext<FormContextType | null>(null);

// Session storage keys
const STORAGE_KEYS = {
  FORM_DATA: 'membershipForm_data',
  SESSION_ID: 'membershipForm_sessionId',
  LAST_SAVED: 'membershipForm_lastSaved'
};

// Provider component
interface FormProviderProps {
  children: React.ReactNode;
}

export const FormProvider: React.FC<FormProviderProps> = ({ children }) => {
  const [state, dispatch] = useReducer(formReducer, initialState);

  // Load data from session storage on mount
  useEffect(() => {
    try {
      const savedData = localStorage.getItem(STORAGE_KEYS.FORM_DATA);
      const savedSessionId = localStorage.getItem(STORAGE_KEYS.SESSION_ID);
      const lastSaved = localStorage.getItem(STORAGE_KEYS.LAST_SAVED);

      if (savedData && savedSessionId && lastSaved) {
        const parsedData = JSON.parse(savedData);
        const lastSavedTime = new Date(lastSaved);
        const now = new Date();
        const hoursDiff = (now.getTime() - lastSavedTime.getTime()) / (1000 * 60 * 60);

        // Check if session is still valid (4 hours timeout)
        if (hoursDiff < 4) {
          dispatch({
            type: 'LOAD_FROM_SESSION',
            payload: {
              formData: parsedData,
              sessionId: savedSessionId
            }
          });
        } else {
          // Clear expired session
          localStorage.removeItem(STORAGE_KEYS.FORM_DATA);
          localStorage.removeItem(STORAGE_KEYS.SESSION_ID);
          localStorage.removeItem(STORAGE_KEYS.LAST_SAVED);
        }
      }
    } catch (error) {
      console.error('Failed to load form data from storage:', error);
    }
  }, []);

  // Save data to session storage whenever form data changes
  useEffect(() => {
    if (state.formData !== initialFormData) {
      try {
        localStorage.setItem(STORAGE_KEYS.FORM_DATA, JSON.stringify(state.formData));
        localStorage.setItem(STORAGE_KEYS.LAST_SAVED, new Date().toISOString());
        
        if (state.sessionId) {
          localStorage.setItem(STORAGE_KEYS.SESSION_ID, state.sessionId);
        }
      } catch (error) {
        console.error('Failed to save form data to storage:', error);
      }
    }
  }, [state.formData, state.sessionId]);

  // Context value methods
  const updateFormData = useCallback((data: Partial<UserFormData>) => {
    dispatch({ type: 'UPDATE_FORM_DATA', payload: data });
  }, []);

  const resetForm = useCallback(() => {
    // Clear storage
    localStorage.removeItem(STORAGE_KEYS.FORM_DATA);
    localStorage.removeItem(STORAGE_KEYS.SESSION_ID);
    localStorage.removeItem(STORAGE_KEYS.LAST_SAVED);
    
    dispatch({ type: 'RESET_FORM' });
  }, []);

  const setCurrentStep = useCallback((step: FormStep) => {
    dispatch({ type: 'SET_CURRENT_STEP', payload: step });
  }, []);

  const setIsLoading = useCallback((loading: boolean) => {
    dispatch({ type: 'SET_LOADING', payload: loading });
  }, []);

  const setErrors = useCallback((errors: FormValidationErrors) => {
    dispatch({ type: 'SET_ERRORS', payload: errors });
  }, []);

  const setSessionId = useCallback((id: string | undefined) => {
    dispatch({ type: 'SET_SESSION_ID', payload: id });
  }, []);

  // Clear specific field error
  const clearFieldError = useCallback((fieldName: string) => {
    if (state.errors[fieldName]) {
      const newErrors = { ...state.errors };
      delete newErrors[fieldName];
      setErrors(newErrors);
    }
  }, [state.errors, setErrors]);

  // Check if form has any data
  const hasFormData = useCallback(() => {
    return Object.values(state.formData).some(value => {
      if (Array.isArray(value)) {
        return value.length > 0;
      }
      return value !== '' && value !== undefined && value !== null;
    });
  }, [state.formData]);

  // Get session timeout warning
  const getSessionTimeoutWarning = useCallback((): { show: boolean; remainingMinutes: number } => {
    const lastSaved = localStorage.getItem(STORAGE_KEYS.LAST_SAVED);
    if (!lastSaved) {
      return { show: false, remainingMinutes: 0 };
    }

    const lastSavedTime = new Date(lastSaved);
    const now = new Date();
    const minutesDiff = (now.getTime() - lastSavedTime.getTime()) / (1000 * 60);
    const remainingMinutes = Math.max(0, 240 - minutesDiff); // 4 hours = 240 minutes

    // Show warning when less than 15 minutes remaining
    const show = remainingMinutes > 0 && remainingMinutes < 15;

    return { show, remainingMinutes: Math.ceil(remainingMinutes) };
  }, []);

  // Extend session
  const extendSession = useCallback(() => {
    localStorage.setItem(STORAGE_KEYS.LAST_SAVED, new Date().toISOString());
  }, []);

  const contextValue: FormContextType = {
    formData: state.formData,
    updateFormData,
    resetForm,
    currentStep: state.currentStep,
    setCurrentStep,
    isLoading: state.isLoading,
    setIsLoading,
    errors: state.errors,
    setErrors,
    sessionId: state.sessionId,
    setSessionId,
    // Additional utility methods
    clearFieldError,
    hasFormData,
    getSessionTimeoutWarning,
    extendSession
  };

  return (
    <FormContext.Provider value={contextValue}>
      {children}
    </FormContext.Provider>
  );
};

// Hook to use form context
export const useFormContext = (): FormContextType => {
  const context = useContext(FormContext);
  if (!context) {
    throw new Error('useFormContext must be used within a FormProvider');
  }
  return context;
};

// Export context for testing purposes
export { FormContext };