# Live Demo: RobotoCore (`Terraform` + `AWS CLI`)

> [← Back to Presentation (`README.md`)](./README.md) | Infrastructure: [`main.tf`](./main.tf)

---

## 1. Start RobotoCore

* Single Rust container on `:4566`

```shell
docker run -d \
  --name robotocore \
  -p 4566:4566 \
  ghcr.io/robotocore/robotocore:latest
```

```shell
curl -s http://localhost:4566/_robotocore/health | python3 -m json.tool
```

---

## 2. Provision with Terraform (`main.tf`)

* **Account ID**: `123456789012` (parsed from `access_key`)
* **Resources (7)**:
  * **S3**: `robotocore-poc-bucket`
  * **SQS**: `robotocore-poc-queue` + `robotocore-poc-dlq`
  * **RDS**: `robotocore-poc-db` (PostgreSQL `15.4` control plane)
  * **IAM**: `robotocore-poc-role` + `robotocore-poc-policy` + attachment

```shell
terraform init
terraform plan
```

```shell
terraform apply -auto-approve
```

---

## 3. Configure AWS CLI Session

* Point CLI to `localhost:4566` and same `123456789012` account

```shell
unset AWS_PROFILE AWS_DEFAULT_PROFILE
export AWS_ENDPOINT_URL=http://localhost:4566
export AWS_ACCESS_KEY_ID=123456789012
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1
export AWS_PAGER=""
```

---

## 4. S3 — Object Storage

* Verify bucket created by Terraform, upload object, list, and download

```shell
aws s3 ls
```

```shell
echo '{"service":"robotocore","status":"verified"}' > payload.json

aws s3api put-object \
  --bucket robotocore-poc-bucket \
  --key events/payload.json \
  --body payload.json

aws s3 ls s3://robotocore-poc-bucket/events/
```

```shell
aws s3api get-object \
  --bucket robotocore-poc-bucket \
  --key events/payload.json \
  downloaded.json

cat downloaded.json
rm -f payload.json downloaded.json
```

---

## 5. SQS — Messaging & Queues

* Resolve queue URL, publish message, consume, and acknowledge (`delete-message`)

```shell
QUEUE_URL=$(aws sqs get-queue-url --queue-name robotocore-poc-queue --query 'QueueUrl' --output text)
echo "Queue URL: $QUEUE_URL"
```

```shell
aws sqs send-message \
  --queue-url "$QUEUE_URL" \
  --message-body '{"event":"user_signup","user_id":42}'
```

```shell
RECEIVE_JSON=$(aws sqs receive-message --queue-url "$QUEUE_URL" --max-number-of-messages 1 --output json)
echo "$RECEIVE_JSON"
```

```shell
RECEIPT_HANDLE=$(echo "$RECEIVE_JSON" | python3 -c "import sys, json; print(json.load(sys.stdin)['Messages'][0]['ReceiptHandle'])")

aws sqs delete-message \
  --queue-url "$QUEUE_URL" \
  --receipt-handle "$RECEIPT_HANDLE"

aws sqs receive-message --queue-url "$QUEUE_URL" --max-number-of-messages 1
```

---

## 6. RDS — Control Plane Emulation

* Instant RDS metadata & endpoint (no 10-minute cloud wait)
* *Trade-off*: Emulates control plane in memory (no Postgres Docker sidecar)

```shell
aws rds describe-db-instances \
  --db-instance-identifier robotocore-poc-db \
  --query 'DBInstances[0].{Identifier:DBInstanceIdentifier,Engine:Engine,Version:EngineVersion,Class:DBInstanceClass,Status:DBInstanceStatus,Endpoint:Endpoint.Address,Port:Endpoint.Port}' \
  --output table
```

---

## 7. IAM, STS & Validation Behavior

* Inspect IAM policy attached by Terraform

```shell
POLICY_ARN=$(aws iam list-attached-role-policies \
  --role-name robotocore-poc-role \
  --query 'AttachedPolicies[0].PolicyArn' \
  --output text)

aws iam get-policy-version \
  --policy-arn "$POLICY_ARN" \
  --version-id v1 \
  --query 'PolicyVersion.Document' \
  --output json
```

* Assume role via **STS** (`AssumeRole`)

```shell
ROLE_ARN=$(aws iam get-role --role-name robotocore-poc-role --query 'Role.Arn' --output text)

aws sts assume-role \
  --role-arn "$ROLE_ARN" \
  --role-session-name "demo-session"
```

* Test duplicate state conflict (`EntityAlreadyExists`)

```shell
aws iam create-role \
  --role-name robotocore-poc-role \
  --assume-role-policy-document '{"Version":"2012-10-17","Statement":[]}'
```

---

## 8. Platform Features: Chaos, Snapshots & Audit Log

* **Chaos Engineering**: Inject a `503 ServiceUnavailableException` fault on S3, trigger it with `aws s3 ls`, then clear the rule

```shell
curl -s -X POST http://localhost:4566/_robotocore/chaos/rules \
  -d '{"service":"s3","operation":"ListBuckets","error_code":"ServiceUnavailableException","status_code":503}' | python3 -m json.tool

# S3 now fails with our injected 503 fault!
aws s3 ls

# Clear chaos rules and verify S3 recovers
curl -s -X POST http://localhost:4566/_robotocore/chaos/rules/clear | python3 -m json.tool
aws s3 ls
```

* **State Snapshots**: Save a named checkpoint and list snapshots

```shell
curl -s -X POST http://localhost:4566/_robotocore/state/save -d '{"name": "demo-checkpoint"}' | python3 -m json.tool
curl -s http://localhost:4566/_robotocore/state/snapshots | python3 -m json.tool
```

* **Audit Log**: Inspect recent API requests recorded by the gateway

```shell
curl -s http://localhost:4566/_robotocore/audit | python3 -m json.tool | head -n 35
```

---

## 9. Teardown

```shell
terraform destroy -auto-approve
docker rm -f robotocore
```

---

### 👉 **[Back to Presentation: 6. Trade-offs (`README.md`)](./README.md#6-trade-offs)**
