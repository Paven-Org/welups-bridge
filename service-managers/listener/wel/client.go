package wel

import (
	"context"
	"fmt"
	"time"

	"github.com/Clownsss/gotron-sdk/pkg/address"
	gotron "github.com/Clownsss/gotron-sdk/pkg/client"
	"github.com/Clownsss/gotron-sdk/pkg/proto/core"
)

//CNodeClient is the wrapper of the nodeClient
type ExtNodeClient struct {
	*gotron.GrpcClient
	grpcTimeout time.Duration
}

//GetContract ...
func (g *ExtNodeClient) GetContract(contractAddress string) (*core.SmartContract, error) {
	var err error
	contractDesc, err := address.Base58ToAddress(contractAddress)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.grpcTimeout)
	defer cancel()

	sm, err := g.Client.GetContract(ctx, gotron.GetMessageBytes(contractDesc))
	if err != nil {
		return nil, err
	}
	if sm == nil {
		return nil, fmt.Errorf("invalid contract abi")
	}

	return sm, nil
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
