// Package service provides business logic layer for the application.
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/octop162/normal-form-app-by-claude/internal/dto"
	"github.com/octop162/normal-form-app-by-claude/internal/model"
	"github.com/octop162/normal-form-app-by-claude/internal/repository"
	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
	"github.com/octop162/normal-form-app-by-claude/pkg/validator"
)

// UserService defines the interface for user business logic
type UserService interface {
	CreateUser(ctx context.Context, req *dto.UserCreateRequest) (*dto.UserCreateResponse, error)
	ValidateUserData(ctx context.Context, req *dto.UserValidateRequest) (*dto.UserValidateResponse, error)
	GetUserByID(ctx context.Context, id int) (*dto.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error)
	UpdateUser(ctx context.Context, id int, req *dto.UserCreateRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id int) error
}

// userService implements UserService
type userService struct {
	userRepo       repository.UserRepository
	userOptionRepo repository.UserOptionRepository
	optionRepo     repository.OptionRepository
	validator      *validator.CustomValidator
	log            *logger.Logger
}

// NewUserService creates a new user service
func NewUserService(
	userRepo repository.UserRepository,
	userOptionRepo repository.UserOptionRepository,
	optionRepo repository.OptionRepository,
	validator *validator.CustomValidator,
	log *logger.Logger,
) UserService {
	return &userService{
		userRepo:       userRepo,
		userOptionRepo: userOptionRepo,
		optionRepo:     optionRepo,
		validator:      validator,
		log:            log,
	}
}

