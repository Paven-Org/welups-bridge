package blockscan

import (
	"bridge/common/consts"

	"github.com/jmoiron/sqlx"
)

// sort of a locator for DAOs
type EthSysDAO struct {
	db *sqlx.DB
}

func (e *EthSysDAO) Get() (*consts.EthDefaultInfo, error) {
	var res = &consts.EthDefaultInfo{}
	err := e.db.Get(res, "SELECT eth_last_scan_block FROM wel_eth_sys LIMIT 1")
	return res, err
}

func (e *EthSysDAO) Create(info *consts.EthDefaultInfo) error {
	return nil
}

func (e *EthSysDAO) Update(info *consts.EthDefaultInfo) error {
	_, err := e.db.NamedExec(`UPDATE wel_eth_sys SET eth_last_scan_block = :last`,
		map[string]interface{}{
			"last": info.LastScannedBlock,
		})
	return err
}

func MkEthSysDao(db *sqlx.DB) *EthSysDAO {
	return &EthSysDAO{
		db: db,
	}
}
