package api

import (
	"github.com/go-chi/chi/v5"
)

/*
Init registers all application routes, splitting them into public endpoints

	and endpoints protected by authentication middleware.
*/
func Init(r chi.Router) {
	r.Get("/api/nextdate", NextDayHandler)
	r.Post("/api/signin", SignInHandler)

	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware)

		r.Get("/api/tasks", tasksHandler)

		r.Route("/api/task", func(r chi.Router) {
			r.Post("/", AddTaskHandler)
			r.Get("/", GetTaskHandler)
			r.Put("/", UpdateTaskHandler)
			r.Delete("/", DeleteTaskHandler)

			r.Post("/done", TaskDoneHandler)
		})
	})
}
