// Form-specific types for the membership registration system

// Main form data structure
export interface UserFormData {
  // Personal information
  lastName: string;
  firstName: string;
  lastNameKana: string;
  firstNameKana: string;
  
  // Phone number (3 parts)
  phone1: string;
  phone2: string;
  phone3: string;
  
  // Postal code (2 parts)
  postalCode1: string;
  postalCode2: string;
  
  // Address
  prefecture: string;
  city: string;
  town?: string;
  chome?: string;
  banchi: string;
  go?: string;
  building?: string;
  room?: string;
  
  // Email
  email: string;
  emailConfirm: string;
  
  // Plan and options
  planType: string;
  optionTypes: string[];
}

// Form step types
export type FormStep = 'input' | 'confirm' | 'complete';

// Form validation error types
export interface FormValidationError {
  field: string;
  message: string;
}

export interface FormValidationErrors {
  [key: string]: string | undefined;
}

// Form context types
export interface FormContextType {
  formData: UserFormData;
  updateFormData: (data: Partial<UserFormData>) => void;
  resetForm: () => void;
  currentStep: FormStep;
  setCurrentStep: (step: FormStep) => void;
  isLoading: boolean;
  setIsLoading: (loading: boolean) => void;
  errors: FormValidationErrors;
  setErrors: (errors: FormValidationErrors) => void;
  sessionId?: string;
  setSessionId: (id: string | undefined) => void;
  // Additional utility methods
  clearFieldError: (fieldName: string) => void;
  hasFormData: () => boolean;
  getSessionTimeoutWarning: () => { show: boolean; remainingMinutes: number };
  extendSession: () => void;
}

// Option availability types
export interface OptionAvailability {
  optionType: string;
  isAvailable: boolean;
  stock?: number;
  isRegionRestricted?: boolean;
  reason?: string;
}

// Plan types
export interface PlanInfo {
  planType: string;
  planName: string;
  description: string;
  basePrice: number;
  availableOptions: string[];
}

// Form field types for validation
export interface FormFieldProps {
  name: string;
  label: string;
  required?: boolean;
  placeholder?: string;
  type?: 'text' | 'email' | 'tel' | 'number';
  maxLength?: number;
  pattern?: string;
  helpText?: string;
  disabled?: boolean;
}

// Form step navigation
export interface FormStepNavigation {
  canGoNext: boolean;
  canGoPrev: boolean;
  nextLabel: string;
  prevLabel: string;
  onNext: () => void;
  onPrev: () => void;
}

// Session management
export interface SessionData {
  sessionId: string;
  expiresAt: string;
  isValid: boolean;
}

export interface SessionTimeoutWarning {
  show: boolean;
  remainingMinutes: number;
  onExtend: () => void;
  onLogout: () => void;
}