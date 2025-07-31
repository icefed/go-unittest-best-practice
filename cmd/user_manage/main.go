package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/pflag"
	"golang.org/x/exp/slog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"go-unittest-best-practice/internal/api"
	"go-unittest-best-practice/internal/config"
	"go-unittest-best-practice/internal/store"
)

func main() {
	flags := pflag.NewFlagSet("test-service", pflag.ExitOnError)
	var conf config.Config
	conf.AddFlags(flags)
	flags.Parse(os.Args)

	slog.Info("load config", "config", conf)

	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
		conf.DBUser, conf.DBPassword, conf.DBHost, conf.DBPort, conf.DBName)), &gorm.Config{})
	if err != nil {
		slog.Error("open database failed", "error", err)
		os.Exit(1)
	}

	svc := api.NewService(store.NewUserRepository(db), &conf)
	apiServer := http.Server{
		Handler: svc,
		Addr:    fmt.Sprintf(":%d", conf.ListenPort),
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("pprof server listening", "pprofAddr", conf.PprofAddr)
		if err := http.ListenAndServe(conf.PprofAddr, nil); err != nil {
			slog.Error("pprof server listen failed", err)
		}
	}()

	go func() {
		<-sigChan
		apiServer.Shutdown(context.Background())
	}()

	slog.Info("server listening", "port", conf.ListenPort)
	if err := apiServer.ListenAndServe(); err != nil {
		slog.Error("server listen failed")
	}
}
