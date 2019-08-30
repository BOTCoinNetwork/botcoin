package babble

import (
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mosaicnetworks/babble/src/babble"
	"github.com/mosaicnetworks/babble/src/hashgraph"
	"github.com/mosaicnetworks/babble/src/proxy"
	"github.com/mosaicnetworks/evm-lite/src/service"
	"github.com/mosaicnetworks/evm-lite/src/state"
	"github.com/sirupsen/logrus"
)

// InmemProxy implements the Babble AppProxy interface
type InmemProxy struct {
	service  *service.Service
	state    *state.State
	babble   *babble.Babble
	submitCh chan []byte
	logger   *logrus.Entry
}

// NewInmemProxy initializes and return a new InmemProxy
func NewInmemProxy(state *state.State,
	service *service.Service,
	babble *babble.Babble,
	submitCh chan []byte,
	logger *logrus.Entry) *InmemProxy {

	return &InmemProxy{
		service:  service,
		state:    state,
		babble:   babble,
		submitCh: submitCh,
		logger:   logger,
	}
}

/*******************************************************************************
Implement Babble AppProxy Interface
*******************************************************************************/

// SubmitCh is the channel through which the Service sends transactions to the
// node.
func (p *InmemProxy) SubmitCh() chan []byte {
	return p.submitCh
}

// CommitBlock applies the block's transactions to the state and commits. It
// also checks the block's internal transactions against the POA smart-contract
// to check if joining peers are authorised to become validators in Babble. It
// returns the resulting state-hash and internal transaction receips.
func (p *InmemProxy) CommitBlock(block hashgraph.Block) (proxy.CommitResponse, error) {

	// XXX get coinbase

	coinbaseAddress := ethCommon.Address{}

	if p.babble != nil {
		babbleValidators, err := p.babble.Node.GetValidators(block.RoundReceived())
		if err != nil {
			return proxy.CommitResponse{}, err
		}

		coinbaseValidator := babbleValidators[block.Index()%len(babbleValidators)]

		coinbasePubKey, err := crypto.UnmarshalPubkey(coinbaseValidator.PubKeyBytes())
		if err != nil {
			p.logger.Warningf("couldn't unmarshal pubkey bytes: %v", err)
		}

		coinbaseAddress = crypto.PubkeyToAddress(*coinbasePubKey)
	}

	p.logger.WithFields(logrus.Fields{
		"coinbase": coinbaseAddress.String(),
		"block":    block.Index(),
	}).Info("Commit")

	// END XXX

	blockHashBytes, err := block.Hash()
	blockHash := ethCommon.BytesToHash(blockHashBytes)

	for i, tx := range block.Transactions() {
		if err := p.state.ApplyTransaction(tx, i, blockHash, coinbaseAddress); err != nil {
			return proxy.CommitResponse{}, err
		}
	}

	hash, err := p.state.Commit()
	if err != nil {
		return proxy.CommitResponse{}, err
	}

	receipts := p.processInternalTransactions(block.InternalTransactions())

	res := proxy.CommitResponse{
		StateHash:                   hash.Bytes(),
		InternalTransactionReceipts: receipts,
	}

	return res, nil
}

// processInternalTransactions decides if InternalTransactions should be
// accepted. For PEER_ADD transactions, it checks the if the peer is authorised
// in the POA smart-contract. All PEER_REMOVE transactions are accepted
func (p *InmemProxy) processInternalTransactions(internalTransactions []hashgraph.InternalTransaction) []hashgraph.InternalTransactionReceipt {
	receipts := []hashgraph.InternalTransactionReceipt{}

	for _, tx := range internalTransactions {
		switch tx.Body.Type {
		case hashgraph.PEER_ADD:
			pk, err := crypto.UnmarshalPubkey(tx.Body.Peer.PubKeyBytes())
			if err != nil {
				p.logger.Warningf("couldn't unmarshal pubkey bytes: %v", err)
			}

			addr := crypto.PubkeyToAddress(*pk)

			ok, err := p.state.CheckAuthorised(addr)

			if err != nil {
				p.logger.WithError(err).Error("Error in checkAuthorised")
				receipts = append(receipts, tx.AsRefused())
			} else {
				if ok {
					p.logger.WithField("addr", addr.String()).Info("Accepted peer")
					receipts = append(receipts, tx.AsAccepted())
				} else {
					p.logger.WithField("addr", addr.String()).Info("Rejected peer")
					receipts = append(receipts, tx.AsRefused())
				}
			}
		case hashgraph.PEER_REMOVE:
			receipts = append(receipts, tx.AsAccepted())
		}
	}

	return receipts
}

//TODO - Implement these two functions
func (p *InmemProxy) GetSnapshot(blockIndex int) ([]byte, error) {
	return []byte{}, nil
}

func (p *InmemProxy) Restore(snapshot []byte) error {
	return nil
}
