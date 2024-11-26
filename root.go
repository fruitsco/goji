package goji

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/fruitsco/goji/conf"
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

type CLIRoot struct {
	CLI *cli.App
}

func NewCommand[C any](params RootParams) *CLIRoot {
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
		// &cli.StringFlag{
		// 	Name: "log-name",
		// 	EnvVars: []string{
		// 		"LOG_NAME",
		// 		fmt.Sprintf("%s.LOG_NAME", params.Prefix),
		// 		fmt.Sprintf("%s__LOG_NAME", params.Prefix),
		// 	},
		// },
	}, params.Flags...)

	cliApp := &cli.App{
		Name:    params.AppName,
		Version: params.Version,
		Usage:   params.Description,
		Flags:   flags,
		Before: func(ctx *cli.Context) error {
			environment := getEnvFromCLI(ctx)

			logLevel := ctx.String("log-level")

			initCtx, err := Init[C](ctx.Context, InitParams{
				AppName:        ctx.App.Name,
				LogLevel:       logLevel,
				Prefix:         params.Prefix,
				Environment:    environment,
				DefaultConfig:  params.DefaultConfig,
				ConfigFileName: params.ConfigFileName,
			})
			if err != nil {
				return err
			}

			ctx.Context = initCtx

			return nil
		},
		After: func(ctx *cli.Context) error {
			log, err := LoggerFromContext(ctx.Context)
			if err != nil {
				return err
			}

			log.Sync()

			return nil
		},
	}

	return &CLIRoot{
		CLI: cliApp,
	}
}

func (r *CLIRoot) AddCommand(cmd *cli.Command) {
	r.CLI.Commands = append(r.CLI.Commands, cmd)
}

func (r *CLIRoot) Run(args []string) {
	r.RunContext(context.Background(), args)
}

func (r *CLIRoot) RunContext(ctx context.Context, args []string) {
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
