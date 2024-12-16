package getByAlias

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

type Request struct {
	Alias string `json:"alias"`
}

type Response struct {
	Status string
	Url    string
}

type URLGetter interface {
	GetUrl(alias string) (string, error)
}

// Стоит назвать New?
func Redirect(log *logrus.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.api.redirect.Redirect"
		log := log.WithFields(logrus.Fields{
			"op":         op,
			"request_id": middleware.GetReqID(r.Context()),
		})

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Error("alias is empty")
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, "invalid request")
			return
		}

		// var req Request
		// err := render.DecodeJSON(r.Body, &req)

		// if err != nil {
		// 	log.WithError(err).Error("failed to decode request body")
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	render.JSON(w, r, "Error failed to decode request body")
		// 	return
		// }

		//Validate alias?

		url, err := urlGetter.GetUrl(alias)

		if err != nil {
			log.WithError(err).Error("Failed to get url")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, "Error failed to get url")
			return
		}

		log.WithFields(logrus.Fields{
			"url": url,
		}).Info("got url")

		// render.JSON(w, r, Response{
		// 	Status: "OK",
		// 	Url:    url,
		// })
		http.Redirect(w, r, url, http.StatusFound)

	}
}
