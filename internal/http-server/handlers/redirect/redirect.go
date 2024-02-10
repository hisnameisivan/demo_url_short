package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/hisnameisivan/demo_url_short/internal/http-server/api"
	"github.com/hisnameisivan/demo_url_short/internal/storage"
)

type UrlGetter interface {
	GetUrl(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter UrlGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		childLog := log.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias") // получаем из строки router.Get("{alias}", redirect.New(log, storage))
		if len(alias) == 0 {
			childLog.Info("alias is empty")
			render.JSON(w, r, api.ResponseError("alias is empty"))

			return
		}

		resultUrl, err := urlGetter.GetUrl(alias)
		if err != nil {
			if errors.Is(err, storage.ErrUrlNotFound) {
				childLog.Info("url not found", slog.String("alias", alias))
				render.JSON(w, r, api.ResponseError(storage.ErrUrlNotFound.Error()))

				return
			} else {
				childLog.Info("failed to get url", slog.String("err", err.Error()))
				render.JSON(w, r, api.ResponseError("internal error"))

				return
			}
		}

		childLog.Info("got url", slog.String("url", resultUrl))
		http.Redirect(w, r, resultUrl, http.StatusFound)
	}
}
