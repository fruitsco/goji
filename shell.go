package goji

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/fx"
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
	log, err := LoggerFromContext(ctx)
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
	s.fxApp = s.createFxApp(shellCtx, config, log, options...)

	// 5. create start context w/ timeout
	startCtx, cancel := context.WithTimeout(shellCtx, s.fxApp.StartTimeout())
	defer cancel()

	// 6. start the application, exit on error
	if err := s.fxApp.Start(startCtx); err != nil {
		return NewShellExitError(1)
	}

	// 7. wait for done signal by OS
	sig := <-s.fxApp.Wait()
	exitCode := sig.ExitCode

	// 8. create shutdown context
	stopCtx, cancel := context.WithTimeout(shellCtx, s.fxApp.StopTimeout())
	defer cancel()

	// 9. gracefully shutdown the app, exit on error
	if err := s.fxApp.Stop(stopCtx); err != nil {
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
	config *RootConfig[C],
	log *zap.Logger,
	options ...fx.Option,
) *fx.App {
	// merge together options
	mergedOptions := make([]fx.Option, 0, len(s.options)+len(options))
	mergedOptions = append(mergedOptions, s.options...)
	mergedOptions = append(mergedOptions, options...)

	// create root module
	fxModule := NewRootModule(ctx, config, log, mergedOptions...)

	// create fx application
	return fx.New(fxModule)
}
