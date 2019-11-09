package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"

	gabs "github.com/Jeffail/gabs/v2"
)

func main() {
	configFile := "/kaniko/.docker/config.json"
	// configFile := "test/config.json"
	jsonFile, err := os.OpenFile(configFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &gabs.Container{}
	if len(byteValue) == 0 {
		cfg, err = gabs.ParseJSON([]byte(`{}`))
	} else {
		cfg, err = gabs.ParseJSON(byteValue)
	}
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range os.Environ() {
		if strings.HasPrefix(v, "ECR_LOGIN_") {
			s := strings.Split(v, "=")
			cfg.Set("ecr-login", "credHelpers", s[1])
		}
	}

	// Docker Auth Configuration
	for _, v := range os.Environ() {
		if strings.HasPrefix(v, "DKR_AUTH_") {
			key := strings.Split(v, "=")[0]
			cnf := strings.Split(key, "__")
			if len(cnf) == 2 {
				if subkey := cnf[1]; len(subkey) != 0 {
					if repo := os.Getenv(cnf[0]); len(repo) != 0 {
						cfg.Set(getValue(key), "auths", repo, strings.ToLower(subkey))
					} else {
						log.Printf("WARN: Unable to find repo for %s", key)
						continue
					}
				}
			}
		}
	}

	// Proxy Configuration
	if v := os.Getenv("DKRCFG_PROXY"); len(v) != 0 {
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

	jsonFile.Truncate(0)
	jsonFile.Seek(0, 0)
	jsonFile.Write([]byte(cfgPretty))

	if v := os.Getenv("DKRCFG_DEBUG"); len(v) != 0 {
		log.Printf("DEBUG: Docker Config: %s\n", configFile)
		log.Println(cfgPretty)
	}
}

func getValue(key string) string {
	val := os.Getenv(key)
	if v := os.Getenv("DKRCFG_ENABLE_AWS_PSTORE"); len(v) != 0 {
		if strings.HasPrefix(val, "arn:aws:ssm:") {
			assumeRole := &AssumeRoleConfig{
				Profile:         getEnv(key+"__PROFILE", ""),
				RoleARN:         getEnv(key+"__ROLE_ARN", ""),
				ExternalID:      getEnv(key+"__EXTERNAL_ID", ""),
				RoleSessionName: getEnv(key+"__SESSION_NAME", ""),
			}
			return getParameter(val, assumeRole)
		}
	}
	return val
}

func getParameter(key string, roleCfg *AssumeRoleConfig) (val string) {
	// Marshal Request
	prm := strings.Split(key, ":parameter")[1]
	region := strings.Split(key, ":")[3]

	// AWS Session
	sess_opts := session.Options{
		Config: *aws.NewConfig().WithRegion(region),
	}
	if roleCfg.Profile != "" {
		sess_opts.Profile = roleCfg.Profile
	}
	sess := session.Must(session.NewSessionWithOptions(sess_opts))

	creds := &credentials.Credentials{}
	if roleCfg.RoleARN != "" {
		// Get AssumeRole Credentials
		creds = stscreds.NewCredentials(sess, roleCfg.RoleARN, func(p *stscreds.AssumeRoleProvider) {
			if roleCfg.RoleSessionName != "" {
				p.RoleSessionName = roleCfg.RoleSessionName
			} else {
				p.RoleSessionName = fmt.Sprintf("kaniko-gitlab-%d", time.Now().UTC().UnixNano())
			}
			if roleCfg.ExternalID != "" {
				p.ExternalID = aws.String(roleCfg.ExternalID)
			}
		})
	} else {
		creds = sess.Config.Credentials
	}

	// SSM Client
	ssmcfg := &aws.Config{
		Credentials: creds,
	}
	ssmclient := ssm.New(sess, ssmcfg)
	resp, err := ssmclient.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(prm),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		log.Fatalf("ERROR: ssm.GetParameter:: %s\n%s", key, err)
	}
	val = *resp.Parameter.Value
	return
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

type DockerConfig struct {
	Auths       map[string]interface{} `json:"auths"`
	CredHelpers map[string]interface{} `json:"credHelpers"`
}

type AssumeRoleConfig struct {
	Profile         string
	RoleARN         string
	ExternalID      string
	RoleSessionName string
}
