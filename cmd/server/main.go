package main

import (
	"fmt"
	"net/http"
	"os"
	"short-url-api/internal/config"
	"short-url-api/storage/postgres"
	"time"

	"short-url-api/internal/http-server/middleware/mwLogger"

	"short-url-api/internal/api/deleteUrl"
	"short-url-api/internal/api/getAll"
	"short-url-api/internal/api/getByAlias"
	getByID "short-url-api/internal/api/getById"
	"short-url-api/internal/api/save"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	os.Setenv("CONFIG_PATH", "C:\\Users\\Aboba\\Desktop\\short-url-api\\config\\local.yaml")
	fmt.Println(os.Getenv("CONFIG_PATH"))
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	//log.Info("starting short-url-api", slog.String("env", cfg.Env))

	log.WithFields(logrus.Fields{
		"cfg.Env": cfg.Env,
	}).Info("starting short-url-api")

	db, err := postgres.NewConnectorPostgreSQL()

	if err != nil {
		log.WithError(err).Error("failed to init storage")

		os.Exit(1)
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID, middleware.RealIP, mwLogger.New(log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Get("/{alias}", getByAlias.Redirect(log, db))
	r.Get("/url/all", getAll.GetAll(log, db))
	r.Get("/url/{id_url}", getByID.GetLinkData(log, db))
	r.Post("/url", save.New(log, db))
	r.Delete("/url/delete/{id_url}", deleteUrl.DeleteById(log, db))

	// TODO http server
	log.WithFields(logrus.Fields{
		"address": "localhost", // брать из cfg
	}).Info("starting server")

	srv := &http.Server{
		Addr:         "127.0.0.1:80",
		Handler:      r,
		ReadTimeout:  time.Second * 4,
		WriteTimeout: time.Second * 4,
		IdleTimeout:  time.Second * 4,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error(err)
	}

	log.Error("server stopped")
	// TODO check postman
}

func setupLogger(env string) *logrus.Logger {
	var log *logrus.Logger = logrus.New()

	//log.SetReportCaller(true)
	switch env {
	case envLocal:
		log.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
			DisableQuote:  true,
		})
		//log = setupPrettySlog()
	case envDev:
		log.SetFormatter(&logrus.JSONFormatter{})
	case envProd:
		log.SetFormatter(&logrus.JSONFormatter{})
	}
	return log
}

//slog logger
// func setupLogger(env string) *slog.Logger {
// 	var log *slog.Logger
// 	switch env {
// 	case envLocal:
// 		log = setupPrettySlog()
// 	case envDev:
// 		log = slog.New(
// 			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
// 		)
// 	case envProd:
// 		log = slog.New(
// 			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
// 		)
// 	}
// 	return log
// }

// func setupPrettySlog() *slog.Logger {
// 	opts := slogpretty.PrettyHandlerOptions{
// 		SlogOpts: &slog.HandlerOptions{
// 			Level: slog.LevelDebug,
// 		},
// 	}

// 	handler := opts.NewPrettyHandler(os.Stdout)

// 	return slog.New(handler)
// }
