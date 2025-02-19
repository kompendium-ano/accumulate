package abci_test

import (
	"crypto/sha256"
	"fmt"
	"testing"
	"time"

	acctesting "github.com/AccumulateNetwork/accumulate/internal/testing"
	"github.com/AccumulateNetwork/accumulate/internal/testing/e2e"
	"github.com/AccumulateNetwork/accumulate/internal/url"
	"github.com/AccumulateNetwork/accumulate/protocol"
	"github.com/AccumulateNetwork/accumulate/types"
	"github.com/AccumulateNetwork/accumulate/types/api/transactions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto"
	randpkg "golang.org/x/exp/rand"
)

var rand = randpkg.New(randpkg.NewSource(0))

type Tx = transactions.Envelope

func TestEndToEndSuite(t *testing.T) {
	suite.Run(t, e2e.NewSuite(func(s *e2e.Suite) e2e.DUT {
		// Recreate the app for each test
		n := createAppWithMemDB(s.T(), crypto.Address{}, true)
		return &e2eDUT{s, n}
	}))
}

func BenchmarkFaucetAndLiteTx(b *testing.B) {
	n := createAppWithMemDB(b, crypto.Address{}, true)

	recipient := generateKey()

	n.Batch(func(send func(*Tx)) {
		tx, err := acctesting.CreateFakeSyntheticDepositTx(recipient)
		require.NoError(b, err)
		send(tx)
	})

	origin := acctesting.NewWalletEntry()
	origin.Nonce = 1
	origin.PrivateKey = recipient.Bytes()
	origin.Addr = acctesting.AcmeLiteAddressTmPriv(recipient).String()

	rwallet := acctesting.NewWalletEntry()

	b.ResetTimer()
	n.Batch(func(send func(*Tx)) {
		for i := 0; i < b.N; i++ {
			exch := new(protocol.SendTokens)
			exch.AddRecipient(n.ParseUrl(rwallet.Addr), 1000)
			tx, err := transactions.New(origin.Addr, 1, func(hash []byte) (*transactions.ED25519Sig, error) {
				return origin.Sign(hash), nil
			}, exch)
			require.NoError(b, err)
			send(tx)
		}
	})
}

func TestCreateLiteAccount(t *testing.T) {
	var count = 11
	n := createAppWithMemDB(t, crypto.Address{}, true)
	originAddr, balances := n.testLiteTx(count)
	require.Equal(t, int64(5e4*acctesting.TokenMx-count*1000), n.GetLiteTokenAccount(originAddr).Balance.Int64())
	for addr, bal := range balances {
		require.Equal(t, bal, n.GetLiteTokenAccount(addr).Balance.Int64())
	}
}

func (n *fakeNode) testLiteTx(count int) (string, map[string]int64) {
	_, recipient, gtx, err := acctesting.BuildTestSynthDepositGenTx()
	require.NoError(n.t, err)

	origin := acctesting.NewWalletEntry()
	origin.Nonce = 1
	origin.PrivateKey = recipient
	origin.Addr = acctesting.AcmeLiteAddressStdPriv(recipient).String()

	recipients := make([]*transactions.WalletEntry, 10)
	for i := range recipients {
		recipients[i] = acctesting.NewWalletEntry()
	}

	n.Batch(func(send func(*transactions.Envelope)) {
		send(gtx)
	})

	balance := map[string]int64{}
	n.Batch(func(send func(*Tx)) {
		for i := 0; i < count; i++ {
			recipient := recipients[rand.Intn(len(recipients))]
			balance[recipient.Addr] += 1000

			exch := new(protocol.SendTokens)
			exch.AddRecipient(n.ParseUrl(recipient.Addr), 1000)
			tx, err := transactions.New(origin.Addr, 1, func(hash []byte) (*transactions.ED25519Sig, error) {
				return origin.Sign(hash), nil
			}, exch)
			require.NoError(n.t, err)
			send(tx)
		}
	})

	return origin.Addr, balance
}

