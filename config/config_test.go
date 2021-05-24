package ART

import (
	"testing"
	"reflect"
)

// Test configuration reading and parsing
func TestParseConfig(t *testing.T){
	expected := &Config {
		DbUser: "username",
		DbPassword: "password",
		DbHost: "host",
		TimeWindow: 60,
		BearerToken: "token",
	}

	configFile := "artconfig_test.yaml"
	config, err := GetConfig(configFile)
	if err != nil{
		t.Errorf("Could not get config. Error: %v ", err)
	}

	if !reflect.DeepEqual(expected, config){
		t.Errorf("Parsed config didn't match. Expected: %v\nGot: %v.", expected, config)
	}
}
