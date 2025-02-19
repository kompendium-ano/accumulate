package api

// GENERATED BY go run ./tools/cmd/genmarshal. DO NOT EDIT.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AccumulateNetwork/accumulate/internal/encoding"
	"github.com/AccumulateNetwork/accumulate/internal/url"
	"github.com/AccumulateNetwork/accumulate/protocol"
	"github.com/AccumulateNetwork/accumulate/types/api/transactions"
)

type ChainIdQuery struct {
	ChainId []byte `json:"chainId,omitempty" form:"chainId" query:"chainId" validate:"required"`
}

type ChainQueryResponse struct {
	Type      string       `json:"type,omitempty" form:"type" query:"type" validate:"required"`
	MainChain *MerkleState `json:"mainChain,omitempty" form:"mainChain" query:"mainChain" validate:"required"`
	Data      interface{}  `json:"data,omitempty" form:"data" query:"data" validate:"required"`
	ChainId   []byte       `json:"chainId,omitempty" form:"chainId" query:"chainId" validate:"required"`
}

type DataEntry struct {
	ExtIds [][]byte `json:"extIds,omitempty" form:"extIds" query:"extIds" validate:"required"`
	Data   []byte   `json:"data,omitempty" form:"data" query:"data" validate:"required"`
}

type DataEntryQuery struct {
	Url       string   `json:"url,omitempty" form:"url" query:"url" validate:"required,acc-url"`
	EntryHash [32]byte `json:"entryHash,omitempty" form:"entryHash" query:"entryHash"`
}

type DataEntryQueryResponse struct {
	EntryHash [32]byte  `json:"entryHash,omitempty" form:"entryHash" query:"entryHash" validate:"required"`
	Entry     DataEntry `json:"entry,omitempty" form:"entry" query:"entry" validate:"required"`
}

type DataEntrySetQuery struct {
	UrlQuery
	QueryPagination
	QueryOptions
}

type DirectoryQuery struct {
	UrlQuery
	QueryPagination
	QueryOptions
}

type KeyPage struct {
	Height uint64 `json:"height,omitempty" form:"height" query:"height" validate:"required"`
	Index  uint64 `json:"index,omitempty" form:"index" query:"index"`
}

type KeyPageIndexQuery struct {
	UrlQuery
	Key []byte `json:"key,omitempty" form:"key" query:"key" validate:"required"`
}

type MerkleState struct {
	Height uint64   `json:"height,omitempty" form:"height" query:"height" validate:"required"`
	Roots  [][]byte `json:"roots,omitempty" form:"roots" query:"roots" validate:"required"`
}

type MetricsQuery struct {
	Metric   string        `json:"metric,omitempty" form:"metric" query:"metric" validate:"required"`
	Duration time.Duration `json:"duration,omitempty" form:"duration" query:"duration" validate:"required"`
}

type MetricsResponse struct {
	Value interface{} `json:"value,omitempty" form:"value" query:"value" validate:"required"`
}

type MultiResponse struct {
	Type       string        `json:"type,omitempty" form:"type" query:"type" validate:"required"`
	Items      []interface{} `json:"items,omitempty" form:"items" query:"items" validate:"required"`
	Start      uint64        `json:"start" form:"start" query:"start" validate:"required"`
	Count      uint64        `json:"count" form:"count" query:"count" validate:"required"`
	Total      uint64        `json:"total" form:"total" query:"total" validate:"required"`
	OtherItems []interface{} `json:"otherItems,omitempty" form:"otherItems" query:"otherItems" validate:"required"`
}

type QueryOptions struct {
	ExpandChains bool `json:"expandChains,omitempty" form:"expandChains" query:"expandChains"`
}

type QueryPagination struct {
	Start uint64 `json:"start,omitempty" form:"start" query:"start"`
	Count uint64 `json:"count,omitempty" form:"count" query:"count"`
}

type Signer struct {
	PublicKey []byte `json:"publicKey,omitempty" form:"publicKey" query:"publicKey" validate:"required"`
	Nonce     uint64 `json:"nonce,omitempty" form:"nonce" query:"nonce" validate:"required"`
}