func TestFaucet(t *testing.T) {
	n := createAppWithMemDB(t, crypto.Address{}, true)
	alice := generateKey()
	aliceUrl := acctesting.AcmeLiteAddressTmPriv(alice).String()

	n.Batch(func(send func(*transactions.Envelope)) {
		body := new(protocol.AcmeFaucet)
		body.Url = aliceUrl
		tx, err := transactions.New(protocol.FaucetUrl.String(), 1, func(hash []byte) (*transactions.ED25519Sig, error) {
			return protocol.FaucetWallet.Sign(hash), nil
		}, body)
		require.NoError(t, err)
		send(tx)
	})

	require.Equal(t, int64(10*protocol.AcmePrecision), n.GetLiteTokenAccount(aliceUrl).Balance.Int64())
}

func TestAnchorChain(t *testing.T) {
	n := createAppWithMemDB(t, crypto.Address{}, true)
	liteAccount := generateKey()
	batch := n.db.Begin()
	require.NoError(n.t, acctesting.CreateLiteTokenAccount(batch, liteAccount, 5e4))
	require.NoError(t, batch.Commit())

	n.Batch(func(send func(*Tx)) {
		adi := new(protocol.CreateIdentity)
		adi.Url = "RoadRunner"
		adi.KeyBookName = "book"
		adi.KeyPageName = "page"

		sponsorUrl := acctesting.AcmeLiteAddressTmPriv(liteAccount).String()
		tx, err := transactions.New(sponsorUrl, 1, edSigner(liteAccount, 1), adi)
		require.NoError(t, err)

		send(tx)
	})

	// Sanity check
	require.Equal(t, types.String("acc://RoadRunner"), n.GetADI("RoadRunner").ChainUrl)

	// Get the anchor chain manager
	batch = n.db.Begin()
	defer batch.Discard()
	ledger := batch.Record(n.network.NodeUrl().JoinPath(protocol.Ledger))

	// Extract and verify the anchor chain ledgerState
	ledgerState := protocol.NewInternalLedger()
	require.NoError(t, ledger.GetStateAs(ledgerState))
	require.ElementsMatch(t, [][32]byte{
		types.Bytes((&url.URL{Authority: "RoadRunner"}).ResourceChain()).AsBytes32(),
		types.Bytes((&url.URL{Authority: "RoadRunner/book"}).ResourceChain()).AsBytes32(),
		types.Bytes((&url.URL{Authority: "RoadRunner/page"}).ResourceChain()).AsBytes32(),
	}, ledgerState.Records.Chains)

	// Check each anchor
	rootChain, err := ledger.Chain(protocol.MinorRootChain)
	require.NoError(t, err)
	first := rootChain.Height() - int64(len(ledgerState.Records.Chains))
	for i, chain := range ledgerState.Records.Chains {
		mgr, err := batch.RecordByID(chain[:]).Chain(protocol.MainChain)
		require.NoError(t, err)

		root, err := rootChain.Entry(first + int64(i))
		require.NoError(t, err)

		assert.Equal(t, mgr.Anchor(), root, "wrong anchor for %X", chain)
	}
}

