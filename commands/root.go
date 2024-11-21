package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "journalz-ro",
	Short: "A CLI tool for managing journal entries",
}

func Execute() error {
	return rootCmd.Execute()
}
