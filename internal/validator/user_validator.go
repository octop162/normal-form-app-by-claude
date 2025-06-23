package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/octop162/normal-form-app-by-claude/internal/handler"
)

// UserValidator handles validation for user-related data
type UserValidator struct{}

// NewUserValidator creates a new UserValidator instance
func NewUserValidator() *UserValidator {
	return &UserValidator{}
}

// ValidateUserCreation validates user creation data
func (v *UserValidator) ValidateUserCreation(data map[string]interface{}) map[string]string {
	errors := make(map[string]string)

	// Personal information validation
	if err := v.validateName(data, "last_name", "姓", errors); err != nil {
		errors["last_name"] = err.Error()
	}
	if err := v.validateName(data, "first_name", "名", errors); err != nil {
		errors["first_name"] = err.Error()
	}
	if err := v.validateKanaName(data, "last_name_kana", "姓カナ", errors); err != nil {
		errors["last_name_kana"] = err.Error()
	}
	if err := v.validateKanaName(data, "first_name_kana", "名カナ", errors); err != nil {
		errors["first_name_kana"] = err.Error()
	}

	// Phone number validation
	if err := v.validatePhoneNumber(data, errors); err != nil {
		errors["phone"] = err.Error()
	}

	// Postal code validation
	if err := v.validatePostalCode(data, errors); err != nil {
		errors["postal_code"] = err.Error()
	}

	// Address validation
	if err := v.validateAddress(data, errors); err != nil {
		errors["address"] = err.Error()
	}

	// Email validation
	if err := v.validateEmail(data, errors); err != nil {
		errors["email"] = err.Error()
	}

	// Plan and options validation
	if err := v.validatePlanAndOptions(data, errors); err != nil {
		errors["plan_options"] = err.Error()
	}

	return errors
}

// validateName validates name fields (last_name, first_name)
func (v *UserValidator) validateName(data map[string]interface{}, field, fieldName string, errors map[string]string) error {
	value, exists := data[field]
	if !exists {
		return &handler.AppError{
			Code:    handler.ErrorCodeRequiredFieldMissing,
			Message: fieldName + "は必須です",
		}
	}

	str, ok := value.(string)
	if !ok {
		return &handler.AppError{
			Code:    handler.ErrorCodeInvalidFormat,
			Message: fieldName + "は文字列で入力してください",
		}
	}

	str = strings.TrimSpace(str)
	if str == "" {
		return &handler.AppError{
			Code:    handler.ErrorCodeRequiredFieldMissing,
			Message: fieldName + "は必須です",
		}
	}

	if utf8.RuneCountInString(str) > 15 {
		return &handler.AppError{
			Code:    handler.ErrorCodeValueTooLong,
			Message: fieldName + "は15文字以内で入力してください",
		}
	}

	// Check for invalid characters (basic check)
	if matched, _ := regexp.MatchString(`[<>&"'\\]`, str); matched {
		return &handler.AppError{
			Code:    handler.ErrorCodeInvalidFormat,
			Message: fieldName + "に使用できない文字が含まれています",
		}
	}

	return nil
}

// validateKanaName validates katakana name fields
func (v *UserValidator) validateKanaName(data map[string]interface{}, field, fieldName string, errors map[string]string) error {
	value, exists := data[field]
	if !exists {
		return &handler.AppError{
			Code:    handler.ErrorCodeRequiredFieldMissing,
			Message: fieldName + "は必須です",
		}
	}

	str, ok := value.(string)
	if !ok {
		return &handler.AppError{
			Code:    handler.ErrorCodeInvalidFormat,
			Message: fieldName + "は文字列で入力してください",
		}
	}

	str = strings.TrimSpace(str)
	if str == "" {
		return &handler.AppError{
			Code:    handler.ErrorCodeRequiredFieldMissing,
			Message: fieldName + "は必須です",
		}
	}

	if utf8.RuneCountInString(str) > 15 {
		return &handler.AppError{
			Code:    handler.ErrorCodeValueTooLong,
			Message: fieldName + "は15文字以内で入力してください",
		}
	}

	// Check for full-width katakana only
	kanaPattern := regexp.MustCompile(`^[ァ-ヶー\s]+$`)
	if !kanaPattern.MatchString(str) {
		return &handler.AppError{
			Code:    handler.ErrorCodeInvalidFormat,
			Message: fieldName + "は全角カタカナで入力してください",
		}
	}

	return nil
}

