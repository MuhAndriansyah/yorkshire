package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MuhAndriansyah/yorkshire/internal/config"
	"github.com/MuhAndriansyah/yorkshire/internal/database"
	db "github.com/MuhAndriansyah/yorkshire/internal/database/sqlc"
	"github.com/MuhAndriansyah/yorkshire/internal/handler"
	jwtPackage "github.com/MuhAndriansyah/yorkshire/internal/jwt"
	"github.com/MuhAndriansyah/yorkshire/internal/logging"
	"github.com/MuhAndriansyah/yorkshire/internal/server"
	"github.com/MuhAndriansyah/yorkshire/internal/service/auth"
	"github.com/MuhAndriansyah/yorkshire/internal/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func main() {
	os.Exit(run())
}

func run() (exitCode int) {

	cfg, err := config.Load(".env")
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		return 1
	}

	logger, logFileCloser, err := logging.New(cfg.LogFile)
	if err != nil {
		slog.Error("failed to initialize logger",
			"error", err,
			"log_file", cfg.LogFile,
		)
		return 1
	}
	slog.SetDefault(logger)

	defer func() {
		if closeErr := logFileCloser.Close(); closeErr != nil {
			logger.Error("failed to close log file", "error", closeErr)
			exitCode = 1
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	pool, err := database.NewPostgresqlPool(ctx, cfg.DatabaseUrl)
	if err != nil {
		logger.Error("failed to initialize database connection pool", "error", err)
		return 1
	}
	logger.Info("database connection pool initialized")
	defer pool.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisUrl,
		Password:     cfg.RedisPassword,
		DB:           cfg.RedisDB,
		DialTimeout:  time.Second,
		ReadTimeout:  500 * time.Millisecond,
		WriteTimeout: 500 * time.Millisecond,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Error("failed to ping redis connection")
		return 1
	}
	logger.Info("redis connection initialized")
	defer rdb.Close()

	privKeyBytes, err := os.ReadFile(cfg.PrivKeyPath)
	if err != nil {
		logger.Error("failed to read private key", "error", err)
		return 1
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privKeyBytes)
	if err != nil {
		logger.Error("failed to parse from privkey pem", "error", err)
		return 1
	}

	pubKeyBytes, err := os.ReadFile(cfg.PubKeyPath)
	if err != nil {
		logger.Error("failed to read public key", "error", err)
		return 1
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes)
	if err != nil {
		logger.Error("failed to parse from pubkey pem", "error", err)
		return 1
	}

	jwtService := jwtPackage.New(&jwtPackage.NewJWTOptions{
		PrivateKey:    privKey,
		PublicKey:     pubKey,
		Issuer:        cfg.JWTIssuer,
		Audience:      cfg.JWTAudience,
		TokenDuration: cfg.JWTTokenDuration,
	})

	// 1. Service Layer
	authService := auth.NewAuth(db.New(pool), jwtService, cfg.JWTTokenDuration)

	// 2. Handler Layer
	val := validator.NewValidator()
	h := handler.New(logger, val, authService)

	// 3. Router Layer (Mounting routes automatically via oapi-codegen)
	httpServer := &http.Server{
		Addr:              cfg.HttpAddr,
		Handler:           server.NewRouter(h),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}

	serverClosed := make(chan struct{})

	go func() {
		<-ctx.Done()
		logger.Info("shutdown signal received; starting graceful shutdown")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer shutdownCancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("graceful shutdown failed", "error", err)
		}
		close(serverClosed)
	}()

	logger.Info("starting HTTP server", "address", cfg.HttpAddr)
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("HTTP server failed", "error", err)
		return 1
	}

	<-serverClosed
	logger.Info("HTTP server stopped gracefully")
	return 0
}
