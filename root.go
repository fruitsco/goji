package goji

import (
	"context"
	"fmt"
	"os"

	"github.com/fruitsco/goji/conf"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

type RootParams struct {
	AppName        string
	Version        string
	Description    string
	Prefix         string
	Flags          []cli.Flag
	DefaultConfig  conf.DefaultConfig
	ConfigFileName string
}

type Root struct {
	CLI *cli.App
}

func NewCommand[C any](params RootParams) *Root {
	flags := append([]cli.Flag{
		// &cli.BoolFlag{
		// 	Name:               "verbose",
		// 	Aliases:            []string{"v"},
		// 	Count:              &verbosity,
		// 	DisableDefaultText: true,
		// },
		&cli.StringFlag{
			Name:    "env",
			Aliases: []string{"e"},
			EnvVars: []string{
				fmt.Sprintf("%s.ENV", params.Prefix),
				fmt.Sprintf("%s__ENV", params.Prefix),
			},
			Value: "development",
		},
		&cli.StringFlag{
			Name: "log-level",
			EnvVars: []string{
				"LOG_LEVEL",
				fmt.Sprintf("%s.LOG_LEVEL", params.Prefix),
				fmt.Sprintf("%s__LOG_LEVEL", params.Prefix),
			},
		},
		&cli.StringFlag{
			Name: "log-name",
			EnvVars: []string{
				"LOG_NAME",
				fmt.Sprintf("%s.LOG_NAME", params.Prefix),
				fmt.Sprintf("%s__LOG_NAME", params.Prefix),
			},
		},
	}, params.Flags...)

	cliApp := &cli.App{
		Name:    params.AppName,
		Version: params.Version,
		Usage:   params.Description,
		Flags:   flags,
		Before: func(ctx *cli.Context) error {
			// get env from parsed cli flags
			environment := getEnvFromCLI(ctx)

			// create the logger
			log, err := createLogger(ctx, params.AppName, environment)
			if err != nil {
				return err
			}

			// inject logger into cli context
			ctx.Context = contextWithLogger(ctx.Context, log)

			// parse config using env
			cfg, err := conf.Parse[RootConfig[C]](conf.ParseOptions{
				AppName:     params.AppName,
				Environment: string(environment),
				Defaults:    params.DefaultConfig,
				Prefix:      params.Prefix,
				FileName:    params.ConfigFileName,
				Log:         log,
			})
			if err != nil {
				return err
			}

			// inject the config into the cli context
			ctx.Context = contextWithRootConfig(ctx.Context, cfg)

			return nil
		},
		After: func(ctx *cli.Context) error {
			log, err := loggerFromContext(ctx.Context)
			if err != nil {
				return err
			}

			log.Sync()

			return nil
		},
	}

	return &Root{
		CLI: cliApp,
	}
}

func (r *Root) AddCommand(cmd *cli.Command) {
	r.CLI.Commands = append(r.CLI.Commands, cmd)
}

func (r *Root) Run(args []string) {
	r.RunContext(context.Background(), args)
}

func (r *Root) RunContext(ctx context.Context, args []string) {
	err := r.CLI.RunContext(ctx, args)

	// if app exited without error, return
	if err == nil {
		return
	}

	fmt.Printf("exit error: %s\n", err.Error())

	// if app exited with shell.ExitError, exit with given exit code
	if IsExitError(err) {
		os.Exit(err.(*ExitError).ExitCode)
	}

	// otherwise, exit with exit code 1
	os.Exit(1)
}

func init() {
	cli.VersionFlag = &cli.BoolFlag{
		Name:               "version",
		Usage:              "print the version",
		DisableDefaultText: true,
	}
}

func createLogger(
	ctx *cli.Context,
	name string,
	environment Environment,
) (*zap.Logger, error) {
	// get log name from parsed cli flags
	logName := name
	if name := ctx.String("log-name"); name != "" {
		logName = name
	}

	var config zap.Config
	if environment == EnvironmentProduction {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	config.InitialFields = map[string]any{
		"app": logName,
		"env": environment,
	}

	config.Level = getLevelFromCLI(ctx)

	return config.Build()
}
