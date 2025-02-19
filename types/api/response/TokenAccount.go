package response

import (
	"math/big"

	"github.com/AccumulateNetwork/accumulate/protocol"
	"github.com/AccumulateNetwork/accumulate/types"
)

type TokenAccount struct {
	*protocol.CreateTokenAccount
	Balance types.Amount `json:"balance" form:"balance" query:"balance"`
}

func NewTokenAccount(account *protocol.CreateTokenAccount, bal *big.Int) *TokenAccount {
	acct := &TokenAccount{}
	acct.CreateTokenAccount = account
	acct.Balance.Set(bal)
	return acct
}

//
//func (t *TokenAccount) MarshalBinary() ([]byte, error) {
//	t.TokenAccount.URL.MarshalBinary()
//	t.TokenAccount.TokenURL.MarshalBinary()
//}
