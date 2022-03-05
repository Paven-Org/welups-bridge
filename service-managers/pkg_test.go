package manager

import (
	"bridge/common"
	"testing"
)

func TestMain(m *testing.M) {
	cnf := common.Redisconf{
		Network: "tcp",
		Host:    "localhost",
		Port:    6379,
	}

	rm = MkRedisManager(cnf, map[string]int{
		"AuthDB": 15,
	})

	m.Run()

}
