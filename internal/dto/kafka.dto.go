package dto

import (
	"ms_exchange/internal/config"
	"ms_exchange/internal/infrastructure/db"
	"ms_exchange/internal/infrastructure/transport"
	"ms_exchange/pkg/logger"
)

type KafkaListener struct {
	DbState        *db.DbState
	TransportState *transport.TransportState
	Logger         *logger.Logger
	Config         *config.Config
}
