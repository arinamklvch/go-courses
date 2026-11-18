package http

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("GET /load-movies", h.GetMoviesRange)
	mux.HandleFunc("GET /last-movie", h.GetLastMovie)

	return mux
}
