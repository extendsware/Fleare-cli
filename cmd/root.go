package cmd

import (
	"fmt"
	"os"

	"github.com/parashmaity/fleare-cli/handler"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "Fleare-cli",
	Short: "Command-line interface for Fleare",
	Long:  "Fleare-cli is a command-line interface for Fleare, an in-memory database.",
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		username, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")

		conn, err := handler.ConnectWithPassword(host, port, username, password)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
			return
		}

		err = handler.HandleCommand(conn)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	},
}

func init() {
	// Add flags for authentication
	RootCmd.Flags().StringP("host", "", "127.0.0.1", "Hostname for the server")
	RootCmd.Flags().IntP("port", "", 9219, "Port for the server")
	RootCmd.Flags().StringP("user", "u", "", "Username for the account")
	RootCmd.Flags().StringP("password", "p", "", "Password for the account")
}

func Execute() {
	// Execute the root command
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
