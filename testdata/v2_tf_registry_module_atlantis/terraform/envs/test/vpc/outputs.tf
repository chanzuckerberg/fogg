# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.

// module "vpc" outputs
output "azs" {
  value     = module.vpc.azs
  sensitive = false
}
output "cgw_arns" {
  value     = module.vpc.cgw_arns
  sensitive = false
}
output "cgw_ids" {
  value     = module.vpc.cgw_ids
  sensitive = false
}
output "database_internet_gateway_route_id" {
  value     = module.vpc.database_internet_gateway_route_id
  sensitive = false
}
output "database_ipv6_egress_route_id" {
  value     = module.vpc.database_ipv6_egress_route_id
  sensitive = false
}
output "database_nat_gateway_route_ids" {
  value     = module.vpc.database_nat_gateway_route_ids
  sensitive = false
}
output "database_network_acl_arn" {
  value     = module.vpc.database_network_acl_arn
  sensitive = false
}
output "database_network_acl_id" {
  value     = module.vpc.database_network_acl_id
  sensitive = false
}
output "database_route_table_association_ids" {
  value     = module.vpc.database_route_table_association_ids
  sensitive = false
}
output "database_route_table_ids" {
  value     = module.vpc.database_route_table_ids
  sensitive = false
}
output "database_subnet_arns" {
  value     = module.vpc.database_subnet_arns
  sensitive = false
}
output "database_subnet_group" {
  value     = module.vpc.database_subnet_group
  sensitive = false
}
output "database_subnet_group_name" {
  value     = module.vpc.database_subnet_group_name
  sensitive = false
}
output "database_subnets" {
  value     = module.vpc.database_subnets
  sensitive = false
}
output "database_subnets_cidr_blocks" {
  value     = module.vpc.database_subnets_cidr_blocks
  sensitive = false
}
output "database_subnets_ipv6_cidr_blocks" {
  value     = module.vpc.database_subnets_ipv6_cidr_blocks
  sensitive = false
}
output "default_network_acl_id" {
  value     = module.vpc.default_network_acl_id
  sensitive = false
}
output "default_route_table_id" {
  value     = module.vpc.default_route_table_id
  sensitive = false
}
output "default_security_group_id" {
  value     = module.vpc.default_security_group_id
  sensitive = false
}
output "default_vpc_arn" {
  value     = module.vpc.default_vpc_arn
  sensitive = false
}
output "default_vpc_cidr_block" {
  value     = module.vpc.default_vpc_cidr_block
  sensitive = false
}
output "default_vpc_default_network_acl_id" {
  value     = module.vpc.default_vpc_default_network_acl_id
  sensitive = false
}
output "default_vpc_default_route_table_id" {
  value     = module.vpc.default_vpc_default_route_table_id
  sensitive = false
}
output "default_vpc_default_security_group_id" {
  value     = module.vpc.default_vpc_default_security_group_id
  sensitive = false
}
output "default_vpc_enable_dns_hostnames" {
  value     = module.vpc.default_vpc_enable_dns_hostnames
  sensitive = false
}
output "default_vpc_enable_dns_support" {
  value     = module.vpc.default_vpc_enable_dns_support
  sensitive = false
}
output "default_vpc_id" {
  value     = module.vpc.default_vpc_id
  sensitive = false
}
output "default_vpc_instance_tenancy" {
  value     = module.vpc.default_vpc_instance_tenancy
  sensitive = false
}
output "default_vpc_main_route_table_id" {
  value     = module.vpc.default_vpc_main_route_table_id
  sensitive = false
}
output "dhcp_options_id" {
  value     = module.vpc.dhcp_options_id
  sensitive = false
}
output "egress_only_internet_gateway_id" {
  value     = module.vpc.egress_only_internet_gateway_id
  sensitive = false
}
output "elasticache_network_acl_arn" {
  value     = module.vpc.elasticache_network_acl_arn
  sensitive = false
}
output "elasticache_network_acl_id" {
  value     = module.vpc.elasticache_network_acl_id
  sensitive = false
}
output "elasticache_route_table_association_ids" {
  value     = module.vpc.elasticache_route_table_association_ids
  sensitive = false
}
output "elasticache_route_table_ids" {
  value     = module.vpc.elasticache_route_table_ids
  sensitive = false
}
output "elasticache_subnet_arns" {
  value     = module.vpc.elasticache_subnet_arns
  sensitive = false
}
output "elasticache_subnet_group" {
  value     = module.vpc.elasticache_subnet_group
  sensitive = false
}
output "elasticache_subnet_group_name" {
  value     = module.vpc.elasticache_subnet_group_name
  sensitive = false
}
output "elasticache_subnets" {
  value     = module.vpc.elasticache_subnets
  sensitive = false
}
output "elasticache_subnets_cidr_blocks" {
  value     = module.vpc.elasticache_subnets_cidr_blocks
  sensitive = false
}
output "elasticache_subnets_ipv6_cidr_blocks" {
  value     = module.vpc.elasticache_subnets_ipv6_cidr_blocks
  sensitive = false
}
output "igw_arn" {
  value     = module.vpc.igw_arn
  sensitive = false
}
output "igw_id" {
  value     = module.vpc.igw_id
  sensitive = false
}
output "intra_network_acl_arn" {
  value     = module.vpc.intra_network_acl_arn
  sensitive = false
}
output "intra_network_acl_id" {
  value     = module.vpc.intra_network_acl_id
  sensitive = false
}
output "intra_route_table_association_ids" {
  value     = module.vpc.intra_route_table_association_ids
  sensitive = false
}
output "intra_route_table_ids" {
  value     = module.vpc.intra_route_table_ids
  sensitive = false
}
output "intra_subnet_arns" {
  value     = module.vpc.intra_subnet_arns
  sensitive = false
}
output "intra_subnets" {
  value     = module.vpc.intra_subnets
  sensitive = false
}
output "intra_subnets_cidr_blocks" {
  value     = module.vpc.intra_subnets_cidr_blocks
  sensitive = false
}
output "intra_subnets_ipv6_cidr_blocks" {
  value     = module.vpc.intra_subnets_ipv6_cidr_blocks
  sensitive = false
}
output "name" {
  value     = module.vpc.name
  sensitive = false
}
output "nat_ids" {
  value     = module.vpc.nat_ids
  sensitive = false
}
output "nat_public_ips" {
  value     = module.vpc.nat_public_ips
  sensitive = false
}
output "natgw_ids" {
  value     = module.vpc.natgw_ids
  sensitive = false
}
output "outpost_network_acl_arn" {
  value     = module.vpc.outpost_network_acl_arn
  sensitive = false
}
output "outpost_network_acl_id" {
  value     = module.vpc.outpost_network_acl_id
  sensitive = false
}
output "outpost_subnet_arns" {
  value     = module.vpc.outpost_subnet_arns
  sensitive = false
}
output "outpost_subnets" {
  value     = module.vpc.outpost_subnets
  sensitive = false
}
output "outpost_subnets_cidr_blocks" {
  value     = module.vpc.outpost_subnets_cidr_blocks
  sensitive = false
}
output "outpost_subnets_ipv6_cidr_blocks" {
  value     = module.vpc.outpost_subnets_ipv6_cidr_blocks
  sensitive = false
}
output "private_ipv6_egress_route_ids" {
  value     = module.vpc.private_ipv6_egress_route_ids
  sensitive = false
}
output "private_nat_gateway_route_ids" {
  value     = module.vpc.private_nat_gateway_route_ids
  sensitive = false
}
output "private_network_acl_arn" {
  value     = module.vpc.private_network_acl_arn
  sensitive = false
}
output "private_network_acl_id" {
  value     = module.vpc.private_network_acl_id
  sensitive = false
}
output "private_route_table_association_ids" {
  value     = module.vpc.private_route_table_association_ids
  sensitive = false
}
output "private_route_table_ids" {
  value     = module.vpc.private_route_table_ids
  sensitive = false
}
output "private_subnet_arns" {
  value     = module.vpc.private_subnet_arns
  sensitive = false
}
output "private_subnets" {
  value     = module.vpc.private_subnets
  sensitive = false
}
output "private_subnets_cidr_blocks" {
  value     = module.vpc.private_subnets_cidr_blocks
  sensitive = false
}
output "private_subnets_ipv6_cidr_blocks" {
  value     = module.vpc.private_subnets_ipv6_cidr_blocks
  sensitive = false
}
output "public_internet_gateway_ipv6_route_id" {
  value     = module.vpc.public_internet_gateway_ipv6_route_id
  sensitive = false
}
output "public_internet_gateway_route_id" {
  value     = module.vpc.public_internet_gateway_route_id
  sensitive = false
}
output "public_network_acl_arn" {
  value     = module.vpc.public_network_acl_arn
  sensitive = false
}
output "public_network_acl_id" {
  value     = module.vpc.public_network_acl_id
  sensitive = false
}
output "public_route_table_association_ids" {
  value     = module.vpc.public_route_table_association_ids
  sensitive = false
}
output "public_route_table_ids" {
  value     = module.vpc.public_route_table_ids
  sensitive = false
}
output "public_subnet_arns" {
  value     = module.vpc.public_subnet_arns
  sensitive = false
}
output "public_subnets" {
  value     = module.vpc.public_subnets
  sensitive = false
}
output "public_subnets_cidr_blocks" {
  value     = module.vpc.public_subnets_cidr_blocks
  sensitive = false
}
output "public_subnets_ipv6_cidr_blocks" {
  value     = module.vpc.public_subnets_ipv6_cidr_blocks
  sensitive = false
}
output "redshift_network_acl_arn" {
  value     = module.vpc.redshift_network_acl_arn
  sensitive = false
}
output "redshift_network_acl_id" {
  value     = module.vpc.redshift_network_acl_id
  sensitive = false
}
output "redshift_public_route_table_association_ids" {
  value     = module.vpc.redshift_public_route_table_association_ids
  sensitive = false
}
output "redshift_route_table_association_ids" {
  value     = module.vpc.redshift_route_table_association_ids
  sensitive = false
}
output "redshift_route_table_ids" {
  value     = module.vpc.redshift_route_table_ids
  sensitive = false
}
output "redshift_subnet_arns" {
  value     = module.vpc.redshift_subnet_arns
  sensitive = false
}
output "redshift_subnet_group" {
  value     = module.vpc.redshift_subnet_group
  sensitive = false
}
output "redshift_subnets" {
  value     = module.vpc.redshift_subnets
  sensitive = false
}
output "redshift_subnets_cidr_blocks" {
  value     = module.vpc.redshift_subnets_cidr_blocks
  sensitive = false
}
output "redshift_subnets_ipv6_cidr_blocks" {
  value     = module.vpc.redshift_subnets_ipv6_cidr_blocks
  sensitive = false
}
output "this_customer_gateway" {
  value     = module.vpc.this_customer_gateway
  sensitive = false
}
output "vgw_arn" {
  value     = module.vpc.vgw_arn
  sensitive = false
}
output "vgw_id" {
  value     = module.vpc.vgw_id
  sensitive = false
}
output "vpc_arn" {
  value     = module.vpc.vpc_arn
  sensitive = false
}
output "vpc_cidr_block" {
  value     = module.vpc.vpc_cidr_block
  sensitive = false
}
output "vpc_enable_dns_hostnames" {
  value     = module.vpc.vpc_enable_dns_hostnames
  sensitive = false
}
output "vpc_enable_dns_support" {
  value     = module.vpc.vpc_enable_dns_support
  sensitive = false
}
output "vpc_flow_log_cloudwatch_iam_role_arn" {
  value     = module.vpc.vpc_flow_log_cloudwatch_iam_role_arn
  sensitive = false
}
output "vpc_flow_log_destination_arn" {
  value     = module.vpc.vpc_flow_log_destination_arn
  sensitive = false
}
output "vpc_flow_log_destination_type" {
  value     = module.vpc.vpc_flow_log_destination_type
  sensitive = false
}
output "vpc_flow_log_id" {
  value     = module.vpc.vpc_flow_log_id
  sensitive = false
}
output "vpc_id" {
  value     = module.vpc.vpc_id
  sensitive = false
}
output "vpc_instance_tenancy" {
  value     = module.vpc.vpc_instance_tenancy
  sensitive = false
}
output "vpc_ipv6_association_id" {
  value     = module.vpc.vpc_ipv6_association_id
  sensitive = false
}
output "vpc_ipv6_cidr_block" {
  value     = module.vpc.vpc_ipv6_cidr_block
  sensitive = false
}
output "vpc_main_route_table_id" {
  value     = module.vpc.vpc_main_route_table_id
  sensitive = false
}
output "vpc_owner_id" {
  value     = module.vpc.vpc_owner_id
  sensitive = false
}
output "vpc_secondary_cidr_blocks" {
  value     = module.vpc.vpc_secondary_cidr_blocks
  sensitive = false
}
// module "my_module" outputs
