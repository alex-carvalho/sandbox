# Floci https://github.com/floci-io/floci

Builded in java with Quarkus

```shell
docker run -p 4566:4566 floci/floci:latest
```

Test in terminal:
```shell
export AWS_ENDPOINT_URL=http://localhost:4566

aws s3 ls

aws s3api create-bucket --bucket my-bucket-local 

aws s3 ls

echo "hello" > my-file.txt   

aws s3api put-object --bucket my-bucket-local --key myfolder/my-file.txt --body my-file.txt

aws s3api get-object --bucket my-bucket-local --key myfolder/my-file.txt downloaded-file.txt

cat downloaded-file.txt
```

Testing terraform:
```shell
terraform init
terraform plan
terraform apply
terraform plan
terraform destroy
```
