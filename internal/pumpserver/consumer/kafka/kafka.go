// Package kafka implements the `iam-pump/internal/pumpserver/consumer.Consumer` interface.
package kafka

import (
	"context"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/che-kwas/iam-kit/logger"

	"iam-pump/internal/pumpserver/consumer"
)

// Consumer is the kafka consumer group handler.
type Consumer struct {
	ready      chan bool
	msgHandler consumer.MsgHandler
	group      sarama.ConsumerGroup
	topic      string
	poolSize   int
	log        *logger.Logger
}

// NewConsumer returns a kafka consumer.
func NewConsumer(ctx context.Context, msgHandler consumer.MsgHandler) (consumer.Consumer, error) {
	opts, _ := getKafkaOpts()
	log := logger.L()
	log.Debugf("new kafka consumer with options: %+v", opts)

	group, err := newConsumerGroup(opts)
	if err != nil {
		return nil, err
	}

	consumer := &Consumer{
		ready:      make(chan bool),
		msgHandler: msgHandler,
		group:      group,
		topic:      opts.Topic,
		poolSize:   opts.PoolSize,
		log:        log,
	}

	return consumer, nil
}

// Start starts the message consuming loop.
func (c *Consumer) Start(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(c.poolSize)
	for i := 0; i < c.poolSize; i++ {
		go func() {
			defer wg.Done()
			for {
				// `Consume` should be called inside an infinite loop, when a
				// server-side rebalance happens, the consumer session will need to be
				// recreated to get the new claims
				if err := c.group.Consume(ctx, []string{c.topic}, c); err != nil {
					c.log.Errorw("kafka consumer", "error", err)
				}
				// check if context was cancelled, signaling that the consumer should stop
				if ctx.Err() != nil {
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

// Stop stops the message consuming loop.
func (c *Consumer) Stop(ctx context.Context) error {
	return c.group.Close()
}

// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L925-L943
var _ sarama.ConsumerGroupHandler = &Consumer{}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	c.log.Debug("setup kafka consumer")
	// Mark the consumer as ready
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	c.log.Debug("cleanup kafka consumer")
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
			c.log.Infow("consume message", "topic", message.Topic, "key", string(message.Key), "partition", message.Partition, "offset", message.Offset)
			c.msgHandler(session.Context(), message.Value)
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