type TokenDeposit struct {
	Url    string `json:"url,omitempty" form:"url" query:"url" validate:"required"`
	Amount uint64 `json:"amount,omitempty" form:"amount" query:"amount" validate:"required"`
	Txid   []byte `json:"txid,omitempty" form:"txid" query:"txid" validate:"required"`
}

type TokenSend struct {
	From string         `json:"from,omitempty" form:"from" query:"from" validate:"required"`
	To   []TokenDeposit `json:"to,omitempty" form:"to" query:"to" validate:"required"`
}

type TransactionQueryResponse struct {
	Type           string                      `json:"type,omitempty" form:"type" query:"type" validate:"required"`
	MainChain      *MerkleState                `json:"mainChain,omitempty" form:"mainChain" query:"mainChain" validate:"required"`
	Data           interface{}                 `json:"data,omitempty" form:"data" query:"data" validate:"required"`
	Origin         string                      `json:"origin,omitempty" form:"origin" query:"origin" validate:"required"`
	KeyPage        *KeyPage                    `json:"keyPage,omitempty" form:"keyPage" query:"keyPage" validate:"required"`
	Txid           []byte                      `json:"txid,omitempty" form:"txid" query:"txid" validate:"required"`
	Signatures     []*transactions.ED25519Sig  `json:"signatures,omitempty" form:"signatures" query:"signatures" validate:"required"`
	Status         *protocol.TransactionStatus `json:"status,omitempty" form:"status" query:"status" validate:"required"`
	SyntheticTxids [][32]byte                  `json:"syntheticTxids,omitempty" form:"syntheticTxids" query:"syntheticTxids" validate:"required"`
}

type TxHistoryQuery struct {
	UrlQuery
	QueryPagination
}

type TxRequest struct {
	CheckOnly bool        `json:"checkOnly,omitempty" form:"checkOnly" query:"checkOnly"`
	Origin    *url.URL    `json:"origin,omitempty" form:"origin" query:"origin" validate:"required"`
	Signer    Signer      `json:"signer,omitempty" form:"signer" query:"signer" validate:"required"`
	Signature []byte      `json:"signature,omitempty" form:"signature" query:"signature" validate:"required"`
	KeyPage   KeyPage     `json:"keyPage,omitempty" form:"keyPage" query:"keyPage" validate:"required"`
	Payload   interface{} `json:"payload,omitempty" form:"payload" query:"payload" validate:"required"`
}

type TxResponse struct {
	Txid      []byte   `json:"txid,omitempty" form:"txid" query:"txid" validate:"required"`
	Hash      [32]byte `json:"hash,omitempty" form:"hash" query:"hash" validate:"required"`
	Code      uint64   `json:"code,omitempty" form:"code" query:"code" validate:"required"`
	Message   string   `json:"message,omitempty" form:"message" query:"message" validate:"required"`
	Delivered bool     `json:"delivered,omitempty" form:"delivered" query:"delivered" validate:"required"`
}

type TxnQuery struct {
	Txid []byte        `json:"txid,omitempty" form:"txid" query:"txid" validate:"required"`
	Wait time.Duration `json:"wait,omitempty" form:"wait" query:"wait"`
}

type UrlQuery struct {
	Url string `json:"url,omitempty" form:"url" query:"url" validate:"required,acc-url"`
}

func (v *DataEntry) Equal(u *DataEntry) bool {
	if !(len(v.ExtIds) == len(u.ExtIds)) {
		return false
	}

	for i := range v.ExtIds {
		v, u := v.ExtIds[i], u.ExtIds[i]
		if !(bytes.Equal(v, u)) {
			return false
		}

	}

	if !(bytes.Equal(v.Data, u.Data)) {
		return false
	}

	return true
}

func (v *DataEntryQuery) Equal(u *DataEntryQuery) bool {
	if !(v.Url == u.Url) {
		return false
	}

	if !(v.EntryHash == u.EntryHash) {
		return false
	}

	return true
}

