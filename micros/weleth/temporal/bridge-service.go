package welethService

import (
	"bridge/service-managers/logger"
	"context"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const (
	welethQueue      = "TEMPORAL_QUEUE_WELETH"
	pingpongActivity = "PING_ACTIVITY"
	pingpongWorkflow = "PING_WORKFLOW"
)

func RegisterWelethBridgeService(w worker.Worker) {
	// register workflow an activities
	w.RegisterWorkflowWithOptions(PingPongWorkflow, workflow.RegisterOptions{Name: pingpongWorkflow})
	w.RegisterActivityWithOptions(PingPongActivity, activity.RegisterOptions{Name: pingpongActivity})

}

func PingPongActivity(ctx context.Context, ping string) (string, error) {
	logger.Get().Info().Msgf("[activity] KV in context: pkeys=%s", ctx.Value("pkeys"))
	logger.Get().Info().Msg("Received ping: " + ping)
	return "pong", nil
}

func PingPongWorkflow(ctx workflow.Context, ping string) (string, error) {
	workflow.GetLogger(ctx).Info("[workflow] KV in context: " + "pkeys=" + (ctx.Value("pkeys").(string)))
	workflow.GetLogger(ctx).Info("Send ping: " + ping)
	ao := workflow.ActivityOptions{
		TaskQueue:              welethQueue,
		ScheduleToCloseTimeout: time.Second * 60,
		ScheduleToStartTimeout: time.Second * 60,
		StartToCloseTimeout:    time.Second * 60,
		HeartbeatTimeout:       time.Second * 10,
		WaitForCancellation:    false,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	future := workflow.ExecuteActivity(ctx, pingpongActivity, ping)
	var res string
	err := future.Get(ctx, &res)
	if err != nil {
		workflow.GetLogger(ctx).Error("Failed to exec activity", pingpongActivity, "error:", err)
		return "", err
	}
	workflow.GetLogger(ctx).Debug("Result: " + res)
	return res, nil
}
