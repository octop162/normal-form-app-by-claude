// Validation helper utilities
import type { UserFormData, FormValidationErrors } from '../types/form';
import { VALIDATION_PATTERNS, ERROR_MESSAGES } from '../utils/constants';

// Phone number validation helpers
export const validatePhoneNumber = (phone1: string, phone2: string, phone3: string): string | null => {
  // Check individual parts format
  if (!VALIDATION_PATTERNS.PHONE.PHONE1.test(phone1)) {
    return ERROR_MESSAGES.PHONE_FORMAT;
  }
  
  if (!VALIDATION_PATTERNS.PHONE.PHONE2.test(phone2)) {
    return ERROR_MESSAGES.PHONE_FORMAT;
  }
  
  if (!VALIDATION_PATTERNS.PHONE.PHONE3.test(phone3)) {
    return ERROR_MESSAGES.PHONE_FORMAT;
  }

  // Check toll-free numbers
  const isTollFree = VALIDATION_PATTERNS.PHONE.TOLL_FREE.some(pattern => 
    phone1.startsWith(pattern)
  );
  
  if (isTollFree) {
    return ERROR_MESSAGES.TOLL_FREE_NOT_ALLOWED;
  }

  // Check mobile number format for 11-digit numbers
  const fullPhone = phone1 + phone2 + phone3;
  if (fullPhone.length === 11 && !VALIDATION_PATTERNS.PHONE.MOBILE_PREFIX.test(fullPhone)) {
    return ERROR_MESSAGES.MOBILE_PHONE_FORMAT;
  }

  return null;
};

// Email validation helpers
export const validateEmailFormat = (email: string): string | null => {
  if (!VALIDATION_PATTERNS.EMAIL.test(email)) {
    return ERROR_MESSAGES.EMAIL_FORMAT;
  }
  return null;
};

export const validateEmailConfirmation = (email: string, emailConfirm: string): string | null => {
  if (email !== emailConfirm) {
    return ERROR_MESSAGES.EMAIL_MISMATCH;
  }
  return null;
};

// Katakana validation helper
export const validateKatakana = (value: string): string | null => {
  if (!VALIDATION_PATTERNS.KATAKANA.test(value)) {
    return ERROR_MESSAGES.KATAKANA_ONLY;
  }
  return null;
};

// Postal code validation helper
export const validatePostalCode = (postalCode1: string, postalCode2: string): string | null => {
  if (!VALIDATION_PATTERNS.POSTAL_CODE.PART1.test(postalCode1)) {
    return ERROR_MESSAGES.POSTAL_CODE_FORMAT;
  }
  
  if (!VALIDATION_PATTERNS.POSTAL_CODE.PART2.test(postalCode2)) {
    return ERROR_MESSAGES.POSTAL_CODE_FORMAT;
  }
  
  return null;
};

// Required field validation
export const validateRequired = (value: any, fieldName: string): string | null => {
  if (value === null || value === undefined || value === '') {
    return ERROR_MESSAGES.REQUIRED;
  }
  
  if (Array.isArray(value) && value.length === 0) {
    // For optional arrays like options, return null
    return null;
  }
  
  if (typeof value === 'string' && value.trim() === '') {
    return ERROR_MESSAGES.REQUIRED;
  }
  
  return null;
};

// Length validation
export const validateMaxLength = (value: string, maxLength: number): string | null => {
  if (value && value.length > maxLength) {
    return ERROR_MESSAGES.MAX_LENGTH(maxLength);
  }
  return null;
};

// Form-level validation helpers
export const getFieldsWithErrors = (errors: FormValidationErrors): string[] => {
  return Object.keys(errors).filter(field => errors[field]);
};

export const hasFieldErrors = (errors: FormValidationErrors): boolean => {
  return Object.values(errors).some(error => error && error.trim() !== '');
};

export const getFirstErrorField = (errors: FormValidationErrors): string | null => {
  const fieldsWithErrors = getFieldsWithErrors(errors);
  return fieldsWithErrors.length > 0 ? fieldsWithErrors[0] : null;
};

// Field grouping for validation
export const FIELD_GROUPS = {
  PERSONAL_INFO: ['lastName', 'firstName', 'lastNameKana', 'firstNameKana'] as const,
  PHONE: ['phone1', 'phone2', 'phone3'] as const,
  ADDRESS: ['postalCode1', 'postalCode2', 'prefecture', 'city', 'town', 'chome', 'banchi', 'go', 'building', 'room'] as const,
  EMAIL: ['email', 'emailConfirm'] as const,
  PLAN_OPTIONS: ['planType', 'optionTypes'] as const
} as const;

export const validateFieldGroup = (
  groupName: keyof typeof FIELD_GROUPS,
  formData: UserFormData
): FormValidationErrors => {
  const errors: FormValidationErrors = {};
  const fields = FIELD_GROUPS[groupName];
  
  // This would integrate with the main validation schema
  // For now, return empty errors as the main validation handles this
  return errors;
};

// Validation summary helpers
export const getValidationSummary = (errors: FormValidationErrors): {
  totalErrors: number;
  errorsByGroup: Record<string, number>;
  criticalErrors: string[];
} => {
  const totalErrors = getFieldsWithErrors(errors).length;
  const errorsByGroup: Record<string, number> = {};
  const criticalErrors: string[] = [];

  // Count errors by field group
  Object.entries(FIELD_GROUPS).forEach(([groupName, fields]) => {
    const groupErrors = fields.filter(field => errors[field]);
    errorsByGroup[groupName] = groupErrors.length;
  });

  // Identify critical errors (required fields)
  const requiredFields: (keyof UserFormData)[] = [
    'lastName', 'firstName', 'lastNameKana', 'firstNameKana',
    'phone1', 'phone2', 'phone3',
    'postalCode1', 'postalCode2',
    'prefecture', 'city', 'banchi',
    'email', 'emailConfirm',
    'planType'
  ];

  requiredFields.forEach(field => {
    if (errors[field]) {
      criticalErrors.push(field);
    }
  });

  return {
    totalErrors,
    errorsByGroup,
    criticalErrors
  };
};

// Cross-field validation helpers
export const validateCrossFieldDependencies = (formData: UserFormData): FormValidationErrors => {
  const errors: FormValidationErrors = {};

  // Email confirmation validation
  if (formData.email && formData.emailConfirm) {
    const emailError = validateEmailConfirmation(formData.email, formData.emailConfirm);
    if (emailError) {
      errors.emailConfirm = emailError;
    }
  }

  // Phone number cross-validation
  if (formData.phone1 && formData.phone2 && formData.phone3) {
    const phoneError = validatePhoneNumber(formData.phone1, formData.phone2, formData.phone3);
    if (phoneError) {
      errors.phone1 = phoneError; // Display error on first phone field
    }
  }

  // Postal code cross-validation
  if (formData.postalCode1 && formData.postalCode2) {
    const postalError = validatePostalCode(formData.postalCode1, formData.postalCode2);
    if (postalError) {
      errors.postalCode1 = postalError; // Display error on first postal field
    }
  }

  return errors;
};