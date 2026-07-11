package captcha

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"domesv2/config"
	"domesv2/config/logger"
	"domesv2/pkg/errors"
	"go.uber.org/zap"
)

type RecaptchaResponse struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes,omitempty"`
}

func VerifyCaptcha(response string) error {
	cfg := config.AppConfig
	env := strings.ToLower(cfg.Server.Env)
	if !cfg.Captcha.Enabled || env == "local" || env == "development" || env == "dev" {
		return nil
	}

	if response == "" {
		logger.Warn("Captcha verification failed: missing response", zap.String("reason", "empty_response"))
		return errors.NewValidationError("Captcha verification failed", "CAPTCHA_MISSING")
	}

	form := url.Values{}
	form.Add("secret", cfg.Captcha.SecretKey)
	form.Add("response", response)

	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", form)
	if err != nil {
		logger.Error("Failed to verify captcha", zap.Error(err))
		return errors.NewInternalServerError("Captcha verification failed", "CAPTCHA_VERIFICATION_FAILED")
	}
	defer resp.Body.Close()

	var recaptchaResp RecaptchaResponse
	if err := json.NewDecoder(resp.Body).Decode(&recaptchaResp); err != nil {
		logger.Error("Failed to decode recaptcha response", zap.Error(err))
		return errors.NewInternalServerError("Captcha verification failed", "CAPTCHA_DECODE_ERROR")
	}

	if !recaptchaResp.Success {
		logger.Warn("Captcha verification failed", zap.Strings("error_codes", recaptchaResp.ErrorCodes))
		return errors.NewValidationError("Invalid captcha response", "CAPTCHA_INVALID")
	}

	logger.Debug("Captcha verified successfully")
	return nil
}
