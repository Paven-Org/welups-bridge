package model

import "time"

const (
	StatusSuccess = "Verified"
	StatusUnknown = "Verifying"
)

type DoneDepositEvent struct {
	ID        string    `json:"id,omitempty" db:"column:id"`
	DepositID string    `json:"deposit_id,omitempty" db:"column:deposit_id"`
	TxHash    string    `json:"tx_hash,omitempty" db:"tx_hash,omitempty"`
	FromAddr  string    `json:"from_addr,omitempty" db:"from_addr,omitempty"`
	Amount    string    `json:"amount,omitempty" db:"amount,omitempty"`
	Decimal   uint      `json:"decimal,omitempty" db:"decimal,omitempty"`
	Status    string    `json:"status,omitempty" db:"decimal,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at,omitempty"`
}

type DoneClaimEvent struct {
	ID          string `json:"id,omitempty" db:"column:id"`
	TxHash      string `json:"tx_hash,omitempty" db:"tx_hash,omitempty"`
	ClaimedAddr string `json:"claimed_addr,omitempty" db:"claimed_addr,omitempty"`
	Amount      string `json:"amount,omitempty" db:"amount,omitempty"`
	Decimal     uint   `json:"decimal,omitempty" db:"decimal,omitempty"`
	Status      string `json:"status,omitempty" db:"decimal,omitempty"`
	CreatedAt   time.Time
}
