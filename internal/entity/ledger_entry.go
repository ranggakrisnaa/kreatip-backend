package entity

import "time"

// LedgerDirection is either debit or credit.
type LedgerDirection string

const (
	LedgerDebit  LedgerDirection = "debit"
	LedgerCredit LedgerDirection = "credit"
)

// LedgerAccount names the virtual account being debited/credited.
type LedgerAccount string

const (
	AccountCreatorBalance    LedgerAccount = "creator_balance"
	AccountPendingWithdrawal LedgerAccount = "pending_withdrawal"
	AccountPlatformRevenue   LedgerAccount = "platform_revenue"
	AccountPaymentGateway    LedgerAccount = "payment_gateway"
	AccountPGFeeExpense      LedgerAccount = "pg_fee_expense"
)

type LedgerEntry struct {
	ID          string          `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	WalletID    string          `gorm:"type:uuid;not null;index"`
	Account     LedgerAccount   `gorm:"type:varchar(50);not null"`
	Direction   LedgerDirection `gorm:"type:varchar(10);not null"` // debit | credit
	Amount      int64           `gorm:"not null"`
	RefType     string          `gorm:"type:varchar(30);not null;index"` // donation | withdrawal | adjustment
	RefID       string          `gorm:"type:uuid;not null;index"`
	Description string          `gorm:"type:varchar(255)"`
	OccurredAt  time.Time       `gorm:"not null;index"`
}

func (LedgerEntry) TableName() string { return "ledger_entries" }
