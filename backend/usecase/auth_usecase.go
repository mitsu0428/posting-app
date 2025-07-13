package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
	"posting-app/domain"
	"posting-app/infrastructure"
	"posting-app/repository"
)

type AuthUsecase struct {
	userRepo          *repository.UserRepository
	passwordResetRepo *repository.PasswordResetRepository
	jwtService        *infrastructure.JWTService
}

func NewAuthUsecase(
	userRepo *repository.UserRepository,
	passwordResetRepo *repository.PasswordResetRepository,
	jwtService *infrastructure.JWTService,
) *AuthUsecase {
	return &AuthUsecase{
		userRepo:          userRepo,
		passwordResetRepo: passwordResetRepo,
		jwtService:        jwtService,
	}
}

func (u *AuthUsecase) Register(email, password, displayName string) (*domain.User, error) {
	// Check if user already exists
	_, err := u.userRepo.GetByEmail(email)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Email:              email,
		PasswordHash:       string(hashedPassword),
		DisplayName:        displayName,
		Role:               string(domain.UserRoleUser),
		SubscriptionStatus: domain.UserSubscriptionStatusInactive,
		IsActive:           true,
		EmailVerified:      true, // For simplicity, we'll skip email verification
	}

	err = u.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	slog.Info("User registered successfully", "email", email)
	return user, nil
}

func (u *AuthUsecase) Login(email, password string) (*domain.User, string, error) {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, "", errors.New("account is deactivated")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := u.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	slog.Info("User logged in successfully", "email", email)
	return user, token, nil
}

func (u *AuthUsecase) AdminLogin(email, password string) (*domain.User, string, error) {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, "", errors.New("account is deactivated")
	}

	if user.Role != string(domain.UserRoleAdmin) {
		return nil, "", errors.New("admin access required")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := u.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	slog.Info("Admin logged in successfully", "email", email)
	return user, token, nil
}

func (u *AuthUsecase) ForgotPassword(email string) error {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		// For security reasons, we don't reveal if the email exists
		slog.Info("Password reset requested for non-existent email", "email", email)
		return nil
	}

	if !user.IsActive {
		slog.Info("Password reset requested for inactive user", "email", email)
		return nil
	}

	// Generate reset token
	tokenBytes := make([]byte, 32)
	_, err = rand.Read(tokenBytes)
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	passwordReset := &domain.PasswordReset{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour), // 1 hour expiry
	}

	err = u.passwordResetRepo.Create(passwordReset)
	if err != nil {
		return fmt.Errorf("failed to create password reset: %w", err)
	}

	// TODO: Send email with reset token
	slog.Info("Password reset token generated", "email", email, "token", token)

	return nil
}

func (u *AuthUsecase) ResetPassword(token, newPassword string) error {
	passwordReset, err := u.passwordResetRepo.GetByToken(token)
	if err != nil {
		return errors.New("invalid or expired reset token")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user password
	err = u.userRepo.UpdatePassword(passwordReset.UserID, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark reset token as used
	err = u.passwordResetRepo.MarkAsUsed(passwordReset.ID)
	if err != nil {
		return fmt.Errorf("failed to mark reset token as used: %w", err)
	}

	slog.Info("Password reset successfully", "user_id", passwordReset.UserID)
	return nil
}

func (u *AuthUsecase) ChangePassword(userID int, currentPassword, newPassword string) error {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword))
	if err != nil {
		return errors.New("current password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	err = u.userRepo.UpdatePassword(userID, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	slog.Info("Password changed successfully", "user_id", userID)
	return nil
}