func TestCreateADI(t *testing.T) {
	n := createAppWithMemDB(t, crypto.Address{}, true)

	liteAccount := generateKey()
	newAdi := generateKey()
	keyHash := sha256.Sum256(newAdi.PubKey().Address())
	batch := n.db.Begin()
	require.NoError(n.t, acctesting.CreateLiteTokenAccount(batch, liteAccount, 5e4))
	require.NoError(t, batch.Commit())

	wallet := new(transactions.WalletEntry)
	wallet.Nonce = 1
	wallet.PrivateKey = liteAccount.Bytes()
	wallet.Addr = acctesting.AcmeLiteAddressTmPriv(liteAccount).String()

	n.Batch(func(send func(*Tx)) {
		adi := new(protocol.CreateIdentity)
		adi.Url = "RoadRunner"
		adi.PublicKey = keyHash[:]
		adi.KeyBookName = "foo-book"
		adi.KeyPageName = "bar-page"

		sponsorUrl := acctesting.AcmeLiteAddressTmPriv(liteAccount).String()
		tx, err := transactions.New(sponsorUrl, 1, func(hash []byte) (*transactions.ED25519Sig, error) {
			return wallet.Sign(hash), nil
		}, adi)
		require.NoError(t, err)

		send(tx)
	})

	r := n.GetADI("RoadRunner")
	require.Equal(t, types.String("acc://RoadRunner"), r.ChainUrl)

	kg := n.GetKeyBook("RoadRunner/foo-book")
	require.Len(t, kg.Pages, 1)

	ks := n.GetKeyPage("RoadRunner/bar-page")
	require.Len(t, ks.Keys, 1)
	require.Equal(t, keyHash[:], ks.Keys[0].PublicKey)
}

