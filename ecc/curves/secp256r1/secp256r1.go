package secp256r1

import (
	"elliptic-curve-interfaces/math/finitefield"
	"elliptic-curve-interfaces/maybe"
	"math/big"
)

type Curve struct{}

func (Curve) A() *big.Int {
	b := new(big.Int)
	b.SetString("ffffffff00000001000000000000000000000000fffffffffffffffffffffffc", 16)
	return b
}

func (Curve) B() *big.Int {
	b := new(big.Int)
	b.SetString("5ac635d8aa3a93e7b3ebbd55769886bc651d06b0cc53b0f63bce3c3e27d2604b", 16)
	return b
}

func (Curve) P() *big.Int {
	b := new(big.Int)
	b.SetString("ffffffff00000001000000000000000000000000ffffffffffffffffffffffff", 16)
	return b
}

var _ finitefield.CurveFiniteField = Curve{}

func Generator() finitefield.WeierstrassCurvePoint[Curve] {
	x := new(big.Int)
	y := new(big.Int)

	x.SetString("6b17d1f2e12c4247f8bce6e563a440f277037d812deb33a0f4a13945d898c296", 16)
	y.SetString("4fe342e2fe1a7f9b8ee7eb4a7c0f9e162bce33576b315ececbb6406837bf51f5", 16)

	return finitefield.WeierstrassCurvePoint[Curve](maybe.Something(finitefield.NotInfinity[Curve]{x, y}))
}

func GeneratorOrder() *big.Int {
	n := new(big.Int)
	n.SetString("ffffffff00000000ffffffffffffffffbce6faada7179e84f3b9cac2fc632551", 16)
	return n
}
