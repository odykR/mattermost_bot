package callback

import (
	"fmt"

	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"tgbot/internal/model"
	"tgbot/internal/pkg/utils"
	"tgbot/internal/repository"
)

type Service struct {
	log   *zap.Logger
	texts map[string]string
	bot   *tgbotapi.BotAPI
	rdb   *redis.Client
	repo  *repository.PGRepository
}

func NewCallbackService(log *zap.Logger, rdb *redis.Client, repo *repository.PGRepository, bot *tgbotapi.BotAPI, texts map[string]string) *Service {
	return &Service{
		log:   log,
		rdb:   rdb,
		repo:  repo,
		bot:   bot,
		texts: texts,
	}
}

func (c *Service) Yes(s *model.Situation) error {
	err := c.repo.DeleteUserFromTeam(s.User.ID)
	if err != nil {
		return err
	}

	return c.SendMsgToUser(s.User.ID, utils.GetFormatText(c.texts, "you_deleted"))
}

func (c *Service) No(s *model.Situation) error {
	return c.SendMsgToUser(s.User.ID, utils.GetFormatText(c.texts, "choose"))
}

func (c *Service) SendMsgToUser(userID int64, text string) error {
	msg := &tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: userID,
		},
		Text: text,
	}

	_, err := c.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("send msg to user: %w", err)
	}

	return nil
}
