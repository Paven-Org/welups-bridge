package wel

import (
	"bridge/libs"
	"bridge/micros/core/model"
	"bridge/service-managers/logger"
	"fmt"
	"math/big"

	welclient "github.com/Paven-Org/gotron-sdk/pkg/client"
	"github.com/Paven-Org/gotron-sdk/pkg/proto/api"
	"github.com/Paven-Org/gotron-sdk/pkg/proto/core"
	"google.golang.org/protobuf/proto"
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

func (inq *WelInquirer) WRC20balanceOf(contractAddr string, account string) (*big.Int, error) {
	jsonString := fmt.Sprintf(`[{"address":"%s"}]`, account)
	tx, err := inq.TriggerConstantContract(account, contractAddr, "balanceOf(address)", jsonString)
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

func (inq *WelInquirer) GetAccount(address string) (*core.Account, error) {
	account, err := inq.cli.GetAccount(address)
	if err == nil {
		fmt.Println("Account " + address + " found")
		return account, nil
	}

	if err.Error() == "account not found" {
		return nil, model.ErrWelAccountNotFound
	}

	return account, err
}

func (inq *WelInquirer) ActivateAccountIfNotExist(address string, activator string, pkey string) error {
	_, err := inq.cli.GetAccount(address)
	if err == nil {
		fmt.Println("Account " + address + " found, no activation needed")
		return nil
	}

	if err.Error() == "account not found" {
		tx, err := inq.cli.Transfer(activator, address, 1)
		if err != nil {
			logger.Get().Err(err).Msgf("RPC failed")
			return err
		}
		rawData, err := proto.Marshal(tx.Transaction.GetRawData())
		if err != nil {
			logger.Get().Err(err).Msgf("Failed to get transaction's raw data")
			return err
		}
		// signing
		signature, err := libs.SignerH256(rawData, pkey)
		if err != nil {
			logger.Get().Err(err).Msgf("Failed to sign transaction")
			return err
		}
		tx.Transaction.Signature = append(tx.Transaction.Signature, signature)
		// broadcast
		ret, err := inq.cli.Broadcast(tx.Transaction)
		if err != nil {
			logger.Get().Err(err).Msgf("Failed to broadcast signed transaction")
			return err
		}
		logger.Get().Info().Msgf("Account activation transaction result message: %s", ret.Message)

		logger.Get().Info().Msgf("Successfully activated wel account %s", address)

		return nil
	}

	return err
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
