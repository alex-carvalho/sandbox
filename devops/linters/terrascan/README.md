## terrascan

https://runterrascan.io/docs/getting-started


```
-c  <config file path>
-c  <config file path>
-t	Use this to specify individual cloud providers
-i  specify the provider (k8s, heml, kustomize, docker)
-o  output format
--find-vuln can display vulnerabilities for container images present

terrascan scan -i k8s

# remote scan
terrascan scan -t aws -r git -u git@github.com:tenable/KaiMonkey.git//terraform/aws

# scan terraform file
terrascan  scan -i terraform -t aws main.tf -o json

# scan with configuration
terrascan  scan -i terraform -t aws -c terrascan-config.toml main.tf -o json

```

### server

```
terrascan server

# endpoints
GET /health

POST - /v1/{iac}/{iacVersion}/{cloud}/local/file/scan
ex: url -i -F "file=@aws_cloudfront_distribution.tf" localhost:9010/v1/terraform/v14/aws/local/file/scan


```
