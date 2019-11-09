# Overview
This is a utility to programatically build a Docker config file.

### ECR-Helper
Create an environment variable with the prefix `ECR_LOGIN_`
```bash
ECR_LOGIN_MY_ACCOUNT__US_EAST1="123456789876.dkr.ecr.us-east-1.amazonaws.com"
ECR_LOGIN_MY_ACCOUNT__US_WEST2="123456789876.dkr.ecr.us-west-2.amazonaws.com"
ECR_LOGIN_ANOTHERR_1__US_EAST1="567898765432.dkr.ecr.us-east-1.amazonaws.com"
```

### Registry Credentials
```bash
DKR_AUTH_NEXUS="https://nexus.myorg.net"
DKR_AUTH_NEXUS__USER="nexus_user"
DKR_AUTH_NEXUS__PASS="nexus_pass"
DKR_AUTH_GITLAB="https://gitlab.myorg.net"
DKR_AUTH_GITLAB__AUTH="gitlab_token"
```
You can pull values from AWS SSM ParameterStore by setting `KCFG_ENABLE_AWS_PSTORE` and ensuring you pass a valid Parameter ARN
```bash
DKRCFG_ENABLE_AWS_PSTORE=1
DKR_AUTH_REPO=https://repo.myorg.net
# Base Config
DKR_AUTH_REPO____AWS_PSTORE_ROLE_ARN=arn:aws:iam::567898765432:role/my_parameterstore_role
DKR_AUTH_REPO__USER=arn:aws:ssm:us-east-1::parameter/path/to/credential/user
DKR_AUTH_REPO__PASS=arn:aws:ssm:us-east-1::parameter/path/to/credential/pass
# Key-Specific Config
DKR_AUTH_REPO__AUTH=arn:aws:ssm:us-east-1::parameter/path/to/credential/auth
DKR_AUTH_REPO__AUTH__AWS_PSTORE_ROLE_ARN=arn:aws:iam::567898765432:role/my_other_pstore_role
```

### Proxy Configuration
You can inherit proxy configuration into the Docker configuration
```bash
export DKRCFG_PROXY=1
```
