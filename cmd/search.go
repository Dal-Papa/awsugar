package cmd

import (
	"net"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search [type]",
	Short: "Search through various AWS services",
	Long: `Provides some helpers to search through services in AWS.
	
	Allows to search for an IP in Route53.`,
	Args: cobra.MinimumNArgs(1),
	Run:  searchFunc,
}

var searchFlags struct {
	IP   []net.IP
	Tags []string
}

func searchFunc(cmd *cobra.Command, args []string) {
	if len(searchFlags.IP) < 1 {
		rootCmd.Usage()
		return
	}
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().IPSliceVarP(&searchFlags.IP, "ip", "", []net.IP{}, "list of IPs to search")
}
