package impl

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/kreatip/kreatip-backend/internal/entity"
	"github.com/kreatip/kreatip-backend/internal/model"
	"github.com/kreatip/kreatip-backend/internal/model/converter"
	"github.com/kreatip/kreatip-backend/internal/repository"
	"github.com/kreatip/kreatip-backend/internal/usecase"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authUseCase struct {
	db         *gorm.DB
	log        *logrus.Logger
	validate   *validator.Validate
	userRepo   repository.UserRepository
	walletRepo repository.WalletRepository
}

func NewAuthUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	userRepo repository.UserRepository,
	walletRepo repository.WalletRepository,
) usecase.AuthUseCase {
	return &authUseCase{
		db:         db,
		log:        log,
		validate:   validate,
		userRepo:   userRepo,
		walletRepo: walletRepo,
	}
}

func (u *authUseCase) Register(ctx context.Context, req *model.RegisterRequest) (*model.UserResponse, error) {
	// 1. Validate request
	if err := u.validate.StructCtx(ctx, req); err != nil {
		return nil, err
	}

	// 2. Check email not taken
	existing, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		u.log.WithError(err).Error("register: find by email")
		return nil, ErrInternalCheckEmail
	}
	if existing != nil {
		return nil, ErrEmailAlreadyTaken
	}

	// 3. Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		u.log.WithError(err).Error("register: hash password")
		return nil, ErrInternalProcessPassword
	}

	// 4. Persist user + wallet in one transaction
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

	u.log.WithField("user_id", user.ID).Info("user registered")
	return converter.ToUserResponse(&user), nil
}

func (u *authUseCase) VerifyEmail(ctx context.Context, token string) error {
	return ErrNotImplemented
}

func (u *authUseCase) Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error) {
	return nil, ErrNotImplemented
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
