package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/happsie/roundabout/internal"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "roundabout",
		Usage:   "Small reverse proxy",
		Version: "0.1",
		Commands: []*cli.Command{
			{
				Name:    "start",
				Aliases: []string{"s"},
				Usage:   "Starts roundabout based on configuration",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Value:    "config.yml",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "loglevel",
						Value:    "info",
						Required: false,
					},
				},
				Action: func(context *cli.Context) error {
					roundaboutFigure := figure.NewColorFigure("Roundabout", "rounded", "blue", true)
					roundaboutFigure.Print()
					setLogLevel(context.String("loglevel"))
					slog.Info("roundabout is starting up")
					conf, err := internal.LoadConfig(context.String("config"))
					if err != nil {
						slog.Error("could not load config, did you create one?", "error", err)
						return err
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
					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		slog.Error("could not start roundabout cli", "error", err)
		return
	}
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

func setLogLevel(logLevel string) {
	switch logLevel {
	case "info":
		slog.SetLogLoggerLevel(slog.LevelInfo)
	case "debug":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case "error":
		slog.SetLogLoggerLevel(slog.LevelError)
	}
}
