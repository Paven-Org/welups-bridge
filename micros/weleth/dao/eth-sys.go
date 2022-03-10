package dao

import (
	"bridge/common/consts"

	"github.com/jmoiron/sqlx"
)

// sort of a locator for DAOs
type ethSysDAO struct {
	db *sqlx.DB
}

func (e *ethSysDAO) Get() (*consts.EthDefaultInfo, error) {
	var res = &consts.EthDefaultInfo{}
	err := e.db.Select(res, "SELECT eth_last_scan_block FROM wel_eth_sys ORDER BY first_name ASC")
	return res, err
}

func (e *ethSysDAO) Create(info *consts.EthDefaultInfo) error {
	return nil
}

func (e *ethSysDAO) Update(info *consts.EthDefaultInfo) error {
	_, err := e.db.NamedExec(`UPDATE wel_eth_sys SET eth_last_scan_block = :last`,
		map[string]interface{}{
			"last": info.LastScannedBlock,
		})
	return err
}

func MkEthSysDao(db *sqlx.DB) *ethSysDAO {
	return &ethSysDAO{
		db: db,
	}
}
