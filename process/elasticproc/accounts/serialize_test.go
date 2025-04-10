package accounts

import (
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-es-indexer/data"
	"github.com/TerraDharitri/drt-go-chain-es-indexer/process/elasticproc/converters"
	"github.com/stretchr/testify/require"
)

func TestSerializeNFTCreateInfo(t *testing.T) {
	t.Parallel()

	nftsCreateInfo := []*data.TokenInfo{
		{
			Token:      "my-token-0001",
			Identifier: "my-token-001-0f",
			Data: &data.TokenMetaData{
				Creator: "010102",
			},
			Type: core.NonFungibleDCDT,
		},
	}

	buffSlice := data.NewBufferSlice(data.DefaultMaxBulkSize)
	err := (&accountsProcessor{}).SerializeNFTCreateInfo(nftsCreateInfo, buffSlice, "tokens")
	require.NoError(t, err)
	require.Equal(t, 1, len(buffSlice.Buffers()))

	expectedRes := `{ "index" : { "_index":"tokens", "_id" : "my-token-001-0f" } }
{"identifier":"my-token-001-0f","token":"my-token-0001","numDecimals":0,"type":"NonFungibleDCDT","data":{"creator":"010102","nonEmptyURIs":false,"whiteListedStorage":false}}
`
	require.Equal(t, expectedRes, buffSlice.Buffers()[0].String())
}

func TestSerializeAccounts(t *testing.T) {
	t.Parallel()

	accs := map[string]*data.AccountInfo{
		"addr1": {
			Address:         "addr1",
			Nonce:           1,
			Balance:         "50",
			BalanceNum:      0.1,
			IsSmartContract: true,
			IsSender:        true,
		},
	}

	buffSlice := data.NewBufferSlice(data.DefaultMaxBulkSize)
	err := (&accountsProcessor{}).SerializeAccounts(accs, buffSlice, "accounts")
	require.NoError(t, err)
	require.Equal(t, 1, len(buffSlice.Buffers()))

	expectedRes := `{ "update" : {"_index": "accounts", "_id" : "addr1" } }
{"scripted_upsert": true, "script": {"source": "if ('create' == ctx.op) {ctx._source = params.account} else {if ((!ctx._source.containsKey('timestamp')) || (ctx._source.timestamp <= params.account.timestamp) ) {params.account.forEach((key, value) -> {ctx._source[key] = value;});}}","lang": "painless","params": { "account": {"address":"addr1","nonce":1,"balance":"50","balanceNum":0.1,"shardID":0} }},"upsert": {}}
`
	require.Equal(t, expectedRes, buffSlice.Buffers()[0].String())
}

func TestSerializeAccountsDCDTNonceZero(t *testing.T) {
	t.Parallel()

	accs := map[string]*data.AccountInfo{
		"addr1": {
			Address:    "addr1",
			Nonce:      1,
			TokenName:  "token-abcd",
			Properties: "000",
			TokenNonce: 0,
			Balance:    "10000000000000",
			BalanceNum: 1,
			Timestamp:  123,
		},
	}

	buffSlice := data.NewBufferSlice(data.DefaultMaxBulkSize)
	err := (&accountsProcessor{}).SerializeAccountsDCDT(accs, nil, buffSlice, "accountsdcdt")
	require.NoError(t, err)
	require.Equal(t, 1, len(buffSlice.Buffers()))

	expectedRes := `{ "update" : {"_index": "accountsdcdt", "_id" : "addr1-token-abcd-00" } }
{"scripted_upsert": true, "script": {"source": "if ('create' == ctx.op) {ctx._source = params.account} else {if ((!ctx._source.containsKey('timestamp')) || (ctx._source.timestamp <= params.account.timestamp) ) {params.account.forEach((key, value) -> {ctx._source[key] = value;});}}","lang": "painless","params": { "account": {"address":"addr1","nonce":1,"balance":"10000000000000","balanceNum":1,"token":"token-abcd","properties":"000","timestamp":123,"shardID":0} }},"upsert": {}}
`
	require.Equal(t, expectedRes, buffSlice.Buffers()[0].String())
}

func TestSerializeAccountsDCDT(t *testing.T) {
	t.Parallel()

	accs := map[string]*data.AccountInfo{
		"addr1": {
			Address:    "addr1",
			Nonce:      1,
			TokenName:  "token-0001",
			Properties: "000",
			TokenNonce: 5,
			Balance:    "10000000000000",
			BalanceNum: 1,
		},
	}

	buffSlice := data.NewBufferSlice(data.DefaultMaxBulkSize)
	err := (&accountsProcessor{}).SerializeAccountsDCDT(accs, nil, buffSlice, "accountsdcdt")
	require.NoError(t, err)
	require.Equal(t, 1, len(buffSlice.Buffers()))

	expectedRes := `{ "update" : {"_index": "accountsdcdt", "_id" : "addr1-token-0001-05" } }
{"scripted_upsert": true, "script": {"source": "if ('create' == ctx.op) {ctx._source = params.account} else {if ((!ctx._source.containsKey('timestamp')) || (ctx._source.timestamp <= params.account.timestamp) ) {params.account.forEach((key, value) -> {ctx._source[key] = value;});}}","lang": "painless","params": { "account": {"address":"addr1","nonce":1,"balance":"10000000000000","balanceNum":1,"token":"token-0001","tokenNonce":5,"properties":"000","shardID":0} }},"upsert": {}}
`
	require.Equal(t, expectedRes, buffSlice.Buffers()[0].String())
}

