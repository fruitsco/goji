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

type Shell[C any] struct {
	fxApp   *fx.App
	options []fx.Option
}

func New[C any](options ...fx.Option) *Shell[C] {
	return &Shell[C]{
		options: options,
	}
}

func (s *Shell[C]) Run(ctx context.Context, options ...fx.Option) error {
	// 0. get logger from context
	log, err := loggerFromContext(ctx)
	if err != nil {
		return err
	}

	// 1. get config from context
	config, err := rootConfigFromContext[C](ctx)
	if err != nil {
		return err
	}

	// 2. after run ends, flush the logger
	defer log.Sync()

	// 3. create execution context
	shellCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 4. create fx application
	fxApp := s.createFxApp(shellCtx, config, log, options...)
	s.fxApp = fxApp

	// 5. create start context w/ timeout
	startCtx, cancel := context.WithTimeout(shellCtx, fxApp.StartTimeout())
	defer cancel()

	// 6. start the application, exit on error
	if err := fxApp.Start(startCtx); err != nil {
		return NewShellExitError(1)
	}

	// 7. wait for done signal by OS
	sig := <-fxApp.Wait()
	exitCode := sig.ExitCode

	// 8. create shutdown context
	stopCtx, cancel := context.WithTimeout(shellCtx, fxApp.StopTimeout())
	defer cancel()

	// 9. gracefully shutdown the app, exit on error
	if err := fxApp.Stop(stopCtx); err != nil {
		return NewShellExitError(1)
	}

	// 10. return with 0 exit code
	return NewShellExitError(exitCode)
}

func (s *Shell[C]) PrintGraph(ctx context.Context) error {
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

func (s *Shell[C]) createFxApp(
	ctx context.Context,
	config *config[C],
	log *zap.Logger,
	options ...fx.Option,
) *fx.App {
	// 1. create fx application
	return fx.New(
		// 2. inject global execution context
		fx.Supply(fx.Annotate(ctx, fx.As(new(context.Context)))),

		// 3. inject the logger
		fx.Supply(log),

		// 4. inject the app config
		fx.Supply(config.App),

		// 5. inject the log config
		fx.Supply(config.Log),

		// 6. inject the child config
		fx.Supply(&config.Child),

		// 7. use the logger also for fx' logs
		fx.WithLogger(func() fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log.Named("fx")}
		}),

		// 8. provide user-provided options
		fx.Options(s.options...),

		// 9. provide user-provided run options
		fx.Options(options...),
	)
}
