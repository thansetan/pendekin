package dto

import "errors"

type URLRequest struct {
	OriginalURL string `json:"original_url"`
}

type URLResponse struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

func (r URLRequest) Validate() error {
	if r.OriginalURL == "" {
		return errors.New("original_url can't be empty")
	}

	return nil
}