func TestSerializeAccountsNFTWithMedaData(t *testing.T) {
	t.Parallel()

	accs := map[string]*data.AccountInfo{
		"addr1": {
			Address:         "addr1",
			Nonce:           1,
			TokenName:       "token-0001",
			Properties:      "000",
			TokenNonce:      22,
			Balance:         "10000000000000",
			BalanceNum:      1,
			TokenIdentifier: "token-0001-5",
			Data: &data.TokenMetaData{
				Name:      "nft",
				Creator:   "010101",
				Royalties: 1,
				Hash:      []byte("hash"),
				URIs: [][]byte{
					[]byte("uri"),
				},
				Attributes:   []byte("tags:test,free,fun;description:This is a test description for an awesome nft"),
				Tags:         converters.ExtractTagsFromAttributes([]byte("tags:test,free,fun;description:This is a test description for an awesome nft")),
				MetaData:     "metadata-test",
				NonEmptyURIs: true,
			},
		},
	}

	buffSlice := data.NewBufferSlice(data.DefaultMaxBulkSize)
	err := (&accountsProcessor{}).SerializeAccountsDCDT(accs, nil, buffSlice, "accountsdcdt")
	require.NoError(t, err)
	require.Equal(t, 1, len(buffSlice.Buffers()))

	expectedRes := `{ "update" : {"_index": "accountsdcdt", "_id" : "addr1-token-0001-16" } }
{"scripted_upsert": true, "script": {"source": "if ('create' == ctx.op) {ctx._source = params.account} else {if ((!ctx._source.containsKey('timestamp')) || (ctx._source.timestamp <= params.account.timestamp) ) {params.account.forEach((key, value) -> {ctx._source[key] = value;});}}","lang": "painless","params": { "account": {"address":"addr1","nonce":1,"balance":"10000000000000","balanceNum":1,"token":"token-0001","identifier":"token-0001-5","tokenNonce":22,"properties":"000","data":{"name":"nft","creator":"010101","royalties":1,"hash":"aGFzaA==","uris":["dXJp"],"tags":["test","free","fun"],"attributes":"dGFnczp0ZXN0LGZyZWUsZnVuO2Rlc2NyaXB0aW9uOlRoaXMgaXMgYSB0ZXN0IGRlc2NyaXB0aW9uIGZvciBhbiBhd2Vzb21lIG5mdA==","metadata":"metadata-test","nonEmptyURIs":true,"whiteListedStorage":false},"shardID":0} }},"upsert": {}}
`
	require.Equal(t, expectedRes, buffSlice.Buffers()[0].String())
}

func TestSerializeAccountsDCDTDelete(t *testing.T) {
	t.Parallel()

	accs := map[string]*data.AccountInfo{
		"addr1": {
			Address:    "addr1",
			Nonce:      1,
			TokenName:  "token-0001",
			Properties: "000",
			Balance:    "0",
			BalanceNum: 0,
		},
	}

	buffSlice := data.NewBufferSlice(data.DefaultMaxBulkSize)
	err := (&accountsProcessor{}).SerializeAccountsDCDT(accs, nil, buffSlice, "accountsdcdt")
	require.NoError(t, err)
	require.Equal(t, 1, len(buffSlice.Buffers()))

	expectedRes := `{ "update" : {"_index":"accountsdcdt", "_id" : "addr1-token-0001-00" } }
{"scripted_upsert": true, "script": {"source": "if ('create' == ctx.op) {ctx.op = 'noop'} else {if (ctx._source.containsKey('timestamp')) {if (ctx._source.timestamp <= params.timestamp) {ctx.op = 'delete'}} else {ctx.op = 'delete'}}","lang": "painless","params": {"timestamp": 0}},"upsert": {}}
`
	require.Equal(t, expectedRes, buffSlice.Buffers()[0].String())
}

func TestSerializeAccountsHistory(t *testing.T) {
	t.Parallel()

	accsh := map[string]*data.AccountBalanceHistory{
		"account1": {
			Address:         "account1",
			Timestamp:       10,
			Balance:         "123",
			Token:           "token-0001",
			IsSender:        true,
			IsSmartContract: true,
		},
	}

	buffSlice := data.NewBufferSlice(data.DefaultMaxBulkSize)
	err := (&accountsProcessor{}).SerializeAccountsHistory(accsh, buffSlice, "accountshistory")
	require.NoError(t, err)
	require.Equal(t, 1, len(buffSlice.Buffers()))

	expectedRes := `{ "index" : { "_index":"accountshistory", "_id" : "account1-token-0001-00-10" } }
{"address":"account1","timestamp":10,"balance":"123","token":"token-0001","isSender":true,"isSmartContract":true,"shardID":0}
`
	require.Equal(t, expectedRes, buffSlice.Buffers()[0].String())
}
