package caigo

import (
	"fmt"
	"math/big"
	"testing"
)

func BenchmarkSignatureVerify(b *testing.B) {
	private, _ := Curve.GetRandomPrivateKey()
	x, y, _ := Curve.PrivateToPoint(private)

	hash, _ := Curve.PedersenHash(
		[]*big.Int{
			HexToBN("0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbddd"),
			HexToBN("0x7f15c38ea577a26f4f553282fcfe4f1feeb8ecfaad8f221ae41abf8224cbdde"),
		})

	r, s, _ := Curve.Sign(hash, private)

	b.Run(fmt.Sprintf("sign_input_size_%d", hash.BitLen()), func(b *testing.B) {
		Curve.Sign(hash, private)
	})
	b.Run(fmt.Sprintf("verify_input_size_%d", hash.BitLen()), func(b *testing.B) {
		Curve.Verify(hash, r, s, x, y)
	})
}

func TestComputeHashOnElements(t *testing.T) {
	hashEmptyArray, err := Curve.ComputeHashOnElements([]*big.Int{})
	expectedHashEmmptyArray := HexToBN("0x49ee3eba8c1600700ee1b87eb599f16716b0b1022947733551fde4050ca6804")
	if err != nil {
		t.Errorf("Could no hash an empty array %v\n", err)
	}
	if hashEmptyArray.Cmp(expectedHashEmmptyArray) != 0 {
		t.Errorf("Hash empty array wrong value. Expected %v got %v\n", expectedHashEmmptyArray, hashEmptyArray)
	}

	hashFilledArray, err := Curve.ComputeHashOnElements([]*big.Int{
		big.NewInt(123782376),
		big.NewInt(213984),
		big.NewInt(128763521321),
	})
	expectedHashFilledArray := HexToBN("0x7b422405da6571242dfc245a43de3b0fe695e7021c148b918cd9cdb462cac59")

	if err != nil {
		t.Errorf("Could no hash an array with values %v\n", err)
	}
	if hashFilledArray.Cmp(expectedHashFilledArray) != 0 {
		t.Errorf("Hash filled array wrong value. Expected %v got %v\n", expectedHashFilledArray, hashFilledArray)
	}
}
