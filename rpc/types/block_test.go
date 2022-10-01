package types

import (
	_ "embed"
	"encoding/json"
	"errors"
	"testing"
)

func TestBlockID_Marshal(t *testing.T) {
	blockNumber := uint64(420)
	for _, tc := range []struct {
		id      BlockID
		want    string
		wantErr error
	}{{
		id: BlockID{
			Tag: "latest",
		},
		want: `{"block_tag":"latest"}`,
	}, {
		id: BlockID{
			Tag: "pending",
		},
		want: `{"block_tag":"pending"}`,
	}, {
		id: BlockID{
			Tag: "bad tag",
		},
		wantErr: ErrInvalidBlockID,
	}, {
		id: BlockID{
			Number: &blockNumber,
		},
		want: `{"block_number":420}`,
	}, {
		id: func() BlockID {
			h := HexToHash("0xdead")
			return BlockID{
				Hash: &h,
			}
		}(),
		want: `{"block_hash":"0x000000000000000000000000000000000000000000000000000000000000dead"}`,
	}} {
		b, err := tc.id.MarshalJSON()
		if err != nil && tc.wantErr == nil {
			t.Errorf("marshalling block id: %v", err)
		} else if err != nil && !errors.Is(err, tc.wantErr) {
			t.Errorf("block error mismatch, want: %v, got: %v", tc.wantErr, err)
		}

		if string(b) != tc.want {
			t.Errorf("block id mismatch, want: %s, got: %s", tc.want, b)
		}
	}
}

func TestBlockStatus(t *testing.T) {
	for _, tc := range []struct {
		status string
		want   Status
	}{{
		status: "PENDING",
		want:   Status_Pending,
	}, {
		status: "ACCEPTED_ON_L2",
		want:   Status_AcceptedOnL2,
	}, {
		status: "ACCEPTED_ON_L1",
		want:   Status_AcceptedOnL1,
	}, {
		status: "REJECTED",
		want:   Status_Rejected,
	}} {
		_, err := stringToStatus(&tc.status)
		if err != nil {
			t.Errorf("%s", err)
		}
	}
}

//go:embed testdata/block.json
var rawBlock []byte

func TestBlock_Unmarshal(t *testing.T) {
	b := Block{}
	if err := json.Unmarshal(rawBlock, &b); err != nil {
		t.Fatalf("Unmarshalling block: %v", err)
	}
}
