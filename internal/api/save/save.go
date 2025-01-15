package save

import (
	"errors"
	"net/http"
	resp "short-url-api/internal/api/response"
	"short-url-api/pkg/generateShort"
	"short-url-api/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}
type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

// TODO: move to conf
//const aliasLength = 6

type URLSaver interface {
	SaveUrl(urlToSave string, alias string) (storage.Url, error)
}

func New(log *logrus.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.api.save.New"
		log := log.WithFields(logrus.Fields{
			"op":         op,
			"request_id": middleware.GetReqID(r.Context()),
		})

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.WithError(err).Error("failed to decode request body")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}
		log.WithField("request", req).Info("request body decoded")

		if err := validator.New().Struct(req); err != nil {
			log.WithError(err).Error("invalid request")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))

			return
		}
		alias := req.Alias
		if alias == "" {
			alias = generateShort.GenerateShortKey()
		}
		url, err := urlSaver.SaveUrl(req.URL, alias)

		if errors.Is(err, storage.ErrURLExists) {
			log.WithField("url", req.URL).Info("url already exists")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("url already exists"))

			return
		}
		if err != nil {
			log.WithError(err).Error("failed to add url")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}
		log.WithFields(logrus.Fields{
			"id": url.ID,
		}).Info("url added")

		render.JSON(w, r, url)
	}
}
