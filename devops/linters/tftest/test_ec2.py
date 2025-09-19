import tftest
import pytest

def test_ec2_configuration():
    # Initialize terraform
    tf = tftest.TerraformTest('.')
    
    # Plan terraform
    tf.setup()
    plan = tf.plan(output=True)
    
    # Get the planned values
    planned_values = plan['planned_values']
    resources = planned_values['root_module']['resources']
    
    # Test VPC configuration
    vpc = next(r for r in resources if r['type'] == 'aws_vpc' and r['name'] == 'main')
    assert vpc['values']['cidr_block'] == '10.0.0.0/16'
    assert vpc['values']['enable_dns_hostnames'] is True
    assert vpc['values']['enable_dns_support'] is True
    
    # Test subnet configuration
    subnet = next(r for r in resources if r['type'] == 'aws_subnet' and r['name'] == 'public')
    assert subnet['values']['cidr_block'] == '10.0.1.0/24'
    assert subnet['values']['map_public_ip_on_launch'] is True
    assert subnet['values']['availability_zone'] == 'us-west-2a'
    
    # Test security group configuration
    sg = next(r for r in resources if r['type'] == 'aws_security_group' and r['name'] == 'web')
    
    # Verify ingress rules
    ingress_rules = sg['values']['ingress']
    http_rule = next(rule for rule in ingress_rules if rule['from_port'] == 80)
    assert http_rule['to_port'] == 80
    assert http_rule['protocol'] == 'tcp'
    assert http_rule['cidr_blocks'] == ['0.0.0.0/0']
    
    ssh_rule = next(rule for rule in ingress_rules if rule['from_port'] == 22)
    assert ssh_rule['to_port'] == 22
    assert ssh_rule['protocol'] == 'tcp'
    assert ssh_rule['cidr_blocks'] == ['0.0.0.0/0']
    
    # Verify egress rules
    egress_rule = sg['values']['egress'][0]
    assert egress_rule['from_port'] == 0
    assert egress_rule['to_port'] == 0
    assert egress_rule['protocol'] == '-1'
    assert egress_rule['cidr_blocks'] == ['0.0.0.0/0']
    
    # Test EC2 instance configuration
    ec2 = next(r for r in resources if r['type'] == 'aws_instance' and r['name'] == 'web')
    assert ec2['values']['instance_type'] == 't2.micro'
    assert ec2['values']['ami'] == 'ami-0735c191cf914754d'
    assert ec2['values']['associate_public_ip_address'] is True
    
    # Test EC2 root block device
    root_block_device = ec2['values']['root_block_device'][0]
    assert root_block_device['volume_size'] == 8
    assert root_block_device['volume_type'] == 'gp2'
    assert root_block_device['encrypted'] is True
    
    # Test EC2 tags
    assert ec2['values']['tags'] == {
        'Name': 'web-server',
        'Environment': 'test',
        'Project': 'TFTest Example'
    }

def test_networking_configuration():
    # Initialize terraform
    tf = tftest.TerraformTest('.')
    
    # Plan terraform
    tf.setup()
    plan = tf.plan(output=True)
    
    # Get the planned values
    planned_values = plan['planned_values']
    resources = planned_values['root_module']['resources']
    
    # Test Internet Gateway configuration
    igw = next(r for r in resources if r['type'] == 'aws_internet_gateway' and r['name'] == 'main')
    assert igw['values']['tags']['Name'] == 'test-igw'
    
    # Test Route Table configuration
    route_table = next(r for r in resources if r['type'] == 'aws_route_table' and r['name'] == 'main')
    routes = route_table['values']['route']
    internet_route = next(route for route in routes if route['cidr_block'] == '0.0.0.0/0')
    assert internet_route is not None

if __name__ == '__main__':
    pytest.main([__file__])