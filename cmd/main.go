package main

import (
	"log/slog"
	"project/internal/server"
	"project/pkg/config"
)

func main() {
	logger := slog.Default()

	config.Init(*logger)

	var s *server.Server
	switch {
	case config.Config.Server_address == "":
		s = server.NewServer(config.Config.Server_port, logger)
	default:
		s = server.NewServer(config.Config.Server_address, logger)
	}

	err := s.Start()
	if err != nil {
		logger.Error("server has been stopped", "error", err)
	}
}
