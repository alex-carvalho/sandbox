## checkov

Python base policy-as-code

- Suppress or skip policies
- Scan credentials and secrets
- Scan Kubernetes clusters
- Scan Terraform plan output and 3rd party modules
- Many output formats


```
pip3 install checkov
```

```
checkov -o json --file main.tf
```

It supports junit_xml format, so it can disply like Jenkins failed itens
