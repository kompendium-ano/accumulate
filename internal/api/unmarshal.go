package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/AccumulateNetwork/accumulate/protocol"
	"github.com/AccumulateNetwork/accumulate/smt/common"
	"github.com/AccumulateNetwork/accumulate/smt/storage"
	"github.com/AccumulateNetwork/accumulate/types"
	"github.com/AccumulateNetwork/accumulate/types/api"
	"github.com/AccumulateNetwork/accumulate/types/api/query"
	"github.com/AccumulateNetwork/accumulate/types/api/response"
	"github.com/AccumulateNetwork/accumulate/types/api/transactions"
	"github.com/AccumulateNetwork/accumulate/types/state"
	tm "github.com/tendermint/tendermint/abci/types"
)

func responseIsError(rQuery tm.ResponseQuery) error {
	if rQuery.Code == 0 {
		return nil
	}

	switch {
	case rQuery.Code == protocol.CodeNotFound:
		return storage.ErrNotFound
	case rQuery.Log != "":
		return errors.New(rQuery.Log)
	case rQuery.Info != "":
		return errors.New(rQuery.Info)
	default:
		return fmt.Errorf("query failed with code %d", rQuery.Code)
	}
}

func unmarshalAs(rQuery tm.ResponseQuery, typ string, as func([]byte) (interface{}, error)) (*api.APIDataResponse, error) {
	if err := responseIsError(rQuery); err != nil {
		return nil, err
	}

	if rQuery.Code != 0 {
		data, err := json.Marshal(rQuery.Value)
		if err != nil {
			return nil, err
		}

		rAPI := new(api.APIDataResponse)
		rAPI.Type = types.String(typ)
		rAPI.Data = (*json.RawMessage)(&data)
		return rAPI, nil
	}

	obj := state.Object{}
	err := obj.UnmarshalBinary(rQuery.Value)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling state object %v", err)
	}

	v, err := as(obj.Entry)
	if err != nil {
		return nil, err
	}

	return respondWith(&obj, v, typ)
}

func respondWith(obj *state.Object, v interface{}, typ string) (*api.APIDataResponse, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	rAPI := new(api.APIDataResponse)
	rAPI.Type = types.String(typ)
	rAPI.MerkleState = new(api.MerkleState)
	rAPI.MerkleState.Count = obj.Height
	rAPI.MerkleState.Roots = make([]types.Bytes, len(obj.Roots))
	for i, r := range obj.Roots {
		rAPI.MerkleState.Roots[i] = r
	}
	rAPI.Data = (*json.RawMessage)(&data)
	return rAPI, nil
}

func unmarshalTokenTx(sigInfo *transactions.Header, txPayload []byte, txId types.Bytes, txSynthTxIds types.Bytes) (*api.APIDataResponse, error) {
	tx := protocol.SendTokens{}
	err := tx.UnmarshalBinary(txPayload)
	if err != nil {
		return nil, accumulateError(err)
	}

	txResp := response.TokenTx{}
	txResp.From = types.String(sigInfo.Origin.String())
	txResp.TxId = txId

	if len(txSynthTxIds)/32 != len(tx.To) {
		return nil, fmt.Errorf("number of synthetic tx, does not match number of outputs")
	}

	//should receive tx,unmarshal to output accounts
	for i, v := range tx.To {
		j := i * 32
		synthTxId := txSynthTxIds[j : j+32]
		txStatus := response.TokenTxOutputStatus{}
		txStatus.URL = types.String(v.Url)
		txStatus.TokenRecipient.Amount = v.Amount
		txStatus.SyntheticTxId = synthTxId

		txResp.ToAccount = append(txResp.ToAccount, txStatus)
	}

	data, err := json.Marshal(&txResp)
	if err != nil {
		return nil, err
	}

	resp := api.APIDataResponse{}
	resp.Type = types.String(types.TxTypeSendTokens.Name())
	resp.Data = (*json.RawMessage)(&data)
	return &resp, err
}

