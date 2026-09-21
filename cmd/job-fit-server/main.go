package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"job-fit/internal/application"
	"job-fit/internal/httpapi"
	"job-fit/internal/jev"
)

func run(ctx context.Context, args []string, stderr io.Writer, getenv func(string) string) int {
	logger := slog.New(slog.NewJSONHandler(stderr, nil))
	flags := flag.NewFlagSet("job-fit-server", flag.ContinueOnError)
	flags.SetOutput(stderr)
	profile := flags.String("profile", "", "canonical profile file (required; loaded once at startup)")
	listen := flags.String("listen", "127.0.0.1:8080", "loopback host:port")
	timeout := flags.Duration("evaluation-timeout", 5*time.Minute, "total evaluation deadline (positive, at most 1h)")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	fail := func(code string) int { logger.Error("server startup failed", "error", code); return 1 }
	if *profile == "" || flags.NArg() != 0 || *timeout <= 0 || *timeout > time.Hour {
		return fail("INVALID_SERVER_CONFIGURATION")
	}
	host, _, err := net.SplitHostPort(*listen)
	if err != nil || net.ParseIP(host) == nil || !net.ParseIP(host).IsLoopback() {
		return fail("LOOPBACK_ADDRESS_REQUIRED")
	}
	client, err := jev.New(getenv("TYPESAFE_API_KEY"))
	if err != nil {
		return fail("INTERNAL_CONFIGURATION_ERROR")
	}
	name, model := client.Identity()
	service, err := application.New(*profile, client, name, model)
	if err != nil {
		return fail("INTERNAL_CONFIGURATION_ERROR")
	}
	listener, err := net.Listen("tcp", *listen)
	if err != nil {
		return fail("LISTEN_FAILED")
	}
	srv := &http.Server{Handler: httpapi.New(service, *timeout, logger), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: *timeout + 30*time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 16 * 1024, BaseContext: func(net.Listener) context.Context { return ctx }, ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError)}
	logger.Info("server listening", "address", listener.Addr().String(), "evaluationTimeout", timeout.String())
	done := make(chan error, 1)
	go func() { done <- srv.Serve(listener) }()
	select {
	case err := <-done:
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server stopped", "error", "SERVE_FAILED")
			return 1
		}
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdown); err != nil {
			_ = srv.Close()
		}
		<-done
	}
	logger.Info("server stopped")
	return 0
}
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	os.Exit(run(ctx, os.Args[1:], os.Stderr, os.Getenv))
}
