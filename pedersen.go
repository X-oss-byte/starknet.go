package caigo

import (
	"fmt"
	"math/big"

	"github.com/dontpanicdao/caigo/types"
)

/*
	Hashes the contents of a given array using a golang Pedersen Hash implementation.

	(ref: https://github.com/seanjameshan/starknet.js/blob/main/src/utils/ellipticCurve.ts)
*/
func (sc StarkCurve) HashElements(elems []*types.Felt) (hash *types.Felt, err error) {
	if len(elems) == 0 {
		elems = append(elems, &types.Felt{big.NewInt(0)})
	}

	hash = &types.Felt{big.NewInt(0)}
	for _, h := range elems {
		hash, err = sc.PedersenHash([]*big.Int{hash, h})
		if err != nil {
			return hash, err
		}
	}
	return hash, err
}

/*
	Hashes the contents of a given array with its size using a golang Pedersen Hash implementation.

	(ref: https://github.com/starkware-libs/cairo-lang/blob/13cef109cd811474de114925ee61fd5ac84a25eb/src/starkware/cairo/common/hash_state.py#L6)
*/
func (sc StarkCurve) ComputeHashOnElements(elems []*types.Felt) (hash *types.Felt, err error) {
	elems = append(elems, &types.Felt{new(big.Int).SetUint64(len(elems))})
	return Curve.HashElements((elems))
}

/*
	Provides the pedersen hash of given array of big integers.
	NOTE: This function assumes the curve has been initialized with contant points

	(ref: https://github.com/seanjameshan/starknet.js/blob/main/src/utils/ellipticCurve.ts)
*/
func (sc StarkCurve) PedersenHash(elems []*types.Felt) (hash *types.Felt, err error) {
	if len(sc.ConstantPoints) == 0 {
		return hash, fmt.Errorf("must initiate precomputed constant points")
	}

	ptx := new(big.Int).Set(sc.Gx)
	pty := new(big.Int).Set(sc.Gy)
	for i, elem := range elems {
		x := new(big.Int).Set(elem.Big())

		if x.Cmp(big.NewInt(0)) != -1 && x.Cmp(sc.P) != -1 {
			return ptx, fmt.Errorf("invalid x: %v", x)
		}

		for j := 0; j < 252; j++ {
			idx := 2 + (i * 252) + j
			xin := new(big.Int).Set(sc.ConstantPoints[idx][0])
			yin := new(big.Int).Set(sc.ConstantPoints[idx][1])
			if xin.Cmp(ptx) == 0 {
				return hash, fmt.Errorf("constant point duplication: %v %v", ptx, xin)
			}
			if x.Bit(0) == 1 {
				ptx, pty = sc.Add(ptx, pty, xin, yin)
			}
			x = x.Rsh(x, 1)
		}
	}

	return &types.Felt{ptx}, nil
}
