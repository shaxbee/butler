package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/shaxbee/butler/internal/routes"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	fs := flag.NewFlagSet("server", flag.ExitOnError)
	addr := fs.String("addr", ":8080", "listen address")

	_ = fs.Parse(os.Args[1:])

	routes := routes.Routes{}

	errg, ctx := errgroup.WithContext(ctx)

	err := setupServer(context.Background(), errg, logger, *addr, routes.Handler())
	if err != nil {
		os.Exit(1)
	}

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
		logger.Error("listen", slog.String("error", err.Error()), slog.String("addr", addr))
		return err
	}

	logger.Info("listen", slog.String("addr", lis.Addr().String()))

	errg.Go(func() error {
		if err := server.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("serve", slog.String("error", err.Error()))
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
			logger.Error("shutdown", slog.String("error", err.Error()))
		}

		return err
	})

	return nil
}
