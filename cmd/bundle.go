/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/apodhrad/iib-cli/utils"
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
	var bundle Bundle

	if len(args) < 3 {
		return errors.New("Specify a csv, pkg and package name!")
	}

	out, err = bundleCmdGrpc(args[0], args[1], args[2])
	if err == nil {
		bundle, err = bundleCmdUnmarshal(out)
		if err == nil {
			if output == "json" {
				out, err = bundleCmdToJson(bundle)
			} else {
				out, err = bundleCmdToText(bundle)
			}
			if err == nil {
				fmt.Println(out)
			}
		}
	}
	return err
}

func bundleCmdGrpc(csv string, pkg string, channel string) (string, error) {
	var err error
	var out string

	err = utils.GrpcStartSafely()
	if err == nil {
		method := "api.Registry/GetBundle"
		data := fmt.Sprintf(`{"csvName":"%s","pkgName":"%s","channelName":"%s"}`, csv, pkg, channel)
		grpcArg := utils.GrpcArgMethodWithData(method, data)
		out, err = utils.GrpcExec(grpcArg)
	}
	utils.GrpcStopSafely()
	return out, err
}

func bundleCmdUnmarshal(input string) (Bundle, error) {
	var bundle Bundle
	err := json.Unmarshal([]byte(input), &bundle)
	return bundle, err
}

func bundleCmdToText(bundle Bundle) (string, error) {
	tbl := NewTable("Csv", "Package", "Channel")
	tbl.AddRow(bundle.CsvName, bundle.PackageName, bundle.ChannelName)
	return TableToString(tbl), nil
}

func bundleCmdToJson(bundle Bundle) (string, error) {
	out, err := json.Marshal(bundle)
	return string(out), err
}
