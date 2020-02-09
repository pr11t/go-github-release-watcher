package cmd

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

// YAML config
type repoConfig struct {
	Owner         string
	Name          string
	AssetPattern  string
	ChecksumAsset string
	TargetDir     string
}

type configFile struct {
	Repositories      []repoConfig
	WaitBetweenChecks int64
}

func LoadConfigFile(file string) *configFile {
	cfg := configFile{}
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("%v Failed to read file %s", err, file)
		return &cfg
	}
	err = yaml.UnmarshalStrict(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("%v failed to unmarshal", err)
	}
	return &cfg
}

// final config with environment variables

type finalConfig struct {
	GithubAccessToken string
	WaitBetweenChecks int64
	Repositories      []repoConfig
}

func loadFromEnvironment(varname string) string {
	token := os.Getenv(varname)
	if token == "" {
		log.Fatalf("Environment variable %s empty", varname)
	}
	return os.Getenv(varname)
}

func LoadConfig(filename string) *finalConfig {
	yamlconf := LoadConfigFile(filename)
	final := finalConfig{
		Repositories:      yamlconf.Repositories,
		WaitBetweenChecks: yamlconf.WaitBetweenChecks,
		GithubAccessToken: loadFromEnvironment("GITHUB_ACCESS_TOKEN"),
	}
	return &final
}
