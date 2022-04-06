package dao

import (
	"bridge/micros/weleth/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type IWelEthTransDAO interface {
	CreateWelEthTrans(t *model.WelEthEvent) error
	CreateEthWelTrans(t *model.WelEthEvent) error

	UpdateDepositWelEthConfirmed(depositTxHash, welWalletAddr, amount, fee string) error
	UpdateDepositEthWelConfirmed(depositTxHash, ethWalletAddr, amount string) error

	UpdateClaimWelEth(id, claimTxHash, status string) error
	UpdateClaimEthWel(id, claimTxHash, fee, status string) error

	SelectTransByDepositTxHash(txHash string) (*model.WelEthEvent, error)
	SelectTransById(id string) (*model.WelEthEvent, error)
}

// sort of a locator for DAOs
type welEthTransDAO struct {
	db *sqlx.DB
}

func (w *welEthTransDAO) CreateWelEthTrans(t *model.WelEthEvent) error {
	_, err := w.db.NamedExec(`INSERT INTO wel_eth_trans(id, wel_eth, deposit_tx_hash, wel_token_addr, eth_token_addr,eth_wallet_addr, wel_wallet_addr, network_id, amount, fee, deposit_at, deposit_status) VALUES (:id, :wel_eth, :deposit_tx_hash, :wel_token_addr, :eth_token_addr, :eth_wallet_addr, :wel_wallet_addr, :network_id, :amount, :fee, :deposit_at, :deposit_status)`,
		map[string]interface{}{
			"id":              t.ID,
			"wel_eth":         true,
			"deposit_tx_hash": t.DepositTxHash,
			"eth_wallet_addr": t.EthWalletAddr,
			"wel_wallet_addr": t.WelWalletAddr,
			"eth_token_addr":  t.EthTokenAddr,
			"wel_token_addr":  t.WelTokenAddr,
			"network_id":      t.NetworkID,
			"amount":          t.Amount,
			"fee":             t.Fee,
			"deposit_at":      t.DepositAt,
			"deposit_status":  t.DepositStatus,
		})
	return err
}

func (w *welEthTransDAO) CreateEthWelTrans(t *model.WelEthEvent) error {
	_, err := w.db.NamedExec(`INSERT INTO wel_eth_trans(id, wel_eth, deposit_tx_hash, wel_token_addr, eth_token_addr, eth_wallet_addr, wel_wallet_addr, network_id, amount, deposit_at, deposit_status) VALUES (:id, :wel_eth, :deposit_tx_hash, :wel_token_addr, :eth_token_addr, :eth_wallet_addr, :wel_wallet_addr, :network_id, :amount, :fee, :deposit_at, :deposit_status)`,
		map[string]interface{}{
			"id":              t.ID,
			"wel_eth":         false,
			"deposit_tx_hash": t.DepositTxHash,
			"wel_token_addr":  t.WelTokenAddr,
			"eth_token_addr":  t.EthTokenAddr,
			"eth_wallet_addr": t.EthWalletAddr,
			"wel_wallet_addr": t.EthWalletAddr,
			"network_id":      t.NetworkID,
			"amount":          t.Amount,
			"fee":             t.Fee,
			"deposit_at":      time.Now(),
			"deposit_status":  t.DepositStatus,
		})

	return err
}

func (w *welEthTransDAO) UpdateDepositWelEthConfirmed(depositTxHash, welWalletAddr, amount, fee string) error {
	_, err := w.db.NamedExec(`UPDATE wel_eth_trans SET deposit_status = :deposit_status, wel_wallet_addr = :wel_wallet_addr, amount = :amount, fee = :fee WHERE deposit_tx_hash = :deposit_tx_hash`,
		map[string]interface{}{
			"deposit_status":  model.StatusSuccess,
			"wel_wallet_addr": welWalletAddr,
			"amount":          amount,
			"fee":             fee,
			"deposit_tx_hash": depositTxHash,
		})
	return err
}

func (w *welEthTransDAO) UpdateDepositEthWelConfirmed(depositTxHash, ethWalletAddr, amount string) error {
	_, err := w.db.NamedExec(`UPDATE wel_eth_trans SET deposit_status = :deposit_status, eth_wallet_addr = :eth_wallet_addr, amount = :amount WHERE deposit_tx_hash = :deposit_tx_hash`,
		map[string]interface{}{
			"deposit_status":  model.StatusSuccess,
			"eth_wallet_addr": ethWalletAddr,
			"amount":          amount,
			"deposit_tx_hash": depositTxHash,
		})
	return err
}

func (w *welEthTransDAO) UpdateClaimWelEth(id, claimTxHash, status string) error {
	_, err := w.db.NamedExec(`UPDATE wel_eth_trans SET claim_tx_hash = :claim_tx_hash, claim_status = :claim_status WHERE id= :id`,
		map[string]interface{}{
			"claim_tx_hash": claimTxHash,
			"status":        status,
			"id":            id,
		})

	return err
}

func (w *welEthTransDAO) UpdateClaimEthWel(id, claimTxHash, fee, status string) error {
	_, err := w.db.NamedExec(`UPDATE wel_eth_trans SET claim_tx_hash = :claim_tx_hash, claim_status = :claim_status, fee = :fee WHERE id= :id`,
		map[string]interface{}{
			"claim_tx_hash": claimTxHash,
			"status":        status,
			"fee":           fee,
			"id":            id,
		})

	return err
}

func (w *welEthTransDAO) SelectTransByDepositTxHash(txHash string) (*model.WelEthEvent, error) {
	var t = &model.WelEthEvent{}
	err := w.db.Get(t, "SELECT * FROM wel_eth_trans WHERE deposit_tx_hash = $1", txHash)
	return t, err
}

func (w *welEthTransDAO) SelectTransById(id string) (*model.WelEthEvent, error) {
	var t = &model.WelEthEvent{}
	err := w.db.Get(t, "SELECT * FROM wel_eth_trans WHERE id = $1", id)
	return t, err
}

func MkWelEthTransDao(db *sqlx.DB) *welEthTransDAO {
	return &welEthTransDAO{
		db: db,
	}
}
