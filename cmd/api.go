/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/apodhrad/iib-cli/utils"
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
	RunE: apiCmdRunE,
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

func apiCmdRunE(cmd *cobra.Command, args []string) error {
	var err error
	var out string
	var services []Service

	if len(args) == 0 {
		services, err = apiCmdGetServices()
		if err == nil {
			if output == "json" {
				apiCmdToJson(services)
			} else {
				apiCmdToText(services)
			}
		}
	} else {
		out, err = utils.GrpcExec(utils.GrpcArgApi("describe " + args[0]))
		if err == nil {
			fmt.Println(out)
		}
	}
	return err
}

func apiCmdGetServices() ([]Service, error) {
	var services []Service

	utils.GrpcStartSafely()

	out, err := utils.GrpcExec(utils.GrpcArgApi("list"))
	if err != nil {
		return services, err
	}

	serviceScanner := bufio.NewScanner(strings.NewReader(out))
	for serviceScanner.Scan() {
		serviceName := serviceScanner.Text()
		service := Service{Name: serviceName}

		out, err = utils.GrpcExec(utils.GrpcArgApi("list " + serviceName))
		if err != nil {
			return services, err
		}
		methodScanner := bufio.NewScanner(strings.NewReader(out))
		for methodScanner.Scan() {
			method := methodScanner.Text()
			service.Methods = append(service.Methods, method)
		}
		services = append(services, service)
	}

	utils.GrpcStopSafely()

	return services, err
}

func apiCmdToText(services []Service) (string, error) {
	tbl := NewTable("Service", "Method")
	for _, service := range services {
		for _, method := range service.Methods {
			tbl.AddRow(service.Name, method)
		}
	}
	return TableToString(tbl), nil
}

func apiCmdToJson(services []Service) (string, error) {
	out, err := json.Marshal(services)
	return string(out), err
}
