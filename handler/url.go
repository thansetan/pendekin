package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	_ "embed"

	"github.com/thansetan/pendekin/dto"
	"github.com/thansetan/pendekin/helper"
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
		helper.ResponseBuilder[error](w, http.StatusUnprocessableEntity, err.Error(), nil)
		return
	}

	err = data.Validate()
	if err != nil {
		helper.ResponseBuilder[error](w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	shortURL, err := h.uc.Save(ctx, data.OriginalURL)
	if err != nil {
		var errResp helper.ResponseError
		if errors.As(err, &errResp) {
			helper.ResponseBuilder[error](w, errResp.Code(), errResp.Error(), nil)
			return
		}
		helper.ResponseBuilder[error](w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.ResponseBuilder[dto.URLResponse](w, http.StatusCreated, "", dto.URLResponse{
		OriginalURL: data.OriginalURL,
		ShortURL:    shortURL,
	})
}

func (h *urlHandlerImpl) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slug := r.PathValue("shortURL")
	origURL, err := h.uc.Get(ctx, slug)

	if err != nil {
		var errResp helper.ResponseError
		if errors.As(err, &errResp) {
			helper.ResponseBuilder[error](w, errResp.Code(), errResp.Error(), nil)
			return
		}
		helper.ResponseBuilder[error](w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	http.Redirect(w, r, origURL, http.StatusTemporaryRedirect)
}

func (h *urlHandlerImpl) Home(html []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(html)
	}
}