func (v *DataEntryQueryResponse) Equal(u *DataEntryQueryResponse) bool {
	if !(v.EntryHash == u.EntryHash) {
		return false
	}

	if !(v.Entry.Equal(&u.Entry)) {
		return false
	}

	return true
}

func (v *DataEntry) BinarySize() int {
	var n int

	n += encoding.UvarintBinarySize(uint64(len(v.ExtIds)))

	for _, v := range v.ExtIds {
		n += encoding.BytesBinarySize(v)

	}

	n += encoding.BytesBinarySize(v.Data)

	return n
}

func (v *DataEntryQuery) BinarySize() int {
	var n int

	n += encoding.StringBinarySize(v.Url)

	n += encoding.ChainBinarySize(&v.EntryHash)

	return n
}

func (v *DataEntryQueryResponse) BinarySize() int {
	var n int

	n += encoding.ChainBinarySize(&v.EntryHash)

	n += v.Entry.BinarySize()

	return n
}

func (v *DataEntry) MarshalBinary() ([]byte, error) {
	var buffer bytes.Buffer

	buffer.Write(encoding.UvarintMarshalBinary(uint64(len(v.ExtIds))))
	for i, v := range v.ExtIds {
		_ = i
		buffer.Write(encoding.BytesMarshalBinary(v))

	}

	buffer.Write(encoding.BytesMarshalBinary(v.Data))

	return buffer.Bytes(), nil
}

func (v *DataEntryQuery) MarshalBinary() ([]byte, error) {
	var buffer bytes.Buffer

	buffer.Write(encoding.StringMarshalBinary(v.Url))

	buffer.Write(encoding.ChainMarshalBinary(&v.EntryHash))

	return buffer.Bytes(), nil
}

func (v *DataEntryQueryResponse) MarshalBinary() ([]byte, error) {
	var buffer bytes.Buffer

	buffer.Write(encoding.ChainMarshalBinary(&v.EntryHash))

	if b, err := v.Entry.MarshalBinary(); err != nil {
		return nil, fmt.Errorf("error encoding Entry: %w", err)
	} else {
		buffer.Write(b)
	}

	return buffer.Bytes(), nil
}

func (v *DataEntry) UnmarshalBinary(data []byte) error {
	var lenExtIds uint64
	if x, err := encoding.UvarintUnmarshalBinary(data); err != nil {
		return fmt.Errorf("error decoding ExtIds: %w", err)
	} else {
		lenExtIds = x
	}
	data = data[encoding.UvarintBinarySize(lenExtIds):]

	v.ExtIds = make([][]byte, lenExtIds)
	for i := range v.ExtIds {
		if x, err := encoding.BytesUnmarshalBinary(data); err != nil {
			return fmt.Errorf("error decoding ExtIds[%d]: %w", i, err)
		} else {
			v.ExtIds[i] = x
		}
		data = data[encoding.BytesBinarySize(v.ExtIds[i]):]

	}

	if x, err := encoding.BytesUnmarshalBinary(data); err != nil {
		return fmt.Errorf("error decoding Data: %w", err)
	} else {
		v.Data = x
	}
	data = data[encoding.BytesBinarySize(v.Data):]

	return nil
}

func (v *DataEntryQuery) UnmarshalBinary(data []byte) error {
	if x, err := encoding.StringUnmarshalBinary(data); err != nil {
		return fmt.Errorf("error decoding Url: %w", err)
	} else {
		v.Url = x
	}
	data = data[encoding.StringBinarySize(v.Url):]

	if x, err := encoding.ChainUnmarshalBinary(data); err != nil {
		return fmt.Errorf("error decoding EntryHash: %w", err)
	} else {
		v.EntryHash = x
	}
	data = data[encoding.ChainBinarySize(&v.EntryHash):]

	return nil
}