func TestCreateAdiDataAccount(t *testing.T) {

	t.Run("Data Account w/ Default Key Book and no Manager Key Book", func(t *testing.T) {
		n := createAppWithMemDB(t, crypto.Address{}, true)
		adiKey := generateKey()
		batch := n.db.Begin()
		require.NoError(t, acctesting.CreateADI(batch, adiKey, "FooBar"))
		require.NoError(t, batch.Commit())

		n.Batch(func(send func(*transactions.Envelope)) {
			tac := new(protocol.CreateDataAccount)
			tac.Url = "FooBar/oof"
			tx, err := transactions.New("FooBar", 1, edSigner(adiKey, 1), tac)
			require.NoError(t, err)
			send(tx)
		})

		r := n.GetDataAccount("FooBar/oof")
		require.Equal(t, types.ChainTypeDataAccount, r.Type)
		require.Equal(t, types.String("acc://FooBar/oof"), r.ChainUrl)

		require.Contains(t, n.GetDirectory("FooBar"), n.ParseUrl("FooBar/oof").String())
	})

	t.Run("Data Account w/ Custom Key Book and Manager Key Book Url", func(t *testing.T) {
		n := createAppWithMemDB(t, crypto.Address{}, true)
		adiKey, pageKey := generateKey(), generateKey()
		batch := n.db.Begin()
		require.NoError(t, acctesting.CreateADI(batch, adiKey, "FooBar"))
		require.NoError(t, acctesting.CreateKeyPage(batch, "acc://FooBar/foo/page1", pageKey.PubKey().Bytes()))
		require.NoError(t, acctesting.CreateKeyBook(batch, "acc://FooBar/foo/book1", "acc://FooBar/foo/page1"))
		require.NoError(t, acctesting.CreateKeyPage(batch, "acc://FooBar/mgr/page1", pageKey.PubKey().Bytes()))
		require.NoError(t, acctesting.CreateKeyBook(batch, "acc://FooBar/mgr/book1", "acc://FooBar/mgr/page1"))
		require.NoError(t, batch.Commit())

		n.Batch(func(send func(*transactions.Envelope)) {
			cda := new(protocol.CreateDataAccount)
			cda.Url = "FooBar/oof"
			cda.KeyBookUrl = "acc://FooBar/foo/book1"
			cda.ManagerKeyBookUrl = "acc://FooBar/mgr/book1"
			tx, err := transactions.New("FooBar", 1, edSigner(adiKey, 1), cda)
			require.NoError(t, err)
			send(tx)
		})

		u := n.ParseUrl("acc://FooBar/foo/book1")

		r := n.GetDataAccount("FooBar/oof")
		require.Equal(t, types.ChainTypeDataAccount, r.Type)
		require.Equal(t, types.String("acc://FooBar/oof"), r.ChainUrl)
		require.Equal(t, types.String("acc://FooBar/mgr/book1"), r.ManagerKeyBook)
		require.Equal(t, types.String(u.String()), r.KeyBook)

	})

	t.Run("Data Account data entry", func(t *testing.T) {
		n := createAppWithMemDB(t, crypto.Address{}, true)
		adiKey := generateKey()
		batch := n.db.Begin()
		require.NoError(t, acctesting.CreateADI(batch, adiKey, "FooBar"))
		require.NoError(t, batch.Commit())

		n.Batch(func(send func(*transactions.Envelope)) {
			tac := new(protocol.CreateDataAccount)
			tac.Url = "FooBar/oof"
			tx, err := transactions.New("FooBar", 1, edSigner(adiKey, 1), tac)
			require.NoError(t, err)
			send(tx)
		})

		r := n.GetDataAccount("FooBar/oof")
		require.Equal(t, types.ChainTypeDataAccount, r.Type)
		require.Equal(t, types.String("acc://FooBar/oof"), r.ChainUrl)
		require.Contains(t, n.GetDirectory("FooBar"), n.ParseUrl("FooBar/oof").String())

		wd := new(protocol.WriteData)
		n.Batch(func(send func(*transactions.Envelope)) {
			for i := 0; i < 10; i++ {
				wd.Entry.ExtIds = append(wd.Entry.ExtIds, []byte(fmt.Sprintf("test id %d", i)))
			}

			wd.Entry.Data = []byte("thequickbrownfoxjumpsoverthelazydog")

			tx, err := transactions.New("FooBar/oof", 1, edSigner(adiKey, 2), wd)
			require.NoError(t, err)
			send(tx)
		})

		// Without the sleep, this test fails on Windows and macOS
		time.Sleep(3 * time.Second)

		// Test getting the data by URL
		r2 := n.GetChainDataByUrl("FooBar/oof")
		if r2 == nil {
			t.Fatalf("error getting chain data by URL")
		}

		if r2.Data == nil {
			t.Fatalf("no data returned")
		}

		rde := protocol.ResponseDataEntry{}

		err := rde.UnmarshalJSON(*r2.Data)
		if err != nil {
			t.Fatal(err)
		}

		if !rde.Entry.Equal(&wd.Entry) {
			t.Fatalf("data query does not match what was entered")
		}

		//now test query by entry hash.
		r3 := n.GetChainDataByEntryHash("FooBar/oof", wd.Entry.Hash())

		if r3.Data == nil {
			t.Fatalf("no data returned")
		}

		rde2 := protocol.ResponseDataEntry{}

		err = rde2.UnmarshalJSON(*r3.Data)
		if err != nil {
			t.Fatal(err)
		}

		if !rde.Entry.Equal(&rde2.Entry) {
			t.Fatalf("data query does not match what was entered")
		}

		//now test query by entry set
		r4 := n.GetChainDataSet("FooBar/oof", 0, 1, true)

		if r4.Data == nil {
			t.Fatalf("no data returned")
		}

		if len(r4.Data) != 1 {
			t.Fatalf("insufficent data return from set query")
		}
		rde3 := protocol.ResponseDataEntry{}
		err = rde3.UnmarshalJSON(*r4.Data[0].Data)
		if err != nil {
			t.Fatal(err)
		}

		if !rde.Entry.Equal(&rde3.Entry) {
			t.Fatalf("data query does not match what was entered")
		}

	})
}

