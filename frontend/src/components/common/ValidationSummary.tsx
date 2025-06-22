// Validation summary component for displaying errors
import React from 'react';
import { getValidationSummary, FIELD_GROUPS } from '../../validation/validationHelpers';
import type { FormValidationErrors } from '../../types/form';

interface ValidationSummaryProps {
  errors: FormValidationErrors;
  showSummary?: boolean;
  className?: string;
  'data-testid'?: string;
}

const FIELD_GROUP_LABELS = {
  PERSONAL_INFO: '個人情報',
  PHONE: '電話番号',
  ADDRESS: '住所',
  EMAIL: 'メールアドレス',
  PLAN_OPTIONS: 'プラン・オプション'
} as const;

const FIELD_LABELS: Record<string, string> = {
  lastName: '姓',
  firstName: '名',
  lastNameKana: '姓カナ',
  firstNameKana: '名カナ',
  phone1: '市外局番',
  phone2: '市内局番',
  phone3: '契約番号',
  postalCode1: '郵便番号（前3桁）',
  postalCode2: '郵便番号（後4桁）',
  prefecture: '都道府県',
  city: '市区町村',
  town: '町名',
  chome: '丁目',
  banchi: '番地',
  go: '号',
  building: '建物名',
  room: '部屋番号',
  email: 'メールアドレス',
  emailConfirm: 'メールアドレス（確認用）',
  planType: 'プラン',
  optionTypes: 'オプション'
};

const ValidationSummary: React.FC<ValidationSummaryProps> = ({
  errors,
  showSummary = true,
  className = '',
  'data-testid': testId
}) => {
  const errorFields = Object.keys(errors).filter(field => errors[field]);
  
  if (errorFields.length === 0) {
    return null;
  }

  const summary = getValidationSummary(errors);

  return (
    <div 
      className={`validation-summary ${className}`}
      role="alert"
      aria-live="polite"
      data-testid={testId || 'validation-summary'}
    >
      <div className="validation-summary__header">
        <div className="validation-summary__icon" aria-hidden="true">
          ⚠️
        </div>
        <h3 className="validation-summary__title">
          入力内容を確認してください
        </h3>
      </div>

      {showSummary && (
        <div className="validation-summary__overview">
          <p className="validation-summary__count">
            {summary.totalErrors}件のエラーがあります
          </p>
          
          {summary.criticalErrors.length > 0 && (
            <p className="validation-summary__critical">
              必須項目に{summary.criticalErrors.length}件の未入力または不正な入力があります
            </p>
          )}
        </div>
      )}

      <div className="validation-summary__errors">
        {Object.entries(FIELD_GROUPS).map(([groupName, fields]) => {
          const groupErrors = fields.filter(field => errors[field]);
          
          if (groupErrors.length === 0) return null;
          
          return (
            <div key={groupName} className="validation-summary__group">
              <h4 className="validation-summary__group-title">
                {FIELD_GROUP_LABELS[groupName as keyof typeof FIELD_GROUP_LABELS]}
              </h4>
              <ul className="validation-summary__error-list">
                {groupErrors.map(field => (
                  <li key={field} className="validation-summary__error-item">
                    <button
                      type="button"
                      className="validation-summary__error-link"
                      onClick={() => {
                        // Focus the field with error
                        const element = document.getElementById(`field-${field}`) || 
                                       document.querySelector(`[name="${field}"]`);
                        if (element) {
                          element.focus();
                          element.scrollIntoView({ behavior: 'smooth', block: 'center' });
                        }
                      }}
                      data-testid={`error-link-${field}`}
                    >
                      <span className="validation-summary__field-name">
                        {FIELD_LABELS[field] || field}
                      </span>
                      <span className="validation-summary__error-message">
                        {errors[field]}
                      </span>
                    </button>
                  </li>
                ))}
              </ul>
            </div>
          );
        })}

        {/* Handle ungrouped errors */}
        {errorFields.filter(field => 
          !Object.values(FIELD_GROUPS).flat().includes(field as any)
        ).map(field => (
          <div key={field} className="validation-summary__ungrouped-error">
            <span className="validation-summary__field-name">
              {FIELD_LABELS[field] || field}
            </span>
            <span className="validation-summary__error-message">
              {errors[field]}
            </span>
          </div>
        ))}
      </div>

      <div className="validation-summary__actions">
        <p className="validation-summary__instruction">
          エラーのある項目をクリックして修正してください
        </p>
      </div>
    </div>
  );
};

export default ValidationSummary;