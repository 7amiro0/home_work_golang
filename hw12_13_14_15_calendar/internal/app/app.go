package app

import (
	"context"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage"
	"time"
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

type ConnCloser interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
}

type StorageQueue interface {
	Clear(ctx context.Context) error
	ListByNotify(ctx context.Context, until time.Duration) ([]storage.Event, error)
}

type BaseStorage interface {
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, event *storage.Event) error
	Add(ctx context.Context, event *storage.Event) error
	List(ctx context.Context, userName string) ([]storage.Event, error)
	ListUpcoming(ctx context.Context, userName string, until time.Duration) ([]storage.Event, error)
}

type Storage interface {
	BaseStorage
	ConnCloser
	StorageQueue
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

func (a *App) Add(ctx context.Context, event *storage.Event) error {
	return a.Storage.Add(ctx, event)
}

func (a *App) Delete(ctx context.Context, id int64) error {
	return a.Storage.Delete(ctx, id)
}

func (a *App) Update(ctx context.Context, newEvent *storage.Event) error {
	return a.Storage.Update(ctx, newEvent)
}

func (a *App) List(ctx context.Context, userName string) ([]storage.Event, error) {
	return a.Storage.List(ctx, userName)
}

func (a *App) ListUpcoming(ctx context.Context, userName string, until time.Duration) ([]storage.Event, error) {
	return a.Storage.ListUpcoming(ctx, userName, until)
}
