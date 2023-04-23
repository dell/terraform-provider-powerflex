# Command to run this tf file : terraform init && terraform plan && terraform apply.
# Create, Update, Delete is supported for this resource.
# To import, check sdc_volumes_mapping_resource_import.tf for more info.
# To create/update, either SDC ID or SDC name must be provided.
# volume_list attribute is optional. 
# To check which attributes of the sdc_volumes_mappping resource can be updated, please refer Product Guide in the documentation

resource "powerflex_sdc_volumes_mapping" "mapping-test" {
  id = "e3ce1fb600000001"
  volume_list = [
    {
      volume_id        = "edb2059700000002"
      limit_iops       = 140
      limit_bw_in_mbps = 19
      access_mode      = "ReadOnly"
    },
    {
      volume_name      = "terraform-vol"
      access_mode      = "ReadWrite"
      limit_iops       = 120
      limit_bw_in_mbps = 25
    }
  ]
}
