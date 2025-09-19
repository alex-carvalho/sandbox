import tftest
import pytest

def test_s3_bucket_creation():
    # Initialize terraform
    tf = tftest.TerraformTest('.')
    
    # Plan terraform
    tf.setup()
    plan = tf.plan(output=True)
    
    # Get the planned changes
    planned_values = plan['planned_values']
    s3_resources = planned_values['root_module']['resources']
    
    # Find the S3 bucket in the plan
    s3_bucket = next(r for r in s3_resources if r['type'] == 'aws_s3_bucket' and r['name'] == 'example')
    
    # Verify bucket configuration
    assert s3_bucket['values']['bucket'] == 'my-test-bucket-name'
    assert s3_bucket['values']['tags'] == {
        'Environment': 'Test',
        'Project': 'TFTest Example'
    }
    
    # Find the versioning configuration
    versioning = next(r for r in s3_resources if r['type'] == 'aws_s3_bucket_versioning' and r['name'] == 'example')
    assert versioning['values']['versioning_configuration'][0]['status'] == 'Enabled'
    
    # Find the encryption configuration
    encryption = next(r for r in s3_resources if r['type'] == 'aws_s3_bucket_server_side_encryption_configuration' and r['name'] == 'example')
    assert encryption['values']['rule'][0]['apply_server_side_encryption_by_default'][0]['sse_algorithm'] == 'AES256'

if __name__ == '__main__':
    pytest.main([__file__])