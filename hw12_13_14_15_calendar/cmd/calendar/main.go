package main

import (
	"context"
	"flag"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/server/http"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var err error
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig()
	logg := logger.New(config.Logger.Level)
	store := storage.New(config.Logger.Level)
	calendarApp := app.New(logg, store)

	logg.Info("[INFO] App has been created")

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	httpServer := internalhttp.NewHTTPServer(ctx, logg, calendarApp, net.JoinHostPort(config.HTTP.Host, config.HTTP.Port))
	grpcServer := internalgrpc.NewGRPCServer(ctx, logg, calendarApp, net.JoinHostPort(config.GRPC.Host, config.GRPC.Port))

	logg.Info("[INFO] Servers has been created")

	go func() {
		if err = httpServer.Start(ctx); err != nil {
			logg.Error("[ERR] Failed to start http server: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()

	go func() {
		if err = grpcServer.Start(ctx); err != nil {
			logg.Error("[ERR] Failed to start grpc server: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()

	logg.Info("[INFO] Calendar is running")

	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err = httpServer.Stop(ctx); err != nil {
		logg.Error("[ERR] Failed to stop http server: " + err.Error())
	}

	if err = grpcServer.Stop(ctx); err != nil {
		logg.Error("[ERR] Failed to stop http server: " + err.Error())
	}

	logg.Info("[INFO] Calendar has been stopped")
}
