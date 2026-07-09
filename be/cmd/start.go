/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"harun1804/e-commerce/app"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the e-commerce application",
	Run: func(cmd *cobra.Command, args []string) {
		app.RunApplication()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
