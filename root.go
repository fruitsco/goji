package goji

import (
	"context"
	"fmt"
	"os"

	"github.com/fruitsco/goji/x/conf"
	"github.com/fruitsco/goji/x/logging"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

type RootParams struct {
	AppName        string
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
				fmt.Sprintf("%s_ENV", params.Prefix),
			},
			Value: "development",
		},
		&cli.StringFlag{
			Name: "log-level",
			EnvVars: []string{
				"LOG_LEVEL",
				fmt.Sprintf("%s.LOG_LEVEL", params.Prefix),
				fmt.Sprintf("%s_LOG_LEVEL", params.Prefix),
			},
		},
	}, params.Flags...)

	cliApp := &cli.App{
		Name:  params.AppName,
		Usage: params.Description,
		Flags: flags,
		Before: func(ctx *cli.Context) error {
			// get env from parsed cli flags
			environment := GetEnvFromCLI(ctx)

			// create the logger
			log, err := createLogger(ctx, params.AppName, environment)
			if err != nil {
				return err
			}

			// inject logger into cli context
			ctx.Context = logging.ContextWithLogger(ctx.Context, log)

			// parse config using env
			cfg, err := conf.Parse[C](conf.ParseOptions{
				Environment: environment,
				Defaults:    params.DefaultConfig,
				Prefix:      params.Prefix,
				FileName:    params.ConfigFileName,
				Log:         log,
			})
			if err != nil {
				return err
			}

			// inject the config into the cli context
			ctx.Context = conf.ContextWithConfig(ctx.Context, cfg)

			return nil
		},
		After: func(ctx *cli.Context) error {
			log, err := logging.LoggerFromContext(ctx.Context)
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

func (r *Root) Run(args []string) error {
	return r.RunContext(context.Background(), args)
}

func (r *Root) RunContext(ctx context.Context, args []string) error {
	err := r.CLI.RunContext(ctx, args)

	// if app exited without error, return
	if err == nil {
		return nil
	}

	fmt.Printf("exit error: %s\n", err.Error())

	// if app exited with shell.ExitError, exit with given exit code
	if IsExitError(err) {
		os.Exit(err.(*ExitError).ExitCode)
	}

	// otherwise, exit with exit code 1
	os.Exit(1)

	return nil
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
	environment conf.Environment,
) (*zap.Logger, error) {
	level := GetLevelFromCLI(ctx)

	var config zap.Config
	if environment == conf.EnvironmentProduction {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	config.InitialFields = map[string]any{
		"app": name,
		"env": environment,
	}

	config.Level = level

	return config.Build()
}
