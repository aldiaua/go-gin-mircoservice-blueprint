package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/example/go-gin-blueprint/internal/config"
	"github.com/example/go-gin-blueprint/internal/database"
	"github.com/example/go-gin-blueprint/internal/httpmiddleware"
	"github.com/example/go-gin-blueprint/internal/transport/httpresponse"
	"github.com/example/go-gin-blueprint/internal/users"
)

type HTTPServer struct {
	httpServer *http.Server
	database   *gorm.DB
	log        *logrus.Logger
}

func NewHTTPServer(cfg config.Config, db *gorm.DB, log *logrus.Logger) *HTTPServer {
	router := gin.New()
	router.Use(gin.Recovery(), httpmiddleware.RequestID(), httpmiddleware.AccessLog(log))
	registerSystemRoutes(router, db)
	registerUserRoutes(router, db)

	return &HTTPServer{
		httpServer: &http.Server{Addr: ":" + cfg.Port, Handler: router},
		database:   db,
		log:        log,
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErrors := make(chan error, 1)
	go func() {
		s.log.WithField("address", s.httpServer.Addr).Info("users service started")
		serverErrors <- s.httpServer.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		return err
	}
	return nil
}

func (s *HTTPServer) Close() error {
	return database.Close(s.database)
}

func registerSystemRoutes(router *gin.Engine, db *gorm.DB) {
	router.GET("/health", func(c *gin.Context) {
		httpresponse.Success(c, http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/ready", func(c *gin.Context) {
		if err := database.Ping(c.Request.Context(), db); err != nil {
			httpresponse.Failure(c, http.StatusServiceUnavailable, "DATABASE_NOT_READY", "database is not ready", nil)
			return
		}
		httpresponse.Success(c, http.StatusOK, gin.H{"status": "ready"})
	})
}

func registerUserRoutes(router *gin.Engine, db *gorm.DB) {
	repository := users.NewRepository(db)
	service := users.NewService(repository)
	handler := users.NewHandler(service)
	users.RegisterRoutes(router.Group("/api/v1"), handler)
}
