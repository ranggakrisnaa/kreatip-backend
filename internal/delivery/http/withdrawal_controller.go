package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kreatip/kreatip-backend/internal/usecase"
	"github.com/sirupsen/logrus"
)

type WithdrawalController struct {
	UseCase usecase.WithdrawalUseCase
	Log     *logrus.Logger
}

func NewWithdrawalController(uc usecase.WithdrawalUseCase, log *logrus.Logger) *WithdrawalController {
	return &WithdrawalController{UseCase: uc, Log: log}
}

// POST /api/v1/me/withdrawals
func (ctrl *WithdrawalController) Request(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// GET /api/v1/me/withdrawals
func (ctrl *WithdrawalController) ListMy(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// DELETE /api/v1/me/withdrawals/:id  (cancel before approved)
func (ctrl *WithdrawalController) Cancel(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// POST /api/v1/me/bank-accounts
func (ctrl *WithdrawalController) CreateBankAccount(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// GET /api/v1/me/bank-accounts
func (ctrl *WithdrawalController) ListBankAccounts(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// DELETE /api/v1/me/bank-accounts/:id
func (ctrl *WithdrawalController) DeleteBankAccount(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}
