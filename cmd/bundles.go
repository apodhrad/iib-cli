/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/apodhrad/iib-cli/format"
	"github.com/apodhrad/iib-cli/grpc"
	"github.com/spf13/cobra"
)

// bundlesCmd represents the bundles command
var bundlesCmd = &cobra.Command{
	Use:   "bundles",
	Short: "List all bundles",
	Long:  `List all bundles`,
	RunE:  bundlesCmdRunE,
}

func init() {
	getCmd.AddCommand(bundlesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bundlesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bundlesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func bundlesCmdRunE(cmd *cobra.Command, args []string) error {
	var err error
	var out string

	address := grpc.GrpcStart()
	defer grpc.GrpcStop()
	client, err := grpc.NewClient(address)
	defer client.Close()

	bundles, err := client.GetBundles()

	if output == "json" {
		out, err = format.Json(bundles, true)
	} else {
		data := [][]string{}
		// headers
		data = append(data, []string{"PACKAGE", "CHANNEL", "CSV", "REPLACES"})
		// items
		for _, bundle := range bundles {
			data = append(data, []string{
				bundle.PackageName,
				bundle.ChannelName,
				bundle.CsvName,
				bundle.Replaces,
			})
		}
		out, err = format.Table(data)
	}

	printOutput(out)

	return err
}
