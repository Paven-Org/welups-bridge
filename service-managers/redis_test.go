package manager

import (
	"bridge/service-managers/logger"
	"context"
	"testing"
	"time"
)

var rm *RedisManager

func TestRedis(t *testing.T) {

	cli0, err := rm.GetRedisClient("AuthDB")

	if err != nil {
		t.Fatalf("Unable to open connection to redis, error: %s ", err.Error())
	}

	ctx := context.Background()

	if err := cli0.Set(ctx, "test", "somethingsomething", time.Minute*2).Err(); err != nil {
		t.Fatalf("Unable to set redis key, error: %s", err.Error())
	}

	if res := cli0.Get(ctx, "test"); res.Err() != nil {
		t.Fatalf("Unable to set redis key, error: %s", res.Err())
	} else {
		logger.Get().Info().Msgf("key \"test\": %s", res.String())
	}

	if err := cli0.Del(ctx, "test").Err(); err != nil {
		t.Fatalf("Unable to delete redis key, error: %s", err.Error())
	}
}
