package main

import (
	"fmt"
	"github.com/happsie/roundabout/internal"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Info("roundabout is starting up!")
	conf, err := internal.LoadConfig()
	if err != nil {
		slog.Error("could not load config", "error", err)
		return
	}
	go func() {
		err := internal.NewReverseProxy(conf)
		if err != nil {
			slog.Error("could not start reverse proxy", "error", err)
			os.Exit(1)
		}
	}()
	go func() {
		err := verifyHealth(conf)
		if err != nil {
			slog.Error("could not verify that proxy is running correctly, shutting down...", "error", err)
			os.Exit(1)
		}
	}()
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
}

func verifyHealth(conf *internal.Config) error {
	time.Sleep(5 * time.Second)
	res, err := http.Get(fmt.Sprintf("http://localhost:%s/health", conf.Port))
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("reverse proxy is not responding")
	}
	slog.Info("roundabout reverse proxy started correctly", "port", conf.Port)
	return nil
}
