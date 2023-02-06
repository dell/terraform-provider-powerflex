# Command to run this tf file : terraform init && terraform plan && terraform apply
# Create, Update, Delete is supported for this resource
# To import , check volume_resource_import.tf for more info
# To create / update, either storage_pool_id or storage_pool_name must be provided
# Also , to create / update, either protection_domain_id or protection_domain_name must be provided
# name, size is the required parameter to create or update
# other  atrributes like : capacity_unit, volume_type, use_rm_cache, compression_method, access_mode, remove_mode, sdc_list are optional 
# To check which attributes of the snapshot can be updated, please refer Product Guide in the documentation


resource "powerflex_volume" "avengers-volume-create"{
	name = "avengers-volume-create"
	protection_domain_name = "domain1"
	storage_pool_name = "pool1" #pool1 have medium granularity
	size = 8
	use_rm_cache = true 
	volume_type = "ThickProvisioned" 
	access_mode = "ReadWrite"
	sdc_list = [
	  {
			   sdc_name = "sdc_01"
			   limit_iops = 119
			   limit_bw_in_mbps = 19
			   access_mode = "ReadOnly"
		   },
	]
  }


# General guidlines for furnishing this resource block  
# resource "powerflex_volume" "avengers-volume-create"{
# 	name = "<volume-name>"
# 	protection_domain_name = "<protection-domain-name>"
# 	storage_pool_name = "<storage-pool-name>"
# 	size = "<size in int>"
# 	capacity_unit = "<GB/TB capacity unit>"
# 	use_rm_cache = "true/false for use rm cache" 
# 	volume_type = "<ThickProvisioned/ThinProvisioned volume type>" 
# 	access_mode = "<ReadWrite/ReadOnly volume access mode>"
# 	compression_method = "<None/Normal compression method>"
# 	sdc_list = [
# 	  		{
# 			   sdc_name = "<sdc name>"
# 			   limit_iops = "<iops limit in int>"
# 			   limit_bw_in_mbps = "<bandwidth limit in mbps>"
# 			   access_mode = "<ReadWrite/ReadOnly/Noaccess sdc access mode>"
# 		   },
# 	]
# }