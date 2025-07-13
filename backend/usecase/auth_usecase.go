package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"posting-app/domain"
	"posting-app/infrastructure"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepo        domain.UserRepository
	passwordRepo    domain.PasswordResetRepository
	jwtManager      *infrastructure.JWTManager
}

func NewAuthUsecase(userRepo domain.UserRepository, passwordRepo domain.PasswordResetRepository, jwtManager *infrastructure.JWTManager) domain.AuthUsecase {
	return &authUsecase{
		userRepo:        userRepo,
		passwordRepo:    passwordRepo,
		jwtManager:      jwtManager,
	}
}

func (u *authUsecase) Register(username, email, password string) (*domain.User, error) {
	existing, _ := u.userRepo.GetByEmail(email)
	if existing != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Username:           username,
		Email:              email,
		PasswordHash:       string(hashedPassword),
		SubscriptionStatus: "inactive",
		IsAdmin:            false,
		IsActive:           true,
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (u *authUsecase) Login(email, password string) (string, *domain.User, error) {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	if user.SubscriptionStatus != "active" {
		return "", nil, fmt.Errorf("subscription required")
	}

	token, err := u.jwtManager.GenerateToken(user)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, nil
}

func (u *authUsecase) AdminLogin(email, password string) (string, *domain.User, error) {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	if !user.IsAdmin {
		return "", nil, fmt.Errorf("access denied")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	token, err := u.jwtManager.GenerateToken(user)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, nil
}

func (u *authUsecase) ValidateToken(token string) (*domain.User, error) {
	claims, err := u.jwtManager.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	user, err := u.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}

func (u *authUsecase) Logout(token string) error {
	return nil
}

func (u *authUsecase) ForgotPassword(email string) error {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// 既存のパスワードリセットトークンを削除
	u.passwordRepo.DeleteByUserID(user.ID)

	// 新しいトークンを生成
	token, err := u.generateResetToken()
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// 24時間の有効期限
	expiresAt := time.Now().Add(24 * time.Hour)

	reset := &domain.PasswordReset{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	if err := u.passwordRepo.Create(reset); err != nil {
		return fmt.Errorf("failed to create password reset: %w", err)
	}

	// 実際のアプリケーションでは、ここでメール送信ロジックを実装
	// 今回はログに出力
	fmt.Printf("Password reset token for %s: %s\n", email, token)

	return nil
}

func (u *authUsecase) ResetPassword(token, newPassword string) error {
	reset, err := u.passwordRepo.GetByToken(token)
	if err != nil {
		return fmt.Errorf("invalid or expired token")
	}

	user, err := u.userRepo.GetByID(reset.UserID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.PasswordHash = string(hashedPassword)
	if err := u.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// パスワードリセットトークンを削除
	u.passwordRepo.DeleteByUserID(user.ID)

	return nil
}

func (u *authUsecase) generateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}