func (v *DataEntryQueryResponse) UnmarshalBinary(data []byte) error {
	if x, err := encoding.ChainUnmarshalBinary(data); err != nil {
		return fmt.Errorf("error decoding EntryHash: %w", err)
	} else {
		v.EntryHash = x
	}
	data = data[encoding.ChainBinarySize(&v.EntryHash):]

	if err := v.Entry.UnmarshalBinary(data); err != nil {
		return fmt.Errorf("error decoding Entry: %w", err)
	}
	data = data[v.Entry.BinarySize():]

	return nil
}

func (v *ChainIdQuery) MarshalJSON() ([]byte, error) {
	u := struct {
		ChainId *string `json:"chainId,omitempty"`
	}{}
	u.ChainId = encoding.BytesToJSON(v.ChainId)
	return json.Marshal(&u)
}

func (v *ChainQueryResponse) MarshalJSON() ([]byte, error) {
	u := struct {
		Type        string       `json:"type,omitempty"`
		MainChain   *MerkleState `json:"mainChain,omitempty"`
		MerkleState *MerkleState `json:"merkleState,omitempty"`
		Data        interface{}  `json:"data,omitempty"`
		ChainId     *string      `json:"chainId,omitempty"`
	}{}
	u.Type = v.Type
	u.MainChain = v.MainChain
	u.MerkleState = v.MainChain
	u.Data = v.Data
	u.ChainId = encoding.BytesToJSON(v.ChainId)
	return json.Marshal(&u)
}

func (v *DataEntry) MarshalJSON() ([]byte, error) {
	u := struct {
		ExtIds []*string `json:"extIds,omitempty"`
		Data   *string   `json:"data,omitempty"`
	}{}
	u.ExtIds = make([]*string, len(v.ExtIds))
	for i, x := range v.ExtIds {
		u.ExtIds[i] = encoding.BytesToJSON(x)
	}
	u.Data = encoding.BytesToJSON(v.Data)
	return json.Marshal(&u)
}

func (v *DataEntryQuery) MarshalJSON() ([]byte, error) {
	u := struct {
		Url       string `json:"url,omitempty"`
		EntryHash string `json:"entryHash,omitempty"`
	}{}
	u.Url = v.Url
	u.EntryHash = encoding.ChainToJSON(v.EntryHash)
	return json.Marshal(&u)
}

func (v *DataEntryQueryResponse) MarshalJSON() ([]byte, error) {
	u := struct {
		EntryHash string    `json:"entryHash,omitempty"`
		Entry     DataEntry `json:"entry,omitempty"`
	}{}
	u.EntryHash = encoding.ChainToJSON(v.EntryHash)
	u.Entry = v.Entry
	return json.Marshal(&u)
}

func (v *KeyPageIndexQuery) MarshalJSON() ([]byte, error) {
	u := struct {
		UrlQuery
		Key *string `json:"key,omitempty"`
	}{}
	u.UrlQuery = v.UrlQuery
	u.Key = encoding.BytesToJSON(v.Key)
	return json.Marshal(&u)
}

func (v *MerkleState) MarshalJSON() ([]byte, error) {
	u := struct {
		Height uint64    `json:"height,omitempty"`
		Count  uint64    `json:"count,omitempty"`
		Roots  []*string `json:"roots,omitempty"`
	}{}
	u.Height = v.Height
	u.Count = v.Height
	u.Roots = make([]*string, len(v.Roots))
	for i, x := range v.Roots {
		u.Roots[i] = encoding.BytesToJSON(x)
	}
	return json.Marshal(&u)
}

func (v *MetricsQuery) MarshalJSON() ([]byte, error) {
	u := struct {
		Metric   string      `json:"metric,omitempty"`
		Duration interface{} `json:"duration,omitempty"`
	}{}
	u.Metric = v.Metric
	u.Duration = encoding.DurationToJSON(v.Duration)
	return json.Marshal(&u)
}

func (v *Signer) MarshalJSON() ([]byte, error) {
	u := struct {
		PublicKey *string `json:"publicKey,omitempty"`
		Nonce     uint64  `json:"nonce,omitempty"`
	}{}
	u.PublicKey = encoding.BytesToJSON(v.PublicKey)
	u.Nonce = v.Nonce
	return json.Marshal(&u)
}

