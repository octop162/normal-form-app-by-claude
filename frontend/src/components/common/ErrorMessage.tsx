// Error message component
import React from 'react';
import type { ApiError } from '../../types/api';

interface ErrorMessageProps {
  error: ApiError | Error | string | null;
  className?: string;
  'data-testid'?: string;
}

const ErrorMessage: React.FC<ErrorMessageProps> = ({
  error,
  className = '',
  'data-testid': testId,
}) => {
  if (!error) return null;

  const getErrorMessage = (error: ApiError | Error | string): string => {
    if (typeof error === 'string') {
      return error;
    }
    
    if ('message' in error) {
      return error.message;
    }
    
    return 'エラーが発生しました';
  };

  const getErrorCode = (error: ApiError | Error | string): string | undefined => {
    if (typeof error === 'object' && 'code' in error) {
      return error.code;
    }
    return undefined;
  };

  const message = getErrorMessage(error);
  const code = getErrorCode(error);

  return (
    <div
      className={`error-message ${className}`}
      role="alert"
      data-testid={testId || 'error-message'}
    >
      <div className="error-message__icon" aria-hidden="true">
        ⚠️
      </div>
      <div className="error-message__content">
        <div className="error-message__text">
          {message}
        </div>
        {code && (
          <div className="error-message__code">
            エラーコード: {code}
          </div>
        )}
      </div>
    </div>
  );
};

export default ErrorMessage;