func TestCreateAdiTokenAccount(t *testing.T) {
	t.Run("Default Key Book", func(t *testing.T) {
		n := createAppWithMemDB(t, crypto.Address{}, true)
		adiKey := generateKey()
		batch := n.db.Begin()
		require.NoError(t, acctesting.CreateADI(batch, adiKey, "FooBar"))
		require.NoError(t, batch.Commit())

		n.Batch(func(send func(*transactions.Envelope)) {
			tac := new(protocol.CreateTokenAccount)
			tac.Url = "FooBar/Baz"
			tac.TokenUrl = protocol.AcmeUrl().String()
			tx, err := transactions.New("FooBar", 1, edSigner(adiKey, 1), tac)
			require.NoError(t, err)
			send(tx)
		})

		r := n.GetTokenAccount("FooBar/Baz")
		require.Equal(t, types.ChainTypeTokenAccount, r.Type)
		require.Equal(t, types.String("acc://FooBar/Baz"), r.ChainUrl)
		require.Equal(t, protocol.AcmeUrl().String(), r.TokenUrl)

		require.Equal(t, []string{
			n.ParseUrl("FooBar").String(),
			n.ParseUrl("FooBar/book0").String(),
			n.ParseUrl("FooBar/page0").String(),
			n.ParseUrl("FooBar/Baz").String(),
		}, n.GetDirectory("FooBar"))
	})

	t.Run("Custom Key Book", func(t *testing.T) {
		n := createAppWithMemDB(t, crypto.Address{}, true)
		adiKey, pageKey := generateKey(), generateKey()
		batch := n.db.Begin()
		require.NoError(t, acctesting.CreateADI(batch, adiKey, "FooBar"))
		require.NoError(t, acctesting.CreateKeyPage(batch, "foo/page1", pageKey.PubKey().Bytes()))
		require.NoError(t, acctesting.CreateKeyBook(batch, "foo/book1", "foo/page1"))
		require.NoError(t, batch.Commit())

		n.Batch(func(send func(*transactions.Envelope)) {
			tac := new(protocol.CreateTokenAccount)
			tac.Url = "FooBar/Baz"
			tac.TokenUrl = protocol.AcmeUrl().String()
			tac.KeyBookUrl = "foo/book1"
			tx, err := transactions.New("FooBar", 1, edSigner(adiKey, 1), tac)
			require.NoError(t, err)
			send(tx)
		})

		u := n.ParseUrl("foo/book1")

		r := n.GetTokenAccount("FooBar/Baz")
		require.Equal(t, types.ChainTypeTokenAccount, r.Type)
		require.Equal(t, types.String("acc://FooBar/Baz"), r.ChainUrl)
		require.Equal(t, protocol.AcmeUrl().String(), r.TokenUrl)
		require.Equal(t, types.String(u.String()), r.KeyBook)
	})
}

func TestLiteAccountTx(t *testing.T) {
	n := createAppWithMemDB(t, crypto.Address{}, true)
	alice, bob, charlie := generateKey(), generateKey(), generateKey()
	batch := n.db.Begin()
	require.NoError(n.t, acctesting.CreateLiteTokenAccount(batch, alice, 5e4))
	require.NoError(n.t, acctesting.CreateLiteTokenAccount(batch, bob, 0))
	require.NoError(n.t, acctesting.CreateLiteTokenAccount(batch, charlie, 0))
	require.NoError(t, batch.Commit())

	aliceUrl := acctesting.AcmeLiteAddressTmPriv(alice).String()
	bobUrl := acctesting.AcmeLiteAddressTmPriv(bob).String()
	charlieUrl := acctesting.AcmeLiteAddressTmPriv(charlie).String()

	n.Batch(func(send func(*transactions.Envelope)) {
		exch := new(protocol.SendTokens)
		exch.AddRecipient(acctesting.MustParseUrl(bobUrl), 1000)
		exch.AddRecipient(acctesting.MustParseUrl(charlieUrl), 2000)

		tx, err := transactions.New(aliceUrl, 2, edSigner(alice, 1), exch)
		require.NoError(t, err)
		send(tx)
	})

	require.Equal(t, int64(5e4*acctesting.TokenMx-3000), n.GetLiteTokenAccount(aliceUrl).Balance.Int64())
	require.Equal(t, int64(1000), n.GetLiteTokenAccount(bobUrl).Balance.Int64())
	require.Equal(t, int64(2000), n.GetLiteTokenAccount(charlieUrl).Balance.Int64())
}

