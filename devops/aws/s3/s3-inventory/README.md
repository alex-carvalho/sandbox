# 🧾 AWS S3 Inventory — Explained

## 📘 Overview

**Amazon S3 Inventory** is a feature that provides **scheduled reports** listing objects and their metadata in an S3 bucket.  
It helps you **audit and manage** the contents of your buckets **without needing to run List or HeadObject requests** (which can be expensive and time-consuming for large buckets).

Instead, you get a **daily or weekly CSV, ORC, or Parquet file** delivered to a destination bucket.

---

## 🎯 Key Use Cases

- **Auditing and Compliance:** Verify object encryption status, replication, or tagging compliance.
- **Cost Optimization:** Identify old or untagged objects.
- **Security Monitoring:** Check which objects lack encryption or ownership controls.
- **Automation:** Use the report as input for scripts that take bulk actions on objects.

---

## ⚙️ How It Works

1. You **configure an inventory** on a source S3 bucket.
2. S3 generates a report at a chosen **frequency (Daily or Weekly)**.
3. The report is stored in a **destination S3 bucket** (can be in the same or another account).
4. Each report lists:
   - Object name (key)
   - Size
   - Storage class
   - Last modified date
   - Encryption status
   - Replication status
   - Object lock status
   - Custom metadata (optional)

---

## 📂 Report File Format

You can choose among:
- **CSV** — Simple text-based list (easier for humans to read)
- **ORC** — Optimized Row Columnar format for efficient queries in **Amazon Athena** or **EMR**
- **Parquet** — Columnar format used for **data analytics** workloads

---

## 🧩 Important Details (AWS Certification Focus)

| Topic | Key Details |
|--------|--------------|
| **Frequency** | Daily or Weekly |
| **Delivery** | To an S3 bucket (can be cross-account) |
| **Formats** | CSV, ORC, Parquet |
| **Encryption** | Inventory reports can be encrypted with **SSE-S3**, **SSE-KMS**, or **SSE-C** |
| **Filters** | You can limit inventory to a specific **prefix** or **object tag** |
| **Metadata Fields** | Choose which metadata to include (e.g., size, ETag, replication status, etc.) |
| **Permissions** | Destination bucket must grant S3 permission for inventory delivery |
| **Integration** | Works with **Amazon Athena**, **AWS Glue**, and **Amazon Redshift Spectrum** for querying reports |

---

## 🪣 Example — Using AWS CLI

### 1. Create an Inventory Configuration

```bash
aws s3api put-bucket-inventory-configuration   --bucket my-source-bucket   --id my-inventory-config   --inventory-configuration '{
    "Id": "my-inventory-config",
    "IsEnabled": true,
    "IncludedObjectVersions": "Current",
    "Schedule": { "Frequency": "Daily" },
    "Destination": {
      "S3BucketDestination": {
        "AccountId": "123456789012",
        "Bucket": "arn:aws:s3:::my-destination-bucket",
        "Format": "CSV",
        "Prefix": "inventory-reports/",
        "Encryption": { "SSES3": {} }
      }
    },
    "OptionalFields": [
      "Size",
      "StorageClass",
      "LastModifiedDate",
      "ETag",
      "EncryptionStatus",
      "ReplicationStatus"
    ]
  }'
```

---

### 2. List Existing Inventory Configurations

```bash
aws s3api list-bucket-inventory-configurations   --bucket my-source-bucket
```

---

### 3. Get a Specific Inventory Configuration

```bash
aws s3api get-bucket-inventory-configuration   --bucket my-source-bucket   --id my-inventory-config
```

---

### 4. Delete an Inventory Configuration

```bash
aws s3api delete-bucket-inventory-configuration   --bucket my-source-bucket   --id my-inventory-config
```

---

## 🧠 Exam Tips (AWS Certifications)

✅ **Remember:** S3 Inventory is **not real-time** — it provides **periodic snapshots** of objects.  
✅ It can replace **S3 LIST operations** for large-scale analysis.  
✅ Reports can be queried with **Amazon Athena** for quick data analysis.  
✅ **Cross-region and cross-account** inventory delivery are supported.  
✅ Inventory supports **both current and noncurrent object versions** (if versioning is enabled).  
✅ Common question topic: **Encryption audit** (check which objects are missing encryption).

---

## 📊 Example Report (CSV)

| Bucket | Key | Size | StorageClass | LastModifiedDate | EncryptionStatus |
|---------|-----|------|---------------|------------------|------------------|
| my-source-bucket | images/cat.jpg | 14532 | STANDARD | 2025-10-01T10:25:00Z | SSE-S3 |
| my-source-bucket | logs/2025/10/01.log | 1048576 | GLACIER | 2025-10-01T00:00:00Z | None |

---

## 🔗 Related AWS Services

- **Amazon Athena:** Query inventory reports directly from S3.  
- **AWS Glue:** Catalog inventory data for ETL jobs.  
- **AWS Lambda:** Automate actions (e.g., delete unencrypted objects).  
- **Amazon S3 Storage Lens:** Provides deeper analytics, but S3 Inventory gives **object-level detail**.
