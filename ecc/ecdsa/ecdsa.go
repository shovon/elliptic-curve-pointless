package ecdsa

import (
	"crypto/rand"
	"crypto/sha256"
	"elliptic-curve-interfaces/math/finitefield"
	"elliptic-curve-interfaces/math/primitives"
	"elliptic-curve-interfaces/maybe"
	"errors"
	"math/big"
)

type GeneratorOrder *big.Int
type PrivateKey *big.Int

func Sign[C finitefield.CurveFiniteField](
	generator finitefield.WeierstrassCurvePoint[C],
	generatorOrder GeneratorOrder,
	privateKey PrivateKey,
	message []byte,
) ([2]*big.Int, error) {
	hash := sha256.New()
	hash.Write(message)
	hashBytes := hash.Sum(nil)

	z := big.NewInt(0)
	z.SetBytes(hashBytes)

	zero := big.NewInt(0)

	var (
		k   *big.Int
		err error
	)
	k = big.NewInt(0)
	s := big.NewInt(0)
	r := big.NewInt(0)

	for s.Cmp(zero) == 0 {
		for k.Cmp(zero) == 0 {
			k, err = rand.Int(rand.Reader, generatorOrder)
			if err != nil {
				panic(err)
			}
		}

		p, ok := maybe.Extract(maybe.Maybe[finitefield.NotInfinity[C]](generator.ScalarMultiply(k)))
		if !ok {
			return [2]*big.Int{}, errors.New("failed to acquire a non-infinity point in the curve")
		}

		r = primitives.Point2D[*big.Int](p).X().Mod(r, generatorOrder)
		kInverse := k.ModInverse(k, generatorOrder)
		s = s.Mul(r, privateKey).Add(s, z).Mul(s, kInverse)
	}

	return [2]*big.Int{r, s}, nil
}
