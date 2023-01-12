//go:build !embed
// +build !embed

package cli

import "github.com/spf13/cobra"

func server() *cobra.Command {
	return &cobra.Command{
		Use: "server",
	}
}
