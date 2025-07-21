package boot

import (
	"errors"
	"fmt"
	stocksControllers "ms_exchange/internal/app/stocks/controller"
	"ms_exchange/internal/config"
	"ms_exchange/internal/dto"
	"ms_exchange/pkg/logger"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	_ "github.com/lib/pq"
)

type KafkaServer struct {
	Handlers      map[string]func(msg *kafka.Message) error
	KafkaConsumer *kafka.Consumer
	Topics        []string
	TopicLocks    map[string]*sync.Mutex
	logger        *logger.Logger
}

func Serve(kafkaListener *dto.KafkaListener, loggerIns *logger.Logger, config *config.Config) {
	topics := []string{
		"prices",
		"storages",
		// "categories_b2b",
		// "categories_b2c",
		// "brands",
		// "products",
		// "product_prop_types",
		// "product_prop_values_create",
		// "product_prop_values_update",
	}

	kafkaServer := &KafkaServer{
		KafkaConsumer: kafkaListener.TransportState.KafkaConsumer,
		Topics:        topics,
		Handlers:      buildHandlers(kafkaListener, config),
		logger:        loggerIns,
	}

	if err := kafkaServer.listenTopics(); err != nil {
		loggerIns.Logger.Err(err)
	}
}

// сделать обработчик ошибок
func (ks *KafkaServer) listenTopics() error {
	if err := ks.KafkaConsumer.SubscribeTopics(ks.Topics, nil); err != nil {
		return fmt.Errorf("ошибка подписки на топик: %w", err)
	}

	ks.createLocks()
	done := make(chan any)

	go func() {
		defer close(done)

		for {
			msg, err := ks.KafkaConsumer.ReadMessage(-1)
			if err != nil {
				ks.logger.Error(err, "ошибка чтения сообщения", map[string]any{})

				if msg != nil && msg.TopicPartition.Topic == nil {
					ks.logger.Error(err, "ошибка чтения из несуществующего топика", map[string]any{})

					return
				}

				continue
			}

			lock, exists := ks.TopicLocks[*msg.TopicPartition.Topic]
			if !exists {
				ks.logger.Error(errors.New("mutex для темы не найден"), "", map[string]any{"топик": *msg.TopicPartition.Topic})

				continue
			}

			go ks.runHandler(msg, ks.Handlers[*msg.TopicPartition.Topic], lock)
		}
	}()

	<-done

	return nil
}

func (ks *KafkaServer) runHandler(
	msg *kafka.Message,
	handler func(*kafka.Message) error,
	lock *sync.Mutex,
) {
	lock.Lock()
	defer lock.Unlock()

	if err := handler(msg); err != nil {
		return
	}

	if _, err := ks.KafkaConsumer.CommitMessage(msg); err != nil {
		ks.logger.Error(err, "ошибка подтверждения обработки сообщения", map[string]any{
			*msg.TopicPartition.Topic: "",
		})
	}
}

func (ks *KafkaServer) createLocks() {
	ks.TopicLocks = make(map[string]*sync.Mutex)

	for _, topic := range ks.Topics {
		ks.TopicLocks[topic] = &sync.Mutex{}
	}
}

func buildHandlers(kafkaListener *dto.KafkaListener, config *config.Config) map[string]func(msg *kafka.Message) error {
	stocksController := stocksControllers.GetStocksController(
		kafkaListener.Logger.Logger,
		kafkaListener.DbState.GormIns,
		kafkaListener.TransportState.KafkaConsumer,
	)
	// productCategoriesController := productsControllers.GetProductCategoriesController(
	// 	kafkaListener.Logger.Logger,
	// 	kafkaListener.DbState.GormIns,
	// 	kafkaListener.TransportState.KafkaConsumer,
	// )
	// productsController := productsControllers.GetProductsController(
	// 	kafkaListener.Logger.Logger,
	// 	kafkaListener.DbState.GormIns,
	// 	kafkaListener.TransportState.KafkaConsumer,
	// )
	// productPropertiesController := productsControllers.GetProductPropertiesController(
	// 	kafkaListener.Logger.Logger,
	// 	kafkaListener.DbState.GormIns,
	// 	kafkaListener.TransportState.KafkaConsumer,
	// )
	// brandsController := brandsControllers.GetBrandsController(
	// 	kafkaListener.Logger.Logger,
	// 	kafkaListener.DbState.GormIns,
	// 	kafkaListener.TransportState.KafkaConsumer,
	// )

	return map[string]func(msg *kafka.Message) error{
		"prices":   func(msg *kafka.Message) error { return stocksController.UpdateBatchPrices(msg.Value) },
		"storages": func(msg *kafka.Message) error { return stocksController.UpdateBatchRemainingStock(msg.Value) },
		// "categories_b2c": func(msg *kafka.Message) error {
		// 	return productCategoriesController.UpdateBatchProductCategoriesB2c(msg.Value)
		// },
		// "categories_b2b": func(msg *kafka.Message) error {
		// 	return productCategoriesController.UpdateBatchProductCategoriesB2b(msg.Value)
		// },
		// "brands": func(msg *kafka.Message) error { return brandsController.UpdateBatchBrands(msg.Value) },
		// "products": func(msg *kafka.Message) error {
		// 	return productsController.UpdateBatchProducts(msg.Value, config.External.BaseUrl1C)
		// },
		// "product_prop_types": func(msg *kafka.Message) error {
		// 	return productPropertiesController.UpdateOrCreateBatchProductPropertyTypes(msg.Value)
		// },
		// "product_prop_values_create": func(msg *kafka.Message) error {
		// 	return productPropertiesController.CreateBatchProductPropertyValues(msg.Value)
		// },
		// "product_prop_values_update": func(msg *kafka.Message) error {
		// 	return productPropertiesController.UpdateOrCreateBatchProductPropertyValues(msg.Value)
		// },
	}
}
