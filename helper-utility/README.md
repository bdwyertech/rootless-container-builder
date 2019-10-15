# Overview
This is a utility to programatically build a Docker config file.

### ECR-Helper
Create an environment variable with the prefix `ECR_LOGIN_`
```shell
ECR_LOGIN_MY_ACCOUNT__US_EAST1 = "123456789876.dkr.ecr.us-east-1.amazonaws.com"
ECR_LOGIN_MY_ACCOUNT__US_WEST2 = "123456789876.dkr.ecr.us-west-2.amazonaws.com"
ECR_LOGIN_ANOTHERR_1__US_EAST1 = "567898765432.dkr.ecr.us-east-1.amazonaws.com"
```
