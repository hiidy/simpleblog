package app

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hiidy/simpleblog/cmd/sb-apiserver/app/options"
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
			if err := viper.Unmarshal(opts); err != nil {
				return err
			}

			if err := opts.Validate(); err != nil {
				return err
			}

			fmt.Printf("ServerMode from ServerOptions: %s\n", opts.JWTKey)
			fmt.Printf("ServerMode from VIper: %s\n\n", viper.GetString("jwt-key"))

			jsonData, _ := json.MarshalIndent(opts, "", " ")
			fmt.Println(string(jsonData))
			return nil
		},
		Args: cobra.NoArgs,
	}

	cobra.OnInitialize(onInitialize)

	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", filePath(), "Path to the simpleblog config file.")

	opts.AddFlags(cmd.PersistentFlags())
	return cmd
}
