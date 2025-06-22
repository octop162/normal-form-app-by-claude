// User Complete Page - Complete implementation
import React, { useEffect } from 'react';
import { useFormContext } from '../contexts/FormContext';
import FormStepIndicator from '../components/layout/FormStepIndicator';
import Button from '../components/common/Button';
import { PAGE_TITLES, SUCCESS_MESSAGES } from '../utils/constants';

const UserComplete: React.FC = () => {
  const { currentStep, setCurrentStep, resetForm } = useFormContext();

  // Set current step to complete when component mounts
  useEffect(() => {
    if (currentStep !== 'complete') {
      setCurrentStep('complete');
    }
  }, [currentStep, setCurrentStep]);

  // Set page title
  useEffect(() => {
    document.title = PAGE_TITLES.COMPLETE;
  }, []);

  const handleNewRegistration = () => {
    resetForm();
    window.location.href = '/';
  };

  const handleGoHome = () => {
    // Navigate to company website or application home
    // For now, just reset the form and go to start
    resetForm();
    window.location.href = '/';
  };

  return (
    <div className="user-complete-page">
      <div className="page-container">
        <header className="page-header">
          <h1>会員登録</h1>
          <FormStepIndicator currentStep={currentStep} />
        </header>

        <main className="page-content">
          <div className="completion-content">
            <div className="completion-icon">
              <span role="img" aria-label="完了">✅</span>
            </div>
            
            <h2 className="completion-title">
              会員登録が完了しました
            </h2>
            
            <p className="completion-message">
              ご登録いただきありがとうございます。<br />
              入力されたメールアドレスに確認メールをお送りしました。
            </p>

            <div className="completion-details">
              <h3>今後の流れ</h3>
              <ol className="completion-steps">
                <li>
                  <strong>確認メールの確認</strong>
                  <p>ご登録いただいたメールアドレスに確認メールをお送りしています。メールが届かない場合は、迷惑メールフォルダもご確認ください。</p>
                </li>
                <li>
                  <strong>サービス開始のご連絡</strong>
                  <p>審査完了後、サービス開始に関するご連絡をさせていただきます。通常、1-2営業日以内にご連絡いたします。</p>
                </li>
                <li>
                  <strong>サービス利用開始</strong>
                  <p>ご利用開始の準備が整い次第、詳細なご案内をお送りいたします。</p>
                </li>
              </ol>
            </div>

            <div className="completion-notice">
              <h4>お問い合わせについて</h4>
              <p>
                ご質問やお困りのことがございましたら、カスタマーサポートまでお気軽にお問い合わせください。
              </p>
              <div className="support-info">
                <p><strong>カスタマーサポート</strong></p>
                <p>電話：0120-XXX-XXX（平日 9:00-18:00）</p>
                <p>メール：support@example.com</p>
              </div>
            </div>

            <div className="completion-actions">
              <Button
                variant="primary"
                size="large"
                onClick={handleGoHome}
                data-testid="go-home-button"
              >
                ホームに戻る
              </Button>
              
              <Button
                variant="outline"
                size="large"
                onClick={handleNewRegistration}
                data-testid="new-registration-button"
              >
                新規登録を行う
              </Button>
            </div>
          </div>
        </main>
      </div>
    </div>
  );
};

export default UserComplete;