func TestAdiAccountTx(t *testing.T) {
	n := createAppWithMemDB(t, crypto.Address{}, true)
	fooKey, barKey := generateKey(), generateKey()
	batch := n.db.Begin()
	require.NoError(t, acctesting.CreateADI(batch, fooKey, "foo"))
	require.NoError(t, acctesting.CreateTokenAccount(batch, "foo/tokens", protocol.AcmeUrl().String(), 1, false))
	require.NoError(t, acctesting.CreateADI(batch, barKey, "bar"))
	require.NoError(t, acctesting.CreateTokenAccount(batch, "bar/tokens", protocol.AcmeUrl().String(), 0, false))
	require.NoError(t, batch.Commit())

	n.Batch(func(send func(*transactions.Envelope)) {
		exch := new(protocol.SendTokens)
		exch.AddRecipient(n.ParseUrl("bar/tokens"), 68)

		tx, err := transactions.New("foo/tokens", 1, edSigner(fooKey, 1), exch)
		require.NoError(t, err)
		send(tx)
	})

	require.Equal(t, int64(acctesting.TokenMx-68), n.GetTokenAccount("foo/tokens").Balance.Int64())
	require.Equal(t, int64(68), n.GetTokenAccount("bar/tokens").Balance.Int64())
}

func TestSendCreditsFromAdiAccountToMultiSig(t *testing.T) {
	n := createAppWithMemDB(t, crypto.Address{}, true)
	fooKey := generateKey()
	batch := n.db.Begin()
	require.NoError(t, acctesting.CreateADI(batch, fooKey, "foo"))
	require.NoError(t, acctesting.CreateTokenAccount(batch, "foo/tokens", protocol.AcmeUrl().String(), 1e2, false))
	require.NoError(t, batch.Commit())

	n.Batch(func(send func(*transactions.Envelope)) {
		ac := new(protocol.AddCredits)
		ac.Amount = 55
		ac.Recipient = "foo/page0"

		tx, err := transactions.New("foo/tokens", 1, edSigner(fooKey, 1), ac)
		require.NoError(t, err)
		send(tx)
	})

	ks := n.GetKeyPage("foo/page0")
	acct := n.GetTokenAccount("foo/tokens")
	require.Equal(t, int64(55), ks.CreditBalance.Int64())
	require.Equal(t, int64(protocol.AcmePrecision*1e2-protocol.AcmePrecision/protocol.CreditsPerFiatUnit*55), acct.Balance.Int64())
}

func TestCreateKeyPage(t *testing.T) {
	n := createAppWithMemDB(t, crypto.Address{}, true)
	fooKey, testKey := generateKey(), generateKey()
	batch := n.db.Begin()
	require.NoError(t, acctesting.CreateADI(batch, fooKey, "foo"))
	require.NoError(t, batch.Commit())

	n.Batch(func(send func(*transactions.Envelope)) {
		cms := new(protocol.CreateKeyPage)
		cms.Url = "foo/keyset1"
		cms.Keys = append(cms.Keys, &protocol.KeySpecParams{
			PublicKey: testKey.PubKey().Bytes(),
		})

		tx, err := transactions.New("foo", 1, edSigner(fooKey, 1), cms)
		require.NoError(t, err)
		send(tx)
	})

	spec := n.GetKeyPage("foo/keyset1")
	require.Len(t, spec.Keys, 1)
	key := spec.Keys[0]
	require.Equal(t, types.String(""), spec.KeyBook)
	require.Equal(t, uint64(0), key.Nonce)
	require.Equal(t, testKey.PubKey().Bytes(), key.PublicKey)
}

