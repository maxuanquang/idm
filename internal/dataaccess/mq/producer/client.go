package producer

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/maxuanquang/idm/internal/configs"
	"github.com/maxuanquang/idm/internal/utils"
	"go.uber.org/zap"
)

type Client interface {
	Produce(ctx context.Context, queueName string, payload []byte) error
}

func NewClient(mqConfig configs.MQ, logger *zap.Logger) (Client, error) {
	producer, err := sarama.NewSyncProducer(mqConfig.Addresses, newSaramaConfig(mqConfig))
	if err != nil {
		return nil, fmt.Errorf("failed to create sarama sync producer: %w", err)
	}

	return &client{
		saramaSyncProducer: producer,
		logger:             logger,
	}, nil
}

type client struct {
	saramaSyncProducer sarama.SyncProducer
	logger             *zap.Logger
}

// Produce implements Client.
func (c *client) Produce(ctx context.Context, queueName string, payload []byte) error {
	logger := utils.LoggerWithContext(ctx, c.logger).
		With(zap.String("queue_name", queueName)).
		With(zap.ByteString("payload", payload))

	_, _, err := c.saramaSyncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: queueName,
		Value: sarama.ByteEncoder(payload),
	})
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to produce message")
		return err
	}

	logger.Info("payload sent to broker")
	return nil
}

func newSaramaConfig(mqConfig configs.MQ) *sarama.Config {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.ClientID = mqConfig.ClientID
	return saramaConfig
}
