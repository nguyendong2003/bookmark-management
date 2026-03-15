package service

import (
	"context"
	"errors"

	"github.com/nguyendong2003/bookmark-management/internal/repository"
	"github.com/redis/go-redis/v9"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name ShortenURL --filename shortenurl.go
type ShortenURL interface {
	ShortenURL(ctx context.Context, url string) (string, error)
	GetURL(ctx context.Context, code string) (string, error)
}

type shortenURLService struct {
	urlStorage repository.URLStorage
	codeGen    Password
}

func NewShortenURL(urlStorage repository.URLStorage, codeGen Password) ShortenURL {
	return &shortenURLService{
		urlStorage: urlStorage,
		codeGen:    codeGen,
	}
}

func (s *shortenURLService) ShortenURL(ctx context.Context, url string) (string, error) {
	// gen code
	code, err := s.codeGen.GeneratePassword()
	if err != nil {
		return "", err
	}

	// call repo to save code and url
	err = s.urlStorage.StoreURL(ctx, code, url)
	if err != nil {
		return "", err
	}

	return code, nil
}

var ErrCodeNotExist = errors.New("code not exist")

func (s *shortenURLService) GetURL(ctx context.Context, code string) (string, error) {
	url, err := s.urlStorage.GetURL(ctx, code)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrCodeNotExist
		}
		return "", err
	}

	return url, nil
}
