package main

type AutocreatorConfig struct {
	Projects []ProjectConfig `yaml:"projects"`
	Slack    struct {
		Name    string `yaml:"name"`
		Channel string `yaml:"channel"`
		Webhook string `yaml:"webhook"`
	} `yaml:"slack"`
}

type ProjectConfig struct {
	ProjectID string         `yaml:"project_id"`
	BigQuery  BigQueryConfig `yaml:"bigquery"`
}

type BigQueryConfig struct {
	Dataset    string `yaml:"dataset"`
	Prefix     string `yaml:"prefix"`
	SchemaPath string `yaml:"schema_path"`
}
