import type { ApiError } from '../types/api';

// Error codes from backend
export const ERROR_CODES = {
  // Generic errors
  INTERNAL_SERVER_ERROR: 'INTERNAL_SERVER_ERROR',
  BAD_REQUEST: 'BAD_REQUEST',
  NOT_FOUND: 'NOT_FOUND',
  UNAUTHORIZED: 'UNAUTHORIZED',
  FORBIDDEN: 'FORBIDDEN',
  CONFLICT: 'CONFLICT',
  TOO_MANY_REQUESTS: 'TOO_MANY_REQUESTS',

  // Validation errors
  VALIDATION_FAILED: 'VALIDATION_FAILED',
  REQUIRED_FIELD_MISSING: 'REQUIRED_FIELD_MISSING',
  INVALID_FORMAT: 'INVALID_FORMAT',
  INVALID_EMAIL: 'INVALID_EMAIL',
  INVALID_PHONE_NUMBER: 'INVALID_PHONE_NUMBER',
  INVALID_POSTAL_CODE: 'INVALID_POSTAL_CODE',
  EMAIL_CONFIRMATION_FAILED: 'EMAIL_CONFIRMATION_FAILED',

  // Business logic errors
  USER_ALREADY_EXISTS: 'USER_ALREADY_EXISTS',
  SESSION_EXPIRED: 'SESSION_EXPIRED',
  INVENTORY_NOT_AVAILABLE: 'INVENTORY_NOT_AVAILABLE',
  REGION_NOT_SUPPORTED: 'REGION_NOT_SUPPORTED',
  OPTION_NOT_AVAILABLE: 'OPTION_NOT_AVAILABLE',
  ADDRESS_NOT_FOUND: 'ADDRESS_NOT_FOUND',

  // External API errors
  EXTERNAL_API_ERROR: 'EXTERNAL_API_ERROR',
  INVENTORY_API_ERROR: 'INVENTORY_API_ERROR',
  ADDRESS_API_ERROR: 'ADDRESS_API_ERROR',
  REGION_API_ERROR: 'REGION_API_ERROR',
  EXTERNAL_API_TIMEOUT: 'EXTERNAL_API_TIMEOUT',

  // Security errors
  CSRF_TOKEN_MISSING: 'CSRF_TOKEN_MISSING',
  CSRF_TOKEN_INVALID: 'CSRF_TOKEN_INVALID',
  RATE_LIMIT_EXCEEDED: 'RATE_LIMIT_EXCEEDED',

  // Network errors
  NETWORK_ERROR: 'NETWORK_ERROR',
  TIMEOUT: 'TIMEOUT',
} as const;

export type ErrorCode = typeof ERROR_CODES[keyof typeof ERROR_CODES];

