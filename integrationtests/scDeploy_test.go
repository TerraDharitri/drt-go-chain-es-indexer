//go:build integrationtests

package integrationtests

import (
	"context"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	dataBlock "github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain-core/data/outport"
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	indexerData "github.com/TerraDharitri/drt-go-chain-es-indexer/process/dataindexer"
	"github.com/stretchr/testify/require"
)

func TestTransactionWithSCDeploy(t *testing.T) {
	setLogLevelDebug()

	esClient, err := createESClient(esURL)
	require.Nil(t, err)

	esProc, err := CreateElasticProcessor(esClient)
	require.Nil(t, err)

	txHash := []byte("scDeployHash")
	header := &dataBlock.Header{
		Round:     50,
		TimeStamp: 5040,
		ShardID:   2,
	}
	body := &dataBlock.Body{
		MiniBlocks: dataBlock.MiniBlockSlice{
			{
				Type:            dataBlock.TxBlock,
				SenderShardID:   2,
				ReceiverShardID: 2,
				TxHashes:        [][]byte{txHash},
			},
		},
	}
	sndAddress := "drt12m3x8jp6dl027pj5f2nw6ght2cyhhjfrs86cdwsa8xn83r375qfq7jkw93"
	tx := &transaction.Transaction{
		Nonce:    1,
		SndAddr:  decodeAddress(sndAddress),
		RcvAddr:  decodeAddress("drt1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq85hk5z"),
		GasLimit: 1000000000,
		GasPrice: 2000000,
		Data:     []byte("0061736d01000000010d036000006000017f60027f7f00023e0303656e760f6765744e756d417267756d656e7473000103656e760b7369676e616c4572726f72000203656e760e636865636b4e6f5061796d656e74000003030200000503010003060f027f00419980080b7f0041a080080b073705066d656d6f7279020004696e697400030863616c6c4261636b00040a5f5f646174615f656e6403000b5f5f686561705f6261736503010a180212001002100004404180800841191001000b0b0300010b0b210100418080080b1977726f6e67206e756d626572206f6620617267756d656e7473@0500@0502"),
		Value:    big.NewInt(0),
	}

	txInfo := &outport.TxInfo{
		Transaction: tx,
		FeeInfo: &outport.FeeInfo{
			GasUsed:        1130820,
			Fee:            big.NewInt(764698200000000),
			InitialPaidFee: big.NewInt(773390000000000),
		},
		ExecutionOrder: 0,
	}

	pool := &outport.TransactionPool{
		Transactions: map[string]*outport.TxInfo{
			hex.EncodeToString(txHash): txInfo,
		},
		Logs: []*outport.LogData{
			{
				TxHash: hex.EncodeToString([]byte("h1")),
				Log: &transaction.Log{
					Address: decodeAddress(sndAddress),
					Events: []*transaction.Event{
						{
							Address:    decodeAddress(sndAddress),
							Identifier: []byte(core.SCDeployIdentifier),
							Topics:     [][]byte{decodeAddress("drt1qqqqqqqqqqqqqpgq4t2tqxpst9a6qttpak8cz8wvz6a0nses5qfqyrdq56"), decodeAddress("drt12m3x8jp6dl027pj5f2nw6ght2cyhhjfrs86cdwsa8xn83r375qfq7jkw93"), []byte("codeHash")},
						},
						nil,
					},
				},
			},
		},
	}
	err = esProc.SaveTransactions(createOutportBlockWithHeader(body, header, pool, nil, testNumOfShards))
	require.Nil(t, err)

	ids := []string{hex.EncodeToString(txHash)}
	genericResponse := &GenericResponse{}
	err = esClient.DoMultiGet(context.Background(), ids, indexerData.TransactionsIndex, true, genericResponse)
	require.Nil(t, err)

	require.JSONEq(t,
		readExpectedResult("./testdata/scDeploy/tx-sc-deploy.json"),
		string(genericResponse.Docs[0].Source),
	)

	ids = []string{"drt1qqqqqqqqqqqqqpgq4t2tqxpst9a6qttpak8cz8wvz6a0nses5qfqyrdq56"}
	err = esClient.DoMultiGet(context.Background(), ids, indexerData.SCDeploysIndex, true, genericResponse)
	require.Nil(t, err)

	require.JSONEq(t,
		readExpectedResult("./testdata/scDeploy/deploy.json"),
		string(genericResponse.Docs[0].Source),
	)

	// UPGRADE contract
	header.TimeStamp = 6000
	pool = &outport.TransactionPool{
		Logs: []*outport.LogData{
			{
				TxHash: hex.EncodeToString([]byte("h2")),
				Log: &transaction.Log{
					Address: decodeAddress(sndAddress),
					Events: []*transaction.Event{
						{
							Address:    decodeAddress(sndAddress),
							Identifier: []byte(core.SCUpgradeIdentifier),
							Topics:     [][]byte{decodeAddress("drt1qqqqqqqqqqqqqpgq4t2tqxpst9a6qttpak8cz8wvz6a0nses5qfqyrdq56"), decodeAddress("drt12m3x8jp6dl027pj5f2nw6ght2cyhhjfrs86cdwsa8xn83r375qfq7jkw93"), []byte("secondCodeHash")},
						},
						nil,
					},
				},
			},
		},
	}
	err = esProc.SaveTransactions(createOutportBlockWithHeader(body, header, pool, nil, testNumOfShards))
	require.Nil(t, err)

	ids = []string{"drt1qqqqqqqqqqqqqpgq4t2tqxpst9a6qttpak8cz8wvz6a0nses5qfqyrdq56"}
	err = esClient.DoMultiGet(context.Background(), ids, indexerData.SCDeploysIndex, true, genericResponse)
	require.Nil(t, err)

	require.JSONEq(t,
		readExpectedResult("./testdata/scDeploy/deploy-after-upgrade.json"),
		string(genericResponse.Docs[0].Source),
	)

	// CHANGE owner first
	header.TimeStamp = 7000
	pool = &outport.TransactionPool{
		Logs: []*outport.LogData{
			{
				TxHash: hex.EncodeToString([]byte("h3")),
				Log: &transaction.Log{
					Address: decodeAddress(sndAddress),
					Events: []*transaction.Event{
						{
							Address:    decodeAddress("drt1qqqqqqqqqqqqqpgq4t2tqxpst9a6qttpak8cz8wvz6a0nses5qfqyrdq56"),
							Identifier: []byte(core.BuiltInFunctionChangeOwnerAddress),
							Topics:     [][]byte{decodeAddress("drt1d942l8w4yvgjffpqacs8vdwl0mndsv0zn0uxa80hxc3xmq4477eqwcejwf")},
						},
						nil,
					},
				},
			},
		},
	}
	err = esProc.SaveTransactions(createOutportBlockWithHeader(body, header, pool, nil, testNumOfShards))
	require.Nil(t, err)

	ids = []string{"drt1qqqqqqqqqqqqqpgq4t2tqxpst9a6qttpak8cz8wvz6a0nses5qfqyrdq56"}
	err = esClient.DoMultiGet(context.Background(), ids, indexerData.SCDeploysIndex, true, genericResponse)
	require.Nil(t, err)

	require.JSONEq(t,
		readExpectedResult("./testdata/scDeploy/deploy-after-upgrade-and-change-owner.json"),
		string(genericResponse.Docs[0].Source),
	)

	// CHANGE owner second
	header.TimeStamp = 8000
	pool = &outport.TransactionPool{
		Logs: []*outport.LogData{
			{
				TxHash: hex.EncodeToString([]byte("h4")),
				Log: &transaction.Log{
					Address: decodeAddress(sndAddress),
					Events: []*transaction.Event{
						{
							Address:    decodeAddress("drt1qqqqqqqqqqqqqpgq4t2tqxpst9a6qttpak8cz8wvz6a0nses5qfqyrdq56"),
							Identifier: []byte(core.BuiltInFunctionChangeOwnerAddress),
							Topics:     [][]byte{decodeAddress("drt1y78ds2tvzw6ntcggldjld2vk96wgq0mj47mk6auny0nkvn242e3ssfhpa9")},
						},
						nil,
					},
				},
			},
		},
	}
	err = esProc.SaveTransactions(createOutportBlockWithHeader(body, header, pool, nil, testNumOfShards))
	require.Nil(t, err)

	ids = []string{"drt1qqqqqqqqqqqqqpgq4t2tqxpst9a6qttpak8cz8wvz6a0nses5qfqyrdq56"}
	err = esClient.DoMultiGet(context.Background(), ids, indexerData.SCDeploysIndex, true, genericResponse)
	require.Nil(t, err)

	require.JSONEq(t,
		readExpectedResult("./testdata/scDeploy/deploy-after-upgrade-and-change-owner-second.json"),
		string(genericResponse.Docs[0].Source),
	)
}

