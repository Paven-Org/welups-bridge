package dao

import (
	"bridge/micros/weleth/model"
	"bridge/service-managers/logger"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type IWelCashoutEthTransDAO interface {
	CreateWelCashoutEthTrans(t *model.WelCashoutEthTrans) (int64, error)

	UpdateWelCashoutEthTx(t *model.WelCashoutEthTrans) error

	SelectTransByWithdrawTxHash(txHash string) (*model.WelCashoutEthTrans, error)
	SelectTransByDisperseTxHash(txHash string) ([]*model.WelCashoutEthTrans, error)
	SelectTransById(id string) (*model.WelCashoutEthTrans, error)

	SelectTransByDisperseTxHashEthAddrAmount(txHash, ethWalletAddr, amount string) ([]*model.WelCashoutEthTrans, error)
	SelectTrans(sender, receiver, status string) ([]model.WelCashoutEthTrans, error)
}

// sort of a locator for DAOs
type welCashoutEthTransDAO struct {
	db *sqlx.DB
}

func (w *welCashoutEthTransDAO) CreateWelCashoutEthTrans(t *model.WelCashoutEthTrans) (int64, error) {
	log := logger.Get()
	tx, err := w.db.Beginx()
	if err != nil {
		log.Err(err).Msg("Can't start transaction")
		return -1, err
	}

	qCreate := tx.
		Rebind(
			`INSERT INTO wel_cashout_eth_trans(
				eth_disperse_tx_hash,
				wel_withdraw_tx_hash,
				eth_token_addr,
				wel_token_addr,
				network_id,
				eth_wallet_addr,
				wel_wallet_addr,
				total,
				amount,
				commission_fee,
				cashout_status,
				disperser_status) VALUES (?,?,?,?,?,?,?,?,?,?,?,?) RETURNING id`)
	var id int64
	err = tx.
		Get(&id,
			qCreate,
			t.EthDisperseTxHash,
			t.WelWithdrawTxHash,
			t.EthTokenAddr,
			t.WelTokenAddr,
			t.NetworkID,
			t.EthWalletAddr,
			t.WelWalletAddr,
			t.Total,
			t.Amount,
			t.CommissionFee,
			t.CashoutStatus,
			t.DisperseStatus)

	if err != nil {
		log.Err(err).Msgf("Error while inserting WelCashoutEth tx with eth tx hash %s", t.EthDisperseTxHash)
		tx.Rollback()
		return id, err
	}

	qUpdateTx2Treasury := tx.Rebind(`UPDATE tx_to_treasury SET status='isCashin' WHERE tx_id = ?`)
	_, err = tx.Exec(qUpdateTx2Treasury, t.EthDisperseTxHash)
	if err != nil {
		log.Err(err).Msgf("Error while inserting WelCashoutEth tx with eth tx hash %s", t.EthDisperseTxHash)
		tx.Rollback()
		return id, err
	}
	tx.Commit()

	return id, nil
}

func (w *welCashoutEthTransDAO) UpdateWelCashoutEthTx(t *model.WelCashoutEthTrans) error {
	db := w.db
	log := logger.Get()

	q := db.
		Rebind(
			`UPDATE wel_cashout_eth_trans SET
		    eth_tx_hash = ?,
		    wel_issue_tx_hash = ?,
		    eth_token_addr = ?,
		    wel_token_addr = ?,
		    network_id = ?,
		    eth_wallet_addr = ?,
		    wel_wallet_addr = ?,
		    amount = ?,
		    commission_fee = ?,
		    cashout_status = ?, 
		    disperse_status = ? 
		    WHERE id = ?`)
	_, err := db.
		Exec(q,
			t.EthDisperseTxHash,
			t.WelWithdrawTxHash,
			t.EthTokenAddr,
			t.WelTokenAddr,
			t.NetworkID,
			t.EthWalletAddr,
			t.WelWalletAddr,
			t.Amount,
			t.CommissionFee,
			t.CashoutStatus,
			t.DisperseStatus,
			t.ID)

	if err != nil {
		log.Err(err).Msgf("Error while inserting WelCashoutEth tx with eth tx hash %s", t.EthDisperseTxHash)
		return err
	}
	return nil
}

func (w *welCashoutEthTransDAO) SelectTransByWithdrawTxHash(txHash string) (*model.WelCashoutEthTrans, error) {
	var t = &model.WelCashoutEthTrans{}
	err := w.db.Get(t, "SELECT * FROM wel_cashout_eth_trans WHERE eth_tx_hash = $1", txHash)
	return t, err
}

func (w *welCashoutEthTransDAO) SelectTransByDisperseTxHash(txHash string) ([]*model.WelCashoutEthTrans, error) {
	var txs = []*model.WelCashoutEthTrans{}
	err := w.db.Select(txs, "SELECT * FROM wel_cashout_eth_trans WHERE wel_issue_tx_hash = $1", txHash)
	return txs, err
}

func (w *welCashoutEthTransDAO) SelectTransByDisperseTxHashEthAddrAmount(txHash, ethWalletAddr, amount string) ([]*model.WelCashoutEthTrans, error) {
	var txs = []*model.WelCashoutEthTrans{}
	err := w.db.Select(txs, "SELECT * FROM wel_cashout_eth_trans WHERE wel_issue_tx_hash = $1 AND eth_wallet_addr = $2 AND amount = $3", txHash, ethWalletAddr, amount)
	return txs, err
}

func (w *welCashoutEthTransDAO) SelectTransById(id string) (*model.WelCashoutEthTrans, error) {
	var t = &model.WelCashoutEthTrans{}
	err := w.db.Get(t, "SELECT * FROM wel_cashout_eth_trans WHERE id = $1", id)
	return t, err
}

func (w *welCashoutEthTransDAO) SelectTrans(sender, receiver, status string) ([]model.WelCashoutEthTrans, error) {
	// building query
	mapper := make(map[string]string)
	if len(sender) > 0 {
		mapper["eth_wallet_addr"] = sender
	}
	if len(receiver) > 0 {
		mapper["wel_wallet_addr"] = receiver
	}
	if len(status) > 0 {
		mapper["status"] = status
	}

	whereClauses := []string{}
	params := []interface{}{}
	for k, v := range mapper {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", k))
		params = append(params, v)
	}

	q := w.db.Rebind("SELECT * FROM wel_cashout_eth_trans WHERE " + strings.Join(whereClauses, " AND "))

	// querying...
	txs := []model.WelCashoutEthTrans{}
	err := w.db.Select(txs, q, params...)

	return txs, err
}

func MkWelCashoutEthTransDao(db *sqlx.DB) *welCashoutEthTransDAO {
	return &welCashoutEthTransDAO{
		db: db,
	}
}
