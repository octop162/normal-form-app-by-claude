// Checkbox group component for options selection
import React from 'react';

export interface CheckboxOption {
  value: string;
  label: string;
  disabled?: boolean;
  description?: string;
}

interface CheckboxGroupProps {
  name: string;
  label: string;
  options: CheckboxOption[];
  values: string[];
  onChange: (values: string[]) => void;
  error?: string | undefined;
  helpText?: string;
  required?: boolean;
  className?: string;
  'data-testid'?: string;
}

const CheckboxGroup: React.FC<CheckboxGroupProps> = ({
  name,
  label,
  options,
  values,
  onChange,
  error,
  helpText,
  required = false,
  className = '',
  'data-testid': testId,
}) => {
  const groupId = `checkbox-group-${name}`;
  const errorId = `error-${name}`;
  const helpId = `help-${name}`;

  const handleCheckboxChange = (optionValue: string, checked: boolean) => {
    if (checked) {
      // Add option to selected values
      onChange([...values, optionValue]);
    } else {
      // Remove option from selected values
      onChange(values.filter(value => value !== optionValue));
    }
  };

  return (
    <fieldset
      className={`checkbox-group ${className} ${error ? 'checkbox-group--error' : ''}`}
      aria-invalid={error ? 'true' : 'false'}
      aria-describedby={
        [helpText && helpId, error && errorId].filter(Boolean).join(' ') || undefined
      }
      data-testid={testId || `${name}-checkbox-group`}
    >
      <legend className="checkbox-group__legend">
        {label}
        {required && <span className="checkbox-group__required">*</span>}
      </legend>
      
      <div className="checkbox-group__options" role="group" aria-labelledby={groupId}>
        {options.map((option) => {
          const isChecked = values.includes(option.value);
          const checkboxId = `${name}-${option.value}`;
          
          return (
            <div key={option.value} className="checkbox-group__option">
              <input
                type="checkbox"
                id={checkboxId}
                name={`${name}[]`}
                value={option.value}
                checked={isChecked}
                disabled={option.disabled}
                onChange={(e) => handleCheckboxChange(option.value, e.target.checked)}
                className="checkbox-group__input"
                data-testid={`${name}-${option.value}`}
              />
              <label htmlFor={checkboxId} className="checkbox-group__label">
                <span className="checkbox-group__label-text">
                  {option.label}
                </span>
                {option.description && (
                  <span className="checkbox-group__description">
                    {option.description}
                  </span>
                )}
              </label>
            </div>
          );
        })}
      </div>
      
      {helpText && (
        <div id={helpId} className="checkbox-group__help">
          {helpText}
        </div>
      )}
      
      {error && (
        <div id={errorId} className="checkbox-group__error" role="alert" data-testid={`${name}-error`}>
          {error}
        </div>
      )}
    </fieldset>
  );
};

export default CheckboxGroup;