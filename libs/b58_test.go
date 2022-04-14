package libs

import (
	"fmt"
	"testing"
)

const (
	b58test = "WDReBjymEzH5Bi4avyfUqXa4rhyvAT7DY2"
	hextest = "0x4125e8370e0e2cf3943ad75e768335c892434bd090"
)

func TestB58toHex(t *testing.T) {
	hex, err := B58toHex(b58test)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println("hex: ", hex)

	stdhex, err := B58toStdHex(b58test)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println("stdhex: ", stdhex)
}

func TestHexToB58(t *testing.T) {
	b58, err := HexToB58(hextest)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println(b58)
}
