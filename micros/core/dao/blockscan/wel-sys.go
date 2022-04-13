package blockscan

import (
	"bridge/common/consts"

	"github.com/jmoiron/sqlx"
)

// sort of a locator for DAOs
type WelSysDAO struct {
	db *sqlx.DB
}

func (e *WelSysDAO) Get() (*consts.WelDefaultInfo, error) {
	var res = &consts.WelDefaultInfo{}
	err := e.db.Get(res, "SELECT wel_last_scan_block FROM wel_eth_sys LIMIT 1")
	return res, err
}

func (e *WelSysDAO) Create(info *consts.WelDefaultInfo) error {
	return nil
}

func (e *WelSysDAO) Update(info *consts.WelDefaultInfo) error {
	_, err := e.db.NamedExec(`UPDATE wel_eth_sys SET wel_last_scan_block = :last`,
		map[string]interface{}{
			"last": info.LastScannedBlock,
		})
	return err
}

func MkWelSysDao(db *sqlx.DB) *WelSysDAO {
	return &WelSysDAO{
		db: db,
	}
}
