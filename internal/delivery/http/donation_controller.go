package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kreatip/kreatip-backend/internal/usecase"
	"github.com/sirupsen/logrus"
)

type DonationController struct {
	UseCase usecase.DonationUseCase
	Log     *logrus.Logger
}

func NewDonationController(uc usecase.DonationUseCase, log *logrus.Logger) *DonationController {
	return &DonationController{UseCase: uc, Log: log}
}

// POST /api/v1/donations
func (ctrl *DonationController) Create(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// GET /api/v1/donations/:id
func (ctrl *DonationController) GetByID(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// GET /api/v1/me/donations
func (ctrl *DonationController) ListMy(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}
