package api

// GENERATED BY go run ./internal/cmd/genmarshal. DO NOT EDIT.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AccumulateNetwork/accumulated/internal/encoding"
)

type IdRequest struct {
	Id []byte `json:"id" form:"id" query:"id" validate:"required"`
}

type KeyPage struct {
	Height uint64 `json:"height" form:"height" query:"height" validate:"required"`
	Index  uint64 `json:"index" form:"index" query:"index" validate:"required"`
}

type MetricsRequest struct {
	Metric   string        `json:"metric" form:"metric" query:"metric" validate:"required"`
	Duration time.Duration `json:"duration" form:"duration" query:"duration" validate:"required"`
}

type MetricsResponse struct {
	Value interface{} `json:"value" form:"value" query:"value" validate:"required"`
}

type QueryMultiResponse struct {
	Items []*QueryResponse `json:"items" form:"items" query:"items" validate:"required"`
	Start uint64           `json:"start" form:"start" query:"start" validate:"required"`
	Count uint64           `json:"count" form:"count" query:"count" validate:"required"`
	Total uint64           `json:"total" form:"total" query:"total" validate:"required"`
}

type QueryResponse struct {
	Type    string      `json:"type" form:"type" query:"type" validate:"required"`
	MdRoot  []byte      `json:"mdRoot" form:"mdRoot" query:"mdRoot" validate:"required"`
	Data    interface{} `json:"data" form:"data" query:"data" validate:"required"`
	Sponsor string      `json:"sponsor" form:"sponsor" query:"sponsor" validate:"required"`
	KeyPage KeyPage     `json:"keyPage" form:"keyPage" query:"keyPage" validate:"required"`
	Txid    []byte      `json:"txid" form:"txid" query:"txid" validate:"required"`
	Signer  Signer      `json:"signer" form:"signer" query:"signer" validate:"required"`
	Sig     []byte      `json:"sig" form:"sig" query:"sig" validate:"required"`
	Status  interface{} `json:"status" form:"status" query:"status" validate:"required"`
}

type Signer struct {
	PublicKey []byte `json:"publicKey" form:"publicKey" query:"publicKey" validate:"required"`
	Nonce     uint64 `json:"nonce" form:"nonce" query:"nonce" validate:"required"`
}

type TokenDeposit struct {
	Url    string `json:"url" form:"url" query:"url" validate:"required"`
	Amount uint64 `json:"amount" form:"amount" query:"amount" validate:"required"`
	Txid   []byte `json:"txid" form:"txid" query:"txid" validate:"required"`
}

type TokenSend struct {
	From string         `json:"from" form:"from" query:"from" validate:"required"`
	To   []TokenDeposit `json:"to" form:"to" query:"to" validate:"required"`
}

type TxRequest struct {
	Sponsor        string      `json:"sponsor" form:"sponsor" query:"sponsor" validate:"required,acc-url"`
	Signer         Signer      `json:"signer" form:"signer" query:"signer" validate:"required"`
	Signature      []byte      `json:"signature" form:"signature" query:"signature" validate:"required"`
	KeyPage        KeyPage     `json:"keyPage" form:"keyPage" query:"keyPage" validate:"required"`
	WaitForDeliver bool        `json:"waitForDeliver" form:"waitForDeliver" query:"waitForDeliver" validate:"required"`
	Payload        interface{} `json:"payload" form:"payload" query:"payload" validate:"required"`
}

type TxResponse struct {
	Txid      []byte         `json:"txid" form:"txid" query:"txid" validate:"required"`
	Hash      [32]byte       `json:"hash" form:"hash" query:"hash" validate:"required"`
	Code      uint64         `json:"code" form:"code" query:"code" validate:"required"`
	Message   string         `json:"message" form:"message" query:"message" validate:"required"`
	Delivered bool           `json:"delivered" form:"delivered" query:"delivered" validate:"required"`
	Synthetic []*TxSynthetic `json:"synthetic" form:"synthetic" query:"synthetic" validate:"required"`
}

type TxSynthetic struct {
	Type string `json:"type" form:"type" query:"type" validate:"required"`
	Txid string `json:"txid" form:"txid" query:"txid" validate:"required"`
	Hash string `json:"hash" form:"hash" query:"hash" validate:"required"`
	Url  string `json:"url" form:"url" query:"url" validate:"required"`
}

