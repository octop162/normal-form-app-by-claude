// User Confirm Page - Complete implementation
import React, { useEffect, useMemo } from 'react';
import { useFormContext } from '../contexts/FormContext';
import { useFormValidation } from '../hooks/useFormValidation';
import { useFormSubmission } from '../hooks/useFormSubmission';
import { useGetPlans, useGetOptions } from '../hooks/useApi';
import FormStepIndicator from '../components/layout/FormStepIndicator';
import FormNavigation from '../components/layout/FormNavigation';
import SessionTimeoutWarning from '../components/layout/SessionTimeoutWarning';
import ValidationSummary from '../components/common/ValidationSummary';
import ErrorMessage from '../components/common/ErrorMessage';
import LoadingSpinner from '../components/common/LoadingSpinner';
import { 
  formatPhoneNumber, 
  formatPostalCode, 
  formatFullAddress, 
  formatPlanName, 
  formatOptionNames 
} from '../utils/apiTransformers';
import { PAGE_TITLES, LOADING_MESSAGES } from '../utils/constants';

const UserConfirm: React.FC = () => {
  const { formData, currentStep, errors, setCurrentStep } = useFormContext();
  const { hasErrors, isValidating } = useFormValidation();
  const { 
    proceedToComplete, 
    returnToInput, 
    canSubmitForm, 
    isSubmitting, 
    submissionError 
  } = useFormSubmission();

  // API hooks for display data
  const plansApi = useGetPlans();
  const optionsApi = useGetOptions();

  // Set current step to confirm when component mounts
  useEffect(() => {
    if (currentStep !== 'confirm') {
      setCurrentStep('confirm');
    }
  }, [currentStep, setCurrentStep]);

  // Set page title
  useEffect(() => {
    document.title = PAGE_TITLES.CONFIRM;
  }, []);

  // Load reference data for display
  useEffect(() => {
    plansApi.execute();
    optionsApi.execute();
  }, []);

  // Format data for display
  const displayData = useMemo(() => {
    const plans = plansApi.data?.plans || [];
    const options = optionsApi.data?.options || [];

    return {
      name: `${formData.lastName} ${formData.firstName}`,
      nameKana: `${formData.lastNameKana} ${formData.firstNameKana}`,
      phone: formatPhoneNumber(formData.phone1, formData.phone2, formData.phone3),
      address: formatFullAddress(formData),
      email: formData.email,
      plan: formatPlanName(formData.planType, plans),
      options: formatOptionNames(formData.optionTypes, options)
    };
  }, [formData, plansApi.data, optionsApi.data]);

  const handleSubmit = async () => {
    await proceedToComplete();
  };

  const handleGoBack = () => {
    returnToInput();
  };

  return (
    <div className="user-confirm-page">
      <div className="page-container">
        <header className="page-header">
          <h1>会員登録</h1>
          <FormStepIndicator currentStep={currentStep} />
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

          {isValidating && (
            <LoadingSpinner message="入力内容を確認中..." />
          )}

          <div className="confirmation-content">
            <h2>入力内容の確認</h2>
            <p className="confirmation-instruction">
              以下の内容で会員登録を行います。内容に間違いがないかご確認ください。
            </p>

            <div className="confirmation-sections">
              <section className="confirmation-section">
                <h3>個人情報</h3>
                <div className="confirmation-grid">
                  <div className="confirmation-item">
                    <span className="confirmation-label">お名前</span>
                    <span className="confirmation-value">{displayData.name}</span>
                  </div>
                  <div className="confirmation-item">
                    <span className="confirmation-label">お名前（カナ）</span>
                    <span className="confirmation-value">{displayData.nameKana}</span>
                  </div>
                  <div className="confirmation-item">
                    <span className="confirmation-label">電話番号</span>
                    <span className="confirmation-value">{displayData.phone}</span>
                  </div>
                  <div className="confirmation-item">
                    <span className="confirmation-label">メールアドレス</span>
                    <span className="confirmation-value">{displayData.email}</span>
                  </div>
                </div>
              </section>

              <section className="confirmation-section">
                <h3>住所</h3>
                <div className="confirmation-grid">
                  <div className="confirmation-item confirmation-item--full">
                    <span className="confirmation-label">ご住所</span>
                    <span className="confirmation-value">{displayData.address}</span>
                  </div>
                </div>
              </section>

              <section className="confirmation-section">
                <h3>プラン・オプション</h3>
                <div className="confirmation-grid">
                  <div className="confirmation-item">
                    <span className="confirmation-label">選択プラン</span>
                    <span className="confirmation-value">{displayData.plan}</span>
                  </div>
                  <div className="confirmation-item">
                    <span className="confirmation-label">選択オプション</span>
                    <span className="confirmation-value">
                      {displayData.options || 'なし'}
                    </span>
                  </div>
                </div>
              </section>
            </div>

            <div className="confirmation-notice">
              <h4>ご注意</h4>
              <ul>
                <li>登録後の内容変更は、カスタマーサポートまでお問い合わせください。</li>
                <li>入力されたメールアドレスに確認メールをお送りします。</li>
                <li>登録情報は適切に管理され、第三者に提供されることはありません。</li>
              </ul>
            </div>
          </div>
        </main>

        <footer className="page-footer">
          <FormNavigation
            currentStep={currentStep}
            canGoNext={canSubmitForm && !hasErrors}
            canGoPrev={true}
            onNext={handleSubmit}
            onPrev={handleGoBack}
            isLoading={isSubmitting || isValidating}
          />
        </footer>
      </div>

      <SessionTimeoutWarning />
    </div>
  );
};

export default UserConfirm;