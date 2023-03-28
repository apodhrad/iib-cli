/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/apodhrad/iib-cli/utils"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

var output string

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("Cmd ", cmd)
		// fmt.Println("Args ", args)
		// fmt.Println("Output ", output)
		if len(args) == 0 {
			table, json, err := listApi()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			if output == "json" {
				fmt.Println(json)
			} else {
				table.Print()
			}
		} else {
			out, err := utils.GrpcExec(utils.GrpcArgApi("describe " + args[0]))
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			fmt.Println(out)
		}
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// apiCmd.PersistentFlags().String("foo", "", "A help for foo")
	apiCmd.PersistentFlags().StringVarP(&output, "output", "o", "text", "Output format")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// apiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type Service struct {
	Name    string
	Methods []string
}

func listApi() (tblOut table.Table, jsonOut string, err error) {
	out, err := utils.GrpcExec(utils.GrpcArgApi("list"))
	if err != nil {
		return nil, "", err
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
			return nil, "", err
		}
		methodScanner := bufio.NewScanner(strings.NewReader(out2))
		for methodScanner.Scan() {
			method := methodScanner.Text()
			service.Methods = append(service.Methods, method)
			tbl.AddRow(serviceName, method)
		}
		services = append(services, service)
	}

	json, err := json.Marshal(services)
	return tbl, string(json), err
}