func (v *TokenDeposit) MarshalJSON() ([]byte, error) {
	u := struct {
		Url    string  `json:"url,omitempty"`
		Amount uint64  `json:"amount,omitempty"`
		Txid   *string `json:"txid,omitempty"`
	}{}
	u.Url = v.Url
	u.Amount = v.Amount
	u.Txid = encoding.BytesToJSON(v.Txid)
	return json.Marshal(&u)
}

func (v *TransactionQueryResponse) MarshalJSON() ([]byte, error) {
	u := struct {
		Type           string                      `json:"type,omitempty"`
		MainChain      *MerkleState                `json:"mainChain,omitempty"`
		MerkleState    *MerkleState                `json:"merkleState,omitempty"`
		Data           interface{}                 `json:"data,omitempty"`
		Origin         string                      `json:"origin,omitempty"`
		Sponsor        string                      `json:"sponsor,omitempty"`
		KeyPage        *KeyPage                    `json:"keyPage,omitempty"`
		Txid           *string                     `json:"txid,omitempty"`
		Signatures     []*transactions.ED25519Sig  `json:"signatures,omitempty"`
		Status         *protocol.TransactionStatus `json:"status,omitempty"`
		SyntheticTxids []string                    `json:"syntheticTxids,omitempty"`
	}{}
	u.Type = v.Type
	u.MainChain = v.MainChain
	u.MerkleState = v.MainChain
	u.Data = v.Data
	u.Origin = v.Origin
	u.Sponsor = v.Origin
	u.KeyPage = v.KeyPage
	u.Txid = encoding.BytesToJSON(v.Txid)
	u.Signatures = v.Signatures
	u.Status = v.Status
	u.SyntheticTxids = encoding.ChainSetToJSON(v.SyntheticTxids)
	return json.Marshal(&u)
}

func (v *TxRequest) MarshalJSON() ([]byte, error) {
	u := struct {
		CheckOnly bool        `json:"checkOnly,omitempty"`
		Origin    *url.URL    `json:"origin,omitempty"`
		Sponsor   *url.URL    `json:"sponsor,omitempty"`
		Signer    Signer      `json:"signer,omitempty"`
		Signature *string     `json:"signature,omitempty"`
		KeyPage   KeyPage     `json:"keyPage,omitempty"`
		Payload   interface{} `json:"payload,omitempty"`
	}{}
	u.CheckOnly = v.CheckOnly
	u.Origin = v.Origin
	u.Sponsor = v.Origin
	u.Signer = v.Signer
	u.Signature = encoding.BytesToJSON(v.Signature)
	u.KeyPage = v.KeyPage
	u.Payload = v.Payload
	return json.Marshal(&u)
}

func (v *TxResponse) MarshalJSON() ([]byte, error) {
	u := struct {
		Txid      *string `json:"txid,omitempty"`
		Hash      string  `json:"hash,omitempty"`
		Code      uint64  `json:"code,omitempty"`
		Message   string  `json:"message,omitempty"`
		Delivered bool    `json:"delivered,omitempty"`
	}{}
	u.Txid = encoding.BytesToJSON(v.Txid)
	u.Hash = encoding.ChainToJSON(v.Hash)
	u.Code = v.Code
	u.Message = v.Message
	u.Delivered = v.Delivered
	return json.Marshal(&u)
}

func (v *TxnQuery) MarshalJSON() ([]byte, error) {
	u := struct {
		Txid *string     `json:"txid,omitempty"`
		Wait interface{} `json:"wait,omitempty"`
	}{}
	u.Txid = encoding.BytesToJSON(v.Txid)
	u.Wait = encoding.DurationToJSON(v.Wait)
	return json.Marshal(&u)
}

func (v *ChainIdQuery) UnmarshalJSON(data []byte) error {
	u := struct {
		ChainId *string `json:"chainId,omitempty"`
	}{}
	u.ChainId = encoding.BytesToJSON(v.ChainId)
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	if x, err := encoding.BytesFromJSON(u.ChainId); err != nil {
		return fmt.Errorf("error decoding ChainId: %w", err)
	} else {
		v.ChainId = x
	}
	return nil
}

