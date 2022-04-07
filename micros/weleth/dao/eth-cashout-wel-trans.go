package dao

import (
	"bridge/micros/weleth/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type IEthCashoutWelTransDAO interface {
	CreateEthCashoutWelTrans(t *model.EthCashoutWelTrans) error

	UpdateDepositEthCashoutWelConfirmed(depositTxHash, ethWalletAddr, amount string) error

	UpdateClaimEthCashoutWel(id, claimTxHash, fee, status string) error

	SelectTransByDepositTxHash(txHash string) (*model.EthCashoutWelTrans, error)
	SelectTransById(id string) (*model.EthCashoutWelTrans, error)
}

// sort of a locator for DAOs
type ethCashoutWelTransDAO struct {
	db *sqlx.DB
}

func (w *ethCashoutWelTransDAO) CreateEthCashoutWelTrans(t *model.EthCashoutWelTrans) error {
	_, err := w.db.NamedExec(`INSERT INTO eth_cashout_wel_trans(id, wel_eth, deposit_tx_hash, wel_token_addr, eth_token_addr, eth_wallet_addr, wel_wallet_addr, network_id, amount, deposit_at, deposit_status) VALUES (:id, :wel_eth, :deposit_tx_hash, :wel_token_addr, :eth_token_addr, :eth_wallet_addr, :wel_wallet_addr, :network_id, :amount, :fee, :deposit_at, :deposit_status)`,
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

func (w *ethCashoutWelTransDAO) UpdateDepositEthCashoutWelConfirmed(depositTxHash, ethWalletAddr, amount string) error {
	_, err := w.db.NamedExec(`UPDATE eth_cashout_wel_trans SET deposit_status = :deposit_status, eth_wallet_addr = :eth_wallet_addr, amount = :amount WHERE deposit_tx_hash = :deposit_tx_hash`,
		map[string]interface{}{
			"deposit_status":  model.StatusSuccess,
			"eth_wallet_addr": ethWalletAddr,
			"amount":          amount,
			"deposit_tx_hash": depositTxHash,
		})
	return err
}

func (w *ethCashoutWelTransDAO) UpdateClaimEthCashoutWel(id, claimTxHash, fee, status string) error {
	_, err := w.db.NamedExec(`UPDATE eth_cashout_wel_trans SET claim_tx_hash = :claim_tx_hash, claim_status = :claim_status, fee = :fee WHERE id= :id`,
		map[string]interface{}{
			"claim_tx_hash": claimTxHash,
			"status":        status,
			"fee":           fee,
			"id":            id,
		})

	return err
}

func (w *ethCashoutWelTransDAO) SelectTransByDepositTxHash(txHash string) (*model.EthCashoutWelTrans, error) {
	var t = &model.EthCashoutWelTrans{}
	err := w.db.Get(t, "SELECT * FROM eth_cashout_wel_trans WHERE deposit_tx_hash = $1", txHash)
	return t, err
}

func (w *ethCashoutWelTransDAO) SelectTransById(id string) (*model.EthCashoutWelTrans, error) {
	var t = &model.EthCashoutWelTrans{}
	err := w.db.Get(t, "SELECT * FROM eth_cashout_wel_trans WHERE id = $1", id)
	return t, err
}

func MkEthCashoutWelTransDao(db *sqlx.DB) *ethCashoutWelTransDAO {
	return &ethCashoutWelTransDAO{
		db: db,
	}
}