func TestCreateKeyBook(t *testing.T) {
	n := createAppWithMemDB(t, crypto.Address{}, true)
	fooKey, testKey := generateKey(), generateKey()
	batch := n.db.Begin()
	require.NoError(t, acctesting.CreateADI(batch, fooKey, "foo"))
	require.NoError(t, acctesting.CreateKeyPage(batch, "foo/page1", testKey.PubKey().Bytes()))
	require.NoError(t, batch.Commit())

	specUrl := n.ParseUrl("foo/page1")

	groupUrl := n.ParseUrl("foo/book1")

	n.Batch(func(send func(*transactions.Envelope)) {
		csg := new(protocol.CreateKeyBook)
		csg.Url = "foo/book1"
		csg.Pages = append(csg.Pages, specUrl.String())

		tx, err := transactions.New("foo", 1, edSigner(fooKey, 1), csg)
		require.NoError(t, err)
		send(tx)
	})

	group := n.GetKeyBook("foo/book1")
	require.Len(t, group.Pages, 1)
	require.Equal(t, specUrl.String(), group.Pages[0])

	spec := n.GetKeyPage("foo/page1")
	require.Equal(t, spec.KeyBook, types.String(groupUrl.String()))
}

func TestAddKeyPage(t *testing.T) {
	n := createAppWithMemDB(t, crypto.Address{}, true)
	fooKey, testKey1, testKey2 := generateKey(), generateKey(), generateKey()

	u := n.ParseUrl("foo/book1")

	batch := n.db.Begin()
	require.NoError(t, acctesting.CreateADI(batch, fooKey, "foo"))
	require.NoError(t, acctesting.CreateKeyPage(batch, "foo/page1", testKey1.PubKey().Bytes()))
	require.NoError(t, acctesting.CreateKeyBook(batch, "foo/book1", "foo/page1"))
	require.NoError(t, batch.Commit())

	// Sanity check
	require.Equal(t, types.String(u.String()), n.GetKeyPage("foo/page1").KeyBook)

	n.Batch(func(send func(*transactions.Envelope)) {
		cms := new(protocol.CreateKeyPage)
		cms.Url = "foo/page2"
		cms.Keys = append(cms.Keys, &protocol.KeySpecParams{
			PublicKey: testKey2.PubKey().Bytes(),
		})

		tx, err := transactions.New("foo/book1", 2, edSigner(testKey1, 1), cms)
		require.NoError(t, err)
		send(tx)
	})

	spec := n.GetKeyPage("foo/page2")
	require.Len(t, spec.Keys, 1)
	key := spec.Keys[0]
	require.Equal(t, types.String(u.String()), spec.KeyBook)
	require.Equal(t, uint64(0), key.Nonce)
	require.Equal(t, testKey2.PubKey().Bytes(), key.PublicKey)
}

func TestAddKey(t *testing.T) {
	n := createAppWithMemDB(t, crypto.Address{}, true)
	fooKey, testKey := generateKey(), generateKey()

	batch := n.db.Begin()
	require.NoError(t, acctesting.CreateADI(batch, fooKey, "foo"))
	require.NoError(t, acctesting.CreateKeyPage(batch, "foo/page1", testKey.PubKey().Bytes()))
	require.NoError(t, acctesting.CreateKeyBook(batch, "foo/book1", "foo/page1"))
	require.NoError(t, batch.Commit())

	newKey := generateKey()
	n.Batch(func(send func(*transactions.Envelope)) {
		body := new(protocol.UpdateKeyPage)
		body.Operation = protocol.AddKey
		body.NewKey = newKey.PubKey().Bytes()

		tx, err := transactions.New("foo/page1", 2, edSigner(testKey, 1), body)
		require.NoError(t, err)
		send(tx)
	})

	spec := n.GetKeyPage("foo/page1")
	require.Len(t, spec.Keys, 2)
	require.Equal(t, newKey.PubKey().Bytes(), spec.Keys[1].PublicKey)
}

