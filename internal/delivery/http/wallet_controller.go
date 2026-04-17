package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kreatip/kreatip-backend/internal/usecase"
	"github.com/sirupsen/logrus"
)

type WalletController struct {
	UseCase usecase.WalletUseCase
	Log     *logrus.Logger
}

func NewWalletController(uc usecase.WalletUseCase, log *logrus.Logger) *WalletController {
	return &WalletController{UseCase: uc, Log: log}
}

// GET /api/v1/me/wallet
func (ctrl *WalletController) GetWallet(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// GET /api/v1/me/ledger
func (ctrl *WalletController) ListLedger(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}
