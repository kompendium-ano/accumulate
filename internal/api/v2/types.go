package api

import (
	"context"
	"time"

	"github.com/tendermint/tendermint/libs/bytes"
	core "github.com/tendermint/tendermint/rpc/core/types"
	tm "github.com/tendermint/tendermint/types"
)

//go:generate go run ../../../tools/cmd/gentypes --package api types.yml
//go:generate go run ../../../tools/cmd/genapi --package api methods.yml
//go:generate go run github.com/golang/mock/mockgen -source types.go -destination ../../mock/api/types.go

type Querier interface {
	QueryUrl(url string) (interface{}, error)
	QueryDirectory(url string, pagination QueryPagination, opts QueryOptions) (*MultiResponse, error)
	QueryChain(id []byte) (*ChainQueryResponse, error)
	QueryTx(id []byte, wait time.Duration) (*TransactionQueryResponse, error)
	QueryTxHistory(url string, start, count uint64) (*MultiResponse, error)
	QueryData(url string, entryHash [32]byte) (*ChainQueryResponse, error)
	QueryDataSet(url string, pagination QueryPagination, opts QueryOptions) (*MultiResponse, error)
	QueryKeyPageIndex(url string, key []byte) (*ChainQueryResponse, error)
}

// ABCIQueryClient is a subset of from TM/rpc/client.ABCIClient for sending
// queries.
type ABCIQueryClient interface {
	ABCIQuery(ctx context.Context, path string, data bytes.HexBytes) (*core.ResultABCIQuery, error)
}

// ABCIBroadcastClient is a subset of from TM/rpc/client.ABCIClient for
// broadcasting transactions.
type ABCIBroadcastClient interface {
	CheckTx(ctx context.Context, tx tm.Tx) (*core.ResultCheckTx, error)
	BroadcastTxAsync(context.Context, tm.Tx) (*core.ResultBroadcastTx, error)
	BroadcastTxSync(context.Context, tm.Tx) (*core.ResultBroadcastTx, error)
}
