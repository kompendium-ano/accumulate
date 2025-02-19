package chain

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/AccumulateNetwork/accumulate/internal/url"
	"github.com/AccumulateNetwork/accumulate/protocol"
	"github.com/AccumulateNetwork/accumulate/smt/storage"
	"github.com/AccumulateNetwork/accumulate/types"
	"github.com/AccumulateNetwork/accumulate/types/api/transactions"
)

type AcmeFaucet struct{}

func (AcmeFaucet) Type() types.TxType { return types.TxTypeAcmeFaucet }

func (AcmeFaucet) Validate(st *StateManager, tx *transactions.Envelope) error {
	// Unmarshal the TX payload
	body := new(protocol.AcmeFaucet)
	err := tx.As(body)
	if err != nil {
		return fmt.Errorf("invalid payload: %v", err)
	}

	u, err := url.Parse(body.Url)
	if err != nil {
		return fmt.Errorf("invalid recipient URL: %v", err)
	}

	// Check the recipient
	account := new(protocol.LiteTokenAccount)
	err = st.LoadUrlAs(u, account)
	switch {
	case err == nil:
		// If the recipient exists, it must be an ACME lite account
		u, err := account.ParseTokenUrl()
		if err != nil {
			return fmt.Errorf("invalid record: bad token URL: %v", err)
		}

		if !protocol.AcmeUrl().Equal(u) {
			return fmt.Errorf("invalid recipient: %q is not an ACME account", u)
		}

	case errors.Is(err, storage.ErrNotFound):
		// If the recipient doesn't exist, ensure it is an ACME lite address
		addr, tok, err := protocol.ParseLiteAddress(u)
		switch {
		case err != nil:
			return fmt.Errorf("error parsing lite address %q: %v", u, err)
		case addr == nil:
			return fmt.Errorf("invalid recipient: %q is not a lite address", u)
		case !protocol.AcmeUrl().Equal(tok):
			return fmt.Errorf("invalid recipient: %q is not an ACME account", u)
		}

	default:
		return fmt.Errorf("invalid recipient: %v", err)
	}

	// Load the faucet state
	faucet := new(protocol.LiteTokenAccount)
	err = st.LoadUrlAs(protocol.FaucetUrl, faucet)
	if err != nil {
		return fmt.Errorf("failed to load faucet: %v", err)
	}

	// Attach this TX to the faucet (don't bother debiting)
	st.Update(faucet)

	// Submit a synthetic deposit token TX
	amount := new(big.Int).SetUint64(10 * protocol.AcmePrecision)
	deposit := new(protocol.SyntheticDepositTokens)
	copy(deposit.Cause[:], tx.Transaction.Hash())
	deposit.Token = protocol.ACME
	deposit.Amount = *amount
	st.Submit(u, deposit)

	// deposit := synthetic.NewTokenTransactionDeposit(txid[:], types.String(protocol.FaucetUrl.String()), types.String(u.String()))
	// err = deposit.SetDeposit(protocol.ACME, amount)

	return nil
}
