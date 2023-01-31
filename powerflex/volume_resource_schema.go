package powerflex

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// VolumeMarkdownDescription add notes for resource
const VolumeMarkdownDescription = `Manages Volume in powerflex.
Note: Volume creation or update is not atomic. In case of partially completed operations, terraform can mark the resource as tainted.
One can manually remove the taint and try applying the configuration (after making necessary adjustments).
Warning: If the taint is not removed, terraform will destroy and recreate the resource.
`

// VolumeResourceSchema variable to define schema for the volume resource
var VolumeResourceSchema schema.Schema = schema.Schema{
	Description:         "Manages volume resource.",
	MarkdownDescription: VolumeMarkdownDescription,
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description:         "The name of the volume.",
			Required:            true,
			MarkdownDescription: "The name of the volume.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"storage_pool_id": schema.StringAttribute{
			Description:         "storage pool id",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "storage pool id",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.ExactlyOneOf(path.MatchRoot("storage_pool_name")),
			},
		},
		"storage_pool_name": schema.StringAttribute{
			Description:         "Storage Pool Name",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Storage Pool Name",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.ExactlyOneOf(path.MatchRoot("storage_pool_id")),
			},
		},
		"protection_domain_id": schema.StringAttribute{
			Description:         "Protection Domain ID.",
			MarkdownDescription: "Protection Domain ID.",
			Computed:            true,
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.ExactlyOneOf(path.MatchRoot("protection_domain_name")),
			},
		},
		"protection_domain_name": schema.StringAttribute{
			Description:         "Protection Domain Name.",
			MarkdownDescription: "Protection Domain Name.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
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
		"volume_type": schema.StringAttribute{
			Description:         "volume type",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "volume type",
			Validators: []validator.String{stringvalidator.OneOf(
				"ThickProvisioned",
				"ThinProvisioned",
			)},
			PlanModifiers: []planmodifier.String{
				stringDefault("ThinProvisioned"),
			},
		},
		"use_rm_cache": schema.BoolAttribute{
			Description:         "use rm cache",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "use rm cache",
		},
		"compression_method": schema.StringAttribute{
			Description: "Compression Method on the volume.",
			// Optional:            true,
			Computed:            true,
			MarkdownDescription: "Compression Method the volume.",
			Validators: []validator.String{stringvalidator.OneOf(
				"None",
				"Normal",
			)},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"id": schema.StringAttribute{
			Description:         "The ID of the volume.",
			Computed:            true,
			MarkdownDescription: "The ID of the volume.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"size_in_kb": schema.Int64Attribute{
			Description:         "Size in KB",
			Computed:            true,
			MarkdownDescription: "Size in KB",
		},
		"access_mode": schema.StringAttribute{
			Description:         "The Access mode of Volume",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "The Access mode of Volume",
			Validators: []validator.String{stringvalidator.OneOf(
				"ReadOnly",
				"ReadWrite",
			)},
			PlanModifiers: []planmodifier.String{
				stringDefault("ReadOnly"),
			},
		},
		"remove_mode": schema.StringAttribute{
			Description:         "remove mode of Volume",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "remove mode of Volume",
			Validators: []validator.String{stringvalidator.OneOf(
				"ONLY_ME",
				"INCLUDING_DESCENDANTS",
			)},
			PlanModifiers: []planmodifier.String{
				stringDefault("ONLY_ME"),
			},
		},
		"sdc_list": schema.SetNestedAttribute{
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
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("sdc_name")),
						},
					},
					"sdc_name": schema.StringAttribute{
						Description:         "The Name of the SDC",
						Computed:            true,
						Optional:            true,
						MarkdownDescription: "The Name of the SDC",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("sdc_id")),
						},
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
						Validators: []validator.String{stringvalidator.OneOf(
							"ReadOnly",
							"ReadWrite",
							"NoAccess",
						)},
						PlanModifiers: []planmodifier.String{
							stringDefault("ReadOnly"),
						},
					},
				},
			},
		},
	},
}
