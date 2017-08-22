// Package util some common utils
package util

import (
	"math/big"
)

const (
	// one ether in wei
	OneEtherInWei = 1000000000000000000
)

// WeiToEther convert wei to ether in float
func WeiToEther(value *big.Int) *big.Float {
	coef := big.NewInt(OneEtherInWei)
	// Make conversion
	v := new(big.Float).SetInt(value)
	c := new(big.Float).SetInt(coef)
	converted := new(big.Float).Quo(v, c)

	// Round
	//converted = Round(converted, 5)
	// Return results
	return converted
}
