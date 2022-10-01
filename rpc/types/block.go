package types

import (
	"encoding/json"
	"errors"
	"strconv"
)

var ErrInvalidBlockID = errors.New("invalid blockid")

// BlockHashAndNumberOutput is a struct that is returned by BlockHashAndNumber.
type BlockHashAndNumberOutput struct {
	BlockNumber uint64 `json:"block_number,omitempty"`
	BlockHash   string `json:"block_hash,omitempty"`
}

// BlockID is an unexposed struct that is used in a OneOf for
// starknet_getBlockWithTxHashes.
type BlockID struct {
	Number *uint64 `json:"block_number,omitempty"`
	Hash   *Hash   `json:"block_hash,omitempty"`
	Tag    string  `json:"block_tag,omitempty"`
}

func (b BlockID) MarshalJSON() ([]byte, error) {
	if b.Tag == "pending" || b.Tag == "latest" {
		return []byte(strconv.Quote(b.Tag)), nil
	}

	type Alias BlockID
	if b.Tag != "" && (b.Tag != "pending" && b.Tag != "latest") {
		return nil, ErrInvalidBlockID
	}

	return json.Marshal((Alias)(b))
}

type Block struct {
	BlockHeader
	Status Status `json:"status"`
	// Transactions The hashes of the transactions included in this block
	Transactions Transactions `json:"transactions"`
}

type BlockHeader struct {
	// BlockHash The hash of this block
	BlockHash Hash `json:"block_hash"`
	// ParentHash The hash of this block's parent
	ParentHash Hash `json:"parent_hash"`
	// BlockNumber the block number (its height)
	BlockNumber uint64 `json:"block_number"`
	// NewRoot The new global state root
	NewRoot string `json:"new_root"`
	// Timestamp the time in which the block was created, encoded in Unix time
	Timestamp uint64 `json:"timestamp"`
	// SequencerAddress the StarkNet identity of the sequencer submitting this block
	SequencerAddress string `json:"sequencer_address"`
}
