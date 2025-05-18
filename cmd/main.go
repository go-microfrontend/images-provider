package main

import (
	"log/slog"
	"net/http"
	"os"

	"go.temporal.io/sdk/client"

	"github.com/go-microfrontend/images-provider/internal/handlers"
)

const addr = ":8080"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	slog.SetLogLoggerLevel(slog.LevelDebug)

	clt, err := client.Dial(client.Options{HostPort: os.Getenv("TEMPORAL_ADDR"), Logger: logger})
	if err != nil {
		slog.Error("failed to connect to temporal", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer clt.Close()

	imageHandler := handlers.NewImage(
		clt,
		&client.StartWorkflowOptions{TaskQueue: os.Getenv("TASK_QUEUE")},
	)

	server := compose(imageHandler)
	slog.Info("started")
	server.ListenAndServe()
}

func compose(imageHandler http.Handler) *http.Server {
	mux := http.NewServeMux()
	mux.Handle(handlers.ImageEndpoint, imageHandler)

	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}
