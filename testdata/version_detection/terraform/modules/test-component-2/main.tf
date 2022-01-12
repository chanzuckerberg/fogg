module "test-module" {
  create      = var.create_security_group
  source      = "terraform-aws-modules/security-group/aws"
  version     = "3.1.0"
  name        = "${local.name}-alb"
  description = "Security group for ${var.internal ? "internal" : "internet facing"} ALB"
  vpc_id      = var.vpc_id
  tags        = local.tags

  ingress_cidr_blocks = var.ingress_cidrs
  egress_cidr_blocks  = var.egress_cidrs
  ingress_rules       = ["https-443-tcp", "http-80-tcp"]
  egress_rules        = ["all-all"]
}
