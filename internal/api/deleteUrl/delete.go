package deleteUrl

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

type LinksDeleter interface {
	DeleteUrl(id int64) error
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

func DeleteById(log *logrus.Logger, linksDeleter LinksDeleter) http.HandlerFunc {
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

		err = linksDeleter.DeleteUrl(int64(idInt))
		if err != nil {
			log.WithError(err).Error("Failed to get url")
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, "Error failed to get url")
			return
		}

		w.WriteHeader(http.StatusNoContent)
		log.Info("url deleted id ", id, "by client ", r.RemoteAddr)

	}
}
