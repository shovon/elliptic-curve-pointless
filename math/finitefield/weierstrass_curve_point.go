package finitefield

import (
	"elliptic-curve-interfaces/maybe"
	"fmt"
	"math/big"
)

type (
	NotInfinity[T CurveFiniteField]           [2]*big.Int
	WeierstrassCurvePoint[T CurveFiniteField] maybe.Maybe[NotInfinity[T]]
)

func (c WeierstrassCurvePoint[T]) Equal(i WeierstrassCurvePoint[T]) bool {
	p1, ok1 := maybe.Extract(maybe.Maybe[NotInfinity[T]](c))
	p2, ok2 := maybe.Extract(maybe.Maybe[NotInfinity[T]](i))
	if !ok1 && !ok2 {
		return true
	}
	if !ok1 || !ok2 {
		return false
	}
	return p1.equal(p2)
}

func (p1 NotInfinity[T]) equal(p2 NotInfinity[T]) bool {
	return p1[0].Cmp(p2[0]) == 0 && p1[1].Cmp(p2[1]) == 0
}

func (c NotInfinity[T]) modulo() NotInfinity[T] {
	var t T
	x := new(big.Int)
	y := new(big.Int)

	x.SetBytes(c[0].Bytes()).Mod(x, t.P())
	y.SetBytes(c[1].Bytes()).Mod(y, t.P())

	return NotInfinity[T]{x, y}
}

func (c WeierstrassCurvePoint[T]) Negate() WeierstrassCurvePoint[T] {
	if c == WeierstrassCurvePoint[T](maybe.Nothing[NotInfinity[T]]()) {
		return c
	}
	return WeierstrassCurvePoint[T](maybe.Then(maybe.Maybe[NotInfinity[T]](c), func(value NotInfinity[T]) maybe.Maybe[NotInfinity[T]] {
		y := big.NewInt(-1)
		return maybe.Something(NotInfinity[T]{value[0], y.Mul(y, value[1])}.modulo())
	}))
}

func modInverse(a, m *big.Int) *big.Int {
	return new(big.Int).ModInverse(a, m)
}

func (c WeierstrassCurvePoint[T]) Add(i WeierstrassCurvePoint[T]) WeierstrassCurvePoint[T] {
	var curve T

	p1, ok1 := maybe.Extract(maybe.Maybe[NotInfinity[T]](c))
	p2, ok2 := maybe.Extract(maybe.Maybe[NotInfinity[T]](i))

	if !ok1 {
		return i
	}

	if !ok2 {
		return c
	}

	m := new(big.Int)
	x := new(big.Int)
	y := new(big.Int)

	if p1.equal(p2) {
		// fmt.Printf("Equal! %v %v\n", c, i)
		if p1[1].Cmp(big.NewInt(0)) == 0 {
			return WeierstrassCurvePoint[T]{}
		}
		// Calculate the slope (m) of the tangent line
		numerator := new(big.Int).Mul(big.NewInt(3), new(big.Int).Mul(p1[0], p1[0]))
		numerator.Add(numerator, curve.A())
		denominator := new(big.Int).Mul(big.NewInt(2), p1[1])
		m.Mul(numerator, modInverse(denominator, curve.P()))
		m.Mod(m, curve.P())
	} else {
		// fmt.Printf("Not equal! %v %v\n", c, i)
		// Calculate the slope (m) of the secant line
		numerator := new(big.Int).Sub(p2[1], p1[1])
		denominator := new(big.Int).Sub(p2[0], p1[0])
		if denominator.Cmp(big.NewInt(0)) == 0 {
			return WeierstrassCurvePoint[T]{}
		}
		m.Mul(numerator, modInverse(denominator, curve.P()))
		m.Mod(m, curve.P())
	}

	// Calculate x3 = m^2 - 2x1 (mod p)
	x.Mul(m, m)
	x.Sub(x, p1[0])
	x.Sub(x, p2[0])
	x.Mod(x, curve.P())

	// Calculate y3 = m(x1 - x3) - y1 (mod p)
	y.Sub(p1[0], x)
	y.Mul(y, m)
	y.Sub(y, p1[1])
	y.Mod(y, curve.P())

	return WeierstrassCurvePoint[T](maybe.Something(NotInfinity[T]{x, y}))
}

func (c WeierstrassCurvePoint[T]) ScalarMultiply(n *big.Int) WeierstrassCurvePoint[T] {
	result := WeierstrassCurvePoint[T]{}
	temp := c
	for i := 0; i < n.BitLen(); i++ {
		if n.Bit(i) == 1 {
			result = result.Add(temp)
		}
		temp = temp.Add(temp)
	}
	return result
}

// Implement the fmt.Formatter interface for MyType
func (m WeierstrassCurvePoint[T]) Format(f fmt.State, c rune) {
	switch c {
	case 'v':
		n, ok := maybe.Extract(maybe.Maybe[NotInfinity[T]](m))
		if !ok {
			fmt.Fprintf(f, "Point at infinity")
			return
		} else {
			// Default behavior, similar to %v or %s
			fmt.Fprintf(f, "[%s, %s]", n[0], n[1])
		}
	default:
		// Fallback to the default formatting for any other verbs
		fmt.Fprintf(f, "%%!%c(%T=%+v)", c, m, m)
	}
}
