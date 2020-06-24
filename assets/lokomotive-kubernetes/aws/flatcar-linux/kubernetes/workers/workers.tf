# Workers AutoScaling Group
resource "aws_autoscaling_group" "workers" {
  count = var.worker_count > 0 ? 1 : 0

  name = "${var.pool_name}-worker"

  # count
  desired_capacity          = var.worker_count
  min_size                  = var.worker_count
  max_size                  = var.worker_count + 2
  default_cooldown          = 30
  health_check_grace_period = 30

  # network
  vpc_zone_identifier = var.subnet_ids

  # template
  launch_configuration = aws_launch_configuration.worker[0].name

  # target groups to which instances should be added
  target_group_arns = flatten([
    aws_lb_target_group.workers_http.id,
    aws_lb_target_group.workers_https.id,
    var.target_groups,
  ])

  lifecycle {
    # override the default destroy and replace update behavior
    create_before_destroy = true
  }

  # Waiting for instance creation delays adding the ASG to state. If instances
  # can't be created (e.g. spot price too low), the ASG will be orphaned.
  # Orphaned ASGs escape cleanup, can't be updated, and keep bidding if spot is
  # used. Disable wait to avoid issues and align with other clouds.
  wait_for_capacity_timeout = "0"

  tags = flatten([
    [
      {
        key                 = "Name"
        value               = "${var.cluster_name}-${var.pool_name}-worker"
        propagate_at_launch = true
      },
    ],
    [
      for tag in keys(var.tags) :
      {
        key                 = tag == "Name" ? "X-Name" : tag
        value               = var.tags[tag]
        propagate_at_launch = true
      }
    ],
  ])
}

# Worker template
resource "aws_launch_configuration" "worker" {
  count = var.worker_count

  name_prefix       = "${var.cluster_name}-${var.pool_name}-"
  image_id          = local.ami_id
  instance_type     = var.instance_type
  spot_price        = var.spot_price
  enable_monitoring = false

  user_data = data.ct_config.worker-ignition[count.index].rendered

  # storage
  root_block_device {
    volume_type = var.disk_type
    volume_size = var.disk_size
    iops        = var.disk_iops
    encrypted   = true
  }

  # network
  security_groups = var.security_groups

  lifecycle {
    // Override the default destroy and replace update behavior
    create_before_destroy = true
    ignore_changes        = [image_id]
  }
}

# Worker Ignition config
data "ct_config" "worker-ignition" {
  count = var.worker_count

  content      = data.template_file.worker-config[count.index].rendered
  pretty_print = false
  snippets     = var.clc_snippets
}

# Worker Container Linux config
data "template_file" "worker-config" {
  count = var.worker_count

  template = file("${path.module}/cl/worker.yaml.tmpl")

  vars = {
    kubeconfig     = indent(10, templatefile("${path.module}/cl/bootstrap-kubeconfig.yaml.tmpl", {
      token_id     = random_string.bootstrap-token-id[count.index].result
      token_secret = random_string.bootstrap-token-secret[count.index].result
      ca_cert      = var.ca_cert
      server       = "https://${var.apiserver}:6443"
    }))

    ssh_keys               = jsonencode(var.ssh_keys)
    cluster_dns_service_ip = cidrhost(var.service_cidr, 10)
    cluster_domain_suffix  = var.cluster_domain_suffix
  }
}
