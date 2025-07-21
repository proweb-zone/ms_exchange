package app

import (
	"ms_exchange/internal/boot"
	"ms_exchange/internal/config"
	"ms_exchange/internal/dto"
	"ms_exchange/internal/infrastructure/db"
	"ms_exchange/internal/infrastructure/transport"
	"ms_exchange/pkg/logger"
)

// Инициализация
func InitApp(cfgEnv *config.Config, logger *logger.Logger) {
	var dbState db.DbState
	var transportState transport.TransportState

	kafkaListener := &dto.KafkaListener{
		DbState:        dbState.InitDb(cfgEnv, logger),
		TransportState: transportState.InitTransport(cfgEnv, logger),
		Logger:         logger,
		Config:         cfgEnv,
	}

	Run(kafkaListener, logger, cfgEnv)
}

// Запуск
func Run(kafkaListener *dto.KafkaListener, logger *logger.Logger, cfgEnv *config.Config) {
	boot.Serve(kafkaListener, logger, cfgEnv)
}
