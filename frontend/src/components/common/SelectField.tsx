// Reusable select field component
import React, { forwardRef } from 'react';

interface SelectOption {
  value: string;
  label: string;
  disabled?: boolean;
}

interface SelectFieldProps {
  name: string;
  label: string;
  value: string;
  onChange: (e: React.ChangeEvent<HTMLSelectElement>) => void;
  onBlur?: (e: React.FocusEvent<HTMLSelectElement>) => void;
  options: SelectOption[];
  required?: boolean;
  disabled?: boolean;
  error?: string | undefined;
  helpText?: string;
  placeholder?: string;
  className?: string;
  'data-testid'?: string;
}

const SelectField = forwardRef<HTMLSelectElement, SelectFieldProps>(
  (
    {
      name,
      label,
      value,
      onChange,
      onBlur,
      options,
      required = false,
      disabled = false,
      error,
      helpText,
      placeholder,
      className = '',
      'data-testid': testId,
      ...props
    },
    ref
  ) => {
    const fieldId = `select-${name}`;
    const errorId = `error-${name}`;
    const helpId = `help-${name}`;

    return (
      <div className={`form-field ${className} ${error ? 'form-field--error' : ''}`}>
        <label htmlFor={fieldId} className="form-field__label">
          {label}
          {required && <span className="form-field__required">*</span>}
        </label>
        
        <select
          ref={ref}
          id={fieldId}
          name={name}
          value={value}
          onChange={onChange}
          onBlur={onBlur}
          disabled={disabled}
          required={required}
          className={`form-field__select ${error ? 'form-field__select--error' : ''}`}
          aria-invalid={error ? 'true' : 'false'}
          aria-describedby={
            [helpText && helpId, error && errorId].filter(Boolean).join(' ') || undefined
          }
          data-testid={testId || `${name}-select`}
          {...props}
        >
          {placeholder && (
            <option value="" disabled>
              {placeholder}
            </option>
          )}
          {options.map((option) => (
            <option
              key={option.value}
              value={option.value}
              disabled={option.disabled}
            >
              {option.label}
            </option>
          ))}
        </select>
        
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

SelectField.displayName = 'SelectField';

export default SelectField;