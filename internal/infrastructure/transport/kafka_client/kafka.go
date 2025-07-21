package kafka_client

import (
	"fmt"
	"ms_exchange/internal/config"
	"ms_exchange/pkg/logger"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func Connect(kafkaConfig *kafka.ConfigMap, logger *logger.Logger) *kafka.Consumer {
	consumer, err := kafka.NewConsumer(kafkaConfig)
	if err != nil {
		panic("не удалось создать consumer kafka: " + err.Error())
	}

	fmt.Println("сonsumer kafka инициализирован")

	return consumer
}

func Create(config *config.Config) *kafka.ConfigMap {
	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers":  config.Kafka.KafkaServer,
		"group.id":           config.Kafka.KafkaGroupId,
		"auto.offset.reset":  config.Kafka.KafkaOffsetResetType,
		"enable.auto.commit": "false",
		// "debug":             "broker,topic,msg",
		// "log_level":         "7"
	}

	return kafkaConfig
}
