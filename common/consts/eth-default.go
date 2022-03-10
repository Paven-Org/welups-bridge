package consts

type EthDefaultInfo struct {
	LastScannedBlock int64 `db:"column:eth_last_scan_block"`
}

type IEthInfoRepo interface {
	Get() (*EthDefaultInfo, error)
	Create(info *EthDefaultInfo) error
	Update(info *EthDefaultInfo) error
}
