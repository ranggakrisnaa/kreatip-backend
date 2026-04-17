package model

import "time"

type WalletResponse struct {
	ID               string     `json:"id"`
	Currency         string     `json:"currency"`
	BalanceAvailable int64      `json:"balance_available"`
	BalancePending   int64      `json:"balance_pending"`
	LastReconciledAt *time.Time `json:"last_reconciled_at,omitempty"`
}

type LedgerEntryResponse struct {
	ID          string    `json:"id"`
	Account     string    `json:"account"`
	Direction   string    `json:"direction"`
	Amount      int64     `json:"amount"`
	RefType     string    `json:"ref_type"`
	RefID       string    `json:"ref_id"`
	Description string    `json:"description"`
	OccurredAt  time.Time `json:"occurred_at"`
}
