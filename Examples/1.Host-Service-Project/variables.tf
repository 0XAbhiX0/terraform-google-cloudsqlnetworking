variable "host_project_id" {
  type        = string
  description = "Project Id of the Host GCP Project."
}

variable "service_project_id" {
  type        = string
  description = "Project Id of the Service GCP Project attached to the Host GCP project."
}

variable "cloudsql_instance_name" {
  type        = string
  description = "Name of the cloud sql instance which will be created."
}

variable "region" {
  type        = string
  description = "Name of a GCP region."
}

variable "zone" {
  type        = string
  description = "Name of a GCP zone, should be in the same region as specified in the region variable."
}

variable "database_version" {
  type        = string
  description = "Database version of the mysql in Cloud SQL ."
}

variable "network_id" {
  type        = string
  default     = ""
  description = "Complete network Id. This is required when var.create_network is set of false. e.g. : projects/pm-singleproject-20/global/networks/cloudsql-easy"
}

variable "subnetwork_id" {
  type        = string
  default     = ""
  description = "Complete subnetwork Id. This is required when var.create_subnetwork is set of false. e.g. : projects/pm-singleproject-20/regions/us-central1/subnetworks/cloudsql-easy-subnet"
}

variable "subnetwork_ip_cidr" {
  type        = string
  description = "CIDR range for the subnet to be created if var.create_subnetwork is set to true."

}

variable "gce_tags" {
  type        = list(string)
  default     = ["cloudsql"]
  description = "List of tags to be applied to gce instance."
}

variable "network_tier" {
  type        = string
  default     = "STANDARD"
  description = "Networking tier to be used."
}

variable "network_routing_mode" {
  type        = string
  default     = "GLOBAL"
  description = "Network Routing Mode to be used, Could be REGIONAL or GLOBAL."
}

variable "target_size" {
  type        = number
  default     = 1
  description = "Number of GCE instances to be created."
}

variable "network_name" {
  type        = string
  default     = "cloudsql-easy"
  description = "Name of the VPC network to be created if var.create_network is marked as true."
}

variable "subnetwork_name" {
  type        = string
  default     = "cloudsql-easy-subnet"
  description = "Name of the sub network to be created if var.create_subnetwork is marked as true."
}

variable "create_network" {
  type        = bool
  default     = true
  description = "Variable to determine if a new network should be created or not."
}

variable "create_subnetwork" {
  type        = bool
  default     = true
  description = "Variable to determine if a new sub network should be created or not."
}

