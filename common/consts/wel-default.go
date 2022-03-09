package consts

type WelDefaultInfo struct {
	LastScannedBlock int64 `gorm:"column:last_scanned_block"`
}

type IWelInfoRepo interface {
	Get() (*EthDefaultInfo, error)
	Create(info *EthDefaultInfo) error
	Update(info *EthDefaultInfo) error
}
