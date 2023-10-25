package handler

import (
	"strconv"
	"strings"

	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"tgbot/internal/assets"
	"tgbot/internal/model"
	rdb "tgbot/internal/redis"
	"tgbot/internal/repository"
	"tgbot/internal/service/callback"
	"tgbot/internal/service/message"
)

const (
	commandsPath   = "internal/assets/commands"
	jsonFormatName = ".json"
)

type Reader struct {
	texts    map[string]string
	bot      *tgbotapi.BotAPI
	logger   *zap.Logger
	rdb      *redis.Client
	msg      *MessageHandlers
	callback *CallBackHandlers
}

func NewReader(log *zap.Logger, rdb *redis.Client, repo *repository.PGRepository, bot *tgbotapi.BotAPI, texts map[string]string) *Reader {
	return &Reader{
		logger:   log,
		rdb:      rdb,
		bot:      bot,
		msg:      newMessagesHandler(message.NewMessageService(log, rdb, repo, bot, texts)),
		callback: newCallbackHandler(callback.NewCallbackService(log, rdb, repo, bot, texts)),
		texts:    texts,
	}
}

func (r *Reader) ReadUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		go r.updateActions(update)
	}
}

func (r *Reader) updateActions(update tgbotapi.Update) {
	if update.Message != nil {
		if strings.Contains(update.Message.Text, "new_team_user_") {
			teamID, err := strconv.Atoi(strings.ReplaceAll(update.Message.Text, "/start new_team_user_", ""))
			if err != nil {
				return
			}
			s := setMessageTeamSituation(update.Message, teamID)

			handler := r.msg.GetHandler("/add_user_team")
			if handler != nil {
				err := handler(s)
				if err != nil {
					r.logger.Error("failed to get handler", zap.Error(err))
				}

				return
			}

			return
		}
		s := setMessageSituation(update.Message)

		handler := r.msg.GetHandler(update.Message.Text)
		if handler != nil {
			err := handler(s)
			if err != nil {
				r.logger.Error("failed to get handler", zap.Error(err))
			}

			return
		}

		path := rdb.GetPath(r.logger, r.rdb, s.Message.Chat.ID)

		handler = r.msg.GetHandler(path)
		if handler != nil {
			err := handler(s)
			if err != nil {
				r.logger.Error("failed to get handler", zap.Error(err))
			}

			return
		}

		command, err := r.GetFromCommands(update.Message.Text)
		if err != nil {
			return
		}

		handler = r.msg.GetHandler(command)
		if handler != nil {
			err = handler(s)
			if err != nil {
				r.logger.Error("failed to get handler", zap.Error(err))
			}

			return
		}

		handler = r.msg.GetHandler("/unrecognized")

		return
	}

	if update.CallbackQuery != nil {
		s := setCallbackSituation(update.CallbackQuery)

		handler := r.callback.GetHandler(update.CallbackQuery.Data)
		err := handler(s)
		if err != nil {
			r.logger.Error("failed to get handler", zap.Error(err))
		}
	}
}

func (r *Reader) GetFromCommands(text string) (string, error) {
	commands, err := assets.LoadJSON(commandsPath + jsonFormatName)
	if err != nil {
		return "", err
	}

	for key, val := range r.texts {
		if val == text {
			for k, v := range commands {
				if key == k {
					return v, nil
				}
			}
		}
	}

	return "", nil
}

func setMessageSituation(message *tgbotapi.Message) *model.Situation {
	return &model.Situation{
		Message: message,
		User:    &model.User{ID: message.Chat.ID},
	}
}

func setMessageTeamSituation(message *tgbotapi.Message, teamId int) *model.Situation {
	return &model.Situation{
		Message: message,
		User:    &model.User{ID: message.Chat.ID},
		TeamID:  teamId,
	}
}

func setCallbackSituation(callback *tgbotapi.CallbackQuery) *model.Situation {
	return &model.Situation{
		CallbackQuery: callback,
		User:          &model.User{ID: callback.Message.Chat.ID},
	}
}

func newMessagesHandler(srv *message.Service) *MessageHandlers {
	handle := MessageHandlers{
		Handlers: map[string]model.Handler{},
	}

	handle.Init(srv)
	return &handle
}

func newCallbackHandler(srv *callback.Service) *CallBackHandlers {
	handle := CallBackHandlers{
		Handlers: map[string]model.Handler{},
	}

	handle.Init(srv)
	return &handle
}
