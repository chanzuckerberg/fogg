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