func (v *ChainQueryResponse) UnmarshalJSON(data []byte) error {
	u := struct {
		Type        string       `json:"type,omitempty"`
		MainChain   *MerkleState `json:"mainChain,omitempty"`
		MerkleState *MerkleState `json:"merkleState,omitempty"`
		Data        interface{}  `json:"data,omitempty"`
		ChainId     *string      `json:"chainId,omitempty"`
	}{}
	u.Type = v.Type
	u.MainChain = v.MainChain
	u.MerkleState = v.MainChain
	u.Data = v.Data
	u.ChainId = encoding.BytesToJSON(v.ChainId)
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	v.Type = u.Type
	var zeroMainChain *MerkleState
	if u.MainChain != zeroMainChain {
		v.MainChain = u.MainChain
	} else {
		v.MainChain = u.MerkleState
	}
	v.Data = u.Data
	if x, err := encoding.BytesFromJSON(u.ChainId); err != nil {
		return fmt.Errorf("error decoding ChainId: %w", err)
	} else {
		v.ChainId = x
	}
	return nil
}

func (v *DataEntry) UnmarshalJSON(data []byte) error {
	u := struct {
		ExtIds []*string `json:"extIds,omitempty"`
		Data   *string   `json:"data,omitempty"`
	}{}
	u.ExtIds = make([]*string, len(v.ExtIds))
	for i, x := range v.ExtIds {
		u.ExtIds[i] = encoding.BytesToJSON(x)
	}
	u.Data = encoding.BytesToJSON(v.Data)
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	v.ExtIds = make([][]byte, len(u.ExtIds))
	for i, x := range u.ExtIds {
		if x, err := encoding.BytesFromJSON(x); err != nil {
			return fmt.Errorf("error decoding ExtIds[%d]: %w", i, err)
		} else {
			v.ExtIds[i] = x
		}
	}
	if x, err := encoding.BytesFromJSON(u.Data); err != nil {
		return fmt.Errorf("error decoding Data: %w", err)
	} else {
		v.Data = x
	}
	return nil
}

func (v *DataEntryQuery) UnmarshalJSON(data []byte) error {
	u := struct {
		Url       string `json:"url,omitempty"`
		EntryHash string `json:"entryHash,omitempty"`
	}{}
	u.Url = v.Url
	u.EntryHash = encoding.ChainToJSON(v.EntryHash)
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	v.Url = u.Url
	if x, err := encoding.ChainFromJSON(u.EntryHash); err != nil {
		return fmt.Errorf("error decoding EntryHash: %w", err)
	} else {
		v.EntryHash = x
	}
	return nil
}

func (v *DataEntryQueryResponse) UnmarshalJSON(data []byte) error {
	u := struct {
		EntryHash string    `json:"entryHash,omitempty"`
		Entry     DataEntry `json:"entry,omitempty"`
	}{}
	u.EntryHash = encoding.ChainToJSON(v.EntryHash)
	u.Entry = v.Entry
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	if x, err := encoding.ChainFromJSON(u.EntryHash); err != nil {
		return fmt.Errorf("error decoding EntryHash: %w", err)
	} else {
		v.EntryHash = x
	}
	v.Entry = u.Entry
	return nil
}

func (v *KeyPageIndexQuery) UnmarshalJSON(data []byte) error {
	u := struct {
		UrlQuery
		Key *string `json:"key,omitempty"`
	}{}
	u.UrlQuery = v.UrlQuery
	u.Key = encoding.BytesToJSON(v.Key)
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	v.UrlQuery = u.UrlQuery
	if x, err := encoding.BytesFromJSON(u.Key); err != nil {
		return fmt.Errorf("error decoding Key: %w", err)
	} else {
		v.Key = x
	}
	return nil
}

