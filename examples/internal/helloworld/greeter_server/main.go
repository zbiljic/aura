// Package main implements a server for Greeter service.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sean-/seed"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/zbiljic/aura/go/pkg/aurafx"
	"github.com/zbiljic/aura/go/pkg/cmd"
	aconfig "github.com/zbiljic/aura/go/pkg/config"
	"github.com/zbiljic/aura/go/pkg/logger"
)

func init() {
	seed.Init() //nolint
}

func main() {
	err := execute()
	cmd.ExitIfErr(os.Stderr, err)
}

func execute() error {
	var conf config
	if err := aconfig.Load("", "helloworld", &conf); err != nil {
		return fmt.Errorf("failed to load configuration: %v", err)
	}

	return start(conf)
}

func start(conf config) error {
	log, err := logger.New(conf.Logger)
	if err != nil {
		return fmt.Errorf("error creating logger: %w", err)
	}

	onErrorCh := make(chan error)
	defer close(onErrorCh)

	app := fx.New(
		fx.Supply(conf),
		fx.Logger(aurafx.NewFxLogger(log)),
		fx.Provide(func() *zap.SugaredLogger { return log }),
		fx.Provide(func() chan error { return onErrorCh }),
		fxModule,
	)

	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		log.Error(err)
		return err
	}

	log.Info("started")

	// Wait for termination of interrupt signals.
	// Block main goroutine until it is interrupted.
	select {
	case sig := <-app.Done():
		log.Debugf("received signal: %v", sig)
	case err := <-onErrorCh:
		log.Debugf("lifecycle error: %v", err)
	}

	log.Info("shutting down")

	if err := app.Stop(ctx); err != nil {
		log.Errorf("shutdown: %v", err)
		return err
	}

	return nil
}
