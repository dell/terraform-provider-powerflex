package volume

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var VolumeResourceSchema schema.Schema = schema.Schema{
	Description: "Manages an volume.",
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description: "name",
			Required:    true,
		},
		"storage_pool_id": schema.StringAttribute{
			Description: "storage pool id",
			Required:    true,
		},
		"protection_domain_id": schema.StringAttribute{
			Description: "protection domain id",
			Required:    true,
		},
		"volume_size_in_kb": schema.StringAttribute{
			Description: "volume size in kb",
			Required:    true,
		},
		"volume_type": schema.StringAttribute{
			Description: "volume type",
			Optional:    true,
			Computed:    true,
		},
		"use_rm_cache": schema.BoolAttribute{
			Description: "use rm cache",
			Optional:    true,
			Computed:    true,
		},
		"id": schema.StringAttribute{
			Description: "ID",
			Computed:    true,
		},
		"creation_time": schema.Int64Attribute{
			Description: "Creation Time",
			Computed:    true,
		},
		"size_in_kb": schema.Int64Attribute{
			Description: "Size in KB",
			Computed:    true,
		},
		"ancestor_volume_id": schema.StringAttribute{
			Description: "ancestor volume id",
			Computed:    true,
		},
		"vtree_id": schema.StringAttribute{
			Description: "v tree id",
			Computed:    true,
		},
		"consistency_group_id": schema.StringAttribute{
			Description: "consistency group id",
			Computed:    true,
		},
		"data_layout": schema.StringAttribute{
			Description: "data layout",
			Computed:    true,
		},
		"not_genuine_snapshot": schema.BoolAttribute{
			Description: "not genuine snapshot",
			Computed:    true,
		},
		"access_mode_limit": schema.StringAttribute{
			Description: "access mode limit",
			Computed:    true,
		},
		"secure_snapshot_exp_time": schema.Int64Attribute{
			Description: "secure snapshot exp time",
			Computed:    true,
		},
		"managed_by": schema.StringAttribute{
			Description: "manged by",
			Computed:    true,
		},
		"locked_auto_snapshot": schema.BoolAttribute{
			Description: "locked auto snapshot",
			Computed:    true,
		},
		"locked_auto_snapshot_marked_for_removal": schema.BoolAttribute{
			Description: "locked auto snapshot marked for removal",
			Computed:    true,
		},
		"compression_method": schema.StringAttribute{
			Description: "compression method",
			Computed:    true,
		},
		"time_stamp_is_accurate": schema.BoolAttribute{
			Description: "time stamp is accurate",
			Computed:    true,
		},
		"original_expiry_time": schema.Int64Attribute{
			Description: "original expiry time",
			Computed:    true,
		},
		"volume_replication_state": schema.StringAttribute{
			Description: "volume replication state",
			Computed:    true,
		},
		"replication_journal_volume": schema.BoolAttribute{
			Description: "replication journal volume",
			Computed:    true,
		},
		"replication_time_stamp": schema.Int64Attribute{
			Description: "replication time stamp",
			Computed:    true,
		},
		"mapping_to_all_sdcs_enabled": schema.BoolAttribute{
			Description: "mapping to all sdcs enabled",
			Computed:    true,
		},
		"is_obfuscated": schema.BoolAttribute{
			Description: "is obfuscated",
			Computed:    true,
		},
		"mapped_scsi_initiator_info": schema.StringAttribute{
			Description: "mapped scsi initiator info",
			Computed:    true,
		},
		"links": schema.ListNestedAttribute{
			Description: "",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"rel": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"href": schema.StringAttribute{
						Description: "Numeric identifier of the coffee ingredient.",
						Computed:    true,
					},
				},
			},
		},
		"mapped_sdc_info": schema.ListNestedAttribute{
			Description: "",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"sdc_id": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"sdc_ip": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"limit_iops": schema.Int64Attribute{
						Description: "",
						Computed:    true,
					},
					"limit_bw_in_mbps": schema.Int64Attribute{
						Description: "",
						Computed:    true,
					},
					"sdc_name": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"access_mode": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"is_direct_buffer_mapping": schema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
				},
			},
		},
	},
}
