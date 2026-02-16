from tikv_client import RawClient

# Connect to TiKV cluster
client = RawClient.connect(["localhost:2379"])

# Put key-value pair
client.put(b"key", b"value")

# Get value
value = client.get(b"key")
print(f"Value: {value}")

# Delete key
client.delete(b"key")

print("✓ Test completed successfully!")