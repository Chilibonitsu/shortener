package getByID

import (
	"errors"
	"net/http"
	"short-url-api/storage"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

type GetById interface {
	GetById(id int64) (storage.Url, error)
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

func GetLinkData(log *logrus.Logger, getLink GetById) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.api.redirect.Redirect"
		log := log.WithFields(logrus.Fields{
			"op":         op,
			"request_id": middleware.GetReqID(r.Context()),
		})

		id := chi.URLParam(r, "id_url")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			log.WithFields(logrus.Fields{
				"id_url": id,
			}).WithError(err).Error("Failed to convert id_url to int")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, "Bad id_url")

			return
		}

		url, err := getLink.GetById(int64(idInt))
		if errors.Is(err, storage.ErrURLNotFound) {
			log.WithField("url with id", id).Error("url not found")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, ("url not found"))

			return
		}

		if err != nil {
			log.WithError(err).Error("Failed to get url")
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, "Error failed to get url")
			return
		}

		render.JSON(w, r, url)
		log.WithFields(logrus.Fields{
			"id_url": id,
		}).Info("url info sent to the client", r.RemoteAddr)

	}
}
