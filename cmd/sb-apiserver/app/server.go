package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewSimpleBlog() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "sb-apiserver",
		Short:        "simple go blog",
		Long:         `my simple go blog`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Hello Blog")
			return nil
		},
		Args: cobra.NoArgs,
	}

	return cmd
}
