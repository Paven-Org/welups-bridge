package wel

import (
	"time"

	gotron "github.com/Clownsss/gotron-sdk/pkg/client"
)

//CNodeClient is the wrapper of the nodeClient
type ExtNodeClient struct {
	*gotron.GrpcClient
	grpcTimeout time.Duration
}

// NewExtNodeClient create grpc controller
func NewExtNodeClient(address string) *ExtNodeClient {
	client := &ExtNodeClient{
		GrpcClient:  gotron.NewGrpcClient(address),
		grpcTimeout: 5 * time.Second,
	}
	return client
}

// NewExtNodeClientWithTimeout create grpc controller
func NewExtNodeClientWithTimeout(address string, timeout time.Duration) *ExtNodeClient {
	client := &ExtNodeClient{
		GrpcClient:  gotron.NewGrpcClientWithTimeout(address, timeout),
		grpcTimeout: timeout,
	}
	return client
}

func NewExtNodeClientFromCli(cli *gotron.GrpcClient, timeout time.Duration) *ExtNodeClient {
	client := &ExtNodeClient{
		GrpcClient:  cli,
		grpcTimeout: timeout,
	}
	return client
}
