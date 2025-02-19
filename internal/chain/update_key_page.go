package chain

import (
	"bytes"
	"fmt"

	"github.com/AccumulateNetwork/accumulate/internal/url"
	"github.com/AccumulateNetwork/accumulate/protocol"
	"github.com/AccumulateNetwork/accumulate/types"
	"github.com/AccumulateNetwork/accumulate/types/api/transactions"
)

type UpdateKeyPage struct{}

func (UpdateKeyPage) Type() types.TxType {
	return types.TxTypeUpdateKeyPage
}

func (UpdateKeyPage) Validate(st *StateManager, tx *transactions.Envelope) error {
	body := new(protocol.UpdateKeyPage)
	err := tx.As(body)
	if err != nil {
		return fmt.Errorf("invalid payload: %v", err)
	}

	page, ok := st.Origin.(*protocol.KeyPage)
	if !ok {
		return fmt.Errorf("invalid origin record: want chain type %v, got %v", types.ChainTypeKeyPage, st.Origin.Header().Type)
	}

	// We're changing the height of the key page, so reset all the nonces
	for _, key := range page.Keys {
		key.Nonce = 0
	}

	// Find the old key
	var oldKey *protocol.KeySpec
	var index int
	if len(body.Key) > 0 && body.Operation != protocol.AddKey {
		for i, key := range page.Keys {
			// We could support (double) SHA-256 but I think it's fine to
			// require the user to provide an exact match.
			if bytes.Equal(key.PublicKey, body.Key) {
				oldKey, index = key, i
				break
			}
		}
	}

	var book *protocol.KeyBook
	var priority = -1
	if page.KeyBook != "" {
		book = new(protocol.KeyBook)
		u, err := url.Parse(*page.KeyBook.AsString())
		if err != nil {
			return fmt.Errorf("invalid key book url : %s", *page.KeyBook.AsString())
		}
		err = st.LoadUrlAs(u, book)
		if err != nil {
			return fmt.Errorf("invalid key book: %v", err)
		}

		for i, p := range book.Pages {
			u, err := url.Parse(p)
			if err != nil {
				return fmt.Errorf("invalid key page url : %s", p)
			}
			if u.ResourceChain32() == st.OriginChainId {
				priority = i
			}
		}
		if priority < 0 {
			return fmt.Errorf("cannot find %q in key book with ID %X", st.OriginUrl, page.KeyBook)
		}

		// 0 is the highest priority, followed by 1, etc
		if tx.Transaction.KeyPageIndex > uint64(priority) {
			return fmt.Errorf("cannot modify %q with a lower priority key page", st.OriginUrl)
		}
	}

	if body.Owner != "" {
		_, err := url.Parse(body.Owner)
		if err != nil {
			return fmt.Errorf("invalid key book url : %s", body.Owner)
		}
	}

	switch body.Operation {
	case protocol.AddKey:
		if len(body.Key) > 0 {
			return fmt.Errorf("trying to add a new key but you gave me an existing key")
		}
		key := &protocol.KeySpec{
			PublicKey: body.NewKey,
		}
		if body.Owner != "" {
			key.Owner = body.Owner
		}
		page.Keys = append(page.Keys, key)

	case protocol.UpdateKey:
		if len(body.Key) == 0 {
			return fmt.Errorf("trying to update a new key but you didn't give me an existing key")
		}
		if oldKey == nil {
			return fmt.Errorf("no matching key found")
		}

		oldKey.PublicKey = body.NewKey
		if body.Owner != "" {
			oldKey.Owner = body.Owner
		}

	case protocol.RemoveKey:
		if len(body.Key) == 0 {
			return fmt.Errorf("trying to update a new key but you didn't give me an existing key")
		}
		if oldKey == nil {
			return fmt.Errorf("no matching key found")
		}

		page.Keys = append(page.Keys[:index], page.Keys[index+1:]...)

		if len(page.Keys) == 0 && priority == 0 {
			return fmt.Errorf("cannot delete last key of the highest priority page of a key book")
		}

	default:
		return fmt.Errorf("invalid operation: %v", body.Operation)
	}

	st.Update(page)
	return nil
}

func (UpdateKeyPage) CheckTx(st *StateManager, tx *transactions.Envelope) error {
	return UpdateKeyPage{}.Validate(st, tx)
}

func (UpdateKeyPage) DeliverTx(st *StateManager, tx *transactions.Envelope) error {
	return UpdateKeyPage{}.Validate(st, tx)
}
