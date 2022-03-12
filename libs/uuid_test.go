package libs

import (
	"fmt"
	"testing"
)

func TestUniq(t *testing.T) {
	fmt.Println("Uniq() = ", Uniq())
	fmt.Println("UniqN(32) = ", UniqN(32))
}
