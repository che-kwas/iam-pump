// package pump defines the core business logic.
package pump

import (
	"context"
	"time"

	"github.com/che-kwas/iam-kit/logger"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/vmihailenco/msgpack/v5"

	rdb "iam-pump/internal/pkg/redis"
	"iam-pump/internal/pumpserver/store"
)

const (
	queueName = "iam-authz-audit"

	redlockName   = "iam-pump"
	redlockExpiry = time.Duration(time.Minute)
)

// Pump defines the structure of a pump.
type Pump struct {
	interval time.Duration
	mutex    *redsync.Mutex
	ctx      context.Context
	log      *logger.Logger
}

var pumpIns *Pump

// InitPump initializes the global pump and returns it.
func InitPump(ctx context.Context, opts *PumpOptions) *Pump {
	log := logger.L()
	log.Debugf("building pump with options: %+v", opts)

	rs := redsync.New(goredis.NewPool(rdb.Client()))
	mutex := rs.NewMutex(redlockName, redsync.WithExpiry(redlockExpiry))

	pumpIns = &Pump{
		interval: opts.PumpInterval,
		mutex:    mutex,
		ctx:      ctx,
		log:      log,
	}

	return pumpIns
}

// GetPump returns the global pump.
func GetPump() *Pump {
	return pumpIns
}

// Start starts the pump.
func (p *Pump) Start() {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.pump()
		case <-p.ctx.Done():
			p.log.Info("the pump has got canceled.")
		}
	}
}

// pump transfers AuditRecord from redis to mongo.
func (p *Pump) pump() {
	if err := p.mutex.Lock(); err != nil {
		p.log.Debugf("there is already a worker pumping: %s", err.Error())
		return
	}

	defer func() {
		if _, err := p.mutex.Unlock(); err != nil {
			p.log.Errorw("could not release iam-pump lock", "error", err.Error())
		}
	}()

	records := p.purgeQueue()
	if err := p.writeToStore(records); err != nil {
		p.log.Errorw("failed to persist records", "error", err.Error())
	}
}

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

func (p *Pump) purgeQueue() []interface{} {
	var lrange *redis.StringSliceCmd
	_, err := rdb.Client().TxPipelined(p.ctx, func(pipe redis.Pipeliner) error {
		lrange = pipe.LRange(p.ctx, queueName, 0, -1)
		pipe.Del(p.ctx, queueName)

		return nil
	})
	if err != nil {
		p.log.Errorw("pop all from queue error", "error", err.Error())
		return nil
	}

	vals := lrange.Val()

	result := make([]interface{}, len(vals))
	for i, v := range vals {
		var ar AuditRecord
		if err := msgpack.Unmarshal([]byte(v), &ar); err != nil {
			p.log.Warnw("invalid audit record", "auditRecord", v)
			continue
		}
		result[i] = ar
	}

	return result
}

func (p *Pump) writeToStore(records []interface{}) error {
	if len(records) == 0 {
		return nil
	}

	p.log.Debugf("pump len = %d", len(records))
	return store.Client().InsertMany(p.ctx, records)
}