func (v *MerkleState) UnmarshalJSON(data []byte) error {
	u := struct {
		Height uint64    `json:"height,omitempty"`
		Count  uint64    `json:"count,omitempty"`
		Roots  []*string `json:"roots,omitempty"`
	}{}
	u.Height = v.Height
	u.Count = v.Height
	u.Roots = make([]*string, len(v.Roots))
	for i, x := range v.Roots {
		u.Roots[i] = encoding.BytesToJSON(x)
	}
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	var zeroHeight uint64
	if u.Height != zeroHeight {
		v.Height = u.Height
	} else {
		v.Height = u.Count
	}
	v.Roots = make([][]byte, len(u.Roots))
	for i, x := range u.Roots {
		if x, err := encoding.BytesFromJSON(x); err != nil {
			return fmt.Errorf("error decoding Roots[%d]: %w", i, err)
		} else {
			v.Roots[i] = x
		}
	}
	return nil
}

func (v *MetricsQuery) UnmarshalJSON(data []byte) error {
	u := struct {
		Metric   string      `json:"metric,omitempty"`
		Duration interface{} `json:"duration,omitempty"`
	}{}
	u.Metric = v.Metric
	u.Duration = encoding.DurationToJSON(v.Duration)
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	v.Metric = u.Metric
	if x, err := encoding.DurationFromJSON(u.Duration); err != nil {
		return fmt.Errorf("error decoding Duration: %w", err)
	} else {
		v.Duration = x
	}
	return nil
}

func (v *Signer) UnmarshalJSON(data []byte) error {
	u := struct {
		PublicKey *string `json:"publicKey,omitempty"`
		Nonce     uint64  `json:"nonce,omitempty"`
	}{}
	u.PublicKey = encoding.BytesToJSON(v.PublicKey)
	u.Nonce = v.Nonce
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	if x, err := encoding.BytesFromJSON(u.PublicKey); err != nil {
		return fmt.Errorf("error decoding PublicKey: %w", err)
	} else {
		v.PublicKey = x
	}
	v.Nonce = u.Nonce
	return nil
}

func (v *TokenDeposit) UnmarshalJSON(data []byte) error {
	u := struct {
		Url    string  `json:"url,omitempty"`
		Amount uint64  `json:"amount,omitempty"`
		Txid   *string `json:"txid,omitempty"`
	}{}
	u.Url = v.Url
	u.Amount = v.Amount
	u.Txid = encoding.BytesToJSON(v.Txid)
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	v.Url = u.Url
	v.Amount = u.Amount
	if x, err := encoding.BytesFromJSON(u.Txid); err != nil {
		return fmt.Errorf("error decoding Txid: %w", err)
	} else {
		v.Txid = x
	}
	return nil
}

func (v *TransactionQueryResponse) UnmarshalJSON(data []byte) error {
	u := struct {
		Type           string                      `json:"type,omitempty"`
		MainChain      *MerkleState                `json:"mainChain,omitempty"`
		MerkleState    *MerkleState                `json:"merkleState,omitempty"`
		Data           interface{}                 `json:"data,omitempty"`
		Origin         string                      `json:"origin,omitempty"`
		Sponsor        string                      `json:"sponsor,omitempty"`
		KeyPage        *KeyPage                    `json:"keyPage,omitempty"`
		Txid           *string                     `json:"txid,omitempty"`
		Signatures     []*transactions.ED25519Sig  `json:"signatures,omitempty"`
		Status         *protocol.TransactionStatus `json:"status,omitempty"`
		SyntheticTxids []string                    `json:"syntheticTxids,omitempty"`
	}{}
	u.Type = v.Type
	u.MainChain = v.MainChain
	u.MerkleState = v.MainChain
	u.Data = v.Data
	u.Origin = v.Origin
	u.Sponsor = v.Origin
	u.KeyPage = v.KeyPage
	u.Txid = encoding.BytesToJSON(v.Txid)
	u.Signatures = v.Signatures
	u.Status = v.Status
	u.SyntheticTxids = encoding.ChainSetToJSON(v.SyntheticTxids)
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	v.Type = u.Type
	var zeroMainChain *MerkleState
	if u.MainChain != zeroMainChain {
		v.MainChain = u.MainChain
	} else {
		v.MainChain = u.MerkleState
	}
	v.Data = u.Data
	var zeroOrigin string
	if u.Origin != zeroOrigin {
		v.Origin = u.Origin
	} else {
		v.Origin = u.Sponsor
	}
	v.KeyPage = u.KeyPage
	if x, err := encoding.BytesFromJSON(u.Txid); err != nil {
		return fmt.Errorf("error decoding Txid: %w", err)
	} else {
		v.Txid = x
	}
	v.Signatures = u.Signatures
	v.Status = u.Status
	if x, err := encoding.ChainSetFromJSON(u.SyntheticTxids); err != nil {
		return fmt.Errorf("error decoding SyntheticTxids: %w", err)
	} else {
		v.SyntheticTxids = x
	}
	return nil
}

