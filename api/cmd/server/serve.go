package main

import (
	"github.com/spf13/cobra"
)

func handleServe(cmd *cobra.Command, _ []string, app *Application) error {
	s := app.GetServer()

	return s.Start()
}
