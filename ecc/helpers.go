package ecc

import (
	"elliptic-curve-interfaces/math/finitefield"
	"elliptic-curve-interfaces/maybe"
	"errors"
)

func SerializeUncompressedECPointFormat[T finitefield.CurveFiniteField](
	p finitefield.WeierstrassCurvePoint[T],
) ([]byte, error) {
	point, ok := maybe.Extract(maybe.Maybe[finitefield.NotInfinity[T]](p))

	if !ok {
		return nil, errors.New("can't do much with a point at infinity, unforunately")
	}

	result := []byte{}
	result = append(result, []byte{0x04}...)
	result = append(result, point[0].Bytes()...)
	result = append(result, point[1].Bytes()...)
	return result, nil
}
