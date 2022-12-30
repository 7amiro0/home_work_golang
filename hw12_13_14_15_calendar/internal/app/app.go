package app

import (
	"context"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage"
	"time"
)

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

type App struct {
	Logger  Logger
	Storage StorageCalendar
}

type StorageScheduler interface {
	ListByNotify(ctx context.Context, until time.Duration) (storage.SliceEvents, error)
	Clear(ctx context.Context) error
	ConnCloser
}

type StorageCalendar interface {
	ListUpcoming(ctx context.Context, userName string, until time.Duration) (storage.SliceEvents, error)
	List(ctx context.Context, userName string) (storage.SliceEvents, error)
	Update(ctx context.Context, event *storage.Event) error
	Add(ctx context.Context, event *storage.Event) error
	Delete(ctx context.Context, id int64) error
	ConnCloser
}

func New(logger Logger, storage StorageCalendar) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (c *App) Connect(ctx context.Context) error {
	return c.Storage.Connect(ctx)
}

func (c *App) Close(ctx context.Context) error {
	return c.Storage.Close(ctx)
}

func (c *App) Add(ctx context.Context, event *storage.Event) error {
	return c.Storage.Add(ctx, event)
}

func (c *App) Delete(ctx context.Context, id int64) error {
	return c.Storage.Delete(ctx, id)
}

func (c *App) Update(ctx context.Context, newEvent *storage.Event) error {
	return c.Storage.Update(ctx, newEvent)
}

func (c *App) List(ctx context.Context, userName string) (storage.SliceEvents, error) {
	return c.Storage.List(ctx, userName)
}

func (c *App) ListUpcoming(ctx context.Context, userName string, until time.Duration) (storage.SliceEvents, error) {
	return c.Storage.ListUpcoming(ctx, userName, until)
}
