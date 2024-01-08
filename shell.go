package goji

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

type ExitError struct {
	ExitCode int
}

func (e *ExitError) Error() string {
	return fmt.Sprintf("shell exited with %d", e.ExitCode)
}

func NewShellExitError(exitCode int) *ExitError {
	return &ExitError{ExitCode: exitCode}
}

func IsExitError(err error) bool {
	if err == nil {
		return false
	}

	var exitErr *ExitError
	return errors.As(err, &exitErr)
}

type Shell struct {
	log     *zap.Logger
	fxApp   *fx.App
	options []fx.Option
}

func New(log *zap.Logger, options ...fx.Option) *Shell {
	return &Shell{
		log:     log,
		options: options,
	}
}

func (s *Shell) Run(ctx context.Context, options ...fx.Option) error {
	// 0. after run ends, flush the logger
	defer s.log.Sync()

	// 1. create execution context
	shellCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 2. create fx application
	fxApp := s.createFxApp(shellCtx, options...)
	s.fxApp = fxApp

	// 3. create start context w/ timeout
	startCtx, cancel := context.WithTimeout(shellCtx, fxApp.StartTimeout())
	defer cancel()

	// 4. start the application, exit on error
	if err := fxApp.Start(startCtx); err != nil {
		return NewShellExitError(1)
	}

	// 5. wait for done signal by OS
	sig := <-fxApp.Wait()
	exitCode := sig.ExitCode

	// 6. create shutdown context
	stopCtx, cancel := context.WithTimeout(shellCtx, fxApp.StopTimeout())
	defer cancel()

	// 7. gracefully shutdown the app, exit on error
	if err := fxApp.Stop(stopCtx); err != nil {
		return NewShellExitError(1)
	}

	// 8. return with 0 exit code
	return NewShellExitError(exitCode)
}

func (s *Shell) PrintGraph(ctx context.Context) error {
	return s.Run(ctx, fx.Options(
		fx.NopLogger,
		fx.Invoke(func(graph fx.DotGraph, shutdown fx.Shutdowner) {
			fmt.Println()
			fmt.Println(graph)
			fmt.Println()
			shutdown.Shutdown()
		})),
	)
}

func (s *Shell) createFxApp(ctx context.Context, options ...fx.Option) *fx.App {
	// 1. create fx application
	return fx.New(
		// 2. inject global execution context
		fx.Supply(fx.Annotate(ctx, fx.As(new(context.Context)))),

		// 3. inject the logger
		fx.Supply(s.log),

		// 4. use the logger also for fx' logs
		fx.WithLogger(func() fxevent.Logger {
			return &fxevent.ZapLogger{Logger: s.log.Named("fx")}
		}),

		// 5. provide user-provided options
		fx.Options(s.options...),

		// 5. provide user-provided run options
		fx.Options(options...),
	)
}
