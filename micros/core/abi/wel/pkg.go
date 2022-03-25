package wel

import "crypto/ecdsa"

type CallOpts struct {
	From          string
	Prikey        *ecdsa.PrivateKey
	Fee_limit     int64
	T_amount      int64
	T_tokenID     string
	T_tokenAmount int64
}
