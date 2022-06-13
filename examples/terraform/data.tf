data "vault_aws_access_credentials" "telemetry" {
  type    = "sts"
  backend = "account/dynamic/aws/your-account"
  role    = "default"
}

data "template_file" "rM_app" {
  template = file("./templates/ecs/repoMan_app.json.tpl")

  vars = {
    app_image      = var.app_image
    app_port       = var.app_port
    fargate_cpu    = var.fargate_cpu
    fargate_memory = var.fargate_memory
    aws_region     = var.aws_region
  }
}

# Fetch AZs in the current region
data "aws_availability_zones" "available" {
}
