package intializer

import (
	"github.com/IstarVin/manga-reader-go/global"
	json "github.com/json-iterator/go"
	"log"
	"os"
	"reflect"
)

func LoadConfigFile() {
	// Read the config file
	configFileObj := *global.DefaultConfig
	configFile, err := os.ReadFile(global.ConfigFilePath)
	if err != nil {
		if os.IsExist(err) {
			log.Fatal("Error opening config file")
		}
	} else {
		// Unmarshal the config file
		err = json.Unmarshal(configFile, &configFileObj)
		if err != nil {
			log.Fatal("Error unmarshalling config file")
		}
	}

	// Check if default config and cli args and if there are changes in the config file are same
	if reflect.DeepEqual(global.DefaultConfig, global.Config) &&
		!reflect.DeepEqual(global.DefaultConfig, configFileObj) {
		global.Config = &configFileObj
	}

	if global.Config == nil {
		global.Config = &configFileObj
	}

	updateConf(global.Config)
}

func updateConf(config *global.Configuration) {
	configMarshal, err := json.MarshalIndent(*config, "", "    ")
	if err != nil {
		log.Fatal("Error marshalling config file")
	}

	err = os.WriteFile(global.ConfigFilePath, configMarshal, 0644)
	if err != nil {
		log.Fatal("Error writing config file", err)
	}
}
