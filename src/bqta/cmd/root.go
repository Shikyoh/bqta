package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	slack "github.com/ashwanthkumar/slack-go-webhook"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	config AutocreatorConfig
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "bqta",
	Short: "Create BigQuery tables and send notifications to Slack",
	Long: `Example:
	./bin/bqta create all --day today
	
	./bin/bqta create one --day tomorrow --name data2`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func validateDay(day string) error {
	if day != "today" && day != "tomorrow" {
		return errors.New("Day must be [today] or [tomorrow]")
	}
	return nil
}

var (
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create table",
	}
	createOneCmd = &cobra.Command{
		Use:   "one",
		Short: "Create table for a specific project",
		Run: func(cmd *cobra.Command, args []string) {
			if err := validateDay(createCmdDay); err != nil {
				log.Fatal(err.Error())
			}
			if createCmdName == "" {
				log.Fatal("Please specify --name")
			}
			var thisProject *ProjectConfig
			for _, proj := range config.Projects {
				if proj.Name == createCmdName {
					thisProject = &proj
					break
				}
			}
			if thisProject == nil {
				log.Fatalf("Cannot find given project - %s", createCmdName)
			}

			att, err := createTable(*thisProject, createCmdDay)
			if err == nil {
				attachSuccess(att)
			} else {
				attachFailure(att, err.Error())
				log.Println(err)
			}
			slackSend([]slack.Attachment{*att})
		},
	}
	createAllCmd = &cobra.Command{
		Use:   "all",
		Short: "Create table for all projects",
		Run: func(cmd *cobra.Command, args []string) {
			if err := validateDay(createCmdDay); err != nil {
				log.Fatal(err.Error())
			}
			attachments := make([]slack.Attachment, len(config.Projects))

			for i, proj := range config.Projects {
				att, err := createTable(proj, createCmdDay)
				if err == nil {
					attachSuccess(att)
				} else {
					attachFailure(att, err.Error())
					log.Println(err)
				}
				attachments[i] = *att
			}
			slackSend(attachments)
		},
	}
)

var (
	createCmdName, createCmdDay string
)

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bq-table-autocreator.yaml)")

	createOneCmd.Flags().StringVar(&createCmdName, "name", "", "BigQuery bqta project name")
	createOneCmd.Flags().StringVar(&createCmdDay, "day", "", "Table date [today, tomorrow]")
	createAllCmd.Flags().StringVar(&createCmdName, "name", "", "BigQuery bqta project name")
	createAllCmd.Flags().StringVar(&createCmdDay, "day", "", "Table date [today, tomorrow]")

	createCmd.AddCommand(createOneCmd)
	createCmd.AddCommand(createAllCmd)
	RootCmd.AddCommand(createCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("bqta") // name of config file (without extension)
	viper.AddConfigPath("$PWD/configs")
	viper.AddConfigPath("/etc/bqta")
	viper.AddConfigPath("$HOME")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	viper.Unmarshal(&config)
	setupBigQuery()
}
