package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	gabs "github.com/Jeffail/gabs/v2"
)

func main() {
	configFile := "/kaniko/.docker/config.json"
	// configFile := "../docker-manifest/config.json"
	jsonFile, err := os.Open(configFile)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	cfg, err := gabs.ParseJSON(byteValue)

	for _, v := range os.Environ() {
		if strings.HasPrefix(v, "ECR_LOGIN_") {
			s := strings.Split(v, "=")
			cfg.Set("ecr-login", "credHelpers", s[1])
		}
	}

	cfgPretty := cfg.StringIndent("", "  ")

	ioutil.WriteFile(configFile, []byte(cfgPretty), 0644)

	// DEBUG
	fmt.Printf("DEBUG: Docker Config: %s\n", configFile)
	fmt.Println(cfgPretty)
}

type DockerConfig struct {
	Auths       map[string]interface{} `json:"auths"`
	CredHelpers map[string]interface{} `json:"credHelpers"`
}
