// config 配置相关包
package config


import (
    "io/ioutil"
    "log"
    "os"
    "gopkg.in/yaml.v3"
)


type DBConfig struct {
    DSN      string       `yaml:"dsn"`
    Driver   string       `yaml:"driver"`
}


type PKGConfig struct {
    DB  DBConfig    `yaml:"db"`
}


var Config = PKGConfig{}


func init() {
    file, err := os.Open("pkg_config.yaml")
	if err != nil {
		log.Fatalf("config.init, fail to parse 'pkg_config.yaml': %v", err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("fail to read 'pkg_config.yaml': %v", err)
	}
	if err := yaml.Unmarshal([]byte(content), &Config); err != nil {
		log.Fatalf("fail to yaml Unmarshal config: %v", err)
	}
}
