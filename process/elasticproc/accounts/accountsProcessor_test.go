package accounts

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data/alteredAccount"
	"github.com/TerraDharitri/drt-go-chain-es-indexer/data"
	"github.com/TerraDharitri/drt-go-chain-es-indexer/mock"
	"github.com/TerraDharitri/drt-go-chain-es-indexer/process/dataindexer"
	"github.com/TerraDharitri/drt-go-chain-es-indexer/process/elasticproc/converters"
	"github.com/TerraDharitri/drt-go-chain-es-indexer/process/elasticproc/tags"
	"github.com/stretchr/testify/require"
)

var balanceConverter, _ = converters.NewBalanceConverter(10)

func TestNewAccountsProcessor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		argsFunc func() (core.PubkeyConverter, dataindexer.BalanceConverter)
		exError  error
	}{
		{
			name: "NilBalanceConverter",
			argsFunc: func() (core.PubkeyConverter, dataindexer.BalanceConverter) {
				return &mock.PubkeyConverterMock{}, nil
			},
			exError: dataindexer.ErrNilBalanceConverter,
		},
		{
			name: "NilPubKeyConverter",
			argsFunc: func() (core.PubkeyConverter, dataindexer.BalanceConverter) {
				return nil, balanceConverter
			},
			exError: dataindexer.ErrNilPubkeyConverter,
		},
		{
			name: "ShouldWork",
			argsFunc: func() (core.PubkeyConverter, dataindexer.BalanceConverter) {
				return &mock.PubkeyConverterMock{}, balanceConverter
			},
			exError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAccountsProcessor(tt.argsFunc())
			require.True(t, errors.Is(err, tt.exError))
		})
	}
}

func TestAccountsProcessor_GetAccountsWithNil(t *testing.T) {
	t.Parallel()

	ap, _ := NewAccountsProcessor(mock.NewPubkeyConverterMock(32), balanceConverter)

	regularAccounts, dcdtAccounts := ap.GetAccounts(nil)
	require.Len(t, regularAccounts, 0)
	require.Len(t, dcdtAccounts, 0)
}

func TestAccountsProcessor_PrepareRegularAccountsMapWithNil(t *testing.T) {
	t.Parallel()

	ap, _ := NewAccountsProcessor(mock.NewPubkeyConverterMock(32), balanceConverter)

	accountsInfo := ap.PrepareRegularAccountsMap(0, nil, 0)
	require.Len(t, accountsInfo, 0)
}

func TestGetDCDTInfo(t *testing.T) {
	t.Parallel()

	ap, _ := NewAccountsProcessor(mock.NewPubkeyConverterMock(32), balanceConverter)
	require.NotNil(t, ap)

	tokenIdentifier := "token-001"
	wrapAccount := &data.AccountDCDT{
		Account: &alteredAccount.AlteredAccount{
			Address: "",
			Balance: "1000",
			Nonce:   0,
			Tokens: []*alteredAccount.AccountTokenData{
				{
					Identifier: tokenIdentifier,
					Balance:    "1000",
					Properties: "6f6b",
				},
			},
		},
		TokenIdentifier: tokenIdentifier,
	}
	balance, prop, _, err := ap.getDCDTInfo(wrapAccount)
	require.Nil(t, err)
	require.Equal(t, big.NewInt(1000), balance)
	require.Equal(t, "6f6b", prop)
}

func TestGetDCDTInfoNFT(t *testing.T) {
	t.Parallel()

	ap, _ := NewAccountsProcessor(mock.NewPubkeyConverterMock(32), balanceConverter)
	require.NotNil(t, ap)

	tokenIdentifier := "token-001"
	wrapAccount := &data.AccountDCDT{
		Account: &alteredAccount.AlteredAccount{
			Address: "",
			Balance: "1",
			Nonce:   10,
			Tokens: []*alteredAccount.AccountTokenData{
				{
					Identifier: tokenIdentifier,
					Balance:    "1",
					Nonce:      10,
					Properties: "6f6b",
				},
			},
		},
		TokenIdentifier: tokenIdentifier,
		IsNFTOperation:  true,
		NFTNonce:        10,
	}
	balance, prop, _, err := ap.getDCDTInfo(wrapAccount)
	require.Nil(t, err)
	require.Equal(t, big.NewInt(1), balance)
	require.Equal(t, "6f6b", prop)
}

