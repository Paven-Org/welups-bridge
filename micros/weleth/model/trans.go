package model

import "time"

const (
	StatusSuccess = "confirmed"
	StatusUnknown = "unconfirmed"
)

type WelEthEvent = WelCashinEthTrans
type WelCashinEthTrans struct {
	ID string `json:"id,omitempty" db:"id,omitempty"`

	// if return = true -> it is the request from wel -> eth, else it is the request from eth -> wel
	WelEth bool `json:"wel_eth" db:"wel_eth"`

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

	DepositAt time.Time `json:"deposit_at" db:"deposit_at"`
	ClaimAt   time.Time `json:"claim_at" db:"claim_at"`
}

type EthWelEvent = EthCashoutWelTrans
type EthCashoutWelTrans struct {
	ID string `json:"id,omitempty" db:"id,omitempty"`

	// if return = true -> it is the request from wel -> eth, else it is the request from eth -> wel
	WelEth bool `json:"wel_eth" db:"wel_eth"`

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

	DepositAt time.Time `json:"deposit_at" db:"deposit_at"`
	ClaimAt   time.Time `json:"claim_at" db:"claim_at"`
}
