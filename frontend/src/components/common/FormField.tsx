// Reusable form field component
import React, { forwardRef } from 'react';
import type { FormFieldProps } from '../../types/form';

interface FormFieldComponentProps extends FormFieldProps {
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onBlur?: (e: React.FocusEvent<HTMLInputElement>) => void;
  error?: string | undefined;
  className?: string;
  autoComplete?: string;
  'data-testid'?: string;
}

const FormField = forwardRef<HTMLInputElement, FormFieldComponentProps>(
  (
    {
      name,
      label,
      required = false,
      placeholder,
      type = 'text',
      maxLength,
      pattern,
      helpText,
      disabled = false,
      value,
      onChange,
      onBlur,
      error,
      className = '',
      autoComplete,
      'data-testid': testId,
      ...props
    },
    ref
  ) => {
    const fieldId = `field-${name}`;
    const errorId = `error-${name}`;
    const helpId = `help-${name}`;

    return (
      <div className={`form-field ${className} ${error ? 'form-field--error' : ''}`}>
        <label htmlFor={fieldId} className="form-field__label">
          {label}
          {required && <span className="form-field__required">*</span>}
        </label>
        
        <input
          ref={ref}
          id={fieldId}
          name={name}
          type={type}
          value={value}
          onChange={onChange}
          onBlur={onBlur}
          placeholder={placeholder}
          maxLength={maxLength}
          pattern={pattern}
          disabled={disabled}
          required={required}
          autoComplete={autoComplete}
          className={`form-field__input ${error ? 'form-field__input--error' : ''}`}
          aria-invalid={error ? 'true' : 'false'}
          aria-describedby={
            [helpText && helpId, error && errorId].filter(Boolean).join(' ') || undefined
          }
          data-testid={testId || `${name}-input`}
          {...props}
        />
        
        {helpText && (
          <div id={helpId} className="form-field__help">
            {helpText}
          </div>
        )}
        
        {error && (
          <div id={errorId} className="form-field__error" role="alert" data-testid={`${name}-error`}>
            {error}
          </div>
        )}
      </div>
    );
  }
);

FormField.displayName = 'FormField';

export default FormField;