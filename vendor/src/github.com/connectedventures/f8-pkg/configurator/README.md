# Configurator
This package is used to load YAML configuration files and parse them into
structs.

## Example

```go
import "github.com/connectedventures/f8-pkg/configurator"

type cloudConfig struct {
	ProjectID string       `yaml:"project_id"`
}

type loggingConfig struct {
	Level string `yaml:"level"`
}

type myserviceConfig struct {
	Cloud   cloudConfig   `yaml:"cloud"`
	Logging loggingConfig `yaml:"logging"`
}

var config myserviceConfig


func main() {
  err := configurator.Parse("/etc/f8-uploader/uploader.yaml", &config)
  if err != nil {
    log.Fatal(err)
  }
  log.Println(config.Cloud.ProjectID)
}
```
