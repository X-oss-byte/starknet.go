package types

import (
	"testing"
)

func TestTransactionStatus(t *testing.T) {
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
