package store

import (
	"context"

	"github.com/Lolik232/luximapp-api/internal/domain"
)

//INewsRepository combines all interfaces for working with the news
type INewsRepository interface {
	INewsFinder
	INewsCreator
}

//INewsFinder combines methods for find news
type INewsFinder interface {
	Fetch(ctx context.Context, offset, count uint) (*[]domain.News, int64, error)
	FindByID(ctx context.Context, id string) (*domain.News, error)
	FindAll(ctx context.Context) (*[]domain.News, int64, error)
	Count(ctx context.Context) (int64, error)
}

//INewsCreator combines methods for create news
type INewsCreator interface {
	Create(ctx context.Context, news *domain.News) (string, error)
}

type ILogRepository interface {
	Create(ctx context.Context, log *domain.Log) error
}
