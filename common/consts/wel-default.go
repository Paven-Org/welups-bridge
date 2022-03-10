package consts

type WelDefaultInfo struct {
	LastScannedBlock int64 `db:"column:wel_last_scan_block"`
}

type IWelInfoRepo interface {
	Get() (*WelDefaultInfo, error)
	Create(info *WelDefaultInfo) error
	Update(info *WelDefaultInfo) error
}
