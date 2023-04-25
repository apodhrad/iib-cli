/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/apodhrad/iib-cli/grpc"
	"github.com/spf13/cobra"
)

// packageCmd represents the package command
var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "List a specific package",
	Long:  `List a specific package`,
	RunE:  packageRunE,
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

func packageRunE(cmd *cobra.Command, args []string) error {
	var err error
	var out string
	var pkg Package

	if len(args) == 0 {
		return errors.New("Specify a package name!")
	}

	name := args[0]
	out, err = packageCmdGrpc(name)
	if err == nil {
		pkg, err = packageCmdUnmarshal(out)
		if err == nil {
			if output == "json" {
				out, err = packageToJson(pkg)
			} else {
				out, err = packageToText(pkg)
			}
			if err == nil {
				fmt.Println(out)
			}
		}
	}

	return err
}

func packageCmdGrpc(name string) (string, error) {
	var err error
	var out string

	grpc.GrpcStart()
	method := "api.Registry/GetPackage"
	data := fmt.Sprintf(`{"name":"%s"}`, name)
	grpcArg := grpc.GrpcArgMethodWithData(method, data)
	out, err = grpc.GrpcExec(grpcArg)
	grpc.GrpcStop()
	return out, err
}

func packageCmdUnmarshal(input string) (Package, error) {
	var pkg Package

	err := json.Unmarshal([]byte(input), &pkg)

	return pkg, err
}

func packageToText(pkg Package) (string, error) {
	tbl := NewTable("Package", "Channel", "Csv", "Default")
	for _, value := range pkg.Channels {
		isDefault := ""
		if value.Name == pkg.DefaultChannelName {
			isDefault = "true"
		}
		tbl.AddRow(pkg.Name, value.Name, value.CsvName, isDefault)
	}
	return TableToString(tbl), nil
}

func packageToJson(pkg Package) (string, error) {
	out, err := json.Marshal(pkg)
	return string(out), err
}
