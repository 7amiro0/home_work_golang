package server

import (
	"context"
	"fmt"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/app"
	"net"
	"net/http"
	"time"
)

var (
	ErrEndBeforeStart = "end \"%v\" before start \"%v\""
)

type Logger interface {
	app.Logger
}

type Application interface {
	app.StorageCalendar
}

type middle struct{}

func (middle) LoggingMiddleware(next http.Handler, logger Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			logger.Fatal(err)
		}

		date := time.Now()
		method := r.Method
		path := r.URL.Path
		version := r.Proto
		answer := "200"
		browser := r.UserAgent()
		latency := time.Since(start)
		logger.Info(createLog(ip, method, path, version, answer, browser, latency, date))
	})
}

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

func createLog(ip, method, path, version, answer, browser string, latency time.Duration, date time.Time) string {
	return fmt.Sprintf("%s [%s] %s %s %s %s %s %s",
		ip, date.Format("02-01-2006 15:04:05"), method,
		path, version, answer, latency, browser)
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
