package app

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hiidy/simpleblog/cmd/sb-apiserver/app/options"
	"github.com/hiidy/simpleblog/internal/pkg/log"
)

var configFile string

func NewSimpleBlogCommand() *cobra.Command {
	opts := options.NewServerOptions()
	cmd := &cobra.Command{
		Use:          "sb-apiserver",
		Short:        "simple go blog",
		Long:         `my simple go blog`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(opts)
		},
		Args: cobra.NoArgs,
	}

	cobra.OnInitialize(onInitialize)

	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", filePath(), "Path to the simpleblog config file.")

	opts.AddFlags(cmd.PersistentFlags())
	return cmd
}

func run(opts *options.ServerOptions) error {
	log.Init(logOptions())

	if err := viper.Unmarshal(opts); err != nil {
		return err
	}

	if err := opts.Validate(); err != nil {
		return err
	}

	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	server, err := cfg.NewUnionServer()
	if err != nil {
		return err
	}

	return server.Run()
}

func logOptions() *log.Options {
	opts := log.NewOptions()
	if viper.IsSet("log.disable-caller") {
		opts.DisableCaller = viper.GetBool("log.disable-caller")
	}
	if viper.IsSet("log.disable-stacktrace") {
		opts.DisableStacktrace = viper.GetBool("log.disable-stacktrace")
	}
	if viper.IsSet("log.level") {
		opts.Level = viper.GetString("log.level")
	}
	if viper.IsSet("log.format") {
		opts.Format = viper.GetString("log.format")
	}
	if viper.IsSet("log.output-paths") {
		opts.OutputPaths = viper.GetStringSlice("log.output-paths")
	}
	return opts
}

