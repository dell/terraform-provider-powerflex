resource "powerflex_volume" "avengers-volume-create"{
	name = "<volume-name>"
	protection_domain_name = "<protection-domain-name>"
	storage_pool_name = "<storage-pool-name>"
	size = "<size in int>"
	capacity_unit = "<GB/TB capacity unit>"
	use_rm_cache = "true/false for use rm cache" 
	volume_type = "<ThickProvisioned/ThinProvisioned volume type>" 
	access_mode = "<ReadWrite/ReadOnly access mode>"
	compression_method = "<None/Normal compression method>"
	sdc_list = [
	  		{
			   sdc_name = "<sdc name>"
			   limit_iops = "<iops limit in int>"
			   limit_bw_in_mbps = "<bandwidth limit in mbps>"
			   access_mode = "<ReadWrite/ReadOnly/Noaccess mode>"
		   },
	]
}