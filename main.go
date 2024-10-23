package main

import (
	"context"
	"log/slog"
	"os"
)

func main() {
	ctx := context.Background()

	dbUrl := os.Getenv("POSTGRES_URL")
	domain := os.Getenv("DOMAIN")
	srvAddr := os.Getenv("SERVER_ADDR")

	storage, err := NewPostgresStorage(ctx, dbUrl)
	if err != nil {
		panic(err)
	}
	slog.Info("connected to postgres")

	urlShortener := NewUrlShortener(domain, storage)

	api := NewAPI(urlShortener)

	server := NewServer(srvAddr, api)

	slog.Info("start listening")
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
