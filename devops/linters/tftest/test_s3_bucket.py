import tftest
import pytest

@pytest.fixture(scope='module')
def resources():
    # Initialize terraform
    tf = tftest.TerraformTest('.')
    
    # Plan terraform
    tf.setup()
    plan = tf.plan(output=True)
    
    # Get the planned values 
    planned_values = plan['planned_values']
    resources = planned_values['root_module']['resources']
    return resources

def test_s3_bucket_creation(resources):
    s3_bucket = next(r for r in resources if r['type'] == 'aws_s3_bucket' and r['name'] == 'example')
    
    assert s3_bucket['values']['bucket'] == 'my-test-bucket-name'
    assert s3_bucket['values']['tags'] == {
        'Environment': 'Test',
        'Project': 'TFTest Example'
    }

def test_s3_bucket_versioning(resources):
    for resource in resources:
        if resource['type'] == 'aws_s3_bucket_versioning':
            assert resource['values']['versioning_configuration'][0]['status'] == 'Enabled'
    
def test_s3_bucket_encryption(resources):
     for resource in resources:
        if resource['type'] == 'aws_s3_bucket_server_side_encryption_configuration':
            assert resource['values']['rule'][0]['apply_server_side_encryption_by_default'][0]['sse_algorithm'] == 'AES256'

if __name__ == '__main__':
    pytest.main([__file__])