func TestGetDCDTInfoNFTWithMetaData(t *testing.T) {
	t.Parallel()

	pubKeyConverter := mock.NewPubkeyConverterMock(32)
	ap, _ := NewAccountsProcessor(pubKeyConverter, balanceConverter)
	require.NotNil(t, ap)

	nftName := "Test-nft"
	creator := "010101"

	tokenIdentifier := "token-001"
	wrapAccount := &data.AccountDCDT{
		Account: &alteredAccount.AlteredAccount{
			Address: "",
			Balance: "1",
			Nonce:   1,
			Tokens: []*alteredAccount.AccountTokenData{
				{
					Identifier: tokenIdentifier,
					Balance:    "1",
					Properties: "6f6b",
					Nonce:      10,
					MetaData: &alteredAccount.TokenMetaData{
						Nonce:     10,
						Name:      nftName,
						Creator:   creator,
						Royalties: 2,
					},
				},
			},
		},
		TokenIdentifier: tokenIdentifier,
		IsNFTOperation:  true,
		NFTNonce:        10,
	}
	balance, prop, metaData, err := ap.getDCDTInfo(wrapAccount)
	require.Nil(t, err)
	require.Equal(t, big.NewInt(1), balance)
	require.Equal(t, "6f6b", prop)
	require.Equal(t, &data.TokenMetaData{
		Name:       nftName,
		Creator:    creator,
		Royalties:  2,
		Attributes: make([]byte, 0),
	}, metaData)
}

func TestAccountsProcessor_GetAccountsREWAAccounts(t *testing.T) {
	t.Parallel()

	addr := "aaaabbbb"
	acc := &alteredAccount.AlteredAccount{}
	alteredAccountsMap := map[string]*alteredAccount.AlteredAccount{
		addr: acc,
	}
	ap, _ := NewAccountsProcessor(mock.NewPubkeyConverterMock(32), balanceConverter)
	require.NotNil(t, ap)

	accounts, dcdtAccounts := ap.GetAccounts(alteredAccountsMap)
	require.Equal(t, 0, len(dcdtAccounts))
	require.Equal(t, []*data.Account{
		{
			UserAccount: acc,
		},
	}, accounts)
}

func TestAccountsProcessor_GetAccountsDCDTAccount(t *testing.T) {
	t.Parallel()

	addr := "aaaabbbb"
	acc := &alteredAccount.AlteredAccount{
		Address: addr,
		Balance: "1",
		Tokens: []*alteredAccount.AccountTokenData{
			{
				Identifier: "token",
			},
		},
	}
	alteredAccountsMap := map[string]*alteredAccount.AlteredAccount{
		addr: acc,
	}
	ap, _ := NewAccountsProcessor(mock.NewPubkeyConverterMock(32), balanceConverter)
	require.NotNil(t, ap)

	accounts, dcdtAccounts := ap.GetAccounts(alteredAccountsMap)
	require.Equal(t, 0, len(accounts))
	require.Equal(t, []*data.AccountDCDT{
		{Account: acc, TokenIdentifier: "token"},
	}, dcdtAccounts)
}

func TestAccountsProcessor_GetAccountsDCDTAccountNewAccountShouldBeInRegularAccounts(t *testing.T) {
	t.Parallel()

	addr := "aaaabbbb"
	acc := &alteredAccount.AlteredAccount{
		Tokens: []*alteredAccount.AccountTokenData{
			{
				Identifier: "token",
			},
		},
	}
	alteredAccountsMap := map[string]*alteredAccount.AlteredAccount{
		addr: acc,
	}
	ap, _ := NewAccountsProcessor(mock.NewPubkeyConverterMock(32), balanceConverter)
	require.NotNil(t, ap)

	accounts, dcdtAccounts := ap.GetAccounts(alteredAccountsMap)
	require.Equal(t, 1, len(accounts))
	require.Equal(t, []*data.AccountDCDT{
		{Account: acc, TokenIdentifier: "token"},
	}, dcdtAccounts)

	require.Equal(t, []*data.Account{
		{UserAccount: acc, IsSender: false},
	}, accounts)
}

