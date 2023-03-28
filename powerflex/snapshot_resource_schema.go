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
const SnapshotMarkdownDescription = `This resource can be used to manage snapshots of volumes on a PowerFlex array.
Note: Snapshot creation or update is not atomic. In case of partially completed operations, terraform can mark the resource as tainted.
One can manually remove the taint and try applying the configuration (after making necessary adjustments).
Warning: If the taint is not removed, terraform will destroy and recreate the resource.
`

// SnapshotResourceSchema variable to define schema for the snapshot resource
var SnapshotResourceSchema schema.Schema = schema.Schema{
	Description:         "This resource can be used to manage snapshots of volumes on a PowerFlex array.",
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
			Description: "The ID of the volume from which snapshot is to be created." +
				" Must be provided if and only if 'volume_name' is not provided." +
				" Cannot be updated.",
			Optional: true,
			Computed: true,
			MarkdownDescription: "The ID of the volume from which snapshot is to be created." +
				" Must be provided if and only if `volume_name` is not provided." +
				" Cannot be updated.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.ExactlyOneOf(path.MatchRoot("volume_name")),
			},
		},
		"volume_name": schema.StringAttribute{
			Description: "The volume name for which snapshot is created." +
				" Must be provided if and only if 'volume_id' is not provided." +
				" Cannot be updated.",
			Optional: true,
			Computed: true,
			MarkdownDescription: "The volume name for which snapshot is created." +
				" Must be provided if and only if `volume_id` is not provided." +
				" Cannot be updated.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.ExactlyOneOf(path.MatchRoot("volume_id")),
			},
		},
		"access_mode": schema.StringAttribute{
			Description:         "The Access mode of snapshot. Valid values are 'ReadOnly' and 'ReadWrite'. Default value is 'ReadOnly'.",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "The Access mode of snapshot. Valid values are `ReadOnly` and `ReadWrite`. Default value is `ReadOnly`.",
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
			Description:         "Size of the snapshot. The unit of size is defined by 'capacity_unit'. The storage capacity of a snapshot cannot be decreased.",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Size of the snapshot. The unit of size is defined by `capacity_unit`. The storage capacity of a snapshot cannot be decreased.",
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
			Description: "The minimum amount of time that the snapshot should be retained on the array starting at the time of apply." +
				" The unit is defined by 'retention_unit'." +
				" Cannot be decreased.",
			Optional: true,
			MarkdownDescription: "The minimum amount of time that the snapshot should be retained on the array starting at the time of apply." +
				" The unit is defined by `retention_unit`." +
				" Cannot be decreased.",
		},
		"retention_unit": schema.StringAttribute{
			Description:         "Retention unit of the snapshot. Valid values are 'hours' and 'days'. Default value is 'hours'.",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Retention unit of the snapshot. Valid values are 'hours' and 'days'. Default value is 'hours'.",
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
			Description:         "Remove mode of the snapshot. Valid values are 'ONLY_ME' and 'INCLUDING_DESCENDANTS'. Default value is 'ONLY_ME'.",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Remove mode of the snapshot. Valid values are `ONLY_ME` and `INCLUDING_DESCENDANTS`. Default value is `ONLY_ME`.",
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
