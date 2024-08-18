package dummy

import (
	"elliptic-curve-interfaces/math/finitefield"
	"elliptic-curve-interfaces/maybe"
	"math/big"
)

type Curve struct{}

func (Curve) A() *big.Int {
	b := new(big.Int)
	b.SetString("0", 16)
	return b
}

func (Curve) B() *big.Int {
	b := new(big.Int)
	b.SetString("7", 16)
	return b
}

func (Curve) P() *big.Int {
	b := new(big.Int)
	b.SetString("17", 16)
	return b
}

var _ finitefield.CurveFiniteField = Curve{}

func Generator() finitefield.WeierstrassCurvePoint[Curve] {
	x := new(big.Int)
	y := new(big.Int)

	x.SetString("15", 16)
	y.SetString("13", 16)

	return finitefield.WeierstrassCurvePoint[Curve](maybe.Something(finitefield.NotInfinity[Curve]{x, y}))
}

func GeneratorOrder() *big.Int {
	n := new(big.Int)
	n.SetString("18", 16)
	return n
}