func TestAccountsProcessor_PrepareAccountsMapREWA(t *testing.T) {
	t.Parallel()

	addrBytes := bytes.Repeat([]byte{0}, 32)
	addr := hex.EncodeToString(addrBytes)

	account := &alteredAccount.AlteredAccount{
		Address: addr,
		Balance: "1000",
		Nonce:   1,
		AdditionalData: &alteredAccount.AdditionalAccountData{
			CodeHash:     []byte("code"),
			CodeMetadata: []byte("metadata"),
			RootHash:     []byte("root"),
		},
	}

	rewaAccount := &data.Account{
		UserAccount: account,
		IsSender:    false,
	}

	ap, _ := NewAccountsProcessor(mock.NewPubkeyConverterMock(32), balanceConverter)
	require.NotNil(t, ap)

	balanceNum, _ := balanceConverter.ComputeBalanceAsFloat(big.NewInt(1000))

	res := ap.PrepareRegularAccountsMap(123, []*data.Account{rewaAccount}, 0)
	require.Equal(t, &data.AccountInfo{
		Address:         addr,
		Nonce:           1,
		Balance:         "1000",
		BalanceNum:      balanceNum,
		IsSmartContract: true,
		Timestamp:       time.Duration(123),
		CodeHash:        []byte("code"),
		CodeMetadata:    []byte("metadata"),
		RootHash:        []byte("root"),
	},
		res[addr])
}

func TestAccountsProcessor_PrepareAccountsMapDCDT(t *testing.T) {
	t.Parallel()

	addr := "aaaabbbb"

	account := &alteredAccount.AlteredAccount{
		Address: hex.EncodeToString([]byte(addr)),
		Tokens: []*alteredAccount.AccountTokenData{
			{
				Balance:    "1000",
				Identifier: "token",
				Nonce:      15,
				Properties: "3032",
				MetaData: &alteredAccount.TokenMetaData{
					Creator: "creator",
				},
			},
			{
				Balance:    "1000",
				Identifier: "token",
				Nonce:      16,
				Properties: "3032",
				MetaData: &alteredAccount.TokenMetaData{
					Creator: "creator",
				},
			},
		},
	}
	ap, _ := NewAccountsProcessor(mock.NewPubkeyConverterMock(32), balanceConverter)
	require.NotNil(t, ap)

	accountsDCDT := []*data.AccountDCDT{
		{Account: account, TokenIdentifier: "token", IsNFTOperation: true, NFTNonce: 15},
		{Account: account, TokenIdentifier: "token", IsNFTOperation: true, NFTNonce: 16},
	}

	tagsCount := tags.NewTagsCount()
	res, _ := ap.PrepareAccountsMapDCDT(123, accountsDCDT, tagsCount, 0)
	require.Len(t, res, 2)

	balanceNum, _ := balanceConverter.ComputeBalanceAsFloat(big.NewInt(1000))
	require.Equal(t, &data.AccountInfo{
		Address:         hex.EncodeToString([]byte(addr)),
		Balance:         "1000",
		BalanceNum:      balanceNum,
		TokenName:       "token",
		TokenIdentifier: "token-0f",
		Properties:      "3032",
		TokenNonce:      15,
		Data: &data.TokenMetaData{
			Creator:    "creator",
			Attributes: make([]byte, 0),
		},
		Timestamp: time.Duration(123),
	}, res[hex.EncodeToString([]byte(addr))+"-token-15"])

	require.Equal(t, &data.AccountInfo{
		Address:         hex.EncodeToString([]byte(addr)),
		Balance:         "1000",
		BalanceNum:      balanceNum,
		TokenName:       "token",
		TokenIdentifier: "token-10",
		Properties:      "3032",
		TokenNonce:      16,
		Data: &data.TokenMetaData{
			Creator:    "creator",
			Attributes: make([]byte, 0),
		},
		Timestamp: time.Duration(123),
	}, res[hex.EncodeToString([]byte(addr))+"-token-16"])
}

func TestAccountsProcessor_PrepareAccountsHistory(t *testing.T) {
	t.Parallel()

	accounts := map[string]*data.AccountInfo{
		"addr1": {
			Address:    "addr1",
			Balance:    "112",
			TokenName:  "token-112",
			TokenNonce: 10,
			IsSender:   true,
		},
	}

	ap, _ := NewAccountsProcessor(mock.NewPubkeyConverterMock(32), balanceConverter)

	res := ap.PrepareAccountsHistory(100, accounts, 0)
	accountBalanceHistory := res["addr1-token-112-10"]
	require.Equal(t, &data.AccountBalanceHistory{
		Address:    "addr1",
		Timestamp:  100,
		Balance:    "112",
		Token:      "token-112",
		IsSender:   true,
		TokenNonce: 10,
		Identifier: "token-112-0a",
	}, accountBalanceHistory)
}

