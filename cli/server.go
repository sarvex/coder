//go:build !embed
// +build !embed

package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func serverFunc() ServerFunc {
	return nil
}

// nolint:gocyclo
func Server(vip *viper.Viper, newAPI ServerFunc) *cobra.Command {
	root := &cobra.Command{
		Use:   "server",
		Short: "Start a Coder server",
	}

	return root
}
