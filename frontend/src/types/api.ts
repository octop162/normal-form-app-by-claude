// API response and request types for the membership registration system

// Common API response wrapper
export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: ApiError;
}

export interface ApiError {
  code: string;
  message: string;
  details?: Record<string, string>;
}

// User-related types
export interface UserCreateRequest {
  last_name: string;
  first_name: string;
  last_name_kana: string;
  first_name_kana: string;
  phone1: string;
  phone2: string;
  phone3: string;
  postal_code1: string;
  postal_code2: string;
  prefecture: string;
  city: string;
  town?: string;
  chome?: string;
  banchi: string;
  go?: string;
  building?: string;
  room?: string;
  email: string;
  email_confirm: string;
  plan_type: string;
  option_types: string[];
}

export interface UserValidateRequest {
  last_name: string;
  first_name: string;
  last_name_kana: string;
  first_name_kana: string;
  phone1: string;
  phone2: string;
  phone3: string;
  postal_code1: string;
  postal_code2: string;
  prefecture: string;
  city: string;
  town?: string;
  chome?: string;
  banchi: string;
  go?: string;
  building?: string;
  room?: string;
  email: string;
  email_confirm: string;
  plan_type: string;
  option_types: string[];
}

export interface UserCreateResponse {
  user_id: number;
  message: string;
}

export interface UserValidateResponse {
  valid: boolean;
  errors?: Record<string, string>;
}

// Session management types
export interface SessionCreateRequest {
  user_data: UserCreateRequest;
  expires_at?: string;
}

export interface SessionCreateResponse {
  session_id: string;
  expires_at: string;
}

export interface SessionGetResponse {
  session_id: string;
  user_data: UserCreateRequest;
  expires_at: string;
}

// Address and prefecture types
export interface AddressSearchRequest {
  postal_code: string;
}

export interface AddressSearchResponse {
  found: boolean;
  prefecture?: string;
  city?: string;
  town?: string;
  postal_code?: string;
}

export interface PrefectureResponse {
  id: number;
  prefecture_code: string;
  prefecture_name: string;
  region: string;
}

export interface PrefecturesGetResponse {
  prefectures: PrefectureResponse[];
}

// Option and plan types
export interface OptionResponse {
  option_type: string;
  option_name: string;
  description: string;
  is_active: boolean;
  price?: number;
}

export interface OptionsGetResponse {
  options: OptionResponse[];
}

export interface PlanResponse {
  plan_type: string;
  plan_name: string;
  description: string;
  base_price: number;
  is_active: boolean;
}

export interface PlansGetResponse {
  plans: PlanResponse[];
}

// Inventory and region check types
export interface InventoryCheckRequest {
  option_types: string[];
}

export interface InventoryCheckResponse {
  inventory: Record<string, number>;
}

export interface RegionCheckRequest {
  prefecture: string;
  city: string;
  option_types: string[];
}

export interface RegionCheckResponse {
  restrictions: Record<string, boolean>;
}

// Health check types
export interface HealthCheckResponse {
  status: string;
  service: string;
  version: string;
  timestamp: string;
  checks: Record<string, string>;
}