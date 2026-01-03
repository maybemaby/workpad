package api

import (
	"log/slog"
)

type services struct {
}

func newServices(logger *slog.Logger) *services {

	return &services{}
}
