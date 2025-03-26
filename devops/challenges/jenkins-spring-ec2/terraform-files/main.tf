resource "aws_instance" "example" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t2.micro"

  tags = {
    Name = "App EC2"
  }

  user_data = <<-EOF
    #!/bin/bash
    sudo yum update -y
    sudo amazon-linux-extras enable corretto8
    sudo yum install -y java-1.8.0-amazon-corretto

    # Create an application directory
    mkdir -p /home/ec2-user/app
    cd /home/ec2-user/app

    # Download JAR file from S3 (replace with your S3 bucket and JAR file name)
    aws s3 cp s3://my-artifacts/app.jar /home/ec2-user/app/app.jar

    # Run the JAR file
    nohup java -jar app.jar > app.log 2>&1 &
  EOF
}
