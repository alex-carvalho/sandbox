Jenkins DSL that that deploy a simple Spring Boot app (build, run tests, run sonnar and deploy in EC2) using terraform

Using localstack and tflocal

```shell
pip install terraform-local # install tflocal
docker-compose up -d # Run localstack and Jenkins
```

Java app code: https://github.com/alex-carvalho/sandbox/tree/master/samples/java/spring-web-3-j21


AWS CLI for localstack
```
aws configure set aws_access_key_id test
aws configure set aws_secret_access_key test
aws configure set region us-east-1
aws configure set output json
export AWS_ENDPOINT_URL=http://localhost:4566
```

Create a S3 bucket
aws s3 mb s3://my-artifacts


## Required to run
- Create in jenkins the secret `sonar-token`

