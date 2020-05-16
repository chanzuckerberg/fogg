# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.

module terraform-aws-vpc {
  source                             = "github.com/terraform-aws-modules/terraform-aws-vpc?ref=v1.30.0"
  azs                                = local.azs
  cidr                               = local.cidr
  create_database_subnet_group       = local.create_database_subnet_group
  create_vpc                         = local.create_vpc
  database_subnet_tags               = local.database_subnet_tags
  database_subnets                   = local.database_subnets
  default_route_table_tags           = local.default_route_table_tags
  default_vpc_enable_classiclink     = local.default_vpc_enable_classiclink
  default_vpc_enable_dns_hostnames   = local.default_vpc_enable_dns_hostnames
  default_vpc_enable_dns_support     = local.default_vpc_enable_dns_support
  default_vpc_name                   = local.default_vpc_name
  default_vpc_tags                   = local.default_vpc_tags
  dhcp_options_domain_name           = local.dhcp_options_domain_name
  dhcp_options_domain_name_servers   = local.dhcp_options_domain_name_servers
  dhcp_options_netbios_name_servers  = local.dhcp_options_netbios_name_servers
  dhcp_options_netbios_node_type     = local.dhcp_options_netbios_node_type
  dhcp_options_ntp_servers           = local.dhcp_options_ntp_servers
  dhcp_options_tags                  = local.dhcp_options_tags
  elasticache_subnet_tags            = local.elasticache_subnet_tags
  elasticache_subnets                = local.elasticache_subnets
  enable_dhcp_options                = local.enable_dhcp_options
  enable_dns_hostnames               = local.enable_dns_hostnames
  enable_dns_support                 = local.enable_dns_support
  enable_dynamodb_endpoint           = local.enable_dynamodb_endpoint
  enable_nat_gateway                 = local.enable_nat_gateway
  enable_s3_endpoint                 = local.enable_s3_endpoint
  enable_vpn_gateway                 = local.enable_vpn_gateway
  external_nat_ip_ids                = local.external_nat_ip_ids
  instance_tenancy                   = local.instance_tenancy
  manage_default_vpc                 = local.manage_default_vpc
  map_public_ip_on_launch            = local.map_public_ip_on_launch
  name                               = local.name
  private_route_table_tags           = local.private_route_table_tags
  private_subnet_tags                = local.private_subnet_tags
  private_subnets                    = local.private_subnets
  propagate_private_route_tables_vgw = local.propagate_private_route_tables_vgw
  propagate_public_route_tables_vgw  = local.propagate_public_route_tables_vgw
  public_route_table_tags            = local.public_route_table_tags
  public_subnet_tags                 = local.public_subnet_tags
  public_subnets                     = local.public_subnets
  redshift_subnet_tags               = local.redshift_subnet_tags
  redshift_subnets                   = local.redshift_subnets
  reuse_nat_ips                      = local.reuse_nat_ips
  single_nat_gateway                 = local.single_nat_gateway
  tags                               = local.tags
  vpc_tags                           = local.vpc_tags
  vpn_gateway_id                     = local.vpn_gateway_id

}
