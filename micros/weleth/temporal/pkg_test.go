package welethService

import (
	"bridge/common"
	manager "bridge/service-managers"
	"bridge/service-managers/logger"
	"context"
	"fmt"
	"sync"
	"testing"

	"gitlab.com/rwxrob/uniq"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var c client.Client

func TestMain(m *testing.M) {
	var err error
	key := string(uniq.Bytes(32))
	fmt.Printf("[Test] Random AES key used in this test: %o", []byte(key))
	c, err = manager.MkTemporalClient(common.TemporalCliconf{
		Host:      "localhost",
		Port:      7233,
		Namespace: "devWelbridge",
		Secret:    key,
	}, []string{"pkeys", "ekey"})
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to connect to temporal backend")
	}
	defer c.Close()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if err := manager.SpawnTemporalWorker(ctx, c, pingpongQueue, worker.Options{}, RegisterPingPongService); err != nil {
			logger.Get().Err(err).Msg("Unable to spawn worker")
		}
		wg.Done()
		return
	}()
	m.Run()
	cancel()
	wg.Wait()
}

func TestPing(t *testing.T) {
	wo := client.StartWorkflowOptions{
		TaskQueue: pingpongQueue,
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, "pkeys", "123123")
	fmt.Printf("[Test] value for pkeys in context: %v\n", ctx.Value("pkeys"))
	we, err := c.ExecuteWorkflow(ctx, wo, pingpongWorkflow, "ping")
	if err != nil {
		t.Fatal("Unable to execute workflow, error: " + err.Error())
	}
	t.Log("Workflow", we.GetID(), "runID=", we.GetRunID(), "dispatched")

	var res string
	if err := we.Get(ctx, &res); err != nil {
		t.Fatal("Unable to get workflow's result, error: ", err.Error())
	}
	logger.Get().Info().Msg("result: " + res)
}
