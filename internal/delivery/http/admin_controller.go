package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kreatip/kreatip-backend/internal/usecase"
	"github.com/sirupsen/logrus"
)

type AdminController struct {
	WithdrawalUC usecase.WithdrawalUseCase
	Log          *logrus.Logger
}

func NewAdminController(withdrawalUC usecase.WithdrawalUseCase, log *logrus.Logger) *AdminController {
	return &AdminController{WithdrawalUC: withdrawalUC, Log: log}
}

// GET /api/v1/admin/withdrawals
func (ctrl *AdminController) ListWithdrawals(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// POST /api/v1/admin/withdrawals/:id/approve
func (ctrl *AdminController) ApproveWithdrawal(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// POST /api/v1/admin/withdrawals/:id/reject
func (ctrl *AdminController) RejectWithdrawal(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// GET /api/v1/admin/users
func (ctrl *AdminController) ListUsers(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// POST /api/v1/admin/users/:id/suspend
func (ctrl *AdminController) SuspendUser(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}
