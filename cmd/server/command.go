package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "loginhub",
		Short: "loginhub is a oauth2, oidc app",
		Run: func(cmd *cobra.Command, args []string) {
			runApp()
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
