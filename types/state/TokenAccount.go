package state

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/AccumulateNetwork/accumulated/types"
)

type tokenAccount struct {
	Chain
	TokenUrl types.UrlChain `json:"tokenUrl"` //need to know who issued tokens, this can be condensed maybe back to adi chain path
	Balance  big.Int        `json:"balance"`  //store the balance as a big int.
}

type TokenAccount struct {
	Entry
	tokenAccount
}

//NewTokenAccount create a new token account.  Requires the identity/chain id's and coinbase if applicable
func NewTokenAccount(accountUrl string, tokenUrl string) *TokenAccount {
	tas := TokenAccount{}

	tas.SetHeader(types.String(accountUrl), types.ChainTypeTokenAccount)
	tas.TokenUrl.String = types.String(tokenUrl)

	return &tas
}

//Set will copy a token account state
func (app *TokenAccount) Set(accountState *TokenAccount) {
	if accountState == nil {
		return
	}
	app.Balance.Set(&accountState.Balance)
	app.ChainUrl = accountState.ChainUrl
	app.TokenUrl = accountState.TokenUrl
	app.Type = accountState.Type
}

// GetType is an implemented interface that returns the chain type of the object
func (app *TokenAccount) GetType() uint64 {
	return app.Type
}

// GetChainUrl returns the chain path for the object in the chain.
func (app *TokenAccount) GetChainUrl() string {
	return app.Chain.GetChainUrl()
}

// CanTransact returns true/false if there is a sufficient balance
func (app *TokenAccount) CanTransact(amt *big.Int) bool {
	//make sure the user has enough in the account to perform the transaction
	//if the balance is greater than or equal to the amount, then we are good.
	return app.GetBalance().Cmp(amt) >= 0
}

//SubBalance will subtract a balance form the account.  If this is a coinbase account,
//the balance will not be subtracted.
func (app *TokenAccount) SubBalance(amt *big.Int) error {
	if amt == nil {
		return fmt.Errorf("invalid input amount specified to subtract from balance")
	}

	if app.Balance.Cmp(amt) < 0 {
		return fmt.Errorf("insufficient balance, amount available : %d, requested, %d", app.GetBalance(), amt)
	}

	app.Balance.Sub(&app.Balance, amt)
	return nil
}

//GetBalance will return the balance of the account
func (app *TokenAccount) GetBalance() *big.Int {
	return &app.Balance
}

//AddBalance will add an amount to the balance. It only accepts a positive balance.
func (app *TokenAccount) AddBalance(amt *big.Int) error {
	if amt == nil {
		return fmt.Errorf("invalid input amount specified to add to balance")
	}

	if amt.Sign() <= 0 {
		return fmt.Errorf("amount to add to balance must be a positive amount")
	}

	app.Balance.Add(&app.Balance, amt)
	return nil
}

//MarshalBinary creates a byte array of the state object needed for storage
func (app *TokenAccount) MarshalBinary() (ret []byte, err error) {
	var buffer bytes.Buffer

	header, err := app.Chain.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buffer.Write(header)

	tokenUrlData, err := app.TokenUrl.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("cannot marshal binary for token URL in TokenAccount, %v", err)
	}
	buffer.Write(tokenUrlData)

	var buff types.Bytes32
	app.Balance.FillBytes(buff[:])
	buffer.Write(buff[:])

	return buffer.Bytes(), nil
}

//UnmarshalBinary will deserialize a byte array
func (app *TokenAccount) UnmarshalBinary(data []byte) error {
	err := app.Chain.UnmarshalBinary(data)
	if err != nil {
		return err
	}

	i := app.GetHeaderSize()

	err = app.TokenUrl.UnmarshalBinary(data[i:])
	if err != nil {
		return fmt.Errorf("unable to unmarshal binary for token account, %v", err)
	}

	i += app.TokenUrl.Size(nil)

	if len(data) < i+32 {
		return fmt.Errorf("invalid data buffer to unmarshal account state")
	}

	app.Balance.SetBytes(data[i : i+32])

	return nil
}
