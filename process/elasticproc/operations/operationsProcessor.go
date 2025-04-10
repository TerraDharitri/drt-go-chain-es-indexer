package operations

import (
	"encoding/hex"
	"strings"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/TerraDharitri/drt-go-chain-es-indexer/data"
	vmcommon "github.com/TerraDharitri/drt-go-chain-vm-common"
)

type operationsProcessor struct {
}

// NewOperationsProcessor will create a new instance of operationsProcessor
func NewOperationsProcessor() (*operationsProcessor, error) {
	return &operationsProcessor{}, nil
}

// ProcessTransactionsAndSCRs will prepare transactions and smart contract results to be indexed
func (op *operationsProcessor) ProcessTransactionsAndSCRs(
	txs []*data.Transaction,
	scrs []*data.ScResult,
	isImportDB bool,
	selfShardID uint32,
) ([]*data.Transaction, []*data.ScResult) {
	newTxsSlice := make([]*data.Transaction, 0)
	newScrsSlice := make([]*data.ScResult, 0)

	for idx, tx := range txs {
		if !op.shouldIndex(txs[idx].ReceiverShard, isImportDB, selfShardID) {
			continue
		}

		copiedTx := *tx
		copiedTx.SmartContractResults = nil
		copiedTx.Type = string(transaction.TxTypeNormal)
		newTxsSlice = append(newTxsSlice, &copiedTx)
	}

	for idx := 0; idx < len(scrs); idx++ {
		if !op.shouldIndex(scrs[idx].ReceiverShard, isImportDB, selfShardID) {
			continue
		}

		copiedScr := *scrs[idx]
		copiedScr.Type = string(transaction.TxTypeUnsigned)

		setCanBeIgnoredField(&copiedScr)
		if selfShardID == copiedScr.ReceiverShard {
			copiedScr.Status = transaction.TxStatusSuccess.String()
		} else {
			copiedScr.Status = transaction.TxStatusPending.String()
		}

		newScrsSlice = append(newScrsSlice, &copiedScr)
	}

	return newTxsSlice, newScrsSlice
}

func (op *operationsProcessor) shouldIndex(destinationShardID uint32, isImportDB bool, selfShardID uint32) bool {
	if !isImportDB {
		return true
	}

	return selfShardID == destinationShardID
}

func setCanBeIgnoredField(scr *data.ScResult) {
	dataFieldStr := string(scr.Data)
	hasOkPrefix := strings.HasPrefix(dataFieldStr, data.AtSeparator+hex.EncodeToString([]byte(vmcommon.Ok.String())))
	isRefundForRelayed := scr.ReturnMessage == data.GasRefundForRelayerMessage && dataFieldStr == ""
	if hasOkPrefix || isRefundForRelayed {
		scr.CanBeIgnored = true
		return
	}

	isNFTTransferOrMultiTransfer := core.BuiltInFunctionDCDTNFTTransfer == scr.Operation || core.BuiltInFunctionMultiDCDTNFTTransfer == scr.Operation
	isSCAddr := core.IsSmartContractAddress(scr.SenderAddressBytes)
	if isNFTTransferOrMultiTransfer && !isSCAddr {
		scr.CanBeIgnored = true
		return
	}
}
