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

```