package logsevents

import (
	"github.com/TerraDharitri/drt-go-chain-es-indexer/data"
	"github.com/TerraDharitri/drt-go-chain-es-indexer/process/elasticproc/converters"
	"github.com/TerraDharitri/drt-go-chain-es-indexer/process/elasticproc/tokeninfo"
)

type logsData struct {
	timestamp               uint64
	txHashStatusInfoProc    txHashStatusInfoHandler
	tokens                  data.TokensHandler
	tokensSupply            data.TokensHandler
	txsMap                  map[string]*data.Transaction
	scrsMap                 map[string]*data.ScResult
	scDeploys               map[string]*data.ScDeployInfo
	changeOwnerOperations   map[string]*data.OwnerData
	delegators              map[string]*data.Delegator
	tokensInfo              []*data.TokenInfo
	nftsDataUpdates         []*data.NFTDataUpdate
	tokenRolesAndProperties *tokeninfo.TokenRolesAndProperties
}

func newLogsData(
	timestamp uint64,
	txs []*data.Transaction,
	scrs []*data.ScResult,
) *logsData {
	ld := &logsData{}

	ld.txsMap = converters.ConvertTxsSliceIntoMap(txs)
	ld.scrsMap = converters.ConvertScrsSliceIntoMap(scrs)
	ld.tokens = data.NewTokensInfo()
	ld.tokensSupply = data.NewTokensInfo()
	ld.timestamp = timestamp
	ld.scDeploys = make(map[string]*data.ScDeployInfo)
	ld.tokensInfo = make([]*data.TokenInfo, 0)
	ld.delegators = make(map[string]*data.Delegator)
	ld.changeOwnerOperations = make(map[string]*data.OwnerData)
	ld.nftsDataUpdates = make([]*data.NFTDataUpdate, 0)
	ld.tokenRolesAndProperties = tokeninfo.NewTokenRolesAndProperties()
	ld.txHashStatusInfoProc = newTxHashStatusInfoProcessor()

	return ld
}
