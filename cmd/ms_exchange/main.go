package main

import (
	"ms_exchange/internal/app"
	"ms_exchange/internal/config"
	"ms_exchange/internal/utils"
	"ms_exchange/pkg/logger"
)

func main() {
	currentDir := utils.GetProjectPath()

	configPath := config.ParseConfigPathFromCl(currentDir)
	cfgEnv := config.MustInit(configPath)

	zerologLogger := logger.ConfigureLogger()
	loggerIns := logger.NewLogger(zerologLogger)

	app.InitApp(cfgEnv, loggerIns)
}
