// Application constants
export const FORM_STEPS = {
  INPUT: 'input' as const,
  CONFIRM: 'confirm' as const,
  COMPLETE: 'complete' as const
} as const;

export const PLAN_TYPES = {
  A: 'A',
  B: 'B'
} as const;

export const OPTION_TYPES = {
  AA: 'AA',
  BB: 'BB', 
  AB: 'AB'
} as const;

// Plan-specific available options
export const PLAN_AVAILABLE_OPTIONS = {
  [PLAN_TYPES.A]: [OPTION_TYPES.AA, OPTION_TYPES.AB],
  [PLAN_TYPES.B]: [OPTION_TYPES.BB, OPTION_TYPES.AB]
} as const;

// Session configuration
export const SESSION_CONFIG = {
  TIMEOUT_HOURS: 4,
  WARNING_MINUTES: 15,
  STORAGE_KEYS: {
    FORM_DATA: 'membershipForm_data',
    SESSION_ID: 'membershipForm_sessionId',
    LAST_SAVED: 'membershipForm_lastSaved'
  }
} as const;

// Validation patterns
export const VALIDATION_PATTERNS = {
  PHONE: {
    PHONE1: /^(0\d{1,4})$/,
    PHONE2: /^\d{1,4}$/,
    PHONE3: /^\d{4}$/,
    MOBILE_PREFIX: /^0[789]0/,
    TOLL_FREE: ['0120', '0800', '0570', '0990']
  },
  POSTAL_CODE: {
    PART1: /^\d{3}$/,
    PART2: /^\d{4}$/
  },
  KATAKANA: /^[ァ-ヶー]+$/,
  EMAIL: /^[^\s@]+@[^\s@]+\.[^\s@]+$/
} as const;

// Form field limits
export const FIELD_LIMITS = {
  NAME: 15,
  NAME_KANA: 15,
  CITY: 50,
  TOWN: 50,
  CHOME: 10,
  BANCHI: 10,
  GO: 10,
  BUILDING: 100,
  ROOM: 20,
  EMAIL: 256
} as const;

// Error messages
export const ERROR_MESSAGES = {
  REQUIRED: '必須項目です',
  EMAIL_FORMAT: 'メールアドレスの形式が正しくありません',
  EMAIL_MISMATCH: 'メールアドレスが一致しません',
  KATAKANA_ONLY: '全角カタカナで入力してください',
  PHONE_FORMAT: '電話番号の形式が正しくありません',
  POSTAL_CODE_FORMAT: '郵便番号の形式が正しくありません',
  MOBILE_PHONE_FORMAT: '11桁の電話番号は携帯電話番号の形式で入力してください',
  TOLL_FREE_NOT_ALLOWED: 'フリーダイヤル等の番号は使用できません',
  MAX_LENGTH: (limit: number) => `${limit}文字以内で入力してください`,
  NETWORK_ERROR: 'ネットワークエラーが発生しました',
  API_ERROR: 'システムエラーが発生しました',
  SESSION_EXPIRED: 'セッションが無効です。最初からやり直してください。',
  INVENTORY_UNAVAILABLE: '選択されたオプションは在庫切れです',
  REGION_RESTRICTED: '選択されたオプションはお住まいの地域では選択できません'
} as const;

// Success messages
export const SUCCESS_MESSAGES = {
  USER_CREATED: '会員登録が完了しました',
  SESSION_SAVED: '入力内容を保存しました',
  ADDRESS_FOUND: '住所を自動入力しました'
} as const;

// Loading messages
export const LOADING_MESSAGES = {
  CREATING_USER: '登録処理中...',
  VALIDATING: '入力内容を確認中...',
  SEARCHING_ADDRESS: '住所を検索中...',
  CHECKING_INVENTORY: '在庫を確認中...',
  LOADING_OPTIONS: 'オプション情報を読み込み中...',
  SAVING_SESSION: '保存中...'
} as const;

// Button labels
export const BUTTON_LABELS = {
  NEXT: '次へ',
  PREV: '戻る',
  SUBMIT: '申し込む',
  SEARCH_ADDRESS: '住所検索',
  RESET: 'リセット',
  SAVE: '保存',
  EXTEND_SESSION: 'セッション延長'
} as const;

// Page titles
export const PAGE_TITLES = {
  INPUT: '会員登録 - 入力',
  CONFIRM: '会員登録 - 確認',
  COMPLETE: '会員登録 - 完了'
} as const;

// API endpoints (for reference)
export const API_ENDPOINTS = {
  HEALTH: '/health',
  PING: '/api/v1/ping',
  USERS: '/api/v1/users',
  USERS_VALIDATE: '/api/v1/users/validate',
  SESSIONS: '/api/v1/sessions',
  OPTIONS: '/api/v1/options',
  OPTIONS_INVENTORY: '/api/v1/options/check-inventory',
  PLANS: '/api/v1/plans',
  ADDRESS_SEARCH: '/api/v1/address/search',
  REGION_CHECK: '/api/v1/region/check',
  PREFECTURES: '/api/v1/prefectures'
} as const;

// Environment variables
export const ENV = {
  API_BASE_URL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  IS_DEV: import.meta.env.DEV,
  IS_PROD: import.meta.env.PROD
} as const;