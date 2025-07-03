package http

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, deps ChatDeps, topic string) {
	chat := NewChatHandler(deps.ChatService, topic)
	binding := NewChatBindingHandler(deps.ChatBindingService)

	r.Route("/chats", func(r chi.Router) {
		r.Post("/", chat.CreateChat)
		r.Get("/", chat.GetChatByID)
		r.Post("/{id}/advice", chat.RequestAdvice)
		r.Post("/{id}/send-invite", chat.SendInvite)
		r.Post("/{id}/accept-invite", chat.AcceptInvite)
	})

	r.Route("/bindings", func(r chi.Router) {
		r.Post("/", binding.CreateOrUpdateBinding)
		r.Get("/", binding.GetBinding)
		r.Delete("/", binding.DeleteBinding)
	})

	r.Get("/invites", chat.GetPendingInvites)
}
