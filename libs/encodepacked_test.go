package libs

import (
	"encoding/hex"
	"fmt"
)

func TestPacked() {
	// bytes32 stateRoots
	stateRoots := "3a53dc4890241dbe03e486e785761577d1c369548f6b09aa38017828dcdf5c27"
	// uint256[2] calldata signatures
	signatures := []string{
		"3402053321874964899321528271743396700217057178612185975187363512030360053932",
		"1235124644010117237054094970590473241953434069965207718920579820322861537001",
	}
	// uint256 feeReceivers,
	feeReceivers := "0"
	// bytes calldata txss
	txss := "000000000000000100010000"

	result := encodePacked(
		encodeBytesString(stateRoots),
		encodeUint256Array(signatures),
		encodeUint256(feeReceivers),
		encodeBytesString(txss),
	)

	got := hex.EncodeToString(result)
	want := "3a53dc4890241dbe03e486e785761577d1c369548f6b09aa38017828dcdf5c2707857e73108d077c5b7ef89540d6493f70d940f1763a9d34c9d98418a39d28ac02bb0e4743a7d0586711ee3dd6311256579ab7abcd53c9c76f040bfde4d6d6e90000000000000000000000000000000000000000000000000000000000000000000000000000000100010000"
	fmt.Println(got == want) // true
}
