// Custom hook for form validation
import { useCallback, useMemo } from 'react';
import { useFormContext } from '../contexts/FormContext';
import { validateForm, validateField, userFormSchema } from '../validation/schemas';
import { useValidateUser } from './useApi';
import { transformFormDataToValidateRequest } from '../utils/apiTransformers';
import type { UserFormData, FormValidationErrors } from '../types/form';
import type { ApiError } from '../types/api';

interface UseFormValidationReturn {
  // Validation functions
  validateSingleField: (fieldName: keyof UserFormData, value: any) => string | null;
  validateCurrentFormData: () => FormValidationErrors;
  validateFormStep: (step: 'input' | 'confirm') => Promise<boolean>;
  
  // Validation state
  hasErrors: boolean;
  isValidating: boolean;
  validationError: ApiError | null;
  
  // Helper functions
  canProceedToNext: boolean;
  getFieldError: (fieldName: keyof UserFormData) => string | undefined;
  clearAllErrors: () => void;
}

export const useFormValidation = (): UseFormValidationReturn => {
  const { formData, errors, setErrors, clearFieldError } = useFormContext();
  const serverValidation = useValidateUser();

  // Validate a single field
  const validateSingleField = useCallback((
    fieldName: keyof UserFormData,
    value: any
  ): string | null => {
    return validateField(fieldName, value, formData);
  }, [formData]);

  // Validate all current form data (client-side only)
  const validateCurrentFormData = useCallback((): FormValidationErrors => {
    return validateForm(formData);
  }, [formData]);

  // Validate form step (includes server-side validation for confirm step)
  const validateFormStep = useCallback(async (step: 'input' | 'confirm'): Promise<boolean> => {
    // First, run client-side validation
    const clientErrors = validateCurrentFormData();
    
    if (Object.keys(clientErrors).length > 0) {
      setErrors(clientErrors);
      return false;
    }

    // For confirm step, also run server-side validation
    if (step === 'confirm') {
      try {
        const validateRequest = transformFormDataToValidateRequest(formData);
        const result = await serverValidation.execute(validateRequest);
        
        if (!result.valid && result.errors) {
          // Transform server errors to match frontend field names
          const serverErrors: FormValidationErrors = {};
          Object.entries(result.errors).forEach(([serverField, message]) => {
            // Convert snake_case to camelCase for frontend
            const frontendField = serverField.replace(/_([a-z])/g, (_, letter) => letter.toUpperCase());
            serverErrors[frontendField] = message;
          });
          
          setErrors(serverErrors);
          return false;
        }
        
        // Clear errors if validation passed
        setErrors({});
        return true;
        
      } catch (error) {
        // Handle server validation error
        const apiError = error as ApiError;
        setErrors({ general: apiError.message || 'サーバーバリデーションエラーが発生しました' });
        return false;
      }
    }
    
    // Clear errors if client validation passed
    setErrors({});
    return true;
  }, [formData, validateCurrentFormData, setErrors, serverValidation]);

  // Check if there are any validation errors
  const hasErrors = useMemo(() => {
    return Object.keys(errors).length > 0;
  }, [errors]);

  // Check if we can proceed to next step
  const canProceedToNext = useMemo(() => {
    // Check if all required fields are filled
    const requiredFields: (keyof UserFormData)[] = [
      'lastName', 'firstName', 'lastNameKana', 'firstNameKana',
      'phone1', 'phone2', 'phone3',
      'postalCode1', 'postalCode2',
      'prefecture', 'city', 'banchi',
      'email', 'emailConfirm',
      'planType'
    ];

    const allRequiredFilled = requiredFields.every(field => {
      const value = formData[field];
      if (Array.isArray(value)) {
        return true; // optionTypes is not required to have values
      }
      return value && value.toString().trim() !== '';
    });

    return allRequiredFilled && !hasErrors;
  }, [formData, hasErrors]);

  // Get error for a specific field
  const getFieldError = useCallback((fieldName: keyof UserFormData): string | undefined => {
    return errors[fieldName];
  }, [errors]);

  // Clear all validation errors
  const clearAllErrors = useCallback(() => {
    setErrors({});
  }, [setErrors]);

  return {
    // Validation functions
    validateSingleField,
    validateCurrentFormData,
    validateFormStep,
    
    // Validation state
    hasErrors,
    isValidating: serverValidation.isLoading,
    validationError: serverValidation.error,
    
    // Helper functions
    canProceedToNext,
    getFieldError,
    clearAllErrors
  };
};