func TestScDeployWithSignalErrorAndCompleteTxEvent(t *testing.T) {
	setLogLevelDebug()

	esClient, err := createESClient(esURL)
	require.Nil(t, err)

	esProc, err := CreateElasticProcessor(esClient)
	require.Nil(t, err)

	txHash := []byte("scDeployHashWithErrorAndCompleteTxEvent")
	header := &dataBlock.Header{
		Round:     50,
		TimeStamp: 5040,
		ShardID:   2,
	}
	body := &dataBlock.Body{
		MiniBlocks: dataBlock.MiniBlockSlice{
			{
				Type:            dataBlock.TxBlock,
				SenderShardID:   2,
				ReceiverShardID: 2,
				TxHashes:        [][]byte{txHash},
			},
		},
	}
	sndAddress := "drt12m3x8jp6dl027pj5f2nw6ght2cyhhjfrs86cdwsa8xn83r375qfq7jkw93"
	tx := &transaction.Transaction{
		Nonce:    1,
		SndAddr:  decodeAddress(sndAddress),
		RcvAddr:  decodeAddress("drt1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq85hk5z"),
		GasLimit: 1000000000,
		GasPrice: 2000000,
		Data:     []byte("0061736d01000000010d036000006000017f60027f7f00023e0303656e760f6765744e756d417267756d656e7473000103656e760b7369676e616c4572726f72000203656e760e636865636b4e6f5061796d656e74000003030200000503010003060f027f00419980080b7f0041a080080b073705066d656d6f7279020004696e697400030863616c6c4261636b00040a5f5f646174615f656e6403000b5f5f686561705f6261736503010a180212001002100004404180800841191001000b0b0300010b0b210100418080080b1977726f6e67206e756d626572206f6620617267756d656e7473@0500@0502"),
		Value:    big.NewInt(0),
	}

	txInfo := &outport.TxInfo{
		Transaction: tx,
		FeeInfo: &outport.FeeInfo{
			GasUsed:        1130820,
			Fee:            big.NewInt(764698200000000),
			InitialPaidFee: big.NewInt(773390000000000),
		},
		ExecutionOrder: 0,
	}

	pool := &outport.TransactionPool{
		Transactions: map[string]*outport.TxInfo{
			hex.EncodeToString(txHash): txInfo,
		},
		Logs: []*outport.LogData{
			{
				TxHash: hex.EncodeToString(txHash),
				Log: &transaction.Log{
					Address: decodeAddress(sndAddress),
					Events: []*transaction.Event{
						{
							Address:    decodeAddress(sndAddress),
							Identifier: []byte(core.SignalErrorOperation),
							Topics:     [][]byte{[]byte("h1"), []byte("h1"), []byte("h1")},
						},
						{
							Address:    decodeAddress(sndAddress),
							Identifier: []byte(core.CompletedTxEventIdentifier),
							Topics:     [][]byte{[]byte("h2"), []byte("h2"), []byte("h2")},
						},
					},
				},
			},
		},
	}
	err = esProc.SaveTransactions(createOutportBlockWithHeader(body, header, pool, nil, testNumOfShards))
	require.Nil(t, err)

	genericResponse := &GenericResponse{}
	ids := []string{hex.EncodeToString(txHash)}
	err = esClient.DoMultiGet(context.Background(), ids, indexerData.OperationsIndex, true, genericResponse)
	require.Nil(t, err)

	require.JSONEq(t,
		readExpectedResult("./testdata/scDeploy/tx-deploy-with-error-event.json"),
		string(genericResponse.Docs[0].Source),
	)
}
