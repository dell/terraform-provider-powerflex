package powerflex

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// SnapshotResourceSchema variable to define schema for the snapshot resource
var SnapshotResourceSchema schema.Schema = schema.Schema{
	Description:         "Manages snapshot resource.",
	MarkdownDescription: "Manages snapshot resource",
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description:         "The name of the snapshot.",
			Required:            true,
			MarkdownDescription: "The name of the snapshot.",
		},
		"volume_id": schema.StringAttribute{
			Description:         "The volume id for which snapshot is created.",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "The volume id for which snapshot is created",
		},
		"volume_name": schema.StringAttribute{
			Description:         "The volume name for which snapshot is created.",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "The volume name for which snapshot is created",
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(path.MatchRoot("volume_id")),
			},
		},
		"access_mode": schema.StringAttribute{
			Description:         "The Access mode of snapshot",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "The Access mode of snapshot",
			Validators: []validator.String{stringvalidator.OneOf(
				"ReadOnly",
				"ReadWrite",
			)},
			PlanModifiers: []planmodifier.String{
				stringDefault("ReadOnly"),
			},
		},
		"id": schema.StringAttribute{
			Description:         "The ID of the snapshot.",
			Computed:            true,
			MarkdownDescription: "The ID of the snapshot.",
		},
		"size": schema.Int64Attribute{
			Description:         "volume size",
			Optional:            true,
			MarkdownDescription: "volume size",
		},
		"capacity_unit": schema.StringAttribute{
			Description:         "capacity unit",
			Optional:            true,
			Computed:            true,
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
		"size_in_kb": schema.Int64Attribute{
			Description:         "Size in KB",
			Computed:            true,
			MarkdownDescription: "Size in KB",
		},
		"locked_auto_snapshot": schema.BoolAttribute{
			Description:         "locked auto snapshot",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "locake auto snapshot",
		},
		// "map_sdcs_id": schema.ListAttribute{
		// 	ElementType:         types.StringType,
		// 	Optional:            true,
		// 	MarkdownDescription: "map sdcs id",
		// },
		"sdc_list": schema.ListNestedAttribute{
			Description:         "mapped sdc info",
			Computed:            true,
			Optional:            true,
			MarkdownDescription: "mapped sdc info",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"sdc_id": schema.StringAttribute{
						Description:         "The ID of the SDC",
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The ID of the SDC",
					},
					"sdc_ip": schema.StringAttribute{
						Description:         "The IP of the SDC",
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The IP of the SDC",
						/* 						Validators: []validator.String{
							stringvalidator.AtLeastOneOf(path.Expressions{
								path.MatchRoot("sdc_id"),
								path.MatchRoot("sdc_name"),
							}...),
						}, */
					},
					"sdc_name": schema.StringAttribute{
						Description:         "The Name of the SDC",
						Computed:            true,
						Optional:            true,
						MarkdownDescription: "The Name of the SDC",
						/* 						Validators: []validator.String{
							stringvalidator.AtLeastOneOf(path.Expressions{
								path.MatchRoot("sdc_id"),
								path.MatchRoot("sdc_ip"),
							}...),
						}, */
					},
					"limit_iops": schema.Int64Attribute{
						Description:         "limit iops",
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "limit iops",
					},
					"limit_bw_in_mbps": schema.Int64Attribute{
						Description:         "limit bw in mbps",
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "limit bw in mbps",
					},
					"access_mode": schema.StringAttribute{
						Description:         "The Access Mode of the SDC",
						Computed:            true,
						Optional:            true,
						MarkdownDescription: "The Access Mode of the SDC",
					},
					"is_direct_buffer_mapping": schema.BoolAttribute{
						Description:         "is direct buffer mapping",
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "is direct buffer mapping",
					},
				},
			},
		},
	},
}
