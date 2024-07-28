package gojitest

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/fruitsco/goji"
	"github.com/fruitsco/goji/conf"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

type ConfigFn[C any] func(cfg C) C

type ModuleFn[C any] func(cfg C) fx.Option

type Params[C any] struct {
	AppName        string
	DefaultConfig  conf.DefaultConfig
	ConfigFileName string
	EnvPrefix      string
	ConfigFn       ConfigFn[C]
	ModuleFn       ModuleFn[C]
}

type Shell[C any] struct {
	tb      testing.TB
	params  *Params[C]
	options []fx.Option

	mu    sync.Mutex
	mocks map[reflect.Type]any

	app *fxtest.App
	log *zap.Logger
}

func New[C any](
	tb testing.TB,
	params *Params[C],
	options ...fx.Option,
) *Shell[C] {
	return &Shell[C]{
		tb:      tb,
		params:  params,
		options: options,
		mocks:   make(map[reflect.Type]any),
		log:     createTestLogger(),
	}
}

func (s *Shell[C]) MockType(mockType reflect.Type, mock any) {
	if s.app != nil {
		s.tb.Fatalf("cannot add mock after app has started")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.mocks[mockType] = mock
}

func (s *Shell[C]) Mock(mock any) {
	s.MockType(reflect.TypeOf(mock), mock)
}

func (s *Shell[C]) GetMock(mockType reflect.Type) (any, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	mock, ok := s.mocks[mockType]
	return mock, ok
}

func (s *Shell[C]) Start() error {
	return s.StartCtx(context.Background())
}

func (s *Shell[C]) StartCtx(ctx context.Context) error {
	cfg, err := s.parseConfig()
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// create list of options for mocks
	mocks := make([]any, 0, len(s.mocks))
	for mockType, mock := range s.mocks {
		// create new empty value of type and it as an `any`
		mockValue := reflect.New(mockType).Interface()
		mocks = append(mocks, fx.Annotate(mock, fx.As(mockValue)))
	}

	fxOptions := []fx.Option{}

	// add the user options
	fxOptions = append(fxOptions, s.options...)

	// wrap the mocks in a fx.Replace option
	fxOptions = append(fxOptions, fx.Replace(mocks...))

	// create the user module
	if s.params != nil && s.params.ModuleFn != nil {
		fxOptions = append(fxOptions, s.params.ModuleFn(cfg.Child))
	}

	// create the root module by composing the mocks and the user module
	fxRootModule := goji.NewRootModule(ctx, cfg, s.log, fxOptions...)

	// create the fx application
	s.app = fxtest.New(s.tb, fxRootModule)

	// start the application
	s.app.RequireStart()

	return nil
}

func (s *Shell[C]) Stop() {
	defer s.log.Sync()

	if s.app != nil {
		s.app.RequireStop()
	}
}

func (s *Shell[C]) parseConfig() (*goji.RootConfig[C], error) {
	appName := "Goji Test Shell"
	if s.params != nil && s.params.AppName != "" {
		appName = s.params.AppName
	}

	defaultConfig := conf.DefaultConfig{}
	if s.params != nil && s.params.DefaultConfig != nil {
		defaultConfig = s.params.DefaultConfig
	}

	envPrefix := ""
	if s.params != nil && s.params.EnvPrefix != "" {
		envPrefix = s.params.EnvPrefix
	}

	configFileName := ""
	if s.params != nil && s.params.ConfigFileName == "" {
		configFileName = s.params.ConfigFileName
	}

	// parse config using env
	cfg, err := conf.Parse[goji.RootConfig[C]](conf.ParseOptions{
		AppName:     appName,
		Environment: string(goji.EnvironmentTest),
		Defaults:    defaultConfig,
		Prefix:      envPrefix,
		FileName:    configFileName,
		Log:         s.log,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// apply user config function
	if s.params != nil && s.params.ConfigFn != nil {
		cfg.Child = s.params.ConfigFn(cfg.Child)
	}

	return cfg, nil
}

func GetMock[T any, C any](t *Shell[C]) (T, error) {
	var ret T

	mockType := reflect.TypeOf(ret)

	mock, ok := t.GetMock(mockType)
	if !ok {
		return ret, fmt.Errorf("no mock found for %T", ret)
	}

	if ret, ok = mock.(T); !ok {
		// this should never happen
		return ret, fmt.Errorf("mock is not of type %T", ret)
	}

	return ret, nil
}

func createTestLogger() *zap.Logger {
	// TODO: logger
	return zap.NewNop()
}
