
resource "powerflex_snapshot" "snapshots-create" {
	name = "snapshots-create"
	volume_id = "<volume-id>"
}

resource "powerflex_snapshot" "snapshots-create-1" {
	name = "<snapshot name>"
	volume_name = "<volume name>"
	access_mode = "<access mode options are ReadOnly/ReadWrite, default value ReadOnly>"
	size = "<size[int] associated with capacity unit>"
	capacity_unit =  "<capacity unit options are gb/tb, default value gb>"
	lock_auto_snapshot = "<lock auto snapshot, snapshot which are created by snapshot policy can be locked.>"
	desired_retention = "<desired retention[int] associated with retention unit>"
	retention_unit = "<retention unit options are hours/days, default value hours>"
	remove_mode = "<remove mode options are ONLY_ME/INCLUDING_DESCENDANTS, default value ONLY_ME>"
	sdc_list = [
	   {
		sdc_id = "<sdc id either of sdc_id/sdc_name should be present to map snapshot to sdc>"
		limit_iops = "<limit iops setting on mapped sdc>"
		limit_bw_in_mbps = "<limit bw in mbps setting on mapped sdc>"
		access_mode = "<access mode options are ReadOnly/ReadWrite, default value ReadOnly>"
	   },
	   	   {
		sdc_name = "<sdc name either of sdc_id/sdc_name should be present to map snapshot to sdc>"
		limit_iops = "<limit iops setting on mapped sdc>"
		limit_bw_in_mbps = "<limit bw in mbps setting on mapped sdc>"
		access_mode = "<access mode options are ReadOnly/ReadWrite, default value ReadOnly>"
	   }
	]
}