package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shaxbee/butler/internal/root"
	"github.com/shaxbee/butler/product"
	"golang.org/x/sync/errgroup"

	_ "modernc.org/sqlite"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	cfg, err := parseConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	db, err := sql.Open("sqlite", cfg.DB)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "open",
			slog.String("error", err.Error()),
			slog.String("db", cfg.DB),
		)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "ping",
			slog.String("error", err.Error()),
			slog.String("db", cfg.DB),
		)
		os.Exit(1)
	}

	products := product.NewService(db)

	mux := http.NewServeMux()
	root.NewRoutes(logger, products).Register(mux)

	errg, ctx := errgroup.WithContext(ctx)
	err = setupServer(ctx, errg, logger, cfg.Addr, mux)
	if err != nil {
		os.Exit(1)
	}

	logger.LogAttrs(ctx, slog.LevelInfo, "started")

	if err := errg.Wait(); err != nil {
		os.Exit(1)
	}
}

func setupServer(ctx context.Context, errg *errgroup.Group, logger *slog.Logger, addr string, handler http.Handler) error {
	logger = logger.With(slog.String("component", "server"))

	server := &http.Server{
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "listen", slog.String("error", err.Error()), slog.String("addr", addr))
		return err
	}

	logger.LogAttrs(ctx, slog.LevelInfo, "listen", slog.String("addr", lis.Addr().String()))

	errg.Go(func() error {
		if err := server.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.LogAttrs(ctx, slog.LevelError, "serve", slog.String("error", err.Error()))
			return err
		}

		return nil
	})

	errg.Go(func() error {
		<-ctx.Done()

		sctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := server.Shutdown(sctx)
		if err != nil {
			logger.LogAttrs(ctx, slog.LevelError, "shutdown", slog.String("error", err.Error()))
		}

		return err
	})

	return nil
}
