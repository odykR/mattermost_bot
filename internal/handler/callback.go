package handler

import (
	"tgbot/internal/model"
	"tgbot/internal/service/callback"
)

type CallBackHandlers struct {
	Handlers map[string]model.Handler
}

func (h *CallBackHandlers) GetHandler(command string) model.Handler {
	return h.Handlers[command]
}

func (h *CallBackHandlers) Init(cs *callback.Service) {
	h.OnCommand("/yes", cs.Yes)
	h.OnCommand("/no", cs.No)
	// Start commands
}

func (h *CallBackHandlers) OnCommand(command string, handler model.Handler) {
	h.Handlers[command] = handler
}
