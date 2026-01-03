package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/maybemaby/workpad/api"
)

type Args struct {
	Port   string
	DbPath string
	TZ     string
}

func argParse() Args {
	var args Args
	flag.StringVar(&args.Port, "port", "8000", "port to listen on")
	flag.StringVar(&args.DbPath, "db", "app.db", "path to sqlite db")
	flag.StringVar(&args.TZ, "tz", "", "timezone for date handling")
	flag.Parse()

	return args
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// location, err := time.LoadLocation("UTC")
	if err != nil {
		log.Println("Error loading location")
	}

	// time.Local = location
}

func main() {
	args := argParse()

	// OS Signals
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	loadEnv()

	// Dates will be based off specified timezone, uses local timezone of computer if none specified
	if args.TZ != "" {
		location, err := time.LoadLocation(args.TZ)

		if err != nil {
			panic(err)
		}

		time.Local = location
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	// Otel
	otelShutdown, err := api.SetupOtel(ctx, api.OtelConfig{
		TraceExporter:   api.OtlpGrpcExporter,
		MetricsExporter: api.OtlpGrpcExporter,
		TraceEnabled:    true,
		MetricsEnabled:  true,
		LoggerEnabled:   false,
	})

	if err != nil {
		log.Fatalf("Error setting up otel: %v", err)
	}

	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// Server
	appEnv := os.Getenv("APP_ENV")

	isDebug := appEnv == "development"

	server, err := api.NewServer(!isDebug)

	if err != nil {
		log.Fatalf("Error creating server: %v", err)
		os.Exit(1)
	}

	server.WithPort(args.Port)

	go func() {
		err := server.Start(ctx)

		if err != nil {
			log.Println(fmt.Printf("Error starting server: %v", err))
		}
	}()

	<-ctx.Done()
	stop()

}
