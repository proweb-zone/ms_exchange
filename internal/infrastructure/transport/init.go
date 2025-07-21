package transport

import (
	"ms_exchange/internal/config"
	"ms_exchange/internal/infrastructure/transport/kafka_client"
	"ms_exchange/pkg/logger"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type TransportState struct {
	KafkaConsumer *kafka.Consumer
}

func (t *TransportState) InitTransport(сfgEnv *config.Config, logger *logger.Logger) *TransportState {
	t.KafkaConsumer = GetConsumerIns(сfgEnv, logger)

	return t
}

func GetConsumerIns(сfgEnv *config.Config, logger *logger.Logger) *kafka.Consumer {
	kafkaCfg := kafka_client.Create(сfgEnv)
	return kafka_client.Connect(kafkaCfg, logger)
}