func TestAccountsProcessor_PutTokenMedataDataInTokens(t *testing.T) {
	t.Parallel()

	t.Run("no tokens with missing data or nonce higher than 0", func(t *testing.T) {
		t.Parallel()

		ap, _ := NewAccountsProcessor(mock.NewPubkeyConverterMock(32), balanceConverter)

		oldCreator := "old creator"
		tokensInfo := []*data.TokenInfo{
			{Data: nil}, {Nonce: 5, Data: &data.TokenMetaData{Creator: oldCreator}},
		}
		emptyAlteredAccounts := map[string]*alteredAccount.AlteredAccount{}
		ap.PutTokenMedataDataInTokens(tokensInfo, emptyAlteredAccounts)
		require.Empty(t, emptyAlteredAccounts)
		require.Empty(t, tokensInfo[0].Data)
		require.Equal(t, oldCreator, tokensInfo[1].Data.Creator)
	})

	t.Run("error loading token, should not update metadata", func(t *testing.T) {
		t.Parallel()

		ap, _ := NewAccountsProcessor(mock.NewPubkeyConverterMock(32), balanceConverter)

		tokensInfo := []*data.TokenInfo{
			{
				Name:  "token0",
				Data:  nil,
				Nonce: 5,
			},
		}

		alteredAccounts := map[string]*alteredAccount.AlteredAccount{
			"addr": {Tokens: []*alteredAccount.AccountTokenData{}},
		}
		ap.PutTokenMedataDataInTokens(tokensInfo, alteredAccounts)
		require.Empty(t, tokensInfo[0].Data)
	})

	t.Run("should work and update metadata", func(t *testing.T) {
		t.Parallel()

		ap, _ := NewAccountsProcessor(mock.NewPubkeyConverterMock(32), balanceConverter)

		metadata0, metadata1 := &alteredAccount.TokenMetaData{Creator: "creator 0"}, &alteredAccount.TokenMetaData{Creator: "creator 1"}
		tokensInfo := []*data.TokenInfo{
			{
				Nonce:      5,
				Token:      "token0-5t6y7u",
				Identifier: "token0-5t6y7u-05",
			},
			{
				Nonce:      10,
				Token:      "token1-999ddd",
				Identifier: "token1-999ddd-0a",
			},
		}

		alteredAccounts := map[string]*alteredAccount.AlteredAccount{
			"addr0": {
				Tokens: []*alteredAccount.AccountTokenData{
					{
						Identifier: "token0-5t6y7u",
						Nonce:      5,
						MetaData:   metadata0,
					},
					{
						Identifier: "token1-999ddd",
						Nonce:      10,
						MetaData:   metadata1,
					},
				},
			},
		}

		ap.PutTokenMedataDataInTokens(tokensInfo, alteredAccounts)
		require.Equal(t, metadata0.Creator, tokensInfo[0].Data.Creator)
		require.Equal(t, metadata1.Creator, tokensInfo[1].Data.Creator)
	})
}

func TestAddAdditionalDataIntoAccounts(t *testing.T) {
	t.Parallel()

	ap, _ := NewAccountsProcessor(mock.NewPubkeyConverterMock(32), balanceConverter)

	account := &data.AccountInfo{}
	ap.addAdditionalDataInAccount(&alteredAccount.AdditionalAccountData{
		DeveloperRewards: "10000",
	}, account)
	require.Equal(t, "10000", account.DeveloperRewards)
	require.Equal(t, 0.000001, account.DeveloperRewardsNum)

	account = &data.AccountInfo{}
	ap.addAdditionalDataInAccount(&alteredAccount.AdditionalAccountData{
		DeveloperRewards: "",
	}, account)
	require.Equal(t, "", account.DeveloperRewards)
	require.Equal(t, float64(0), account.DeveloperRewardsNum)

	account = &data.AccountInfo{
		Address: "addr",
	}
	ap.addAdditionalDataInAccount(&alteredAccount.AdditionalAccountData{
		DeveloperRewards: "wrong",
	}, account)
	require.Equal(t, "", account.DeveloperRewards)
	require.Equal(t, float64(0), account.DeveloperRewardsNum)
}

func TestIsFrozen(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	properties := []byte{1}
	result := isFrozen(hex.EncodeToString(properties))
	require.True(result)

	result = isFrozen("invalid")
	require.False(result)

	emptyHex := ""
	result = isFrozen(emptyHex)
	require.False(result)

	properties = []byte{3}
	result = isFrozen(hex.EncodeToString(properties))
	require.True(result)

	properties = []byte{4}
	result = isFrozen(hex.EncodeToString(properties))
	require.False(result)
}
