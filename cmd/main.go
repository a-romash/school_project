package main

import (
	"log/slog"
	"os"
	"project/internal/server"
	"project/pkg/config"
	"project/pkg/database/postgresql"
)

func main() {
	config.Init()

	var opts *slog.HandlerOptions

	if config.Config.Debug == "true" {
		opts = &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
	} else {
		opts = nil
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	slog.SetDefault(logger)

	db, err := postgresql.Connect()
	if err != nil {
		slog.Error("failed connect to postgresql")
		os.Exit(1)
	}
	defer db.Close()

	var s *server.Server
	switch {
	case config.Config.Server_address == "":
		s = server.NewServer(config.Config.Server_port, db)
	default:
		s = server.NewServer(config.Config.Server_address, db)
	}

	err = s.Start()
	if err != nil {
		slog.Error("server has been stopped", "error", err)
	}
}
