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

type MockRegistry interface {
	T() testing.TB
	Mock(mock any)
	MockType(mockType reflect.Type, mock any)
	GetMock(mockType reflect.Type) (any, bool)
}

type ConfigFn[C any] func(cfg C) C

type ModuleFn[C any] func(cfg C) fx.Option

type Params[C any] struct {
	// AppName is the name of the application
	AppName string

	// DefaultConfig is the default configuration
	// used when parsing the configuration
	DefaultConfig conf.DefaultConfig

	// ConfigFileName is the name of the configuration file
	ConfigFileName string

	// EnvPrefix is the prefix used for environment variables
	EnvPrefix string

	// ConfigFn is a function that can be used to make
	// changes to the configuration before it is injected
	ConfigFn ConfigFn[C]

	// ModuleFn is a function that is used to create the
	// application module on-demand
	ModuleFn ModuleFn[C]
}

type Bench[C any] struct {
	tb       testing.TB
	moduleFn ModuleFn[C]
	options  []fx.Option
	config   *goji.RootConfig[C]

	mu    sync.Mutex
	mocks map[reflect.Type]any

	app *fxtest.App
	log *zap.Logger
}

var _ = MockRegistry(&Bench[any]{})

func New[C any](
	tb testing.TB,
	params *Params[C],
	options ...fx.Option,
) *Bench[C] {
	cfg, err := parseConfig(params)
	if err != nil {
		tb.Errorf("failed to parse config: %v", err)
	}

	return &Bench[C]{
		tb:       tb,
		moduleFn: params.ModuleFn,
		options:  options,
		config:   cfg,
		mocks:    make(map[reflect.Type]any),
		log:      createTestLogger(),
	}
}

func (s *Bench[C]) T() testing.TB {
	return s.tb
}

func (s *Bench[C]) Config() C {
	return s.config.Child
}

func (s *Bench[C]) MockType(mockType reflect.Type, mock any) {
	if s.app != nil {
		s.tb.Fatalf("cannot add mock after app has started")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.mocks[mockType] = mock
}

func (s *Bench[C]) Mock(mock any) {
	s.MockType(reflect.TypeOf(mock), mock)
}

func (s *Bench[C]) GetMock(mockType reflect.Type) (any, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	mock, ok := s.mocks[mockType]
	return mock, ok
}

func (s *Bench[C]) Start() {
	err := s.StartCtx(context.Background())
	if err != nil {
		s.tb.Fatal(err)
	}
}

func (s *Bench[C]) StartCtx(ctx context.Context) error {
	// create list of options for mocks
	mocks := make([]any, 0, len(s.mocks))
	for mockType, mock := range s.mocks {
		// create new empty value of type and it as an `any`
		mockValue := reflect.New(mockType).Interface()
		mocks = append(mocks, fx.Annotate(mock, fx.As(mockValue)))
	}

	fxOptions := []fx.Option{}

	// add the application options
	fxOptions = append(fxOptions, s.options...)

	// wrap the mocks in a fx.Replace option
	fxOptions = append(fxOptions, fx.Replace(mocks...))

	// create the application module
	if s.moduleFn != nil {
		fxOptions = append(fxOptions, s.moduleFn(s.config.Child))
	}

	// create the root module by composing the mocks and the application module
	fxRootModule := goji.NewRootModule(ctx, s.config, s.log, fxOptions...)

	// create the fx application
	s.app = fxtest.New(s.tb, fxRootModule)

	// start the application
	s.app.RequireStart()

	return nil
}

func (s *Bench[C]) Stop() {
	defer s.log.Sync()

	if s.app != nil {
		s.app.RequireStop()
	}
}

func parseConfig[C any](params *Params[C]) (*goji.RootConfig[C], error) {
	if params == nil {
		params = &Params[C]{}
	}

	appName := "Goji Testbench"
	if params.AppName != "" {
		appName = params.AppName
	}

	defaultConfig := conf.DefaultConfig{}
	if params.DefaultConfig != nil {
		defaultConfig = params.DefaultConfig
	}

	envPrefix := ""
	if params.EnvPrefix != "" {
		envPrefix = params.EnvPrefix
	}

	configFileName := ""
	if params.ConfigFileName == "" {
		configFileName = params.ConfigFileName
	}

	// parse config using env
	cfg, err := conf.Parse[goji.RootConfig[C]](conf.ParseOptions{
		AppName:     appName,
		Environment: string(goji.EnvironmentTest),
		Defaults:    defaultConfig,
		Prefix:      envPrefix,
		FileName:    configFileName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// apply user config function
	if params.ConfigFn != nil {
		cfg.Child = params.ConfigFn(cfg.Child)
	}

	return cfg, nil
}

func GetMock[T, M any](r MockRegistry) *M {

	mockType := reflect.TypeFor[T]()

	mock, ok := r.GetMock(mockType)
	if !ok {
		var t T
		r.T().Fatalf("no mock found for %T", t)
	}

	ret, ok := mock.(*M)
	if !ok {
		var ret *M
		r.T().Fatalf("mock is not of type %T", ret)
	}

	return ret
}

func createTestLogger() *zap.Logger {
	// TODO: logger
	return zap.NewNop()
}
