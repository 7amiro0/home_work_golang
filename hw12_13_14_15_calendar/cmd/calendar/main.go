package main

import (
	"context"
	"flag"
	internalgrpc "github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/server/http"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/logger"
	storageMemory "github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage/memory"
	storageSql "github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage/sql"
)

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig()
	logg := logger.New(config.Logger.Level)

	var storage app.Storage
	switch config.Storage {
	case "sql":
		storage = storageSql.New(logg)
	case "memory":
		storage = storageMemory.New()
	}

	calendar := app.New(logg, storage)
	logg.Info("[INFO] New app created")

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	httpServer := internalhttp.NewHTTPServer(ctx, logg, calendar, net.JoinHostPort(config.HTTP.Host, config.HTTP.Port))
	grpcServer, err := internalgrpc.NewGRPCServer(ctx, logg, calendar, net.JoinHostPort(config.GRPC.Host, config.GRPC.Port))
	if err != nil {
		logg.Fatal(err)
	}

	logg.Info("[INFO] Created new servers")

	go func() {
		<-ctx.Done()

		ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err = httpServer.Stop(ctx); err != nil {
			logg.Error("[ERR] Failed to stop http server: " + err.Error())
		}
		if err = grpcServer.Stop(ctx); err != nil {
			logg.Error("[ERR] Failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("[INFO] Calendar is running")

	go func() {
		if err = httpServer.Start(ctx); err != nil {
			logg.Error("[ERR] Failed to start http server: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()

	if err = grpcServer.Start(ctx); err != nil {
		logg.Error("[ERR] Failed to start grpc server: " + err.Error())
		cancel()
		os.Exit(1)
	}
}
