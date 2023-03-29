/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/apodhrad/iib-cli/utils"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Explore API inside the IIB",
	Long: `Explore API inside the IIB (i.e. gRPC services, methods and params)

Examples:
  iib-cli api
  iib-cli api api.Registry/ListPackages
`,
	RunE: apiRunE,
}

func apiRunE(cmd *cobra.Command, args []string) error {
	err := utils.GrpcStartSafely()
	if err != nil {
		return exitE(err)
	}
	if len(args) == 0 {
		table, json, err := listApi()
		if err != nil {
			return exitE(err)
		}
		if output == "json" {
			fmt.Println(json)
		} else {
			fmt.Println(table)
		}
	} else {
		out, err := utils.GrpcExec(utils.GrpcArgApi("describe " + args[0]))
		if err != nil {
			exitE(err)
		}
		fmt.Println(out)
	}
	return exitE(nil)
}

func init() {
	rootCmd.AddCommand(apiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// apiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// apiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type Service struct {
	Name    string
	Methods []string
}

func listApi() (tblOut string, jsonOut string, err error) {
	out, err := utils.GrpcExec(utils.GrpcArgApi("list"))
	if err != nil {
		return "", "", err
	}

	table.DefaultHeaderFormatter = func(format string, vals ...interface{}) string {
		return strings.ToUpper(fmt.Sprintf(format, vals...))
	}

	tbl := table.New("Service", "Method")
	services := []Service{}

	serviceScanner := bufio.NewScanner(strings.NewReader(out))
	for serviceScanner.Scan() {
		serviceName := serviceScanner.Text()
		service := Service{Name: serviceName}

		out2, err := utils.GrpcExec(utils.GrpcArgApi("list " + serviceName))
		if err != nil {
			return "", "", err
		}
		methodScanner := bufio.NewScanner(strings.NewReader(out2))
		for methodScanner.Scan() {
			method := methodScanner.Text()
			service.Methods = append(service.Methods, method)
			tbl.AddRow(serviceName, method)
		}
		services = append(services, service)
	}

	// b := new(bytes.Buffer)
	var tblBuf bytes.Buffer
	tbl.WithWriter(&tblBuf)
	tbl.Print()

	json, err := json.Marshal(services)
	return tblBuf.String(), string(json), err
}
