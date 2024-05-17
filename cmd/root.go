/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/louvri/gosl-gen/internal/process"
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
func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	cfg := ""
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "initialize gosl powered project",
		Long:  `create & copy all necessary files required by gosl to run in the project`,
		Run: func(cmd *cobra.Command, args []string) {
			if cfg != "" {
				runner := process.New()
				err := runner.Initialize(cfg)
				if err != nil {
					fmt.Printf("initialization failed %v\n", err)
				}
			} else {
				fmt.Println("config is not set")
			}

		},
	}
	var genCmd = &cobra.Command{
		Use:   "gen",
		Short: `generate golang helper library at host project`,
		Long:  `generate model, helper, and query golang modules. Built based on the stored configs`,
		Run: func(cmd *cobra.Command, args []string) {
			runner := process.New()
			if err := runner.IsInitiated(); err == nil {
				err := runner.Generate(cfg)
				if err != nil {
					fmt.Printf("gosl-gen failed %v\n", err)
				} else {
					fmt.Println("gosl is generated")
				}
			} else {
				fmt.Println(err.Error())
			}
		},
	}
	rootCmd.PersistentFlags().StringVarP(&cfg, "config", "c", "", "config file")
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(genCmd)
}
