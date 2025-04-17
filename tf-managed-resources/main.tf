provider "aws" {
  region = "us-east-2"
}

resource "tls_private_key" "deployer" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "local_file" "private_key" {
  filename        = "id_rsa"
  content         = tls_private_key.deployer.private_key_pem
  file_permission = "0600"
}

resource "aws_key_pair" "deployer" {
  key_name   = "deployer-key"
  public_key = tls_private_key.deployer.public_key_openssh
}

resource "aws_security_group" "pods_security_group" {
  name_prefix = "pods-security-group"

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "slave" {
  ami             = "ami-0100e595e1cc1ff7f"
  instance_type   = "t2.micro"
  key_name        = aws_key_pair.deployer.key_name
  security_groups = [aws_security_group.pods_security_group.name]
  
  root_block_device {
    volume_size = 32
  }
  
  tags = {
    Name = "slave"
  }
}

resource "aws_instance" "slave2" {
  ami             = "ami-0100e595e1cc1ff7f"
  instance_type   = "t2.micro"
  key_name        = aws_key_pair.deployer.key_name
  security_groups = [aws_security_group.pods_security_group.name]
  
  root_block_device {
    volume_size = 32
  }
  
  tags = {
    Name = "slave2"
  }
}

resource "aws_instance" "slave3" {
  ami             = "ami-0100e595e1cc1ff7f"
  instance_type   = "t2.micro"
  key_name        = aws_key_pair.deployer.key_name
  security_groups = [aws_security_group.pods_security_group.name]
  
  root_block_device {
    volume_size = 32
  }
  
  tags = {
    Name = "slave3"
  }
}

resource "aws_instance" "slave4" {
  ami             = "ami-0100e595e1cc1ff7f"
  instance_type   = "t2.micro"
  key_name        = aws_key_pair.deployer.key_name
  security_groups = [aws_security_group.pods_security_group.name]
  
  root_block_device {
    volume_size = 32
  }
  
  tags = {
    Name = "slave4"
  }
}

resource "aws_instance" "slave5" {
  ami             = "ami-0100e595e1cc1ff7f"
  instance_type   = "t2.micro"
  key_name        = aws_key_pair.deployer.key_name
  security_groups = [aws_security_group.pods_security_group.name]
  
  root_block_device {
    volume_size = 32
  }
  
  tags = {
    Name = "slave5"
  }
}

resource "aws_instance" "slave6" {
  ami             = "ami-0100e595e1cc1ff7f"
  instance_type   = "t2.micro"
  key_name        = aws_key_pair.deployer.key_name
  security_groups = [aws_security_group.pods_security_group.name]
  
  root_block_device {
    volume_size = 32
  }
  
  tags = {
    Name = "slave6"
  }
}

resource "aws_instance" "master" {
  ami             = "ami-0100e595e1cc1ff7f"
  instance_type   = "t2.micro"
  key_name        = aws_key_pair.deployer.key_name
  security_groups = [aws_security_group.pods_security_group.name]
  tags = {
    Name = "master"
  }
}

output "master_ip" {
  value = aws_instance.master.public_ip
}

output "slave_ips" {
  value = [
    aws_instance.slave.public_ip,
    aws_instance.slave2.public_ip,
    aws_instance.slave3.public_ip,
    aws_instance.slave4.public_ip,
    aws_instance.slave5.public_ip,
    aws_instance.slave6.public_ip
  ]
}
