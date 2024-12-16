package getAll

import (
	"net/http"
	"short-url-api/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

type LinksGetter interface {
	GetAll() ([]storage.Url, error)
}

type Url struct {
	ID    int64  `gorm:"primaryKey"`
	Alias string `gorm:"not null;unique;"`
	Url   string `gorm:"not null"`
}

type Response struct {
	Status string
	Url    string
}

func GetAll(log *logrus.Logger, linksGetter LinksGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.api.redirect.Redirect"
		log := log.WithFields(logrus.Fields{
			"op":         op,
			"request_id": middleware.GetReqID(r.Context()),
		})

		urls, err := linksGetter.GetAll()
		if err != nil {
			log.WithError(err).Error("Failed to get urls")
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, "Error failed to get all urls")
			return
		}

		render.JSON(w, r, urls)
		log.Info("all urls sent to the client ", r.RemoteAddr)

	}
}
