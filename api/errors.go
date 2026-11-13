package api

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
)

type ErrResponse struct {
	HTTPStatusCode int    `json:"-"`
	ErrorText      string `json:"error"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusBadRequest,
		ErrorText:      err.Error(),
	}
}

func ErrInternal(err error) render.Renderer {
	log.Println("internal error:", err)
	return &ErrResponse{
		HTTPStatusCode: http.StatusInternalServerError,
		ErrorText:      http.StatusText(http.StatusInternalServerError),
	}
}

func ErrNotFound() render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusNotFound,
		ErrorText:      http.StatusText(http.StatusNotFound),
	}
}

func ErrUnauthorized() render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusUnauthorized,
		ErrorText:      http.StatusText(http.StatusUnauthorized),
	}
}
