
resource "powerflex_snapshot" "snapshots-create" {
	name = "snapshots-create"
	volume_id = "<volume-id>"
}

resource "powerflex_snapshot" "snapshots-create-1" {
	name = "<snapshot-name>"
	volume_id = "<volume-id>"
	access_mode = "<access-mode>"
	map_sdcs_id = ["<sdc-id-1>","<sdc-id-2>"]
}