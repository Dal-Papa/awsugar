package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/cobra"
)

var sess *session.Session

var rootCmd = &cobra.Command{
	Use:   "awsugar",
	Short: "AWS Working Sugar",
	Long: `AWS Working Sugar provides a set of useful tools for
	your day to day AWS duties.`,
	Version: "0.0.1",
}

var rootFlags struct {
	DryRun bool
	Region string
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initSession() {
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region: aws.String(rootFlags.Region),
		},
	}))
}

func init() {
	cobra.OnInitialize(initSession)
	rootCmd.PersistentFlags().BoolVarP(&rootFlags.DryRun, "dry-run", "d", false,
		"Toggle a list-only mode without executing any action.")
	rootCmd.PersistentFlags().StringVarP(&rootFlags.Region, "region", "r", "us-west-2",
		"Choose the region to execute the actions in")
}
