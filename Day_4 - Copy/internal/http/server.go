package http

import (
	"api-mvc/config"
	"api-mvc/internal/factory"
	"api-mvc/internal/middleware"
	"api-mvc/pkg/logger"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
)

func Start() error {
	f := factory.NewFactory()
	e := echo.New()

	middleware.LogMiddleware(e)

	NewHttp(e, f)

	idleConnsClosed := make(chan struct{})
	go func() {
		defer close(idleConnsClosed)

		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint

		err := e.Shutdown(context.Background())
		if err != nil {
			logger.Log().Err(err).Msg("failed to shutdown server")
		}
	}()

	err := e.Start(fmt.Sprintf(":%v", config.Cfg().APPPort))
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	logger.Log().Info().Msgf("starting server on %s", e.Server.Addr)

	<-idleConnsClosed

	logger.Log().Info().Msg("stopped server gracefully")

	return nil
}
