// API data transformers for frontend-backend communication
import type { UserFormData } from '../types/form';
import type { UserCreateRequest, UserValidateRequest } from '../types/api';

// Transform frontend form data to backend API format
export const transformFormDataToApiRequest = (formData: UserFormData): UserCreateRequest => {
  return {
    last_name: formData.lastName,
    first_name: formData.firstName,
    last_name_kana: formData.lastNameKana,
    first_name_kana: formData.firstNameKana,
    phone1: formData.phone1,
    phone2: formData.phone2,
    phone3: formData.phone3,
    postal_code1: formData.postalCode1,
    postal_code2: formData.postalCode2,
    prefecture: formData.prefecture,
    city: formData.city,
    town: formData.town || undefined,
    chome: formData.chome || undefined,
    banchi: formData.banchi,
    go: formData.go || undefined,
    building: formData.building || undefined,
    room: formData.room || undefined,
    email: formData.email,
    email_confirm: formData.emailConfirm,
    plan_type: formData.planType,
    option_types: formData.optionTypes
  };
};

// Transform frontend form data to validation request format
export const transformFormDataToValidateRequest = (formData: UserFormData): UserValidateRequest => {
  return transformFormDataToApiRequest(formData) as UserValidateRequest;
};

// Transform backend API response to frontend form data
export const transformApiResponseToFormData = (apiData: UserCreateRequest): UserFormData => {
  return {
    lastName: apiData.last_name,
    firstName: apiData.first_name,
    lastNameKana: apiData.last_name_kana,
    firstNameKana: apiData.first_name_kana,
    phone1: apiData.phone1,
    phone2: apiData.phone2,
    phone3: apiData.phone3,
    postalCode1: apiData.postal_code1,
    postalCode2: apiData.postal_code2,
    prefecture: apiData.prefecture,
    city: apiData.city,
    town: apiData.town || '',
    chome: apiData.chome || '',
    banchi: apiData.banchi,
    go: apiData.go || '',
    building: apiData.building || '',
    room: apiData.room || '',
    email: apiData.email,
    emailConfirm: apiData.email_confirm,
    planType: apiData.plan_type,
    optionTypes: apiData.option_types
  };
};

// Format phone number for display
export const formatPhoneNumber = (phone1: string, phone2: string, phone3: string): string => {
  return `${phone1}-${phone2}-${phone3}`;
};

// Format postal code for display
export const formatPostalCode = (postalCode1: string, postalCode2: string): string => {
  return `〒${postalCode1}-${postalCode2}`;
};

// Format address for display
export const formatAddress = (formData: UserFormData): string => {
  const parts = [
    formData.prefecture,
    formData.city,
    formData.town,
    formData.chome,
    formData.banchi,
    formData.go,
    formData.building,
    formData.room
  ].filter(part => part && part.trim() !== '');
  
  return parts.join('');
};

// Format full address with postal code for display
export const formatFullAddress = (formData: UserFormData): string => {
  const postalCode = formatPostalCode(formData.postalCode1, formData.postalCode2);
  const address = formatAddress(formData);
  return `${postalCode} ${address}`;
};

// Validate postal code format
export const isValidPostalCodeFormat = (postalCode1: string, postalCode2: string): boolean => {
  return /^\d{3}$/.test(postalCode1) && /^\d{4}$/.test(postalCode2);
};

// Format option names for display
export const formatOptionNames = (optionTypes: string[], availableOptions: Array<{ option_type: string; option_name: string }>): string => {
  const optionNames = optionTypes.map(type => {
    const option = availableOptions.find(opt => opt.option_type === type);
    return option ? option.option_name : type;
  });
  
  return optionNames.join('、');
};

// Format plan name for display
export const formatPlanName = (planType: string, availablePlans: Array<{ plan_type: string; plan_name: string }>): string => {
  const plan = availablePlans.find(p => p.plan_type === planType);
  return plan ? plan.plan_name : planType;
};