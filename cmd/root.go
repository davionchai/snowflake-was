package cmd

import (
	"log"
	"os"

	"github.com/davionchai/snowflake-was/cmd/was"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var yamlFile string

var rootCmd = &cobra.Command{
	Use:   "snowflake-was",
	Short: "Caller for DE snowflake-was tools",
	Long:  `Caller for Data Engineering snowflake warehouse auto scaling tools.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && viper.ConfigFileUsed() == "" {
			cmd.Help()
		}
	},
	DisableFlagsInUseLine: true,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("There was an error while executing CLI [%s]", err)
	}
}

func init() {
	cobra.OnInitialize(getYaml)

	// subcommands tree here
	rootCmd.AddCommand(was.WasCmd)

	// global cli that can be shared across subcommands
	rootCmd.PersistentFlags().StringVarP(
		&yamlFile,
		"file",
		"f",
		"",
		"config yaml file (default is $(pwd)/config.yaml)",
	)
}

func getYaml() {
	if yamlFile != "" {
		// use config file from the flag
		viper.SetConfigFile(yamlFile)
	} else {
		// find relative directory
		currentPath, err := os.Getwd()
		cobra.CheckErr(err)

		// search config in relative path with name "./config.yaml"
		viper.AddConfigPath(currentPath)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// if a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Found config file: %s", viper.ConfigFileUsed())
	}
}
