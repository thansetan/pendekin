package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/thansetan/pendekin/repository"
	"github.com/thansetan/pendekin/utils"
)

type URLUsecase interface {
	Save(context.Context, string) (string, error)
	Get(context.Context, string) (string, error)
}

type urlUsecaseImpl struct {
	repo   repository.URLRepository
	logger *slog.Logger
}

func NewURLUsecase(repo repository.URLRepository, logger *slog.Logger) URLUsecase {
	return &urlUsecaseImpl{
		repo:   repo,
		logger: logger,
	}
}

func (uc *urlUsecaseImpl) Save(ctx context.Context, longURL string) (string, error) {
	var shortURL string
	err := utils.ValidateURL(longURL)
	if err != nil {
		uc.logger.Error(err.Error())
		return shortURL, fmt.Errorf("invalid URL: %s, perhaps you forgot the scheme(http/https) or the URL itself is invalid (?)", longURL)
	}

	shortURL, err = utils.Shorten(longURL, 5)
	if err != nil {
		uc.logger.Error(err.Error())
		return shortURL, errors.New("sorry, there was an error on our side")
	}

	err = uc.repo.Save(longURL, shortURL)
	if err != nil {
		uc.logger.Error(err.Error())
		return shortURL, errors.New("sorry, there was an error on our side")
	}

	uc.logger.Info("new short URL created", "original_url", longURL, "short_url", shortURL, "creator", ctx.Value("user_ip"))
	return shortURL, nil
}

func (uc *urlUsecaseImpl) Get(ctx context.Context, shortURL string) (string, error) {
	var longURL string

	longURL, err := uc.repo.Get(shortURL)
	if err != nil {
		uc.logger.Error(err.Error())
		return longURL, err
	}

	uc.logger.Info("a short URL is accessed", "short_url", shortURL, "original_url", longURL, "accessor", ctx.Value("user_ip"))

	return longURL, nil
}
