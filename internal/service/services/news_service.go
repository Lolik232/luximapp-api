package services

import (
	"context"

	"github.com/Lolik232/luximapp-api/internal/context_keys"
	"github.com/Lolik232/luximapp-api/internal/domain"
	"github.com/Lolik232/luximapp-api/internal/store"
	"github.com/Lolik232/luximapp-api/pkg/errors"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	minimumCount = 1
	maximumCount = 40
)

type NewsService struct {
	st store.IStore
}

func NewNewsService(store store.IStore) *NewsService {
	return &NewsService{
		st: store,
	}
}

func (ns *NewsService) handleError(ctx context.Context, err error) error {
	data := ctx.Value(context_keys.InputDataContextKey).(context_keys.InputData)
	device := ctx.Value(context_keys.DeviceInfoContextKey).(string)

	log := domain.Log{
		Msg:    err.Error(),
		Data:   data,
		Device: device,
	}

	ns.st.Logs().Create(ctx, &log)

	switch err {
	case mongo.ErrNoDocuments:
		return errors.ErrInvalidArgument.New("")
	default:
		return errors.NoType.Wrap(err, "")
	}
}

func (ns *NewsService) Fetch(ctx context.Context, offset uint, count uint) (*[]domain.News, int64, error) {

	if count > maximumCount {
		return nil, 0, errors.ErrInvalidArgument.Newf("Maximum count is %d your count is %d.", maximumCount, count)
	}

	inputData := context_keys.InputData{
		"method": "NewsService.Fetch",
		"offset": offset,
		"count":  count,
	}

	ctx = context.WithValue(ctx, context_keys.InputDataContextKey, inputData)
	news, fieldsCount, err := ns.st.News().Fetch(ctx, offset, count)

	if err != nil {
		err := ns.handleError(ctx, err)
		return nil, 0, err
	}
	return news, fieldsCount, nil
}

func (ns *NewsService) FindByID(ctx context.Context, id string) (*domain.News, error) {

	inputData := context_keys.InputData{
		"method": "NewsService.FindByID",
		"id":     id,
	}

	ctx = context.WithValue(ctx, context_keys.InputDataContextKey, inputData)

	news, err := ns.st.News().FindByID(ctx, id)

	if err != nil {
		err = ns.handleError(ctx, err)
		if errors.GetType(err) == errors.ErrInvalidArgument {
			return nil, errors.ErrInvalidArgument.Newf("Invalid News ID %s", id)
		}
		return nil, err
	}
	return news, nil
}

func (ns *NewsService) FindAll(ctx context.Context) (*[]domain.News, int64, error) {
	inputData := context_keys.InputData{
		"method": "NewsService.FindAll",
	}
	ctx = context.WithValue(ctx, context_keys.InputDataContextKey, inputData)

	news, count, err := ns.st.News().FindAll(ctx)
	if err != nil {
		err = ns.handleError(ctx, err)
		return nil, 0, err
	}
	return news, count, err
}

func (ns *NewsService) Create(ctx context.Context, news *domain.News) (string, error) {
	news.Sanitize()
	inputData := context_keys.InputData{
		"method": "NewsService.Create",
		"news":   news,
	}
	ctx = context.WithValue(ctx, context_keys.InputDataContextKey, inputData)

	id, err := ns.st.News().Create(ctx, news)
	if err != nil {
		err = ns.handleError(ctx, err)
		return "", err
	}
	return id, nil
}

func (ns *NewsService) Count(ctx context.Context) (int64, error) {
	count, err := ns.st.News().Count(ctx)
	return count, err
}
