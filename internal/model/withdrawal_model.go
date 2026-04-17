package model

import "time"

type CreateWithdrawalRequest struct {
	BankAccountID string `json:"bank_account_id" validate:"required,uuid"`
	Amount        int64  `json:"amount"          validate:"required,min=1"`
	TOTPCode      string `json:"totp_code"       validate:"omitempty,len=6"`
}

type WithdrawalResponse struct {
	ID            string    `json:"id"`
	Amount        int64     `json:"amount"`
	NetAmount     int64     `json:"net_amount"`
	Status        string    `json:"status"`
	BankAccountID string    `json:"bank_account_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type CreateBankAccountRequest struct {
	Type              string `json:"type"               validate:"required,oneof=bank ewallet"`
	BankCode          string `json:"bank_code"          validate:"required"`
	AccountNumber     string `json:"account_number"     validate:"required,min=5,max=30"`
	AccountHolderName string `json:"account_holder_name" validate:"required,min=2,max=100"`
	IsDefault         bool   `json:"is_default"`
}

type BankAccountResponse struct {
	ID                string `json:"id"`
	Type              string `json:"type"`
	BankCode          string `json:"bank_code"`
	AccountNumber     string `json:"account_number"`
	AccountHolderName string `json:"account_holder_name"`
	IsVerified        bool   `json:"is_verified"`
	IsDefault         bool   `json:"is_default"`
}
