package service

import (
	"context"

	"github.com/Lolik232/luximapp-api/internal/domain"
)

type INewsService interface {
	INewsServiceFinder
	INewsServiceCreater
}

type INewsServiceFinder interface {
	Fetch(ctx context.Context, offset, count uint) (*[]domain.News, int64, error)
	FindByID(ctx context.Context, id string) (*domain.News, error)
	FindAll(ctx context.Context) (*[]domain.News, int64, error)
	Count(ctx context.Context) (int64, error)
}

type INewsServiceCreater interface {
	Create(ctx context.Context, news *domain.News) (string, error)
}
