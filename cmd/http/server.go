package http

import (
	"net/http"
	"phonebook/cmd/container"
	"phonebook/internal/global"

	phonebookhttp "phonebook/internal/phonebook/transport/http"
	kitxserver "phonebook/pkg/httperror"

	"github.com/go-chi/chi"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
)

func MakeHandler(container *container.Container, logger kitlog.Logger) http.Handler {
	router := chi.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(kitxserver.EncodeError),
	}

	registerPhoneBookHandler(router, container, logger, opts)

	return router
}

func HealthyCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := global.DB().Exec("SELECT 1")
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
}

func registerPhoneBookHandler(r *chi.Mux, container *container.Container, logger kitlog.Logger, opts []kithttp.ServerOption) {
	r.Get("/phonebook", phonebookhttp.Get(
		container.PhoneBook,
		logger,
		append(opts)).ServeHTTP)
	r.Post("/phonebook", phonebookhttp.Create(
		container.PhoneBook,
		logger,
		append(opts)).ServeHTTP)
	r.Put("/phonebook/{id}", phonebookhttp.Update(
		container.PhoneBook,
		logger,
		append(opts)).ServeHTTP)
	r.Delete("/phonebook/{id}", phonebookhttp.Remove(
		container.PhoneBook,
		logger,
		append(opts)).ServeHTTP)
}
