package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/SmoothWay/gophermart/internal/api"
	"github.com/SmoothWay/gophermart/internal/config"
	"github.com/SmoothWay/gophermart/internal/logger"
	postgresrepo "github.com/SmoothWay/gophermart/internal/repository/postgres"
	"github.com/SmoothWay/gophermart/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	middleware "github.com/oapi-codegen/nethttp-middleware"
)

type Server struct {
	srv *http.Server
}

func NewServer(addr string) *Server {
	return &Server{
		srv: &http.Server{
			Addr: addr,
		},
	}
}

func (s *Server) RegisterHandlers(cfg *config.ServerConfig, svc api.Service) {
	swagger, err := api.GetSwagger()
	if err != nil {
		logger.Log().Info("GetSwagger error", slog.String("error", err.Error()))
		return
	}

	swagger.Servers = nil

	r := chi.NewRouter()
	r.Use(middleware.OapiRequestValidator(swagger))
	r.Use(api.Authenticate([]byte(cfg.Secret)))
	r.Use(api.LogRequest())

	gophermart := api.NewGophermart(svc, []byte(cfg.Secret))
	strictHandler := api.NewStrictHandler(gophermart, nil)
	h := api.HandlerFromMux(strictHandler, r)
	s.srv.Handler = h
}

func Run() {
	cfg := config.NewServerConfig()
	logger.InitSlog(cfg.LogLevel)
	db, err := ConnectDB(cfg.DSN)
	if err != nil {
		logger.Log().Info("DB connection err", slog.String("error", err.Error()))
		return
	}
	defer db.Close()

	repo := postgresrepo.New(db)
	client := &http.Client{}
	svc := service.New(repo, client, []byte(cfg.Secret), cfg.AccuralSysAddr)

	server := NewServer(cfg.Host)
	server.RegisterHandlers(cfg, svc)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	var wg = new(sync.WaitGroup)

	wg.Add(1)

	go server.HandleShutdown(ctx, wg)

	orders := repo.ScanOrders(ctx)
	go svc.FetchOrders(ctx, orders)

	logger.Log().Info("Server is listening on", slog.String("address", cfg.Host))
	err = server.srv.ListenAndServe()
	if err != nil {
		logger.Log().Info("Server encountered error", slog.String("error", err.Error()))
		return
	}

	wg.Wait()
}

func ConnectDB(dsn string) (*sql.DB, error) {
	var connection *sql.DB
	var counts int
	var err error
	for {
		connection, err = openDB(dsn)
		if err != nil {
			logger.Log().Info("Database not ready...", slog.String("error", err.Error()))
			counts++
		} else {
			logger.Log().Info("Connected to database")
			break
		}
		if counts > 2 {
			return nil, err
		}
		logger.Log().Info(fmt.Sprintf("Retrying to connect after %d seconds\n", counts+2))
		time.Sleep(time.Duration(2+counts) * time.Second)
	}

	instance, err := postgres.WithInstance(connection, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", instance)
	if err != nil {
		return nil, err
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	return connection, nil
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (s *Server) HandleShutdown(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	<-ctx.Done()
	logger.Log().Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		logger.Log().Info("Shutdown server error", slog.String("error", err.Error()))
		return
	}

	logger.Log().Info("Server stopped gracefully")
}
