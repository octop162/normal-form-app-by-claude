// User Input Page - Complete implementation
import React, { useEffect } from 'react';
import { useFormContext } from '../contexts/FormContext';
import { useFormValidation } from '../hooks/useFormValidation';
import { useFormSubmission } from '../hooks/useFormSubmission';
import { useSessionRecovery } from '../hooks/useSessionRecovery';
import PersonalInfoForm from '../components/forms/PersonalInfoForm';
import AddressForm from '../components/forms/AddressForm';
import PlanOptionForm from '../components/forms/PlanOptionForm';
import FormStepIndicator from '../components/layout/FormStepIndicator';
import FormNavigation from '../components/layout/FormNavigation';
import SessionTimeoutWarning from '../components/layout/SessionTimeoutWarning';
import SessionRecoveryDialog from '../components/layout/SessionRecoveryDialog';
import ValidationSummary from '../components/common/ValidationSummary';
import ErrorMessage from '../components/common/ErrorMessage';
import SessionStatus from '../components/common/SessionStatus';
import { PAGE_TITLES } from '../utils/constants';

const UserInput: React.FC = () => {
  const { currentStep, errors, setCurrentStep } = useFormContext();
  const { hasErrors, isValidating } = useFormValidation();
  const { 
    proceedToConfirm, 
    canProceedToConfirm, 
    isSubmitting, 
    submissionError 
  } = useFormSubmission();
  const {
    showRecoveryDialog,
    recoverSession,
    discardSession,
    closeRecoveryDialog
  } = useSessionRecovery();

  // Set current step to input when component mounts
  useEffect(() => {
    if (currentStep !== 'input') {
      setCurrentStep('input');
    }
  }, [currentStep, setCurrentStep]);

  // Set page title
  useEffect(() => {
    document.title = PAGE_TITLES.INPUT;
  }, []);

  const handleNext = async () => {
    await proceedToConfirm();
  };

  return (
    <div className="user-input-page">
      <div className="page-container">
        <header className="page-header">
          <h1>会員登録</h1>
          <FormStepIndicator currentStep={currentStep} />
          <SessionStatus showDetails={false} />
        </header>

        <main className="page-content">
          {submissionError && (
            <ErrorMessage 
              error={submissionError} 
              data-testid="submission-error"
            />
          )}

          {hasErrors && (
            <ValidationSummary 
              errors={errors} 
              showSummary={true}
              data-testid="validation-summary"
            />
          )}

          <form className="registration-form" noValidate>
            <section className="form-section">
              <PersonalInfoForm />
            </section>

            <section className="form-section">
              <AddressForm />
            </section>

            <section className="form-section">
              <PlanOptionForm />
            </section>
          </form>
        </main>

        <footer className="page-footer">
          <FormNavigation
            currentStep={currentStep}
            canGoNext={canProceedToConfirm}
            canGoPrev={false}
            onNext={handleNext}
            onPrev={() => {}}
            isLoading={isSubmitting || isValidating}
          />
        </footer>
      </div>

      <SessionTimeoutWarning />
      
      {showRecoveryDialog && (
        <SessionRecoveryDialog
          onRecover={recoverSession}
          onDiscard={discardSession}
          onClose={closeRecoveryDialog}
        />
      )}
    </div>
  );
};

export default UserInput;