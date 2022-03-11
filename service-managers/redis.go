package manager

import (
	"bridge/common"
	"bridge/common/consts"
	"bridge/service-managers/logger"
	"context"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	StdAuthDBName = "AuthDB"
	StdTestDBName = "TestDB"
)
var StdDbMap map[string]int = map[string]int{
	StdAuthDBName: 1,
	StdTestDBName: 15,
}

type RedisManager struct {
	cnf      common.Redisconf
	lock     sync.RWMutex
	rClients map[int]*redis.Client
	dbMap    map[string]int
}

func MkRedisManager(cnf common.Redisconf, dbMap map[string]int) *RedisManager {
	return &RedisManager{
		cnf:      cnf,
		lock:     sync.RWMutex{},
		rClients: make(map[int]*redis.Client),
		dbMap:    dbMap,
	}
}

func (r *RedisManager) GetRedisClient(dbname string) (*redis.Client, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	dbnum, ok := r.dbMap[dbname]
	if !ok {
		err := consts.ErrRedisDBNotExist
		logger.Get().Err(err).Msgf("Redis db %s does not exist", dbname)
		return nil, err
	}

	if cli, ok := r.rClients[dbnum]; ok {
		return cli, nil
	}

	cnf := r.cnf
	cli := redis.NewClient(&redis.Options{
		Network:  cnf.Network,
		Addr:     fmt.Sprintf("%s:%d", cnf.Host, cnf.Port),
		Username: cnf.Username,
		Password: cnf.Password,
		DB:       dbnum,
	})

	if err := cli.Ping(context.Background()).Err(); err != nil {
		logger.Get().Err(err).Msgf("Unable to connect to redis db %d, trying to close connection...", dbnum)
		for {
			if err := cli.Close(); err != nil {
				logger.Get().Err(err).Msgf("Failed to close client to redis db %d, retrying...", dbnum)
			} else {
				logger.Get().Info().Msgf("Closed client to redis db %d", dbnum)
				break
			}
		}
		return nil, err
	}

	r.rClients[dbnum] = cli

	return cli, nil
}

func (r *RedisManager) Flush(dbname string) error {
	cli, err := r.GetRedisClient(dbname)
	if err != nil {
		logger.Get().Err(err).Msgf("Failed to get redis connection")
		return err
	}

	if err = cli.FlushDB(context.Background()).Err(); err != nil {
		logger.Get().Err(err).Msgf("Failed to flush redis db %s, db number: %d", dbname, r.dbMap[dbname])
		return err
	}

	logger.Get().Info().Msgf("Flushed redis db %s, db number: %d", dbname, r.dbMap[dbname])
	return nil
}

func (r *RedisManager) CloseAll() {
	r.lock.Lock()
	defer r.lock.Unlock()

	for db, cli := range r.rClients {
		for {
			if err := cli.Close(); err != nil {
				logger.Get().Err(err).Msgf("Failed to close client to redis db %d, retrying...", db)
			} else {
				logger.Get().Info().Msgf("Closed client to redis db %d", db)
				break
			}
		}
	}
}
