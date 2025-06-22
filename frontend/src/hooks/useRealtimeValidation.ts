// Custom hook for real-time field validation
import { useCallback, useRef } from 'react';
import { useFormContext } from '../contexts/FormContext';
import { validateField } from '../validation/schemas';
import type { UserFormData } from '../types/form';

interface UseRealtimeValidationReturn {
  // Event handlers for form fields
  createFieldChangeHandler: (fieldName: keyof UserFormData) => (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => void;
  createFieldBlurHandler: (fieldName: keyof UserFormData) => (e: React.FocusEvent<HTMLInputElement | HTMLSelectElement>) => void;
  createCheckboxChangeHandler: (fieldName: 'optionTypes') => (values: string[]) => void;
  
  // Manual validation triggers
  validateFieldManually: (fieldName: keyof UserFormData, value: any) => void;
  clearFieldErrorManually: (fieldName: keyof UserFormData) => void;
}

export const useRealtimeValidation = (): UseRealtimeValidationReturn => {
  const { formData, updateFormData, errors, setErrors, clearFieldError } = useFormContext();
  
  // Debounce validation to avoid excessive validation calls
  const validationTimeouts = useRef<Record<string, NodeJS.Timeout>>({});

  // Set field error
  const setFieldError = useCallback((fieldName: keyof UserFormData, error: string) => {
    setErrors({ ...errors, [fieldName]: error });
  }, [errors, setErrors]);

  // Create change handler for input fields
  const createFieldChangeHandler = useCallback((fieldName: keyof UserFormData) => {
    return (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
      const value = e.target.value;
      
      // Update form data immediately
      updateFormData({ [fieldName]: value });
      
      // Clear existing error when user starts typing
      if (errors[fieldName]) {
        clearFieldError(fieldName);
      }
      
      // Clear any pending validation
      if (validationTimeouts.current[fieldName]) {
        clearTimeout(validationTimeouts.current[fieldName]);
      }
      
      // Debounced validation for better UX
      validationTimeouts.current[fieldName] = setTimeout(() => {
        const error = validateField(fieldName, value, { ...formData, [fieldName]: value });
        if (error) {
          setFieldError(fieldName, error);
        }
      }, 500); // 500ms debounce
    };
  }, [formData, updateFormData, errors, clearFieldError, setFieldError]);

  // Create blur handler for input fields (immediate validation)
  const createFieldBlurHandler = useCallback((fieldName: keyof UserFormData) => {
    return (e: React.FocusEvent<HTMLInputElement | HTMLSelectElement>) => {
      const value = e.target.value;
      
      // Clear any pending debounced validation
      if (validationTimeouts.current[fieldName]) {
        clearTimeout(validationTimeouts.current[fieldName]);
        delete validationTimeouts.current[fieldName];
      }
      
      // Immediate validation on blur
      const currentFormData = { ...formData, [fieldName]: value };
      const error = validateField(fieldName, value, currentFormData);
      
      if (error) {
        setFieldError(fieldName, error);
      } else {
        clearFieldError(fieldName);
      }
      
      // Special handling for email confirmation
      if (fieldName === 'email' && formData.emailConfirm) {
        const emailConfirmError = validateField('emailConfirm', formData.emailConfirm, currentFormData);
        if (emailConfirmError) {
          setFieldError('emailConfirm', emailConfirmError);
        } else {
          clearFieldError('emailConfirm');
        }
      }
      
      if (fieldName === 'emailConfirm') {
        const emailConfirmError = validateField('emailConfirm', value, currentFormData);
        if (emailConfirmError) {
          setFieldError('emailConfirm', emailConfirmError);
        } else {
          clearFieldError('emailConfirm');
        }
      }
    };
  }, [formData, clearFieldError, setFieldError]);

  // Create change handler for checkbox groups
  const createCheckboxChangeHandler = useCallback((fieldName: 'optionTypes') => {
    return (values: string[]) => {
      // Update form data immediately
      updateFormData({ [fieldName]: values });
      
      // Clear existing error
      if (errors[fieldName]) {
        clearFieldError(fieldName);
      }
      
      // Validate immediately for checkboxes (no debounce needed)
      const error = validateField(fieldName, values, { ...formData, [fieldName]: values });
      if (error) {
        setFieldError(fieldName, error);
      }
    };
  }, [formData, updateFormData, errors, clearFieldError, setFieldError]);

  // Manual validation trigger
  const validateFieldManually = useCallback((fieldName: keyof UserFormData, value: any) => {
    const error = validateField(fieldName, value, formData);
    if (error) {
      setFieldError(fieldName, error);
    } else {
      clearFieldError(fieldName);
    }
  }, [formData, setFieldError, clearFieldError]);

  // Manual error clearing
  const clearFieldErrorManually = useCallback((fieldName: keyof UserFormData) => {
    clearFieldError(fieldName);
  }, [clearFieldError]);

  return {
    createFieldChangeHandler,
    createFieldBlurHandler,
    createCheckboxChangeHandler,
    validateFieldManually,
    clearFieldErrorManually
  };
};