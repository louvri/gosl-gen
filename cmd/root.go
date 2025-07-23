/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/louvri/gosl-gen/internal/process"
	"github.com/spf13/cobra"
)

const VERSION string = "v0.2.16"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gosl-gen",
	Short: "query and model generator built on top of gosl",
	Long:  `gosl-gen is a query helper that tries to mimic orm capabilities but built using generator methodology `,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	cfg := ""
	var genCmd = &cobra.Command{
		Use:   "compile",
		Short: `generate golang helper library at host project`,
		Long:  `generate model, helper, and query golang modules. Built based on the stored configs`,
		Run: func(cmd *cobra.Command, args []string) {
			runner := process.New()
			if err := runner.Initialize(cfg); err != nil {
				fmt.Printf("gosl-gen init failed %v\n", err)
			} else {
				fmt.Println("gosl is initiated")
			}
			if err := runner.Generate(cfg); err != nil {
				fmt.Printf("gosl-gen gen failed %v\n", err)
			} else {
				fmt.Println("gosl is generated")
			}

		},
	}
	rootCmd.PersistentFlags().StringVarP(&cfg, "config", "c", "", "config file")
	rootCmd.AddCommand(genCmd)
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: `gosl-gen version`,
		Long:  `gosl-gen version`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(VERSION)
		},
	}
	rootCmd.AddCommand(versionCmd)
}
