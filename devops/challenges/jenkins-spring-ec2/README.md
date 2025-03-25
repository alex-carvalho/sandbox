Jenkins DSL that that deploy a simple Spring Boot app (build, run tests, run sonnar and deploy in EC2) using terraform

Using localstack and tflocal

```shell
pip install terraform-local # install tflocal
docker-compose up -d # Run localstack and Jenkins
```

Java app code: https://github.com/alex-carvalho/sandbox/tree/master/samples/java/spring-web-3-j21



## Required to run
- Create in jenkins the secret `sonar-token`