func TestUpdateKey(t *testing.T) {
	n := createAppWithMemDB(t, crypto.Address{}, true)
	fooKey, testKey := generateKey(), generateKey()

	batch := n.db.Begin()
	require.NoError(t, acctesting.CreateADI(batch, fooKey, "foo"))
	require.NoError(t, acctesting.CreateKeyPage(batch, "foo/page1", testKey.PubKey().Bytes()))
	require.NoError(t, acctesting.CreateKeyBook(batch, "foo/book1", "foo/page1"))
	require.NoError(t, batch.Commit())

	newKey := generateKey()
	n.Batch(func(send func(*transactions.Envelope)) {
		body := new(protocol.UpdateKeyPage)
		body.Operation = protocol.UpdateKey
		body.Key = testKey.PubKey().Bytes()
		body.NewKey = newKey.PubKey().Bytes()

		tx, err := transactions.New("foo/page1", 2, edSigner(testKey, 1), body)
		require.NoError(t, err)
		send(tx)
	})

	spec := n.GetKeyPage("foo/page1")
	require.Len(t, spec.Keys, 1)
	require.Equal(t, newKey.PubKey().Bytes(), spec.Keys[0].PublicKey)
}

func TestRemoveKey(t *testing.T) {
	n := createAppWithMemDB(t, crypto.Address{}, true)
	fooKey, testKey1, testKey2 := generateKey(), generateKey(), generateKey()

	batch := n.db.Begin()
	require.NoError(t, acctesting.CreateADI(batch, fooKey, "foo"))
	require.NoError(t, acctesting.CreateKeyPage(batch, "foo/page1", testKey1.PubKey().Bytes(), testKey2.PubKey().Bytes()))
	require.NoError(t, acctesting.CreateKeyBook(batch, "foo/book1", "foo/page1"))
	require.NoError(t, batch.Commit())

	n.Batch(func(send func(*transactions.Envelope)) {
		body := new(protocol.UpdateKeyPage)
		body.Operation = protocol.RemoveKey
		body.Key = testKey1.PubKey().Bytes()

		tx, err := transactions.New("foo/page1", 2, edSigner(testKey2, 1), body)
		require.NoError(t, err)
		send(tx)
	})

	spec := n.GetKeyPage("foo/page1")
	require.Len(t, spec.Keys, 1)
	require.Equal(t, testKey2.PubKey().Bytes(), spec.Keys[0].PublicKey)
}

func TestSignatorHeight(t *testing.T) {
	n := createAppWithMemDB(t, crypto.Address{}, true)
	liteKey, fooKey := generateKey(), generateKey()

	liteUrl, err := protocol.LiteAddress(liteKey.PubKey().Bytes(), "ACME")
	require.NoError(t, err)
	tokenUrl, err := url.Parse("foo/tokens")
	require.NoError(t, err)
	keyPageUrl, err := url.Parse("foo/page0")
	require.NoError(t, err)

	batch := n.db.Begin()
	require.NoError(t, acctesting.CreateLiteTokenAccount(batch, liteKey, 1))
	require.NoError(t, batch.Commit())

	getHeight := func(u *url.URL) uint64 {
		batch := n.db.Begin()
		defer batch.Discard()
		chain, err := batch.Record(u).Chain(protocol.MainChain)
		require.NoError(t, err)
		return uint64(chain.Height())
	}

	liteHeight := getHeight(liteUrl)

	n.Batch(func(send func(*transactions.Envelope)) {
		adi := new(protocol.CreateIdentity)
		adi.Url = "foo"
		adi.PublicKey = fooKey.PubKey().Bytes()
		adi.KeyBookName = "book"
		adi.KeyPageName = "page0"

		tx, err := transactions.New(liteUrl.String(), 1, edSigner(liteKey, 1), adi)
		require.NoError(t, err)
		send(tx)
	})

	require.Equal(t, liteHeight, getHeight(liteUrl), "Lite account height changed")

	keyPageHeight := getHeight(keyPageUrl)

	n.Batch(func(send func(*transactions.Envelope)) {
		tac := new(protocol.CreateTokenAccount)
		tac.Url = tokenUrl.String()
		tac.TokenUrl = protocol.AcmeUrl().String()
		tx, err := transactions.New("foo", 1, edSigner(fooKey, 1), tac)
		require.NoError(t, err)
		send(tx)
	})

	require.Equal(t, keyPageHeight, getHeight(keyPageUrl), "Key page height changed")
}
