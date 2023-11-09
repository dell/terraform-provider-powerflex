# resource "powerflex_protection_domain" "pd" {
#   name = "domain_a2"
# }

# resource "powerflex_sds" "sds1" {
#   name                 = "sds_a1"
#   protection_domain_id = powerflex_protection_domain.pd.id
#   ip_list = [
#     {
#       ip   = "10."
#       role = "all"
#     }
#   ]
#   depends_on = [powerflex_protection_domain.pd]
# }

# resource "powerflex_sds" "sds2" {
#   name                 = "sds_a2"
#   protection_domain_id = powerflex_protection_domain.pd.id
#   ip_list = [
#     {
#       ip   = "10."
#       role = "all"
#     }
#   ]
#   depends_on = [powerflex_protection_domain.pd]
# }

# resource "powerflex_sds" "sds3" {
#   name                 = "sds_a3"
#   protection_domain_id = powerflex_protection_domain.pd.id
#   ip_list = [
#     {
#       ip   = "10."
#       role = "all"
#     }
#   ]
#   depends_on = [powerflex_protection_domain.pd]
# }

# resource "powerflex_storage_pool" "sp" {
#   name                 = "SP"
#   protection_domain_name = powerflex_protection_domain.pd.name
#   media_type           = "SSD"
#   use_rmcache          = true
#   use_rfcache          = true
# }

# resource "powerflex_device" "device1" {
#   name                       = "device_a1"
#   device_path                = "/dev/sdb"
#   sds_id                     = powerflex_sds.sds1.id
#   storage_pool_id            = powerflex_storage_pool.sp.id
#   media_type                 = "HDD"
#   external_acceleration_type = "ReadAndWrite"
#   depends_on                 = [powerflex_storage_pool.sp]
# }

# resource "powerflex_device" "device2" {
#   name                       = "device_a2"
#   device_path                = "/dev/sdb"
#   sds_id                     = powerflex_sds.sds2.id
#   storage_pool_id            = powerflex_storage_pool.sp.id
#   media_type                 = "HDD"
#   external_acceleration_type = "ReadAndWrite"
#   depends_on                 = [powerflex_storage_pool.sp]
# }

# resource "powerflex_device" "device3" {
#   name                       = "device_a3"
#   device_path                = "/dev/sdb"
#   sds_id                     = powerflex_sds.sds3.id
#   storage_pool_id            = powerflex_storage_pool.sp.id
#   media_type                 = "HDD"
#   external_acceleration_type = "ReadAndWrite"
#   depends_on                 = [powerflex_storage_pool.sp]
# }

resource "powerflex_volume" "volume" {
  name                 = "volume2"
  protection_domain_name = "domain1"
  storage_pool_name      = "pool1"
  size                 = 8
  capacity_unit        = "GB"
  volume_type          = "ThinProvisioned"
  #depends_on           = [powerflex_storage_pool.sp, powerflex_protection_domain.pd,powerflex_device.device1,powerflex_device.device2,powerflex_device.device3]
}

resource "powerflex_snapshot" "snap" {
  name                 = "snap1"
  volume_name     = powerflex_volume.volume.name
  access_mode   = "ReadWrite"
  size          = 8
  capacity_unit = "GB"
}

