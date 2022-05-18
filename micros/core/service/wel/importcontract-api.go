package welService

const (
	ImportContractQueue = "WelImportContractService"

	Issue = "Issue"

	WatchForTx2TreasuryWF = "WatchForTx2TreasuryWF"
	// signal
	BatchIssueSignal = "BatchedIssueSignal"

	BatchIssueID = "BatchIssueWFOnlyInstance"
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
//func MkImportContractService(client *welclient.GrpcClient, tempCli client.Client, daos *dao.DAOs, contractAddr string) (*ImportContractService, error) {
//	imp := welImport.MkWelImport(client, contractAddr)
//
//	return &ImportContractService{cli: client, tempCli: tempCli, imp: imp, dao: daos.Wel, defaultFeelimit: 8000000}, nil
//}
