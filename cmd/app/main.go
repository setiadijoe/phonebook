package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-kit/kit/log"
	_ "github.com/lib/pq"
	"github.com/oklog/oklog/pkg/group"

	"phonebook/cmd/container"
	httpService "phonebook/cmd/http"
	"phonebook/config"
	"phonebook/internal/global"
	"phonebook/pkg/queryable"
)

func main() {
	con := global.DB()
	defer con.Close()

	queryable.MigrateAndSeed(con.DB)

	cfg, err := config.Get()

	if err != nil {
		panic(err)
	}

	logger := global.InitLogger()

	containerHTTP := container.CreateContainer("http", logger)

	var g group.Group

	initHTTP(cfg.HTTPAddress, containerHTTP, &g, logger)

	logger.Log("exit", g.Run())
}

func initHTTP(
	HTTPAddress string,
	container *container.Container,
	g *group.Group,
	logger log.Logger,
) {
	httpLogger := log.With(logger, "component", "http")
	router := chi.NewRouter()
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Client-ID", "Client-Secret"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(corsHandler.Handler)
	router.Handle("/healthy", httpService.HealthyCheck())
	router.Mount("/v1", httpService.MakeHandler(container, httpLogger))
	g.Add(func() error {
		logger.Log("transport", "debug/HTTP", "addr", HTTPAddress)
		return http.ListenAndServe(HTTPAddress, router)
	}, func(err error) {
		if nil != err {
			logger.Log("transport", "debug/HTTP", "addr", HTTPAddress, "error", err)
			panic(err)
		}
	})
}