// User-friendly error messages
const ERROR_MESSAGES: Record<ErrorCode, string> = {
  [ERROR_CODES.INTERNAL_SERVER_ERROR]: 'サーバーでエラーが発生しました。時間をおいて再度お試しください。',
  [ERROR_CODES.BAD_REQUEST]: '入力内容に不備があります。',
  [ERROR_CODES.NOT_FOUND]: '要求されたリソースが見つかりません。',
  [ERROR_CODES.UNAUTHORIZED]: '認証が必要です。',
  [ERROR_CODES.FORBIDDEN]: 'アクセスが拒否されました。',
  [ERROR_CODES.CONFLICT]: '既に存在するデータです。',
  [ERROR_CODES.TOO_MANY_REQUESTS]: 'リクエスト数が多すぎます。時間をおいて再度お試しください。',

  [ERROR_CODES.VALIDATION_FAILED]: '入力内容を確認してください。',
  [ERROR_CODES.REQUIRED_FIELD_MISSING]: '必須項目が入力されていません。',
  [ERROR_CODES.INVALID_FORMAT]: '入力形式が正しくありません。',
  [ERROR_CODES.INVALID_EMAIL]: 'メールアドレスの形式が正しくありません。',
  [ERROR_CODES.INVALID_PHONE_NUMBER]: '電話番号の形式が正しくありません。',
  [ERROR_CODES.INVALID_POSTAL_CODE]: '郵便番号の形式が正しくありません。',
  [ERROR_CODES.EMAIL_CONFIRMATION_FAILED]: 'メールアドレスが一致しません。',

  [ERROR_CODES.USER_ALREADY_EXISTS]: '既に登録されているメールアドレスです。',
  [ERROR_CODES.SESSION_EXPIRED]: 'セッションが期限切れです。最初からやり直してください。',
  [ERROR_CODES.INVENTORY_NOT_AVAILABLE]: '選択されたオプションは在庫切れです。',
  [ERROR_CODES.REGION_NOT_SUPPORTED]: 'お住まいの地域ではご利用いただけないオプションが含まれています。',
  [ERROR_CODES.OPTION_NOT_AVAILABLE]: '選択されたオプションは現在ご利用いただけません。',
  [ERROR_CODES.ADDRESS_NOT_FOUND]: '入力された郵便番号では住所が見つかりません。',

  [ERROR_CODES.EXTERNAL_API_ERROR]: '外部サービスとの連携でエラーが発生しました。',
  [ERROR_CODES.INVENTORY_API_ERROR]: '在庫確認でエラーが発生しました。',
  [ERROR_CODES.ADDRESS_API_ERROR]: '住所検索でエラーが発生しました。',
  [ERROR_CODES.REGION_API_ERROR]: '地域確認でエラーが発生しました。',
  [ERROR_CODES.EXTERNAL_API_TIMEOUT]: '外部サービスとの通信がタイムアウトしました。',

  [ERROR_CODES.CSRF_TOKEN_MISSING]: 'セキュリティトークンが不足しています。ページを更新してください。',
  [ERROR_CODES.CSRF_TOKEN_INVALID]: 'セキュリティトークンが無効です。ページを更新してください。',
  [ERROR_CODES.RATE_LIMIT_EXCEEDED]: 'アクセス数が上限に達しました。時間をおいて再度お試しください。',

  [ERROR_CODES.NETWORK_ERROR]: 'ネットワークエラーが発生しました。接続を確認してください。',
  [ERROR_CODES.TIMEOUT]: '通信がタイムアウトしました。時間をおいて再度お試しください。',
};

// Error severity levels
export enum ErrorSeverity {
  LOW = 'low',
  MEDIUM = 'medium',
  HIGH = 'high',
  CRITICAL = 'critical',
}

// Error categories
export enum ErrorCategory {
  VALIDATION = 'validation',
  BUSINESS = 'business',
  TECHNICAL = 'technical',
  SECURITY = 'security',
  NETWORK = 'network',
}

interface ErrorConfig {
  severity: ErrorSeverity;
  category: ErrorCategory;
  retryable: boolean;
  userAction?: string;
}

