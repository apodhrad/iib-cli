/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/apodhrad/iib-cli/utils"
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
	var bundles []Bundle

	out, err = bundlesCmdGrpc()
	if err == nil {
		bundles, err = bundlesCmdUnmarshal(out)
		if err == nil {
			if output == "json" {
				out, err = bundlesCmdToJson(bundles)
			} else {
				out, err = bundlesCmdToText(bundles)
			}
			if err == nil {
				fmt.Println(out)
			}
		}
	}
	return err
}

func bundlesCmdGrpc() (string, error) {
	var err error
	var out string

	utils.GrpcStart()
	method := "api.Registry/ListBundles"
	grpcArg := utils.GrpcArgMethod(method)
	out, err = utils.GrpcExec(grpcArg)
	utils.GrpcStop()
	return out, err
}

func bundlesCmdUnmarshal(input string) ([]Bundle, error) {
	// grpcurl produces improper json array
	// so, we will fix it to a proper json array
	// numOfPackages := strings.Count(input, "}")
	// input = strings.Replace(input, "}", "},", numOfPackages-1)
	input = strings.ReplaceAll(input, "}\n{", "},\n{")
	input = "[" + input + "]"
	// now, we can unmarshal
	var bundles []Bundle
	err := json.Unmarshal([]byte(input), &bundles)
	return bundles, err
}

func bundlesCmdToText(bundles []Bundle) (string, error) {
	tbl := NewTable("Csv", "Package", "Channel", "Replaces")
	for _, value := range bundles {
		tbl.AddRow(value.CsvName, value.PackageName, value.ChannelName, value.Replaces)
	}
	return TableToString(tbl), nil
}

func bundlesCmdToJson(bundles []Bundle) (string, error) {
	out, err := json.Marshal(bundles)
	return string(out), err
}
