package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/connectedventures/f8-pkg/configurator"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/bigquery/v2"
)

func newTableDefinition(project ProjectConfig) *bigquery.Table {
	date := time.Now().AddDate(0, 0, 1).Format("20060102")
	tableID := project.BigQuery.Prefix + date

	schema := readSchema(project)

	return &bigquery.Table{
		FriendlyName: "Automatically created by bq-table-autocreator",
		TableReference: &bigquery.TableReference{
			DatasetId: project.BigQuery.Dataset,
			ProjectId: project.ProjectID,

			TableId: tableID,
		},

		Schema: schema,
	}
}

var (
	config AutocreatorConfig
)

func readSchema(project ProjectConfig) *bigquery.TableSchema {
	var schema bigquery.TableSchema
	// Load the schema from disk.
	f, err := ioutil.ReadFile(project.BigQuery.SchemaPath)
	if err != nil {
		log.Fatal("Could not read BigQuery table schema: " + err.Error())
	}
	if err = json.Unmarshal(f, &schema); err != nil {
		log.Fatal("Could not parse BigQuery table schema: " + err.Error())
	}

	return &schema
}

func slackSend(text string, err error) {
	att := slack.Attachment{}
	if err == nil {
		att.AddField(slack.Field{Title: "Status", Value: "Success"})
	} else {
		att.AddField(slack.Field{Title: "Status", Value: "Failure"})
		att.AddField(slack.Field{Title: "Error", Value: err.Error()})
	}
	payload := slack.Payload(text, config.Slack.Name,
		"",
		config.Slack.Channel,
		[]slack.Attachment{att})

	slack.Send(config.Slack.Webhook, "", payload)
}

func main() {
	err := configurator.Parse("/etc/bq-table-autocreator/config.yaml", &config)
	if err != nil {
		log.Fatal(err)
	}

	client, err := google.DefaultClient(oauth2.NoContext, bigquery.BigqueryInsertdataScope)
	if err != nil {
		log.Fatal(err)
	}

	service, err := bigquery.New(client)
	if err != nil {
		log.Fatal(err)
	}

	for _, project := range config.Projects {
		tableService := bigquery.NewTablesService(service)
		table := newTableDefinition(project)
		call := tableService.Insert(project.ProjectID, project.BigQuery.Dataset, table)
		table, err := call.Do()

		if err != nil {
			log.Println(err)
			slackSend("Could not create table :here", err)
		} else {
			slackSend("Successfully created table", nil)
		}
	}
}
