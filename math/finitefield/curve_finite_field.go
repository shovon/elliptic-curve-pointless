package finitefield

import (
	"elliptic-curve-interfaces/math/primitives"
	"math/big"
)

type CurveFiniteField interface {
	P() *big.Int
	A() *big.Int
	B() *big.Int
}

type DummyCurve struct {
}

func (d DummyCurve) P() *big.Int {
	return big.NewInt(0)
}

func (d DummyCurve) A() *big.Int {
	return big.NewInt(0)
}

func (d DummyCurve) B() *big.Int {
	return big.NewInt(0)
}

var _ primitives.ABValues[*big.Int] = DummyCurve{}
var _ FiniteField = DummyCurve{}
var _ CurveFiniteField = DummyCurve{}
