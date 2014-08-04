package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

type Config struct {
	ConcurentRequest int             `json:"concurentRequest,omitempty"`
	Retries          int             `json:"retries,omitempty"`
	EtcdConfigPath   string          `json:"etcdConfigPath,omitempty"`
	DumpFilePath     string          `json:"dumpFilePath,omitempty"`
	BackupStrategy   *BackupStrategy `json:"backupStrategy,omitempty"`
	LogFatal         func(v ...interface{})
	LogPrintln       func(v ...interface{})
}

type BackupStrategy struct {
	Keys      []string `json:"keys,omitempty"`
	Sorted    bool     `json:"sorted,omitempty"`
	Recursive bool     `json:"recursive,omitempty"`
}

func (config *Config) ToString() string {
	stringVersion := "ConcurentRequest: " + fmt.Sprintf("%#v", config.ConcurentRequest)
	stringVersion += ", Retries: " + fmt.Sprintf("%#v", config.Retries)
	stringVersion += ", EtcdConfigPath: " + config.EtcdConfigPath
	stringVersion += ", DumpFilePath: " + config.DumpFilePath
	stringVersion += ", BackupStrategy: " + fmt.Sprintf("%#v", config.BackupStrategy)

	return stringVersion
}

var (
	config *Config
)

func init() {
	config = LoadConfig(parseCommandLineOptions())
	config.LogPrintln("Current configuration: ", config.ToString())
}

func LoadConfig(configPath *string, commandLineConfig *Config) *Config {
	currentConfig := loadConfigFile(configPath)
	currentConfig.LogPrintln = func(v ...interface{}) { log.Println(v...) }
	currentConfig.LogFatal = func(v ...interface{}) { log.Fatal(v...) }

	configNilValueOverride(currentConfig, commandLineConfig)
	return currentConfig
}

func parseCommandLineOptions() (*string, *Config) {
	config := Config{}
	configPath := flag.String("config", "backup-configuration.json", "Location of the backup configuration file")
	config.ConcurentRequest = *flag.Int("concurent-request", 10, "Maximum limit of goroutines talking to etcd at a time")
	config.Retries = *flag.Int("retries", 5, "Number of retries before the program give up on failing request")
	config.EtcdConfigPath = *flag.String("etcd-config", "etcd-configuration.json", "Location of the etcd config file")
	config.DumpFilePath = *flag.String("dump", "etcd-dump.json", "Location of the etcd dump file")
	config.BackupStrategy = &BackupStrategy{[]string{"/"}, true, true}

	flag.Parse()
	return configPath, &config
}

func loadConfigFile(configPath *string) *Config {
	file, error := os.Open(*configPath)
	defer file.Close()
	if error != nil {
		config.LogPrintln("Default options: ")
		flag.PrintDefaults()
		config.LogFatal("Error when trying to open the configuration file `"+*configPath+"`. Error: ", error)
	}

	currentConfig := &Config{}
	jsonParser := json.NewDecoder(file)
	if err := jsonParser.Decode(currentConfig); err != nil {
		config.LogFatal("Error when trying to load config file set into json. Error: ", err)
	}

	return currentConfig
}

func configNilValueOverride(currentConfig *Config, defaultValue *Config) {

	if currentConfig.ConcurentRequest == 0 {
		currentConfig.ConcurentRequest = defaultValue.ConcurentRequest
	}

	if currentConfig.ConcurentRequest == 0 {
		currentConfig.ConcurentRequest = defaultValue.ConcurentRequest
	}

	if currentConfig.Retries == 0 {
		currentConfig.Retries = defaultValue.Retries
	}

	if currentConfig.EtcdConfigPath == "" {
		currentConfig.EtcdConfigPath = defaultValue.EtcdConfigPath
	}

	if currentConfig.DumpFilePath == "" {
		currentConfig.DumpFilePath = defaultValue.DumpFilePath
	}

	if currentConfig.BackupStrategy == nil {
		currentConfig.BackupStrategy = defaultValue.BackupStrategy
	}

	if currentConfig.LogFatal == nil {
		currentConfig.LogFatal = defaultValue.LogFatal
	}

	if currentConfig.LogPrintln == nil {
		currentConfig.LogPrintln = defaultValue.LogPrintln
	}
}
