package dao

import (
	"bridge/micros/weleth/model"

	"github.com/jmoiron/sqlx"
)

type IWelCashinEthTransDAO interface {
	CreateWelCashinEthTrans(t *model.WelCashinEthTrans) error

	UpdateDepositWelCashinEthConfirmed(depositTxHash, welWalletAddr, amount, fee string) error

	UpdateClaimWelCashinEth(id, claimTxHash, status string) error

	SelectTransByDepositTxHash(txHash string) (*model.WelCashinEthTrans, error)
	SelectTransById(id string) (*model.WelCashinEthTrans, error)
}

// sort of a locator for DAOs
type welCashinEthTransDAO struct {
	db *sqlx.DB
}

func (w *welCashinEthTransDAO) CreateWelCashinEthTrans(t *model.WelCashinEthTrans) error {
	_, err := w.db.NamedExec(`INSERT INTO wel_cashin_eth_trans(id, wel_eth, deposit_tx_hash, wel_token_addr, eth_token_addr,eth_wallet_addr, wel_wallet_addr, network_id, amount, fee, deposit_at, deposit_status) VALUES (:id, :wel_eth, :deposit_tx_hash, :wel_token_addr, :eth_token_addr, :eth_wallet_addr, :wel_wallet_addr, :network_id, :amount, :fee, :deposit_at, :deposit_status)`,
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

func (w *welCashinEthTransDAO) UpdateDepositWelCashinEthConfirmed(depositTxHash, welWalletAddr, amount, fee string) error {
	_, err := w.db.NamedExec(`UPDATE wel_cashin_eth_trans SET deposit_status = :deposit_status, wel_wallet_addr = :wel_wallet_addr, amount = :amount, fee = :fee WHERE deposit_tx_hash = :deposit_tx_hash`,
		map[string]interface{}{
			"deposit_status":  model.StatusSuccess,
			"wel_wallet_addr": welWalletAddr,
			"amount":          amount,
			"fee":             fee,
			"deposit_tx_hash": depositTxHash,
		})
	return err
}

func (w *welCashinEthTransDAO) UpdateClaimWelCashinEth(id, claimTxHash, status string) error {
	_, err := w.db.NamedExec(`UPDATE wel_cashin_eth_trans SET claim_tx_hash = :claim_tx_hash, claim_status = :claim_status WHERE id= :id`,
		map[string]interface{}{
			"claim_tx_hash": claimTxHash,
			"status":        status,
			"id":            id,
		})

	return err
}

func (w *welCashinEthTransDAO) SelectTransByDepositTxHash(txHash string) (*model.WelCashinEthTrans, error) {
	var t = &model.WelCashinEthTrans{}
	err := w.db.Get(t, "SELECT * FROM wel_cashin_eth_trans WHERE deposit_tx_hash = $1", txHash)
	return t, err
}

func (w *welCashinEthTransDAO) SelectTransById(id string) (*model.WelCashinEthTrans, error) {
	var t = &model.WelCashinEthTrans{}
	err := w.db.Get(t, "SELECT * FROM wel_cashin_eth_trans WHERE id = $1", id)
	return t, err
}

func MkWelCashinEthTransDao(db *sqlx.DB) *welCashinEthTransDAO {
	return &welCashinEthTransDAO{
		db: db,
	}
}
