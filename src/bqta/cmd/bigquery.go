package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/ashwanthkumar/slack-go-webhook"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/bigquery/v2"
)

var (
	bqService *bigquery.Service

	colourSuccess = "#00AA00"
	colourFailure = "#AA0000"
)

func newTableDefinition(project ProjectConfig, day string) (*bigquery.Table, error) {
	if err := validateDay(createCmdDay); err != nil {
		return nil, err
	}
	date := time.Now()
	if day == "tomorrow" {
		date = date.AddDate(0, 0, 1)
	}

	tableID := project.BigQuery.Prefix + date.Format("20060102")

	schema := readSchema(project)

	return &bigquery.Table{
		Description: "Created by bq-table-autocreator",
		TableReference: &bigquery.TableReference{
			DatasetId: project.BigQuery.Dataset,
			ProjectId: project.ProjectID,

			TableId: tableID,
		},

		Schema: schema,
	}, nil
}

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

func newAttachment(table *bigquery.Table) (att slack.Attachment) {
	return
}

func attachSuccess(att *slack.Attachment) {
	att.Color = &colourSuccess
}

func attachFailure(att *slack.Attachment, errText string) {
	att.Color = &colourFailure
	att.Title = &errText
}

func slackSend(attachments []slack.Attachment) {
	payload := slack.Payload("", config.Slack.Name,
		"",
		config.Slack.Channel,
		attachments)

	slack.Send(config.Slack.Webhook, "", payload)
}

func createTable(project ProjectConfig, day string) (*slack.Attachment, error) {
	table, err := newTableDefinition(project, day)
	if err != nil {
		return nil, err
	}
	att := &slack.Attachment{}
	att.AddField(slack.Field{Title: "Name", Value: project.Name, Short: true})
	att.AddField(slack.Field{Title: "Table", Value: table.TableReference.TableId, Short: true})
	att.AddField(slack.Field{Title: "Dataset", Value: table.TableReference.DatasetId, Short: true})
	att.AddField(slack.Field{Title: "Project", Value: table.TableReference.ProjectId, Short: true})

	tableService := bigquery.NewTablesService(bqService)
	if err != nil {
		return att, err
	}
	call := tableService.Insert(project.ProjectID, project.BigQuery.Dataset, table)
	_, err = call.Do()
	if err != nil {
		return att, err
	}
	return att, nil
}

func setupBigQuery() {
	client, err := google.DefaultClient(oauth2.NoContext, bigquery.BigqueryInsertdataScope)
	if err != nil {
		log.Fatal(err)
	}

	bqService, err = bigquery.New(client)
	if err != nil {
		log.Fatal(err)
	}
}
