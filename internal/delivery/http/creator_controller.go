package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kreatip/kreatip-backend/internal/usecase"
	"github.com/sirupsen/logrus"
)

type CreatorController struct {
	UseCase usecase.CreatorUseCase
	Log     *logrus.Logger
}

func NewCreatorController(uc usecase.CreatorUseCase, log *logrus.Logger) *CreatorController {
	return &CreatorController{UseCase: uc, Log: log}
}

// GET /api/v1/me/profile
func (ctrl *CreatorController) GetMyProfile(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// PUT /api/v1/me/profile
func (ctrl *CreatorController) UpdateMyProfile(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// POST /api/v1/me/avatar
func (ctrl *CreatorController) UploadAvatar(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// GET /api/v1/creators/:username
func (ctrl *CreatorController) GetPublicProfile(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// GET /api/v1/me/alert/token
func (ctrl *CreatorController) GetAlertToken(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// POST /api/v1/me/alert/token  (rotate)
func (ctrl *CreatorController) RotateAlertToken(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// PUT /api/v1/me/alert/settings
func (ctrl *CreatorController) SaveAlertSettings(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}
