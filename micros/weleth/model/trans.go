package model

import (
	"database/sql"
	"fmt"
	"time"
)

var (
	EthTokenFromWel map[string]string = map[string]string{
		"0x4272ffC0682d68aCF5eEbD2ABFDc38d721BCF55a": "W9yD14Nj9j7xAB4dbGeiX9h8unkKHxuTtb",
	}
	WelTokenFromEth map[string]string = map[string]string{
		"W9yD14Nj9j7xAB4dbGeiX9h8unkKHxuTtb": "0x4272ffC0682d68aCF5eEbD2ABFDc38d721BCF55a",
	}
)

const (
	StatusSuccess = "confirmed"
	StatusUnknown = "unconfirmed"
	StatusPending = "pending"

	RequestDoubleClaimed = "doubleclaimed"
	RequestExpired       = "expired"
	RequestSuccess       = "success"
)

var (
	ErrAlreadyClaimed     = fmt.Errorf("Already claimed")
	ErrRequestPending     = fmt.Errorf("Request pending")
	ErrUnrecognizedStatus = fmt.Errorf("Unrecognized transaction status")
)

type ClaimRequest struct {
	Txid   int64  `db:"tx_id"`
	ReqID  string `db:"request_id"`
	Status string `db:"status"`
}

type WelEthEvent = WelCashinEthTrans
type WelCashinEthTrans struct {
	ID    int64  `json:"id,omitempty" db:"id,omitempty"`
	ReqID string `json:"request_id,omitempty" db:"request_id,omitempty"`

	DepositTxHash string `json:"deposit_tx_hash" db:"deposit_tx_hash"`
	ClaimTxHash   string `json:"claim_tx_hash" db:"claim_tx_hash"`

	WelTokenAddr string `json:"wel_token_addr" db:"wel_token_addr,omitempty"`
	EthTokenAddr string `json:"eth_token_addr" db:"eth_token_addr,omitempty"`

	WelWalletAddr string `json:"wel_wallet_addr,omitempty" db:"wel_wallet_addr"`
	EthWalletAddr string `json:"eth_wallet_addr" db:"eth_wallet_addr"`

	NetworkID string `json:"network_id" db:"network_id"`

	Amount string `json:"amount" db:"amount"`

	Fee string `json:"fee" db:"fee"`

	DepositStatus string `json:"deposit_status" db:"deposit_status"`
	ClaimStatus   string `json:"claim_status" db:"claim_status"`

	DepositAt time.Time    `json:"deposit_at" db:"deposit_at"`
	ClaimAt   sql.NullTime `json:"claim_at" db:"claim_at"`
}

type EthWelEvent = EthCashoutWelTrans
type EthCashoutWelTrans struct {
	ID    int64  `json:"id,omitempty" db:"id,omitempty"`
	ReqID string `json:"request_id,omitempty" db:"request_id,omitempty"`

	DepositTxHash string `json:"deposit_tx_hash" db:"deposit_tx_hash"`
	ClaimTxHash   string `json:"claim_tx_hash" db:"claim_tx_hash"`

	WelTokenAddr string `json:"wel_token_addr" db:"wel_token_addr,omitempty"`
	EthTokenAddr string `json:"eth_token_addr" db:"eth_token_addr,omitempty"`

	WelWalletAddr string `json:"wel_wallet_addr,omitempty" db:"wel_wallet_addr"`
	EthWalletAddr string `json:"eth_wallet_addr" db:"eth_wallet_addr"`

	NetworkID string `json:"network_id" db:"network_id"`

	Amount string `json:"amount" db:"amount"`

	Fee string `json:"fee" db:"fee"`

	DepositStatus string `json:"deposit_status" db:"deposit_status"`
	ClaimStatus   string `json:"claim_status" db:"claim_status"`

	DepositAt time.Time    `json:"deposit_at" db:"deposit_at"`
	ClaimAt   sql.NullTime `json:"claim_at" db:"claim_at"`
}
