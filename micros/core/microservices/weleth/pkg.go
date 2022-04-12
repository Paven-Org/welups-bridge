package msweleth

import (
	welethService "bridge/micros/weleth/temporal"
)

const (
	TaskQueue = welethService.WelethServiceQueue

	GetWelToEthCashinByTxHash  = "GetWelToEthCashinByTxHashWF"
	GetEthToWelCashoutByTxHash = "GetEthToWelCashoutByTxHashWF"

	GetEthToWelCashinByTxHash  = "GetEthToWelCashinByTxHashWF"
	GetWelToEthCashoutByTxHash = "GetWelToEthCashoutByTxHashWF"

	CreateW2ECashinClaimRequestWF  = "CreateW2ECashinClaimRequestWF"
	CreateE2WCashoutClaimRequestWF = "CreateE2WCashoutClaimRequestWF"

	WaitForPendingW2ECashinClaimRequestWF  = "WaitForPendingW2ECashinClaimRequestWF"
	WaitForPendingE2WCashoutClaimRequestWF = "WaitForPendingE2WCashoutClaimRequestWF"
)
