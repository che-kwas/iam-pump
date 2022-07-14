// package pump defines the core business logic.
package pump

import (
	"context"

	"github.com/che-kwas/iam-kit/logger"
	"github.com/vmihailenco/msgpack"

	"iam-pump/internal/pumpserver/consumer"
	"iam-pump/internal/pumpserver/store"
)

func TransferAuditRecord(ctx context.Context, record []byte) {
	var ar AuditRecord
	if err := msgpack.Unmarshal(record, &ar); err != nil {
		return
	}
	logger.L().Debugw("pump", "AuditRecord", ar)

	store.Client().InsertOne(ctx, ar)
}

// Pump defines the structure of a pump.
type Pump struct {
	ctx      context.Context
	consumer consumer.Consumer
}

var pumpIns *Pump

// InitPump initializes the global pump and returns it.
func InitPump(ctx context.Context, consumer consumer.Consumer) *Pump {
	pumpIns = &Pump{ctx: ctx, consumer: consumer}

	return pumpIns
}

// GetPump returns the global pump.
func GetPump() *Pump {
	return pumpIns
}

// Start starts the pump.
func (p *Pump) Start() { p.consumer.Start(p.ctx) }

// AuditRecord defines the details of a authorization request.
type AuditRecord struct {
	Timestamp  int64
	Username   string
	Effect     string
	Conclusion string
	Request    string
	Policies   string
	Deciders   string
}
