// Enhanced validation schemas using Zod
import { z } from 'zod';
import { securityService } from '../services/securityService';

// Enhanced validation patterns
const JAPANESE_NAME_PATTERN = /^[ひらがなカタカナ漢字ａ-ｚＡ-Ｚ０-９\s\-ー]+$/;
const KATAKANA_PATTERN = /^[ァ-ヶー\s]+$/;
const PHONE_PATTERNS = {
  phone1: /^(0\d{1,4})$/, // Area code: 2-5 digits starting with 0
  phone2: /^\d{1,4}$/, // Local exchange: 1-4 digits
  phone3: /^\d{4}$/, // Number: 4 digits
};
const POSTAL_CODE_PATTERN = /^\d{3}$/;
const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

// Security validation functions
const validateNoSuspiciousContent = (value: string): boolean => {
  if (!value) return true;
  return !securityService.containsSuspiciousPatterns(value);
};

const validateJapaneseName = (value: string): boolean => {
  if (!value) return false;
  const trimmed = value.trim();
  return trimmed.length > 0 && JAPANESE_NAME_PATTERN.test(trimmed);
};

const validateKatakanaOnly = (value: string): boolean => {
  if (!value) return false;
  const trimmed = value.trim();
  return trimmed.length > 0 && KATAKANA_PATTERN.test(trimmed);
};

const validatePhoneNumberFormat = (phone1: string, phone2: string, phone3: string): boolean => {
  return securityService.isValidPhoneNumber(phone1, phone2, phone3);
};

const validatePostalCodeFormat = (code1: string, code2: string): boolean => {
  return securityService.isValidPostalCode(code1, code2);
};

const validateEmailFormat = (email: string): boolean => {
  return securityService.isValidEmail(email);
};

// Base validation schemas
export const userFormSchema = z.object({
  // Personal information with enhanced validation
  lastName: z.string()
    .min(1, '姓は必須です')
    .max(15, '姓は15文字以内で入力してください')
    .refine(validateJapaneseName, '有効な日本語の名前を入力してください')
    .refine(validateNoSuspiciousContent, '使用できない文字が含まれています'),
  
  firstName: z.string()
    .min(1, '名は必須です')
    .max(15, '名は15文字以内で入力してください')
    .refine(validateJapaneseName, '有効な日本語の名前を入力してください')
    .refine(validateNoSuspiciousContent, '使用できない文字が含まれています'),
  
  lastNameKana: z.string()
    .min(1, '姓カナは必須です')
    .max(15, '姓カナは15文字以内で入力してください')
    .refine(validateKatakanaOnly, '姓カナは全角カタカナで入力してください')
    .refine(validateNoSuspiciousContent, '使用できない文字が含まれています'),
  
  firstNameKana: z.string()
    .min(1, '名カナは必須です')
    .max(15, '名カナは15文字以内で入力してください')
    .refine(validateKatakanaOnly, '名カナは全角カタカナで入力してください')
    .refine(validateNoSuspiciousContent, '使用できない文字が含まれています'),
  
  // Phone number (3 parts)
  phone1: z.string()
    .min(1, '市外局番は必須です')
    .regex(PHONE_PATTERNS.phone1, '市外局番の形式が正しくありません'),
  
  phone2: z.string()
    .min(1, '市内局番は必須です')
    .regex(PHONE_PATTERNS.phone2, '市内局番の形式が正しくありません'),
  
  phone3: z.string()
    .min(1, '契約番号は必須です')
    .regex(PHONE_PATTERNS.phone3, '契約番号は4桁で入力してください'),
  
  // Postal code (2 parts)
  postalCode1: z.string()
    .min(1, '郵便番号（前3桁）は必須です')
    .regex(POSTAL_CODE_PATTERN, '郵便番号（前3桁）は3桁の数字で入力してください'),
  
  postalCode2: z.string()
    .min(1, '郵便番号（後4桁）は必須です')
    .regex(/^\d{4}$/, '郵便番号（後4桁）は4桁の数字で入力してください'),
  
  // Address
  prefecture: z.string()
    .min(1, '都道府県は必須です'),
  
  city: z.string()
    .min(1, '市区町村は必須です')
    .max(50, '市区町村は50文字以内で入力してください'),
  
  town: z.string()
    .max(50, '町名は50文字以内で入力してください')
    .optional(),
  
  chome: z.string()
    .max(10, '丁目は10文字以内で入力してください')
    .optional(),
  
  banchi: z.string()
    .min(1, '番地は必須です')
    .max(10, '番地は10文字以内で入力してください'),
  
  go: z.string()
    .max(10, '号は10文字以内で入力してください')
    .optional(),
  
  building: z.string()
    .max(100, '建物名は100文字以内で入力してください')
    .optional(),
  
  room: z.string()
    .max(20, '部屋番号は20文字以内で入力してください')
    .optional(),
  
  // Email with enhanced validation
  email: z.string()
    .min(1, 'メールアドレスは必須です')
    .max(256, 'メールアドレスは256文字以内で入力してください')
    .refine(validateEmailFormat, 'メールアドレスの形式が正しくありません')
    .refine(validateNoSuspiciousContent, '使用できない文字が含まれています'),
  
  emailConfirm: z.string()
    .min(1, 'メールアドレス（確認用）は必須です')
    .refine(validateEmailFormat, 'メールアドレスの形式が正しくありません')
    .refine(validateNoSuspiciousContent, '使用できない文字が含まれています'),
  
  // Plan and options
  planType: z.string()
    .min(1, 'プランを選択してください'),
  
  optionTypes: z.array(z.string())
    .default([])
}).refine((data) => {
  // Email confirmation validation
  return data.email === data.emailConfirm;
}, {
  message: 'メールアドレスが一致しません',
  path: ['emailConfirm']
}).refine((data) => {
  // Enhanced phone number validation
  return validatePhoneNumberFormat(data.phone1, data.phone2, data.phone3);
}, {
  message: '電話番号の形式が正しくありません。フリーダイヤルは使用できません。',
  path: ['phone3']
}).refine((data) => {
  // Enhanced postal code validation
  return validatePostalCodeFormat(data.postalCode1, data.postalCode2);
}, {
  message: '郵便番号の形式が正しくありません',
  path: ['postalCode2']
}).refine((data) => {
  // Validate options for selected plan
  if (!data.optionTypes || data.optionTypes.length === 0) {
    return true; // No options selected is valid
  }
  
  const validOptions = {
    'A': ['AA', 'AB'],
    'B': ['BB', 'AB']
  };
  
  const allowedOptions = validOptions[data.planType as keyof typeof validOptions];
  return allowedOptions && data.optionTypes.every(option => allowedOptions.includes(option));
}, {
  message: '選択されたオプションは指定されたプランでは利用できません',
  path: ['optionTypes']
});

