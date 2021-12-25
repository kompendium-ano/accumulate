package database

// GENERATED BY go run ./internal/cmd/genmarshal. DO NOT EDIT.

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/AccumulateNetwork/accumulate/internal/encoding"
	"github.com/AccumulateNetwork/accumulate/types/api/transactions"
)

type txSignatures struct {
	Signatures []*transactions.ED25519Sig `json:"signatures,omitempty" form:"signatures" query:"signatures" validate:"required"`
}

type txSyntheticTxns struct {
	Txids [][32]byte `json:"txids,omitempty" form:"txids" query:"txids" validate:"required"`
}

func (v *txSignatures) Equal(u *txSignatures) bool {
	if !(len(v.Signatures) == len(u.Signatures)) {
		return false
	}

	for i := range v.Signatures {
		v, u := v.Signatures[i], u.Signatures[i]
		if !(v.Equal(u)) {
			return false
		}

	}

	return true
}

func (v *txSyntheticTxns) Equal(u *txSyntheticTxns) bool {
	if !(len(v.Txids) == len(u.Txids)) {
		return false
	}

	for i := range v.Txids {
		if v.Txids[i] != u.Txids[i] {
			return false
		}
	}

	return true
}

func (v *txSignatures) BinarySize() int {
	var n int

	n += encoding.UvarintBinarySize(uint64(len(v.Signatures)))

	for _, v := range v.Signatures {
		n += v.BinarySize()

	}

	return n
}

func (v *txSyntheticTxns) BinarySize() int {
	var n int

	n += encoding.ChainSetBinarySize(v.Txids)

	return n
}

func (v *txSignatures) MarshalBinary() ([]byte, error) {
	var buffer bytes.Buffer

	buffer.Write(encoding.UvarintMarshalBinary(uint64(len(v.Signatures))))
	for i, v := range v.Signatures {
		_ = i
		if b, err := v.MarshalBinary(); err != nil {
			return nil, fmt.Errorf("error encoding Signatures[%d]: %w", i, err)
		} else {
			buffer.Write(b)
		}

	}

	return buffer.Bytes(), nil
}

func (v *txSyntheticTxns) MarshalBinary() ([]byte, error) {
	var buffer bytes.Buffer

	buffer.Write(encoding.ChainSetMarshalBinary(v.Txids))

	return buffer.Bytes(), nil
}

func (v *txSignatures) UnmarshalBinary(data []byte) error {
	var lenSignatures uint64
	if x, err := encoding.UvarintUnmarshalBinary(data); err != nil {
		return fmt.Errorf("error decoding Signatures: %w", err)
	} else {
		lenSignatures = x
	}
	data = data[encoding.UvarintBinarySize(lenSignatures):]

	v.Signatures = make([]*transactions.ED25519Sig, lenSignatures)
	for i := range v.Signatures {
		x := new(transactions.ED25519Sig)
		if err := x.UnmarshalBinary(data); err != nil {
			return fmt.Errorf("error decoding Signatures[%d]: %w", i, err)
		}
		data = data[x.BinarySize():]

		v.Signatures[i] = x
	}

	return nil
}

func (v *txSyntheticTxns) UnmarshalBinary(data []byte) error {
	if x, err := encoding.ChainSetUnmarshalBinary(data); err != nil {
		return fmt.Errorf("error decoding Txids: %w", err)
	} else {
		v.Txids = x
	}
	data = data[encoding.ChainSetBinarySize(v.Txids):]

	return nil
}

func (v *txSyntheticTxns) MarshalJSON() ([]byte, error) {
	u := struct {
		Txids []string `json:"txids,omitempty"`
	}{}
	u.Txids = encoding.ChainSetToJSON(v.Txids)
	return json.Marshal(&u)
}

func (v *txSyntheticTxns) UnmarshalJSON(data []byte) error {
	u := struct {
		Txids []string `json:"txids,omitempty"`
	}{}
	u.Txids = encoding.ChainSetToJSON(v.Txids)
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	if x, err := encoding.ChainSetFromJSON(u.Txids); err != nil {
		return fmt.Errorf("error decoding Txids: %w", err)
	} else {
		v.Txids = x
	}
	return nil
}
