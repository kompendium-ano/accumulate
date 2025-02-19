package chain

import (
	"fmt"

	"github.com/AccumulateNetwork/accumulate/protocol"
	"github.com/AccumulateNetwork/accumulate/types"
	"github.com/AccumulateNetwork/accumulate/types/api/transactions"
)

type SyntheticDepositCredits struct{}

func (SyntheticDepositCredits) Type() types.TxType { return types.TxTypeSyntheticDepositCredits }

func (SyntheticDepositCredits) Validate(st *StateManager, tx *transactions.Envelope) error {
	body := new(protocol.SyntheticDepositCredits)
	err := tx.As(body)
	if err != nil {
		return fmt.Errorf("invalid payload: %v", err)
	}

	var account creditChain
	switch origin := st.Origin.(type) {
	case *protocol.LiteTokenAccount:
		account = origin

	case *protocol.KeyPage:
		account = origin

	default:
		return fmt.Errorf("invalid origin record: want chain type %v or %v, got %v", types.ChainTypeLiteTokenAccount, types.ChainTypeKeyPage, st.Origin.Header().Type)
	}

	account.CreditCredits(body.Amount)
	st.Update(account)
	return nil
}