// validatePhoneNumber validates phone number (3-part format)
func (v *UserValidator) validatePhoneNumber(data map[string]interface{}, errors map[string]string) error {
	phone1, exists1 := data["phone1"]
	phone2, exists2 := data["phone2"]
	phone3, exists3 := data["phone3"]

	if !exists1 || !exists2 || !exists3 {
		return &handler.AppError{
			Code:    handler.ErrorCodeRequiredFieldMissing,
			Message: "電話番号は必須です",
		}
	}

	p1, ok1 := phone1.(string)
	p2, ok2 := phone2.(string)
	p3, ok3 := phone3.(string)

	if !ok1 || !ok2 || !ok3 {
		return &handler.AppError{
			Code:    handler.ErrorCodeInvalidFormat,
			Message: "電話番号は正しい形式で入力してください",
		}
	}

	p1 = strings.TrimSpace(p1)
	p2 = strings.TrimSpace(p2)
	p3 = strings.TrimSpace(p3)

	if p1 == "" || p2 == "" || p3 == "" {
		return &handler.AppError{
			Code:    handler.ErrorCodeRequiredFieldMissing,
			Message: "電話番号は必須です",
		}
	}

	// Validate numeric characters only
	numberPattern := regexp.MustCompile(`^\d+$`)
	if !numberPattern.MatchString(p1) || !numberPattern.MatchString(p2) || !numberPattern.MatchString(p3) {
		return &handler.AppError{
			Code:    handler.ErrorCodeInvalidFormat,
			Message: "電話番号は数字のみで入力してください",
		}
	}

	fullNumber := p1 + p2 + p3

	// Check for free dial numbers (not allowed)
	freeDialPrefixes := []string{"0120", "0800", "0570"}
	for _, prefix := range freeDialPrefixes {
		if strings.HasPrefix(fullNumber, prefix) {
			return &handler.AppError{
				Code:    handler.ErrorCodeInvalidPhoneNumber,
				Message: "フリーダイヤル番号は使用できません",
			}
		}
	}

	// Validate length and format
	if len(fullNumber) == 11 {
		// Mobile number: must start with 0X0 (070, 080, 090)
		mobilePattern := regexp.MustCompile(`^0[789]0\d{8}$`)
		if !mobilePattern.MatchString(fullNumber) {
			return &handler.AppError{
				Code:    handler.ErrorCodeInvalidPhoneNumber,
				Message: "携帯電話番号の形式が正しくありません",
			}
		}
	} else if len(fullNumber) == 10 {
		// Landline number
		landlinePattern := regexp.MustCompile(`^0[1-9]\d{8}$`)
		if !landlinePattern.MatchString(fullNumber) {
			return &handler.AppError{
				Code:    handler.ErrorCodeInvalidPhoneNumber,
				Message: "固定電話番号の形式が正しくありません",
			}
		}
	} else {
		return &handler.AppError{
			Code:    handler.ErrorCodeInvalidPhoneNumber,
			Message: "電話番号は10桁または11桁で入力してください",
		}
	}

	// Validate part lengths
	if len(p1) < 2 || len(p1) > 5 {
		return &handler.AppError{
			Code:    handler.ErrorCodeInvalidPhoneNumber,
			Message: "市外局番は2-5桁で入力してください",
		}
	}
	if len(p2) < 1 || len(p2) > 4 {
		return &handler.AppError{
			Code:    handler.ErrorCodeInvalidPhoneNumber,
			Message: "市内局番は1-4桁で入力してください",
		}
	}
	if len(p3) != 4 {
		return &handler.AppError{
			Code:    handler.ErrorCodeInvalidPhoneNumber,
			Message: "契約番号は4桁で入力してください",
		}
	}

	return nil
}

// validatePostalCode validates postal code (2-part format)
func (v *UserValidator) validatePostalCode(data map[string]interface{}, errors map[string]string) error {
	postal1, exists1 := data["postal_code1"]
	postal2, exists2 := data["postal_code2"]

	if !exists1 || !exists2 {
		return &handler.AppError{
			Code:    handler.ErrorCodeRequiredFieldMissing,
			Message: "郵便番号は必須です",
		}
	}

	p1, ok1 := postal1.(string)
	p2, ok2 := postal2.(string)

	if !ok1 || !ok2 {
		return &handler.AppError{
			Code:    handler.ErrorCodeInvalidFormat,
			Message: "郵便番号は正しい形式で入力してください",
		}
	}

	p1 = strings.TrimSpace(p1)
	p2 = strings.TrimSpace(p2)

	if p1 == "" || p2 == "" {
		return &handler.AppError{
			Code:    handler.ErrorCodeRequiredFieldMissing,
			Message: "郵便番号は必須です",
		}
	}

	// Validate format: 3 digits + 4 digits
	if len(p1) != 3 || len(p2) != 4 {
		return &handler.AppError{
			Code:    handler.ErrorCodeInvalidPostalCode,
			Message: "郵便番号は3桁-4桁の形式で入力してください",
		}
	}

	numberPattern := regexp.MustCompile(`^\d+$`)
	if !numberPattern.MatchString(p1) || !numberPattern.MatchString(p2) {
		return &handler.AppError{
			Code:    handler.ErrorCodeInvalidPostalCode,
			Message: "郵便番号は数字のみで入力してください",
		}
	}

	return nil
}

