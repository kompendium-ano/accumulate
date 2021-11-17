// Package abci implements the Accumulate ABCI applications.
//
// Transaction Processing
//
// Tendermint processes transactions in the following phases:
//
//  • BeginBlock
//  • [CheckTx]
//  • [DeliverTx]
//  • EndBlock
//  • Commit
package abci

import (
	"github.com/AccumulateNetwork/accumulate/types/api/transactions"
	types2 "github.com/tendermint/tendermint/abci/types"
	"time"

	"github.com/AccumulateNetwork/accumulate/types"
	apiQuery "github.com/AccumulateNetwork/accumulate/types/api/query"
	"github.com/AccumulateNetwork/accumulate/types/state"
)

//go:generate go run github.com/golang/mock/mockgen -source abci.go -destination ../mock/abci/abci.go

// Version is the version of the ABCI applications.
const Version uint64 = 0x1

type BeginBlockRequest struct {
	IsLeader bool
	Height   int64
	Time     time.Time
}

type EndBlockRequest struct{}

type Chain interface {
	Query(*apiQuery.Query) types2.ResponseQuery

	BeginBlock(BeginBlockRequest) types2.ResponseBeginBlock
	CheckTx(*transactions.GenTransaction) types2.ResponseCheckTx
	DeliverTx(*transactions.GenTransaction) types2.ResponseDeliverTx
	EndBlock(EndBlockRequest) types2.ResponseEndBlock
	Commit() types2.ResponseCommit
}

type State interface {
	// BlockIndex returns the current block index/height of the chain
	BlockIndex() int64

	// RootHash returns the root hash of the chain
	RootHash() []byte

	// AddStateEntry only used for genesis
	AddStateEntry(chainId *types.Bytes32, txHash *types.Bytes32, object *state.Object)

	// TODO I think this can be removed
	EnsureRootHash() []byte
}
