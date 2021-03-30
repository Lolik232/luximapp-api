package store

type IStore interface {
	News() INewsRepository
	Logs() ILogRepository
}