func unmarshalTxAs(payload []byte, v protocol.TransactionPayload) (*api.APIDataResponse, error) {
	err := v.UnmarshalBinary(payload)
	if err != nil {
		return nil, accumulateError(err)
	}

	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	resp := api.APIDataResponse{}
	resp.Type = types.String(v.GetType().Name())
	resp.Data = (*json.RawMessage)(&data)
	return &resp, err
}

//unmarshalTransaction will unpack the transaction stored on-chain and marshal it into a response
func unmarshalTransaction(sigInfo *transactions.Header, txPayload []byte, txId []byte, txSynthTxIds []byte) (resp *api.APIDataResponse, err error) {
	txType, _ := common.BytesUint64(txPayload)
	payload, err := protocol.NewTransaction(types.TransactionType(txType))
	if err != nil {
		err = fmt.Errorf("unable to unmarshal txn %x: %v", txId, err)
	}

	switch payload := payload.(type) {
	case *protocol.SendTokens:
		resp, err = unmarshalTokenTx(sigInfo, txPayload, txId, txSynthTxIds)
	default:
		resp, err = unmarshalTxAs(txPayload, payload)
	}
	if err != nil {
		return nil, err
	}

	resp.Origin = types.String(sigInfo.Origin.String())
	return resp, err
}

func unmarshalQueryResponse(rQuery tm.ResponseQuery, expect ...types.ChainType) (*api.APIDataResponse, error) {
	if err := responseIsError(rQuery); err != nil {
		return nil, err
	}

	switch typ := string(rQuery.Key); typ {
	case "tx":
		rid := query.ResponseByTxId{}
		err := rid.UnmarshalBinary(rQuery.Value)
		if err != nil {
			return nil, err
		}

		return packTransactionQuery(rid.TxId[:], rid.TxState, rid.TxPendingState, rid.TxSynthTxIds)
	case "chain":
		// OK
	default:
		return nil, fmt.Errorf("want tx or chain, got %q", typ)
	}

	obj := new(state.Object)
	err := obj.UnmarshalBinary(rQuery.Value)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling chain state object %v", err)
	}

	sChain, err := protocol.UnmarshalChain(obj.Entry)
	if err != nil {
		return nil, fmt.Errorf("invalid state object: %v", err)
	}

	if len(expect) > 0 {
		if err := isExpected(expect, sChain.Header().Type); err != nil {
			return nil, err
		}
	}

	switch sChain := sChain.(type) {
	case *protocol.ADI:
		rAdi := new(response.ADI)
		rAdi.Url = *sChain.ChainUrl.AsString()
		return respondWith(obj, rAdi, sChain.Type.String())

	case *protocol.TokenAccount:
		ta := new(protocol.CreateTokenAccount)
		ta.Url = string(sChain.ChainUrl)
		ta.TokenUrl = sChain.TokenUrl
		rAccount := response.NewTokenAccount(ta, &sChain.Balance)
		return respondWith(obj, rAccount, sChain.Type.String())

	case *protocol.LiteTokenAccount:
		rAccount := new(response.LiteTokenAccount)
		rAccount.CreateTokenAccount = new(protocol.CreateTokenAccount)
		rAccount.Url = string(sChain.ChainUrl)
		rAccount.TokenUrl = string(sChain.TokenUrl)
		rAccount.Balance = types.Amount{Int: sChain.Balance}
		rAccount.CreditBalance = types.Amount{Int: sChain.CreditBalance}
		rAccount.Nonce = sChain.Nonce
		return respondWith(obj, rAccount, sChain.Type.String())

	case *state.Transaction:
		return respondWith(obj, sChain, "tx")

	default:
		return respondWith(obj, sChain, sChain.Header().Type.String())
	}
}

func isExpected(expect []types.ChainType, typ types.ChainType) error {
	for _, e := range expect {
		if e == typ {
			return nil
		}
	}

	if len(expect) == 1 {
		return fmt.Errorf("want %v, got %v", expect[0], typ)
	}

	s := make([]string, len(expect))
	for i, e := range expect {
		s[i] = e.String()
	}
	return fmt.Errorf("want one of %s; got %v", strings.Join(s, ", "), typ)
}
