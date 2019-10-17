package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	gabs "github.com/Jeffail/gabs/v2"
)

func main() {
	configFile := "/kaniko/.docker/config.json"
	// configFile := "test/config.json"
	jsonFile, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
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

	// Docker Auth Configuration
	for _, v := range os.Environ() {
		if strings.HasPrefix(v, "DKR_AUTH_") {
			s := strings.Split(v, "=")
			cnf := strings.Split(s[0], "__")
			val := s[1]
			if len(cnf) == 2 {
				if key := cnf[1]; len(key) != 0 {
					if repo := os.Getenv(cnf[0]); len(repo) != 0 {
						cfg.Set(val, "auths", repo, strings.ToLower(key))
					} else {
						log.Printf("WARN: Unable to find repo for %s", s[0])
						continue
					}
				}
			}
		}
	}

	// Proxy Configuration
	if v := os.Getenv("KCFG_PROXY"); len(v) != 0 {
		if v := os.Getenv("http_proxy"); len(v) != 0 {
			cfg.SetP(v, "proxies.default.httpProxy")
		}
		if v := os.Getenv("https_proxy"); len(v) != 0 {
			cfg.SetP(v, "proxies.default.httpsProxy")
		}
		if v := os.Getenv("no_proxy"); len(v) != 0 {
			cfg.SetP(v, "proxies.default.noProxy")
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
