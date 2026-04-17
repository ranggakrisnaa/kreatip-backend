package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kreatip/kreatip-backend/internal/usecase"
	"github.com/sirupsen/logrus"
)

type WebhookController struct {
	PaymentUC    usecase.PaymentUseCase
	WithdrawalUC usecase.WithdrawalUseCase
	Log          *logrus.Logger
}

func NewWebhookController(paymentUC usecase.PaymentUseCase, withdrawalUC usecase.WithdrawalUseCase, log *logrus.Logger) *WebhookController {
	return &WebhookController{PaymentUC: paymentUC, WithdrawalUC: withdrawalUC, Log: log}
}

// POST /api/v1/webhooks/xendit  (donation payment)
func (ctrl *WebhookController) XenditPayment(c *fiber.Ctx) error {
	// TODO: verify X-CALLBACK-TOKEN, parse payload, call PaymentUC.HandleDonationWebhook
	return fiber.ErrNotImplemented
}

// POST /api/v1/webhooks/disbursement  (withdrawal payout)
func (ctrl *WebhookController) XenditDisbursement(c *fiber.Ctx) error {
	// TODO: verify signature, parse payload, call WithdrawalUC.HandleDisbursementWebhook
	return fiber.ErrNotImplemented
}
