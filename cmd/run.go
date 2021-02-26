package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/imedvedec/hls-server/hls"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const (
	hlsServerAddress string = "localhost:8090"

	shutdownTimeout time.Duration = 5 * time.Second
)

type application struct {
	logger *zerolog.Logger

	ctx    context.Context
	cancel context.CancelFunc
}

func New() *application {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	consoleWriter := zerolog.NewConsoleWriter()
	logger := zerolog.New(consoleWriter).With().Timestamp().Stack().Logger()

	ctx, cancel := context.WithCancel(context.Background())

	return &application{
		logger: &logger,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (app *application) Run() {
	defer app.cancel()

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, os.Kill)
	go func() {
		sig := <-signals
		app.logger.Info().Msg(fmt.Sprintf("OS signal caught: %v", sig))
		app.cancel()
	}()

	app.serverLifeCycle()

	app.logger.Info().Msg(fmt.Sprintf("Application has finished successfully"))
}

func (app *application) serverLifeCycle() {
	server := hls.New(hlsServerAddress)

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			app.logger.Error().Stack().Err(err).Msg(fmt.Sprintf("Server error on listen and server (%v): %v", hlsServerAddress, err))
			app.cancel()
		}
	}()
	app.logger.Info().Msg(fmt.Sprintf("HLS server started on: %v", hlsServerAddress))

	<-app.ctx.Done()

	app.logger.Info().Msg(fmt.Sprintf("HLS server on '%v' is shutting down", hlsServerAddress))

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		app.logger.Error().Stack().Err(err).Msg(fmt.Sprintf("Server error on shutdown (%v): %v", hlsServerAddress, err))
	}
}
