package save

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/hisnameisivan/demo_url_short/internal/http-server/api"
	"github.com/hisnameisivan/demo_url_short/internal/lib/random"
	"github.com/hisnameisivan/demo_url_short/internal/storage"
)

const aliasLenght = 6

type Request struct {
	Url string `json:"url" validate:"required,url"` // url добавляет проверку на корректность урла в валидаторе
	// Url   string `json:"url" validate:"required"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	// Status string `json:"status"`
	// Error  string `json:"error,omitempty"`
	api.CommonResponse
	Alias string `json:"alias,omitempty"`
}

type UrlSaver interface {
	SaveUrl(urlToSave string, alias string) (int64, error)
}

// FIXME: kill New()
func New(log *slog.Logger, urlSaver UrlSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// FIXME: Тузов использует log - много записей request_id, к тому же они разные
		childLog := log.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		// FIXME: change to unmarshal?
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			childLog.Error("failed to decode request body", slog.String("err", err.Error()))
			render.JSON(w, r, api.ResponseError("failed to decode request"))

			return
		}

		err = validator.New().Struct(req)
		if err != nil {
			validatorErr := err.(validator.ValidationErrors)

			childLog.Error("invalid request", slog.String("err", err.Error()))
			render.JSON(w, r, api.ValidationError(validatorErr))

			return
		}

		alias := req.Alias
		if len(alias) == 0 {
			alias = random.NewRandomString(aliasLenght)
			// FIXME: duplicate
		}

		id, err := urlSaver.SaveUrl(req.Url, alias)
		if err != nil {
			if errors.Is(err, storage.ErrUrlExists) {
				childLog.Info("url already exists", slog.String("url", req.Url))
				render.JSON(w, r, api.ResponseError("url already exists"))

				return
			} else {
				childLog.Error("failed to add url", slog.String("err", err.Error()))
				render.JSON(w, r, api.ResponseError("failed to add url"))

				return
			}
		}

		childLog.Info("url added", slog.Int64("id", id))
		render.JSON(w, r, Response{
			CommonResponse: api.CommonResponse{
				Status: api.StatusOk,
			},
			Alias: alias,
		})
	}
}
