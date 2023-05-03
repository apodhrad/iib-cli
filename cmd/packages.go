/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/apodhrad/iib-cli/format"
	"github.com/apodhrad/iib-cli/grpc"
	"github.com/spf13/cobra"
)

// packagesCmd represents the packages command
var packagesCmd = &cobra.Command{
	Use:   "packages",
	Short: "List all packages",
	Long:  `List all packages`,
	RunE: func(cmd *cobra.Command, args []string) error {
		funcArgs := PackagesCmdArgs{Output: output}
		out, err := packagesCmdFunc(funcArgs)
		fmt.Println(out)
		return err
	},
}

func init() {
	getCmd.AddCommand(packagesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// packagesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// packagesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type PackagesCmdArgs struct {
	Output string
}

func packagesCmdFunc(args PackagesCmdArgs) (string, error) {
	var out string
	var err error

	address := grpc.GrpcStart()
	client, err := grpc.NewClient(address)
	packageNames, err := client.GetPackageNames()

	if args.Output == "json" {
		out, err = format.Json(packageNames, true)
	} else {
		data := [][]string{}
		// headers
		data = append(data, []string{"PACKAGE_NAME"})
		// items
		for _, packageName := range packageNames {
			data = append(data, []string{packageName.Name})
		}
		out, err = format.Table(data)
	}

	return out, err
}
