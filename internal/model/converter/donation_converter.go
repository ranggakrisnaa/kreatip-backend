package converter

import (
	"github.com/kreatip/kreatip-backend/internal/entity"
	"github.com/kreatip/kreatip-backend/internal/model"
)

func ToDonationResponse(d *entity.Donation) *model.DonationResponse {
	resp := &model.DonationResponse{
		ID:        d.ID,
		CreatorID: d.CreatorID,
		DonorName: d.DonorName,
		Amount:    d.Amount,
		NetAmount: d.NetAmount,
		Message:   d.Message,
		Status:    string(d.Status),
		PaidAt:    d.PaidAt,
		CreatedAt: d.CreatedAt,
	}
	if d.Payment != nil {
		resp.CheckoutURL = d.Payment.CheckoutURL
	}
	return resp
}
