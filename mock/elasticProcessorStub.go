package mock

import (
	coreData "github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain-core/data/outport"
)

// ElasticProcessorStub -
type ElasticProcessorStub struct {
	SaveHeaderCalled                 func(outportBlockWithHeader *outport.OutportBlockWithHeader) error
	RemoveHeaderCalled               func(header coreData.HeaderHandler) error
	RemoveMiniblocksCalled           func(header coreData.HeaderHandler, body *block.Body) error
	RemoveTransactionsCalled         func(header coreData.HeaderHandler, body *block.Body) error
	SaveMiniblocksCalled             func(header coreData.HeaderHandler, miniBlocks []*block.MiniBlock) error
	SaveTransactionsCalled           func(outportBlockWithHeader *outport.OutportBlockWithHeader) error
	SaveValidatorsRatingCalled       func(validatorsRating *outport.ValidatorsRating) error
	SaveRoundsInfoCalled             func(infos *outport.RoundsInfo) error
	SaveShardValidatorsPubKeysCalled func(validators *outport.ValidatorsPubKeys) error
	SaveAccountsCalled               func(accountsData *outport.Accounts) error
	RemoveAccountsDCDTCalled         func(headerTimestamp uint64) error
}

// RemoveAccountsDCDT -
func (eim *ElasticProcessorStub) RemoveAccountsDCDT(headerTimestamp uint64, _ uint32) error {
	if eim.RemoveAccountsDCDTCalled != nil {
		return eim.RemoveAccountsDCDTCalled(headerTimestamp)
	}

	return nil
}

// SaveHeader -
func (eim *ElasticProcessorStub) SaveHeader(obh *outport.OutportBlockWithHeader) error {
	if eim.SaveHeaderCalled != nil {
		return eim.SaveHeaderCalled(obh)
	}
	return nil
}

// RemoveHeader -
func (eim *ElasticProcessorStub) RemoveHeader(header coreData.HeaderHandler) error {
	if eim.RemoveHeaderCalled != nil {
		return eim.RemoveHeaderCalled(header)
	}
	return nil
}

// RemoveMiniblocks -
func (eim *ElasticProcessorStub) RemoveMiniblocks(header coreData.HeaderHandler, body *block.Body) error {
	if eim.RemoveMiniblocksCalled != nil {
		return eim.RemoveMiniblocksCalled(header, body)
	}
	return nil
}

// RemoveTransactions -
func (eim *ElasticProcessorStub) RemoveTransactions(header coreData.HeaderHandler, body *block.Body) error {
	if eim.RemoveMiniblocksCalled != nil {
		return eim.RemoveTransactionsCalled(header, body)
	}
	return nil
}

// SaveMiniblocks -
func (eim *ElasticProcessorStub) SaveMiniblocks(header coreData.HeaderHandler, miniBlocks []*block.MiniBlock) error {
	if eim.SaveMiniblocksCalled != nil {
		return eim.SaveMiniblocksCalled(header, miniBlocks)
	}
	return nil
}

// SaveTransactions -
func (eim *ElasticProcessorStub) SaveTransactions(outportBlockWithHeader *outport.OutportBlockWithHeader) error {
	if eim.SaveTransactionsCalled != nil {
		return eim.SaveTransactionsCalled(outportBlockWithHeader)
	}
	return nil
}

// SaveValidatorsRating -
func (eim *ElasticProcessorStub) SaveValidatorsRating(validatorsRating *outport.ValidatorsRating) error {
	if eim.SaveValidatorsRatingCalled != nil {
		return eim.SaveValidatorsRatingCalled(validatorsRating)
	}
	return nil
}

// SaveRoundsInfo -
func (eim *ElasticProcessorStub) SaveRoundsInfo(info *outport.RoundsInfo) error {
	if eim.SaveRoundsInfoCalled != nil {
		return eim.SaveRoundsInfoCalled(info)
	}
	return nil
}

// SaveShardValidatorsPubKeys -
func (eim *ElasticProcessorStub) SaveShardValidatorsPubKeys(validatorsPubKeys *outport.ValidatorsPubKeys) error {
	if eim.SaveShardValidatorsPubKeysCalled != nil {
		return eim.SaveShardValidatorsPubKeysCalled(validatorsPubKeys)
	}
	return nil
}

// SaveAccounts -
func (eim *ElasticProcessorStub) SaveAccounts(accounts *outport.Accounts) error {
	if eim.SaveAccountsCalled != nil {
		return eim.SaveAccountsCalled(accounts)
	}

	return nil
}

// SetOutportConfig -
func (eim *ElasticProcessorStub) SetOutportConfig(_ outport.OutportConfig) error {
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (eim *ElasticProcessorStub) IsInterfaceNil() bool {
	return eim == nil
}
