package powerflex

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// SnapshotMarkdownDescription add notes for resource
const SnapshotMarkdownDescription = `Manages Snapshot in powerflex.
Note: Snapshot creation or update is not atomic. In case of partially completed operations, terraform can mark the resource as tainted.
One can manually remove the taint and try applying the configuration (after making necessary adjustments).
Warning: If the taint is not removed, terraform will destroy and recreate the resource.
`

// SnapshotResourceSchema variable to define schema for the snapshot resource
var SnapshotResourceSchema schema.Schema = schema.Schema{
	Description:         "Manages snapshot resource.",
	MarkdownDescription: SnapshotMarkdownDescription,
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description:         "The name of the snapshot.",
			Required:            true,
			MarkdownDescription: "The name of the snapshot.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"volume_id": schema.StringAttribute{
			Description:         "The volume id for which snapshot is created.",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "The volume id for which snapshot is created - Either of Volume ID/Name is Required.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.ExactlyOneOf(path.MatchRoot("volume_name")),
			},
		},
		"volume_name": schema.StringAttribute{
			Description:         "The volume name for which snapshot is created.",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "The volume name for which snapshot is created - Either of Volume ID/Name is Required.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
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
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"size": schema.Int64Attribute{
			Description:         "snapshot size",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "snapshot size",
		},
		"capacity_unit": schema.StringAttribute{
			Description:         "capacity unit",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "capacity unit",
			Validators: []validator.String{stringvalidator.OneOf(
				"GB",
				"TB",
			),
				stringvalidator.AlsoRequires(path.MatchRoot("size")),
			},
			PlanModifiers: []planmodifier.String{
				stringDefault("GB"),
			},
		},
		"size_in_kb": schema.Int64Attribute{
			Description:         "Size in KB",
			Computed:            true,
			MarkdownDescription: "Size in KB",
		},
		"lock_auto_snapshot": schema.BoolAttribute{
			Description:         "lock auto snapshot",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "lock auto snapshot",
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"desired_retention": schema.Int64Attribute{
			Description:         "desired retention of snapshot",
			Optional:            true,
			MarkdownDescription: "desired retention of snapshot",
		},
		"retention_unit": schema.StringAttribute{
			Description:         "retention unit of snapshot",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "retention unit of snapshot",
			Validators: []validator.String{stringvalidator.OneOf(
				"hours",
				"days",
			),
				stringvalidator.AlsoRequires(path.MatchRoot("desired_retention")),
			},
			PlanModifiers: []planmodifier.String{
				stringDefault("hours"),
			},
		},
		"retention_in_min": schema.StringAttribute{
			Description:         "retention of snapshot in min",
			Computed:            true,
			MarkdownDescription: "retention of snapshot in min",
		},
		"remove_mode": schema.StringAttribute{
			Description:         "remove mode of snapshot",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "remove mode of snapshot",
			Validators: []validator.String{stringvalidator.OneOf(
				"ONLY_ME",
				"INCLUDING_DESCENDANTS",
			)},
			PlanModifiers: []planmodifier.String{
				stringDefault("ONLY_ME"),
			},
		},
		"sdc_list": sdcListSchema,
	},
}
