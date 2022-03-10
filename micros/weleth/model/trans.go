package model

import "time"

const (
	StatusSuccess = "Verified"
	StatusUnknown = "Verifying"
)

type DoneDepositEvent struct {
	ID           string    `json:"id,omitempty" db:"column:id"`
	DepositID    string    `json:"deposit_id,omitempty" db:"column:deposit_id"`
	TxHash       string    `json:"tx_hash,omitempty" db:"tx_hash"`
	WelTokenAddr string    `json:"wel_token_addr" db:"wel_token_addr,omitempty"`
	FromAddr     string    `json:"from_addr,omitempty" db:"from_addr,omitempty"`
	EthTokenAddr string    `json:"eth_token_addr" db:"eth_token_addr,omitempty"`
	NetworkID    uint64    `json:"network_id" db:"network_id,omitempty"`
	Amount       string    `json:"amount,omitempty" db:"amount"`
	Fee          string    `json:"fee,omitempty" db:"fee"`
	Status       string    `json:"status,omitempty" db:"status,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty" db:"created_at,omitempty"`
}

type DoneClaimEvent struct {
	ID          string `json:"id,omitempty" db:"column:id"`
	TxHash      string `json:"tx_hash,omitempty" db:"tx_hash,omitempty"`
	ClaimedAddr string `json:"claimed_addr,omitempty" db:"claimed_addr,omitempty"`
	Amount      string `json:"amount,omitempty" db:"amount,omitempty"`
	Status      string `json:"status,omitempty" db:"decimal,omitempty"`
	CreatedAt   time.Time
}
