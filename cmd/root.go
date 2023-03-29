/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/apodhrad/iib-cli/utils"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "iib-cli",
	Short: "Index Image Bundle (IIB) cli tool",
	Long: `Index Image Bundle (IIB) cli tool allows you to explore APIs, packages and bundles in IIBs.

Examples:
  iib-cli api
  iib-cli get packages
  iib-cli get bundles
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var output string

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.iib-cli.yaml)")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "text", "Output format")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func exitE(err error) error {
	utils.GrpcStopSafely()
	return err
}

func NewTable(headers ...string) table.Table {
	table.DefaultHeaderFormatter = func(format string, vals ...interface{}) string {
		return strings.ToUpper(fmt.Sprintf(format, vals...))
	}
	newHeaders := make([]interface{}, len(headers))
	for i, v := range headers {
		newHeaders[i] = v
	}
	return table.New(newHeaders...)
}

func TableToString(tbl table.Table) string {
	var tblBuf bytes.Buffer
	tbl.WithWriter(&tblBuf)
	tbl.Print()
	return tblBuf.String()
}
