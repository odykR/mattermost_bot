package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"tgbot/config"
	"tgbot/internal/assets"
	"tgbot/internal/handler"
	"tgbot/internal/redis"
	"tgbot/internal/repository"
)

func main() {
	cfg := config.LoadConfig()

	logger, _ := zap.NewProduction()
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		logger.Panic("create bot instance", zap.Error(err))
	}

	db := repository.NewDB(logger, cfg)
	rdbClient, err := redis.NewClient(cfg.RedisDB.Host + ":" + cfg.RedisDB.Port)
	if err != nil {
		logger.Panic("failed to ping redis client", zap.Error(err))
	}
	repo := repository.NewPgRepository(db)

	texts, err := assets.LoadJSON(cfg.TextsPath)
	if err != nil {
		logger.Panic("failed to load texts", zap.Error(err))
	}

	logger.Info("All Databases connected successful!")
	logger.Info("Authorized on account", zap.String("account", bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)

	updates := bot.GetUpdatesChan(u)
	r := handler.NewReader(logger, rdbClient, repo, bot, texts)

	logger.Info("All services are running!")
	r.ReadUpdates(updates)
}
