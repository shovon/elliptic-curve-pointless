package finitefield

import "math/big"

type FiniteField interface {
	P() *big.Int
}
