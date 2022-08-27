package dao

import (
	"bridge/micros/weleth/model"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type IWelCashinEthTransDAO interface {
	CreateWelCashinEthTrans(t *model.WelCashinEthTrans) error

	UpdateDepositWelCashinEthConfirmed(depositTxHash, welWalletAddr, amount, fee string) error

	UpdateClaimWelCashinEth(id int64, reqID, reqStatus, claimTxHash, status string) error

	SelectTransByDepositTxHash(txHash string) (*model.WelCashinEthTrans, error)
	SelectTransById(id string) (*model.WelCashinEthTrans, error)
	SelectTrans(sender, receiver, status string, offset, size uint64) ([]model.WelCashinEthTrans, error)

	CreateClaimRequest(requestID string, txID int64, status string, expiredAt time.Time) error
	SelectTransByRqId(rid string) (*model.WelCashinEthTrans, error)
	UpdateClaimRequest(reqID, status string) error
	GetClaimRequest(reqID string) (*model.ClaimRequest, error)
}

// sort of a locator for DAOs
type welCashinEthTransDAO struct {
	db *sqlx.DB
}

func (w *welCashinEthTransDAO) CreateWelCashinEthTrans(t *model.WelCashinEthTrans) error {
	_, err := w.db.NamedExec(`INSERT INTO wel_cashin_eth_trans(deposit_tx_hash, wel_token_addr, eth_token_addr,eth_wallet_addr, wel_wallet_addr, network_id, amount, fee, deposit_at, deposit_status) VALUES (:deposit_tx_hash, :wel_token_addr, :eth_token_addr, :eth_wallet_addr, :wel_wallet_addr, :network_id, :amount, :fee, :deposit_at, :deposit_status)`,
		map[string]interface{}{
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

func (w *welCashinEthTransDAO) UpdateClaimWelCashinEth(id int64, reqID, reqStatus, claimTxHash, status string) error {
	tx, err := w.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(`UPDATE wel_cashin_eth_trans SET request_id = :request_id, claim_tx_hash = :claim_tx_hash, claim_status = :claim_status WHERE id= :id`,
		map[string]interface{}{
			"claim_tx_hash": claimTxHash,
			"request_id":    reqID,
			"claim_status":  status,
			"id":            id,
		})
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("UPDATE wel_cashin_eth_req SET status = $1 WHERE request_id = $2", reqStatus, reqID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (w *welCashinEthTransDAO) UpdateClaimRequest(reqID, status string) error {
	_, err := w.db.Exec("UPDATE wel_cashin_eth_req SET status = $1 WHERE request_id = $2", status, reqID)
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

func (w *welCashinEthTransDAO) CreateClaimRequest(requestID string, txID int64, status string, expiredAt time.Time) error {
	tx, err := w.db.Beginx()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO wel_cashin_eth_req(request_id, tx_id, status, expired_at) VALUES ($1, $2, $3, $4)`, requestID, txID, status, expiredAt)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.NamedExec(`UPDATE wel_cashin_eth_trans SET request_id = :request_id, claim_status = :status WHERE id= :id`,
		map[string]interface{}{
			"request_id": requestID,
			"status":     status,
			"id":         txID,
		})
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (w *welCashinEthTransDAO) SelectTransByRqId(rid string) (*model.WelCashinEthTrans, error) {
	var t = &model.WelCashinEthTrans{}
	err := w.db.Get(t,
		`SELECT t.* FROM 
					wel_cashin_eth_trans as t 
					JOIN wel_cashin_eth_req as r 
					ON t.id = r.tx_id 
					WHERE r.request_id = $1`, rid)
	return t, err
}

func (w *welCashinEthTransDAO) SelectTrans(sender, receiver, status string, offset, size uint64) ([]model.WelCashinEthTrans, error) {
	// building query
	mapper := make(map[string]string)
	if len(sender) > 0 {
		mapper["wel_wallet_addr"] = sender
	}
	if len(receiver) > 0 {
		mapper["eth_wallet_addr"] = receiver
	}
	if len(status) > 0 {
		mapper["deposit_status"] = status
	}

	whereClauses := []string{}
	params := []interface{}{}
	for k, v := range mapper {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", k))
		params = append(params, v)
	}

	limitClause := ""
	if size > 0 {
		limitClause = fmt.Sprintf(" OFFSET %d LIMIT %d", offset, size)
	}

	q := "SELECT * FROM wel_cashin_eth_trans"
	if len(whereClauses) > 0 {
		q = w.db.Rebind(q + " WHERE " + strings.Join(whereClauses, " AND "))
	}
	q = q + limitClause

	// querying...
	txs := []model.WelCashinEthTrans{}
	err := w.db.Select(&txs, q, params...)

	return txs, err
}

func (w *welCashinEthTransDAO) GetClaimRequest(reqID string) (*model.ClaimRequest, error) {
	var req = &model.ClaimRequest{}
	err := w.db.Get(req, `SELECT * FROM wel_cashin_eth_req WHERE request_id = $1`, reqID)
	return req, err
}

func MkWelCashinEthTransDao(db *sqlx.DB) *welCashinEthTransDAO {
	return &welCashinEthTransDAO{
		db: db,
	}
}
