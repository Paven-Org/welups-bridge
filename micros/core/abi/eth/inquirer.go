package eth

import (
	"context"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthInquirer struct {
	client *ethclient.Client
}

func MkEthInquirer(c *ethclient.Client) *EthInquirer {
	return &EthInquirer{client: c}
}

func (inq *EthInquirer) BalanceAt(account string) (*big.Int, error) {
	ctx := context.Background()
	address := common.HexToAddress(account)
	block, err := inq.client.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}
	return inq.client.BalanceAt(ctx, address, (&big.Int{}).SetUint64(block))
}

func (inq *EthInquirer) BalanceOf(contract string, account string) (*big.Int, error) {
	contractAddr := common.HexToAddress(contract)
	method := "balanceOf"
	accountAddr := common.HexToAddress(account)

	abiJson := `[{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`
	abi, err := abi.JSON(strings.NewReader(abiJson))
	if err != nil {
		return nil, err
	}

	payload, err := abi.Pack(method, accountAddr)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{From: accountAddr, To: &contractAddr, Data: payload}

	ctx := context.Background()
	block, err := inq.client.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}

	res, err := inq.client.CallContract(ctx, msg, (&big.Int{}).SetUint64(block))
	balance := &big.Int{}
	balance.SetBytes(res)
	if err != nil {
		return balance, err
	}

	return balance, nil
}
