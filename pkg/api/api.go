package api

import (
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Password string
}

func NewHandler(password string) *Handler {
	return &Handler{
		Password: password,
	}
}

/*
Init registers all application routes, splitting them into public endpoints

	and endpoints protected by authentication middleware.
*/
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/api/nextdate", h.NextDayHandler)
	r.Post("/api/signin", h.SignInHandler)

	r.Group(func(r chi.Router) {
		r.Use(h.AuthMiddleware)

		r.Get("/api/tasks", h.tasksHandler)

		r.Route("/api/task", func(r chi.Router) {
			r.Post("/", h.AddTaskHandler)
			r.Get("/", h.GetTaskHandler)
			r.Put("/", h.UpdateTaskHandler)
			r.Delete("/", h.DeleteTaskHandler)

			r.Post("/done", h.TaskDoneHandler)
		})
	})
}
