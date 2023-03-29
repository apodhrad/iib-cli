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

// packagesCmd represents the packages command
var packagesCmd = &cobra.Command{
	Use:   "packages",
	Short: "List all packages",
	Long:  `List all packages`,
	RunE:  packagesRunE,
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

func packagesExitE(err error) error {
	utils.GrpcStop()
	return err
}

type Package struct {
	Name string `json:"name"`
}

func packagesToText(packages []Package) (string, error) {
	tbl := NewTable("Package")
	for _, value := range packages {
		tbl.AddRow(value.Name)
	}
	return TableToString(tbl), nil
}

func packagesToJson(packages []Package) (string, error) {
	out, err := json.Marshal(packages)
	return string(out), err
}

func getPackages(input string) ([]Package, error) {
	// grpcurl produces improper json array
	// so, we will fix it to a proper json array
	numOfPackages := strings.Count(input, "}")
	input = strings.Replace(input, "}", "},", numOfPackages-1)
	input = "[" + input + "]"
	// now, we can unmarshal
	var packages []Package
	err := json.Unmarshal([]byte(input), &packages)
	return packages, err
}

func packagesRunE(cmd *cobra.Command, args []string) error {
	var err error
	var out string

	err = utils.GrpcStartSafely()
	if err == nil {
		grpcArg := utils.GrpcArgMethod("api.Registry/ListPackages")
		out, err = utils.GrpcExec(grpcArg)
		if err == nil {
			var packages []Package
			packages, err = getPackages(out)
			if err == nil {
				var result string
				if output == "json" {
					result, err = packagesToJson(packages)
				} else {
					result, err = packagesToText(packages)
				}
				if err == nil {
					fmt.Println(result)
				}
			}
		}
	}
	return exitE(err)
}
