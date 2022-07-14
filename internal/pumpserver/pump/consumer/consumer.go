// Package consumer is the kafka consumer builder.
package consumer

import (
	"context"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/che-kwas/iam-kit/logger"
)

// MsgHandler is the callback for handling the received message.
type MsgHandler func(topic string, key, value []byte)

// Consumer is the kafka consumer group handler.
type Consumer struct {
	ready      chan bool
	msgHandler MsgHandler
	ctx        context.Context
	group      sarama.ConsumerGroup
	topic      string
	poolSize   int
	log        *logger.Logger
}

// NewConsumer returns a kafka consumer.
func NewConsumer(ctx context.Context, msgHandler MsgHandler, opts *KafkaOptions) (*Consumer, error) {
	log := logger.L()
	log.Debugf("building kafka consumer with options: %+v", opts)

	group, err := newConsumerGroup(opts)
	if err != nil {
		return nil, err
	}

	consumer := &Consumer{
		ready:      make(chan bool),
		msgHandler: msgHandler,
		ctx:        ctx,
		group:      group,
		topic:      opts.Topic,
		poolSize:   opts.PoolSize,
		log:        log,
	}

	return consumer, nil
}

// Start starts the consume loop.
func (c *Consumer) Start() {
	wg := &sync.WaitGroup{}
	wg.Add(c.poolSize)
	for i := 0; i < c.poolSize; i++ {
		go func() {
			defer wg.Done()
			for {
				// `Consume` should be called inside an infinite loop, when a
				// server-side rebalance happens, the consumer session will need to be
				// recreated to get the new claims
				if err := c.group.Consume(c.ctx, []string{c.topic}, c); err != nil {
					c.log.Errorw("kafka consumer", "error", err)
				}
				// check if context was cancelled, signaling that the consumer should stop
				if c.ctx.Err() != nil {
					return
				}

				c.ready = make(chan bool)
			}
		}()
	}

	// Await until the consumer has been setup
	<-c.ready
	c.log.Info("kafka consumers are up and running!")

	wg.Wait()
}

// Close closes the consume group.
func (c *Consumer) Close(ctx context.Context) error {
	return c.group.Close()
}

// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L925-L943
var _ sarama.ConsumerGroupHandler = &Consumer{}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	logger.L().Debug("setup kafka consumer")
	// Mark the consumer as ready
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	logger.L().Debug("cleanup kafka consumer")
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message := <-claim.Messages():
			c.msgHandler(message.Topic, message.Key, message.Value)
			session.MarkMessage(message, "")

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}

func newConsumerGroup(opts *KafkaOptions) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	version, _ := sarama.ParseKafkaVersion(opts.Version)
	config.Version = version
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	if opts.FromBeginning {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	return sarama.NewConsumerGroup(opts.Brokers, opts.GroupID, config)
}
