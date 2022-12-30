package internalhttp

import (
	"context"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/server"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type HTTPServer struct {
	server     *server.Server
	httpServer *http.Server
}

func NewHTTPServer(ctx context.Context, log server.Logger, app server.Application, addr string) *HTTPServer {
	httpServer := &HTTPServer{
		server: server.New(ctx, log, app, addr),
	}

	httpServer.httpServer = &http.Server{
		Addr:    httpServer.server.Addr,
		Handler: httpServer.server.Middle.LoggingMiddleware(httpServer.createRouter(), httpServer.server.Logger),
	}

	return httpServer
}

func (s *HTTPServer) Start(ctx context.Context) error {
	if err := s.server.App.Connect(ctx); err != nil {
		s.server.Logger.Error("[ERR] App don`t connect: ", err)
		return err
	}

	return s.httpServer.ListenAndServe()
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	if err := s.server.App.Close(ctx); err != nil {
		s.server.Logger.Error("[ERR] App don`t close: ", err)
		return err
	}

	return s.httpServer.Shutdown(ctx)
}

func (s *HTTPServer) createRouter() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc("GET", "/add", s.add)
	router.HandlerFunc("GET", "/list", s.list)
	router.HandlerFunc("GET", "/update", s.update)
	router.HandlerFunc("GET", "/delete", s.delete)
	router.HandlerFunc("GET", "/listUpcoming", s.listUpcoming)
	return router
}
