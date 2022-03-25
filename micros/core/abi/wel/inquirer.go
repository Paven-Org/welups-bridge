package wel

import (
	"fmt"
	"math/big"

	welclient "github.com/Clownsss/gotron-sdk/pkg/client"
	"github.com/Clownsss/gotron-sdk/pkg/proto/api"
)

// Contract caller for the ERC20 interface

type WelInquirer struct {
	cli *welclient.GrpcClient
}

func MkWelInquirer(welcli *welclient.GrpcClient) *WelInquirer {
	return &WelInquirer{cli: welcli}
}

func (inq *WelInquirer) TriggerConstantContract(from, contractAddress, method, jsonString string) (*api.TransactionExtention, error) {
	return inq.cli.TriggerConstantContract(from, contractAddress, method, jsonString)
}

func (inq *WelInquirer) TriggerContract(from, contractAddress, method, jsonString string, feeLimit, tAmount int64, tTokenID string, tTokenAmount int64) (*api.TransactionExtention, error) {
	return inq.cli.TriggerContract(from, contractAddress, method, jsonString, feeLimit, tAmount, tTokenID, tTokenAmount)
}

func (inq *WelInquirer) WRC20balanceOf(opts *CallOpts, contractAddr string, account string) (*big.Int, error) {
	jsonString := fmt.Sprintf(`[{"address":"%s"}]`, account)
	tx, err := inq.TriggerConstantContract(opts.From, contractAddr, "balanceOf(address)", jsonString)
	if err != nil {
		fmt.Println("BalanceOf failed, error: ", err.Error())
		return big.NewInt(0), err
	}

	fmt.Printf("%x\n", tx.Txid)
	if len(tx.GetConstantResult()) < 1 {
		return big.NewInt(0), fmt.Errorf("No result for read transaction BalanceOf")
	}

	balance := &big.Int{}
	balance.SetBytes(tx.GetConstantResult()[0])

	return balance, nil
}

func (inq *WelInquirer) GetAssets(account string) (map[string]int64, error) {
	tx, err := inq.cli.GetAccountDetailed(account)
	if err != nil {
		fmt.Println("getNativeBalace failed, error: ", err.Error())
		return nil, err
	}

	return tx.Assets, nil
}

func (inq *WelInquirer) GetNativeBalance(account string) (int64, error) {
	tx, err := inq.cli.GetAccountDetailed(account)
	if err != nil {
		fmt.Println("getNativeBalace failed, error: ", err.Error())
		return 0, err
	}

	return tx.Balance, nil
}
