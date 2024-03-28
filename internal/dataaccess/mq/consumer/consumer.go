package consumer

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
	"github.com/maxuanquang/idm/internal/configs"
	"github.com/maxuanquang/idm/internal/utils"
	"go.uber.org/zap"
)

type Consumer interface {
	Start(ctx context.Context) error
	RegisterHandler(queueName string, handlerFunc HandlerFunc)
}

func NewConsumer(
	mqConfig configs.MQ,
	logger *zap.Logger,
) (Consumer, error) {
	saramaConsumer, err := sarama.NewConsumer(mqConfig.Addresses, newSaramaConfig(mqConfig))
	if err != nil {
		logger.With(zap.Error(err)).Error("can not create sarama consumer")
		return nil, err
	}

	return &consumer{
		saramaConsumer:            saramaConsumer,
		logger:                    logger,
		queueNameToHandlerFuncMap: make(map[string]HandlerFunc),
	}, nil
}

type HandlerFunc func(ctx context.Context, payload []byte) error

type consumer struct {
	logger                    *zap.Logger
	queueNameToHandlerFuncMap map[string]HandlerFunc
	saramaConsumer            sarama.Consumer
}

// RegisterHandler implements Consumer.
func (c *consumer) RegisterHandler(queueName string, handlerFunc HandlerFunc) {
	c.queueNameToHandlerFuncMap[queueName] = handlerFunc
}

// Start implements Consumer.
func (c *consumer) Start(ctx context.Context) error {
	logger := utils.LoggerWithContext(ctx, c.logger)

	exitSignalChannel := make(chan os.Signal, 1)
	signal.Notify(exitSignalChannel, os.Interrupt)

	for queueName, handlerFunc := range c.queueNameToHandlerFuncMap {
		go func(queueName string, handerFunc HandlerFunc) {
			err := c.consume(queueName, handerFunc, exitSignalChannel)
			if err != nil {
				logger.With(zap.String("queueName", queueName)).With(zap.Error(err)).Error("failed to consume message from queue")
			}
		}(queueName, handlerFunc)
	}

	<-exitSignalChannel
	return nil
}

func (c *consumer) consume(queueName string, handlerFunc HandlerFunc, exitSignalChannel chan os.Signal) error {
	logger := c.logger.With(zap.String("queueName", queueName))

	partitionConsumer, err := c.saramaConsumer.ConsumePartition(queueName, 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("failed to create sarama partition consumer: %w", err)
	}

	for {
		select {
		case message := <-partitionConsumer.Messages():
			err = handlerFunc(context.Background(), message.Value)
			if err != nil {
				logger.With(zap.Error(err)).Error("failed to handle message")
			}
		case <-exitSignalChannel:
			return nil
		}
	}
}

func newSaramaConfig(mqConfig configs.MQ) *sarama.Config {
	config := sarama.NewConfig()
	config.ClientID = mqConfig.ClientID
	config.Metadata.Full = true
	return config
}
