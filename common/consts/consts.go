package consts

import (
	"bridge/common"
	"math/big"
)

var (
	EthereumMainnet = big.NewInt(1)
	EthereumGoerli  = big.NewInt(5)
)

var EthChainFromEnv = map[string](*big.Int){
	common.LocalEnv:      EthereumGoerli,
	common.DevEnv:        EthereumGoerli,
	common.StagingEnv:    EthereumGoerli,
	common.ProductionEnv: EthereumMainnet,
}

const (
	EthereumTk = "0x0000000000000000000000000000000000000000"
	WelupsTk   = "W9yD14Nj9j7xAB4dbGeiX9h8unkKHxuTtb"
)
