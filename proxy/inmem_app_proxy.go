package proxy

import (
	"github.com/newdag/crypto"
	"github.com/newdag/ledger"
	"github.com/newdag/store"
)

//InmemProxy is used for testing
type InmemAppProxy struct {
	submitCh              chan ledger.Transaction
	stateHash             []byte
	committedTransactions []ledger.Transaction
	logger                *Logger
}

func NewInmemAppProxy(logger *Logger) *InmemAppProxy {
	if logger == nil {
		logger = store.Logger_new("")
	}
	return &InmemAppProxy{
		submitCh:              make(chan ledger.Transaction),
		stateHash:             []byte{},
		committedTransactions: []ledger.Transaction{},
		logger:                logger,
	}
}

func (iap *InmemAppProxy) commit(block ledger.Block) ([]byte, error) {

	//??iap.committedTransactions = append(iap.committedTransactions, block.Transactions...)

	hash := iap.stateHash
	for _, t := range block.Transactions {
		tHash := crypto.SHA256(t.Serialize())
		hash = crypto.SimpleHashFromTwoHashes(hash, tHash)
	}

	iap.stateHash = hash

	return iap.stateHash, nil

}

//------------------------------------------------------------------------------
//Implement AppProxy Interface

func (p *InmemAppProxy) SubmitCh() chan ledger.Transaction {
	return p.submitCh
}

func (p *InmemAppProxy) CommitBlock(block ledger.Block) (stateHash []byte, err error) {
	p.logger.Log("debug", "InmemProxy CommitBlock", LogFields{"round_received": block.RoundReceived, "txs": len(block.Transactions)})
	return p.commit(block)
}

//------------------------------------------------------------------------------

func (p *InmemAppProxy) SubmitTx(tx ledger.Transaction) {
	p.submitCh <- tx
}

func (p *InmemAppProxy) GetCommittedTransactions() []ledger.Transaction {
	return p.committedTransactions
}
