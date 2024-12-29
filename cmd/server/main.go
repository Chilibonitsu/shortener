package main

import (
	"log"
	"net/http"
	"os"
	"short-url-api/storage/postgres"

	"short-url-api/internal/config"
	"short-url-api/internal/http-server/middleware/mwLogger"

	"short-url-api/internal/api/deleteUrl"
	"short-url-api/internal/api/getAll"
	"short-url-api/internal/api/getByAlias"
	getByID "short-url-api/internal/api/getById"
	"short-url-api/internal/api/save"

	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("No config file .env", err)
	}

	cfg := config.Config{}
	err = env.Parse(&cfg)

	if err != nil {
		log.Fatal("Failed to parse config ", err)
	}

	log := setupLogger(cfg.Env)

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

	log.WithFields(logrus.Fields{
		"address": cfg.HTTPServerAddress, // брать из cfg
	}).Info("starting server")

	srv := &http.Server{
		Addr:         cfg.HTTPServerAddress,
		Handler:      r,
		ReadTimeout:  cfg.HTTPServerTimeout,
		WriteTimeout: cfg.HTTPServerTimeout,
		IdleTimeout:  cfg.HTTPServerIdleTimeout,
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

	case envDev:
		log.SetFormatter(&logrus.JSONFormatter{})
	case envProd:
		log.SetFormatter(&logrus.JSONFormatter{})
	}
	return log
}
