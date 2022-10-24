package app

import (
	"context"

	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Logger  Logger
	Storage Storage
}

type Logger interface {
	Error(msg ...any)
	Debug(msg ...any)
	Fatal(msg ...any)
	Info(msg ...any)
	Warn(msg ...any)
}

type Storage interface {
	Update(ctx context.Context, event storage.Event) error
	Add(ctx context.Context, event storage.Event) error
	List(ctx context.Context, idUser int64) []storage.Event
	Delete(ctx context.Context, id int64) error
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
}

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) Connect(ctx context.Context) error {
	return a.Storage.Connect(ctx)
}

func (a *App) Close(ctx context.Context) error {
	return a.Storage.Close(ctx)
}

func (a *App) Add(ctx context.Context, event storage.Event) error {
	return a.Storage.Add(ctx, event)
}

func (a *App) Delete(ctx context.Context, id int64) error {
	return a.Storage.Delete(ctx, id)
}

func (a *App) Update(ctx context.Context, newEvent storage.Event) error {
	return a.Storage.Update(ctx, newEvent)
}

func (a *App) List(ctx context.Context, idUser int64) (result []storage.Event) {
	return a.Storage.List(ctx, idUser)
}
