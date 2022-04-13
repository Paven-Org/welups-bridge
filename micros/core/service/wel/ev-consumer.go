package welService

import (
	"bridge/libs"
	"bridge/micros/core/dao"
	welDAO "bridge/micros/core/dao/wel-account"
	welListener "bridge/service-managers/listener/wel"
	"context"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
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
	weldao       welDAO.IWelDAO
	tempCli      client.Client
	abi          abi.ABI
}

func NewGovEvConsumer(addr string, daos *dao.DAOs, tempCli client.Client) *GovEvConsumer {
	exportAbiJSON, err := os.Open("abi/wel/Governance.json")
	if err != nil {
		panic(err)
	}

	defer exportAbiJSON.Close()

	abi, err := abi.JSON(exportAbiJSON)
	if err != nil {
		panic(err)
	}

	return &GovEvConsumer{
		ContractAddr: addr,
		weldao:       daos.Wel,
		tempCli:      tempCli,
		abi:          abi,
	}
}

func (ge *GovEvConsumer) GetConsumer() ([]*welListener.EventConsumer, error) {
	return []*welListener.EventConsumer{
		{
			Address: ge.ContractAddr,
			Topic: crypto.Keccak256Hash(
				[]byte(ge.abi.Events["RoleGranted"].Sig),
			),
			ParseEvent: ge.RoleGranted,
		},
		{
			Address: ge.ContractAddr,
			Topic: crypto.Keccak256Hash(
				[]byte(ge.abi.Events["RoleRevoked"].Sig),
			),

			ParseEvent: ge.RoleRevoked,
		},
	}, nil
}

func (ge *GovEvConsumer) RoleGranted(t *welListener.Transaction, logpos int) error {
	fmt.Println("Consume rolegranted")
	roleRaw := "0x" + common.Bytes2Hex(t.Log[logpos].Topics[1])
	address, _ := libs.HexToB58("0x41" + common.Bytes2Hex(t.Log[logpos].Topics[2][12:]))
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

func (ge *GovEvConsumer) RoleRevoked(t *welListener.Transaction, logpos int) error {
	roleRaw := "0x" + common.Bytes2Hex(t.Log[logpos].Topics[1])
	address, _ := libs.HexToB58("0x41" + common.Bytes2Hex(t.Log[logpos].Topics[2][12:]))
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
