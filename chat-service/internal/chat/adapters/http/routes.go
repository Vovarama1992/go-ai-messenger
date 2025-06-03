package http

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, deps ChatDeps) {
	chat := NewChatHandler(deps.ChatService)
	binding := NewChatBindingHandler(deps.ChatBindingService)

	r.Route("/chats", func(r chi.Router) {
		r.Post("/", chat.CreateChat)
		r.Get("/", chat.GetChatByID)
	})

	r.Route("/bindings", func(r chi.Router) {
		r.Post("/", binding.CreateOrUpdateBinding)
		r.Get("/", binding.GetBinding)
		r.Delete("/", binding.DeleteBinding)
	})
}
