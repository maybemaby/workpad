package api

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Server struct {
	logger   *slog.Logger
	port     string
	srv      *http.Server
	db       *sql.DB
	sqliteDB *sqlx.DB
	services *services
	prod     bool
}

func NewServer(isProd bool) (*Server, error) {

	server := &Server{
		port: "8000",
		prod: isProd,
	}

	server.WithLogger(isProd)

	// Initialize SQLite connection
	sqliteDB, sqlDB, err := NewSqliteDB(context.Background(), !isProd)
	if err != nil {
		return nil, err
	}

	server.sqliteDB = sqliteDB
	server.db = sqlDB

	services := newServices(server.logger)
	server.services = services

	return server, nil
}

func (s *Server) Start(ctx context.Context) error {

	s.MountRoutesOapi()

	s.logger.Info("Server started at http://localhost:" + s.port)
	s.logger.Info(fmt.Sprintf("Server is running in production mode: %t", s.prod))
	s.logger.Debug("Server is running in debug mode")

	return s.srv.ListenAndServe()
}

func (s *Server) WithLogger(isProd bool) {
	format := JSONFormat
	level := slog.LevelInfo

	if !isProd {
		level = slog.LevelDebug
		format = TEXTFormat
	}

	s.logger = BootstrapLogger(level, format, !isProd)
}

func (s *Server) WithPort(port string) {
	s.port = port
}
