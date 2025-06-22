// Custom hook for form submission logic
import { useCallback, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useFormContext } from '../contexts/FormContext';
import { useFormValidation } from './useFormValidation';
import { useCreateUser, useCreateSession, useUpdateSession } from './useApi';
import { transformFormDataToApiRequest } from '../utils/apiTransformers';
import { SUCCESS_MESSAGES, ERROR_MESSAGES } from '../utils/constants';
import type { ApiError } from '../types/api';

interface UseFormSubmissionReturn {
  // Navigation functions
  proceedToConfirm: () => Promise<boolean>;
  proceedToComplete: () => Promise<boolean>;
  returnToInput: () => void;
  
  // Submission state
  isSubmitting: boolean;
  submissionError: ApiError | null;
  
  // Session management
  saveSession: () => Promise<boolean>;
  
  // Helper functions
  canProceedToConfirm: boolean;
  canSubmitForm: boolean;
}

export const useFormSubmission = (): UseFormSubmissionReturn => {
  const navigate = useNavigate();
  const { 
    formData, 
    currentStep, 
    setCurrentStep, 
    isLoading, 
    setIsLoading,
    sessionId,
    setSessionId,
    resetForm 
  } = useFormContext();
  
  const { validateFormStep, canProceedToNext } = useFormValidation();
  const createUser = useCreateUser();
  const createSession = useCreateSession();
  const updateSession = useUpdateSession();
  
  const [submissionError, setSubmissionError] = useState<ApiError | null>(null);

  // Check if we can proceed to confirm step
  const canProceedToConfirm = canProceedToNext;

  // Check if we can submit the form
  const canSubmitForm = currentStep === 'confirm' && canProceedToNext;

  // Save session data
  const saveSession = useCallback(async (): Promise<boolean> => {
    try {
      setIsLoading(true);
      
      const apiRequest = transformFormDataToApiRequest(formData);
      const sessionData = {
        user_data: apiRequest
      };

      if (sessionId) {
        // Update existing session
        await updateSession.execute(sessionId, sessionData);
      } else {
        // Create new session
        const result = await createSession.execute(sessionData);
        setSessionId(result.session_id);
      }

      return true;
    } catch (error) {
      console.error('Session save failed:', error);
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [formData, sessionId, setSessionId, setIsLoading, createSession, updateSession]);

  // Proceed to confirmation step
  const proceedToConfirm = useCallback(async (): Promise<boolean> => {
    try {
      setSubmissionError(null);
      setIsLoading(true);

      // Validate form data
      const isValid = await validateFormStep('input');
      if (!isValid) {
        return false;
      }

      // Save session data
      const sessionSaved = await saveSession();
      if (!sessionSaved) {
        setSubmissionError({
          code: 'SESSION_SAVE_FAILED',
          message: 'セッションの保存に失敗しました'
        });
        return false;
      }

      // Navigate to confirm step
      setCurrentStep('confirm');
      navigate('/confirm');
      return true;

    } catch (error) {
      const apiError = error as ApiError;
      setSubmissionError(apiError);
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [validateFormStep, saveSession, setCurrentStep, navigate, setIsLoading]);

  // Proceed to completion step (submit form)
  const proceedToComplete = useCallback(async (): Promise<boolean> => {
    try {
      setSubmissionError(null);
      setIsLoading(true);

      // Final validation including server-side validation
      const isValid = await validateFormStep('confirm');
      if (!isValid) {
        return false;
      }

      // Submit user data
      const apiRequest = transformFormDataToApiRequest(formData);
      await createUser.execute(apiRequest);

      // Clean up session after successful submission
      if (sessionId) {
        try {
          // Note: We don't await this as it's cleanup
          // and shouldn't block the success flow
          import('../services/apiClient').then(({ ApiService }) => {
            ApiService.deleteSession(sessionId);
          });
        } catch (cleanupError) {
          // Log but don't fail on cleanup error
          console.warn('Session cleanup failed:', cleanupError);
        }
      }

      // Navigate to complete step
      setCurrentStep('complete');
      navigate('/complete');
      return true;

    } catch (error) {
      const apiError = error as ApiError;
      setSubmissionError(apiError);
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [validateFormStep, formData, createUser, sessionId, setCurrentStep, navigate, setIsLoading]);

  // Return to input step
  const returnToInput = useCallback(() => {
    setCurrentStep('input');
    navigate('/');
  }, [setCurrentStep, navigate]);

  return {
    // Navigation functions
    proceedToConfirm,
    proceedToComplete,
    returnToInput,
    
    // Submission state
    isSubmitting: isLoading || createUser.isLoading,
    submissionError: submissionError || createUser.error,
    
    // Session management
    saveSession,
    
    // Helper functions
    canProceedToConfirm,
    canSubmitForm
  };
};