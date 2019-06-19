package consensus

import (
	"context"

	"github.com/filecoin-project/go-filecoin/chain"
	"github.com/filecoin-project/go-filecoin/porcelain"
	"github.com/filecoin-project/go-filecoin/types"
)

// FaultMonitor checks each new tipset for faults
type FaultMonitor struct {
	porcelain porcelain.API
}

// HandleNewTipSet receives an iterator over the current chain, and a new tipset
// and checks the new tipset for fault errors, iteratoring over chnIter
func (fm *FaultMonitor) HandleNewTipSet(ctx context.Context, iter chain.TipsetIterator, newTs types.TipSet) {

}
