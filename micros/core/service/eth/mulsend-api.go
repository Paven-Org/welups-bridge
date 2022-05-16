package ethService

const (
	MulsendContractQueue = "EthMulsendContractService"

	Withdraw = "Withdraw"
	Disperse = "Disperse"

	// signal
	BatchDisperseSignal = "BatchedDisperseSignal"

	BatchDisperseID = "BatchDisperseWFOnlyInstance"
)

//const (
//	ImportContractQueue = implement.ImportContractQueue
//
//	// going to be deprecated
//	WithdrawWorkflow = implement.WithdrawWorkflow
//	IssueWorkflow    = implement.IssueWorkflow
//
//	// signal
//	BatchIssueSignal = implement.BatchIssueSignal
//
//	BatchIssueID = implement.BatchIssueID
//)
//
//type ImportContractService = implement.ImportContractService
//
//func MkImportContractService(client *ethclient.GrpcClient, tempCli client.Client, daos *dao.DAOs, contractAddr string) (*ImportContractService, error) {
//	imp := ethImport.MkEthImport(client, contractAddr)
//
//	return &ImportContractService{cli: client, tempCli: tempCli, imp: imp, dao: daos.Eth, defaultFeelimit: 8000000}, nil
//}
