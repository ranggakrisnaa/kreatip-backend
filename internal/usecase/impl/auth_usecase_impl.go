package impl

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hibiken/asynq"
	"github.com/kreatip/kreatip-backend/internal/entity"
	gatewaymsg "github.com/kreatip/kreatip-backend/internal/gateway/messaging"
	"github.com/kreatip/kreatip-backend/internal/model"
	"github.com/kreatip/kreatip-backend/internal/model/converter"
	"github.com/kreatip/kreatip-backend/internal/repository"
	"github.com/kreatip/kreatip-backend/internal/usecase"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type JWTConfig struct {
	Secret     string
	AccessTTL  string
	RefreshTTL string
	AppURL     string // base URL for email verification links, e.g. https://app.kreatip.id
}

type jwtClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	Role  string `json:"role"`
}

type authUseCase struct {
	db               *gorm.DB
	log              *logrus.Logger
	validate         *validator.Validate
	jwtCfg           JWTConfig
	producer         *gatewaymsg.Producer
	userRepo         repository.UserRepository
	walletRepo       repository.WalletRepository
	refreshTokenRepo repository.RefreshTokenRepository
}

func NewAuthUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	jwtCfg JWTConfig,
	producer *gatewaymsg.Producer,
	userRepo repository.UserRepository,
	walletRepo repository.WalletRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
) usecase.AuthUseCase {
	return &authUseCase{
		db:               db,
		log:              log,
		validate:         validate,
		jwtCfg:           jwtCfg,
		producer:         producer,
		userRepo:         userRepo,
		walletRepo:       walletRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (u *authUseCase) Register(ctx context.Context, req *model.RegisterRequest) (*model.UserResponse, error) {
	if err := u.validate.StructCtx(ctx, req); err != nil {
		return nil, err
	}

	existing, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		u.log.WithError(err).Error("register: find by email")
		return nil, ErrInternalCheckEmail
	}
	if existing != nil {
		return nil, ErrEmailAlreadyTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		u.log.WithError(err).Error("register: hash password")
		return nil, ErrInternalProcessPassword
	}

	var user entity.User
	err = u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		user = entity.User{
			Email:        req.Email,
			PasswordHash: string(hash),
			Role:         "user",
		}
		if err := u.userRepo.Create(ctx, &user, tx); err != nil {
			u.log.WithError(err).Error("register: create user")
			return ErrInternalCreateUser
		}

		wallet := entity.Wallet{
			UserID:   user.ID,
			Currency: "IDR",
		}
		if err := u.walletRepo.Create(ctx, &wallet, tx); err != nil {
			u.log.WithError(err).Error("register: create wallet")
			return ErrInternalCreateWallet
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Enqueue verification email — fire and forget, don't block registration
	u.enqueueVerificationEmail(user.ID, user.Email)

	u.log.WithField("user_id", user.ID).Info("user registered")
	return converter.ToUserResponse(&user), nil
}

func (u *authUseCase) enqueueVerificationEmail(userID, email string) {
	token, err := u.generateVerifyToken(userID, email)
	if err != nil {
		u.log.WithError(err).Error("register: generate verify token")
		return
	}

	verifyURL := fmt.Sprintf("%s/api/v1/auth/verify-email?token=%s", u.jwtCfg.AppURL, token)
	payload := model.VerificationEmailPayload{
		UserID:          userID,
		Email:           email,
		VerificationURL: verifyURL,
	}

	if err := u.producer.Enqueue(model.TypeEmailVerification, payload, asynq.Queue("default")); err != nil {
		u.log.WithError(err).Error("register: enqueue verification email")
	}
}

func (u *authUseCase) VerifyEmail(ctx context.Context, token string) error {
	claims, err := u.parseVerifyToken(token)
	if err != nil {
		return ErrInvalidVerifyToken
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return ErrInvalidVerifyToken
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ErrInvalidVerifyToken
	}

	if user.EmailVerified {
		return nil // idempotent — already verified
	}

	if err := u.userRepo.MarkEmailVerified(ctx, userID); err != nil {
		u.log.WithError(err).Error("verify email: mark verified")
		return ErrInternalCheckEmail
	}

	u.log.WithField("user_id", userID).Info("email verified")
	return nil
}

func (u *authUseCase) ResendVerification(ctx context.Context, email string) error {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// Don't reveal whether email exists — return success either way
		return nil
	}

	if user.EmailVerified {
		return nil
	}

	u.enqueueVerificationEmail(user.ID, user.Email)
	return nil
}

func (u *authUseCase) Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error) {
	if err := u.validate.StructCtx(ctx, req); err != nil {
		return nil, err
	}

	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredential
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredential
	}

	if !user.EmailVerified {
		return nil, ErrEmailNotVerified
	}

	accessToken, err := u.generateAccessToken(user)
	if err != nil {
		u.log.WithError(err).Error("login: generate access token")
		return nil, ErrInternalSaveToken
	}

	rawRefresh, hashRefresh, err := generateRefreshToken()
	if err != nil {
		u.log.WithError(err).Error("login: generate refresh token")
		return nil, ErrInternalSaveToken
	}

	refreshTTL, err := time.ParseDuration(u.jwtCfg.RefreshTTL)
	if err != nil {
		refreshTTL = 7 * 24 * time.Hour
	}

	rt := &entity.RefreshToken{
		UserID:    user.ID,
		TokenHash: hashRefresh,
		ExpiresAt: time.Now().Add(refreshTTL),
	}
	if err := u.refreshTokenRepo.Create(ctx, rt, nil); err != nil {
		u.log.WithError(err).Error("login: save refresh token")
		return nil, ErrInternalSaveToken
	}

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
	}, nil
}

func (u *authUseCase) generateAccessToken(user *entity.User) (string, error) {
	ttl, err := time.ParseDuration(u.jwtCfg.AccessTTL)
	if err != nil {
		ttl = 15 * time.Minute
	}

	claims := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Email: user.Email,
		Role:  user.Role,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(u.jwtCfg.Secret))
}

// generateVerifyToken creates a short-lived signed JWT for email verification.
func (u *authUseCase) generateVerifyToken(userID, email string) (string, error) {
	claims := jwt.MapClaims{
		"sub":    userID,
		"email":  email,
		"action": "verify_email",
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
		"iat":    time.Now().Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(u.jwtCfg.Secret))
}

func (u *authUseCase) parseVerifyToken(token string) (jwt.MapClaims, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(u.jwtCfg.Secret), nil
	})
	if err != nil || !parsed.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok || claims["action"] != "verify_email" {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func generateRefreshToken() (raw, hash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return
	}
	raw = hex.EncodeToString(b)
	sum := sha256.Sum256([]byte(raw))
	hash = hex.EncodeToString(sum[:])
	return
}

func (u *authUseCase) Refresh(ctx context.Context, refreshToken string) (*model.AuthResponse, error) {
	return nil, ErrNotImplemented
}

func (u *authUseCase) ForgotPassword(ctx context.Context, req *model.ForgotPasswordRequest) error {
	return ErrNotImplemented
}

func (u *authUseCase) ResetPassword(ctx context.Context, req *model.ResetPasswordRequest) error {
	return ErrNotImplemented
}

func (u *authUseCase) Logout(ctx context.Context, userID string, jti string) error {
	return ErrNotImplemented
}
