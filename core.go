package main

import (
	"context"
	"fmt"
	"math/big"
)

type Storage interface {
	NextSeqValue(ctx context.Context) (*big.Int, error)
	Save(ctx context.Context, url string, slug string) error
	Get(ctx context.Context, slug string) (string, error)
}

type UrlShortener struct {
	domain  string
	storage Storage
}

func NewUrlShortener(domain string, storage Storage) *UrlShortener {
	return &UrlShortener{
		domain:  domain,
		storage: storage,
	}
}

func (u *UrlShortener) Create(ctx context.Context, url string) (string, error) {
	seq, err := u.storage.NextSeqValue(ctx)
	if err != nil {
		return "", err
	}

	slug := seq.Text(62)
	err = u.storage.Save(ctx, url, slug)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://%s/%s", u.domain, slug), nil
}

func (u *UrlShortener) Get(ctx context.Context, slug string) (string, error) {
	return u.storage.Get(ctx, slug)
}
