package consts

type WelDefaultInfo struct {
	LastScannedBlock int64 `gorm:"column:last_scanned_block"`
}

type IWelInfoRepo interface {
	Get() (*WelDefaultInfo, error)
	Create(info *WelDefaultInfo) error
	Update(info *WelDefaultInfo) error
}
