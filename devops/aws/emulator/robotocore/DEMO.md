# Live Demo: RobotoCore (`Terraform` + `AWS CLI`)

> [← Back to Presentation (`README.md`)](./README.md) | Infrastructure: [`main.tf`](./main.tf)

---

## 1. Start RobotoCore

* Single container on `:4566`
* `ENFORCE_IAM=1` enables real-time IAM policy evaluation

```shell
docker run -d \
  --name robotocore \
  -p 4566:4566 \
  -e ENFORCE_IAM=1 \
  ghcr.io/robotocore/robotocore:latest
```

```shell
curl -s http://localhost:4566/_robotocore/health
```

---

## 2. Provision with Terraform (`main.tf`)

* **Account ID**: `123456789012` (parsed from `access_key`)
* **Resources (7)**:
  * **S3**: `robotocore-poc-bucket`
  * **SQS**: `robotocore-poc-queue` + `robotocore-poc-dlq`
  * **RDS**: `robotocore-poc-db` (PostgreSQL `15.4` control plane)
  * **IAM**: `robotocore-poc-role` + `robotocore-poc-policy` (Allow S3/SQS/RDS, explicit **Deny** on delete bucket/queue)

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
* *Trade-off*: Control-plane & Data API in memory (no Postgres Docker sidecar)

```shell
aws rds describe-db-instances \
  --db-instance-identifier robotocore-poc-db \
  --query 'DBInstances[0].{Identifier:DBInstanceIdentifier,Engine:Engine,Version:EngineVersion,Class:DBInstanceClass,Status:DBInstanceStatus,Endpoint:Endpoint.Address,Port:Endpoint.Port}' \
  --output table
```

---

## 7. IAM, STS & Runtime Enforcement (`ENFORCE_IAM=1`)

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

* Assume role via **STS** to get temporary credentials

```shell
ROLE_ARN=$(aws iam get-role --role-name robotocore-poc-role --query 'Role.Arn' --output text)

ASSUME_JSON=$(aws sts assume-role \
  --role-arn "$ROLE_ARN" \
  --role-session-name "demo-session" \
  --output json)

echo "$ASSUME_JSON" | python3 -m json.tool
```

* Test **Allowed** (`s3:ListBucket`) vs **Explicitly Denied** (`s3:DeleteBucket` → `AccessDenied`)

```shell
ROLE_AK=$(echo "$ASSUME_JSON" | python3 -c "import sys, json; print(json.load(sys.stdin)['Credentials']['AccessKeyId'])")
ROLE_SK=$(echo "$ASSUME_JSON" | python3 -c "import sys, json; print(json.load(sys.stdin)['Credentials']['SecretAccessKey'])")
ROLE_ST=$(echo "$ASSUME_JSON" | python3 -c "import sys, json; print(json.load(sys.stdin)['Credentials']['SessionToken'])")

# Allowed by policy
AWS_ACCESS_KEY_ID="$ROLE_AK" AWS_SECRET_ACCESS_KEY="$ROLE_SK" AWS_SESSION_TOKEN="$ROLE_ST" \
  aws s3 ls s3://robotocore-poc-bucket

# Blocked by ENFORCE_IAM=1 (Explicit Deny)
AWS_ACCESS_KEY_ID="$ROLE_AK" AWS_SECRET_ACCESS_KEY="$ROLE_SK" AWS_SESSION_TOKEN="$ROLE_ST" \
  aws s3api delete-bucket --bucket robotocore-poc-bucket
```

---

## 8. Validation: State Conflicts vs. Policy Syntax

* **State Conflict (`EntityAlreadyExists`)**: Enforced by RobotoCore

```shell
aws iam create-role \
  --role-name robotocore-poc-role \
  --assume-role-policy-document '{"Version":"2012-10-17","Statement":[]}'
```

* **Policy Schema Validation**:
  * RobotoCore stores `PolicyDocument` as raw string on `CreatePolicy` and evaluates at **runtime** (`ENFORCE_IAM=1`).
  * Creation-time syntax checks are enforced by **Terraform** (`data "aws_iam_policy_document"`).

---

## 9. Teardown

```shell
terraform destroy -auto-approve
docker rm -f robotocore
```

---

### 👉 **[Back to Presentation: 6. Trade-offs (`README.md`)](./README.md#6-trade-offs)**
