package manager

import (
	"bridge/common"
	"bridge/service-managers/logger"
	"context"
	"fmt"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func MkTemporalClient(cnf common.TemporalCliconf) (client.Client, error) {
	return client.NewClient(client.Options{
		HostPort:  cnf.Host + ":" + fmt.Sprintf("%d", cnf.Port),
		Namespace: cnf.Namespace,
		Logger:    logger.StdLogger(),
	})
}

//Spawn new temporal worker, blocking
func SpawnTemporalWorker(ctx context.Context, client client.Client, taskQueue string, opts worker.Options, register func(w worker.Worker)) error {
	w := worker.New(client, taskQueue, opts)
	logger.Get().Info().Msg("New temporal worker created")
	register(w)

	logger.Get().Info().Msg("Starting temporal worker...")

	var err error
	var ch = make(chan interface{})
	defer close(ch)
	go func() {
		if err = w.Run(ch); err != nil {
			logger.Get().Err(err).Msg("Temporal worker failed")
			return
		}
		return
	}()

	ch <- (<-ctx.Done())

	return err
}
