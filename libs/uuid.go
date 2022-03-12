package libs

import (
	"gitlab.com/rwxrob/uniq"
)

func Uniq() string {
	return uniq.Hex(16)
}

func UniqN(size int) string {
	return uniq.Hex(size)
}
