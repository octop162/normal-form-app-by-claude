// Plan and option selection form component
import React, { useCallback, useEffect, useState } from 'react';
import { useFormContext } from '../../contexts/FormContext';
import { useRealtimeValidation } from '../../hooks/useRealtimeValidation';
import { useGetPlans, useGetOptions, useOptionAvailability } from '../../hooks/useApi';
import SelectField from '../common/SelectField';
import CheckboxGroup from '../common/CheckboxGroup';
import ErrorMessage from '../common/ErrorMessage';
import LoadingSpinner from '../common/LoadingSpinner';
import { PLAN_AVAILABLE_OPTIONS, LOADING_MESSAGES } from '../../utils/constants';
import type { UserFormData } from '../../types/form';
import type { CheckboxOption } from '../common/CheckboxGroup';

const PlanOptionForm: React.FC = () => {
  const { formData, updateFormData, errors } = useFormContext();
  const { createCheckboxChangeHandler } = useRealtimeValidation();
  const [availabilityChecked, setAvailabilityChecked] = useState(false);
  
  // API hooks
  const plansApi = useGetPlans();
  const optionsApi = useGetOptions();
  const availabilityApi = useOptionAvailability();

  // Load plans and options on mount
  useEffect(() => {
    plansApi.execute();
    optionsApi.execute();
  }, []);

  // Check availability when plan, options, or address changes
  useEffect(() => {
    const shouldCheckAvailability = 
      formData.planType &&
      formData.optionTypes.length > 0 &&
      formData.prefecture &&
      formData.city;

    if (shouldCheckAvailability) {
      availabilityApi.checkAvailability(
        formData.optionTypes,
        formData.prefecture,
        formData.city
      );
      setAvailabilityChecked(true);
    } else {
      setAvailabilityChecked(false);
    }
  }, [formData.planType, formData.optionTypes, formData.prefecture, formData.city]);

  const handlePlanChange = useCallback((e: React.ChangeEvent<HTMLSelectElement>) => {
    const planType = e.target.value;
    updateFormData({ 
      planType,
      // Clear options when plan changes
      optionTypes: []
    });
  }, [updateFormData]);

  const handleOptionsChange = createCheckboxChangeHandler('optionTypes');

  // Plan options
  const planOptions = plansApi.data?.plans
    ?.filter(plan => plan.is_active)
    ?.map(plan => ({
      value: plan.plan_type,
      label: `${plan.plan_name} (¥${plan.base_price.toLocaleString()})`
    })) || [];

  // Available options based on selected plan
  const availableOptionTypes = formData.planType 
    ? PLAN_AVAILABLE_OPTIONS[formData.planType as keyof typeof PLAN_AVAILABLE_OPTIONS] || []
    : [];

  // Option checkbox options
  const optionCheckboxOptions: CheckboxOption[] = optionsApi.data?.options
    ?.filter(option => 
      option.is_active && 
      availableOptionTypes.includes(option.option_type)
    )
    ?.map(option => {
      const inventory = availabilityApi.data.inventory[option.option_type];
      const isRestricted = availabilityApi.data.restrictions[option.option_type] === false;
      const isOutOfStock = inventory !== undefined && inventory === 0;
      const isDisabled = isRestricted || isOutOfStock;
      
      let description = option.description;
      if (option.price) {
        description += ` (¥${option.price.toLocaleString()})`;
      }
      
      if (isOutOfStock) {
        description += ' [在庫切れ]';
      } else if (isRestricted) {
        description += ' [地域制限]';
      } else if (inventory !== undefined) {
        description += ` [在庫: ${inventory}]`;
      }

      return {
        value: option.option_type,
        label: option.option_name,
        description,
        disabled: isDisabled
      };
    }) || [];

  // Check if form is ready for availability check
  const isAvailabilityCheckReady = formData.prefecture && formData.city;

  return (
    <div className="plan-option-form">
      <h3>プラン・オプション選択</h3>
      
      {plansApi.isLoading && (
        <LoadingSpinner message="プラン情報を読み込み中..." />
      )}
      
      {plansApi.error && (
        <ErrorMessage error={plansApi.error} />
      )}
      
      <div className="form-row">
        <SelectField
          name="planType"
          label="プラン"
          value={formData.planType}
          onChange={handlePlanChange}
          options={planOptions}
          error={errors['planType']}
          required
          placeholder="プランを選択してください"
          disabled={plansApi.isLoading}
        />
      </div>
      
      {formData.planType && (
        <>
          {optionsApi.isLoading && (
            <LoadingSpinner message={LOADING_MESSAGES.LOADING_OPTIONS} />
          )}
          
          {optionsApi.error && (
            <ErrorMessage error={optionsApi.error} />
          )}
          
          <div className="form-row">
            <CheckboxGroup
              name="optionTypes"
              label="オプション"
              options={optionCheckboxOptions}
              values={formData.optionTypes}
              onChange={handleOptionsChange}
              error={errors['optionTypes']}
              helpText={
                !isAvailabilityCheckReady
                  ? "住所を入力すると、在庫状況と地域制限を確認できます"
                  : availabilityApi.isLoading
                  ? "在庫状況を確認中..."
                  : "複数選択可能です"
              }
            />
          </div>
          
          {availabilityApi.isLoading && isAvailabilityCheckReady && (
            <LoadingSpinner message={LOADING_MESSAGES.CHECKING_INVENTORY} />
          )}
          
          {availabilityApi.error && (
            <ErrorMessage error={availabilityApi.error} />
          )}
          
          {availabilityChecked && !availabilityApi.isLoading && (
            <div className="availability-info">
              <h4>在庫・地域制限情報</h4>
              {Object.entries(availabilityApi.data.inventory).map(([optionType, stock]) => {
                const option = optionsApi.data?.options?.find(opt => opt.option_type === optionType);
                const isRestricted = availabilityApi.data.restrictions[optionType] === false;
                
                if (!option) return null;
                
                return (
                  <div key={optionType} className="availability-item">
                    <span className="availability-option">{option.option_name}:</span>
                    {isRestricted ? (
                      <span className="availability-status availability-status--restricted">
                        地域制限により選択できません
                      </span>
                    ) : stock === 0 ? (
                      <span className="availability-status availability-status--out-of-stock">
                        在庫切れ
                      </span>
                    ) : (
                      <span className="availability-status availability-status--available">
                        在庫あり ({stock}個)
                      </span>
                    )}
                  </div>
                );
              })}
            </div>
          )}
        </>
      )}
    </div>
  );
};

export default PlanOptionForm;