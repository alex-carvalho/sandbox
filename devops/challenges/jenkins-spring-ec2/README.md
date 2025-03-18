Jenkins DSL that that deploy a simple Spring Boot app (build, run tests, run sonnar and deploy in EC2) using terraform

Using localstack and tflocal

install terraform
```shell
wget -O - https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
sudo apt update && sudo apt install terraform
```

Install tflocal

```shell
pip install terraform-local
```

Run localstack and Jenkins

```shell
docker-compose up
```
