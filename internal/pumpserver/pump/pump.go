// package pump defines the core business logic.
package pump

import (
	"context"
	"sync/atomic"

	"github.com/che-kwas/iam-kit/logger"
)

// Pump defines the structure of a pump.
type Pump struct {
	shouldStop uint32
	ctx        context.Context
	log        *logger.Logger
}

var pumpIns *Pump

// InitPump initializes the global pump and returns it.
func InitPump(ctx context.Context, opts *PumpOptions) *Pump {
	log := logger.L()
	log.Debugf("building pump with options: %+v", opts)

	pumpIns = &Pump{
		ctx: ctx,
		log: log,
	}

	return pumpIns
}

// GetPump returns the global pump.
func GetPump() *Pump {
	return pumpIns
}

// Start starts the pump.
func (p *Pump) Start() {
	atomic.SwapUint32(&p.shouldStop, 0)
}

// Stop flushes the buffer and stop the pump.
func (p *Pump) Stop(ctx context.Context) error {
	// flag to stop sending records into channel
	atomic.SwapUint32(&p.shouldStop, 1)

	return nil
}
