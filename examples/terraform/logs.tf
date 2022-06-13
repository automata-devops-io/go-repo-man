# logs.tf

# Set up CloudWatch group and log stream and retain logs for 30 days
resource "aws_cloudwatch_log_group" "rM_log_group" {
  name              = "/ecs/rM-app"
  retention_in_days = 30

  tags = {
    Name = "rM-log-group"
  }
}

resource "aws_cloudwatch_log_stream" "rM_log_stream" {
  name           = "rM-log-stream"
  log_group_name = aws_cloudwatch_log_group.rM_log_group.name
}
