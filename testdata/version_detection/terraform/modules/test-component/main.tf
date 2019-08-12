locals {
  service_name = "${var.project}-${var.env}-${var.name}"
}

data "aws_iam_policy_document" "assume_role" {
  statement {
    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "role" {
  name               = "${local.service_name}"
  assume_role_policy = "${data.aws_iam_policy_document.assume_role.json}"
}

module "parameters-policy" {
  source = "github.com/chanzuckerberg/cztack//aws-params-reader-policy?ref=v0.15.1"

  project   = "${var.project}"
  env       = "${var.env}"
  service   = "${var.name}"
  region    = "${var.region}"
  role_name = "${aws_iam_role.role.name}"
}

module "container_environment" {
  source  = "github.com/chanzuckerberg/cztack//aws-params-writer?ref=v0.15.1"
  project = "${var.project}"
  env     = "${var.env}"
  service = "${var.name}"
  owner   = "${var.owner}"

  parameters = {
    AWS_REGION         = "${var.region}"
    AWS_DEFAULT_REGION = "${var.region}"
    DATABASE           = "${var.db_name}"
    DATABASE_USER      = "${var.db_user}"
    DATABASE_HOST      = "${var.db_host}"
  }
}

module "sg" {
  create      = "${var.create_security_group}"
  source      = "terraform-aws-modules/security-group/aws"
  version     = "3.1.0"
  name        = "${local.name}-alb"
  description = "Security group for ${var.internal ? "internal" : "internet facing"} ALB"
  vpc_id      = "${var.vpc_id}"
  tags        = "${local.tags}"

  ingress_cidr_blocks = "${var.ingress_cidrs}"
  egress_cidr_blocks  = "${var.egress_cidrs}"
  ingress_rules       = ["https-443-tcp", "http-80-tcp"]
  egress_rules        = ["all-all"]
}