// Partial validation for real-time validation
export const createPartialSchema = (fields: (keyof z.infer<typeof userFormSchema>)[]) => {
  const partialSchema = userFormSchema.partial();
  return partialSchema.pick(
    fields.reduce((acc, field) => {
      acc[field] = true;
      return acc;
    }, {} as Record<keyof z.infer<typeof userFormSchema>, true>)
  );
};

// Individual field validation schemas for real-time validation
export const fieldSchemas = {
  lastName: userFormSchema.shape.lastName,
  firstName: userFormSchema.shape.firstName,
  lastNameKana: userFormSchema.shape.lastNameKana,
  firstNameKana: userFormSchema.shape.firstNameKana,
  phone1: userFormSchema.shape.phone1,
  phone2: userFormSchema.shape.phone2,
  phone3: userFormSchema.shape.phone3,
  postalCode1: userFormSchema.shape.postalCode1,
  postalCode2: userFormSchema.shape.postalCode2,
  prefecture: userFormSchema.shape.prefecture,
  city: userFormSchema.shape.city,
  town: userFormSchema.shape.town,
  chome: userFormSchema.shape.chome,
  banchi: userFormSchema.shape.banchi,
  go: userFormSchema.shape.go,
  building: userFormSchema.shape.building,
  room: userFormSchema.shape.room,
  email: userFormSchema.shape.email,
  emailConfirm: userFormSchema.shape.emailConfirm,
  planType: userFormSchema.shape.planType,
  optionTypes: userFormSchema.shape.optionTypes,
};

// Type inference
export type UserFormData = z.infer<typeof userFormSchema>;
export type ValidationErrors = Record<string, string>;

// Validation helper functions
export const validateField = (
  fieldName: keyof UserFormData,
  value: any,
  allData?: Partial<UserFormData>
): string | null => {
  try {
    // For email confirmation, we need the full data to compare
    if (fieldName === 'emailConfirm' && allData) {
      const result = userFormSchema.pick({ email: true, emailConfirm: true }).parse({
        email: allData.email,
        emailConfirm: value
      });
      return null;
    }
    
    // For other fields, validate individually
    const schema = fieldSchemas[fieldName];
    schema.parse(value);
    return null;
  } catch (error) {
    if (error instanceof z.ZodError) {
      return error.errors[0]?.message || 'バリデーションエラー';
    }
    return 'バリデーションエラー';
  }
};

export const validateForm = (data: UserFormData): ValidationErrors => {
  try {
    userFormSchema.parse(data);
    return {};
  } catch (error) {
    if (error instanceof z.ZodError) {
      const errors: ValidationErrors = {};
      error.errors.forEach((err) => {
        const path = err.path.join('.');
        errors[path] = err.message;
      });
      return errors;
    }
    return { general: 'バリデーションエラーが発生しました' };
  }
};