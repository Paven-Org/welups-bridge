package consts

import "fmt"

var (
	ErrNilDB           = fmt.Errorf("NIL DB initialized")
	ErrRedisDBNotExist = fmt.Errorf("Redis DB does not exist")
)