// Error configuration mapping
const ERROR_CONFIG: Record<ErrorCode, ErrorConfig> = {
  [ERROR_CODES.INTERNAL_SERVER_ERROR]: {
    severity: ErrorSeverity.HIGH,
    category: ErrorCategory.TECHNICAL,
    retryable: true,
    userAction: '時間をおいて再度お試しいただくか、サポートにお問い合わせください。',
  },
  [ERROR_CODES.BAD_REQUEST]: {
    severity: ErrorSeverity.MEDIUM,
    category: ErrorCategory.VALIDATION,
    retryable: false,
  },
  [ERROR_CODES.NOT_FOUND]: {
    severity: ErrorSeverity.LOW,
    category: ErrorCategory.BUSINESS,
    retryable: false,
  },
  [ERROR_CODES.UNAUTHORIZED]: {
    severity: ErrorSeverity.HIGH,
    category: ErrorCategory.SECURITY,
    retryable: false,
  },
  [ERROR_CODES.FORBIDDEN]: {
    severity: ErrorSeverity.HIGH,
    category: ErrorCategory.SECURITY,
    retryable: false,
  },
  [ERROR_CODES.CONFLICT]: {
    severity: ErrorSeverity.MEDIUM,
    category: ErrorCategory.BUSINESS,
    retryable: false,
  },
  [ERROR_CODES.TOO_MANY_REQUESTS]: {
    severity: ErrorSeverity.MEDIUM,
    category: ErrorCategory.SECURITY,
    retryable: true,
    userAction: '少し時間をおいて再度お試しください。',
  },

  [ERROR_CODES.VALIDATION_FAILED]: {
    severity: ErrorSeverity.LOW,
    category: ErrorCategory.VALIDATION,
    retryable: false,
  },
  [ERROR_CODES.REQUIRED_FIELD_MISSING]: {
    severity: ErrorSeverity.LOW,
    category: ErrorCategory.VALIDATION,
    retryable: false,
  },
  [ERROR_CODES.INVALID_FORMAT]: {
    severity: ErrorSeverity.LOW,
    category: ErrorCategory.VALIDATION,
    retryable: false,
  },
  [ERROR_CODES.INVALID_EMAIL]: {
    severity: ErrorSeverity.LOW,
    category: ErrorCategory.VALIDATION,
    retryable: false,
  },
  [ERROR_CODES.INVALID_PHONE_NUMBER]: {
    severity: ErrorSeverity.LOW,
    category: ErrorCategory.VALIDATION,
    retryable: false,
  },
  [ERROR_CODES.INVALID_POSTAL_CODE]: {
    severity: ErrorSeverity.LOW,
    category: ErrorCategory.VALIDATION,
    retryable: false,
  },
  [ERROR_CODES.EMAIL_CONFIRMATION_FAILED]: {
    severity: ErrorSeverity.LOW,
    category: ErrorCategory.VALIDATION,
    retryable: false,
  },

  [ERROR_CODES.USER_ALREADY_EXISTS]: {
    severity: ErrorSeverity.MEDIUM,
    category: ErrorCategory.BUSINESS,
    retryable: false,
    userAction: '別のメールアドレスをご利用ください。',
  },
  [ERROR_CODES.SESSION_EXPIRED]: {
    severity: ErrorSeverity.MEDIUM,
    category: ErrorCategory.BUSINESS,
    retryable: false,
    userAction: 'ページを更新して最初からやり直してください。',
  },
  [ERROR_CODES.INVENTORY_NOT_AVAILABLE]: {
    severity: ErrorSeverity.MEDIUM,
    category: ErrorCategory.BUSINESS,
    retryable: true,
    userAction: '別のオプションをお選びいただくか、時間をおいて再度お試しください。',
  },
  [ERROR_CODES.REGION_NOT_SUPPORTED]: {
    severity: ErrorSeverity.MEDIUM,
    category: ErrorCategory.BUSINESS,
    retryable: false,
    userAction: 'ご利用可能なオプションをお選びください。',
  },
  [ERROR_CODES.OPTION_NOT_AVAILABLE]: {
    severity: ErrorSeverity.MEDIUM,
    category: ErrorCategory.BUSINESS,
    retryable: true,
  },
  [ERROR_CODES.ADDRESS_NOT_FOUND]: {
    severity: ErrorSeverity.LOW,
    category: ErrorCategory.BUSINESS,
    retryable: false,
    userAction: '郵便番号を確認するか、住所を手動で入力してください。',
  },

  [ERROR_CODES.EXTERNAL_API_ERROR]: {
    severity: ErrorSeverity.HIGH,
    category: ErrorCategory.TECHNICAL,
    retryable: true,
  },
  [ERROR_CODES.INVENTORY_API_ERROR]: {
    severity: ErrorSeverity.HIGH,
    category: ErrorCategory.TECHNICAL,
    retryable: true,
  },
  [ERROR_CODES.ADDRESS_API_ERROR]: {
    severity: ErrorSeverity.MEDIUM,
    category: ErrorCategory.TECHNICAL,
    retryable: true,
    userAction: '住所を手動で入力してください。',
  },
  [ERROR_CODES.REGION_API_ERROR]: {
    severity: ErrorSeverity.HIGH,
    category: ErrorCategory.TECHNICAL,
    retryable: true,
  },
  [ERROR_CODES.EXTERNAL_API_TIMEOUT]: {
    severity: ErrorSeverity.MEDIUM,
    category: ErrorCategory.TECHNICAL,
    retryable: true,
  },

  [ERROR_CODES.CSRF_TOKEN_MISSING]: {
    severity: ErrorSeverity.HIGH,
    category: ErrorCategory.SECURITY,
    retryable: true,
    userAction: 'ページを更新してください。',
  },
  [ERROR_CODES.CSRF_TOKEN_INVALID]: {
    severity: ErrorSeverity.HIGH,
    category: ErrorCategory.SECURITY,
    retryable: true,
    userAction: 'ページを更新してください。',
  },
  [ERROR_CODES.RATE_LIMIT_EXCEEDED]: {
    severity: ErrorSeverity.MEDIUM,
    category: ErrorCategory.SECURITY,
    retryable: true,
  },

  [ERROR_CODES.NETWORK_ERROR]: {
    severity: ErrorSeverity.HIGH,
    category: ErrorCategory.NETWORK,
    retryable: true,
    userAction: 'インターネット接続を確認してください。',
  },
  [ERROR_CODES.TIMEOUT]: {
    severity: ErrorSeverity.MEDIUM,
    category: ErrorCategory.NETWORK,
    retryable: true,
  },
};

