package ethService

import (
	"bridge/micros/core/dao"
	ethDAO "bridge/micros/core/dao/eth-account"
	ethListener "bridge/service-managers/listener/eth"
	"context"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"go.temporal.io/sdk/client"
)

var roleMap map[string]string = map[string]string{
	fmt.Sprintf("0x%x", crypto.Keccak256([]byte("MANAGER_ROLE"))):        "MANAGER_ROLE",
	fmt.Sprintf("0x%x", crypto.Keccak256([]byte("AUTHENTICATOR"))):       "AUTHENTICATOR",
	"0x0000000000000000000000000000000000000000000000000000000000000000": "super_admin",
}

type GovEvConsumer struct {
	ContractAddr string
	ethdao       ethDAO.IEthDAO
	abi          abi.ABI
	tempCli      client.Client
}

func NewGovEvConsumer(addr string, daos *dao.DAOs, tempCli client.Client) *GovEvConsumer {
	govAbiJSON, err := os.Open("abi/eth/Governance.json")
	if err != nil {
		panic(err)
	}

	defer govAbiJSON.Close()

	abi, err := abi.JSON(govAbiJSON)
	if err != nil {
		panic(err)
	}

	return &GovEvConsumer{
		ContractAddr: addr,
		ethdao:       daos.Eth,
		abi:          abi,
		tempCli:      tempCli,
	}
}

func (ge *GovEvConsumer) GetConsumer() ([]*ethListener.EventConsumer, error) {
	return []*ethListener.EventConsumer{
		{
			Address: common.HexToAddress(ge.ContractAddr),
			Topic: crypto.Keccak256Hash(
				[]byte(ge.abi.Events["RoleGranted"].Sig),
			),
			ParseEvent: ge.RoleGranted,
		},
		{
			Address: common.HexToAddress(ge.ContractAddr),
			Topic: crypto.Keccak256Hash(
				[]byte(ge.abi.Events["RoleRevoked"].Sig),
			),

			ParseEvent: ge.RoleRevoked,
		},
	}, nil
}

func (ge *GovEvConsumer) GetFilterQuery() ethereum.FilterQuery {
	return ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(ge.ContractAddr)},
		Topics: [][]common.Hash{{
			crypto.Keccak256Hash(
				[]byte(ge.abi.Events["RoleGranted"].Sig),
			),
			crypto.Keccak256Hash(
				[]byte(ge.abi.Events["RoleRevoked"].Sig),
			),
		}},
	}

}

func (ge *GovEvConsumer) RoleGranted(l types.Log) error {
	roleRaw := l.Topics[1].Hex()
	address := "0x" + common.Bytes2Hex(l.Topics[2].Bytes()[12:])
	fmt.Println("[grant] raw role and address: ", roleRaw, address)

	role, ok := roleMap[roleRaw]
	if !ok {
		return fmt.Errorf("Unknown role")
	}
	fmt.Println("[grant] role and address: ", role, address)

	ctx := context.Background()
	wo := client.StartWorkflowOptions{
		TaskQueue: GovContractQueue,
	}
	ge.tempCli.ExecuteWorkflow(ctx, wo, SaveRoleWF, address, role)
	return nil
}

func (ge *GovEvConsumer) RoleRevoked(l types.Log) error {
	roleRaw := l.Topics[1].Hex()
	address := "0x" + common.Bytes2Hex(l.Topics[2].Bytes()[12:])
	fmt.Println("[revoke] raw role and address: ", roleRaw, address)

	role, ok := roleMap[roleRaw]
	if !ok {
		return fmt.Errorf("Unknown role")
	}

	fmt.Println("[revoke] role and address: ", role, address)

	ctx := context.Background()
	wo := client.StartWorkflowOptions{
		TaskQueue: GovContractQueue,
	}
	ge.tempCli.ExecuteWorkflow(ctx, wo, RemoveRoleWF, address, role)
	return nil
}
