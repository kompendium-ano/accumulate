package api

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"time"

	"github.com/AccumulateNetwork/accumulated/internal/url"
	"github.com/AccumulateNetwork/accumulated/protocol"
	"github.com/AccumulateNetwork/accumulated/types/api/transactions"
	"github.com/AccumulateNetwork/jsonrpc2/v15"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/ybbus/jsonrpc/v2"
)

func (m *JrpcMethods) Execute(ctx context.Context, params json.RawMessage) interface{} {
	var payload []byte
	req := new(TxRequest)
	req.Payload = &payload
	err := m.parse(params, req)
	if err != nil {
		return err
	}

	return m.execute(ctx, req, payload)
}

func (m *JrpcMethods) ExecuteWith(newParams func() protocol.TransactionPayload, validateFields ...string) jsonrpc2.MethodFunc {
	return func(ctx context.Context, params json.RawMessage) interface{} {
		payload := newParams()
		req := new(TxRequest)
		req.Payload = payload
		err := m.parse(params, req)
		if err != nil {
			return err
		}

		// validate request data
		if len(validateFields) == 0 {
			if err = m.validate.Struct(payload); err != nil {
				return validatorError(err)
			}
		} else {
			if err = m.validate.StructPartial(payload, validateFields...); err != nil {
				return validatorError(err)
			}
		}

		b, err := payload.MarshalBinary()
		if err != nil {
			return accumulateError(err)
		}

		return m.execute(ctx, req, b)
	}
}

// executeQueue manages queues for batching and dispatch of execute requests.
type executeQueue struct {
	leader  chan struct{}
	enqueue chan *executeRequest
}

// executeRequest captures the state of an execute requests.
type executeRequest struct {
	remote int
	params json.RawMessage
	result interface{}
	done   chan struct{}
}

// execute either executes the request locally, or dispatches it to another BVC
func (m *JrpcMethods) execute(ctx context.Context, req *TxRequest, payload []byte) interface{} {
	u, err := url.Parse(req.Sponsor)
	if err != nil {
		return validatorError(err)
	}

	// Route the request
	i := int(u.Routing() % uint64(len(m.opts.Remote)))
	if i == m.localIndex {
		// We have a local node and the routing number is local, so process the
		// request and broadcast it locally
		return m.executeLocal(ctx, req, payload)
	}

	// Prepare the request for dispatch to a remote BVC
	req.Payload = payload
	ex := new(executeRequest)
	ex.remote = i
	ex.params, err = req.MarshalJSON()
	if err != nil {
		return accumulateError(err)
	}
	ex.done = make(chan struct{})

	// Either send the request to the active dispatcher, or start a new
	// dispatcher
	select {
	case <-ctx.Done():
		// Request was canceled
		return ErrCanceled

	case m.queue.enqueue <- ex:
		// Request was accepted by the leader

	case <-m.queue.leader:
		// We are the leader, start a new dispatcher
		go m.executeBatch(ex)
	}

	// Wait for dispatch to complete
	select {
	case <-ctx.Done():
		// Canceled
		return ErrCanceled

	case <-ex.done:
		// Completed
		return ex.result
	}
}

// executeLocal constructs a TX, broadcasts it to the local node, and waits for
// results.
func (m *JrpcMethods) executeLocal(ctx context.Context, req *TxRequest, payload []byte) interface{} {
	// Build the TX
	tx := new(transactions.GenTransaction)
	tx.Transaction = payload

	tx.SigInfo = new(transactions.SignatureInfo)
	tx.SigInfo.URL = req.Sponsor
	tx.SigInfo.Nonce = req.Signer.Nonce
	tx.SigInfo.MSHeight = req.KeyPage.Height
	tx.SigInfo.PriorityIdx = req.KeyPage.Index

	ed := new(transactions.ED25519Sig)
	ed.Nonce = req.Signer.Nonce
	ed.PublicKey = req.Signer.PublicKey
	ed.Signature = req.Signature
	tx.Signature = append(tx.Signature, ed)

	txb, err := tx.Marshal()
	if err != nil {
		return accumulateError(err)
	}

	// Disable websocket based behavior if it is not enabled
	if !m.opts.EnableSubscribeTx {
		req.WaitForDeliver = false
	}

	var done chan abci.TxResult
	if req.WaitForDeliver {
		done = make(chan abci.TxResult, 1)
	}

	// Broadcast the TX
	r, err := m.opts.Local.BroadcastTxSync(ctx, txb)
	if err != nil {
		return accumulateError(err)
	}

	res := new(TxResponse)
	res.Code = uint64(r.Code)
	res.Txid = tx.TransactionHash()
	res.Hash = sha256.Sum256(txb)

	// Check for errors
	switch {
	case len(r.MempoolError) > 0:
		res.Message = r.MempoolError
		return res
	case len(r.Log) > 0:
		res.Message = r.Log
		return res
	case r.Code != 0:
		return res
	case !req.WaitForDeliver:
		return res
	}

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	// Wait for results
	select {
	case txr := <-done:
		r := txr.Result
		res.Code = uint64(r.Code)
		res.Delivered = true

		// Check for errors
		switch {
		case len(r.Log) > 0:
			res.Message = r.Log
			return res
		case len(r.Info) > 0:
			res.Message = r.Info
			return res
		case r.Code != 0:
			return res
		}

		// Process synthetic TX events
		for _, e := range r.Events {
			if e.Type != "accSyn" {
				continue
			}

			syn := new(TxSynthetic)
			res.Synthetic = append(res.Synthetic, syn)

			for _, a := range e.Attributes {
				switch a.Key {
				case "type":
					syn.Type = a.Value
				case "hash":
					syn.Txid = a.Value
				case "txRef":
					syn.Hash = a.Value
				case "url":
					syn.Url = a.Value
				}
			}
		}
		return res

	case <-timer.C:
		res.Message = "Timed out while waiting for deliver"
		return res
	}
}

// executeBatch accepts execute requests for dispatch, then dispatches requests
// in batches to the appropriate remote BVCs.
func (m *JrpcMethods) executeBatch(queue ...*executeRequest) {
	// Free the leader semaphore once we're done
	defer func() { m.queue.leader <- struct{}{} }()

	timer := time.NewTimer(m.opts.QueueDuration)
	defer timer.Stop()

	// Accept requests until we reach the target depth or the timer fires
loop:
	for {
		select {
		case <-timer.C:
			break loop
		case ex := <-m.queue.enqueue:
			queue = append(queue, ex)
			if len(queue) >= m.opts.QueueDepth {
				break loop
			}
		}
	}

	// Construct batches
	lup := map[*jsonrpc.RPCRequest]*executeRequest{}
	batches := make([]jsonrpc.RPCRequests, len(m.remote))
	for _, ex := range queue {
		rq := &jsonrpc.RPCRequest{
			Method: "execute",
			Params: ex.params,
		}
		lup[rq] = ex
		batches[ex.remote] = append(batches[ex.remote], rq)
	}

	for i, rq := range batches {
		if len(rq) == 0 {
			continue
		}

		// Send batch
		res, err := m.remote[i].CallBatch(rq)

		// Forward results
		for j := range rq {
			ex := lup[rq[j]]
			switch {
			case err != nil:
				ex.result = internalError(err)
			case res[j].Error != nil:
				err := res[j].Error
				ex.result = jsonrpc2.NewError(jsonrpc2.ErrorCode(err.Code), err.Message, err.Data)
			default:
				ex.result = res[j].Result
			}
			close(ex.done)
		}
	}
}
