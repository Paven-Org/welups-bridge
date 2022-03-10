package wel

import (
	"fmt"
	"log"

	GotronCommon "github.com/Clownsss/gotron-sdk/pkg/common"
	"github.com/Clownsss/gotron-sdk/pkg/proto/api"
	CoreProto "github.com/Clownsss/gotron-sdk/pkg/proto/core"
	"github.com/fatih/structs"
	proto "github.com/golang/protobuf/proto"
)

func NewTransHandler(client *ExtNodeClient, threshold int64) *TransHandler {
	return &TransHandler{
		Client:                  client,
		ConfirmedBlockThreshold: threshold,
	}
}

type TransHandler struct {
	Client                  *ExtNodeClient
	ConfirmedBlockThreshold int64 // = 20
}

//Transaction defines the transaction data
type Transaction struct {
	Hash              string                           `json:"hash,omitempty"`
	Timestamp         int64                            `json:"timestamp,omitempty"`
	ContractAddress   string                           `json:"contract_address,omitempty"`
	Log               []*CoreProto.TransactionInfo_Log `json:"log,omitempty"`
	Type              string                           `json:"type,omitempty"`
	FeeLimit          int64                            `json:"fee_limit"`
	EnergyUsage       int64                            `json:"energy_usage,omitempty"`
	OriginEnergyUsage int64                            `json:"origin_energy_usage"`
	EnergyUsageTotal  int64                            `json:"energy_usage_total"`
	NetUsage          int64                            `json:"net_usage,omitempty"`
	RefBlockBytes     string                           `json:"refBlockBytes,omitempty"`
	RefBlockNum       int64                            `json:"refBlockNum,omitempty"`
	RefBlockHash      string                           `json:"refBlockHash,omitempty"`
	Expiration        int64                            `json:"expiration,omitempty"`
	Status            string                           `json:"status,omitempty"`
	Auths             interface{}                      `json:"auths,omitempty"`
	Data              string                           `json:"data,omitempty"`
	Contract          struct {
		Type      string `json:"type,omitempty"`
		Name      string `json:"name,omitempty"`
		Parameter struct {
			TypeURL string                 `json:"type_url,omitempty"`
			Value   []byte                 `json:"value,omitempty"`
			Raw     map[string]interface{} `json:"raw,omitempty"`
		} `json:"parameter,omitempty"`
	} `json:"contract,omitempty"`
	Scripts     string `json:"scripts"`
	BlockNumber int64  `json:"blockNumber"`
	NumOfBlocks int64  `json:"num_of_blocks,omitempty"`
	Result      string `json:"result,omitempty"`
}

//GetTransactionDetails return transaction info details
func (t *TransHandler) GetTransactionDetails(n string) *Transaction {
	tranInfo := &Transaction{}

	tranDetails, err := t.Client.GetTransactionInfoByID(n)
	nowBlock, err := t.Client.GetNowBlock()
	// TODO: get trans topic from tranDetails

	if err != nil {
		log.Println(err)
		return tranInfo
	}
	tranDetails2, err := t.Client.GetTransactionByID(n)
	if err != nil {
		log.Println(err)
		return tranInfo
	}

	blockId := nowBlock.BlockHeader.GetRawData().GetNumber()

	resultStatus := "unconfirmed"
	diffBlock := blockId - tranDetails.GetBlockNumber()
	if diffBlock >= t.ConfirmedBlockThreshold {
		resultStatus = "confirmed"
	}

	tranInfo = t.ToTransaction(tranDetails2, nowBlock)
	tranInfo.Hash = GotronCommon.Bytes2Hex(tranDetails.Id)
	tranInfo.ContractAddress = GotronCommon.EncodeCheck(tranDetails.GetContractAddress())
	tranInfo.Log = tranDetails.Log
	tranInfo.Result = tranDetails.GetResult().String()
	tranInfo.EnergyUsage = tranDetails.Receipt.EnergyUsage
	tranInfo.OriginEnergyUsage = tranDetails.Receipt.OriginEnergyUsage
	tranInfo.EnergyUsageTotal = tranDetails.Receipt.EnergyUsageTotal
	tranInfo.NetUsage = tranDetails.Receipt.NetUsage
	tranInfo.Status = resultStatus
	tranInfo.NumOfBlocks = diffBlock
	return tranInfo
}

//GetInfoTransactionRange ....
func (t *TransHandler) GetInfoListTransactionRange(blocknum int64, limit int64, sortString string, output chan *Transaction, errChan chan error) {
	if sortString == "inc" {
		blocksList, err := t.Client.GetBlockByLimitNext(blocknum+1, blocknum+limit+1)
		if err != nil {
			errChan <- err
			return
		}
		for _, b := range blocksList.Block {
			tempTrans := b.GetTransactions()
			for _, tx := range tempTrans {
				aTran := t.GetTransactionDetails(GotronCommon.Bytes2Hex(tx.GetTxid()))
				output <- aTran
			}
		}
	} else {
		temp := blocknum - limit + 1
		if temp < 0 {
			temp = 0
		}
		blocksList, err := t.Client.GetBlockByLimitNext(temp, blocknum+1)
		if err != nil {
			errChan <- err
			return
		}
		for _, b := range blocksList.Block {
			tempTrans := b.GetTransactions()
			for _, tx := range tempTrans {
				aTran := t.GetTransactionDetails(GotronCommon.Bytes2Hex(tx.GetTxid()))
				output <- aTran
			}
		}
	}
	return
}

