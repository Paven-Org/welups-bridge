package consts

type EthDefaultInfo struct {
	LastScannedBlock int64 `gorm:"column:last_scanned_block"`
}

type IEthInfoRepo interface {
	Get() (*EthDefaultInfo, error)
	Create(info *EthDefaultInfo) error
	Update(info *EthDefaultInfo) error
}
