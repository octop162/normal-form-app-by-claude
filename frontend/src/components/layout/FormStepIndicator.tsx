// Form step indicator component
import React from 'react';
import type { FormStep } from '../../types/form';

interface FormStepIndicatorProps {
  currentStep: FormStep;
  className?: string;
}

const STEPS = [
  { key: 'input', label: '入力', number: 1 },
  { key: 'confirm', label: '確認', number: 2 },
  { key: 'complete', label: '完了', number: 3 }
] as const;

const FormStepIndicator: React.FC<FormStepIndicatorProps> = ({
  currentStep,
  className = ''
}) => {
  const currentStepNumber = STEPS.find(step => step.key === currentStep)?.number || 1;

  return (
    <div className={`step-indicator ${className}`} role="progressbar" aria-valuemin={1} aria-valuemax={3} aria-valuenow={currentStepNumber}>
      <div className="step-indicator__container">
        {STEPS.map((step) => {
          const isActive = step.key === currentStep;
          const isCompleted = step.number < currentStepNumber;
          const stepClass = [
            'step-indicator__step',
            isActive && 'step-indicator__step--active',
            isCompleted && 'step-indicator__step--completed'
          ].filter(Boolean).join(' ');

          return (
            <div key={step.key} className={stepClass}>
              <div className="step-indicator__step-circle">
                {isCompleted ? (
                  <span className="step-indicator__check" aria-hidden="true">✓</span>
                ) : (
                  <span className="step-indicator__number">{step.number}</span>
                )}
              </div>
              <div className="step-indicator__step-label">
                {step.label}
              </div>
            </div>
          );
        })}
      </div>
      <div className="step-indicator__progress-bar">
        <div 
          className="step-indicator__progress-fill"
          style={{ width: `${((currentStepNumber - 1) / (STEPS.length - 1)) * 100}%` }}
        />
      </div>
    </div>
  );
};

export default FormStepIndicator;