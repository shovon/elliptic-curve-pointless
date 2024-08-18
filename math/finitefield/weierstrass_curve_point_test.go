package finitefield

import (
	"elliptic-curve-interfaces/maybe"
	"math/big"
	"testing"
)

type Dummy struct{}

func (Dummy) A() *big.Int {
	b := new(big.Int)
	b.SetString("0", 16)
	return b
}

func (Dummy) B() *big.Int {
	b := new(big.Int)
	b.SetString("7", 16)
	return b
}

func (Dummy) P() *big.Int {
	b := new(big.Int)
	b.SetString("17", 16)
	return b
}

func Generator() WeierstrassCurvePoint[Dummy] {
	x := new(big.Int)
	y := new(big.Int)

	x.SetString("15", 16)
	y.SetString("13", 16)

	return WeierstrassCurvePoint[Dummy](maybe.Something(NotInfinity[Dummy]{x, y}))
}

func TestDefaultWeierstrassCurve(t *testing.T) {
	var c WeierstrassCurvePoint[Dummy]
	_, ok := maybe.Extract(maybe.Maybe[NotInfinity[Dummy]](c))
	if ok {
		t.Error("A default WeierstrassCurvePoint should be the point at infinity")
	}
}

func TestScalarMultiplyZero(t *testing.T) {
	result := Generator().ScalarMultiply(big.NewInt(0))
	_, ok := maybe.Extract(maybe.Maybe[NotInfinity[Dummy]](result))
	if ok {
		t.Error("A WeierstrassCurvePoint scalar-multiplied by zero should be the point at infinity")
	}
}

func TestScalarMultiplyOne(t *testing.T) {
	result := Generator().ScalarMultiply(big.NewInt(1))
	if !result.Equal(Generator()) {
		t.Error("Scalar multiplying by zero should yield the original curve point")
	}
}

func TestScalarMultiplyTwo(t *testing.T) {
	result := Generator().ScalarMultiply(big.NewInt(2))
	expected := WeierstrassCurvePoint[Dummy](maybe.Something(NotInfinity[Dummy]{big.NewInt(2), big.NewInt(10)}))
	if !result.Equal(expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
