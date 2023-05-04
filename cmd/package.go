/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/apodhrad/iib-cli/format"
	"github.com/apodhrad/iib-cli/grpc"
	"github.com/spf13/cobra"
)

// packageCmd represents the package command
var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "List a specific package",
	Long:  `List a specific package`,
	RunE: func(cmd *cobra.Command, args []string) error {
		funcArgs := PackageCmdArgs{Name: args[1], Output: output}
		out, err := packageCmdFunc(funcArgs)
		fmt.Println(out)
		return err
	},
}

func init() {
	getCmd.AddCommand(packageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// packageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// packageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type PackageCmdArgs struct {
	Name   string
	Output string
}

func packageCmdFunc(args PackageCmdArgs) (string, error) {
	var out string
	var err error

	if args.Name == "" {
		return "", errors.New("Specify a package name!")
	}

	address := grpc.GrpcStart()
	defer grpc.GrpcStop()

	client, err := grpc.NewClient(address)
	defer client.Close()

	pkg, err := client.GetPackage(args.Name)

	if args.Output == "json" {
		out, err = format.Json(pkg, true)
	} else {
		data := [][]string{}
		// headers
		data = append(data, []string{"PACKAGE_NAME", "CHANNEL", "CSV", "DEFAULT"})
		// items
		for _, channel := range pkg.Channels {
			isDefault := ""
			if channel.Name == pkg.DefaultChannelName {
				isDefault = "true"
			}
			data = append(data, []string{pkg.Name, channel.Name, channel.CsvName, isDefault})
		}
		out, err = format.Table(data)
	}

	return out, err
}
