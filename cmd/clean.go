package cmd

import (
	"fmt"
	"log"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"

	"github.com/Dal-Papa/awsugar/aws"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean [type]",
	Short: "Clean your AWS account in various places",
	Long: `Clean your AWS account in various places including:
	
	- Soft kill an EC2 instance with a snapshot first
	- Remove deprecated ELB without target instances
	- Remove available volumes and snapshot them
	- Release unattached Elastic IPs and Network Interfaces
	- Remove unused Security Groups
	- Remove unused Launch Configurations`,
	Args: cobra.MinimumNArgs(1),
	Run:  cleanFunc,
}

var cleanFlags struct {
	SweetClean bool
}

func cleanFunc(cmd *cobra.Command, args []string) {
	resource := args[0]
	switch resource {
	case "elb":
		cleanELB()
	case "ebs":
		cleanEBS()
	case "network-interface":
		cleanNetworkInterfaces()
	default:
		fmt.Println("Resource type not supported")
	}
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	cleanCmd.PersistentFlags().BoolVarP(&cleanFlags.SweetClean, "sweet-clean",
		"s", true, "allow some preparation before cleaning (snapshot, etc.)")
}

func cleanAbstractList(list []aws.Deletable) error {
	var retErr *multierror.Error
	for _, d := range list {
		fmt.Printf("%s [%s] to be deleted...\n", d.Type(), d.Name())
		if !rootFlags.DryRun {
			if err := d.Delete(sess); err != nil {
				retErr = multierror.Append(retErr, err)
			} else {
				fmt.Printf("%s [%s] deleted successfully!\n", d.Type(), d.Name())
			}
		}
	}
	return retErr.ErrorOrNil()
}

func sweetenList(list []aws.Sweetener) error {
	var retErr *multierror.Error
	for _, d := range list {
		if !rootFlags.DryRun {
			if err := d.Sweeten(sess); err != nil {
				retErr = multierror.Append(retErr, err)
			}
		}
	}
	return retErr.ErrorOrNil()
}

func cleanELB() {
	res, err := aws.ListInactiveLoadBalancers(sess)
	if err != nil {
		log.Fatal(err)
	}
	deletableList := make([]aws.Deletable, len(res))
	for i, d := range res {
		deletableList[i] = d
	}
	if err := cleanAbstractList(deletableList); err != nil {
		log.Println(err)
	}
}

func cleanNetworkInterfaces() {
	res, err := aws.ListUnattachedNetworkInterfaces(sess)
	if err != nil {
		log.Fatal(err)
	}
	deletableList := make([]aws.Deletable, len(res))
	for i, d := range res {
		deletableList[i] = d
	}
	if err := cleanAbstractList(deletableList); err != nil {
		log.Println(err)
	}
}

func cleanEBS() {
	res, err := aws.ListAvailableEBS(sess)
	if err != nil {
		log.Fatal(err)
	}
	deletableList := make([]aws.Deletable, len(res))
	for i, d := range res {
		deletableList[i] = d
	}
	if cleanFlags.SweetClean {
		sweetList := make([]aws.Sweetener, len(res))
		for i, d := range res {
			sweetList[i] = d
		}
		// To prevent still deleting if error while sweetening.
		// Need to do at the item level.
		if err := sweetenList(sweetList); err != nil {
			log.Fatal(err)
		}
	}
	if err := cleanAbstractList(deletableList); err != nil {
		log.Println(err)
	}
}
