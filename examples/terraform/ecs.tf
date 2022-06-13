# ecs.tf

resource "aws_ecs_cluster" "main" {
  name = "rM-builder-cluster"
}

resource "aws_ecs_task_definition" "app" {
  family                   = "rM-builder-task"
  execution_role_arn       = "arn:aws:iam::862552916454:role/ecsTaskExecutionRole"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.fargate_cpu
  memory                   = var.fargate_memory
  container_definitions    = data.template_file.la_app.rendered
}

resource "aws_ecs_service" "main" {
  name            = "rM-builder"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.app.arn
  desired_count   = var.app_count
  launch_type     = "FARGATE"

  network_configuration {
    security_groups  = [aws_security_group.ecs_tasks.id]
    subnets          = aws_subnet.private.*.id
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_alb_target_group.app.id
    container_name   = "rM-app"
    container_port   = var.app_port
  }

}

depends_on = [aws_alb_listener.front_end]
