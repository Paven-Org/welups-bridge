package wel

import (
	"bridge/service-managers/logger"
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	"math/big"

	welclient "github.com/Clownsss/gotron-sdk/pkg/client"
	"github.com/Clownsss/gotron-sdk/pkg/proto/api"
	"github.com/ethereum/go-ethereum/crypto"
	"google.golang.org/protobuf/proto"
)

// WelGov is an auto generated Go binding around an Welereum contract.
type WelGov struct {
	address string
	cli     *welclient.GrpcClient
}

type CallOpts struct {
	From          string
	Prikey        *ecdsa.PrivateKey
	Fee_limit     int64
	T_amount      int64
	T_tokenID     string
	T_tokenAmount int64
}

// NewWelGov creates a new instance of WelGov, bound to a specific deployed contract.
func MkWelGov(address string, welcli *welclient.GrpcClient) *WelGov {
	return &WelGov{address: address, cli: welcli}
}

// TtriggerConstantContract and return tx result
func (welgov *WelGov) TriggerConstantContract(from, contractAddress, method, jsonString string) (*api.TransactionExtention, error) {
	return welgov.cli.TriggerConstantContract(from, contractAddress, method, jsonString)
}

func (welgov *WelGov) TriggerContract(from, contractAddress, method, jsonString string, feeLimit, tAmount int64, tTokenID string, tTokenAmount int64) (*api.TransactionExtention, error) {
	return welgov.cli.TriggerContract(from, contractAddress, method, jsonString, feeLimit, tAmount, tTokenID, tTokenAmount)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (welgov *WelGov) HasRole(opts *CallOpts, role [32]byte, account string) (bool, error) {
	//params := [](map[string]interface{}){
	//	{"bytes32": fmt.Sprintf("%x", role)},
	//	{"address": account},
	//}
	//jsonBytes, _ := json.Marshal(params)
	//jsonString := fmt.Sprintf("%s", jsonBytes)
	jsonString := fmt.Sprintf(`[{"bytes32": "%x"},{"address":"%s"}]`, role, account)
	tx, err := welgov.TriggerConstantContract(opts.From, welgov.address, "hasRole(bytes32,address)", jsonString)
	if err != nil {
		fmt.Println("HasRole failed, error: ", err.Error())
		return false, err
	}
	fmt.Printf("%x\n", tx.Txid)
	if len(tx.GetConstantResult()) < 1 {
		return false, fmt.Errorf("No result for read transaction HasRole")
	}

	out := &big.Int{}
	out.SetBytes(tx.GetConstantResult()[0])

	return out.Int64() == 1, nil
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (welgov *WelGov) GrantRole(opts *CallOpts, role [32]byte, account string) (*api.TransactionExtention, error) {
	jsonString := fmt.Sprintf(`[{"bytes32": "%x"},{"address":"%s"}]`, role, account)
	res, err := welgov.cli.TriggerContract(opts.From, welgov.address, "grantRole(bytes32,address)", jsonString, opts.Fee_limit, opts.T_amount, opts.T_tokenID, opts.T_tokenAmount)
	if err != nil {
		logger.Get().Err(err).Msgf("RPC to make contract call failed")
		return nil, err
	}

	// signing
	logger.Get().Info().Msgf("Signing transaction...")

	tx := res.Transaction
	rawData, err := proto.Marshal(tx.GetRawData())
	if err != nil {
		logger.Get().Err(err).Msgf("Failed to sign transaction")
		return nil, err
	}
	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)

	signature, err := crypto.Sign(hash, opts.Prikey)
	if err != nil {
		return nil, err
	}
	tx.Signature = append(tx.Signature, signature)

	// Broadcast
	ret, err := welgov.cli.Broadcast(tx)
	if err != nil {
		logger.Get().Err(err).Msgf("Failed to broadcast signed transaction")
		return nil, err
	}

	fmt.Println(ret)
	if ret.GetCode() != api.Return_SUCCESS {
		err = fmt.Errorf(api.ReturnResponseCode_name[int32(ret.GetCode())])
		logger.Get().Err(err).Msgf("Failed to broadcast signed transaction")
		return nil, err
	}

	// Optionally loop to confirm transaction right here
	// it's better to do this inside of a workflow anyway, so kinda not compelling

	// return
	res.Transaction = tx
	res.Result = ret

	return res, nil
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (welgov *WelGov) RevokeRole(opts *CallOpts, role [32]byte, account string) (*api.TransactionExtention, error) {
	jsonString := fmt.Sprintf(`[{"bytes32": "%x"},{"address":"%s"}]`, role, account)
	res, err := welgov.cli.TriggerContract(opts.From, welgov.address, "revokeRole(bytes32,address)", jsonString, opts.Fee_limit, opts.T_amount, opts.T_tokenID, opts.T_tokenAmount)
	if err != nil {
		logger.Get().Err(err).Msgf("RPC to make contract call failed")
		return nil, err
	}

	// signing
	logger.Get().Info().Msgf("Signing transaction...")

	tx := res.Transaction
	rawData, err := proto.Marshal(tx.GetRawData())
	if err != nil {
		logger.Get().Err(err).Msgf("Failed to sign transaction")
		return nil, err
	}
	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)

	signature, err := crypto.Sign(hash, opts.Prikey)
	if err != nil {
		return nil, err
	}
	tx.Signature = append(tx.Signature, signature)

	// Broadcast
	ret, err := welgov.cli.Broadcast(tx)
	if err != nil {
		logger.Get().Err(err).Msgf("Failed to broadcast signed transaction")
		return nil, err
	}

	fmt.Println(ret)
	if ret.GetCode() != api.Return_SUCCESS {
		err = fmt.Errorf(api.ReturnResponseCode_name[int32(ret.GetCode())])
		logger.Get().Err(err).Msgf("Failed to broadcast signed transaction")
		return nil, err
	}

	// Optionally loop to confirm transaction right here
	// it's better to do this inside of a workflow anyway, so kinda not compelling

	// return
	res.Transaction = tx
	res.Result = ret

	return res, nil
}
