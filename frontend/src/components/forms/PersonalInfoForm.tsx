// Personal information form component
import React from 'react';
import { useFormContext } from '../../contexts/FormContext';
import { useRealtimeValidation } from '../../hooks/useRealtimeValidation';
import FormField from '../common/FormField';
import type { UserFormData } from '../../types/form';

const PersonalInfoForm: React.FC = () => {
  const { formData, errors } = useFormContext();
  const { createFieldChangeHandler, createFieldBlurHandler } = useRealtimeValidation();

  return (
    <div className="personal-info-form">
      <h3>個人情報</h3>
      
      <div className="form-row">
        <FormField
          name="lastName"
          label="姓"
          value={formData.lastName}
          onChange={createFieldChangeHandler('lastName')}
          onBlur={createFieldBlurHandler('lastName')}
          error={errors.lastName}
          required
          maxLength={15}
          placeholder="山田"
          autoComplete="family-name"
        />
        
        <FormField
          name="firstName"
          label="名"
          value={formData.firstName}
          onChange={createFieldChangeHandler('firstName')}
          onBlur={createFieldBlurHandler('firstName')}
          error={errors.firstName}
          required
          maxLength={15}
          placeholder="太郎"
          autoComplete="given-name"
        />
      </div>
      
      <div className="form-row">
        <FormField
          name="lastNameKana"
          label="姓カナ"
          value={formData.lastNameKana}
          onChange={createFieldChangeHandler('lastNameKana')}
          onBlur={createFieldBlurHandler('lastNameKana')}
          error={errors.lastNameKana}
          required
          maxLength={15}
          placeholder="ヤマダ"
          helpText="全角カタカナで入力してください"
          autoComplete="family-name"
        />
        
        <FormField
          name="firstNameKana"
          label="名カナ"
          value={formData.firstNameKana}
          onChange={createFieldChangeHandler('firstNameKana')}
          onBlur={createFieldBlurHandler('firstNameKana')}
          error={errors.firstNameKana}
          required
          maxLength={15}
          placeholder="タロウ"
          helpText="全角カタカナで入力してください"
          autoComplete="given-name"
        />
      </div>
      
      <div className="form-row form-row--phone">
        <FormField
          name="phone1"
          label="電話番号（市外局番）"
          value={formData.phone1}
          onChange={createFieldChangeHandler('phone1')}
          onBlur={createFieldBlurHandler('phone1')}
          error={errors.phone1}
          required
          type="tel"
          maxLength={5}
          placeholder="03"
          helpText="市外局番（2〜5桁）"
          autoComplete="tel-area-code"
        />
        
        <div className="form-separator">-</div>
        
        <FormField
          name="phone2"
          label="市内局番"
          value={formData.phone2}
          onChange={createFieldChangeHandler('phone2')}
          onBlur={createFieldBlurHandler('phone2')}
          error={errors.phone2}
          required
          type="tel"
          maxLength={4}
          placeholder="1234"
          helpText="市内局番（1〜4桁）"
          autoComplete="tel-local-prefix"
        />
        
        <div className="form-separator">-</div>
        
        <FormField
          name="phone3"
          label="契約番号"
          value={formData.phone3}
          onChange={createFieldChangeHandler('phone3')}
          onBlur={createFieldBlurHandler('phone3')}
          error={errors.phone3}
          required
          type="tel"
          maxLength={4}
          placeholder="5678"
          helpText="契約番号（4桁）"
          autoComplete="tel-local-suffix"
        />
      </div>
      
      <div className="form-row">
        <FormField
          name="email"
          label="メールアドレス"
          value={formData.email}
          onChange={createFieldChangeHandler('email')}
          onBlur={createFieldBlurHandler('email')}
          error={errors.email}
          required
          type="email"
          maxLength={256}
          placeholder="example@email.com"
          autoComplete="email"
        />
      </div>
      
      <div className="form-row">
        <FormField
          name="emailConfirm"
          label="メールアドレス（確認用）"
          value={formData.emailConfirm}
          onChange={createFieldChangeHandler('emailConfirm')}
          onBlur={createFieldBlurHandler('emailConfirm')}
          error={errors.emailConfirm}
          required
          type="email"
          maxLength={256}
          placeholder="example@email.com"
          helpText="確認のため、もう一度入力してください"
          autoComplete="email"
        />
      </div>
    </div>
  );
};

export default PersonalInfoForm;