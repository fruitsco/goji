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

func (t *Shell[C]) MockType(mockType reflect.Type, mock any) {
	if t.app != nil {
		t.tb.Fatalf("cannot add mock after app has started")
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	t.mocks[mockType] = mock
}

func (t *Shell[C]) Mock(mock any) {
	t.MockType(reflect.TypeOf(mock), mock)
}

func (t *Shell[C]) GetMock(mockType reflect.Type) (any, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	mock, ok := t.mocks[mockType]
	return mock, ok
}

func (t *Shell[C]) Start() error {
	return t.StartCtx(context.Background())
}

func (t *Shell[C]) StartCtx(ctx context.Context) error {
	cfg, err := t.parseConfig()
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// create list of options for mocks
	mocks := make([]any, 0, len(t.mocks))
	for mockType, mock := range t.mocks {
		mocks = append(mocks, fx.Annotate(mock, fx.As(mockType)))
	}

	fxOptions := []fx.Option{}

	// add the user options
	fxOptions = append(fxOptions, t.options...)

	// wrap the mocks in a fx.Replace option
	fxOptions = append(fxOptions, fx.Replace(mocks...))

	// create the user module
	if t.params != nil && t.params.ModuleFn != nil {
		fxOptions = append(fxOptions, t.params.ModuleFn(cfg.Child))
	}

	// create the root module by composing the mocks and the user module
	fxRootModule := goji.NewRootModule(ctx, cfg, t.log, fxOptions...)

	// create the fx application
	t.app = fxtest.New(t.tb, fxRootModule)

	// start the application
	t.app.RequireStart()

	return nil
}

func (t *Shell[C]) Stop() {
	defer t.log.Sync()

	if t.app != nil {
		t.app.RequireStop()
	}
}

func (t *Shell[C]) parseConfig() (*goji.RootConfig[C], error) {
	appName := "Goji Test Shell"
	if t.params != nil && t.params.AppName != "" {
		appName = t.params.AppName
	}

	defaultConfig := conf.DefaultConfig{}
	if t.params != nil && t.params.DefaultConfig != nil {
		defaultConfig = t.params.DefaultConfig
	}

	envPrefix := ""
	if t.params != nil && t.params.EnvPrefix != "" {
		envPrefix = t.params.EnvPrefix
	}

	configFileName := ""
	if t.params != nil && t.params.ConfigFileName == "" {
		configFileName = t.params.ConfigFileName
	}

	// parse config using env
	cfg, err := conf.Parse[goji.RootConfig[C]](conf.ParseOptions{
		AppName:     appName,
		Environment: string(goji.EnvironmentTest),
		Defaults:    defaultConfig,
		Prefix:      envPrefix,
		FileName:    configFileName,
		Log:         t.log,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// apply user config function
	if t.params != nil && t.params.ConfigFn != nil {
		cfg.Child = t.params.ConfigFn(cfg.Child)
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
