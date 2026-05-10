# SG для ALB — принимает HTTP из интернета
resource "aws_security_group" "alb" {
  name   = "${var.project}-alb-sg"
  vpc_id = aws_vpc.main.id

  ingress { # поток данных из интернета
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]   # весь интернет
  }

  egress {  # поток данных в интернет
    from_port   = 0
    to_port     = 0
    protocol    = "-1"            # -1 = все протоколы
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = { Name = "${var.project}-alb-sg" }
}

# SG для ECS — принимает трафик ТОЛЬКО от ALB
resource "aws_security_group" "ecs" {
  name   = "${var.project}-ecs-sg"
  vpc_id = aws_vpc.main.id

  ingress {
    from_port       = 8080
    to_port         = 8080
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]   # только от ALB
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = { Name = "${var.project}-ecs-sg" }
}