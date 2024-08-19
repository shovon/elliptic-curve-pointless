package ecdsa

import (
	"crypto/rand"
	"elliptic-curve-interfaces/math/finitefield"
	"elliptic-curve-interfaces/math/primitives"
	"elliptic-curve-interfaces/maybe"
	"errors"
	"math/big"
)

type GeneratorOrder *big.Int
type PrivateKey *big.Int

func extractLeftMostBits(num *big.Int, n int) *big.Int {
	// Get the total number of bits in the big.Int number
	totalBits := num.BitLen()

	// If n is greater than or equal to the total number of bits, return the number as is
	if n >= totalBits {
		return new(big.Int).Set(num)
	}

	// Calculate the number of bits to shift right
	bitsToShift := totalBits - n

	// Shift the number to the right to discard the least significant bits
	leftMostBits := new(big.Int).Rsh(num, uint(bitsToShift))

	return leftMostBits
}

func Sign[C finitefield.CurveFiniteField](
	hasher func(input []byte) ([]byte, error),
	generator finitefield.WeierstrassCurvePoint[C],
	generatorOrder GeneratorOrder,
	privateKey PrivateKey,
	message []byte,
) (*big.Int, *big.Int, error) {
	hashBytes, err := hasher(message)
	if err != nil {
		return nil, nil, err
	}

	z := extractLeftMostBits(big.NewInt(0).SetBytes(hashBytes), (*big.Int)(generatorOrder).BitLen())

	var k *big.Int
	k = big.NewInt(0)
	s := big.NewInt(0)
	r := big.NewInt(0)

	for s.Cmp(big.NewInt(0)) == 0 || r.Cmp(big.NewInt(0)) == 0 {
		for k.Cmp(big.NewInt(0)) == 0 {
			k, err = rand.Int(rand.Reader, generatorOrder)
			if err != nil {
				panic(err)
			}
		}

		p, ok := maybe.Extract(maybe.Maybe[finitefield.NotInfinity[C]](generator.ScalarMultiply(k)))
		if !ok {
			return nil, nil, errors.New("failed to acquire a non-infinity point in the curve")
		}

		x := big.NewInt(0).SetBytes(primitives.Point2D[*big.Int](p).X().Bytes())
		r = r.Mod(x, generatorOrder)
		kInverse := k.ModInverse(k, generatorOrder)
		s = s.Mul(r, privateKey).Add(s, z).Mul(s, kInverse)
		s.Mod(s, generatorOrder)
	}

	return r, s, nil
}

func Verify[C finitefield.CurveFiniteField](
	hasher func(input []byte) ([]byte, error),

	generator finitefield.WeierstrassCurvePoint[C],
	generatorOrder GeneratorOrder,
	publicKey finitefield.WeierstrassCurvePoint[C],
	signature [2]*big.Int,

	message []byte,
) (bool, error) {

	// TODO: check that the public key is not the point at infinity
	// TODO: check that the public key lies on the curve
	// TODO: check that generatorOrder*publicKey = point at infinity

	hashBytes, err := hasher(message)
	if err != nil {
		return false, err
	}

	z := extractLeftMostBits(big.NewInt(0).SetBytes(hashBytes), (*big.Int)(generatorOrder).BitLen())
	sInverse := big.NewInt(0).ModInverse(signature[1], generatorOrder)

	u1 := big.NewInt(0).Mod(big.NewInt(0).Mul(z, sInverse), generatorOrder)
	u2 := big.NewInt(0).Mod(big.NewInt(0).Mul(signature[0], sInverse), generatorOrder)

	result, ok := maybe.Extract(maybe.Maybe[finitefield.NotInfinity[C]](generator.ScalarMultiply(u1).Add(publicKey.ScalarMultiply(u2))))
	if !ok {
		return false, nil
	}

	return signature[0].Cmp(result[0]) == 0, nil
}