export interface ProcessedError {
  code: ErrorCode;
  message: string;
  details?: Record<string, string>;
  severity: ErrorSeverity;
  category: ErrorCategory;
  retryable: boolean;
  userAction?: string;
  originalError?: ApiError;
}

class ErrorService {
  // Process API errors into user-friendly format
  processError(error: ApiError): ProcessedError {
    const code = error.code as ErrorCode;
    const config = ERROR_CONFIG[code] || ERROR_CONFIG[ERROR_CODES.INTERNAL_SERVER_ERROR];
    const message = ERROR_MESSAGES[code] || error.message || 'エラーが発生しました';

    return {
      code,
      message,
      details: error.details,
      severity: config.severity,
      category: config.category,
      retryable: config.retryable,
      userAction: config.userAction,
      originalError: error,
    };
  }

  // Get user-friendly message for error code
  getUserMessage(code: ErrorCode): string {
    return ERROR_MESSAGES[code] || 'エラーが発生しました';
  }

  // Check if error is retryable
  isRetryable(error: ApiError): boolean {
    const code = error.code as ErrorCode;
    const config = ERROR_CONFIG[code];
    return config?.retryable || false;
  }

  // Get error severity
  getSeverity(error: ApiError): ErrorSeverity {
    const code = error.code as ErrorCode;
    const config = ERROR_CONFIG[code];
    return config?.severity || ErrorSeverity.MEDIUM;
  }

  // Check if error should be reported to monitoring
  shouldReport(error: ApiError): boolean {
    const severity = this.getSeverity(error);
    return severity === ErrorSeverity.HIGH || severity === ErrorSeverity.CRITICAL;
  }

  // Format validation errors for display
  formatValidationErrors(details?: Record<string, string>): Record<string, string> {
    if (!details) return {};

    const formatted: Record<string, string> = {};
    
    Object.entries(details).forEach(([field, message]) => {
      // Convert backend field names to frontend field names if needed
      const frontendField = this.mapBackendFieldToFrontend(field);
      formatted[frontendField] = message;
    });

    return formatted;
  }

  // Map backend field names to frontend field names
  private mapBackendFieldToFrontend(backendField: string): string {
    const fieldMapping: Record<string, string> = {
      'last_name': 'lastName',
      'first_name': 'firstName',
      'last_name_kana': 'lastNameKana',
      'first_name_kana': 'firstNameKana',
      'phone1': 'phone1',
      'phone2': 'phone2',
      'phone3': 'phone3',
      'postal_code1': 'postalCode1',
      'postal_code2': 'postalCode2',
      'email_confirmation': 'emailConfirmation',
      'plan_type': 'planType',
      'option_types': 'selectedOptions',
    };

    return fieldMapping[backendField] || backendField;
  }

  // Log error for debugging
  logError(error: ProcessedError): void {
    const logLevel = this.getLogLevel(error.severity);
    
    console[logLevel]('Application Error:', {
      code: error.code,
      message: error.message,
      category: error.category,
      severity: error.severity,
      retryable: error.retryable,
      details: error.details,
      originalError: error.originalError,
      timestamp: new Date().toISOString(),
    });
  }

  private getLogLevel(severity: ErrorSeverity): 'error' | 'warn' | 'info' {
    switch (severity) {
      case ErrorSeverity.CRITICAL:
      case ErrorSeverity.HIGH:
        return 'error';
      case ErrorSeverity.MEDIUM:
        return 'warn';
      case ErrorSeverity.LOW:
      default:
        return 'info';
    }
  }
}

export const errorService = new ErrorService();