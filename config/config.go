package ART

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct{
	DbUser string `yaml:"db_user"`
	DbPassword string `yaml:"db_password"`
	DbHost string `yaml:"db_host"`
	TimeWindow int `yaml:"default_time_window_minutes"`
	BearerToken string `yaml:"twitter_bearer_token"`
}
 
func GetConfig(configFile string) (*Config, error){
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to read %s. ", configFile)
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil{
		return nil, err
	}

	return config, nil	
}
