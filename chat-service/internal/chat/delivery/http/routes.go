// @title Chat Service API
// @version 1.0
// @description HTTP API для работы с чатами и AI-привязками.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package http

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, deps ChatDeps, topic string) {
	chat := NewChatHandler(deps.ChatService, topic)
	binding := NewChatBindingHandler(deps.ChatBindingService)

	r.Route("/chats", func(r chi.Router) {
		// @Summary Создать чат
		// @Tags chats
		// @Accept json
		// @Produce json
		// @Param request body createChatRequest true "Тип чата и участники"
		// @Success 200 {object} model.Chat
		// @Failure 400 {string} string "invalid request"
		// @Failure 401 {string} string "unauthorized"
		// @Router /chats [post]
		r.Post("/", chat.CreateChat)

		// @Summary Получить чат по ID
		// @Tags chats
		// @Produce json
		// @Param id query int true "ID чата"
		// @Success 200 {object} model.Chat
		// @Failure 400 {string} string "invalid id"
		// @Failure 404 {string} string "not found"
		// @Router /chats [get]
		r.Get("/", chat.GetChatByID)

		// @Summary Запросить AI-совет по чату
		// @Tags chats
		// @Param id path int true "ID чата"
		// @Success 202
		// @Failure 400 {string} string "invalid chat ID"
		// @Failure 401 {string} string "unauthorized"
		// @Router /chats/{id}/advice [post]
		r.Post("/{id}/advice", chat.RequestAdvice)

		// @Summary Отправить инвайт в чат
		// @Tags chats
		// @Accept json
		// @Param id path int true "ID чата"
		// @Param request body sendInviteRequest true "Список userID"
		// @Success 204
		// @Failure 403 {string} string "forbidden"
		// @Failure 401 {string} string "unauthorized"
		// @Router /chats/{id}/send-invite [post]
		r.Post("/{id}/send-invite", chat.SendInvite)

		// @Summary Принять инвайт в чат
		// @Tags chats
		// @Param id path int true "ID чата"
		// @Success 204
		// @Failure 401 {string} string "unauthorized"
		// @Router /chats/{id}/accept-invite [post]
		r.Post("/{id}/accept-invite", chat.AcceptInvite)
	})

	r.Route("/bindings", func(r chi.Router) {
		// @Summary Создать или обновить AI-привязку
		// @Tags bindings
		// @Accept json
		// @Param chat_id query int true "ID чата"
		// @Param request body bindRequest true "Тип привязки"
		// @Success 204
		// @Failure 400 {string} string "invalid chat_id"
		// @Failure 401 {string} string "unauthorized"
		// @Router /bindings [post]
		r.Post("/", binding.CreateOrUpdateBinding)

		// @Summary Получить текущую AI-привязку
		// @Tags bindings
		// @Param chat_id query int true "ID чата"
		// @Success 200 {object} model.ChatBinding
		// @Failure 404 {string} string "binding not found"
		// @Router /bindings [get]
		r.Get("/", binding.GetBinding)

		// @Summary Удалить AI-привязку
		// @Tags bindings
		// @Param chat_id query int true "ID чата"
		// @Success 204
		// @Router /bindings [delete]
		r.Delete("/", binding.DeleteBinding)
	})

	// @Summary Получить список входящих инвайтов
	// @Tags invites
	// @Produce json
	// @Success 200 {array} model.Chat
	// @Failure 401 {string} string "unauthorized"
	// @Router /invites [get]
	r.Get("/invites", chat.GetPendingInvites)
}
