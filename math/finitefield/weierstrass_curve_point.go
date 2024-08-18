package finitefield

import (
	"elliptic-curve-interfaces/math/primitives"
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
	return p1.equal(p2)
}

func (p1 NotInfinity[T]) equal(p2 NotInfinity[T]) bool {
	return p1[0].Cmp(p2[0]) == 0 && p1[1].Cmp(p2[1]) == 0
}

func (p1 NotInfinity[T]) slope(p2 NotInfinity[T]) maybe.Maybe[*big.Int] {
	var curve T

	m := new(big.Int)

	if !p1.equal(p2) {
		p12d := primitives.Point2D[*big.Int](p1)
		p22d := primitives.Point2D[*big.Int](p2)

		// Dividend
		a := new(big.Int)
		a.Sub(p22d.Y(), p12d.Y())

		// Compute the multiplicative inverse of the divisor
		m.Sub(p22d.X(), p12d.X())

		if m.Cmp(big.NewInt(0)) == 0 {
			return maybe.Nothing[*big.Int]()
		}

		// Compute the divisor
		m.ModInverse(m, curve.P())

		// Compute the slope
		m.Mul(a, m)
	} else {
		// p2d := primitives.Point2D[*big.Int](p1)

		// two := big.NewInt(2)
		// three := big.NewInt(3)

		// numerator := new(big.Int).Mul(three, new(big.Int).Exp(p2d[0], two, curve.P()))
		// numerator.Add(numerator, curve.A())
		// numerator.Mod(numerator, curve.P())

		// denominator := new(big.Int).Mul(two, point.Y)
		// denominator.Mod(denominator, curve.P)

		// // Calculate the modular inverse of the denominator
		// invDenominator := InverseMod(denominator, curve.P)
		// if invDenominator == nil {
		// 	return nil, fmt.Errorf("point is not on the curve or division by zero error")
		// }

		two := big.NewInt(2)
		three := big.NewInt(3)

		p2d := primitives.Point2D[*big.Int](p1)
		if p2d.Y().Cmp(big.NewInt(0)) == 0 {
			return maybe.Nothing[*big.Int]()
		}

		// Dividend
		a := big.NewInt(3)
		x := new(big.Int)
		x.SetBytes(p2d.X().Bytes())
		x.Mul(x, x)
		a.Mul(a, three).Mul(a, x).Add(a, curve.A())

		m.Mul(two, p2d.Y()).ModInverse(m, curve.P()).Mul(m, a)
	}

	return maybe.Something(m)
}

func (p1 NotInfinity[T]) yIntercept(p2 NotInfinity[T]) maybe.Maybe[*big.Int] {
	// d = y1 - m*x1

	return maybe.Then(p1.slope(p2), func(m *big.Int) maybe.Maybe[*big.Int] {
		d := new(big.Int)

		p2d := primitives.Point2D[*big.Int](p1)

		m.Mul(m, p2d.X())

		return maybe.Something(d.Sub(p2d.Y(), m))
	})

}

func (p1 NotInfinity[T]) x3(p2 NotInfinity[T]) maybe.Maybe[*big.Int] {
	// x3 = m^2 - x1 - x2

	return maybe.Then(p1.slope(p2), func(result *big.Int) maybe.Maybe[*big.Int] {
		return maybe.Something(result.
			Mul(result, result).
			Sub(result, primitives.Point2D[*big.Int](p1).X()).
			Sub(result, primitives.Point2D[*big.Int](p2).X()))
	})
}

func (p1 NotInfinity[T]) thirdPoint(p2 NotInfinity[T]) WeierstrassCurvePoint[T] {
	return WeierstrassCurvePoint[T](maybe.Then(p1.slope(p2), func(slope *big.Int) maybe.Maybe[NotInfinity[T]] {
		// Extract the y-intercept
		return maybe.Then(p1.yIntercept(p2), func(intercept *big.Int) maybe.Maybe[NotInfinity[T]] {
			// Extract the
			return maybe.Then(p1.x3(p2), func(x3 *big.Int) maybe.Maybe[NotInfinity[T]] {
				y3 := slope.Mul(slope, x3).Add(slope, intercept)
				return maybe.Something(NotInfinity[T]{x3, y3})
			})
		})
	}))
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

func (c WeierstrassCurvePoint[T]) Add(i WeierstrassCurvePoint[T]) WeierstrassCurvePoint[T] {
	p1, ok1 := maybe.Extract(maybe.Maybe[NotInfinity[T]](c))
	p2, ok2 := maybe.Extract(maybe.Maybe[NotInfinity[T]](i))

	if !ok1 {
		return WeierstrassCurvePoint[T](maybe.Then(maybe.Maybe[NotInfinity[T]](i), func(value NotInfinity[T]) maybe.Maybe[NotInfinity[T]] {
			return maybe.Something(value.modulo())
		}))
	}
	if !ok2 {
		return WeierstrassCurvePoint[T](maybe.Then(maybe.Maybe[NotInfinity[T]](c), func(value NotInfinity[T]) maybe.Maybe[NotInfinity[T]] {
			return maybe.Something(value.modulo())
		}))
	}

	return WeierstrassCurvePoint[T](maybe.Maybe[NotInfinity[T]](p1.thirdPoint(p2).Negate()))
}

func (c WeierstrassCurvePoint[T]) ScalarMultiply(n *big.Int) WeierstrassCurvePoint[T] {
	if n.Cmp(big.NewInt(0)) == 0 {
		return WeierstrassCurvePoint[T]{}
	}
	var r0 WeierstrassCurvePoint[T]
	r1 := c

	for i := n.BitLen() - 1; i >= 0; i-- {
		if n.Bit(i) == 1 {
			r0 = r0.Add(r1)
			r1 = r1.Add(r1)
		} else {
			r1 = r0.Add(r1)
			r0 = r0.Add(r0)
		}
	}
	return r0
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
