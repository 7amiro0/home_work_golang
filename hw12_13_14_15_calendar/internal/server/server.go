package server

import (
	"context"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/app"
	"net/http"
)

var (
	ErrEndBeforeStart = "end \"%v\" before start \"%v\""
)

type Logger interface {
	app.Logger
}

type Application interface {
	app.BaseStorage
	app.ConnCloser
}

type middle struct{}

type middleware interface {
	LoggingMiddleware(next http.Handler, logger Logger) http.Handler
}

type Server struct {
	App    Application
	Logger Logger
	Ctx    context.Context
	Middle middleware
	Addr   string
}

func New(ctx context.Context, logger Logger, app Application, addr string) *Server {
	return &Server{
		Addr:   addr,
		Ctx:    ctx,
		Logger: logger,
		App:    app,
		Middle: middle{},
	}
}
