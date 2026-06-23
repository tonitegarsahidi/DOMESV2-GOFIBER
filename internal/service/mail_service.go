package service

import (
	"domesv2/config/logger"
	"go.uber.org/zap"
)

type MailService interface {
	SendResetPassword(email, token string) error
}

type mailService struct{}

func NewMailService() MailService {
	return &mailService{}
}

func (s *mailService) SendResetPassword(email, token string) error {
	resetURL := "http://localhost:3000/reset-password?token=" + token

	logger.Info("Password reset email",
		zap.String("to", email),
		zap.String("reset_url", resetURL),
	)

	return nil
}