func (v *TxRequest) UnmarshalJSON(data []byte) error {
	u := struct {
		CheckOnly bool        `json:"checkOnly,omitempty"`
		Origin    *url.URL    `json:"origin,omitempty"`
		Sponsor   *url.URL    `json:"sponsor,omitempty"`
		Signer    Signer      `json:"signer,omitempty"`
		Signature *string     `json:"signature,omitempty"`
		KeyPage   KeyPage     `json:"keyPage,omitempty"`
		Payload   interface{} `json:"payload,omitempty"`
	}{}
	u.CheckOnly = v.CheckOnly
	u.Origin = v.Origin
	u.Sponsor = v.Origin
	u.Signer = v.Signer
	u.Signature = encoding.BytesToJSON(v.Signature)
	u.KeyPage = v.KeyPage
	u.Payload = v.Payload
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	v.CheckOnly = u.CheckOnly
	var zeroOrigin *url.URL
	if u.Origin != zeroOrigin {
		v.Origin = u.Origin
	} else {
		v.Origin = u.Sponsor
	}
	v.Signer = u.Signer
	if x, err := encoding.BytesFromJSON(u.Signature); err != nil {
		return fmt.Errorf("error decoding Signature: %w", err)
	} else {
		v.Signature = x
	}
	v.KeyPage = u.KeyPage
	v.Payload = u.Payload
	return nil
}

func (v *TxResponse) UnmarshalJSON(data []byte) error {
	u := struct {
		Txid      *string `json:"txid,omitempty"`
		Hash      string  `json:"hash,omitempty"`
		Code      uint64  `json:"code,omitempty"`
		Message   string  `json:"message,omitempty"`
		Delivered bool    `json:"delivered,omitempty"`
	}{}
	u.Txid = encoding.BytesToJSON(v.Txid)
	u.Hash = encoding.ChainToJSON(v.Hash)
	u.Code = v.Code
	u.Message = v.Message
	u.Delivered = v.Delivered
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	if x, err := encoding.BytesFromJSON(u.Txid); err != nil {
		return fmt.Errorf("error decoding Txid: %w", err)
	} else {
		v.Txid = x
	}
	if x, err := encoding.ChainFromJSON(u.Hash); err != nil {
		return fmt.Errorf("error decoding Hash: %w", err)
	} else {
		v.Hash = x
	}
	v.Code = u.Code
	v.Message = u.Message
	v.Delivered = u.Delivered
	return nil
}

func (v *TxnQuery) UnmarshalJSON(data []byte) error {
	u := struct {
		Txid *string     `json:"txid,omitempty"`
		Wait interface{} `json:"wait,omitempty"`
	}{}
	u.Txid = encoding.BytesToJSON(v.Txid)
	u.Wait = encoding.DurationToJSON(v.Wait)
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	if x, err := encoding.BytesFromJSON(u.Txid); err != nil {
		return fmt.Errorf("error decoding Txid: %w", err)
	} else {
		v.Txid = x
	}
	if x, err := encoding.DurationFromJSON(u.Wait); err != nil {
		return fmt.Errorf("error decoding Wait: %w", err)
	} else {
		v.Wait = x
	}
	return nil
}
