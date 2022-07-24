// config is a web configuration package
package config


import (
    "os"
    "log"
    "io/ioutil"
    "gopkg.in/yaml.v3"
)


type Config struct {
    Mode          string   `yaml:"mode"`
    Port          int      `yaml:"port"`
    FileBaseUri   string   `yaml:"fileBaseUri"`
}


var GlobaConfig = Config{}


// SetupConfig is an initialization web configuration method
func SetupConfig() {
    file, err := os.Open("config.yaml")
    if err != nil {
        log.Fatalf("config setup fail to parse 'config.yaml': %v", err)
    }
    content, err := ioutil.ReadAll(file)
    if err != nil {
        log.Fatalf("failed to read 'config.yaml': %v", err)
    }
    if err := yaml.Unmarshal([]byte(content), &GlobaConfig); err != nil {
        log.Fatalf("failed to yaml Unmarshal err: %v", err)
    }
}
