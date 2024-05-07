/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gosl-gen",
	Short: "query and model generator built on top of gosl",
	Long:  `gosl-gen is a query helper that tries to mimic orm capabilities but built using generator methodology `,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "initialize gosl powered project",
		Long:  `create & copy all necessary files required by gosl to run in the project`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	var cfgCmd = &cobra.Command{
		Use:   "cfg",
		Short: `set gosl config file`,
		Long:  `set gosl config file, such as workdir, dbConnection ,dbType, dbSchema, dbIncludeTables, dbExcludeTables`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	var genCmd = &cobra.Command{
		Use:   "gen",
		Short: `generate golang helper library at host project`,
		Long:  `generate model, helper, and query golang modules. Built based on the stored configs`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(cfgCmd)
	rootCmd.AddCommand(genCmd)
	rootCmd.MarkFlagRequired("cfg")
}
