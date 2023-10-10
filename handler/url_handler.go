package handler

import (
	"encoding/json"
	"net/http"

	"github.com/thansetan/pendekin/dto"
	"github.com/thansetan/pendekin/helpers"
	"github.com/thansetan/pendekin/usecase"
)

type urlHandlerImpl struct {
	uc usecase.URLUsecase
}

func NewURLHandler(uc usecase.URLUsecase) *urlHandlerImpl {
	return &urlHandlerImpl{
		uc: uc,
	}
}

func (h *urlHandlerImpl) Save(w http.ResponseWriter, r *http.Request) {
	var data dto.URLRequest

	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		helpers.ResponseBuilder(w, http.StatusUnprocessableEntity, err.Error(), nil)
		return
	}

	err = data.Validate()
	if err != nil {
		helpers.ResponseBuilder(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	shortURL, err := h.uc.Save(ctx, data.OriginalURL)
	if err != nil {
		helpers.ResponseBuilder(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helpers.ResponseBuilder(w, http.StatusCreated, "", dto.URLResponse{
		OriginalURL: data.OriginalURL,
		ShortURL:    shortURL,
	})
}

func (h *urlHandlerImpl) Get(w http.ResponseWriter, r *http.Request, shortURL []byte) {
	ctx := r.Context()

	origURL, err := h.uc.Get(ctx, string(shortURL))
	if err != nil {
		helpers.ResponseBuilder(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	http.Redirect(w, r, origURL, http.StatusTemporaryRedirect)
}