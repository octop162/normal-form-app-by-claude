// Form navigation component
import React from 'react';
import Button from '../common/Button';
import type { FormStep } from '../../types/form';

interface FormNavigationProps {
  currentStep: FormStep;
  canGoNext: boolean;
  canGoPrev: boolean;
  onNext: () => void;
  onPrev: () => void;
  isLoading?: boolean;
  className?: string;
}

const STEP_LABELS = {
  input: {
    next: '確認画面へ',
    prev: ''
  },
  confirm: {
    next: '申し込む',
    prev: '入力画面に戻る'
  },
  complete: {
    next: '',
    prev: ''
  }
} as const;

const FormNavigation: React.FC<FormNavigationProps> = ({
  currentStep,
  canGoNext,
  canGoPrev,
  onNext,
  onPrev,
  isLoading = false,
  className = ''
}) => {
  const labels = STEP_LABELS[currentStep];

  // Don't show navigation on complete step
  if (currentStep === 'complete') {
    return null;
  }

  return (
    <div className={`form-navigation ${className}`}>
      <div className="form-navigation__buttons">
        {canGoPrev && labels.prev && (
          <Button
            type="button"
            variant="secondary"
            onClick={onPrev}
            disabled={isLoading}
            data-testid="prev-button"
          >
            {labels.prev}
          </Button>
        )}
        
        <div className="form-navigation__spacer" />
        
        {canGoNext && labels.next && (
          <Button
            type="button"
            variant="primary"
            onClick={onNext}
            disabled={!canGoNext}
            loading={isLoading}
            data-testid="next-button"
          >
            {labels.next}
          </Button>
        )}
      </div>
    </div>
  );
};

export default FormNavigation;