package volumedatasource

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// VolumeDataSourceSchema is the schema for reading the volume data
var VolumeDataSourceSchema schema.Schema = schema.Schema{
	Description: ".",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "",
			Optional:    true,
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "",
			Optional:    true,
			Computed:    true,
		},
		"storage_pool_id": schema.StringAttribute{
			Description: "",
			Optional:    true,
			Computed:    true,
		},
		"volumes": schema.ListNestedAttribute{
			Description: "List of volumes",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"name": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"creation_time": schema.Int64Attribute{
						Description: "",
						Computed:    true,
					},
					"size_in_kb": schema.Int64Attribute{
						Description: "",
						Computed:    true,
					},
					"ancestor_volume_id": schema.StringAttribute{
						Description: ".",
						Computed:    true,
					},
					"vtree_id": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"consistency_group_id": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"volume_type": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"use_rm_cache": schema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"storage_pool_id": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"data_layout": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"not_genuine_snapshot": schema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"access_mode_limit": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"secure_snapshot_exp_time": schema.Int64Attribute{
						Description: "",
						Computed:    true,
					},
					"managed_by": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"locked_auto_snapshot": schema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"locked_auto_snapshot_marked_for_removal": schema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"compression_method": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"time_stamp_is_accurate": schema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"original_expiry_time": schema.Int64Attribute{
						Description: "",
						Computed:    true,
					},
					"volume_replication_state": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"replication_journal_volume": schema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"replication_time_stamp": schema.Int64Attribute{
						Description: "",
						Computed:    true,
					},
					"links": schema.ListNestedAttribute{
						Description: "",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"rel": schema.StringAttribute{
									Description: "Numeric identifier of the coffee ingredient.",
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
			},
		},
	},
}
