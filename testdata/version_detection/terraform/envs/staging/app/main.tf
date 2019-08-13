module "test" {
  source = "../../../modules/test-component"

  project = "${var.project}"
  env     = "${var.env}"
  name    = "golinks"
  owner   = "${var.owner}"

  db_host = "${data.terraform_remote_state.golinksdb.endpoint}"
  db_name = "${data.terraform_remote_state.golinksdb.database_name}"
  db_user = "${data.terraform_remote_state.golinksdb.database_username}"

  ecs_cluster_id              = "${data.terraform_remote_state.ecs.cluster_id}"
  egress_cidrs                = "${data.terraform_remote_state.cloud-env.private_subnets_cidr_blocks}"
  ingress_cidrs               = "${data.terraform_remote_state.cloud-env.private_subnets_cidr_blocks}"
  region                      = "${var.region}"
  route53_zone_id             = "${data.terraform_remote_state.route53.gostaging_si_czi_technology}"
  internal_route53_zone_id    = "${data.terraform_remote_state.route53.staging_si_czi_technology}"
  private_subnets             = "${data.terraform_remote_state.cloud-env.private_subnets}"
  public_subnets              = "${data.terraform_remote_state.cloud-env.public_subnets}"
  vpc_id                      = "${data.terraform_remote_state.cloud-env.vpc_id}"
  public_subdomain            = ""
  use_fargate                 = "true"
  registry_secretsmanager_arn = "${data.terraform_remote_state.credentials.dockerhub_czisi_arn}"
}

module "parameters-policy" {
  source = "github.com/chanzuckerberg/cztack//aws-params-reader-policy?ref=v0.15.1"

  project   = "${var.project}"
  env       = "${var.env}"
  service   = "${var.name}"
  region    = "${var.region}"
  role_name = "${aws_iam_role.role.name}"
}

module "mock-module" {
  create      = "${var.create_security_group}"
  source      = "terraform-aws-modules/security-group/aws"
  version     = "3.1.0"
  name        = "${local.name}-alb"
  description = "Security group"
  vpc_id      = "${var.vpc_id}"
  tags        = "${local.tags}"

  ingress_cidr_blocks = "${var.ingress_cidrs}"
  egress_cidr_blocks  = "${var.egress_cidrs}"
  ingress_rules       = ["https-443-tcp", "http-80-tcp"]
  egress_rules        = ["all-all"]
}