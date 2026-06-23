package service

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"domesv2/config"
	"domesv2/internal/model"
	"domesv2/internal/repository"
	"domesv2/pkg/captcha"
	"domesv2/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req *model.RegisterRequest) (*model.AuthResponse, error)
	Login(req *model.LoginRequest) (*model.AuthResponse, error)
	ForgotPassword(req *model.ForgotPasswordRequest) error
	ResetPassword(req *model.ResetPasswordRequest) error
	GetProfile(userID uint) (*model.UserProfileResponse, error)
}

type authService struct {
	userRepo    repository.UserRepository
	mailService MailService
	cfg         *config.Config
}

func NewAuthService(userRepo repository.UserRepository, mailService MailService) AuthService {
	return &authService{
		userRepo:    userRepo,
		mailService: mailService,
		cfg:         config.AppConfig,
	}
}

func (s *authService) Register(req *model.RegisterRequest) (*model.AuthResponse, error) {
	if err := captcha.VerifyCaptcha(req.Captcha); err != nil {
		return nil, err
	}

	if req.FirstName == "" || req.LastName == "" {
		return nil, errors.NewValidationError("First name and last name are required", "VALIDATION_FAILED")
	}
	if req.Email == "" {
		return nil, errors.NewValidationError("Email is required", "VALIDATION_FAILED")
	}
	if req.Password == "" || len(req.Password) < 6 {
		return nil, errors.NewValidationError("Password must be at least 6 characters", "VALIDATION_FAILED")
	}
	if req.Password != req.ConfirmPassword {
		return nil, errors.NewValidationError("Passwords do not match", "VALIDATION_FAILED")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		zap.L().Error("Failed to hash password", zap.Error(err))
		return nil, errors.NewInternalServerError("Failed to create user", "PASSWORD_HASH_ERROR")
	}

	fullName := req.FirstName + " " + req.LastName
	username := strings.Split(req.Email, "@")[0]

	user := &model.User{
		Username:     &username,
		Name:         &fullName,
		FirstName:    &req.FirstName,
		LastName:     &req.LastName,
		Password:     string(hashedPassword),
		Email:        req.Email,
		Position:     &req.Position,
		Organization: &req.Organization,
		PhoneNumber:  &req.PhoneNumber,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *authService) Login(req *model.LoginRequest) (*model.AuthResponse, error) {
	if err := captcha.VerifyCaptcha(req.Captcha); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.NewUnauthorizedError("Invalid credentials", "INVALID_CREDENTIALS")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *authService) ForgotPassword(req *model.ForgotPasswordRequest) error {
	if err := captcha.VerifyCaptcha(req.Captcha); err != nil {
		return err
	}

	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return err
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		zap.L().Error("Failed to generate reset token", zap.Error(err))
		return errors.NewInternalServerError("Failed to process request", "INTERNAL_ERROR")
	}
	token := hex.EncodeToString(tokenBytes)

	expiry := time.Now().Add(1 * time.Hour)
	user.ResetPasswordToken = &token
	user.ResetPasswordExpiry = &expiry

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	if err := s.mailService.SendResetPassword(user.Email, token); err != nil {
		return err
	}

	return nil
}

func (s *authService) ResetPassword(req *model.ResetPasswordRequest) error {
	if req.Password == "" || len(req.Password) < 6 {
		return errors.NewValidationError("Password must be at least 6 characters", "VALIDATION_FAILED")
	}
	if req.Password != req.ConfirmPassword {
		return errors.NewValidationError("Passwords do not match", "VALIDATION_FAILED")
	}

	user, err := s.userRepo.FindByResetToken(req.Token)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		zap.L().Error("Failed to hash password", zap.Error(err))
		return errors.NewInternalServerError("Failed to reset password", "PASSWORD_HASH_ERROR")
	}

	if err := s.userRepo.UpdatePassword(user, string(hashedPassword)); err != nil {
		return err
	}

	return nil
}

func (s *authService) GetProfile(userID uint) (*model.UserProfileResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return &model.UserProfileResponse{
		ID:           user.ID,
		Username:     user.Username,
		Name:         user.Name,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Type:         user.Type,
		Position:     user.Position,
		Organization: user.Organization,
		PhoneNumber:  user.PhoneNumber,
	}, nil
}

func (s *authService) generateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(s.cfg.JWT.ExpiresIn).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		zap.L().Error("Failed to generate JWT token", zap.Error(err))
		return "", errors.NewInternalServerError("Failed to generate token", "TOKEN_GENERATION_ERROR")
	}

	return signedToken, nil
}