// CreateUser creates a new user with validation
func (s *userService) CreateUser(ctx context.Context, req *dto.UserCreateRequest) (*dto.UserCreateResponse, error) {
	// Validate request
	validationResp, err := s.ValidateUserData(ctx, &dto.UserValidateRequest{UserCreateRequest: *req})
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if !validationResp.Valid {
		return nil, fmt.Errorf("validation errors: %v", validationResp.Errors)
	}

	// Check if user already exists
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		s.log.WithError(err).Error("Failed to check user existence")
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Convert DTO to model
	user := s.convertCreateRequestToModel(req)

	// Create user
	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		s.log.WithError(err).Error("Failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create user options if any
	if len(req.OptionTypes) > 0 {
		userOptions := make([]*model.UserOption, 0, len(req.OptionTypes))
		for _, optionType := range req.OptionTypes {
			userOptions = append(userOptions, &model.UserOption{
				UserID:     createdUser.ID,
				OptionType: optionType,
			})
		}

		if err := s.userOptionRepo.CreateBatch(ctx, userOptions); err != nil {
			s.log.WithError(err).Error("Failed to create user options")
			return nil, fmt.Errorf("failed to create user options: %w", err)
		}
	}

	s.log.WithField("user_id", createdUser.ID).Info("User created successfully with options")

	return &dto.UserCreateResponse{
		ID:      createdUser.ID,
		Message: "User created successfully",
	}, nil
}

// ValidateUserData validates user registration data
func (s *userService) ValidateUserData(
	ctx context.Context, req *dto.UserValidateRequest,
) (*dto.UserValidateResponse, error) {
	errors := make(map[string]string)

	// Struct validation
	if err := s.validator.ValidateStruct(req); err != nil {
		s.log.WithError(err).Debug("Struct validation failed")
		// Convert validation errors to map
		// Note: This is a simplified version - production code would parse validation errors properly
		errors["validation"] = err.Error()
	}

	// Business logic validation
	s.validateBusinessRules(ctx, &req.UserCreateRequest, errors)

	valid := len(errors) == 0

	return &dto.UserValidateResponse{
		Valid:  valid,
		Errors: errors,
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(ctx context.Context, id int) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.log.WithError(err).WithField("user_id", id).Error("Failed to get user by ID")
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return s.convertModelToResponse(user), nil
}

// GetUserByEmail retrieves a user by email
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.log.WithError(err).WithField("email", email).Error("Failed to get user by email")
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return s.convertModelToResponse(user), nil
}

// UpdateUser updates an existing user
func (s *userService) UpdateUser(ctx context.Context, id int, req *dto.UserCreateRequest) (*dto.UserResponse, error) {
	// Validate request
	validationResp, err := s.ValidateUserData(ctx, &dto.UserValidateRequest{UserCreateRequest: *req})
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if !validationResp.Valid {
		return nil, fmt.Errorf("validation errors: %v", validationResp.Errors)
	}

	// Get existing user
	existingUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check email uniqueness if email is being changed
	if existingUser.Email != req.Email {
		emailExists, emailErr := s.userRepo.ExistsByEmail(ctx, req.Email)
		if emailErr != nil {
			return nil, fmt.Errorf("failed to check email uniqueness: %w", emailErr)
		}
		if emailExists {
			return nil, fmt.Errorf("user with email %s already exists", req.Email)
		}
	}

	// Update user fields
	s.updateUserFields(existingUser, req)

	// Update user
	updatedUser, err := s.userRepo.Update(ctx, existingUser)
	if err != nil {
		s.log.WithError(err).Error("Failed to update user")
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Update user options
	if err := s.updateUserOptions(ctx, id, req.OptionTypes); err != nil {
		s.log.WithError(err).Error("Failed to update user options")
		return nil, fmt.Errorf("failed to update user options: %w", err)
	}

	s.log.WithField("user_id", id).Info("User updated successfully")

	return s.convertModelToResponse(updatedUser), nil
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(ctx context.Context, id int) error {
	// Delete user options first
	if err := s.userOptionRepo.DeleteByUserID(ctx, id); err != nil {
		s.log.WithError(err).Error("Failed to delete user options")
		return fmt.Errorf("failed to delete user options: %w", err)
	}

	// Delete user
	if err := s.userRepo.Delete(ctx, id); err != nil {
		s.log.WithError(err).Error("Failed to delete user")
		return fmt.Errorf("failed to delete user: %w", err)
	}

	s.log.WithField("user_id", id).Info("User deleted successfully")
	return nil
}

// validateBusinessRules validates business-specific rules
func (s *userService) validateBusinessRules(
	ctx context.Context, req *dto.UserCreateRequest, errors map[string]string,
) {
	// Validate phone number format
	fullPhone := req.Phone1 + req.Phone2 + req.Phone3
	if !validator.IsValidPhone(fullPhone) {
		errors["phone"] = "Invalid phone number format"
	}

	// Validate postal code
	fullPostalCode := req.PostalCode1 + "-" + req.PostalCode2
	if !validator.IsValidPostalCode(fullPostalCode) {
		errors["postal_code"] = "Invalid postal code format"
	}

	// Validate plan type
	if !validator.IsValidPlanType(req.PlanType) {
		errors["plan_type"] = "Invalid plan type"
	}

	// Validate option types
	for _, optionType := range req.OptionTypes {
		if !validator.IsValidOptionType(optionType) {
			errors["option_types"] = "Invalid option type: " + optionType
			break
		}

		// Check if option is compatible with plan
		option, err := s.optionRepo.GetByOptionType(ctx, optionType)
		if err != nil {
			errors["option_types"] = "Option not found: " + optionType
			continue
		}

		if !s.isOptionCompatibleWithPlan(option, req.PlanType) {
			errors["option_types"] = fmt.Sprintf("Option %s is not compatible with plan %s", optionType, req.PlanType)
			break
		}
	}
}

// isOptionCompatibleWithPlan checks if an option is compatible with a plan
func (s *userService) isOptionCompatibleWithPlan(option *model.OptionMaster, planType string) bool {
	switch option.PlanCompatibility {
	case "A":
		return planType == "A"
	case "B":
		return planType == "B"
	case "AB":
		return planType == "A" || planType == "B"
	default:
		return false
	}
}

// convertCreateRequestToModel converts DTO to model
func (s *userService) convertCreateRequestToModel(req *dto.UserCreateRequest) *model.User {
	return &model.User{
		LastName:      req.LastName,
		FirstName:     req.FirstName,
		LastNameKana:  req.LastNameKana,
		FirstNameKana: req.FirstNameKana,
		Phone1:        req.Phone1,
		Phone2:        req.Phone2,
		Phone3:        req.Phone3,
		PostalCode1:   req.PostalCode1,
		PostalCode2:   req.PostalCode2,
		Prefecture:    req.Prefecture,
		City:          req.City,
		Town:          req.Town,
		Chome:         req.Chome,
		Banchi:        req.Banchi,
		Go:            req.Go,
		Building:      req.Building,
		Room:          req.Room,
		Email:         req.Email,
		PlanType:      req.PlanType,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// convertModelToResponse converts model to response DTO
func (s *userService) convertModelToResponse(user *model.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:            user.ID,
		LastName:      user.LastName,
		FirstName:     user.FirstName,
		LastNameKana:  user.LastNameKana,
		FirstNameKana: user.FirstNameKana,
		PhoneNumber:   user.GetPhoneNumber(),
		PostalCode:    user.GetPostalCode(),
		Address:       user.GetFullAddress(),
		Email:         user.Email,
		PlanType:      user.PlanType,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
}

// updateUserFields updates user fields from request
func (s *userService) updateUserFields(user *model.User, req *dto.UserCreateRequest) {
	user.LastName = req.LastName
	user.FirstName = req.FirstName
	user.LastNameKana = req.LastNameKana
	user.FirstNameKana = req.FirstNameKana
	user.Phone1 = req.Phone1
	user.Phone2 = req.Phone2
	user.Phone3 = req.Phone3
	user.PostalCode1 = req.PostalCode1
	user.PostalCode2 = req.PostalCode2
	user.Prefecture = req.Prefecture
	user.City = req.City
	user.Town = req.Town
	user.Chome = req.Chome
	user.Banchi = req.Banchi
	user.Go = req.Go
	user.Building = req.Building
	user.Room = req.Room
	user.Email = req.Email
	user.PlanType = req.PlanType
}

// updateUserOptions updates user options
func (s *userService) updateUserOptions(ctx context.Context, userID int, optionTypes []string) error {
	// Delete existing options
	if err := s.userOptionRepo.DeleteByUserID(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete existing options: %w", err)
	}

	// Create new options
	if len(optionTypes) > 0 {
		userOptions := make([]*model.UserOption, 0, len(optionTypes))
		for _, optionType := range optionTypes {
			userOptions = append(userOptions, &model.UserOption{
				UserID:     userID,
				OptionType: optionType,
			})
		}

		if err := s.userOptionRepo.CreateBatch(ctx, userOptions); err != nil {
			return fmt.Errorf("failed to create new options: %w", err)
		}
	}

	return nil
}
