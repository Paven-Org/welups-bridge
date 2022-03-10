package dao

import (
	"bridge/micros/weleth/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type IWelTransDAO interface {
	CreateTrans(t *model.DoneDepositEvent) error
	UpdateVerified(txHash string) error
	UpdateClaimed(id, depositID string) error
	SelectTransByTxHash(txHash string) (*model.DoneDepositEvent, error)
}

// sort of a locator for DAOs
type welTransDAO struct {
	db *sqlx.DB
}

func (w *welTransDAO) CreateTrans(t *model.DoneDepositEvent) error {
	_, err := w.db.NamedExec(`INSERT INTO wel_eth_deposit_trans(id, deposit_id, tx_hash, from_addr, amount, decimal, status, created_at) VALUES (:id, :deposit_id, :tx_hash, :from_addr, :amount, :decimal, :status, :created_at)`,
		map[string]interface{}{
			"id":         t.ID,
			"deposit_id": t.DepositID,
			"tx_hash":    t.TxHash,
			"from_addr":  t.FromAddr,
			"amount":     t.Amount,
			"decimal":    t.Decimal,
			"status":     model.StatusUnknown,
			"created_at": time.Now(),
		})
	return err
}

func (w *welTransDAO) UpdateVerified(txHash string) error {
	_, err := w.db.NamedExec(`UPDATE wel_eth_deposit_trans SET status = :status WHERE tx_hash = :tx_hash`,
		map[string]interface{}{
			"status":  model.StatusSuccess,
			"tx_hash": txHash,
		})
	return err
}

func (w *welTransDAO) UpdateClaimed(id, depositID string) error {
	_, err := w.db.NamedExec(`UPDATE wel_eth_deposit_trans SET deposit_id = :deposit_id WHERE id = :id`,
		map[string]interface{}{
			"deposit_id": depositID,
			"id":         id,
		})
	return err
}

func (w *welTransDAO) SelectTransByTxHash(txHash string) (*model.DoneDepositEvent, error) {
	var t = &model.DoneDepositEvent{}
	err := w.db.Select(t, "SELECT wel_eth_deposit_trans FROM wel_eth_sys WHERE tx_hash = :tx_hash",
		map[string]interface{}{
			"tx_hash": txHash,
		})
	return t, err
}

func MkWelTransDao(db *sqlx.DB) *welTransDAO {
	return &welTransDAO{
		db: db,
	}
}