func (tr *TransHandler) ToTransaction(tx *CoreProto.Transaction, b *api.BlockExtention) *Transaction {
	trans := &Transaction{
		Timestamp:     tx.GetRawData().GetTimestamp(),
		Type:          tx.GetRawData().Contract[0].GetType().String(),
		FeeLimit:      tx.GetRawData().GetFeeLimit(),
		RefBlockBytes: GotronCommon.EncodeCheck(tx.GetRawData().GetRefBlockBytes()),
		RefBlockNum:   tx.GetRawData().GetRefBlockNum(),
		RefBlockHash:  GotronCommon.EncodeCheck(tx.GetRawData().GetRefBlockHash()),
		Expiration:    tx.GetRawData().GetExpiration(),
		Auths:         tx.GetRawData().GetAuths(),
		Data:          string(tx.GetRawData().GetData()),
		Scripts:       string(tx.GetRawData().GetScripts()),
		BlockNumber:   b.BlockHeader.GetRawData().GetNumber(),
	}
	trans.Contract.Type = tx.RawData.Contract[0].Type.String()
	trans.Contract.Name = string(tx.RawData.Contract[0].ContractName)
	trans.Contract.Parameter.TypeURL = tx.RawData.Contract[0].Parameter.TypeUrl
	trans.Contract.Parameter.Value = tx.RawData.Contract[0].Parameter.Value
	postFuncs := map[string]func(interface{}) string{}
	base58AddrFunc := func(str interface{}) string {
		return GotronCommon.EncodeCheck(str.([]byte))
	}
	strAddrFunc := func(str interface{}) string {
		return fmt.Sprintf("%s", str.([]byte))
	}
	postFuncs["OwnerAddress"] = base58AddrFunc
	postFuncs["ContractAddress"] = base58AddrFunc
	postFuncs["ToAddress"] = base58AddrFunc
	postFuncs["AssetName"] = strAddrFunc
	switch t := tx.RawData.Contract[0].Type; t {
	case CoreProto.Transaction_Contract_TriggerSmartContract:
		var ctDetails CoreProto.TriggerSmartContract
		proto.Unmarshal(trans.Contract.Parameter.Value, &ctDetails)
		trans.Contract.Parameter.Raw = structs.Map(ctDetails)
	case CoreProto.Transaction_Contract_TransferContract:
		var ctDetails CoreProto.TransferContract
		proto.Unmarshal(trans.Contract.Parameter.Value, &ctDetails)
		trans.Contract.Parameter.Raw = structs.Map(ctDetails)
	case CoreProto.Transaction_Contract_TransferAssetContract:
		var ctDetails CoreProto.TransferAssetContract
		proto.Unmarshal(trans.Contract.Parameter.Value, &ctDetails)
		trans.Contract.Parameter.Raw = structs.Map(ctDetails)
		if asset, err := tr.Client.GetAssetIssueByName(string(ctDetails.AssetName)); err == nil {
			trans.Contract.Parameter.Raw["AssetName"] = asset.Name
			trans.Contract.Parameter.Raw["AssetID"] = string(asset.Id)
		}
	case CoreProto.Transaction_Contract_CreateSmartContract:
		var ctDetails CoreProto.CreateSmartContract
		proto.Unmarshal(trans.Contract.Parameter.Value, &ctDetails)
		trans.Contract.Parameter.Raw = structs.Map(ctDetails)
	case CoreProto.Transaction_Contract_AssetIssueContract:
		var ctDetails CoreProto.AssetIssueContract
		proto.Unmarshal(trans.Contract.Parameter.Value, &ctDetails)
		trans.Contract.Parameter.Raw = structs.Map(ctDetails)
		trans.Contract.Parameter.Raw["Abbr"] = string(ctDetails.Abbr)
		trans.Contract.Parameter.Raw["Description"] = string(ctDetails.Description)
		trans.Contract.Parameter.Raw["Name"] = string(ctDetails.Name)
		trans.Contract.Parameter.Raw["Url"] = string(ctDetails.Url)
		if asset, err := tr.Client.GetAssetIssueByName(string(ctDetails.Name)); err == nil {
			trans.Contract.Parameter.Raw["Id"] = string(asset.Id)
		}
	}
	//TODO: need to process all smartcontract types

	for k, f := range postFuncs {
		if _, ok := trans.Contract.Parameter.Raw[k]; ok {
			trans.Contract.Parameter.Raw[k] = f(trans.Contract.Parameter.Raw[k])
		}
	}
	return trans
}
