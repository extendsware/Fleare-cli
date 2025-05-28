package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {

	buildCmd.Flags().Bool("clean", false, "Clean build artifacts before building")

	RootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build [target]",
	Short: "Build the project",
	Long: `Build the project for specific target or the default target
if none is specified. Available targets are defined in build-config.yml.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := "default"
		if len(args) > 0 {
			target = args[0]
		}

		config, err := LoadConfig("build-config.yaml")
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		executor := NewExecutor(config)
		if err := executor.RunPreBuildTasks(); err != nil {
			fmt.Printf("Pre-build tasks failed: %v\n", err)
			return
		}

		if err := executor.BuildTarget(target); err != nil {
			fmt.Printf("Build failed: %v\n", err)
			return
		}

		if err := executor.RunPostBuildTasks(); err != nil {
			fmt.Printf("Post-build tasks failed: %v\n", err)
			return
		}
	},
}
