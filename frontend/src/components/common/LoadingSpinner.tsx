// Loading spinner component
import React from 'react';

interface LoadingSpinnerProps {
  size?: 'small' | 'medium' | 'large';
  message?: string;
  className?: string;
}

const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({
  size = 'medium',
  message,
  className = '',
}) => {
  const spinnerClasses = [
    'loading-spinner',
    `loading-spinner--${size}`,
    className
  ].filter(Boolean).join(' ');

  return (
    <div className={`loading-container ${className}`} role="status" aria-live="polite">
      <div className={spinnerClasses} aria-hidden="true">
        <div className="loading-spinner__circle"></div>
      </div>
      {message && (
        <div className="loading-message">
          {message}
        </div>
      )}
      <span className="sr-only">読み込み中...</span>
    </div>
  );
};

export default LoadingSpinner;