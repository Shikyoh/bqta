package cmd

type AutocreatorConfig struct {
	Projects []ProjectConfig
	Slack    struct {
		Name    string
		Channel string
		Webhook string
	}
}

type ProjectConfig struct {
	// Name is a human readable name for this BigQuery dataset/prefix
	Name      string
	ProjectID string `mapstructure:"project_id"`
	BigQuery  BigQueryConfig
}

type BigQueryConfig struct {
	Dataset    string
	Prefix     string
	SchemaPath string `mapstructure:"schema_path"`
}