type UrlRequest struct {
	Url   string `json:"url" form:"url" query:"url" validate:"required,acc-url"`
	Start uint64 `json:"start" form:"start" query:"start"`
	Count uint64 `json:"count" form:"count" query:"count"`
}

func (v *MetricsRequest) BinarySize() int {
	var n int

	n += encoding.StringBinarySize(v.Metric)

	n += encoding.DurationBinarySize(v.Duration)

	return n
}

func (v *MetricsRequest) MarshalBinary() ([]byte, error) {
	var buffer bytes.Buffer

	buffer.Write(encoding.StringMarshalBinary(v.Metric))

	buffer.Write(encoding.DurationMarshalBinary(v.Duration))

	return buffer.Bytes(), nil
}

func (v *MetricsRequest) UnmarshalBinary(data []byte) error {
	if x, err := encoding.StringUnmarshalBinary(data); err != nil {
		return fmt.Errorf("error decoding Metric: %w", err)
	} else {
		v.Metric = x
	}
	data = data[encoding.StringBinarySize(v.Metric):]

	if x, err := encoding.DurationUnmarshalBinary(data); err != nil {
		return fmt.Errorf("error decoding Duration: %w", err)
	} else {
		v.Duration = x
	}
	data = data[encoding.DurationBinarySize(v.Duration):]

	return nil
}

func (v *IdRequest) MarshalJSON() ([]byte, error) {
	var u struct {
		Id string `json:"id"`
	}
	u.Id = encoding.BytesToJSON(v.Id)
	return json.Marshal(u)
}

func (v *MetricsRequest) MarshalJSON() ([]byte, error) {
	var u struct {
		Metric   string      `json:"metric"`
		Duration interface{} `json:"duration"`
	}
	u.Metric = v.Metric
	u.Duration = encoding.DurationToJSON(v.Duration)
	return json.Marshal(u)
}

func (v *QueryResponse) MarshalJSON() ([]byte, error) {
	var u struct {
		Type    string      `json:"type"`
		MdRoot  string      `json:"mdRoot"`
		Data    interface{} `json:"data"`
		Sponsor string      `json:"sponsor"`
		KeyPage KeyPage     `json:"keyPage"`
		Txid    string      `json:"txid"`
		Signer  Signer      `json:"signer"`
		Sig     string      `json:"sig"`
		Status  interface{} `json:"status"`
	}
	u.Type = v.Type
	u.MdRoot = encoding.BytesToJSON(v.MdRoot)
	u.Data = v.Data
	u.Sponsor = v.Sponsor
	u.KeyPage = v.KeyPage
	u.Txid = encoding.BytesToJSON(v.Txid)
	u.Signer = v.Signer
	u.Sig = encoding.BytesToJSON(v.Sig)
	u.Status = v.Status
	return json.Marshal(u)
}

func (v *Signer) MarshalJSON() ([]byte, error) {
	var u struct {
		PublicKey string `json:"publicKey"`
		Nonce     uint64 `json:"nonce"`
	}
	u.PublicKey = encoding.BytesToJSON(v.PublicKey)
	u.Nonce = v.Nonce
	return json.Marshal(u)
}

func (v *TokenDeposit) MarshalJSON() ([]byte, error) {
	var u struct {
		Url    string `json:"url"`
		Amount uint64 `json:"amount"`
		Txid   string `json:"txid"`
	}
	u.Url = v.Url
	u.Amount = v.Amount
	u.Txid = encoding.BytesToJSON(v.Txid)
	return json.Marshal(u)
}

func (v *TxRequest) MarshalJSON() ([]byte, error) {
	var u struct {
		Sponsor        string      `json:"sponsor"`
		Signer         Signer      `json:"signer"`
		Signature      string      `json:"signature"`
		KeyPage        KeyPage     `json:"keyPage"`
		WaitForDeliver bool        `json:"waitForDeliver"`
		Payload        interface{} `json:"payload"`
	}
	u.Sponsor = v.Sponsor
	u.Signer = v.Signer
	u.Signature = encoding.BytesToJSON(v.Signature)
	u.KeyPage = v.KeyPage
	u.WaitForDeliver = v.WaitForDeliver
	u.Payload = v.Payload
	return json.Marshal(u)
}

