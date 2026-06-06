package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type runEHandlerFunc func(cmd *cobra.Command, args []string) error
type runEHandlerFuncWithApp func(cmd *cobra.Command, args []string, cfg *Application) error
type runHandlerFunc func(cmd *cobra.Command, args []string)

var rootCmd = &cobra.Command{
	Use:   "panoptes-server",
	Short: "Panoptes API server.",
	RunE:  handleGroup,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the current version you're using.",
	Run:   withErrorHandler(handleVersion, 1),
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server for the backend components.",
	Run:   withErrorHandler(withApp(handleServe), 1),
}

func withApp(f runEHandlerFuncWithApp) runEHandlerFunc {
	return func(cmd *cobra.Command, args []string) error {
		a, err := NewApplication(cmd)

		if err != nil {
			return err
		}

		return f(cmd, args, a)
	}
}

func withErrorHandler(f runEHandlerFunc, errorExitCode int) runHandlerFunc {
	return func(cmd *cobra.Command, args []string) {
		if err := f(cmd, args); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(errorExitCode)
		}
	}
}

func handleGroup(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

func init() {
	rootCmd.PersistentFlags().String("log-level", "info", "Log level (debug, info, warn, error).")
	rootCmd.PersistentFlags().String("log-format", "json", "Log format (json, text).")

	versionCmd.Flags().Bool("short", false, "Show only the version, excluding commit and date information.")
	rootCmd.AddCommand(versionCmd)

	rootCmd.AddCommand(serveCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
