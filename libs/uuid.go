package libs

import (
	"gitlab.com/rwxrob/uniq"
)

func Uniq() string {
	return uniq.Hex(16)
}