func (v *TxResponse) MarshalJSON() ([]byte, error) {
	var u struct {
		Txid      string         `json:"txid"`
		Hash      string         `json:"hash"`
		Code      uint64         `json:"code"`
		Message   string         `json:"message"`
		Delivered bool           `json:"delivered"`
		Synthetic []*TxSynthetic `json:"synthetic"`
	}
	u.Txid = encoding.BytesToJSON(v.Txid)
	u.Hash = encoding.ChainToJSON(v.Hash)
	u.Code = v.Code
	u.Message = v.Message
	u.Delivered = v.Delivered
	u.Synthetic = v.Synthetic
	return json.Marshal(u)
}

func (v *IdRequest) UnmarshalJSON(data []byte) error {
	var u struct {
		Id string `json:"id"`
	}
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	if x, err := encoding.BytesFromJSON(u.Id); err != nil {
		return fmt.Errorf("error decoding Id: %w", err)
	} else {
		v.Id = x
	}
	return nil
}

func (v *MetricsRequest) UnmarshalJSON(data []byte) error {
	var u struct {
		Metric   string      `json:"metric"`
		Duration interface{} `json:"duration"`
	}
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

func (v *QueryResponse) UnmarshalJSON(data []byte) error {
	var u struct {
		Type    string      `json:"type"`
		MdRoot  string      `json:"mdRoot"`
		Data    interface{} `json:"data"`
		Sponsor string      `json:"sponsor"`
		KeyPage KeyPage     `json:"keyPage"`
		Txid    string      `json:"txid"`
		Signer  Signer      `json:"signer"`
		Sig     string      `json:"sig"`
		Status  interface{} `json:"status"`
	}
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	v.Type = u.Type
	if x, err := encoding.BytesFromJSON(u.MdRoot); err != nil {
		return fmt.Errorf("error decoding MdRoot: %w", err)
	} else {
		v.MdRoot = x
	}
	v.Data = u.Data
	v.Sponsor = u.Sponsor
	v.KeyPage = u.KeyPage
	if x, err := encoding.BytesFromJSON(u.Txid); err != nil {
		return fmt.Errorf("error decoding Txid: %w", err)
	} else {
		v.Txid = x
	}
	v.Signer = u.Signer
	if x, err := encoding.BytesFromJSON(u.Sig); err != nil {
		return fmt.Errorf("error decoding Sig: %w", err)
	} else {
		v.Sig = x
	}
	v.Status = u.Status
	return nil
}

func (v *Signer) UnmarshalJSON(data []byte) error {
	var u struct {
		PublicKey string `json:"publicKey"`
		Nonce     uint64 `json:"nonce"`
	}
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
	var u struct {
		Url    string `json:"url"`
		Amount uint64 `json:"amount"`
		Txid   string `json:"txid"`
	}
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

func (v *TxRequest) UnmarshalJSON(data []byte) error {
	var u struct {
		Sponsor        string      `json:"sponsor"`
		Signer         Signer      `json:"signer"`
		Signature      string      `json:"signature"`
		KeyPage        KeyPage     `json:"keyPage"`
		WaitForDeliver bool        `json:"waitForDeliver"`
		Payload        interface{} `json:"payload"`
	}
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	v.Sponsor = u.Sponsor
	v.Signer = u.Signer
	if x, err := encoding.BytesFromJSON(u.Signature); err != nil {
		return fmt.Errorf("error decoding Signature: %w", err)
	} else {
		v.Signature = x
	}
	v.KeyPage = u.KeyPage
	v.WaitForDeliver = u.WaitForDeliver
	v.Payload = u.Payload
	return nil
}

func (v *TxResponse) UnmarshalJSON(data []byte) error {
	var u struct {
		Txid      string         `json:"txid"`
		Hash      string         `json:"hash"`
		Code      uint64         `json:"code"`
		Message   string         `json:"message"`
		Delivered bool           `json:"delivered"`
		Synthetic []*TxSynthetic `json:"synthetic"`
	}
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
	v.Synthetic = u.Synthetic
	return nil
}
