// Address form component
import React, { useCallback, useState, useEffect } from 'react';
import { useFormContext } from '../../contexts/FormContext';
import { useRealtimeValidation } from '../../hooks/useRealtimeValidation';
import { useAddressSearch, useGetPrefectures } from '../../hooks/useApi';
import FormField from '../common/FormField';
import SelectField from '../common/SelectField';
import Button from '../common/Button';
import ErrorMessage from '../common/ErrorMessage';
import LoadingSpinner from '../common/LoadingSpinner';
import { isValidPostalCodeFormat } from '../../utils/apiTransformers';
import { LOADING_MESSAGES, SUCCESS_MESSAGES } from '../../utils/constants';
import type { UserFormData } from '../../types/form';

const AddressForm: React.FC = () => {
  const { formData, updateFormData, errors } = useFormContext();
  const { createFieldChangeHandler, createFieldBlurHandler } = useRealtimeValidation();
  const [searchMessage, setSearchMessage] = useState<string>('');
  
  // API hooks
  const addressSearch = useAddressSearch();
  const prefecturesApi = useGetPrefectures();

  // Load prefectures on mount
  useEffect(() => {
    prefecturesApi.execute();
  }, []);

  const handleFieldChange = useCallback((field: keyof UserFormData) => {
    return (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
      const value = e.target.value;
      updateFormData({ [field]: value });
      
      // Clear field error when user starts typing
      if (errors[field]) {
        clearFieldError(field);
      }
      
      // Clear search message when user manually changes address fields
      if (['prefecture', 'city', 'town'].includes(field) && searchMessage) {
        setSearchMessage('');
      }
    };
  }, [updateFormData, errors, clearFieldError, searchMessage]);

  const handleFieldBlur = useCallback((field: keyof UserFormData) => {
    return (e: React.FocusEvent<HTMLInputElement | HTMLSelectElement>) => {
      const value = e.target.value;
      
      // Validate field on blur
      const error = validateField(field, value, formData);
      if (error) {
        // Note: Error handling should be improved with setFieldError in context
      }
    };
  }, [formData]);

  const handleAddressSearch = useCallback(async () => {
    const { postalCode1, postalCode2 } = formData;
    
    // Validate postal code format
    if (!isValidPostalCodeFormat(postalCode1, postalCode2)) {
      setSearchMessage('郵便番号の形式が正しくありません');
      return;
    }

    try {
      setSearchMessage('');
      const result = await addressSearch.searchByPostalCode(postalCode1, postalCode2);
      
      if (result.found && result.prefecture && result.city) {
        // Update address fields
        updateFormData({
          prefecture: result.prefecture,
          city: result.city,
          town: result.town || ''
        });
        setSearchMessage(SUCCESS_MESSAGES.ADDRESS_FOUND);
      } else {
        setSearchMessage('該当する住所が見つかりませんでした');
      }
    } catch (error) {
      setSearchMessage('住所検索に失敗しました');
    }
  }, [formData, addressSearch, updateFormData]);

  // Prefecture options
  const prefectureOptions = prefecturesApi.data?.prefectures?.map(pref => ({
    value: pref.prefecture_name,
    label: pref.prefecture_name
  })) || [];

  return (
    <div className="address-form">
      <h3>住所</h3>
      
      <div className="form-row form-row--postal">
        <FormField
          name="postalCode1"
          label="郵便番号（前3桁）"
          value={formData.postalCode1}
          onChange={createFieldChangeHandler('postalCode1')}
          onBlur={createFieldBlurHandler('postalCode1')}
          error={errors.postalCode1}
          required
          type="tel"
          maxLength={3}
          placeholder="123"
          pattern="[0-9]{3}"
        />
        
        <div className="form-separator">-</div>
        
        <FormField
          name="postalCode2"
          label="郵便番号（後4桁）"
          value={formData.postalCode2}
          onChange={createFieldChangeHandler('postalCode2')}
          onBlur={createFieldBlurHandler('postalCode2')}
          error={errors.postalCode2}
          required
          type="tel"
          maxLength={4}
          placeholder="4567"
          pattern="[0-9]{4}"
        />
        
        <Button
          type="button"
          variant="outline"
          onClick={handleAddressSearch}
          disabled={
            !formData.postalCode1 ||
            !formData.postalCode2 ||
            addressSearch.isLoading
          }
          loading={addressSearch.isLoading}
          data-testid="address-search-button"
        >
          住所検索
        </Button>
      </div>
      
      {addressSearch.isLoading && (
        <LoadingSpinner message={LOADING_MESSAGES.SEARCHING_ADDRESS} />
      )}
      
      {addressSearch.error && (
        <ErrorMessage error={addressSearch.error} />
      )}
      
      {searchMessage && (
        <div className={`search-message ${searchMessage.includes('見つかりませんでした') || searchMessage.includes('失敗') ? 'search-message--error' : 'search-message--success'}`}>
          {searchMessage}
        </div>
      )}
      
      <div className="form-row">
        {prefecturesApi.isLoading ? (
          <LoadingSpinner message="都道府県を読み込み中..." />
        ) : (
          <SelectField
            name="prefecture"
            label="都道府県"
            value={formData.prefecture}
            onChange={createFieldChangeHandler('prefecture')}
            onBlur={createFieldBlurHandler('prefecture')}
            options={prefectureOptions}
            error={errors.prefecture}
            required
            placeholder="選択してください"
          />
        )}
      </div>
      
      <div className="form-row">
        <FormField
          name="city"
          label="市区町村"
          value={formData.city}
          onChange={createFieldChangeHandler('city')}
          onBlur={createFieldBlurHandler('city')}
          error={errors.city}
          required
          maxLength={50}
          placeholder="渋谷区"
        />
      </div>
      
      <div className="form-row">
        <FormField
          name="town"
          label="町名"
          value={formData.town || ''}
          onChange={createFieldChangeHandler('town')}
          onBlur={createFieldBlurHandler('town')}
          error={errors.town}
          maxLength={50}
          placeholder="神南"
          helpText="町名がある場合のみ入力してください"
        />
      </div>
      
      <div className="form-row">
        <FormField
          name="chome"
          label="丁目"
          value={formData.chome || ''}
          onChange={createFieldChangeHandler('chome')}
          onBlur={createFieldBlurHandler('chome')}
          error={errors.chome}
          maxLength={10}
          placeholder="1丁目"
          helpText="任意入力"
        />
        
        <FormField
          name="banchi"
          label="番地"
          value={formData.banchi}
          onChange={createFieldChangeHandler('banchi')}
          onBlur={createFieldBlurHandler('banchi')}
          error={errors.banchi}
          required
          maxLength={10}
          placeholder="2-3"
        />
        
        <FormField
          name="go"
          label="号"
          value={formData.go || ''}
          onChange={createFieldChangeHandler('go')}
          onBlur={createFieldBlurHandler('go')}
          error={errors.go}
          maxLength={10}
          placeholder="4号"
          helpText="任意入力"
        />
      </div>
      
      <div className="form-row">
        <FormField
          name="building"
          label="建物名"
          value={formData.building || ''}
          onChange={createFieldChangeHandler('building')}
          onBlur={createFieldBlurHandler('building')}
          error={errors.building}
          maxLength={100}
          placeholder="サンプルビル"
          helpText="任意入力"
        />
        
        <FormField
          name="room"
          label="部屋番号"
          value={formData.room || ''}
          onChange={createFieldChangeHandler('room')}
          onBlur={createFieldBlurHandler('room')}
          error={errors.room}
          maxLength={20}
          placeholder="101"
          helpText="任意入力"
        />
      </div>
    </div>
  );
};

export default AddressForm;