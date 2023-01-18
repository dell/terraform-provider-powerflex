package powerflex

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// VolumeResourceSchema variable to define schema for the volume resource
var VolumeResourceSchema schema.Schema = schema.Schema{
	Description:         "Manages volume resource.",
	MarkdownDescription: "Manages volume resource",
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description:         "The name of the volume.",
			Required:            true,
			MarkdownDescription: "The name of the volume.",
		},
		"storage_pool_id": schema.StringAttribute{
			Description:         "storage pool id",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "storage pool id",
		},
		"storage_pool_name": schema.StringAttribute{
			Description:         "Storage Pool Name",
			Optional:            true,
			MarkdownDescription: "Storage Pool Name",
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(path.MatchRoot("storage_pool_id")),
			},
		},
		"protection_domain_id": schema.StringAttribute{
			Description:         "Protection Domain ID.",
			MarkdownDescription: "Protection Domain ID.",
			Optional:            true,
		},
		"protection_domain_name": schema.StringAttribute{
			Description:         "Protection Domain Name.",
			MarkdownDescription: "Protection Domain Name.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(path.MatchRoot("protection_domain_id")),
			},
		},
		"size": schema.Int64Attribute{
			Description:         "volume size",
			Required:            true,
			MarkdownDescription: "volume size",
		},
		"capacity_unit": schema.StringAttribute{
			Description:         "capacity unit",
			Computed:            true,
			Optional:            true,
			MarkdownDescription: "capacity unit",
			Validators: []validator.String{stringvalidator.OneOf(
				"GB",
				"TB",
			)},
			PlanModifiers: []planmodifier.String{
				stringDefault("GB"),
			},
		},
		"volume_size_in_kb": schema.StringAttribute{
			Description:         "volume size in kb",
			Computed:            true,
			MarkdownDescription: "volume siz in kb",
		},
		"volume_type": schema.StringAttribute{
			Description:         "volume type",
			Computed:            true,
			MarkdownDescription: "volume type",
		},
		"use_rm_cache": schema.BoolAttribute{
			Description:         "use rm cache",
			Computed:            true,
			MarkdownDescription: "use rm cache",
		},
		"id": schema.StringAttribute{
			Description:         "The ID of the volume.",
			Computed:            true,
			MarkdownDescription: "The ID of the volume.",
		},
		"creation_time": schema.Int64Attribute{
			Description:         "Creation Time",
			Computed:            true,
			MarkdownDescription: "Creation Time",
		},
		"size_in_kb": schema.Int64Attribute{
			Description:         "Size in KB",
			Computed:            true,
			MarkdownDescription: "Size in KB",
		},
		"ancestor_volume_id": schema.StringAttribute{
			Description:         "ancestor volume id",
			Computed:            true,
			MarkdownDescription: "ancestor volume id",
		},
		"vtree_id": schema.StringAttribute{
			Description:         "vtree id",
			Computed:            true,
			MarkdownDescription: "vtree id",
		},
		"consistency_group_id": schema.StringAttribute{
			Description:         "consistency group id",
			Computed:            true,
			MarkdownDescription: "consistency group id",
		},
		"data_layout": schema.StringAttribute{
			Description:         "data layout",
			Computed:            true,
			MarkdownDescription: "data layout",
		},
		"not_genuine_snapshot": schema.BoolAttribute{
			Description:         "not genuine snapshot",
			Computed:            true,
			MarkdownDescription: "not genuine snapshot",
		},
		"access_mode_limit": schema.StringAttribute{
			Description:         "access mode limit",
			Computed:            true,
			MarkdownDescription: "access mode limit",
		},
		"secure_snapshot_exp_time": schema.Int64Attribute{
			Description:         "secure snapshot exp time",
			Computed:            true,
			MarkdownDescription: "secure snapshot exp time",
		},
		"managed_by": schema.StringAttribute{
			Description:         "mansged by",
			Computed:            true,
			MarkdownDescription: "managed by",
		},
		"locked_auto_snapshot": schema.BoolAttribute{
			Description:         "locked auto snapshot",
			Computed:            true,
			MarkdownDescription: "locake auto snapshot",
		},
		"locked_auto_snapshot_marked_for_removal": schema.BoolAttribute{
			Description:         "locked auto snapshot marked for removal",
			Computed:            true,
			MarkdownDescription: "locaked auto snapshot marked for removal",
		},
		"compression_method": schema.StringAttribute{
			Description:         "compression method",
			Computed:            true,
			MarkdownDescription: "compression method",
		},
		"time_stamp_is_accurate": schema.BoolAttribute{
			Description:         "time stamp is accurate",
			Computed:            true,
			MarkdownDescription: "time stamp is accurate",
		},
		"original_expiry_time": schema.Int64Attribute{
			Description:         "original expiry time",
			Computed:            true,
			MarkdownDescription: "original expriry time",
		},
		"volume_replication_state": schema.StringAttribute{
			Description:         "volume replication state",
			Computed:            true,
			MarkdownDescription: "volume replication state",
		},
		"replication_journal_volume": schema.BoolAttribute{
			Description:         "replication journal volume",
			Computed:            true,
			MarkdownDescription: "replication journal volume",
		},
		"replication_time_stamp": schema.Int64Attribute{
			Description:         "replication time stamp",
			Computed:            true,
			MarkdownDescription: "replication time stamp",
		},
		"mapping_to_all_sdcs_enabled": schema.BoolAttribute{
			Description:         "mapping to all sdcs enabled",
			Computed:            true,
			MarkdownDescription: "mapping to all sdcs enabled",
		},
		"is_obfuscated": schema.BoolAttribute{
			Description:         "is obfuscated",
			Computed:            true,
			MarkdownDescription: "is obfuscated",
		},
		"mapped_scsi_initiator_info": schema.StringAttribute{
			Description:         "mapped scsi initiator info",
			Computed:            true,
			MarkdownDescription: "mapped scsi initiator info",
		},
		"map_sdcs_id": schema.ListAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			MarkdownDescription: "map sdcs id",
		},
		"links": schema.ListNestedAttribute{
			Description:         "links for the volume resource",
			Computed:            true,
			MarkdownDescription: "links for the volume resource",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"rel": schema.StringAttribute{
						Description:         "rel",
						Computed:            true,
						MarkdownDescription: "rel",
					},
					"href": schema.StringAttribute{
						Description:         "href",
						Computed:            true,
						MarkdownDescription: "href",
					},
				},
			},
		},
		"mapped_sdc_info": schema.ListNestedAttribute{
			Description:         "mapped sdc info",
			Computed:            true,
			MarkdownDescription: "mapped sdc info",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"sdc_id": schema.StringAttribute{
						Description:         "The ID of the SDC",
						Computed:            true,
						MarkdownDescription: "The ID of the SDC",
					},
					"sdc_ip": schema.StringAttribute{
						Description:         "The IP of the SDC",
						Computed:            true,
						MarkdownDescription: "The IP of the SDC",
					},
					"limit_iops": schema.Int64Attribute{
						Description:         "limit iops",
						Computed:            true,
						MarkdownDescription: "limit iops",
					},
					"limit_bw_in_mbps": schema.Int64Attribute{
						Description:         "limit bw in mbps",
						Computed:            true,
						MarkdownDescription: "limit bw in mbps",
					},
					"sdc_name": schema.StringAttribute{
						Description:         "The Name of the SDC",
						Computed:            true,
						MarkdownDescription: "The Name of the SDC",
					},
					"access_mode": schema.StringAttribute{
						Description:         "The Access Mode of the SDC",
						Computed:            true,
						MarkdownDescription: "The Access Mode of the SDC",
					},
					"is_direct_buffer_mapping": schema.BoolAttribute{
						Description:         "is direct buffer mapping",
						Computed:            true,
						MarkdownDescription: "is direct buffer mapping",
					},
				},
			},
		},
	},
}