// validateAddress validates address fields
func (v *UserValidator) validateAddress(data map[string]interface{}, errors map[string]string) error {
	// Prefecture (required)
	if err := v.validateRequiredField(data, "prefecture", "都道府県"); err != nil {
		return err
	}

	// City (required)
	if err := v.validateRequiredField(data, "city", "市区町村"); err != nil {
		return err
	}

	// Banchi (required)
	if err := v.validateRequiredField(data, "banchi", "番地"); err != nil {
		return err
	}

	// Optional fields validation
	optionalFields := map[string]string{
		"town":     "町名",
		"chome":    "丁目",
		"go":       "号",
		"building": "建物名",
		"room":     "部屋番号",
	}

	for field, fieldName := range optionalFields {
		if value, exists := data[field]; exists {
			if str, ok := value.(string); ok && str != "" {
				if utf8.RuneCountInString(str) > 50 {
					return &handler.AppError{
						Code:    handler.ErrorCodeValueTooLong,
						Message: fieldName + "は50文字以内で入力してください",
					}
				}
			}
		}
	}

	return nil
}

// validateEmail validates email and email confirmation
func (v *UserValidator) validateEmail(data map[string]interface{}, errors map[string]string) error {
	email, emailExists := data["email"]
	emailConfirm, confirmExists := data["email_confirmation"]

	if !emailExists {
		return &handler.AppError{
			Code:    handler.ErrorCodeRequiredFieldMissing,
			Message: "メールアドレスは必須です",
		}
	}

	emailStr, emailOk := email.(string)
	if !emailOk {
		return &handler.AppError{
			Code:    handler.ErrorCodeInvalidFormat,
			Message: "メールアドレスは文字列で入力してください",
		}
	}

	emailStr = strings.TrimSpace(emailStr)
	if emailStr == "" {
		return &handler.AppError{
			Code:    handler.ErrorCodeRequiredFieldMissing,
			Message: "メールアドレスは必須です",
		}
	}

	if len(emailStr) > 256 {
		return &handler.AppError{
			Code:    handler.ErrorCodeValueTooLong,
			Message: "メールアドレスは256文字以内で入力してください",
		}
	}

	// Email format validation (RFC 5322 compliant)
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailPattern.MatchString(emailStr) {
		return &handler.AppError{
			Code:    handler.ErrorCodeInvalidEmail,
			Message: "メールアドレスの形式が正しくありません",
		}
	}

	// Email confirmation validation
	if confirmExists {
		confirmStr, confirmOk := emailConfirm.(string)
		if confirmOk {
			confirmStr = strings.TrimSpace(confirmStr)
			if emailStr != confirmStr {
				return &handler.AppError{
					Code:    handler.ErrorCodeEmailConfirmationFail,
					Message: "メールアドレスが一致しません",
				}
			}
		}
	}

	return nil
}

// validatePlanAndOptions validates plan type and selected options
func (v *UserValidator) validatePlanAndOptions(data map[string]interface{}, errors map[string]string) error {
	planType, exists := data["plan_type"]
	if !exists {
		return &handler.AppError{
			Code:    handler.ErrorCodeRequiredFieldMissing,
			Message: "プランは必須です",
		}
	}

	planStr, ok := planType.(string)
	if !ok {
		return &handler.AppError{
			Code:    handler.ErrorCodeInvalidFormat,
			Message: "プランは文字列で指定してください",
		}
	}

	planStr = strings.TrimSpace(planStr)
	if planStr == "" {
		return &handler.AppError{
			Code:    handler.ErrorCodeRequiredFieldMissing,
			Message: "プランは必須です",
		}
	}

	// Validate plan type
	validPlans := map[string]bool{
		"A": true,
		"B": true,
	}
	if !validPlans[planStr] {
		return &handler.AppError{
			Code:    handler.ErrorCodeInvalidFormat,
			Message: "無効なプランが選択されています",
		}
	}

	// Validate options if provided
	if optionsData, exists := data["option_types"]; exists {
		if options, ok := optionsData.([]interface{}); ok {
			for _, option := range options {
				if optionStr, ok := option.(string); ok {
					if err := v.validateOptionForPlan(optionStr, planStr); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// validateOptionForPlan validates if an option is available for the selected plan
func (v *UserValidator) validateOptionForPlan(option, plan string) error {
	validOptions := map[string]map[string]bool{
		"A": {
			"AA": true,
			"AB": true,
		},
		"B": {
			"BB": true,
			"AB": true,
		},
	}

	planOptions, planExists := validOptions[plan]
	if !planExists {
		return &handler.AppError{
			Code:    handler.ErrorCodeInvalidFormat,
			Message: "無効なプランです",
		}
	}

	if !planOptions[option] {
		return &handler.AppError{
			Code:    handler.ErrorCodeOptionNotAvailable,
			Message: "選択されたオプションは指定されたプランでは利用できません",
		}
	}

	return nil
}

// validateRequiredField validates a required string field
func (v *UserValidator) validateRequiredField(data map[string]interface{}, field, fieldName string) error {
	value, exists := data[field]
	if !exists {
		return &handler.AppError{
			Code:    handler.ErrorCodeRequiredFieldMissing,
			Message: fieldName + "は必須です",
		}
	}

	str, ok := value.(string)
	if !ok {
		return &handler.AppError{
			Code:    handler.ErrorCodeInvalidFormat,
			Message: fieldName + "は文字列で入力してください",
		}
	}

	str = strings.TrimSpace(str)
	if str == "" {
		return &handler.AppError{
			Code:    handler.ErrorCodeRequiredFieldMissing,
			Message: fieldName + "は必須です",
		}
	}

	return nil
}