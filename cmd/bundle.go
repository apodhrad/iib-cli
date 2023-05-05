/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"

	"github.com/apodhrad/iib-cli/format"
	"github.com/apodhrad/iib-cli/grpc"
	"github.com/spf13/cobra"
)

// bundleCmd represents the bundle command
var bundleCmd = &cobra.Command{
	Use:   "bundle",
	Short: "List a specific bundle",
	Long:  "List a specific bundle",
	RunE:  bundleCmdRunE,
}

func init() {
	getCmd.AddCommand(bundleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bundleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bundleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func bundleCmdRunE(cmd *cobra.Command, args []string) error {
	var err error
	var out string

	if len(args) < 3 {
		return errors.New("Specify a package, channel and csv!")
	}

	address := grpc.GrpcStart()
	defer grpc.GrpcStop()
	client, err := grpc.NewClient(address)
	defer client.Close()

	bundle, err := client.GetBundle(args[0], args[1], args[2])

	// TODO: find out why the object cannot be printed - fmt.Println() freezes
	// As a workaround we will skip it
	bundle.Object = nil

	// Just for the sake of better readability
	bundle.CsvJson = ""

	if output == "json" {
		out, err = format.Json(bundle, true)
	} else {
		data := [][]string{}
		// headers
		data = append(data, []string{"PACKAGE", "CHANNEL", "CSV"})
		// items
		data = append(data, []string{
			bundle.PackageName,
			bundle.ChannelName,
			bundle.CsvName,
		})
		out, err = format.Table(data)
	}

	printOutput(out)

	return err
}
