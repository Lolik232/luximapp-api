package mongo_store

import (
	"context"
	"time"

	"github.com/Lolik232/luximapp-api/internal/domain"
	"github.com/Lolik232/luximapp-api/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	mopts "go.mongodb.org/mongo-driver/mongo/options"
)

type DbNews struct {
	ID        primitive.ObjectID  `bson:"_id"`
	Title     string              `bson:"title"`
	Text      string              `bson:"text"`
	CreatedAt *primitive.DateTime `bson:"created_at"`
}

type NewsRepository struct {
	database *mongo.Database
	newsCol  *mongo.Collection
}

func NewNewsRepository(db *mongo.Database) *NewsRepository {
	return &NewsRepository{
		database: db,
		newsCol:  db.Collection(newsCollection),
	}
}

func (n *NewsRepository) fetchOne(ctx context.Context, query bson.M, opts ...*mopts.FindOneOptions) (*DbNews, error) {
	var options *mopts.FindOneOptions
	if opts != nil && len(opts) > 0 {
		options = opts[0]
	} else {
		options = nil
	}
	var news DbNews
	err := n.newsCol.FindOne(ctx, query, options).Decode(&news)
	if err != nil {
		return nil, err
	}
	return &news, nil

}
func (n *NewsRepository) fetchMany(ctx context.Context, query bson.M, opts ...*mopts.FindOptions) (*[]DbNews, error) {
	var options *mopts.FindOptions
	if opts != nil && len(opts) > 0 {
		options = opts[0]
	} else {
		options = mopts.Find()
	}
	sort := bson.M{
		"created_at": -1,
	}
	options = options.SetSort(sort)

	var news []DbNews
	cur, err := n.newsCol.Find(ctx, query, options)

	if err != nil {
		return nil, err
	}
	cur.All(ctx, &news)

	return &news, nil
}

func (n *NewsRepository) Fetch(ctx context.Context, offset uint, count uint) (*[]domain.News, int64, error) {
	opts := options.Find()
	opts.SetSkip(int64(offset))
	opts.SetLimit(int64(count))
	query := bson.M{}

	news, err := n.fetchMany(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}

	docCount, err := n.newsCol.CountDocuments(ctx, bson.M{})

	if err != nil {
		return nil, 0, err
	}
	return ToNewsSlice(news), docCount, nil
}

func (n *NewsRepository) FindByID(ctx context.Context, id string) (*domain.News, error) {
	newsObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.ErrInvalidArgument.Newf("Invalid News ID %s", id)
	}
	query := bson.M{
		"_id": newsObjectID,
	}

	var dbnews *DbNews

	err = n.newsCol.FindOne(ctx, query).Decode(&dbnews)
	if err != nil {
		return nil, err
	}
	return ToNews(dbnews), nil
}

func (n *NewsRepository) FindAll(ctx context.Context) (*[]domain.News, int64, error) {
	news, err := n.fetchMany(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}
	count, err := n.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	return ToNewsSlice(news), count, nil
}

func (n *NewsRepository) Create(ctx context.Context, news *domain.News) (string, error) {
	dbnews := ToDbNews(news)
	createdAt := primitive.NewDateTimeFromTime(time.Now())
	dbnews.CreatedAt = &createdAt
	dbnews.ID = primitive.NewObjectID()
	res, err := n.newsCol.InsertOne(ctx, dbnews)
	if err != nil {
		return "", err
	}
	insID := res.InsertedID.(primitive.ObjectID).Hex()
	return insID, nil

}
func (n *NewsRepository) Count(ctx context.Context) (int64, error) {
	count, err := n.newsCol.CountDocuments(ctx, bson.M{})
	return count, err
}

func ToNewsSlice(dbnews *[]DbNews) *[]domain.News {
	news := make([]domain.News, 0)
	for _, v := range *dbnews {
		n := ToNews(&v)
		news = append(news, *n)
	}
	return &news
}

func ToNews(dbnews *DbNews) *domain.News {
	time := dbnews.CreatedAt.Time()
	return &domain.News{
		ID:        dbnews.ID.Hex(),
		Title:     dbnews.Title,
		Text:      dbnews.Text,
		CreatedAt: &time,
	}
}

func ToDbNews(news *domain.News) *DbNews {
	return &DbNews{
		Title: news.Title,
		Text:  news.Text,
	}
}
