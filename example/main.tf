variable "peering_account_id" {
	type = string
}

variable "peering_network_id" {
	type = string
}

variable "peering_route" {
	type = string
}

provider "eventstorecloud" {}

resource "eventstorecloud_project" "chicken_window" {
	name = "Improved Chicken Window"
}

resource "eventstorecloud_network" "chicken_window" {
	name = "Chicken Window Net"

	project_id = eventstorecloud_project.chicken_window.id

	resource_provider = "aws"
	region = "us-west-2"
	cidr_block = "172.21.0.0/16"
}

resource "eventstorecloud_peering" "peering" {
	name = "Example Peering"

	project_id = eventstorecloud_network.chicken_window.project_id
	network_id = eventstorecloud_network.chicken_window.id

	peer_resource_provider = eventstorecloud_network.chicken_window.resource_provider
	peer_network_region = eventstorecloud_network.chicken_window.region

	peer_account_id = var.peering_account_id
	peer_network_id = var.peering_network_id
	routes = [var.peering_route]
}

resource "eventstorecloud_managed_cluster" "wings" {
	name = "Wings Cluster"

	project_id = eventstorecloud_network.chicken_window.project_id
	network_id = eventstorecloud_network.chicken_window.id

	topology = "single-node"
	instance_type = "F1"
	disk_size = 16
	disk_type = "gp2"
	server_version = "20.6"
}

output "chicken_window_id" {
  value = eventstorecloud_project.chicken_window.id
}

output "chicken_window_net" {
	value = eventstorecloud_network.chicken_window
}

output "wings_cluster_dns_name" {
	value = eventstorecloud_managed_cluster.wings.dns_name
}
