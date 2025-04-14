package save

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/random"
)

type URLSaver interface {
	SaveURL(newUrl string,alias string) (error)
}

type Request struct {
	URL string `json:"url" validate:"required,url"`	
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 8

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_method", middleware.GetReqID(r.Context())),
		)

		var req Request

		err:= render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("error request decode", sl.Err(err))
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}

		log.Info("request decode success", slog.Any("request", req))

		if err:= validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("error validate request", sl.Err(err))
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		log.Info("request validate success", slog.Any("request", req))

		alias := req.Alias

		if alias == "" {
			alias = random.String(aliasLength,random.Alphanumeric)
		}

		if err := urlSaver.SaveURL(req.URL, alias); err != nil {
			if(errors.Is(err, storage.ErrUrlExists)) {
				log.Info("url already exists", slog.String("url", req.URL))
				render.JSON(w, r, resp.Error("url already exists"))
				return
			}


			log.Error("error save url", sl.Err(err))
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}

		log.Info("url saved" )

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias: alias,
		})
	}
}