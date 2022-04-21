package routes

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	create_message_handler "otus_dialog_go/internal/handlers/create-message-handler"
	get_messages_handler "otus_dialog_go/internal/handlers/get-messages-handler"
	mark_as_read_handler "otus_dialog_go/internal/handlers/mark-as-read-handler"
)

func RegisterRouter() *chi.Mux {
	r := chi.NewRouter()
	r.With(
		RequestIDMiddleware,
		cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: []string{"https://*", "http://*", "*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}),
	).Route("/api", func(r chi.Router) {
		r.Route("/messages", func(r chi.Router) {
			r.Post("/", create_message_handler.CreateMessageHandler)
			r.Post("/mark_as_read", mark_as_read_handler.MarkAsReadHandler)
			r.Get("/", get_messages_handler.GetMessagesHandler)
		})
	})
	return r
}
